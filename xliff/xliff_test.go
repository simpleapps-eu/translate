package xliff

import (
	"strings"
	"testing"

	"github.com/simpleapps-eu/translate/xliff/exml"
)

const xliffData = `<?xml version="1.0"?>
<xliff version="1.2">
<file original="Localizable.strings" source-language="en-US" datatype="STRINGS" target-language="es">
<body>
<trans-unit id="JUAgTWluZCBNYXA">
<source>%@ Mind Map</source><target>%@ Mapa Mental</target>
</trans-unit>
<trans-unit id="JWQgZGF5IGFnbw">
<source>%d day ago</source><target>Hace %d día</target>
</trans-unit>
<trans-unit id="JWQgZGF5cyBhZ28">
<source>%d days ago</source><target>Hace %d días</target>
</trans-unit>
<trans-unit id="JWQgaG91ciBhZ28">
<source>%d hour ago</source><target>Hace %d hora</target>
</trans-unit>
<trans-unit id="JWQgaG91cnMgYWdv">
<source>%d hours ago</source><target>Hace %d horas</target>
</trans-unit>
<trans-unit id="JWQgbWludXRlIGFnbw">
<source>%d minute ago</source><target>Hace %d minuto</target>
</trans-unit>
</body>
</file>
</xliff>
`

var expectMap = map[string]struct {
	source string
	target string
}{
	"JUAgTWluZCBNYXA":    {source: "%@ Mind Map", target: "%@ Mapa Mental"},
	"JWQgZGF5IGFnbw":     {source: "%d day ago", target: "Hace %d día"},
	"JWQgZGF5cyBhZ28":    {source: "%d days ago", target: "Hace %d días"},
	"JWQgaG91ciBhZ28":    {source: "%d hour ago", target: "Hace %d hora"},
	"JWQgaG91cnMgYWdv":   {source: "%d hours ago", target: "Hace %d horas"},
	"JWQgbWludXRlIGFnbw": {source: "%d minute ago", target: "Hace %d minuto"},
}

var expectKeys = []string{
	"JUAgTWluZCBNYXA",
	"JWQgZGF5IGFnbw",
	"JWQgZGF5cyBhZ28",
	"JWQgaG91ciBhZ28",
	"JWQgaG91cnMgYWdv",
	"JWQgbWludXRlIGFnbw",
}

func TestParsing(t *testing.T) {

	// just test whether exml will do a proper parse of the test data

	reader := strings.NewReader(xliffData)
	if reader == nil {
		t.Error("expected NewReader got nil")
	}

	decoder := exml.NewDecoder(reader)
	if decoder == nil {
		t.Error("expected NewDecoder got nil")
	}

	decoder.On("xliff/file/body/trans-unit", func(attrs exml.Attrs) {

		resname, err := attrs.Get("id")
		if err != nil {
			t.Error(err)
		} else {

			st, ok := expectMap[resname]
			if !ok {

				t.Errorf("id %q not found in tests table", resname)

			} else {
				decoder.On("source/$text", func(text exml.CharData) {

					if st.source != string(text) {
						t.Errorf("Expected %q got %q", st.source, string(text))
					}
				})

				decoder.On("target/$text", func(text exml.CharData) {

					if st.target != string(text) {
						t.Errorf("Expected %q got %q", st.target, string(text))
					}

				})
			}
		}
	})

	decoder.Run()

}

func TestTranslationUnits(t *testing.T) {

	type expect struct {
		ID     string
		Source string
		Target string
	}
	expecting := make(chan *expect)
	go func() {
		defer close(expecting)
		for _, id := range expectKeys {
			expecting <- &expect{id, expectMap[id].source, expectMap[id].target}
		}
	}()

	reader := strings.NewReader(xliffData)
	if reader == nil {
		t.Error("expected NewReader got nil")
	}
	var testcount int
	tuchan, echan := LoadTranslationUnits(reader)
	done := make(chan struct{})
	go func() {
		for tu := range tuchan {
			// we should do something with the TranslationUnit

			expect, ok := <-expecting
			if !ok {
				t.Error("Expected test to be finished but got more data")
			} else {
				if expect.ID != tu.ID {
					t.Errorf("Expected id %q got %q", expect.ID, tu.ID)
				}
				if expect.Source != tu.Source {
					t.Errorf("Expected source %q got %q", expect.Source, tu.Source)
				}
				if expect.Target != tu.Target {
					t.Errorf("Expected target %q got %q", expect.Target, tu.Target)
				}

				testcount++
			}
		}
		close(done)
	}()
	err, errReceived := <-echan
	if errReceived {
		t.Error(err)
	}
	if _, ok := <-done; ok {
		t.Error("Expected done to be closed")
	}
	if testcount != len(expectKeys) {
		t.Errorf("Expected to test %d cases but only %d where actually tested", len(expectKeys), testcount)
	}
}
