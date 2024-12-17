package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func ApiQuery() {

	var k string = readAPIKey()
	qs := []string{
		"function=TIME_SERIES_DAILY",
		"symbol=IBM",
		"outputsize=compact",
		fmt.Sprintf("apikey=%v", k),
	}

	url := alphaQueryGen(qs)
	fmt.Println(url)

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Print the response
	fmt.Println(string(body))

}

func alphaQueryGen(qs []string) string {
	// Generates a query for the alphavantage API

	url := "https://www.alphavantage.co/query?"
	for _, q := range qs {
		url = fmt.Sprintf("%v%v&", url, q)
	}

	return url

}

func readAPIKey() string {
	//Reads my API key from the text file for alphavantage

	//File Path of the infomation file
	_, filePath, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filePath)
	filePath = fmt.Sprintf("%v/apiKey.txt", dir)

	// Reads the file
	k, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %v: %v", filePath, err)
	}

	return string(k)

}
