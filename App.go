package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
)

var baseDir = flag.String("dir", ".", "get base directory for listing")
var thumbnails = flag.String("thumbnails-dir", "thumbnails", "get thumbnail directory for storage")
var theme = flag.String("theme", "default", "choose templating theme")

func init() {
	flag.Parse()

}

func main() {

	r := gin.Default()
	fmap := sprig.FuncMap()
	fmap["dict"] = dict
	fmap["hrsize"] = prettyBytes
	fmap["urlAddQuery"] = addUrlArg
	fmap["ftype"] = fileType
	fmap["last"] = func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	}
	r.SetFuncMap(fmap)
	r.LoadHTMLGlob("templates/" + *theme + "/*")
	r.Use(static.Serve("/_assets", static.LocalFile("public", false)))

	r.GET("/*path", func(c *gin.Context) {
		path := *baseDir + c.Param("path")
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			c.JSON(404, gin.H{"message": "not found"})
			return
		}
		if !info.IsDir() {
			_, ok := c.GetQuery("thumbnail")
			processFile(c, path, ok)

			return
		}
		sortBy := c.DefaultQuery("sortBy", "name")
		sortDir := c.DefaultQuery("sortDir", "asc")
		layout := c.DefaultQuery("layout", "index")
		listing := getSortedFileListing(path, sortBy, sortDir)
		listing.Info = info
		listing.Path = c.Param("path")
		listing.Request = c.Request
		c.Request.URL.Query()
		c.HTML(http.StatusOK, layout+".html", listing)
	})

	r.Run("0.0.0.0:8000")
}

func processFile(c *gin.Context, path string, thumbnail bool) {
	if thumbnail {
		thumb, _ := LoadThumbnail(path)
		c.Writer.Write(thumb)
		return
	}
	c.File(path)

}

func getSortedFileListing(path string, sortBy string, sortDir string) *Listing {
	l := Listing{
		SortBy:  sortBy,
		SortDir: sortDir,
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(files, func(i, j int) bool {
		switch sortBy {
		case "size":
			if sortDir == "asc" {
				return files[i].Size() < files[j].Size()
			}
			return files[i].Size() > files[j].Size()
		case "date":
			if sortDir == "asc" {
				return files[i].ModTime().Before(files[j].ModTime())
			}
			return files[j].ModTime().Before(files[i].ModTime())
		}
		if sortDir == "asc" {
			return strings.Compare(strings.ToLower(files[i].Name()), strings.ToLower(files[j].Name())) < 1
		}
		return strings.Compare(strings.ToLower(files[i].Name()), strings.ToLower(files[j].Name())) > -1
	})

	l.Files = files

	return &l
}

type Listing struct {
	Files   []fs.FileInfo
	SortBy  string
	SortDir string
	Info    os.FileInfo
	Path    string
	Request *http.Request
}

func fileType(n string) string {
	n = strings.ToLower(n)
	types := map[string][]string{
		"image": {"bmp", "gif", "jpg", "jpeg", "png", "tiff", "svg"},
		"video": {"mp4", "mpeg", "avi", "mkv", "ogv", "webm", "3gp", "3g2", "mov"},
		"audio": {"oga", "mid", "midi", "mp3", "opus", "wav", "weba", "aac"},
		"doc":   {"md", "pdf", "html", "txt"},
	}

	for t, ext := range types {
		for _, e := range ext {
			if strings.HasSuffix(n, "."+e) {
				return t
			}
		}
	}
	return "other"

}

func prettyBytes(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func addUrlArg(u *url.URL, k string, v string) *url.URL {
	query := u.Query()
	query.Del(k)
	query.Add(k, v)
	u.RawQuery = query.Encode()
	return u
}
