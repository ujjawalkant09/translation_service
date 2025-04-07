package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"github.com/ujjawalkant09/translation_service/cli"
)

var wg sync.WaitGroup
var sourceLang string
var targetLang string
var sourceText string

func init() {
	// Flags initialization
	flag.StringVar(&sourceLang, "s", "en", "Source language[en]")
	flag.StringVar(&targetLang, "t", "fr", "Target language[fr]")
	flag.StringVar(&sourceText, "st", "", "Text to translate")
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Options : ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create a channel and start a goroutine for translation
	strChan := make(chan string)
	wg.Add(1)

	reqBody := &cli.RequestBody{
		SourceLang: sourceLang,
		TargetLang: targetLang,
		SourceText: sourceText,
	}


	go cli.RequestTranslate(reqBody, strChan, &wg)

	processedStr := strings.ReplaceAll(<-strChan, " + ", " ")

	fmt.Printf("%s\n", processedStr)

	close(strChan)
	wg.Wait()
}
