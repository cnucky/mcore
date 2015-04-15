package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xsnews/microservice-core/gosanitize/rule"
	"github.com/xsnews/microservice-core/gosanitize/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var SchemaContent map[string][]byte

/* Struct of all input keys */
type TestInput1 struct {
	Code        string    `json:",omitempty"`
	TestDep     string    `json:",omitempty"`
	Conditional string    `json:",omitempty"`
	Int         int       `json:",omitempty"`
	Float       float64   `json:",omitempty"`
	Bool        bool      `json:",omitempty"`
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

	SchemaContent = make(map[string][]byte)
	SchemaContent["test1"] = LoadSchema("./schemas/test1.json")
	SchemaContent["test2"] = LoadSchema("./schemas/test1.json")
}

func LoadSchema(filename string) []byte {
	schema, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return schema
}

func BenchmarkValidate(b *testing.B) {
	var validate bool

	TestValues1 := map[string]interface{}{
		"Code":  "Hello World",
		"Int":   "10",
		"Float": "1.5",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	for i := 0; i < b.N; i++ {
		validate, _ = util.Validate("test1", SchemaContent["test1"], &TestInput1{}, &TestValues1)
		if !validate {
			b.FailNow()
		}
	}
}

func TestInputIsValidOK(t *testing.T) {
	TestValues1 := map[string]interface{}{
		"Code":  "Hello World",
		"Int":   "10",
		"Float": "1.5",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	validate, _ := util.Validate("test1", SchemaContent["test1"], &TestInput1{}, &TestValues1)
	if !validate {
		t.FailNow()
	}
}

func TestDependencyOK(t *testing.T) {
	TestValues1 := map[string]interface{}{
		"Code":        "Hello World",
		"TestDep":     "wat",
		"Conditional": "ok",
		"Int":         "10",
		"Float":       "1.5",
		"Bool":        "true",
		"Email":       "test@gmail.com",
	}

	validate, _ := util.Validate("test1", SchemaContent["test1"], &TestInput1{}, &TestValues1)
	if !validate {
		t.FailNow()
	}
}

func TestDependencyFail(t *testing.T) {
	TestValues1 := map[string]interface{}{
		"Code":    "Hello World",
		"TestDep": "wat",
		"Int":     "10",
		"Float":   "1.5",
		"Bool":    "true",
		"Email":   "test@gmail.com",
	}

	validate, _ := util.Validate("test1", SchemaContent["test1"], &TestInput1{}, &TestValues1)
	if validate {
		t.FailNow()
	}
}

func TestMissingArgFail(t *testing.T) {
	TestValues1 := map[string]interface{}{
		"Int":   "10",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	validate, _ := util.Validate("test1", SchemaContent["test1"], &TestInput1{}, &TestValues1)
	if validate {
		t.FailNow()
	}
}

func TestRegexMatchFail(t *testing.T) {
	TestValues1 := map[string]interface{}{
		"Code":  "Hello World",
		"Int":   "10",
		"Float": "1.5",
		"Bool":  "true",
		"Email": "t est@gmail.com",
	}

	validate, _ := util.Validate("test1", SchemaContent["test1"], &TestInput1{}, &TestValues1)
	if validate {
		t.FailNow()
	}
}

func TestFieldRuleFail(t *testing.T) {
	TestValues1 := map[string]interface{}{
		"Code":  "Hello World",
		"Int":   "10",
		"Float": "1.5",
		"Bool":  "true",
		"Email": "test@no.such.domain",
	}

	validate, _ := util.Validate("test1", SchemaContent["test1"], &TestInput1{}, &TestValues1)
	if validate {
		t.FailNow()
	}
}

type TestInput2 struct {
	Array []int `json:",omitempty"`
}

func TestArray(t *testing.T) {
	TestValues2 := map[string]interface{}{
		"Array": []string{"10", "20"},
	}

	validate, _ := util.Validate("test1", SchemaContent["test2"], &TestInput2{}, &TestValues2)
	if validate {
		t.FailNow()
	}
}

func TestHttpRequestForm(t *testing.T) {
	var validate bool
	var result *util.ValidateResult

	TestValues1 := &url.Values{
		"Code":  {"hello world"},
		"Int":   {"10"},
		"Float": {"1.5"},
		"Bool":  {"true"},
		"Email": {"test@gmail.com"},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validate, result = util.ValidateRequest("test1", SchemaContent["test1"], &TestInput1{}, r)
	}))
	defer ts.Close()

	res, err := http.PostForm(ts.URL, *TestValues1)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !validate {
		fmt.Println(result.SchemaErrors)
		t.FailNow()
	}
}

func TestHttpRequestJson(t *testing.T) {
	var validate bool
	var result *util.ValidateResult

	TestValues1 := &TestInput1{
		Code:  "hello world",
		Int:   10,
		Float: 1.5,
		Bool:  true,
		Email: "test@gmail.com",
	}

	var b bytes.Buffer
	j := json.NewEncoder(&b)
	j.Encode(TestValues1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validate, result = util.ValidateRequest("test1", SchemaContent["test1"], &TestInput1{}, r)
	}))
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL, &b)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !validate {
		fmt.Println(result.SchemaErrors)
		t.FailNow()
	}
}
