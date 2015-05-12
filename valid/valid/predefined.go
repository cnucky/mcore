package valid

import (
	"fmt"
	"reflect"
)

func init() {
	Fns = map[string]FnValidate{
		"len":    FnLen,
		"type":   FnType,
		"reqif":  FnReqif,
		"onlyif": FnOnlyIf,
		"oneof":  FnOneOf,
		"def":    FnDef,
		"count":  FnCount,
	}
}

func FnCount(ctx Context, args FnArgs) bool {
	return true
}

func FnDef(ctx Context, args FnArgs) bool {
	return true
}

func FnOnlyIf(ctx Context, args FnArgs) bool {
	s := reflect.ValueOf(ctx.Ctx).Elem()
	for k, cmp := range args {
		fld := s.FieldByName(k)
		if !fld.IsValid() {
			panic("field missing")
		}

		cmp2 := fmt.Sprintf("%v", fld.Interface())

		enum, _ := FnGetStrSlice(cmp)
		for _, cmp := range enum {
			if cmp == cmp2 {
				return true
			}
		}
	}

	return false
}

func FnOneOf(ctx Context, args FnArgs) bool {
	cmp := ctx.Value.(string)
	enum, _ := FnGetStrSlice(args["enum"])
	for _, cmp2 := range enum {
		if cmp == cmp2 {
			return true
		}
	}

	return false
}

func FnReqif(ctx Context, args FnArgs) bool {
	s := reflect.ValueOf(ctx.Ctx).Elem()
	for k, cmp := range args {
		fld := s.FieldByName(k)
		if !fld.IsValid() {
			panic("field missing")
		}

		cmp2 := fmt.Sprintf("%v", fld.Interface())
		fmt.Println(cmp2, cmp)
		break
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

	max, err := FnGetInt(args["max"])
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
