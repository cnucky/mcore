package config

import (
	"encoding/json"
	"os"
)

func Load(filename string, prefs interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&prefs)
}
