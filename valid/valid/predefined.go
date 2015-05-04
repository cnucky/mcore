package valid

import (
	"fmt"
)

func init() {
	Fns = map[string]FnValidate{
		"len":   FnLen,
		"type":  FnType,
		"reqif": FnReqif,
	}
}

func FnReqif(ctx Context, args FnArgs) bool {
	//fld := reflect.ValueOf(ctx.Ctx)
	//fmt.Println(fld.FieldByName("Test"))
	for k, v := range args {
		fmt.Println(k, v)
	}

	return true
}

func FnType(ctx Context, args FnArgs) bool {
	fmt, _ := FnGetStr(args["fmt"])
	if fmt == "" {
		panic("fmt: missing")
	}

	return true
}

func FnLen(ctx Context, args FnArgs) bool {
	min, err := FnGetInt(args["min"])
	if err != nil {
		panic(err)
	}

	max, err := FnGetInt(args["min"])
	if err != nil {
		panic(err)
	}

	if len(ctx.Value.(string)) < int(min) {
		return false
	}

	if len(ctx.Value.(string)) > int(max) {
		return false
	}

	return true
}
