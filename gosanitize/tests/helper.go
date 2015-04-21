package main

import (
	"errors"
	"fmt"
	"github.com/xsnews/microservice-core/gosanitize/rule"
	"github.com/xsnews/microservice-core/gosanitize/util"
	"io/ioutil"
)

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

func LoadSchema(filename string) []byte {
	schema, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return schema
}

func simpleValidate(id string, schema []byte, values *util.Values, input interface{}, expectFail bool) bool {
	/* Create a validator object called test1 with a json schema */
	v := util.NewValidator("test1", schema, input)

	/* Load test values into validator object */
	if err := v.LoadValues(*values); err != nil {
		fmt.Println("LoadValues error:", err)
		return false
	}

	/* Validate against JSON schema */
	if ok, err := v.Validate(); !ok {
		fmt.Println("Validate:", err)
		if !expectFail {
			return false
		}
	}

	/* Validate against custom rules */
	if ok, err := v.ValidateRules(); !ok {
		fmt.Println("Validate:", err)
		if !expectFail {
			return false
		}
	}

	return true
}
