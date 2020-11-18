package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"github.com/simpleapps-eu/translate/dotstrings"
)

func merge(tmName, tgtName string) (err error) {

	// Perform the merge into an in memory bytes.Buffer
	resultBuf := &bytes.Buffer{}
	n, err := fuzzyMergeTo(tmName, tgtName, resultBuf)
	if err != nil {
		return
	}

	// Open existing translation file for writing
	tmFile, err := os.Create(tmName)
	if err != nil {
		return
	}
	defer tmFile.Close()

	// Write the resultBuf to the tmFile.
	_, err = resultBuf.WriteTo(dotstrings.NewWriterUTF16(tmFile))

	// If there wasn't an error reported during processing, report the status here.
	if err != nil {
		return
	}

	fmt.Printf("%d\tStrings written to %q\n", n, tmName)
	return
}

func fuzzyMergeTo(tmName, tgtName string, writer io.Writer) (n int, err error) {
	// Load the fuzzies file with updated translations
	newTranslations, err := dotstrings.LoadMessagesMapFromFile(tgtName)
	if err != nil {
		return
	}

	// Open existing translation file that needs to be updated
	tmFile, err := os.Open(tmName)
	if err != nil {
		return
	}
	defer tmFile.Close()

	// Asynchronously load the messages from the existing translation
	msgChan, errChan := dotstrings.LoadMessages(dotstrings.NewReaderUTF16(tmFile))

	// Asynchronously replace outdated existing translations with new translations
	msgChan = replaceOutdatedTranslations(msgChan, newTranslations)

	// Now synchronously save the msgChan to the writer
	n = dotstrings.SaveMessages(msgChan, writer)

	err, _ = <-errChan
	if err != nil {
		return
	}

	// Open new translation file again
	tgtFile, err := os.Open(tgtName)
	if err != nil {
		return
	}
	defer tgtFile.Close()

	// Asynchronously load the messages from the new translation
	msgChan, errChan = dotstrings.LoadMessages(dotstrings.NewReaderUTF16(tgtFile))

	// Asynchronously append new translations that were not written during the previous phase.
	msgChan = appendNewTranslations(msgChan, newTranslations)

	// Now synchronously append msgChan entries to the resultBuf
	n += dotstrings.SaveMessages(msgChan, writer)

	err, _ = <-errChan
	return
}

func replaceOutdatedTranslations(srcChan <-chan dotstrings.Message, newTranslations map[string]dotstrings.Message) <-chan dotstrings.Message {
	dstChan := make(chan dotstrings.Message, 3)

	replacer := func(srcChan <-chan dotstrings.Message, newTranslations map[string]dotstrings.Message, dstChan chan<- dotstrings.Message) {
		defer close(dstChan)

		for src := range srcChan {
			if tran, isNewTranslation := newTranslations[src.ID]; isNewTranslation {
				if len(tran.Ctx) == 0 {
					// New translation doesn't have context, use context of original.
					dstChan <- dotstrings.Message{Ctx: src.Ctx, ID: tran.ID, Str: tran.Str}
				} else {
					// Send the new translation as-is. Don't worry about Fuzzy
					// flag as loading it as a translation would have croaked.
					dstChan <- tran
				}
				delete(newTranslations, tran.ID)
			} else {
				// Send the entry from the source translation to the output.
				dstChan <- src
			}
		}
	}

	go replacer(srcChan, newTranslations, dstChan)
	return dstChan
}

func appendNewTranslations(srcChan <-chan dotstrings.Message, newTranslations map[string]dotstrings.Message) <-chan dotstrings.Message {
	dstChan := make(chan dotstrings.Message, 3)

	appender := func(srcChan <-chan dotstrings.Message, newTranslations map[string]dotstrings.Message, dstChan chan<- dotstrings.Message) {
		defer close(dstChan)
		for src := range srcChan {
			// Send new translation that has not been sent yet and remove it from the map.
			if tran, ok := newTranslations[src.ID]; ok {
				dstChan <- tran
				delete(newTranslations, tran.ID)
			}
		}
	}
	go appender(srcChan, newTranslations, dstChan)
	return dstChan
}
