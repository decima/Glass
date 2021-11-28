package models

import (
	"net/http"
	"os"
	"time"
)

type Listing struct {
	Files   []FileInfo
	SortBy  string
	SortDir string
	Info    os.FileInfo
	Path    string
	Request *http.Request
	Mode    string
}
type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}