package utility

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestTrailingSlash(t *testing.T) {
    tests := map[string]struct {
        input  string
        expect string
    }{
        "empty":                  {input: "", expect: "/"},
        "simple missing slash":   {input: "abc", expect: "abc/"},
        "missing trailing slash": {input: "a/b/c", expect: "a/b/c/"},
        "trailing":               {input: "a/b/c/", expect: "a/b/c/"},
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            got := TrailingSlash(tc.input)
            assert.Equal(t, tc.expect, got, name)
        })
    }
}

func TestRandomString(t *testing.T) {
    tests := map[string]struct {
        input  uint
        expect uint
    }{
        "missing trailing slash": {input: 0, expect: 0},
        "empty":                  {input: 8, expect: 8},
        "trailing":               {input: 64, expect: 64},
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            got := RandomString(tc.input)
            assert.Equal(t, tc.expect, uint(len(got)), name)
        })
    }
}
