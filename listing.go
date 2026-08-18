package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func formatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatRelativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case d < 30*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case d < 365*24*time.Hour:
		months := int(d.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(d.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

func getFileExt(name string) string {
	ext := filepath.Ext(name)
	if ext == "" {
		return ""
	}
	return strings.ToLower(ext[1:])
}

func isTextFile(name string) bool {
	textExts := map[string]bool{
		"go": true, "py": true, "js": true, "ts": true, "jsx": true, "tsx": true,
		"c": true, "h": true, "cpp": true, "hpp": true, "cc": true, "cxx": true,
		"rs": true, "java": true, "kt": true, "swift": true,
		"rb": true, "php": true, "lua": true, "pl": true, "r": true,
		"sh": true, "bash": true, "zsh": true, "fish": true, "ps1": true,
		"txt": true, "md": true, "markdown": true, "rst": true, "log": true,
		"json": true, "yaml": true, "yml": true, "toml": true, "ini": true, "xml": true, "csv": true,
		"html": true, "htm": true, "css": true, "scss": true, "less": true,
		"sql": true, "graphql": true, "gql": true,
		"env": true, "cfg": true, "conf": true,
		"gitignore": true, "gitattributes": true, "dockerignore": true,
		"makefile": true, "dockerfile": true, "containerfile": true,
		"license": true, "readme": true, "todo": true, "authors": true, "changelog": true,
		"mod": true, "sum": true, "lock": true,
		"vue": true, "svelte": true, "astro": true,
		"bat": true, "cmd": true,
		"tf": true, "hcl": true,
		"proto": true, "thrift": true,
		"zig": true, "nim": true, "ex": true, "exs": true, "erl": true,
		"hs": true, "ml": true, "elm": true, "clj": true,
		"dart": true, "gradle": true,
		"properties": true,
	}
	ext := getFileExt(name)
	if textExts[ext] {
		return true
	}
	base := strings.ToLower(filepath.Base(name))
	knownFiles := map[string]bool{
		"makefile": true, "dockerfile": true, "containerfile": true,
		"justfile": true, "cmakelists.txt": true,
		".gitignore": true, ".gitattributes": true, ".dockerignore": true, ".editorconfig": true,
		"license": true, "copying": true, "readme": true, "todo": true, "authors": true, "changelog": true,
	}
	return knownFiles[base]
}

func isImageFile(name string) bool {
	imgExts := map[string]bool{
		"jpg": true, "jpeg": true, "png": true, "gif": true,
		"webp": true, "svg": true, "bmp": true, "ico": true, "avif": true,
	}
	return imgExts[getFileExt(name)]
}

func isVideoFile(name string) bool {
	videoExts := map[string]bool{
		"mp4": true, "webm": true, "ogg": true, "mov": true, "avi": true, "mkv": true,
	}
	return videoExts[getFileExt(name)]
}

func calculateDirStats(entries []os.DirEntry) DirStats {
	stats := DirStats{}
	for _, entry := range entries {
		if entry.IsDir() {
			stats.DirCount++
		} else {
			stats.FileCount++
			info, err := entry.Info()
			if err == nil {
				stats.TotalSize += info.Size()
			}
		}
	}
	stats.TotalSizeStr = formatSize(stats.TotalSize)
	return stats
}

func buildFileList(entries []os.DirEntry, cleanURL string) ([]FileData, DirStats) {
	var files []FileData
	stats := DirStats{}

	if cleanURL != "/" && cleanURL != "." {
		parentURL := filepath.Dir(cleanURL)
		if parentURL == "." {
			parentURL = "/"
		}
		files = append(files, FileData{
			Name:       "..",
			IsDir:      true,
			SizeStr:    "-",
			ModTime:    "-",
			ModTimeRel: "-",
			URLPath:    parentURL,
		})
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if entry.IsDir() {
			stats.DirCount++
		} else {
			stats.FileCount++
			stats.TotalSize += info.Size()
		}

		sizeStr := "-"
		var sizeVal int64
		if !entry.IsDir() {
			sizeStr = formatSize(info.Size())
			sizeVal = info.Size()
		}

		fileURLPath := filepath.ToSlash(filepath.Join(cleanURL, entry.Name()))
		if !strings.HasPrefix(fileURLPath, "/") {
			fileURLPath = "/" + fileURLPath
		}

		ext := getFileExt(entry.Name())
		isText := !entry.IsDir() && isTextFile(entry.Name())
		isImg := !entry.IsDir() && isImageFile(entry.Name())
		isVid := !entry.IsDir() && isVideoFile(entry.Name())

		editURL := ""
		deleteURL := "/delete" + fileURLPath
		renameURL := "/rename" + fileURLPath
		zipURL := ""
		if entry.IsDir() {
			zipURL = "/zip" + fileURLPath
			editURL = ""
		} else if isText {
			editURL = "/edit" + fileURLPath
		}

		dlCount := 0
		dlCountsMu.Lock()
		if c, ok := downloadCounts[fileURLPath]; ok {
			dlCount = c
		}
		dlCountsMu.Unlock()

		files = append(files, FileData{
			Name:       entry.Name(),
			IsDir:      entry.IsDir(),
			Size:       sizeVal,
			SizeStr:    sizeStr,
			ModTime:    info.ModTime().Format("2006-01-02 15:04"),
			ModTimeRel: formatRelativeTime(info.ModTime()),
			ModTimeISO: info.ModTime().Format(time.RFC3339),
			URLPath:    fileURLPath,
			Ext:        ext,
			IsText:     isText,
			IsImage:    isImg,
			IsVideo:    isVid,
			DlCount:    dlCount,
			EditURL:    editURL,
			DeleteURL:  deleteURL,
			RenameURL:  renameURL,
			ZipURL:     zipURL,
		})
	}

	sortFiles(files)
	stats.TotalSizeStr = formatSize(stats.TotalSize)
	return files, stats
}

func sortFiles(files []FileData) {
	sort.Slice(files, func(i, j int) bool {
		if files[i].Name == ".." {
			return true
		}
		if files[j].Name == ".." {
			return false
		}
		if files[i].IsDir && !files[j].IsDir {
			return true
		}
		if !files[i].IsDir && files[j].IsDir {
			return false
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})
}

func buildBreadcrumbs(urlPath string) []Breadcrumb {
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")

	var crumbs []Breadcrumb
	crumbs = append(crumbs, Breadcrumb{Name: "home", URLPath: "/", IsLast: len(parts) == 0 || parts[0] == ""})

	if len(parts) == 0 || parts[0] == "" {
		return crumbs
	}

	currentPath := ""
	for i, part := range parts {
		currentPath += "/" + part
		crumbs = append(crumbs, Breadcrumb{
			Name:    part,
			URLPath: currentPath,
			IsLast:  i == len(parts)-1,
		})
	}
	return crumbs
}
