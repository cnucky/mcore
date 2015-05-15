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
		"hash":   FnHash,
		"eq":     FnEq,
	}
}

func FnEq(ctx Context, args FnArgs) bool {
	return true
}

func FnHash(ctx Context, args FnArgs) bool {
	cmp, ok := ctx.Value.(string)
	if !ok {
		panic("expected string")
	}

	t, err := FnGetStr(args["type"])
	if err != nil {
		panic(err)
	}

	switch t {
	case "sha256":
		if len(cmp) != 64 {
			return false
		}

		for i := 0; i < len(cmp); i++ {
			if cmp[i] >= 'A' && cmp[i] <= 'F' {
				continue
			} else if cmp[i] >= '0' && cmp[i] <= '9' {
				continue
			}

			return false
		}

		return true
	default:
		panic(fmt.Sprintf("invalid hash %s", t))
	}
}

func FnCount(ctx Context, args FnArgs) bool {
	k := reflect.TypeOf(ctx.Value).Kind()
	if k != reflect.Slice {
		panic(fmt.Sprintf("expected slice, got %s", k))
	}
	s := reflect.ValueOf(ctx.Value)

	min, err := FnGetInt(args["min"])
	if err != nil {
		panic(err)
	}

	max, err := FnGetInt(args["max"])
	if err != nil {
		panic(err)
	}

	if s.Len() < int(min) {
		return false
	}

	if s.Len() > int(max) {
		return false
	}

	return true
}

func FnDef(ctx Context, args FnArgs) bool {
	t, _ := FnGetStr(args["type"])
	switch t {
	case "ascii":
		cmp, ok := ctx.Value.(string)
		if !ok {
			for i := 0; i < len(cmp); i++ {
				if cmp[i] >= 'A' && cmp[i] <= 'Z' {
					continue
				}

				return false
			}
		}

		return true
	case "uint":
		cmp, ok := ctx.Value.(int64)
		if !ok {
			panic("expected int64")
		}

		if cmp >= 0 {
			return true
		}

		return false
	case "date":
		return true
	case "email":
		return true
	default:
		panic(fmt.Sprintf("type %s not implemented", t))
	}

	return false
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
	cmp, ok := ctx.Value.(string)
	if !ok {
		panic("expected string")
	}

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
	cmp, ok := ctx.Value.(string)
	if !ok {
		panic("expected string")
	}

	min, err := FnGetInt(args["min"])
	if err != nil {
		panic(err)
	}

	max, err := FnGetInt(args["max"])
	if err != nil {
		panic(err)
	}

	if len(cmp) < int(min) {
		return false
	}

	if len(cmp) > int(max) {
		return false
	}

	return true
}
