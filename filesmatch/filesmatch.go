package filesmatch

import (
	"github.com/itshosted/mcore/log"
	"path/filepath"
	"os"
)

// Read files in dir by pattern.
// For a list of patterns: http://golang.org/pkg/path/filepath/#Match
func Match(pathGlob string) (map[string]string, error) {
	out := make(map[string]string)
	basedir := filepath.Dir(pathGlob)
	err := filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if basedir == path {
			// Ignore basedir
			return nil
		}
		match, err := filepath.Match(pathGlob, path)
		if err != nil {
			return err
		}
		if !match {
			// Ignore not matching
			return nil
		}

		log.Debug("Glob match file=%s", path)
		out[f.Name()] = path
		return nil
	})

	return out, err
}
