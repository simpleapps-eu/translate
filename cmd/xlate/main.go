package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"github.com/simpleapps-eu/translate"
	"github.com/simpleapps-eu/translate/dotstrings"
)

var (
	tmName     string
	tmfbName   string
	srcName    string
	tgtName    string
	forcePLIST bool
)

func init() {
	flag.StringVar(&tmName, "tm", "", "file used as translation memory")
	flag.StringVar(&tmfbName, "tmfb", "", "file used as fallback translation memory")
	flag.StringVar(&srcName, "source", "", "file for reading source strings")
	flag.StringVar(&tgtName, "target", "", "file to write the translated target strings to")
	flag.BoolVar(&forcePLIST, "plist", false, "Interpret -source and -target as XML plist files")
}

func main() {
	defer catch()

	// Flag checking
	flag.Parse()
	if flag.NFlag() < 3 || flag.NFlag() > 5 {
		flag.Usage()
		panic(-1)
	}

	// Translation memory file
	tmExt := filepath.Ext(tmName)
	if !strings.EqualFold(tmExt, ".strings") {
		panic(fmt.Errorf("Error: Unsupported -tm file type %q", tmExt))
	}

	// Translation memory fallback file
	if len(tmfbName) > 0 {
		tmbfExt := filepath.Ext(tmfbName)
		if !strings.EqualFold(tmbfExt, ".strings") {
			panic(fmt.Errorf("Error: Unsupported -tmfb file type %q", tmbfExt))
		}
	}

	// Source file
	srcExt := filepath.Ext(srcName)
	if !strings.EqualFold(srcExt, ".strings") && !strings.EqualFold(srcExt, ".tpl") && !strings.EqualFold(srcExt, ".txt") && !strings.EqualFold(srcExt, ".plist") {
		panic(fmt.Errorf("Error: Unsupported -source file type %q", srcExt))
	}

	// Target file
	tgtExt := filepath.Ext(tgtName)
	if !strings.EqualFold(tgtExt, ".strings") && !strings.EqualFold(tgtExt, ".txt") {
		panic(fmt.Errorf("Error: Unsupported -target file type %q", tgtExt))
	}

	if srcName == tgtName {
		panic(fmt.Errorf("Error: -source and -target file cannot be the same"))
	}

	// Read translations from translation memory
	translations, err := dotstrings.LoadMessagesMapFromFile(tmName)
	if err != nil {
		panic(err)
	}

	// Read fallback translations from translation memory
	var translationsFallback map[string]dotstrings.Message
	if len(tmfbName) > 0 {
		translationsFallback, err = dotstrings.LoadMessagesMapFromFile(tmfbName)
		if err != nil {
			panic(err)
		}
		// Merge fallback translations where no primary translation is available
		for key, val := range translationsFallback {
			if _, ok := translations[key]; !ok {
				translations[key] = val
			}
		}
	}

	// open source
	srcFile, err := os.Open(srcName)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	// open destination
	tgtFile, err := os.Create(tgtName)
	if err != nil {
		panic(err)
	}
	defer tgtFile.Close()

	switch srcExt {
	case ".strings":
		if forcePLIST {
			n, err := translate.TranslatePlistFile(srcFile, translations, tgtFile)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Translated %d Plist Entries\n", n)
		} else {
			n, err := translate.TranslateMessagesFile(srcFile, translations, tgtFile)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Translated %d Strings Entries\n", n)
		}
	case ".tpl":
		n, err := translate.TranslateIDsFile(srcFile, translations, translationsFallback, tgtFile)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Translated %d IDs\n", n)
	case ".txt":
		n, err := translate.TranslateTextFile(srcFile, translations, tgtFile)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Translated %d Text Strings\n", n)
	}
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

func trace() {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Println(file, ":", line)
	}
}
