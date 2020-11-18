package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	normal, present, doCount, doExport, doImport bool
	tmName                                       string
	srcName                                      string
	tgtName                                      string
)

func init() {
	flag.BoolVar(&normal, "normal", false, "operate on normal instead of fuzzy strings for -count and -export")
	flag.BoolVar(&present, "present", false, "operate on fuzzy strings not missing from the -tm file for -count and -export")
	flag.BoolVar(&doCount, "count", false, "count the (non)fuzzy strings, needs -source and -tm")
	flag.BoolVar(&doExport, "export", false, "export the (non)fuzzy strings to -target file, needs -source, -tm and -target")
	flag.BoolVar(&doImport, "import", false, "import translated strings from -target file and merge into -tm file")

	flag.StringVar(&tmName, "tm", "", "translation file used to translate source strings into target strings")
	flag.StringVar(&srcName, "source", "", "file to read source strings from")
	flag.StringVar(&tgtName, "target", "", "file to read/write translated strings")
}

func main() {
	defer catch()

	flag.Parse()

	if flag.NArg() != 0 || flag.NFlag() < 2 || flag.NFlag() > 5 {
		flag.Usage()
		panic(1)
	}
	fuzzy := !normal
	missing := !present

	tmExt := filepath.Ext(tmName)
	if !strings.EqualFold(tmExt, ".strings") {
		panic(fmt.Errorf("Error: Unsupported -tm file type %q", tmExt))
	}
	srcExt := filepath.Ext(srcName)
	if len(srcExt) > 0 {
		if !strings.EqualFold(srcExt, ".strings") {
			panic(fmt.Errorf("Error: Unsupported -src file type %q", srcExt))
		}
	}
	tgtExt := filepath.Ext(tgtName)
	if len(tgtExt) > 0 {
		if !strings.EqualFold(tgtExt, ".strings") {
			panic(fmt.Errorf("Error: Unsupported -tgt file type %q", tgtExt))
		}
	}

	if doExport {
		if err := export(fuzzy, missing, srcName, tmName, tgtName); err != nil {
			panic(err)
		}
		return
	}

	if doImport {
		if err := merge(tmName, tgtName); err != nil {
			panic(err)
		}
		return
	}

	// Just count them
	if err := count(fuzzy, missing, srcName, tmName); err != nil {
		panic(err)
	}
	return
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
