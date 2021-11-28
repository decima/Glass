package tools

import (
	"os"
	"path/filepath"
	"sync"
)

type inMemoryItem struct {
	size int64
}

var sizes = map[string]inMemoryItem{}
var mutex = sync.Mutex{}

func GetSize(path string, info os.FileInfo) int64 {
	size := info.Size()
	if info.IsDir() {
		defer mutex.Unlock()
		mutex.Lock()
		if im, ok := sizes[path+info.Name()]; ok {
			return im.size
		}
		size, _ = DirSize(path + info.Name())
		sizes[path+info.Name()] = inMemoryItem{
			size: size,
		}
	}
	return size
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func CalculateSizes(path string) error {
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			go func(p string, info os.FileInfo) {
				res, err := DirSize(p)
				if err != nil {
					res = 0
				}
				mutex.Lock()
				sizes[path+info.Name()] = inMemoryItem{
					size: res,
				}
				mutex.Unlock()
			}(p, info)
		}
		return nil
	})
	return err

}
