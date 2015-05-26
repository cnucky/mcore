package config

import (
	"encoding/json"
	"errors"
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

func files(basedir string) (map[string]string, error) {
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
	// Get list of files
	filelist, err := files(basedir)
	if err != nil {
		return err
	}

	// No files found
	if len(filelist) == 0 {
		return nil
	}

	// Create one big json string where every file is a key
	jsonCollection := ""
	for filename, fullpath := range filelist {
		// Check if we can open file
		fh, err := os.Open(fullpath)
		if err != nil {
			return err
		}
		defer fh.Close()

		// Only load directory.d/file.Extension
		s := strings.Split(filename, Extension)
		if len(s) != 2 {
			return errors.New(fmt.Sprintf("Invalid file '%s' present in config dir.", filename))
		}

		// Load content from file
		data, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return err
		}

		// Add to our json structure with key "filename"
		jsonCollection = jsonCollection + fmt.Sprintf(`"%s": %s`, s[0], data)
		// Add trailing comma
		jsonCollection = jsonCollection + ","
	}

	// Finish json structure, remove trailing comma
	jsonCollection = "{\n" + jsonCollection[:len(jsonCollection)-1] + "}\n"

	// Unmarshal json
	err = json.Unmarshal([]byte(jsonCollection), &x)
	if err != nil {
		return err
	}

	return nil
}
