package cli

import (
	"log"
	"net/http"
	"sync"
	"github.com/Jeffail/gabs"
)

const translateUrl = "https://translate.googleapis.com/translate_a/single"

type RequestBody struct {
	SourceLang string 
	TargetLang string 
	SourceText string 
}

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}

	req, err := http.NewRequest("GET", translateUrl, nil)
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}

	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		str <- "You have been rate limited, Try again later."
		return
	}

	parsedJson, err := gabs.ParseJSONBuffer(res.Body)
	if err != nil {
		log.Fatalf("Error parsing response: %s", err)
	}

	nestOne, err := parsedJson.ArrayElement(0)
	if err != nil {
		log.Fatalf("Error accessing first element: %s", err)
	}

	nestTwo, err := nestOne.ArrayElement(0)
	if err != nil {
		log.Fatalf("Error accessing second element: %s", err)
	}

	translatedStr, err := nestTwo.ArrayElement(0)
	if err != nil {
		log.Fatalf("Error accessing translated string: %s", err)
	}

	str <- translatedStr.Data().(string)
}