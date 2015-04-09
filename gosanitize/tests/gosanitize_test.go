package main

import (
	"errors"
	"github.com/xsnews/microservice-core/gosanitize/rule"
	"github.com/xsnews/microservice-core/gosanitize/util"
	"github.com/xsnews/microservice-core/gosanitize/validate"
	"io/ioutil"
	"testing"
)

/* Struct of all input keys */
type TestInput1 struct {
	Code        string    `json:",omitempty"`
	TestDep     string    `json:",omitempty"`
	Conditional string    `json:",omitempty"`
	Int         int       `json:",omitempty"`
	Bool        bool      `json:",omitempty"`
	Array       []int     `json:",omitempty"`
	Email       string    `json:",omitempty" rule:"RuleField"`   // Rule for extended validation
	_           rule.Rule `json:",omitempty" rule:"RuleGeneric"` // Generic validation rule not specific to a field
}

func RuleField(obj interface{}) error {
	if obj.(string) == "test@no.such.domain" {
		return errors.New("Invalid email address")
	}

	return nil
}

func RuleGeneric(obj interface{}) error {
	return nil
}

func init() {
	rule.AddRule("RuleField", RuleField)
	rule.AddRule("RuleGeneric", RuleGeneric)
}

func BenchmarkValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestValues1 := map[string]string{
			"Code":  "Hello World",
			"Int":   "10",
			"Bool":  "true",
			"Email": "test@gmail.com",
		}

		RunValidate("test1", "./schemas/test1.json", &TestInput1{}, &TestValues1)
	}
}

func RunValidate(id string, filename string, params interface{}, values *map[string]string) bool {
	/* Load schema */
	schema, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	/* Create param struct and get validator */
	v := validate.NewValidator(id, schema, params)

	/* Inject input values from map into params */
	err = util.LoadFromMap(params, *values)

	/* Hack for arrays since loader doesn't support these yet */
	params.(*TestInput1).Array = []int{10, 20}

	if err != nil {
		return false
	}

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

func TestInputIsValidOK(t *testing.T) {
	TestValues1 := map[string]string{
		"Code":  "Hello World",
		"Int":   "10",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	validate := RunValidate("test1", "./schemas/test1.json", &TestInput1{}, &TestValues1)
	if !validate {
		t.FailNow()
	}
}

func TestDependencyOK(t *testing.T) {
	TestValues1 := map[string]string{
		"Code":        "Hello World",
		"TestDep":     "wat",
		"Conditional": "ok",
		"Int":         "10",
		"Bool":        "true",
		"Email":       "test@gmail.com",
	}

	validate := RunValidate("test1", "./schemas/test1.json", &TestInput1{}, &TestValues1)
	if !validate {
		t.FailNow()
	}
}

func TestDependencyFail(t *testing.T) {
	TestValues1 := map[string]string{
		"Code":    "Hello World",
		"TestDep": "wat",
		"Int":     "10",
		"Bool":    "true",
		"Email":   "test@gmail.com",
	}

	validate := RunValidate("test1", "./schemas/test1.json", &TestInput1{}, &TestValues1)
	if validate {
		t.FailNow()
	}
}

func TestMissingArgFail(t *testing.T) {
	TestValues1 := map[string]string{
		"Int":   "10",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	validate := RunValidate("test1", "./schemas/test1.json", &TestInput1{}, &TestValues1)
	if validate {
		t.FailNow()
	}
}

func TestRegexMatchFail(t *testing.T) {
	TestValues1 := map[string]string{
		"Code":  "Hello World",
		"Int":   "10",
		"Bool":  "true",
		"Email": "t est@gmail.com",
	}

	validate := RunValidate("test1", "./schemas/test1.json", &TestInput1{}, &TestValues1)
	if validate {
		t.FailNow()
	}
}

func TestFieldRuleFail(t *testing.T) {
	TestValues1 := map[string]string{
		"Code":  "Hello World",
		"Int":   "10",
		"Bool":  "true",
		"Email": "test@no.such.domain",
	}

	validate := RunValidate("test1", "./schemas/test1.json", &TestInput1{}, &TestValues1)
	if validate {
		t.FailNow()
	}
}
