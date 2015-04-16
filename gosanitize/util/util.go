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

		err := Set(&field, postValue)
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
		j.Decode(v.Data)
	} else {
		/* Parse http form */
		r.ParseForm()
	}

	s := reflect.Indirect(reflect.ValueOf(v.Data))
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

	v.dataOk = true
	return nil
}

type ValidateResult struct {
	GoError      error
	SchemaErrors *gojsonschema.Result
	RuleErrors   []error
}

func (v *ValidateResult) hasError() bool {
	if v.GoError != nil || v.SchemaErrors != nil || len(v.RuleErrors) > 0 {
		return true
	}

	return false
}

/* Load params from map */
func LoadFromMap(params interface{}, values map[string]interface{}) error {
	s := reflect.Indirect(reflect.ValueOf(params))
	for num := 0; num < s.NumField(); num++ {
		field := s.Field(num)
		fieldType := s.Type().Field(num)
		postValue := values[fieldType.Name]
		if postValue == nil {
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

func SetSlice(field *reflect.Value, v interface{}) error {
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

/* Convert input value and set field */
func Set(field *reflect.Value, v interface{}) error {
	var value string

	/* Check if field is settable */
	if !field.CanSet() {
		return errors.New(fmt.Sprintf("Can't set field '%s'", field))
	}

	/* Is it a slice? */
	if reflect.ValueOf(v).Kind() == reflect.Slice {
		return SetSlice(field, v)
	}

	/* Cast value, reflect and set it */
	value = v.(string)

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

/* Helper function to validate a struct (params) with values (values) against a json schema (filename) */
func Validate(id string, schema []byte, params interface{}, values *map[string]interface{}) (bool, *ValidateResult) {
	rs := &ValidateResult{}

	/* Inject input values from map into params */
	rs.GoError = LoadFromMap(params, *values)
	if rs.GoError != nil {
		return false, rs
	}

	/* Create param struct and get validator */
	v := validate.NewValidator(id, schema, params)

	/* Validate input values with the json schema */
	validateOk, schemaErr := v.Validate()
	if !validateOk {
		rs.SchemaErrors = schemaErr
		return false, rs
	}

	r := rule.NewValidator(id, params)
	ruleOk, ruleErr := r.Validate()
	if !ruleOk {
		rs.RuleErrors = ruleErr
		return false, rs
	}

	return true, rs
}

/* Helper function to validate a struct (params) with values (values) against a json schema (filename) */
func ValidateRequest(id string, schema []byte, params interface{}, r *http.Request) (bool, *ValidateResult) {
	rs := &ValidateResult{}

	/* Needs refactoring */
	if r.Header["Content-Type"][0] == "application/json" {
		r.ParseForm()

		j := json.NewDecoder(r.Body)
		j.Decode(params)
	} else {
		rs.GoError = LoadFromRequest(params, r)
		if rs.GoError != nil {
			return false, rs
		}
	}

	/* Create param struct and get validator */
	v := validate.NewValidator(id, schema, params)

	/* Validate input values with the json schema */
	validateOk, schemaErr := v.Validate()
	if !validateOk {
		rs.SchemaErrors = schemaErr
		return false, rs
	}

	rv := rule.NewValidator(id, params)
	ruleOk, ruleErr := rv.Validate()
	if !ruleOk {
		rs.RuleErrors = ruleErr
		return false, rs
	}

	return true, rs
}
