package tools

import (
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strings"

	"github.com/Masterminds/sprig"
)

func FuncMap() template.FuncMap {
	fmap := sprig.FuncMap()
	fmap["dict"] = Dict
	fmap["hrsize"] = PrettyBytes
	fmap["urlAddQuery"] = AddUrlArg
	fmap["changePath"] = ChangePath
	fmap["ftype"] = FileType
	fmap["last"] = func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	}
	return fmap
}

func FileType(n string) string {
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

func PrettyBytes(b int64) string {
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

func Dict(values ...interface{}) (map[string]interface{}, error) {
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

func AddUrlArg(original *url.URL, k string, v string) *url.URL {
	u, _ := url.Parse(original.String())
	query := u.Query()
	query.Del(k)
	query.Add(k, v)
	u.RawQuery = query.Encode()
	return u
}

func ChangePath(original *url.URL, path string) *url.URL {
	u, _ := url.Parse(original.String())
	u.Path = path
	return u
}
