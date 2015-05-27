package valid

import (
	"testing"
)

// TODO: FnEq
// TODO: date

func TestHashBase64(t *testing.T) {
	var c string
	checks := map[string]bool{
		`%invalid`: false,
		`dGVzdA==`: true,
	}
	for val, expect := range checks {
		if valid := FnHash(Context{Ctx: &c, Value: val}, map[string]interface{}{
			"type": "base64",
		}); valid != expect {
			t.Errorf("Base64 check failed, val=%s expect=%t", val, expect)
		}
	}
}

func TestHashSha256(t *testing.T) {
	var c string
	checks := map[string]bool{
		`%invalid`: false,
		`9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08`: true,
	}
	for val, expect := range checks {
		if valid := FnHash(Context{Ctx: &c, Value: val}, map[string]interface{}{
			"type": "sha256",
		}); valid != expect {
			t.Errorf("Sha256 check failed, val=%s expect=%t", val, expect)
		}
	}
}

func TestHashPanic(t *testing.T) {
	defer func() {
		recover()
    }()

	var c string
	FnHash(Context{Ctx: &c, Value: ``}, map[string]interface{}{
		"type": "unsupported",
	})
	t.Error("No panic as it should for hash(type=unsupported)")
}

func TestCount(t *testing.T) {
	var c string
	// min
	if valid := FnCount(Context{Ctx: &c, Value: map[string]string{}}, map[string]interface{}{
		"min": "5",
	}); valid {
		t.Error("Should not accept empty map for count(min=5)")
	}
	if valid := FnCount(Context{Ctx: &c, Value: map[string]string{}}, map[string]interface{}{
		"min": "0",
	}); !valid {
		t.Error("Should accept empty map for count(min=0)")
	}
	if valid := FnCount(Context{Ctx: &c, Value: map[string]string{"a":"a"}}, map[string]interface{}{
		"min": "1",
	}); !valid {
		t.Error("Should accept map for count(min=1)")
	}

	// max
	if valid := FnCount(Context{Ctx: &c, Value: map[string]string{}}, map[string]interface{}{
		"max": "0",
	}); !valid {
		t.Error("Should accept empty map for count(max=0)")
	}
	if valid := FnCount(Context{Ctx: &c, Value: map[string]string{"a":"a"}}, map[string]interface{}{
		"max": "0",
	}); valid {
		t.Error("Should not accept map (with one item) for count(max=0)")
	}
	if valid := FnCount(Context{Ctx: &c, Value: map[string]string{"a":"a", "b":"b"}}, map[string]interface{}{
		"max": "1",
	}); valid {
		t.Error("Should not accept map (with two items) for count(max=1)")
	}

	// Array
	if valid := FnCount(Context{Ctx: &c, Value: []string{"a"}}, map[string]interface{}{
		"min": "1", "max": "1",
	}); !valid {
		t.Error("Should accept array (with one item) for count(min=1,max=1)")
	}
}

func TestCountPanic(t *testing.T) {
	defer func() {
		recover()
    }()

	var c string
	FnCount(Context{Ctx: &c, Value: ""}, map[string]interface{}{})
	t.Error("Should receive panic for missing min/max")
}

func TestSlug(t *testing.T) {
	var c string

	checks := map[string]bool{
		` invalid`: false,
		``: false,
		`standard-slug`: true,
		`standard_slug`: true,
		`PartOfUrl`: false,
		`900`: false,
	}
	for val, expect := range checks {
		if valid := FnDef(Context{Ctx: &c, Value: val}, map[string]interface{}{
			"type": "slug",
		}); valid != expect {
			t.Errorf("Slug check failed, val=%s expect=%t", val, expect)
		}
	}
}

func TestUdecimal(t *testing.T) {
	var c float64

	checks := map[float64]bool{
		float64(-10): false,
		float64(-10.2): false,
		float64(10.2): true,
		float64(900): true,
	}
	for val, expect := range checks {
		if valid := FnDef(Context{Ctx: &c, Value: val}, map[string]interface{}{
			"type": "udecimal",
		}); valid != expect {
			t.Errorf("Udecimal check failed, val=%f expect=%t", val, expect)
		}
	}
}

func TestUint(t *testing.T) {
	var c int64

	checks := map[int64]bool{
		int64(-10): false,
		int64(-10000): false,
		int64(10): true,
		int64(900): true,
	}
	for val, expect := range checks {
		if valid := FnDef(Context{Ctx: &c, Value: val}, map[string]interface{}{
			"type": "uint",
		}); valid != expect {
			t.Errorf("Uint check failed, val=%d expect=%t", val, expect)
		}
	}
}

func TestPlain(t *testing.T) {
	var c string

	checks := map[string]bool{
		` valid `: true,
		`helloworld`: true,
		`HelloWorld`: true,
		`~<>?,./:"|;'\{}[]!@#$%^&*()_+-=`: true,
		`1234567890`: true,
		``: false,
		"\r\n": false,

		// unicode
		"“ ” ‘ ’ – — … ‐ ‒ ° © ® ™ • ½ ¼ ¾ ⅓ ⅔ † ‡ µ ¢ £ € « » ♠ ♣ ♥ ♦ ¿ �": false,
		"僀 僁 僂 僃 僄 僅 僆 僇 僈 僉 僊 僋 僌 働 僎 像 僐 僑 僒 僓 僔 僕 僖 僗 僘 僙 僚 僛 僜 僝 僞 僟": false,
	}
	for val, expect := range checks {
		if valid := FnDef(Context{Ctx: &c, Value: val}, map[string]interface{}{
			"type": "plain",
		}); valid != expect {
			t.Errorf("Plain check failed, val=%s expect=%t", val, expect)
		}
	}
}

func TestEmail(t *testing.T) {
	var c string

	checks := map[string]bool{
		``: false,
		`mark@xsnews.nl`: true,
		`invalid@.`: false,
		`@xsnews.nl`: false,
		`a@b.c`: false,
		`a.@b.cc`: true,
		`fname.lname@test.com`: true,
		`person@newtld.amsterdam`: true,
	}
	for val, expect := range checks {
		if valid := FnDef(Context{Ctx: &c, Value: val}, map[string]interface{}{
			"type": "email",
		}); valid != expect {
			t.Errorf("Email check failed, val=%s expect=%t", val, expect)
		}
	}
}

func TestLen(t *testing.T) {
	var c string
	// min
	if valid := FnLen(Context{Ctx: &c, Value: ""}, map[string]interface{}{
		"min": "5",
	}); valid {
		t.Error("Should not accept empty map for count(min=5)")
	}
	if valid := FnLen(Context{Ctx: &c, Value: ""}, map[string]interface{}{
		"min": "0",
	}); !valid {
		t.Error("Should accept empty map for count(min=0)")
	}
	if valid := FnLen(Context{Ctx: &c, Value: "a"}, map[string]interface{}{
		"min": "1",
	}); !valid {
		t.Error("Should accept map for count(min=1)")
	}

	// max
	if valid := FnLen(Context{Ctx: &c, Value: ""}, map[string]interface{}{
		"max": "0",
	}); !valid {
		t.Error("Should accept empty map for count(max=0)")
	}
	if valid := FnLen(Context{Ctx: &c, Value: "abc"}, map[string]interface{}{
		"max": "0",
	}); valid {
		t.Error("Should not accept map (with one item) for count(max=0)")
	}
	if valid := FnLen(Context{Ctx: &c, Value: "abcd"}, map[string]interface{}{
		"max": "1",
	}); valid {
		t.Error("Should not accept map (with two items) for count(max=1)")
	}

	// Array
	if valid := FnLen(Context{Ctx: &c, Value: "a"}, map[string]interface{}{
		"min": "1", "max": "1",
	}); !valid {
		t.Error("Should accept array (with one item) for count(min=1,max=1)")
	}
}

func TestLenPanic(t *testing.T) {
	defer func() {
		recover()
    }()

	var c string
	FnLen(Context{Ctx: &c, Value: ""}, map[string]interface{}{})
	t.Error("Should receive panic for missing min/max")
}

func TestCsv(t *testing.T) {
	var c string
	if valid := FnCsv(Context{Ctx: &c, Value: "a,b,c"}, map[string]interface{}{
		"sep": ",",
		"type": "plain",
	}); !valid {
		t.Error("csv(plain) failed?")
	}
	if valid := FnCsv(Context{Ctx: &c, Value: "a\n€,b,c"}, map[string]interface{}{
		"sep": ",",
		"type": "plain",
	}); valid {
		t.Error("csv(plain) should failed with newline char")
	}
}