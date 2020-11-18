package plist

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

// <plist version="1.0">
// <dict>
// 	<key>Version</key>
// 	<string>Version</string>

// LoadEntries will read an XML format Property List (Plist) file.
// The contents of the Plist is expected to be an array of key value
// string entries where the key represents an ID and the value represents the
// text translation in the target language. Format is expected to be UTF-8.
func LoadEntries(srcFile io.Reader) (<-chan Entry, <-chan error) {
	entryChan := make(chan Entry, 3)
	errChan := make(chan error, 1)
	reader := func(srcFile io.Reader, entryChan chan<- Entry, errChan chan<- error) {
		defer close(entryChan)
		defer close(errChan)

		bytes, err := ioutil.ReadAll(srcFile)
		if err != nil {
			errChan <- err
			return
		}

		type PList struct {
			XMLName xml.Name `xml:"plist"`
			Keys    []string `xml:"dict>key"`
			Strings []string `xml:"dict>string"`
		}
		l := PList{}

		err = xml.Unmarshal(bytes, &l)
		if err != nil {
			errChan <- err
			return
		}

		if len(l.Keys) != len(l.Strings) {
			errChan <- fmt.Errorf("Number of Keys (%d) and Strings (%d) differ in PList Strings file", len(l.Keys), len(l.Strings))
			return
		}

		for idx, key := range l.Keys {
			entryChan <- Entry{ID: key, Str: l.Strings[idx]}
		}
	}

	go reader(srcFile, entryChan, errChan)
	return entryChan, errChan
}
