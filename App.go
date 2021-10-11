package main

import (
	"Glass/tools"
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

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
	routines()
	r := gin.Default()
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
			path += "/"
			c.Redirect(302, path)
		}

		if hidden, ok := c.GetQuery("hidden"); ok {
			loadPreview(c, path, hidden)
			return
		}
		sortBy := c.DefaultQuery("sortBy", "name")
		sortDir := c.DefaultQuery("sortDir", "asc")
		layout := c.DefaultQuery("layout", "index")
		fullList := c.DefaultQuery("full", "no")
		listing := tools.GetSortedFileListing(path, sortBy, sortDir, fullList != "no")
		listing.Info = info
		listing.Path = c.Param("path")
		listing.Request = c.Request
		listing.Mode = "simple"
		if fullList != "no" {
			listing.Mode = "full"
		}
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