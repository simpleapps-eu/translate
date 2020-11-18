package main

import (
	"flag"
	"os"
)

var (
	outname string
	srcname string
	tgtname string
	xlfname string
)

/*
source,xliff,out => translate source .strigns file using xliff file to out .strings file (which is formatted as a target .strings file)
source,out => normalize source .strings file writing result to out .strings file, normalize entries and report errors
target,out => normalize target .strings file writing result to out .strings file, normalize entries and report errors
xliff,out => convert xliff to a target formatted .strings file using source element as Ctx and target element as Str.

source,target,xliff => combine source,target and write result to xliff file. Only output translation units where id's exist in both files.
source,xliff => convert source to xliff
target,xliff => convert target to xliff
*/

func init() {
	flag.StringVar(&outname, "out", "", ".strings file in target language that is written during translation.")
	flag.StringVar(&srcname, "source", "", ".strings file in source language.")
	flag.StringVar(&tgtname, "target", "", ".strings file in target language.")
	flag.StringVar(&xlfname, "xliff", "", ".xlf file to use for translation (when -out is set) or to be written.")
}

func main() {
	defer catch()

	flag.Parse()

	if len(outname) > 0 {
		// Translate
		if srcname != "" && xlfname != "" {
			Translate(srcname, xlfname, outname)
			return
		}

		// Normalize
		if srcname != "" {
			Normalize(srcname, outname)
			return
		}

		// Normalize
		if tgtname != "" {
			Normalize(tgtname, outname)
			return
		}

		// Convert from xlf to strings
		if xlfname != "" {
			ConvertXliff(xlfname, outname)
			return
		}

	} else {
		// Convert source and target language .strings file to xlf
		if srcname != "" && tgtname != "" && xlfname != "" {
			ConvertSourceAndTarget(srcname, tgtname, xlfname)
			return
		}

		// Convert source language .strings to xlf
		if srcname != "" && xlfname != "" {
			ConvertSource(srcname, xlfname)
			return
		}

		// Convert target language .strings to xlf
		if tgtname != "" && xlfname != "" {
			ConvertTarget(tgtname, xlfname)
			return
		}
	}

	flag.Usage()
	panic(1)
}

func catch() {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case error:
			println(e.Error())
			os.Exit(1)
		case int:
			os.Exit(e)
		default:
			panic(err)
		}
	}
}
