package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xsnews/microservice-core/gosanitize/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttpRequestForm(t *testing.T) {
	var validate bool = true

	TestValues1 := &url.Values{
		"Code":  {"hello world"},
		"Int":   {"10"},
		"Float": {"1.5"},
		"Bool":  {"true"},
		"Email": {"test@gmail.com"},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* Create a validator object called test1 with a json schema */
		v := util.NewValidator("test1", LoadSchema("./schemas/test1.json"), &TestInput1{})

		/* Load test values into validator object */
		if err := v.LoadValuesFromRequest(r); err != nil {
			fmt.Println("LoadValuesFromRequest error:", err)
			validate = false
			return
		}

		/* Validate against JSON schema */
		if ok, err := v.Validate(); !ok {
			fmt.Println("Validate:", err)
			validate = false
			return
		}

		/* Validate against custom rules */
		if ok, err := v.ValidateRules(); !ok {
			fmt.Println("Validate:", err)
			validate = false
			return
		}
	}))
	defer ts.Close()

	res, err := http.PostForm(ts.URL, *TestValues1)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !validate {
		t.FailNow()
	}
}

func TestHttpRequestJson(t *testing.T) {
	var validate bool = true

	TestValues1 := &TestInput1{
		Code:  "hello world",
		Int:   10,
		Float: 1.5,
		Bool:  true,
		Email: "test@gmail.com",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* Create a validator object called test1 with a json schema */
		v := util.NewValidator("test1", LoadSchema("./schemas/test1.json"), &TestInput1{})

		/* Load test values into validator object */
		if err := v.LoadValuesFromRequest(r); err != nil {
			fmt.Println("LoadValuesFromRequest error:", err)
			validate = false
			return
		}

		/* Validate against JSON schema */
		if ok, err := v.Validate(); !ok {
			fmt.Println("Validate:", err)
			validate = false
			return
		}

		/* Validate against custom rules */
		if ok, err := v.ValidateRules(); !ok {
			fmt.Println("Validate:", err)
			validate = false
			return
		}
	}))
	defer ts.Close()

	var b bytes.Buffer
	j := json.NewEncoder(&b)
	j.Encode(TestValues1)

	req, err := http.NewRequest("POST", ts.URL, &b)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !validate {
		t.FailNow()
	}
}
