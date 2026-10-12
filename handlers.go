package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var rootDir string
var liteMode bool

var downloadCounts = map[string]int{}
var dlCountsMu sync.Mutex

const dlCountsFile = ".dlcounts.json"

func loadDownloadCounts() {
	data, err := os.ReadFile(dlCountsFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Failed to load download counts: %v", err)
		}
		return
	}
	dlCountsMu.Lock()
	defer dlCountsMu.Unlock()
	if err := json.Unmarshal(data, &downloadCounts); err != nil {
		log.Printf("Failed to parse %s: %v", dlCountsFile, err)
	}
}

func saveDownloadCounts() {
	dlCountsMu.Lock()
	defer dlCountsMu.Unlock()
	data, err := json.MarshalIndent(downloadCounts, "", "  ")
	if err != nil {
		log.Printf("Failed to save download counts: %v", err)
		return
	}
	if err := os.WriteFile(dlCountsFile, data, 0644); err != nil {
		log.Printf("Failed to save download counts: %v", err)
	}
}

func isLiteRequest(r *http.Request) bool {
	return liteMode || r.URL.Query().Get("lite") == "1"
}

func renderDirListing(w http.ResponseWriter, r *http.Request, fullPath, cleanURL string) {
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	files, stats := buildFileList(entries, cleanURL)
	data := PageData{
		CurrentPath: cleanURL,
		CurrentURL:  escapeURL(cleanURL),
		Breadcrumbs: buildBreadcrumbs(cleanURL),
		Files:       files,
		CSS:         template.CSS(cssContent),
		Stats:       stats,
		Lite:        isLiteRequest(r),
		Version:     Version,
		IsRoot:      cleanURL == "/" || cleanURL == ".",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := tmplLite
	if !data.Lite {
		t = tmplFull
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// escapeURL percent-encodes a slash-separated path for use in links and redirects,
// so names containing '#', '?' or '%' keep working.
func escapeURL(p string) string {
	return (&url.URL{Path: p}).EscapedPath()
}

// parentURL returns the slash-separated parent of relPath, "/" for top-level entries.
func parentURL(relPath string) string {
	dir := path.Dir(relPath)
	if dir == "." {
		return "/"
	}
	return dir
}

// resolvePath maps a slash-separated URL path to a filesystem path inside rootDir.
func resolvePath(relPath string) (string, bool) {
	full := filepath.Join(rootDir, filepath.FromSlash(path.Clean("/"+relPath)))
	if full == rootDir || strings.HasPrefix(full, strings.TrimSuffix(rootDir, string(filepath.Separator))+string(filepath.Separator)) {
		return full, true
	}
	return "", false
}

// validName reports whether name is usable as a single path component.
func validName(name string) bool {
	return name != "" && name != "." && name != ".." && !strings.ContainsAny(name, "/\\")
}

// statError writes the response for a failed os.Stat.
func statError(w http.ResponseWriter, err error) {
	if os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

var actionPrefixes = []string{"edit", "delete", "rename", "zip", "mkdir", "mkfile"}

// splitAction splits "/<action>/<rel>" into its action and "/<rel>".
func splitAction(cleanURL, method string) (action, relPath string) {
	for _, a := range actionPrefixes {
		prefix := "/" + a
		// path.Clean strips the trailing slash of "/mkdir/", which the root
		// directory's forms post to. Only creation is allowed on the root itself.
		if cleanURL == prefix && method == http.MethodPost && (a == "mkdir" || a == "mkfile") {
			return a, "/"
		}
		if strings.HasPrefix(cleanURL, prefix+"/") {
			return a, strings.TrimPrefix(cleanURL, prefix)
		}
	}
	return "", cleanURL
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	cleanURL := path.Clean("/" + r.URL.Path)
	action, relPath := splitAction(cleanURL, r.Method)
	fullPath, ok := resolvePath(relPath)
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	switch action {
	case "edit":
		if r.Method == http.MethodPost {
			handleSaveEdit(w, r, fullPath, relPath)
		} else {
			handleEdit(w, r, fullPath, relPath)
		}

	case "delete":
		if r.Method == http.MethodPost {
			handleDelete(w, r, fullPath, relPath)
		} else {
			handleDeleteConfirm(w, r, fullPath, relPath)
		}

	case "rename":
		if r.Method == http.MethodPost {
			handleRename(w, r, fullPath, relPath)
		} else {
			handleRenameForm(w, r, fullPath, relPath)
		}

	case "zip":
		handleZipDownload(w, r, fullPath)

	case "mkdir":
		handleCreateDir(w, r, fullPath, relPath)

	case "mkfile":
		handleCreateFile(w, r, fullPath, relPath)

	default:
		stat, err := os.Stat(fullPath)
		if err != nil {
			statError(w, err)
			return
		}

		if !stat.IsDir() {
			if r.Method == http.MethodGet {
				dlCountsMu.Lock()
				downloadCounts[cleanURL]++
				dlCountsMu.Unlock()
			}
			http.ServeFile(w, r, fullPath)
			return
		}

		if r.Method == http.MethodPost {
			handleUpload(w, r, fullPath, cleanURL)
			return
		}

		renderDirListing(w, r, fullPath, cleanURL)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request, targetDir string, redirectURL string) {
	// Parts beyond 32 MiB are spooled to temporary files instead of memory.
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "upload too large or invalid", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		files = r.MultipartForm.File["file"]
	}
	if len(files) == 0 {
		http.Error(w, "no files provided", http.StatusBadRequest)
		return
	}

	for _, header := range files {
		name := filepath.Base(header.Filename)
		if !validName(name) {
			http.Error(w, fmt.Sprintf("invalid file name %q", header.Filename), http.StatusBadRequest)
			return
		}
		if err := saveUploadedFile(header, filepath.Join(targetDir, name)); err != nil {
			http.Error(w, fmt.Sprintf("failed to save %q: %v", name, err), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, escapeURL(redirectURL), http.StatusSeeOther)
}

func saveUploadedFile(header *multipart.FileHeader, destPath string) error {
	src, err := header.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, src)
	if cerr := dst.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(destPath)
	}
	return err
}

func handleCreateDir(w http.ResponseWriter, r *http.Request, targetDir, redirectURL string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	dirName := filepath.Base(r.FormValue("dirname"))
	if !validName(dirName) {
		http.Error(w, "valid dirname required", http.StatusBadRequest)
		return
	}
	newDir := filepath.Join(targetDir, dirName)
	if err := os.Mkdir(newDir, 0755); err != nil {
		http.Error(w, "failed to create directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, escapeURL(redirectURL), http.StatusSeeOther)
}

func handleCreateFile(w http.ResponseWriter, r *http.Request, targetDir, redirectURL string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fileName := filepath.Base(r.FormValue("filename"))
	if !validName(fileName) {
		http.Error(w, "valid filename required", http.StatusBadRequest)
		return
	}
	newFile := filepath.Join(targetDir, fileName)
	f, err := os.OpenFile(newFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	f.Close()
	http.Redirect(w, r, escapeURL(redirectURL), http.StatusSeeOther)
}

func handleEdit(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if !isTextFile(fullPath) {
		http.Error(w, "Cannot edit binary files", http.StatusForbidden)
		return
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		statError(w, err)
		return
	}
	cleanURL := parentURL(relPath)
	pageData := EditPageData{
		CurrentPath:  cleanURL,
		Breadcrumbs:  buildBreadcrumbs(cleanURL),
		CSS:          template.CSS(cssContent),
		FileName:     filepath.Base(relPath),
		FilePath:     escapeURL(cleanURL),
		Content:      string(data),
		SaveURL:      escapeURL("/edit" + relPath),
		HighlightJS:  template.JS(hlJS),
		HighlightCSS: template.CSS(hlCSS),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := tmplEditorLite
	if !isLiteRequest(r) {
		t = tmplEditor
	}
	if err := t.Execute(w, pageData); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// maxEditSize bounds the editor's form body; without an explicit limit
// ParseForm stops at 10 MB and FormValue would return an empty string.
const maxEditSize = 64 << 20

func handleSaveEdit(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isTextFile(fullPath) {
		http.Error(w, "Cannot edit binary files", http.StatusForbidden)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxEditSize)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to read form: "+err.Error(), http.StatusBadRequest)
		return
	}
	if _, ok := r.PostForm["content"]; !ok {
		http.Error(w, "content required", http.StatusBadRequest)
		return
	}
	content := r.PostForm.Get("content")
	// Browsers submit textarea line breaks as CRLF; keep LF files LF.
	if old, err := os.ReadFile(fullPath); err == nil && !strings.Contains(string(old), "\r\n") {
		content = strings.ReplaceAll(content, "\r\n", "\n")
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, escapeURL(parentURL(relPath)), http.StatusSeeOther)
}

func handleDeleteConfirm(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	stat, err := os.Stat(fullPath)
	if err != nil {
		statError(w, err)
		return
	}
	name := filepath.Base(relPath)
	cleanURL := parentURL(relPath)
	msg := fmt.Sprintf("Are you sure you want to delete %q?", name)
	if stat.IsDir() {
		msg = fmt.Sprintf("Are you sure you want to delete directory %q and ALL its contents?", name)
	}
	pageData := ConfirmPageData{
		CurrentPath: cleanURL,
		Breadcrumbs: buildBreadcrumbs(cleanURL),
		CSS:         template.CSS(cssContent),
		Title:       "Confirm Delete",
		Message:     msg,
		ActionURL:   escapeURL("/delete" + relPath),
		CancelURL:   escapeURL(cleanURL),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplConfirm.Execute(w, pageData); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if fullPath == rootDir {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	stat, err := os.Lstat(fullPath)
	if err != nil {
		statError(w, err)
		return
	}
	if stat.IsDir() {
		err = os.RemoveAll(fullPath)
	} else {
		err = os.Remove(fullPath)
	}
	if err != nil {
		http.Error(w, "Failed to delete: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, escapeURL(parentURL(relPath)), http.StatusSeeOther)
}

func handleRenameForm(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if _, err := os.Stat(fullPath); err != nil {
		statError(w, err)
		return
	}
	cleanURL := parentURL(relPath)
	pageData := RenamePageData{
		CurrentPath: cleanURL,
		Breadcrumbs: buildBreadcrumbs(cleanURL),
		CSS:         template.CSS(cssContent),
		OldName:     filepath.Base(relPath),
		NewURL:      escapeURL("/rename" + relPath),
		CancelURL:   escapeURL(cleanURL),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplRename.Execute(w, pageData); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func handleRename(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if fullPath == rootDir {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	newName := filepath.Base(r.FormValue("newname"))
	if !validName(newName) {
		http.Error(w, "valid newname required", http.StatusBadRequest)
		return
	}
	newPath := filepath.Join(filepath.Dir(fullPath), newName)

	if _, err := os.Lstat(fullPath); err != nil {
		statError(w, err)
		return
	}
	if newPath != fullPath {
		if _, err := os.Lstat(newPath); err == nil {
			http.Error(w, fmt.Sprintf("%q already exists", newName), http.StatusConflict)
			return
		}
	}

	if err := os.Rename(fullPath, newPath); err != nil {
		http.Error(w, "Failed to rename: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, escapeURL(parentURL(relPath)), http.StatusSeeOther)
}

func handleZipDownload(w http.ResponseWriter, r *http.Request, fullPath string) {
	stat, err := os.Stat(fullPath)
	if err != nil {
		statError(w, err)
		return
	}
	if !stat.IsDir() {
		http.ServeFile(w, r, fullPath)
		return
	}

	name := filepath.Base(fullPath) + ".zip"
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))

	zw := zip.NewWriter(w)
	err = filepath.Walk(fullPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(fullPath, p)
		if err != nil || relPath == "." {
			return err
		}
		if !info.IsDir() && !info.Mode().IsRegular() {
			return nil
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(relPath)
		if info.IsDir() {
			hdr.Name += "/"
			_, err := zw.CreateHeader(hdr)
			return err
		}
		hdr.Method = zip.Deflate
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		fw, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, f)
		return err
	})
	if err != nil {
		// Headers are already sent; abort the connection so the client sees
		// a failed download rather than a silently truncated archive.
		log.Printf("Zip of %s failed: %v", fullPath, err)
		panic(http.ErrAbortHandler)
	}
	if err := zw.Close(); err != nil {
		log.Printf("Zip of %s failed: %v", fullPath, err)
	}
}
