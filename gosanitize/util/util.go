package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"github.com/xsnews/microservice-core/gosanitize/rule"
	"github.com/xsnews/microservice-core/gosanitize/validate"
	"net/http"
	"reflect"
	"strconv"
)

type Values map[string]interface{}

type Validator struct {
	Id     string
	Schema []byte

	Data   interface{}
	dataOk bool

	inputValidator *validate.InputValidator
	ruleValidator  *rule.RuleValidator
}

func NewValidator(id string, schema []byte, data interface{}) *Validator {
	return &Validator{
		Id:     id,
		Schema: schema,
		Data:   data,
		dataOk: false,
	}
}

func (v *Validator) Validate() (bool, *gojsonschema.Result) {
	if !v.dataOk {
		panic("No input data loaded")
	}

	inputValidator := validate.NewValidator(v.Id, v.Schema, v.Data)

	/* Validate input values with the json schema */
	ok, err := inputValidator.Validate()
	if !ok {
		return false, err
	}

	return true, nil
}

func (v *Validator) ValidateRules() (bool, []error) {
	if !v.dataOk {
		panic("No input data loaded")
	}

	ruleValidator := rule.NewValidator(v.Id, v.Data)

	/* Validate input values with the json schema */
	ok, err := ruleValidator.Validate()
	if !ok {
		return false, err
	}

	return true, nil
}

/* Load data from map */
func (v *Validator) LoadValues(values Values) error {
	s := reflect.Indirect(reflect.ValueOf(v.Data))
	for num := 0; num < s.NumField(); num++ {
		field := s.Field(num)
		fieldType := s.Type().Field(num)
		postValue := values[fieldType.Name]
		if postValue == nil {
			continue
		}

		err := setField(&field, postValue)
		if err != nil {
			return err
		}
	}

	v.dataOk = true
	return nil
}

/* Load params from post */
func (v *Validator) LoadValuesFromRequest(r *http.Request) error {
	if r.Header["Content-Type"][0] == "application/json" {
		j := json.NewDecoder(r.Body)
		err := j.Decode(v.Data)
		if err != nil {
			return err
		}

		v.dataOk = true
		return nil
	} else {
		/* Parse http form */
		r.ParseForm()
	}

	s := reflect.Indirect(reflect.ValueOf(v.Data))
	for num := 0; num < s.NumField(); num++ {
		field := s.Field(num)
		fieldType := s.Type().Field(num)

		if len(r.Form[fieldType.Name]) == 0 {
			continue
		}

		postValue := r.Form[fieldType.Name]
		if postValue == nil {
			continue
		}

		err := setField(&field, postValue)
		if err != nil {
			return err
		}
	}

	v.dataOk = true
	return nil
}

func setSlice(field *reflect.Value, v interface{}) error {
	/* Check if field is settable */
	if !field.CanSet() {
		return errors.New(fmt.Sprintf("Can't set field '%s'", field))
	}

	if reflect.ValueOf(v).Kind() != reflect.Slice {
		return errors.New("Not a slice")
	}

	/* Reflect tpe of slice */
	s := reflect.ValueOf(v)

	var newSlice reflect.Value
	var newSliceType reflect.Type
	newSliceElem := field.Type()

	/* TODO: Reflect fields */
	switch newSliceElem {
	case reflect.TypeOf([]string{}):
		newSliceType = reflect.TypeOf([]string{})
	case reflect.TypeOf([]int{}):
		newSliceType = reflect.TypeOf([]int{})
	case reflect.TypeOf([]bool{}):
		newSliceType = reflect.TypeOf([]bool{})
	default:
		panic(fmt.Sprintf("Invalid slice type %s", newSliceElem))
	}

	newSlice = reflect.MakeSlice(newSliceType, s.Len(), s.Cap())

	switch newSliceElem {
	case reflect.TypeOf([]string{}):
		for i := 0; i < s.Len(); i++ {
			newSlice.Index(i).SetString(s.Index(i).String())
		}
	case reflect.TypeOf([]float64{}):
		for i := 0; i < s.Len(); i++ {
			convertedValue, err := strconv.ParseFloat(s.Index(i).String(), 64)
			if err != nil {
				return err
			}

			newSlice.Index(i).SetFloat(convertedValue)
		}
	case reflect.TypeOf([]int{}):
		for i := 0; i < s.Len(); i++ {
			convertedValue, err := strconv.ParseInt(s.Index(i).String(), 0, 0)
			if err != nil {
				return err
			}

			newSlice.Index(i).SetInt(convertedValue)
		}
	case reflect.TypeOf([]bool{}):
		for i := 0; i < s.Len(); i++ {
			convertedValue, err := strconv.ParseBool(s.Index(i).String())
			if err != nil {
				return err
			}

			newSlice.Index(i).SetBool(convertedValue)
		}
	}

	field.Set(newSlice)
	return nil
}

func isSlice(field *reflect.Value) bool {
	switch field.Type() {
	case reflect.TypeOf([]string{}):
		return true
	case reflect.TypeOf([]float64{}):
		return true
	case reflect.TypeOf([]int{}):
		return true
	case reflect.TypeOf([]bool{}):
		return true
	}

	return false
}

/* Convert input value and set field */
func setField(field *reflect.Value, v interface{}) error {
	var value string

	/* Check if field is settable */
	if !field.CanSet() {
		return errors.New(fmt.Sprintf("Can't set field '%s'", field))
	}

	/* Check if supplied field is a slice */
	if isSlice(field) {
		return setSlice(field, v)
	}

	/* Check if supplied v is a slice */
	if reflect.ValueOf(v).Kind() == reflect.Slice {
		/* If the supplied value is a slice we take the 1st value (same way url.Values's get method works) */
		value = v.([]string)[0]
	} else {
		value = v.(string)
	}

	switch field.Kind() {
	case reflect.Float64:
		convertedValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(convertedValue)
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
