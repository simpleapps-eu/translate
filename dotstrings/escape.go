package dotstrings

import (
	"strconv"
	"strings"
)

func StringsUnescape(s string) (t string, err error) {
	return strconv.Unquote(`"` + s + `"`)
}

func StringsEscape(s string) string {
	return strings.Trim(strconv.Quote(s), `"`)
}
