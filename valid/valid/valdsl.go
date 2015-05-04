package valid

import (
	"errors"
	"fmt"
)

/* call context */
type Valdsl struct {
	Ctx interface{}
}

type FnArgs map[string]string              /* validation arguments */
type FnValue interface{}                   /* validation value */
type FnValidate func(Context, FnArgs) bool /* validation function definition */
var Fns map[string]FnValidate              /* validation id -> func map */

/* parses the next available token */
func (v *Valdsl) Next(tokens []*Token, ctx *Context) (int, bool, error) {
	var state int                  /* current parsing state, 0 = parse symbol, 1 = parse arguments */
	var fn string                  /* function id */
	var args FnArgs = make(FnArgs) /* collected arguments */

	/* iterate through tokens */
	var i int
	for i = 0; i < len(tokens); i++ {
		if tokens[i].Type == tokenArgOpen {
			continue
		} else if tokens[i].Type == tokenSep {
			continue
		} else if tokens[i].Type == tokenArgClose {
			/* we're done parsing */
			break
		}

		/* set validation function */
		if tokens[i].Type == tokenSymbol && state == 0 {
			fn = tokens[i].Id
			state = 1
			continue
		}

		/* todo: replace this with a grammar checker */
		if i+2 >= len(tokens) {
			return -1, false, errors.New("Invalid rule")
		}

		/* Parse argument */
		id, _, val := tokens[i].Id, tokens[i+1].Id, tokens[i+2].Id
		i = i + 2
		args[id] = val
		//fmt.Println("Rule:", id, op, val)
	}

	if Fns[fn] == nil {
		fmt.Println(fn, "not implemented")
		return i, false, nil
	}

	ret := Fns[fn](*ctx, args)
	fmt.Println(fn, ret)

	return i, ret, nil
}

func (v *Valdsl) Parse(c interface{}, code string, value FnValue) error {
	l := &Lexer{}
	ctx := &Context{Ctx: c, Value: value}
	tokens := l.Tokenize(code)
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == tokenSep {
			//fmt.Println("NEXT!")
			continue
		}

		skip, _, err := v.Next(tokens[i:], ctx)
		if err != nil {
			//fmt.Println(err)
			return err
		}

		i += skip
	}

	return nil
}
