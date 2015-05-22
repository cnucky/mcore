package valid

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/schema"
)

func ParseForm(input interface{}, r *http.Request) error {
	if e := r.ParseForm(); e != nil {
		return e
	}
	return schema.NewDecoder().Decode(input, r.PostForm)
}

func ParseJson(input interface{}, r *http.Request) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(input)
}