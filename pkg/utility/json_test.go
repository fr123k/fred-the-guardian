package utility

import (
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
