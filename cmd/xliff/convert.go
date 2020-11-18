package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"github.com/simpleapps-eu/translate"
	"github.com/simpleapps-eu/translate/dotstrings"
	"github.com/simpleapps-eu/translate/xliff"
)

func ConvertXliff(xlfName, outName string) {
	// Determine languages for xlfName and outName
	_, inFileName := path.Split(xlfName)
	ilang := strings.SplitN(inFileName, ".", 2)[0]
	_, outFileName := path.Split(outName)
	olang := strings.SplitN(outFileName, ".", 2)[0]
	convertSource := (ilang == olang && (ilang == "en" || ilang == "en-US"))

	xlfFile, err := os.Open(xlfName)
	if err != nil {
		panic(fmt.Errorf("Failed to open -xliff %q (%v)", xlfName, err))
	}
	defer xlfFile.Close()

	outFile, err := os.Create(outName)
	if err != nil {
		panic(fmt.Errorf("Failed to open -out %q (%v)", outName, err))
	}
	defer outFile.Close()

	xlfChan, errChan := xliff.LoadTranslationUnits(xlfFile)

	var msgChan <-chan dotstrings.Message

	if convertSource {
		msgChan = translate.ConvertTranslationUnitsToSourceMessages(xlfChan)
	} else {
		msgChan = translate.ConvertTranslationUnitsToTargetMessages(xlfChan)
	}

	dotstrings.SaveMessages(msgChan, dotstrings.NewWriterUTF16(outFile))

	if err := <-errChan; err != nil {
		panic(err)
	}
}

/*
Implement this as a merge of the srcName source .strings file using
the strings in the tgtName target .strings file. If there is no translation
then skip the generation of the translation unit.
Put the Ctx of the source file as Note in the xlf file.
Mark an xlf entry as fuzzy when the source Str value is different from the
target Str value.
*/
func ConvertSourceAndTarget(srcName, tgtName, xlfName string) {
	fmt.Println("Not implemented yet!")
}

// Convert reads the en.strings and write out a fresh .xlf file to be send
// on to translators.
//
// e.g. xliff -source en.strings -xliff fr.xlf
func ConvertSource(srcName, xlfName string) {
	// Check language of srcName is either en or en-US
	_, srcFileName := path.Split(srcName)
	slang := strings.SplitN(srcFileName, ".", 2)[0]
	if len(slang) > 0 && slang != "en" && slang != "en-US" {
		panic(fmt.Errorf("Invalid source language -source %q (must be \"en\" or \"en-US\", found %q)", srcName, slang))
	}

	// Deduce translation file metadata
	_, xlfFileName := path.Split(xlfName)
	tlang := strings.SplitN(xlfFileName, ".", 2)[0]
	if len(tlang) > 0 && tlang != "en" && tlang != "en-US" {
		fmt.Printf("Converting to Target Language %q\n", tlang)
	} else {
		tlang = ""
	}
	tf := &xliff.TranslationFile{Original: "Localizable.strings", SourceLanguage: "en-US", Datatype: "x-strings", TargetLanguage: tlang}

	inFile, err := os.Open(srcName)
	if err != nil {
		panic(fmt.Errorf("Failed to open -source %q (%v)", srcName, err))
	}
	defer inFile.Close()

	xlfFile, err := os.Create(xlfName)
	if err != nil {
		panic(fmt.Errorf("Failed to create -xliff %q (%v)", xlfName, err))
	}
	defer xlfFile.Close()

	// Read strings from srcName and write xlf to xlfName
	fmt.Printf("Converting strings file %q to xliff file %q\n", srcName, xlfName)
	msgChan, errChan1 := dotstrings.LoadMessages(dotstrings.NewReaderUTF16(inFile))
	unitChan, errChan2 := translate.ConvertSourceMessagesToTranslationUnits(msgChan, tf)
	n := xliff.SaveTranslationUnits(unitChan, xlfFile)
	if err, _ := <-errChan2; err != nil {
		panic(err)
	}
	if err, _ := <-errChan1; err != nil {
		panic(err)
	}
	fmt.Printf("Converted %d strings\n", n)
}

// Convert reads the xx.strings and write out a fresh xx.xlf file to be send
// on to translators. The target language is taken from the text until
// the first dot of the target and xliff filename.
//
// e.g. xliff -target fr.strings -xliff fr.xlf
func ConvertTarget(tgtName, xlfName string) {
	// Check language of tgtName is the same as the language of the xliff file.
	_, tgtFileName := path.Split(tgtName)
	fromtlang := strings.SplitN(tgtFileName, ".", 2)[0]

	// Deduce translation file metadata
	_, xlfFileName := path.Split(xlfName)
	tlang := strings.SplitN(xlfFileName, ".", 2)[0]

	if fromtlang != tlang {
		panic(fmt.Errorf("Mismatching target languages %q and %q ", fromtlang, tlang))
	}

	if tlang == "en" || tlang == "en-US" {
		panic(fmt.Errorf("Invalid language for -target %q (%q is not a valid target language)", tgtName, tlang))
	}

	if len(tlang) > 0 {
		fmt.Printf("Converting to Target Language %q\n", tlang)
	}
	tf := &xliff.TranslationFile{Original: "Localizable.strings", SourceLanguage: "en-US", Datatype: "x-strings", TargetLanguage: tlang}

	tgtFile, err := os.Open(tgtName)
	if err != nil {
		panic(fmt.Errorf("Failed to open -target %q (%v)", tgtName, err))
	}
	defer tgtFile.Close()

	xlfFile, err := os.Create(xlfName)
	if err != nil {
		panic(fmt.Errorf("Failed to create -xliff %q (%v)", xlfName, err))
	}
	defer xlfFile.Close()

	// Read strings from tgtName and write xlf to xlfName
	fmt.Printf("Converting strings file %q to xliff file %q\n", tgtName, xlfName)
	msgChan, errChan1 := dotstrings.LoadMessages(dotstrings.NewReaderUTF16(tgtFile))
	unitChan, errChan2 := translate.ConvertTargetMessagesToTranslationUnits(msgChan, tf)
	n := xliff.SaveTranslationUnits(unitChan, xlfFile)
	if err, _ := <-errChan2; err != nil {
		panic(err)
	}
	if err, _ := <-errChan1; err != nil {
		panic(err)
	}
	fmt.Printf("Converted %d strings\n", n)
}
