package models

import (
	"net/http"
	"os"
	"time"
)

type Listing struct {
	Files   []FileInfo `json:"files"`
	SortBy  string
	SortDir string
	Info    os.FileInfo
	Path    string `json:"path"`
	Request *http.Request
	Mode    string
}
type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"-"`
	ModTime time.Time   `json:"updated"`
	IsDir   bool        `json:"-"`
	Url     string      `json:"url"`
}

func (l *Listing) Compact() []string {
	compactList := []string{}
	for _, v := range l.Files {

		res := v.Url
		if len(res) == 0 {
			res = v.Name
		}
		compactList = append(compactList, res)
	}
	return compactList
}
