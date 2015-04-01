package configloader

import (
	"encoding/json"
	"os"
)

var Prefs *interface{}

func Load(filename string, prefs interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&prefs)
	if err != nil {
		return err
	}
	Prefs = &prefs

	return nil
}
