package config

import (
	"path/filepath"
	"os"
	"github.com/xsnews/mcore/log"
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
