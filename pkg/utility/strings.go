package utility

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

func TrailingSlash(s string) string {
    if strings.HasSuffix(s, "/") {
        return s
    }
    return fmt.Sprintf("%s/", s)
}

//https://stackoverflow.com/a/50581165
func RandomString(len uint) string {
    buff := make([]byte, len)
    rand.Read(buff)
    str := base64.StdEncoding.EncodeToString(buff)
    // Base 64 can be longer than len
    return str[:len]
}
