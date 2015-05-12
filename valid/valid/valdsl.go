package valid

import (
	"fmt"
)

//"errors"

/* call context */
type Valdsl struct {
	Ctx interface{}
}

type FnArgs map[string]interface{}         /* validation arguments */
type FnValue interface{}                   /* validation value */
type FnValidate func(Context, FnArgs) bool /* validation function definition */
var Fns map[string]FnValidate              /* validation id -> func map */

/* Parse next rule */
func (v *Valdsl) Next(tokens []*Token, ctx *Context) (int, bool, error) {
	var state int                  /* current parsing state, 0 = parse symbol, 1 = parse arguments */
	var fn string                  /* function id */
	var args FnArgs = make(FnArgs) /* collected arguments */

	/* iterate through tokens */
	var c int
	var done bool
	for c = 0; c < len(tokens); c++ {
		switch state {
		case 0: /* Set function name */
			if tokens[c].Type != tokenSymbol {
				panic(fmt.Sprintf("expected symbol, found %s", tokens[c].Id))
				break
			}

			fn = tokens[c].Id

			/* Peek at next token */
			j := c + 1
			if j >= len(tokens) {
				/* Exit if we're at the end already */
				done = true
				break
			}

			if tokens[j].Type == tokenSep {
				done = true
				break
			} else if tokens[j].Type == tokenArgOpen {
				state = 1
				c++
			} else {
				panic(fmt.Sprintf("invalid rule, expected , or ( found %s", tokens[j].Id))
			}
		case 1:
			if tokens[c].Type != tokenSymbol {
				/* Done parsing arguments */
				done = true
				break
			}

			/* Set argument vars and increase cursor */
			argname := tokens[c].Id
			_ = tokens[c+1].Id /* currently this is always = */
			argtype := tokens[c+2].Type
			argval := tokens[c+2].Id
			c = c + 3

			/* What kind of arguments are we passing */
			if argtype == tokenArgSliceOpen {
				/* Parse a slice of symbols */
				var j int
				var s []string = make([]string, 0, 4)
				for j = c; j < len(tokens); j++ {
					if tokens[j].Type == tokenArgSliceClose {
						break
					} else if tokens[j].Type == tokenSep {
						continue
					} else if tokens[j].Type != tokenSymbol {
						panic(fmt.Sprintf("expected symbol found %s", tokens[j].Id))
					}

					s = append(s, tokens[j].Id)
				}

				args[argname] = s

				c = j + 1
			} else {
				/* Single symbol */
				args[argname] = argval
			}
		}

		/* Are we done parsing? */
		if done {
			break
		}
	}

	fmt.Println(fn, args)

	if Fns[fn] == nil {
		fmt.Println(fn, "not implemented")
		return c, false, nil
	}

	ret := Fns[fn](*ctx, args)
	fmt.Println(fn, ret)

	return c, ret, nil
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
