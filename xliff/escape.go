package xliff

import (
	"bytes"
	"encoding/xml"
	"html"
	"strings"
)

// XMLEscapeStrict will propertly XML escape the passed in text and
// return the escape text or an error. The result can be used e.g.
// as the value of an xml attribute value.
func XMLEscapeStrict(s string) (t string, err error) {

	buf := &bytes.Buffer{}
	err = xml.EscapeText(buf, []byte(s))
	if err != nil {
		return
	}
	t = buf.String()
	return
}

var xmlTextEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`<`, "&lt;",
)

// XMLEscapeLoose will propertly XML escape the passed in text and return the
// escape text. The loose escaping will only replace the bare minimum so the
// text can be included as the text of an element, but not as an attribute
// value.
func XMLEscapeLoose(s string) string {
	return xmlTextEscaper.Replace(s)
}

// XMLUnescape will remove the XML escaping from the passed in text and return
// the unescaped text.
func XMLUnescape(s string) string {
	return html.UnescapeString(s)
}
