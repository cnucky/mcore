package config

import (
	"path/filepath"
	"os"
	"github.com/xsnews/mcore/log"
)

// Parse JSON's recursively in given dir.
func LoadDir(basedir string, retType interface{}) ([]interface{}, error) {
	var out []interface{}

	err := filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			// Only read dirs
			return nil
		}
		if basedir == path {
			// Ignore
			return nil
		}
		log.Debug("Load config=%s", path)

		if e := Load(path, &retType); e != nil {
			return e
		}
		out = append(out, retType)
		return nil
	})

	return out, err
}
