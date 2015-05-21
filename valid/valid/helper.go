package valid

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"net/http"
	"reflect"
	"strconv"
)

func FnGetInt(i interface{}) (int64, error) {
	min, err := strconv.ParseInt(i.(string), 0, 0)
	if err != nil {
		return 0, err
	}

	return min, nil
}

func FnGetStr(i interface{}) (string, error) {
	return i.(string), nil
}

func FnGetStrSlice(i interface{}) ([]string, error) {
	return i.([]string), nil
}

/* Reflect struct */
func Validate(t interface{}) (bool, map[string][]string) {
	var collectResults map[string][]string = make(map[string][]string)

	/* Loop through each field in given struct */
	s := reflect.Indirect(reflect.ValueOf(t))
	for num := 0; num < s.NumField(); num++ {
		/* Get name of variable and check if it's exportable */
		name := s.Type().Field(num).Name
		if name == "_" || !s.Field(num).CanInterface() {
			continue
		}

		/* Skip pointers with nil values, these are optional fields */
		if s.Field(num).Kind() == reflect.Ptr {
			if s.Field(num).IsNil() {
				continue
			}
		}

		/* Get the validation tag */
		rule := s.Type().Field(num).Tag.Get("validate")
		if len(rule) == 0 {
			continue
		}

		/* This field is OK for validation, get the value from interface */
		var value interface{}
		if s.Field(num).Kind() == reflect.Ptr {
			/* Derefence pointer value */
			value = s.Field(num).Elem().Interface()
		} else {
			/* Just get the interface value */
			value = s.Field(num).Interface()
		}

		/* Create parser for this rule and pass the context to it */
		l := new(Valdsl)
		l.Debug = true
		err, results := l.Parse(t, rule, value)
		if err != nil {
			/* Deverror in rule */
			panic(err)
		}

		if len(results) > 0 {
			for _, v := range results {
				collectResults[name] = append(collectResults[name], v)
			}
		}

		/* Is this a slice? */
		if s.Type().Field(num).Type.Kind() == reflect.Slice {
			/* Loop through slice and validate each item */
			s := reflect.ValueOf(s.Field(num).Interface())
			for i := 0; i < s.Len(); i++ {
				_, results := Validate(s.Index(i).Interface())
				if len(results) > 0 {
					for _, v := range results {
						collectResults[name] = append(collectResults[name], v...)
					}
				}
			}
		} else if s.Type().Field(num).Type.Kind() == reflect.Struct {
			/* If it's a struct, only validate the struct */
			_, results := Validate(s.Field(num).Interface())
			if len(results) > 0 {
				for _, v := range results {
					collectResults[name] = append(collectResults[name], v...)
				}
			}
		}
	}

	if len(collectResults) > 0 {
		return false, collectResults
	}

	return true, collectResults
}

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
