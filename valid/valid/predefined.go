package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Validation regexps for type=...
var Regexps map[string]*regexp.Regexp

func init() {
	Fns = map[string]FnValidate{
		"len":    FnLen,
		"type":   FnType,
		"csv":    FnCsv,
		"reqif":  FnReqif,
		"onlyif": FnOnlyIf,
		"oneof":  FnOneOf,
		"def":    FnDef,
		"count":  FnCount,
		"hash":   FnHash,
		"eq":     FnEq,
	}
	initRegex()
}

func initRegex() {
	Regexps = map[string]*regexp.Regexp{
		"email": regexp.MustCompile(`.+@.+\..{2,}`),
	}
}

func FnEq(ctx Context, args FnArgs) bool {
	return false
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
			if cmp[i] >= 'a' && cmp[i] <= 'f' {
				continue
			} else if cmp[i] >= '0' && cmp[i] <= '9' {
				continue
			}
			return false
		}

		return true

	case "base64":
		for i := 0; i < len(cmp); i++ {
			if cmp[i] >= 'A' && cmp[i] <= 'Z' {
				continue
			} else if cmp[i] >= 'a' && cmp[i] <= 'z' {
				continue
			} else if cmp[i] >= '0' && cmp[i] <= '9' {
				continue
			} else if cmp[i] == '+' || cmp[i] == '/' || cmp[i] == '=' {
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
	if k != reflect.Slice && k != reflect.Map {
		panic(fmt.Sprintf("expected slice/map, got %s", k))
	}
	s := reflect.ValueOf(ctx.Value)
	if args["min"] == nil && args["max"] == nil {
		panic("no min/max given")
	}

	if args["min"] != nil {
		min, err := FnGetInt(args["min"])
		if err != nil {
			panic(err)
		}
		if s.Len() < int(min) {
			return false
		}
	}

	if args["max"] != nil {
		max, err := FnGetInt(args["max"])
		if err != nil {
			panic(err)
		}
		if s.Len() > int(max) {
			return false
		}
	}

	return true
}

func FnDef(ctx Context, args FnArgs) bool {
	t, _ := FnGetStr(args["type"])
	switch t {
	// Slugs are generally entirely lowercase, with accented characters replaced by
	// letters from the English alphabet and whitespace characters replaced by a dash or an underscore
	// http://en.wikipedia.org/wiki/Semantic_URL#Slug
	case "slug":
		cmp, ok := ctx.Value.(string)
		if !ok {
			return false
		}

		if len(cmp) == 0 {
			return false
		}
		for i := 0; i < len(cmp); i++ {
			if cmp[i] >= 'a' && cmp[i] <= 'z' {
				continue
			}
			if cmp[i] == '_' || cmp[i] == '-' {
				continue
			}

			return false
		}
		return true
	case "udecimal":
		cmp, ok := ctx.Value.(float64)
		if !ok {
			panic("expected float64")
		}

		if cmp >= 0 {
			return true
		}

		return false
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
		return false

	// Plain is regular ASCII input from Western keyboard
	case "plain":
		cmp, ok := ctx.Value.(string)
		if !ok {
			return false
		}

		if len(cmp) == 0 {
			return false
		}
		for i := 0; i < len(cmp); i++ {
			if cmp[i] >= ' ' && cmp[i] <= '~' {
				continue
			}

			return false
		}
		return true

	default:
		regx, ok := Regexps[t]
		if !ok {
			panic(fmt.Sprintf("type %s not implemented", t))
		}
		return regx.Match([]byte(ctx.Value.(string)))
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
				/* Todo: generic value reflect */
				if ctx.Value.(int64) > 0 {
					return true
				} else {
					return false
				}
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

	return false
}

func FnLen(ctx Context, args FnArgs) bool {
	cmp, ok := ctx.Value.(string)
	if !ok {
		panic("expected string")
	}
	if args["min"] == nil && args["max"] == nil {
		panic("no min/max given")
	}

	if args["min"] != nil {
		min, err := FnGetInt(args["min"])
		if err != nil {
			panic(err)
		}
		if len(cmp) < int(min) {
			return false
		}
	}

	if args["max"] != nil {
		max, err := FnGetInt(args["max"])
		if err != nil {
			panic(err)
		}
		if len(cmp) > int(max) {
			return false
		}
	}

	return true
}

func FnCsv(ctx Context, args FnArgs) bool {
	cmp, ok := ctx.Value.(string)
	if !ok {
		panic("expected string")
	}
	sep := ","
	if args["sep"] != nil {
		sep = args["sep"].(string)
	}

	allOk := true
	for _, s := range strings.Split(cmp, sep) {
		s = strings.TrimSpace(s)
		ok = FnDef(
			Context{Value: s},
			map[string]interface{}{"type": args["type"]},
		)
		if !ok {
			allOk = false
		}
	}
	return allOk
}
