package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// The the structure of the response from the alpha vantage API
type MetaData struct {
	Information   string    `json:"1. Information"`
	Symbol        string    `json:"2. Symbol"`
	LastRefreshed time.Time `json:"3. Last Refreshed"`
	OutputSize    string    `json:"4. Output Size"`
	TimeZone      string    `json:"5. Time Zone"`
}

type DailyData struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

type ApiResponse struct {
	MetaData        MetaData              `json:"Meta Data"`
	TimeSeriesDaily map[string]*DailyData `json:"Time Series (Daily)"`
}

func ApiQuery(secID string) ApiResponse {
	//Queries Alphavantage

	// Reads API key from directory
	k, err := readAPIKey()
	if err != nil {
		log.Fatalf("Error reading API key: %v", err)
	}

	var secStruct ApiResponse

	qs := []string{
		"function=TIME_SERIES_DAILY",
		fmt.Sprintf("symbol=%v", secID),
		"outputsize=compact",
		fmt.Sprintf("apikey=%v", k),
	}

	url := alphaQueryGen(qs)

	// Gets the response from alphavantage API
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

	// Convert body of type []byte to json

	err = json.Unmarshal(body, &secStruct)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return secStruct

}

func alphaQueryGen(qs []string) string {
	// Generates a query for the alphavantage API

	url := "https://www.alphavantage.co/query?"
	for _, q := range qs {
		url = fmt.Sprintf("%v%v&", url, q)
	}

	return url

}

func readAPIKey() (string, error) {
	//Reads my API key from the text file for alphavantage

	// Loads .env file
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	// Gets the key from it's enviroment variable
	k := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if k == "" {
		return "", fmt.Errorf("no API key found in .env file: %v", err)
	}

	return string(k), nil

}

func (m *MetaData) UnmarshalJSON(data []byte) error {
	// Constructs the structure MetaData into the correct format from JSON

	//MetaData struct before conversion into correct types
	type rawMetaData struct {
		Information   string `json:"1. Information"`
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		OutputSize    string `json:"4. Output Size"`
		TimeZone      string `json:"5. Time Zone"`
	}

	//Get raw output
	var raw rawMetaData
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Convert string to time.TIme
	typedTime, err := time.Parse("2006-01-02", raw.LastRefreshed)
	if err != nil {
		return err
	}

	// Assign new typed values
	m.Information = raw.Information
	m.Symbol = raw.Symbol
	m.LastRefreshed = typedTime
	m.OutputSize = raw.OutputSize
	m.TimeZone = raw.TimeZone

	return nil

}

// Custom unmarshaler for DD
func (o *DailyData) UnmarshalJSON(data []byte) error {
	type RawDD struct {
		Open   string `json:"1. open"`
		High   string `json:"2. high"`
		Low    string `json:"3. low"`
		Close  string `json:"4. close"`
		Volume string `json:"5. volume"`
	}

	var raw RawDD
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Convert string fields to numeric values
	var err error
	o.Open, err = strconv.ParseFloat(raw.Open, 64)
	if err != nil {
		return fmt.Errorf("error parsing Open: %w", err)
	}
	o.High, err = strconv.ParseFloat(raw.High, 64)
	if err != nil {
		return fmt.Errorf("error parsing High: %w", err)
	}
	o.Low, err = strconv.ParseFloat(raw.Low, 64)
	if err != nil {
		return fmt.Errorf("error parsing Low: %w", err)
	}
	o.Close, err = strconv.ParseFloat(raw.Close, 64)
	if err != nil {
		return fmt.Errorf("error parsing Close: %w", err)
	}
	o.Volume, err = strconv.ParseInt(raw.Volume, 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing Volume: %w", err)
	}

	return nil
}
