package utility

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type JsonExample struct {
	Print string
}

func TestPrettyPrint(t *testing.T) {
	str := PrettyPrint(JsonExample{Print: "foo"})
	assert.Equal(t, "{\n\t\"Print\": \"foo\"\n}", str, "")
}


func TestJsonTagName(t *testing.T) {
	name := JsonTagName(reflect.StructField{
		Tag: `json:"request" validate:"required"`,
	})
	assert.Equal(t, "request", name, "")
}

func TestJsonTagNameNoJsonTag(t *testing.T) {
	name := JsonTagName(reflect.StructField{
		Tag: `validate:"required"`,
	})
	assert.Equal(t, "", name, "")
}

func TestJsonTagNameNoTags(t *testing.T) {
	name := JsonTagName(reflect.StructField{})
	assert.Equal(t, "", name, "")
}
