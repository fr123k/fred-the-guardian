package utility

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
