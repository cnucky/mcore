package main

import (
	"github.com/xsnews/microservice-core/gosanitize/rule"
	"github.com/xsnews/microservice-core/gosanitize/util"
	"testing"
)

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

func BenchmarkValidate(b *testing.B) {
	TestValues1 := util.Values{
		"Code":  "Hello World",
		"Int":   "10",
		"Float": "1.5",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	for i := 0; i < b.N; i++ {
		ok := simpleValidate("test1", LoadSchema("./schemas/test1.json"), &TestValues1, &TestInput1{}, false)
		if !ok {
			b.FailNow()
		}
	}
}

func TestInputIsValidOK(t *testing.T) {
	/* Our test values */
	TestValues1 := util.Values{
		"Code":  "Hello World",
		"Int":   "10",
		"Float": "1.5",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	ok := simpleValidate("test1", LoadSchema("./schemas/test1.json"), &TestValues1, &TestInput1{}, false)
	if !ok {
		t.FailNow()
	}
}

func TestDependencyOK(t *testing.T) {
	TestValues1 := util.Values{
		"Code":        "Hello World",
		"TestDep":     "wat",
		"Conditional": "ok",
		"Int":         "10",
		"Float":       "1.5",
		"Bool":        "true",
		"Email":       "test@gmail.com",
	}

	ok := simpleValidate("test1", LoadSchema("./schemas/test1.json"), &TestValues1, &TestInput1{}, false)
	if !ok {
		t.FailNow()
	}
}

func TestDependencyFail(t *testing.T) {
	TestValues1 := util.Values{
		"Code":    "Hello World",
		"TestDep": "wat",
		"Int":     "10",
		"Float":   "1.5",
		"Bool":    "true",
		"Email":   "test@gmail.com",
	}

	ok := simpleValidate("test1", LoadSchema("./schemas/test1.json"), &TestValues1, &TestInput1{}, true)
	if !ok {
		t.FailNow()
	}
}

func TestMissingArgFail(t *testing.T) {
	TestValues1 := util.Values{
		"Int":   "10",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	ok := simpleValidate("test1", LoadSchema("./schemas/test1.json"), &TestValues1, &TestInput1{}, true)
	if !ok {
		t.FailNow()
	}
}

func TestRegexMatchFail(t *testing.T) {
	TestValues1 := util.Values{
		"Code":  "Hello World",
		"Int":   "10",
		"Float": "1.5",
		"Bool":  "true",
		"Email": "t est@gmail.com",
	}

	ok := simpleValidate("test1", LoadSchema("./schemas/test1.json"), &TestValues1, &TestInput1{}, true)
	if !ok {
		t.FailNow()
	}
}

func TestFieldRuleFail(t *testing.T) {
	TestValues1 := util.Values{
		"Code":  "Hello World",
		"Int":   "10",
		"Float": "1.5",
		"Bool":  "true",
		"Email": "test@no.such.domain",
	}

	ok := simpleValidate("test1", LoadSchema("./schemas/test1.json"), &TestValues1, &TestInput1{}, true)
	if !ok {
		t.FailNow()
	}
}

type TestInput2 struct {
	Array []int `json:",omitempty"`
}

func TestArray(t *testing.T) {
	TestValues2 := util.Values{
		"Array": []string{"10", "20"},
	}

	ok := simpleValidate("test2", LoadSchema("./schemas/test2.json"), &TestValues2, &TestInput2{}, false)
	if !ok {
		t.FailNow()
	}
}
