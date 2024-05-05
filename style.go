package svger

import (
	"fmt"
	"strconv"
	"strings"
)

// refString returns a string reference.
func refString(s string) *string {
	x := fmt.Sprint(s)
	return &x
}

// parse a numerical width value. Since 0 is a legitimate width,
// return -1 in case of error.
func parseWidth(val string) float64 {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return -1
	}
	return f
}

// splitStyle unpacks the style string from a path element into a key
// value map.
func splitStyle(style string) map[string]string {
	var r map[string]string
	r = make(map[string]string)
	props := strings.Split(style, ";")

	for _, keyval := range props {
		kv := strings.Split(strings.TrimSpace(keyval), ":")
		if len(kv) >= 2 {
			r[kv[0]] = kv[1]
		}
	}

	return r
}
