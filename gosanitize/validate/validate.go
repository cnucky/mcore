package validate

import (
	"bytes"
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
	"text/template"
)

type InputValidator struct {
	Id            string
	SchemaContent []byte

	Params interface{}
}

/* Convert Field values to Json */
func (i *InputValidator) Json() string {
	w := new(bytes.Buffer)

	enc := json.NewEncoder(w)
	enc.Encode(i.Params)

	return w.String()
}

/* Validate params */
func (i *InputValidator) Validate() (bool, *gojsonschema.Result) {
	t := template.Must(template.New(i.Id).Parse(string(i.SchemaContent)))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, Patterns)
	if err != nil {
		panic(err)
	}

	schemaLoader := gojsonschema.NewStringLoader(buf.String())
	documentLoader := gojsonschema.NewStringLoader(i.Json())

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	return result.Valid(), result
}

func NewValidator(id string, schema []byte, params interface{}) *InputValidator {
	return &InputValidator{
		Id:            id,
		Params:        params,
		SchemaContent: schema,
	}
}
