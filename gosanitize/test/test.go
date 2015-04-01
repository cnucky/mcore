package main

import (
	"errors"
	"fmt"
	"goutils/gosanitize/rule"
	"goutils/gosanitize/util"
	"goutils/gosanitize/validate"
	"io/ioutil"
	"net"
	"strings"
)

/* Struct of all input keys */
type TestInput1 struct {
	Code        string    `json:",omitempty"`
	TestDep     string    `json:",omitempty"`
	Conditional string    `json:",omitempty"`
	Int         int       `json:",omitempty"`
	Bool        bool      `json:",omitempty"`
	Array       []int     `json:",omitempty"`
	Email       string    `json:",omitempty" rule:"RuleLookupMX"` // Rule for extended validation
	_           rule.Rule `json:",omitempty" rule:"RuleGeneric"`  // Generic validation rule not specific to a field
}

/* Just a rule test, we pass the entire struct in case a rule depends on more fields */
func RuleLookupMX(obj interface{}) error {
	/* Get domain from email address */
	domain := strings.SplitN(obj.(string), "@", 2)
	if len(domain) != 2 {
		return errors.New("Invalid email address")
	}

	/* Lookup MX record first */
	_, errMx := net.LookupMX(domain[1])
	if errMx != nil {
		/* Fallback on A record */
		_, errHost := net.LookupHost(domain[1])
		if errHost != nil {
			return errHost
		}

		return errMx
	}

	return nil
}

func RuleGeneric(obj interface{}) error {
	/* Cast interface to struct */
	params := obj.(*TestInput1)

	fmt.Println("RuleGeneric OK", params.Array)

	return nil
}

func init() {
	rule.AddRule("RuleLookupMX", RuleLookupMX)
	rule.AddRule("RuleGeneric", RuleGeneric)
}

func ValidateTest(id string, filename string, params interface{}, values *map[string]string) bool {
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
		fmt.Println(err)
		return false
	}

	/* Validate input values with the json schema */
	fmt.Println("Input:", params)
	ok, result := v.Validate()
	if ok {
		fmt.Println("InputValidator = OK")
	} else {
		fmt.Println("InputValidator = ERR")
		fmt.Println("Details:", result.Errors())
		fmt.Println()
		return false
	}

	r := rule.NewValidator(id, params)
	ok, errors := r.Validate()
	if ok {
		fmt.Println("RuleValidator = OK")
	} else {
		fmt.Println("RuleValidator = ERR")
		fmt.Println("Details:", errors)
		fmt.Println()
		return false
	}

	return true
}

func main() {
	/* Test validation, all OK */
	TestValues1 := map[string]string{
		"Code":  "Hello World",
		"Int":   "10",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	validate := ValidateTest("test1", "./examples/1.json", &TestInput1{}, &TestValues1)
	if validate {
		fmt.Println("PASS\n")
	}

	/* Test dependency */
	TestValues1 = map[string]string{
		"Code":    "Hello World",
		"TestDep": "wat",
		"Int":     "10",
		"Bool":    "true",
		"Email":   "test@gmail.com",
	}

	validate = ValidateTest("test1", "./examples/1.json", &TestInput1{}, &TestValues1)
	if validate {
		fmt.Println("PASS\n")
	}

	/* Test validation, missing arg */
	TestValues1 = map[string]string{
		"Int":   "10",
		"Bool":  "true",
		"Email": "test@gmail.com",
	}

	validate = ValidateTest("test1", "./examples/1.json", &TestInput1{}, &TestValues1)
	if validate {
		fmt.Println("PASS\n")
	}

	/* Test validation, pattern incorrect */
	TestValues1 = map[string]string{
		"Code":  "Hello World",
		"Int":   "10",
		"Bool":  "true",
		"Email": "t est@gmail.com",
	}

	validate = ValidateTest("test1", "./examples/1.json", &TestInput1{}, &TestValues1)
	if validate {
		fmt.Println("PASS\n")
	}

	/* Test validation, LookupMX fails */
	TestValues1 = map[string]string{
		"Code":  "Hello World",
		"Int":   "10",
		"Bool":  "true",
		"Email": "test@no.such.domain",
	}

	validate = ValidateTest("test1", "./examples/1.json", &TestInput1{}, &TestValues1)
	if validate {
		fmt.Println("PASS\n")
	}
}
