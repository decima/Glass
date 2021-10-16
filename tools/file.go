package tools

import (
	"Glass/models"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GetSortedFileListing(path string, sortBy string, sortDir string, fullList bool) *models.Listing {
	l := models.Listing{
		SortBy:  sortBy,
		SortDir: sortDir,
		Files:   []models.FileInfo{},
	}
	if fullList {
		filepath.Walk(path, func(p string, file os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			path = strings.TrimPrefix(path, ".")
			path = strings.Trim(path, "/")
			p = strings.Trim(p, "/")

			p = strings.TrimPrefix(p, path+"/")
			if strings.HasPrefix(file.Name(), ".") {
				return nil
			}
			if !file.IsDir() {
				f := models.FileInfo{
					Name:    p,
					Size:    file.Size(),
					Mode:    file.Mode(),
					ModTime: file.ModTime(),
					IsDir:   file.IsDir(),
				}
				l.Files = append(l.Files, f)
			}
			return nil
		})
	} else {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}
			f := models.FileInfo{
				Name:    file.Name(),
				Size:    GetSize(path, file),
				Mode:    file.Mode(),
				ModTime: file.ModTime(),
				IsDir:   file.IsDir(),
			}
			l.Files = append(l.Files, f)

		}
	}
	sort.Slice(l.Files, func(i, j int) bool {
		switch sortBy {
		case "size":
			if sortDir == "asc" {
				return l.Files[i].Size < l.Files[j].Size
			}
			return l.Files[i].Size > l.Files[j].Size
		case "date":
			if sortDir == "asc" {
				return l.Files[i].ModTime.Before(l.Files[j].ModTime)
			}
			return l.Files[j].ModTime.Before(l.Files[i].ModTime)
		}
		if sortDir == "asc" {
			return strings.Compare(strings.ToLower(l.Files[i].Name), strings.ToLower(l.Files[j].Name)) < 1
		}
		return strings.Compare(strings.ToLower(l.Files[i].Name), strings.ToLower(l.Files[j].Name)) > -1
	})

	return &l
}
