package dotstrings

import (
	"io"
	"text/template"
)

// SaveMessages will take a channel with messages and stream them to character
// stream dstWriter. The function will return when all messages have been sent.
// The goroutine feeding srcChan should close the channel once it has finished.
// The closing of the channel indicates to SaveMessages that it can finish too.
// The function returns the number of messages it has written to the dstWriter.
func SaveMessages(srcChan <-chan Message, dstWriter io.Writer) (n int) {
	entryTpl := template.Must(template.New("strings").Parse("{{if .Fuzzy}}/* Fuzzy */\n{{end}}/* {{.Ctx}} */\n\"{{.ID}}\" = \"{{.Str}}\";\n\n"))
	for src := range srcChan {
		entryTpl.Execute(dstWriter, src)
		n++
	}
	return
}
