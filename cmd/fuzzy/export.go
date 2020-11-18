package main

import (
	"fmt"
	"os"
	"github.com/simpleapps-eu/translate"
	"github.com/simpleapps-eu/translate/dotstrings"
)

func export(fuzzy bool, missing bool, srcName string, tmName string, tgtName string) (err error) {
	// Open source file
	srcFile, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Load map of translations from the tm file.
	translations, err := dotstrings.LoadMessagesMapFromFile(tmName)
	if err != nil {
		return
	}

	// Open target file
	tgtFile, err := os.Create(tgtName)
	if err != nil {
		return
	}
	defer tgtFile.Close()

	// Start loading the messages asynchronously from the srcFile
	msgChan, errChan := dotstrings.LoadMessages(dotstrings.NewReaderUTF16(srcFile))

	// Start translating messages asynchronously
	msgChan = translate.TranslateMessages(msgChan, translations)

	// Filter out any (non)fuzzy messages asynchronously
	msgChan = exportMessages(fuzzy, missing, msgChan)

	// Finally write the fuzzy messages to a file synchronously.
	n := dotstrings.SaveMessages(msgChan, dotstrings.NewWriterUTF16(tgtFile))

	// If there wasn't an error reported during processing, report the status here.

	err, _ = <-errChan
	if err != nil {
		return
	}

	if fuzzy {
		fmt.Printf("%d\tFuzzy strings written to %q\n", n, tgtName)
	} else {
		fmt.Printf("%d\tNormal strings written to %q\n", n, tgtName)
	}
	return
}

func exportMessages(fuzzy bool, missing bool, srcChan <-chan dotstrings.Message) <-chan dotstrings.Message {
	tgtChan := make(chan dotstrings.Message, 3)

	extractor := func(srcChan <-chan dotstrings.Message, tgtChan chan<- dotstrings.Message) {
		defer close(tgtChan)

		// Process the translated messages, extract fuzzies to file, don't mark them as fuzzy though.
		for m := range srcChan {
			if fuzzy == m.Fuzzy && (missing || !m.Missing) {
				m.Fuzzy = false
				tgtChan <- m
			}
		}
	}

	go extractor(srcChan, tgtChan)
	return tgtChan
}
