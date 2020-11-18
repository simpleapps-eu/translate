package xliff

import (
	"errors"
	"io"

	"github.com/simpleapps-eu/translate/xliff/exml"
)

// LoadTranslationMap reads xliff translation units from an xml file and then
// creates a translation map out of them. The mandatory id attribute in the trans-unit element
// is expected to match the id in the strings file. Both ID and Target value from the xliff file
// are unescaped before being written to the translation table.
func LoadTranslationMap(reader io.Reader) (tf *TranslationFile, translation map[string]string, err error) {
	// Read in the translation from the xlf file and store it based on Resname in a map
	translation = make(map[string]string)
	tuchan, echan := LoadTranslationUnits(reader)
	for tu := range tuchan {
		if tf == nil {
			tf = tu.File
		} else {
			if tf != tu.File {
				err = errors.New("Multiple files in a single xlf are currently not supported")
				return
			}
		}
		translation[XMLUnescape(tu.ID)] = XMLUnescape(tu.Target)
	}
	err, _ = <-echan
	return
}

// LoadTranslationUnits returns a channel of TranslationUnit values and will start
// processing the xliff file passed in via the reader argument asynchronously.
// Whenever it it has read a TranslationUnit, this will written to the channel.
func LoadTranslationUnits(reader io.Reader) (<-chan TranslationUnit, <-chan error) {
	tuchan := make(chan TranslationUnit)
	echan := make(chan error, 1)
	go func(reader io.Reader, tuchan chan<- TranslationUnit, echan chan<- error) {

		defer close(tuchan)
		defer close(echan)

		if reader == nil {
			echan <- errors.New("argument reader is nil")
			return
		}

		decoder := exml.NewDecoder(reader)
		if decoder == nil {
			echan <- errors.New("exml.NewDecoder returned nil")
			return
		}

		var tf *TranslationFile
		var tu *TranslationUnit
		decoder.On("xliff/file", func(attrs exml.Attrs) {

			tf = &TranslationFile{}
			original, err := attrs.Get("original")
			if err == nil {
				tf.Original = original
			}
			sourceLanguage, err := attrs.Get("source-language")
			if err == nil {
				tf.SourceLanguage = sourceLanguage
			}
			datatype, err := attrs.Get("datatype")
			if err == nil {
				tf.Datatype = datatype
			}
			targetLanguage, err := attrs.Get("target-language")
			if err == nil {
				tf.TargetLanguage = targetLanguage
			}

			decoder.On("body/trans-unit", func(attrs exml.Attrs) {

				if tu != nil {
					tuchan <- *tu
				}
				tu = &TranslationUnit{File: tf}

				id, err := attrs.Get("id")
				if err != nil {
					decoder.Error(err)
					return
				}
				tu.ID = id

				decoder.On("source/$text", func(text exml.CharData) {
					tu.Source = string(text)
				})

				decoder.On("target/$text", func(text exml.CharData) {
					tu.Target = string(text)
				})

				decoder.On("note/$text", func(text exml.CharData) {
					tu.Note = string(text)
				})
			})
		})
		decoder.OnError(func(err error) {
			echan <- err
		})
		decoder.Run()

		// Make sure the final TranslationUnit is also send
		if tu != nil {
			tuchan <- *tu
		}

	}(reader, tuchan, echan)
	return tuchan, echan
}
