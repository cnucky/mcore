// Relations contains field compare functions.
package valid
import (
	"fmt"
	"reflect"
)
func FnEq(ctx Context, args FnArgs) bool {
	return false
}

func FnOnlyIf(ctx Context, args FnArgs) bool {
	s := reflect.ValueOf(ctx.Ctx) //.Elem()
	for k, cmp := range args {
		fld := s.FieldByName(k)
		if !fld.IsValid() {
			panic("field missing")
		}

		cmp2 := fmt.Sprintf("%v", fld.Interface())

		enum, _ := FnGetStrSlice(cmp)
		for _, cmp := range enum {
			if cmp == cmp2 {
				//if ctx.Value.(int64) > 0 {
					return true
				//} else {
				//	return false
				//}
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