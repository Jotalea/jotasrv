package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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
		return
	}
	dlCountsMu.Lock()
	defer dlCountsMu.Unlock()
	json.Unmarshal(data, &downloadCounts)
}

func saveDownloadCounts() {
	dlCountsMu.Lock()
	defer dlCountsMu.Unlock()
	data, err := json.MarshalIndent(downloadCounts, "", "  ")
	if err != nil {
		log.Printf("Failed to save download counts: %v", err)
		return
	}
	os.WriteFile(dlCountsFile, data, 0644)
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

func handleRequest(w http.ResponseWriter, r *http.Request) {
	cleanURL := filepath.Clean(r.URL.Path)
	fullPath := filepath.Join(rootDir, filepath.FromSlash(cleanURL))

	if !strings.HasPrefix(fullPath, rootDir) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	switch {
	case strings.HasPrefix(cleanURL, "/edit/"):
		relPath := "/" + strings.TrimPrefix(cleanURL, "/edit/")
		editFull := filepath.Join(rootDir, filepath.FromSlash(relPath))
		editFull = filepath.Clean(editFull)
		if !strings.HasPrefix(editFull, rootDir) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Method == http.MethodPost {
			handleSaveEdit(w, r, editFull, relPath)
		} else {
			handleEdit(w, r, editFull, relPath)
		}

	case strings.HasPrefix(cleanURL, "/delete/"):
		relPath := "/" + strings.TrimPrefix(cleanURL, "/delete/")
		deleteFull := filepath.Join(rootDir, filepath.FromSlash(relPath))
		deleteFull = filepath.Clean(deleteFull)
		if !strings.HasPrefix(deleteFull, rootDir) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Method == http.MethodPost {
			handleDelete(w, r, deleteFull, relPath)
		} else {
			handleDeleteConfirm(w, r, deleteFull, relPath)
		}

	case strings.HasPrefix(cleanURL, "/rename/"):
		relPath := "/" + strings.TrimPrefix(cleanURL, "/rename/")
		renameFull := filepath.Join(rootDir, filepath.FromSlash(relPath))
		renameFull = filepath.Clean(renameFull)
		if !strings.HasPrefix(renameFull, rootDir) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Method == http.MethodPost {
			handleRename(w, r, renameFull, relPath)
		} else {
			handleRenameForm(w, r, renameFull, relPath)
		}

	case strings.HasPrefix(cleanURL, "/zip/"):
		relPath := "/" + strings.TrimPrefix(cleanURL, "/zip/")
		zipFull := filepath.Join(rootDir, filepath.FromSlash(relPath))
		zipFull = filepath.Clean(zipFull)
		if !strings.HasPrefix(zipFull, rootDir) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		handleZipDownload(w, r, zipFull)

	case strings.HasPrefix(cleanURL, "/mkdir/"):
		relPath := "/" + strings.TrimPrefix(cleanURL, "/mkdir/")
		mkdirFull := filepath.Join(rootDir, filepath.FromSlash(relPath))
		mkdirFull = filepath.Clean(mkdirFull)
		if !strings.HasPrefix(mkdirFull, rootDir) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		handleCreateDir(w, r, mkdirFull, relPath)

	case strings.HasPrefix(cleanURL, "/mkfile/"):
		relPath := "/" + strings.TrimPrefix(cleanURL, "/mkfile/")
		mkfileFull := filepath.Join(rootDir, filepath.FromSlash(relPath))
		mkfileFull = filepath.Clean(mkfileFull)
		if !strings.HasPrefix(mkfileFull, rootDir) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		handleCreateFile(w, r, mkfileFull, relPath)

	default:
		stat, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !stat.IsDir() {
			dlCountsMu.Lock()
			downloadCounts[cleanURL]++
			dlCountsMu.Unlock()
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
	err := r.ParseMultipartForm(1 << 30)
	if err != nil {
		http.Error(w, "upload too large or invalid", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		fallback, fh, err2 := r.FormFile("file")
		if err2 != nil {
			http.Error(w, "no files provided", http.StatusBadRequest)
			return
		}
		fallback.Close()
		files = []*multipart.FileHeader{fh}
	}

	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			continue
		}
		destPath := filepath.Join(targetDir, filepath.Base(header.Filename))
		destFile, err := os.Create(destPath)
		if err != nil {
			file.Close()
			continue
		}
		io.Copy(destFile, file)
		destFile.Close()
		file.Close()
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func handleCreateDir(w http.ResponseWriter, r *http.Request, targetDir, redirectURL string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	dirName := r.FormValue("dirname")
	if dirName == "" {
		http.Error(w, "dirname required", http.StatusBadRequest)
		return
	}
	dirName = filepath.Base(dirName)
	newDir := filepath.Join(targetDir, dirName)
	if err := os.Mkdir(newDir, 0755); err != nil {
		http.Error(w, "failed to create directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func handleCreateFile(w http.ResponseWriter, r *http.Request, targetDir, redirectURL string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fileName := r.FormValue("filename")
	if fileName == "" {
		http.Error(w, "filename required", http.StatusBadRequest)
		return
	}
	fileName = filepath.Base(fileName)
	newFile := filepath.Join(targetDir, fileName)
	f, err := os.OpenFile(newFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	f.Close()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func handleEdit(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if !isTextFile(fullPath) {
		http.Error(w, "Cannot edit binary files", http.StatusForbidden)
		return
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	cleanURL := filepath.Dir(relPath)
	if cleanURL == "." {
		cleanURL = "/"
	}
	pageData := EditPageData{
		CurrentPath:  cleanURL,
		Breadcrumbs:  buildBreadcrumbs(cleanURL),
		CSS:          template.CSS(cssContent),
		FileName:     filepath.Base(relPath),
		FilePath:     relPath,
		Content:      string(data),
		SaveURL:      "/edit" + relPath,
		HighlightJS:  template.JS(hlJS),
		HighlightCSS: template.CSS(hlCSS),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := tmplEditorLite
	if !isLiteRequest(r) {
		t = tmplEditor
	}
	t.Execute(w, pageData)
}

func handleSaveEdit(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	content := r.FormValue("content")
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, relPath, http.StatusSeeOther)
}

func handleDeleteConfirm(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	stat, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	name := filepath.Base(relPath)
	cleanURL := filepath.Dir(relPath)
	if cleanURL == "." {
		cleanURL = "/"
	}
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
		ActionURL:   "/delete" + relPath,
		CancelURL:   relPath,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmplConfirm.Execute(w, pageData)
}

func handleDelete(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	parentDir := filepath.Dir(relPath)
	if parentDir == "." {
		parentDir = "/"
	}

	stat, _ := os.Stat(fullPath)
	var err error
	if stat.IsDir() {
		err = os.RemoveAll(fullPath)
	} else {
		err = os.Remove(fullPath)
	}
	if err != nil {
		http.Error(w, "Failed to delete: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, parentDir, http.StatusSeeOther)
}

func handleRenameForm(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	cleanURL := filepath.Dir(relPath)
	if cleanURL == "." {
		cleanURL = "/"
	}
	pageData := RenamePageData{
		CurrentPath: cleanURL,
		Breadcrumbs: buildBreadcrumbs(cleanURL),
		CSS:         template.CSS(cssContent),
		OldName:     filepath.Base(relPath),
		NewURL:      "/rename" + relPath,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmplRename.Execute(w, pageData)
}

func handleRename(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	newName := r.FormValue("newname")
	if newName == "" {
		http.Error(w, "newname required", http.StatusBadRequest)
		return
	}
	newName = filepath.Base(newName)
	parentDir := filepath.Dir(fullPath)
	newPath := filepath.Join(parentDir, newName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err := os.Rename(fullPath, newPath); err != nil {
		http.Error(w, "Failed to rename: "+err.Error(), http.StatusInternalServerError)
		return
	}

	relParent := filepath.Dir(relPath)
	if relParent == "." {
		relParent = "/"
	}
	http.Redirect(w, r, relParent, http.StatusSeeOther)
}

func handleZipDownload(w http.ResponseWriter, r *http.Request, fullPath string) {
	stat, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if !stat.IsDir() {
		http.ServeFile(w, r, fullPath)
		return
	}

	name := filepath.Base(fullPath) + ".zip"
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))

	pr, pw := io.Pipe()
	go func() {
		zw := zip.NewWriter(pw)
		filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			relPath, _ := filepath.Rel(fullPath, path)
			if relPath == "." {
				return nil
			}
			relPath = filepath.ToSlash(relPath)
			if info.IsDir() {
				_, err := zw.Create(relPath + "/")
				return err
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			w, err := zw.Create(relPath)
			if err != nil {
				return err
			}
			io.Copy(w, f)
			return nil
		})
		zw.Close()
		pw.Close()
	}()

	io.Copy(w, pr)
}
