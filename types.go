package main

import "html/template"

type FileData struct {
	Name       string
	IsDir      bool
	Size       int64
	SizeStr    string
	ModTime    string
	ModTimeRel string
	ModTimeISO string
	URLPath    string
	Ext        string
	IsText     bool
	IsImage    bool
	IsVideo    bool
	DlCount    int
	EditURL    string
	DeleteURL  string
	RenameURL  string
	ZipURL     string
}

type Breadcrumb struct {
	Name    string
	URLPath string
	IsLast  bool
}

type DirStats struct {
	TotalSize    int64
	TotalSizeStr string
	FileCount    int
	DirCount     int
}

type PageData struct {
	CurrentPath string
	Breadcrumbs []Breadcrumb
	Files       []FileData
	CSS         template.CSS
	Stats       DirStats
	Lite        bool
	Version     string
	IsRoot      bool
}

type EditPageData struct {
	CurrentPath  string
	Breadcrumbs  []Breadcrumb
	CSS          template.CSS
	FileName     string
	FilePath     string
	Content      string
	SaveURL      string
	HighlightJS  template.JS
	HighlightCSS template.CSS
}

type ConfirmPageData struct {
	CurrentPath string
	Breadcrumbs []Breadcrumb
	CSS         template.CSS
	Title       string
	Message     string
	ActionURL   string
	CancelURL   string
}

type RenamePageData struct {
	CurrentPath string
	Breadcrumbs []Breadcrumb
	CSS         template.CSS
	OldName     string
	NewURL      string
}
