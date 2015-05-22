package config

import (
	"encoding/json"
	"fmt"
	"github.com/xsnews/mcore/log"
	"io/ioutil"
	"os"
	"path/filepath"
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

func Json(basedir string, x interface{}) error {
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

	// Create k:v json from files, refactor this to a stream
	content := "{\n"
	for fn, c := range data {
		content = content + fmt.Sprintf("\"%s\": %s,", fn, c)
	}
	content = content[0 : len(content)-1] // Remove trailing ,
	content = content + "}\n"

	// Unmarshal json
	err = json.Unmarshal([]byte(content), &x)
	if err != nil {
		panic(err)
	}

	return nil
}
