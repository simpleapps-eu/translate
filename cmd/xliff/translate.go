package main

import (
	"fmt"
	"os"
	"github.com/simpleapps-eu/translate"
	"github.com/simpleapps-eu/translate/dotstrings"
	"github.com/simpleapps-eu/translate/xliff"
)

// Translate reads the strings from -source <file> .strings file, then translates
// them using the translation from -xliff <file> XLIFF file and writes out
// to -out <file> .strings file.
//
//  e.g. xliff -source en.strings -xliff fr.xlf -out fr.strings
func Translate(inName, xlfName, outName string) {

	inFile, err := os.Open(inName)
	if err != nil {
		panic(fmt.Errorf("Failed to open -in %q (%v)", inName, err))
	}
	defer inFile.Close()

	xlfFile, err := os.Open(xlfName)
	if err != nil {
		panic(fmt.Errorf("Failed to open -xlf %q (%v)", xlfName, err))
	}
	defer xlfFile.Close()

	// Read in the translation from the xlf file and store it based on Resname in a map
	tf, translation, err := xliff.LoadTranslationMap(xlfFile)
	if err != nil {
		panic(fmt.Errorf("Error processing -xlf %q (%v)", xlfName, err))
	}

	outFile, err := os.Create(outName)
	if err != nil {
		panic(fmt.Errorf("Failed to create -out %q (%v)", outName, err))
	}
	defer outFile.Close()

	// Read strings from in and write translated strings to out
	fmt.Printf("Translating %q to %q using %q\n", inName, outName, xlfName)

	msgChan, errChan1 := dotstrings.LoadMessages(dotstrings.NewReaderUTF16(inFile))
	msgChan, errChan2 := translate.TranslateMessagesXLIFF(msgChan, translation)
	transcount := dotstrings.SaveMessages(msgChan, dotstrings.NewWriterUTF16(outFile))
	if err, errorOccurred := <-errChan1; errorOccurred {
		panic(fmt.Errorf("Failure while loading messages from -in %q (%v)", inName, err))
	}
	if err, errorOccurred := <-errChan2; errorOccurred {
		panic(fmt.Errorf("Failure while translating -in %q using -xlf %q (%v)", inName, xlfName, err))
	}

	fmt.Printf("Translated %d strings from %q to %q\n", transcount, tf.SourceLanguage, tf.TargetLanguage)
}
