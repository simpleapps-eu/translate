package dotstrings

import (
	"io"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// NewReaderUTF16 will create a reader that expects data in UTF16 LittleEndian
// mode and it will also expect the matching BOM.
func NewReaderUTF16(fileReader io.Reader) io.Reader {
	encoding := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
	utf16Decoder := encoding.NewDecoder()
	return transform.NewReader(fileReader, utf16Decoder)
}

// NewWriterUTF16 will create a writer that will stream out text in UTF16
// LittleEndian mode with a BOM.
func NewWriterUTF16(fileWriter io.Writer) io.Writer {
	encoding := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
	utf16Encoder := encoding.NewEncoder()
	return transform.NewWriter(fileWriter, utf16Encoder)
}
