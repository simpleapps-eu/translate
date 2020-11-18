package main

import (
	"fmt"
	"os"
	"github.com/simpleapps-eu/translate/dotstrings"
)

// Normalize reads the strings from -in <file> .strings file and then writes
// them out again to the -out <file> .strings file. This detects any errors
// in the .strings file, normalizes the strings and cleans up any formatting
// issues while writing them back out.
//
// e.g. xliff -out normalized.strings -in en.strings
func Normalize(inName, outName string) {
	inFile, err := os.Open(inName)
	if err != nil {
		panic(fmt.Errorf("Failed to open -in %q (%v)", inName, err))
	}
	defer inFile.Close()

	outFile, err := os.Create(outname)
	if err != nil {
		panic(fmt.Errorf("Failed to create -out %q (%v)", outName, err))
	}
	defer outFile.Close()

	// Read strings from inName and write normalized strings to outName
	fmt.Printf("Normalizing %q writing result to %q \n", inName, outName)
	msgChan, errChan := dotstrings.LoadMessages(dotstrings.NewReaderUTF16(inFile))
	n := dotstrings.SaveMessages(msgChan, dotstrings.NewWriterUTF16(outFile))
	if err, errorOccurred := <-errChan; errorOccurred {
		panic(err)
	}
	fmt.Printf("Normalized %d strings\n", n)
}
