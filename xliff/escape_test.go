package xliff

import (
	"strconv"
	"strings"
	"testing"
)

var tvXMLEscaped = [...]string{
	`&amp;`, //0
	`&lt;`,  //1
	`&gt;`,  //2
	`&#39; apostrof needs to be encoded as an entity for attributes using apostrofs`, //3
	`&#34; quote needs to be encoded as an entity for attributes using quotes`,       //4
	`ha&amp;llo`, //5
	`hallo&#xA;`, //6
	`&#x9;hallo`, //7
	`\test`,      //8
	`\nest`,      //9
	`hallo&#xA;dit&#xA;is&#x9;een&#xA;test\test&#34;test&#39;test`,              //10
	`ha&amp;llo&#xA;d&lt;it&#xA;i&gt;s&#x9;een&#xA;test\test&#34;test&#39;test`, //11
}

var tvXMLUnescaped = [...]string{
	`&`, //0
	`<`, //1
	`>`, //2
	`' apostrof needs to be encoded as an entity for attributes using apostrofs`, //3
	`" quote needs to be encoded as an entity for attributes using quotes`,       //4
	`ha&llo`, //5
	`hallo
`, //6
	`	hallo`, //7
	`\test`, //8
	`\nest`, //9
	`hallo
dit
is	een
test\test"test'test`, //10
	`ha&llo
d<it
i>s	een
test\test"test'test`, //11
}

// Some characters can be escaped by multipel entities.
var tvXMLMultiEscapeTest = [...]string{
	`&#39;`,  //0
	`&#34;`,  //1
	`&apos;`, //2
	`&quot;`, //3
}

var tvXMLMultiEscapeExpect = [...]string{
	`'`,
	`"`,
	`'`,
	`"`,
}

func TestXMLUnescaping(t *testing.T) {

	for i := 0; i < len(tvXMLEscaped); i++ {
		s := XMLUnescape(tvXMLEscaped[i])
		if tvXMLUnescaped[i] != s {
			t.Errorf("XMLUnescape(tvXMLEscaped[%d])doesn't match test vector\n", i)
		}
	}

}

func TestXMLMultiEscaped(t *testing.T) {

	for i := 0; i < len(tvXMLMultiEscapeTest); i++ {
		s := XMLUnescape(tvXMLMultiEscapeTest[i])
		if tvXMLMultiEscapeExpect[i] != s {
			t.Errorf("XMLUnescape(tvXMLMultiEscapeTest[%d])doesn't match test vector\n", i)
		}
	}

}

func TestXMLEscaping(t *testing.T) {

	for i := 0; i < len(tvXMLUnescaped); i++ {
		s, err := XMLEscapeStrict(tvXMLUnescaped[i])
		if err != nil {
			t.Errorf("XMLEscape(tvXMLUnescaped[%d]) failed (%v)\n", i, err)
		}

		if tvXMLEscaped[i] != s {
			// t.Log(s)
			// t.Log(tvXMLEscaped[i])
			t.Errorf("XMLEscape(tvXMLUnescaped[%d]) doesn't match test vector\n", i)
		}
	}

}

func stringsUnescape(s string) (t string, err error) {
	return strconv.Unquote(`"` + s + `"`)
}

func stringsEscape(s string) string {
	return strings.Trim(strconv.Quote(s), `"`)
}

func TestCombinedXMLStringsEscaping(t *testing.T) {

	// test whether the handling of quoting from xml back to strings is done correctly
	for i := 0; i < len(tvXMLEscaped); i++ {

		s := XMLUnescape(tvXMLEscaped[i])
		s = stringsEscape(s)
		s, err := stringsUnescape(s)
		if err != nil {
			t.Errorf("stringsUnescape(stringsEscape(XMLUnescape(tvXMLEscaped[%d]))) (%v)\n", i, err)
		}
		s, err = XMLEscapeStrict(s)
		if err != nil {
			t.Errorf("XMLEscape(stringsUnescape(stringsEscape(XMLUnescape(tvXMLEscaped[%d])))) (%v)\n", i, err)
		}
		if tvXMLEscaped[i] != s {
			t.Errorf("Expected symetric escaping to return original string for tvStringsEscaped[%d]\n", i)
		}
	}
}
