package xliff

import (
	"io"
	"text/template"
)

const head = `<?xml version="1.0"?>
<xliff version="1.2">
<file original="{{.Original}}" source-language="{{.SourceLanguage}}"{{with .TargetLanguage}} target-language="{{.}}"{{end}} datatype="{{.Datatype}}">
<body>
`
const unit = `<trans-unit id="{{.ID}}">
{{- with .Source}}
<source>{{.}}</source>
{{- end}}
{{- with .Target}}
<target>{{.}}</target>
{{- end}}
{{- with .Note}}
<note>{{.}}</note>
{{- end}}
</trans-unit>
`
const foot = `</body>
</file>
</xliff>
`

// SaveTranslationUnits will take a channel with translation units and stream
// them to writer in xliff xml format. The function will return when all translation
// units have been written.
// The goroutine feeding srcChan should close the channel once it has finished.
// The closing of the channel indicates to SaveTranslation that it can finish too.
// The function then returns the number of translation units it has written to
// the writer.
func SaveTranslationUnits(srcChan <-chan TranslationUnit, writer io.Writer) (n int) {
	headTpl := template.Must(template.New("head").Parse(head))
	unitTpl := template.Must(template.New("unit").Parse(unit))
	headWritten := false
	for m := range srcChan {
		if !headWritten {
			headWritten = true
			headTpl.Execute(writer, m.File)
		}
		unitTpl.Execute(writer, m)
		n++
	}
	if headWritten {
		writer.Write([]byte(foot))
	}
	return
}
