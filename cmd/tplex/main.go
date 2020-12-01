package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"text/template"
)

/*
  data
  {
    "firstName": "John",
    "surName": "Doe",
    "age": 42,
    "email": ["john@doe.com","john.doe@gmail.com"]
  }

  templ
    This is {{.firstName}} {{.surName}}, {{.age}} years old.
    Send a message to {{range .email}}{{.}}, {{end}}
*/

var (
	templName string
	dataName  string
	outName   string
)

func init() {
	flag.StringVar(&templName, "template", "", "template file to execute")
	flag.StringVar(&dataName, "data", "", "file with json data to execute the template with.")
	flag.StringVar(&outName, "out", "", "file to write the template execution result to")
}

func main() {
	defer catch()

	// Flag checking
	flag.Parse()
	if flag.NFlag() < 3 || flag.NFlag() > 3 {
		flag.Usage()
		panic(-1)
	}

	templFile, err := os.Open(templName)
	if err != nil {
		panic(err)
	}
	defer templFile.Close()

	dataFile, err := os.Open(dataName)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	dataDecoder := json.NewDecoder(dataFile)

	// decode an array value (Message)
	var v interface{}
	err = dataDecoder.Decode(&v)
	if err != nil {
		panic(err)
	}

	templ, err := template.ParseFiles(templName)
	if err != nil {
		panic(err)
	}

	if outName == "-" {
		templ.Execute(os.Stdout, v)
	} else {
		outFile, err := os.Create(outName)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()
		templ.Execute(outFile, v)
	}
}

func catch() {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case error:
			log.Fatalf("Error: %v", e)
		case int:
			os.Exit(e)
		default:
			panic(err)
		}
	}
}
