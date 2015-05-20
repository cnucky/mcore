package valid

import (
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
		name := s.Type().Field(num).Name
		if name == "_" || !s.Field(num).CanInterface() {
			continue
		}

		/* Exported field, specific field rule */
		value := s.Field(num).Interface()

		/* Get validation rule */
		rule := s.Type().Field(num).Tag.Get("validate")
		if len(rule) == 0 {
			continue
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
	r.ParseForm()
	decoder := schema.NewDecoder()
	err := decoder.Decode(input, r.PostForm)
	return err
}
