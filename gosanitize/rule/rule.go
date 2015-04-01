package rule

import (
	"fmt"
	"reflect"
)

type Rule struct{}
type RuleFunc func(interface{}) error

var rules map[string]RuleFunc

type RuleValidator struct {
	Id string

	Params interface{}
}

func init() {
	rules = make(map[string]RuleFunc)

	/* Add pre-defined rules */
	AddRule("Logger", Logger)
}

func AddRule(id string, fn RuleFunc) {
	rules[id] = fn
}

/* Validate params */
func (i *RuleValidator) Validate() (bool, []error) {
	var errCollection []error
	var value interface{}

	s := reflect.Indirect(reflect.ValueOf(i.Params))
	for num := 0; num < s.NumField(); num++ {
		name := s.Type().Field(num).Name
		if name == "_" {
			/* Unexported field, generic struct rule */
			value = i.Params
		} else {
			/* Exported field, specific field rule */
			value = s.Field(num).Interface()
		}
		ruleTag := s.Type().Field(num).Tag.Get("rule")
		if ruleTag == "" {
			continue
		}

		if key, exists := rules[ruleTag]; exists {
			err := key(value)
			if err != nil {
				errCollection = append(errCollection, err)
			}
		} else {
			panic(fmt.Sprintf("validate rule '%s' not found", ruleTag))
		}
	}

	if len(errCollection) > 0 {
		return false, errCollection
	}

	return true, nil
}

func NewValidator(id string, params interface{}) *RuleValidator {
	return &RuleValidator{
		Id:     id,
		Params: params,
	}
}
