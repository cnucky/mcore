package config

import (
	"encoding/json"
	"fmt"
	"github.com/xsnews/mcore/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	Extension = ".json"
)

func Files(basedir string) (map[string]string, error) {
	out := make(map[string]string)
	err := filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			// Only read files
			return nil
		}
		if basedir == path {
			// Ignore
			return nil
		}
		log.Debug("Found config=%s", path)
		out[f.Name()] = path
		return nil
	})

	return out, err
}

func LoadJsonD(basedir string, x interface{}) error {
	files, err := Files(basedir)
	if err != nil {
		panic(err)
	}

	if len(files) == 0 {
		return nil
	}

	// Read content from all files
	data := make(map[string]string)
	for x, y := range files {
		fh, err := os.Open(y)
		if err != nil {
			panic(err)
		}
		defer fh.Close()
		content, err := ioutil.ReadFile(y)
		if err != nil {
			panic(err)
		}
		data[x] = string(content)
	}

	// Create one big json string where every file is a key
	jsonCollection := ""
	i := 0
	for fn, c := range data {
		i++

		// Only load directory.d/file.Extension
		s := strings.Split(fn, Extension)
		if len(s) != 2 {
			panic("found an invalid file in directory")
		}

		// Add to our json structure with key "filename"
		jsonCollection = jsonCollection + fmt.Sprintf("\"%s\": %s", s[0], c)
		if len(data) == i {
			// We're done, don't add trailing comma
			break
		}

		// Add trailing comma
		jsonCollection = jsonCollection + ",\n"
	}

	// Finish json structure
	jsonCollection = "{\n" + jsonCollection + "}\n"

	// Unmarshal json
	err = json.Unmarshal([]byte(jsonCollection), &x)
	if err != nil {
		panic(err)
	}

	return nil
}
