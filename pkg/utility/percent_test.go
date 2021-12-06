package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeventyFivePercent(t *testing.T) {
	tests := map[string]struct {
		input  UInt
		expect uint
	}{
		"of 0":     {input: UInt(0), expect: 0},
		"of 1":     {input: UInt(1), expect: 1},
		"of 2":     {input: UInt(2), expect: 2},
		"of 4":     {input: UInt(4), expect: 3},
		"of 10":    {input: UInt(10), expect: 8},
		"of 133":   {input: UInt(133), expect: 100},
		"of 15262": {input: UInt(15262), expect: 11447},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.SeventyFivePercent()
			assert.Equal(t, int(tc.expect), int(got), name)
		})
	}
}
