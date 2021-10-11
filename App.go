package main

import (
	"Glass/tools"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var baseDir = flag.String("dir", ".", "get base directory for listing")
var thumbnails = flag.String("thumbnails-dir", "thumbnails", "get thumbnail directory for storage")
var theme = flag.String("theme", "default", "choose templating theme")

func init() {
	flag.Parse()

}

func routines() {
	go func() {
		for {
			tools.CalculateSizes(*baseDir)
			time.Sleep(6 * time.Hour)
		}
	}()
}

func main() {
	processStyleAndJs()
	routines()
	r := gin.Default()
	r.Use(location.Default())
	r.SetFuncMap(tools.FuncMap())
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
		if !strings.HasSuffix(path, "/") {
			c.Request.URL.Path += "/"
			c.Redirect(302, c.Request.URL.String())
		}

		if hidden, ok := c.GetQuery("hidden"); ok {
			loadPreview(c, path, hidden)
			return
		}
		sortBy := c.DefaultQuery("sortBy", "name")
		sortDir := c.DefaultQuery("sortDir", "asc")
		layout := c.DefaultQuery("layout", "index")
		fullList := c.DefaultQuery("full", "no")
		api := c.DefaultQuery("api", "no")
		compact := c.DefaultQuery("compact", "no")
		listing := tools.GetSortedFileListing(path, sortBy, sortDir, fullList != "no")
		listing.Info = info
		listing.Path = c.Param("path")
		listing.Request = c.Request
		listing.Mode = "simple"
		if fullList != "no" {
			listing.Mode = "full"
		}
		if api != "no" {
			loc := location.Get(c)
			c.Request.URL.Host = loc.Host
			c.Request.URL.Scheme = loc.Scheme
			for k, v := range listing.Files {
				tmpNewUrl := *c.Request.URL
				tmpNewUrl.Path = tmpNewUrl.Path + v.Name
				listing.Files[k].Url = tmpNewUrl.String()

			}
			if compact == "no" {
				c.JSON(200, gin.H{"path": listing.Path, "files": listing.Files})
			} else {
				c.JSON(200, listing.Compact())

			}
			return
		}
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

func loadPreview(c *gin.Context, path string, hidden string) {
	path += "/." + hidden
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		processFile(c, path, false)
		return
	}
	switch hidden {
	case "dir":

		c.Redirect(302, "/_assets/icons/folder.svg")
	default:
		c.Writer.Write([]byte{})
	}
}

func processStyleAndJs() {
	content, err := ioutil.ReadFile("public/build/manifest.json")
	if err != nil {
		return
	}
	var manifest struct {
		AppJs struct {
			File    string   `json:"file"`
			Src     string   `json:"src"`
			IsEntry bool     `json:"isEntry"`
			CSS     []string `json:"css"`
		} `json:"app.js"`
	}
	err = json.Unmarshal(content, &manifest)
	if err != nil {
		return
	}
	js, err := ioutil.ReadFile("public/build/" + manifest.AppJs.File)
	if err != nil {
		return
	}
	style, err := ioutil.ReadFile("public/build/" + manifest.AppJs.CSS[0])
	if err != nil {
		return
	}
	ioutil.WriteFile("public/build/app.css", style, 0644)
	ioutil.WriteFile("public/build/app.js", js, 0644)

}
