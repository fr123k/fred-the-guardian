package utility

import (
    "reflect"
    "strings"
)

// In case of invalid struct it returns the field name from the json tag instead of the struct variable name.
func JsonTagName(fld reflect.StructField) string {
    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

    if name == "-" {
        return ""
    }

    return name
}
