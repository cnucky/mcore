package main

import (
	"github.com/xsnews/mcore/valid/valid"
	"reflect"
)

type TestParams struct {
	Code string `validate:"type(fmt=ascii),len(min=1,max=255)"`
	Dep  int    `validate:"reqif(Test=5)"`
	Test int    ``
}
type ComplexParam struct {
	Coupon string
	Test   []int
}

func main() {
	/* Initialize params */
	t := &TestParams{}
	t.Code = "Hello World"
	t.Test = 10

	/* Reflect struct */
	s := reflect.Indirect(reflect.ValueOf(t))
	for num := 0; num < s.NumField(); num++ {
		name := s.Type().Field(num).Name
		var value interface{}

		if name == "_" {
			/* Unexported field, generic struct rule */
			continue
		} else {
			/* Exported field, specific field rule */
			value = s.Field(num).Interface()
		}
		validateTag := s.Type().Field(num).Tag.Get("validate")
		if validateTag == "" {
			continue
		}
		//fmt.Println(value)

		//code.Write([]byte(fmt.Sprintf("init(ctx=%s),%s", name, validateTag)))
		//code.Write([]byte(validateTag))
		//err := l.Compile(code.String())
		//if code.Len() > 0 {

		if len(validateTag) > 0 {
			l := new(valid.Valdsl)
			err := l.Parse(t, validateTag, value)
			if err != nil {
				panic(err)
			}
		}
	}
}
