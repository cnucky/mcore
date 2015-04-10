package util

import (
	"errors"
	"fmt"
	"github.com/xsnews/microservice-core/gosanitize/rule"
	"github.com/xsnews/microservice-core/gosanitize/validate"
	"net/http"
	"reflect"
	"strconv"
)

/* Load params from map */
func LoadFromMap(params interface{}, values map[string]string) error {
	s := reflect.Indirect(reflect.ValueOf(params))
	for num := 0; num < s.NumField(); num++ {
		field := s.Field(num)
		fieldType := s.Type().Field(num)
		postValue := values[fieldType.Name]
		if postValue == "" {
			continue
		}

		err := Set(&field, postValue)
		if err != nil {
			return err
		}
	}

	return nil
}

/* Load params from post */
func LoadFromRequest(params interface{}, r *http.Request) error {
	/* Parse http form */
	r.ParseForm()

	s := reflect.Indirect(reflect.ValueOf(params))
	for num := 0; num < s.NumField(); num++ {
		field := s.Field(num)
		fieldType := s.Type().Field(num)
		postValue := r.Form.Get(fieldType.Name)
		if postValue == "" {
			continue
		}

		err := Set(&field, postValue)
		if err != nil {
			return err
		}
	}

	return nil
}

/* Convert input value and set field */
func Set(field *reflect.Value, value string) error {
	if !field.CanSet() {
		return errors.New(fmt.Sprintf("Can't set field '%s'", field))
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		convertedValue, err := strconv.ParseInt(value, 0, 0)
		if err != nil {
			return err
		}
		field.SetInt(convertedValue)
	case reflect.Bool:
		convertedValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(convertedValue)
	}

	return nil
}

/* Helper function to validate a struct (params) with values (values) against a json schema (filename) */
func Validate(id string, schema []byte, params interface{}, values *map[string]string) bool {
	/* Inject input values from map into params */
	err := LoadFromMap(params, *values)
	if err != nil {
		return false
	}

	/* Create param struct and get validator */
	v := validate.NewValidator(id, schema, params)

	/* Validate input values with the json schema */
	validateOk, _ := v.Validate()
	if !validateOk {
		return false
	}

	r := rule.NewValidator(id, params)
	ruleOk, _ := r.Validate()
	if !ruleOk {
		return false
	}

	return true
}
