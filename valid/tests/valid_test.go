package valid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xsnews/mcore/valid/valid"
	"testing"
)

type InputLine struct {
	Price    float64 `validate:"def(type=udecimal)"`         // allow: 5,2 deny: -5,2 5.2
	Quantity int64   `validate:"def(type=uint)"`             // Allow: 4 2 10000 deny: -4 5.3 5,3
	Comment  *string `validate:"opt,strlen(min=10,max=255)"` // optional but if set check strlen
	Barcode  *string `validate:"opt,barcode"`                // Barcode is domain validator (project X only)
	Type     string  `validate:"match(val=[DEBIT,CREDIT])"`  // Value must be DEBIT OR CREDIT
}

type Input struct {
	Id    int64  `validate:"def(type=uint)"`
	Email string `validate:"def(type=email),len(min=10,max=255)"`
	Pass  string `validate:"eq(field=Pass1)"` // Equal value+type like other field
	Pass1 string `validate:"eq(field=Pass)"`
	Hash  string `validate:"hash(type=sha256)"` // panic on invalid type
	Role  string `validate:"oneof(enum=[admin,user])"`

	Customer string `validate:"def(type=slug)"`                                    // Only allow [a-zA-Z]+
	Supplier int64  `validate:"def(type=uint),onlyif(Customer=[OPTION1,OPTION2])"` // supplier must be > 0 if customer field

	Map map[string]string `validate:count(min=1,max=1)`     // require at least 1 map and most 1
	Lines []InputLine `validate:"count(min=2,max=4)"` // require at least 2 lines and most 4
	Line  InputLine   `validate:"def(type=slug)"`     // require at least 2 lines and most 4
	Date  string      `validate:"def(type=date)"`     // date ALWAYS needs a date in a format WITH a timezone!

	Added  *int64 `validate:"def(type=uint)"`
	Added2 *int64 `validate:"def(type=uint)"`
}

func data() *Input {
	s := &Input{}
	s.Email = "hello@world.test"
	s.Role = "admin"
	s.Customer = "OPTION1"
	s.Hash = "AF93BCDEAFAF93BCDEAFAF93BCDEAFAF93BCDEAFAF93BCDEAFAF93BCDEAFFDAA"
	s.Lines = []InputLine{InputLine{Price: 10}, InputLine{Price: 20}}
	s.Line = InputLine{Price: 10}
	s.Map = map[string]string{"KEY": "VALUE"}
	var a int64 = -10
	s.Added = &a

	return s
}

func TestParser(t *testing.T) {
	s := data()
	ok, msg := valid.Validate(s)
	if !ok {
		t.Error("Validation failed: " + fmt.Sprintf("%+v", msg))
	}
}

func TestJson(t *testing.T) {
	s := data()
	var buf bytes.Buffer
	if e := json.NewEncoder(&buf).Encode(&s); e != nil {
		panic(e)
	}
	// TODO: Finish
}