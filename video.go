package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func init() {
}

func getSanitizedBaseDir(path string) string {
	path = strings.ReplaceAll(path, *baseDir, "")
	path = strings.ReplaceAll(path, "/", "_")
	return path
}

func LoadThumbnail(path string) ([]byte, error) {
	h := sha1.New()
	h.Write([]byte(path))
	hash := hex.EncodeToString(h.Sum(nil))
	thumbPath := *thumbnails + "/" + getSanitizedBaseDir(path) + "_" + hash
	if _, err := os.Stat(thumbPath); errors.Is(err, os.ErrNotExist) {
		e := bytes.NewBufferString("")
		for frame := 100; frame >= 0; frame -= 25 {
			e = Thumbnail(path, frame)
			if(len(e.Bytes())>0){
				break
			}
		}
		err := ioutil.WriteFile(thumbPath, e.Bytes(), 0644)
		if err != nil {

			log.Println(err)
			panic("end")
		}
		return e.Bytes(), nil
	}
	return ioutil.ReadFile(thumbPath)
}

func Thumbnail(path string, frameNum int) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	ffmpeg.Input(path).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	return buf
}


