package main

import (
	"fmt"
	"os"
	"github.com/simpleapps-eu/translate"
	"github.com/simpleapps-eu/translate/dotstrings"
)

func count(fuzzy bool, missing bool, srcName string, tmName string) (err error) {
	// Open source file
	srcFile, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Load map of translations from tm file
	translations, err := dotstrings.LoadMessagesMapFromFile(tmName)
	if err != nil {
		return
	}

	// Start loading the messages asynchronously
	msgChan, errChan := dotstrings.LoadMessages(dotstrings.NewReaderUTF16(srcFile))

	// Translate messages asynchronously
	msgChan = translate.TranslateMessages(msgChan, translations)

	// Synchronously receive all messages from the msgChan
	var n uint32
	for m := range msgChan {
		if fuzzy == m.Fuzzy && (missing || !m.Missing) {
			n++
		}
	}

	err, _ = <-errChan
	if err != nil {
		return
	}

	if fuzzy {
		fmt.Printf("%d\tFuzzy strings\n", n)
	} else {
		fmt.Printf("%d\tNormal strings\n", n)
	}
	return
}
