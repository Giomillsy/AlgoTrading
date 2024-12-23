package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestApiQuery(t *testing.T) {
	/*
		Test cases:
		- Happy path
		- Invalid JSON output from API
		- No response from server
	*/

	// Test cases
	tests := []struct {
		name         string
		mockResponse string
		mockStatus   int
		mockError    error
		want         ApiResponse
		wantErr      bool
	}{
		{
			name: "Successful API Response",
			mockResponse: `{
                "Meta Data": {
                    "1. Information": "Daily Prices",
                    "2. Symbol": "AAPL",
                    "3. Last Refreshed": "2024-12-21",
                    "4. Output Size": "Compact",
                    "5. Time Zone": "US/Eastern"
                },
                "Time Series (Daily)": {
                    "2024-12-21": {
                        "1. open": "100",
                        "2. high": "200",
                        "3. low": "50",
                        "4. close": "124",
                        "5. volume": "1000"
                    }
                }
            }`,
			mockStatus: http.StatusOK,
			mockError:  nil,
			want: ApiResponse{
				MetaData: MetaData{
					Information:   "Daily Prices",
					Symbol:        "AAPL",
					LastRefreshed: time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
					OutputSize:    "Compact",
					TimeZone:      "US/Eastern",
				},
				TimeSeriesDaily: map[string]*DailyData{
					"2024-12-21": {
						Open:   100,
						High:   200,
						Low:    50,
						Close:  124,
						Volume: 1000,
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "Invalid JSON Response",
			mockResponse: `{"invalid": "json}`,
			mockStatus:   http.StatusOK,
			mockError:    nil,
			want:         ApiResponse{},
			wantErr:      true,
		},
		{
			name:         "Network Error",
			mockResponse: "",
			mockStatus:   0,
			mockError:    fmt.Errorf("network error"),
			want:         ApiResponse{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Mock the getApiResponse
			getApiResponse = func(secID string) (*http.Response, error) {

				// APIresponse threw an error
				if tt.mockError != nil {
					return nil, tt.mockError
				}

				// APIReponse didn't throw an error
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
				}, nil // No error
			}

			// Call ApiQuery
			got, err := ApiQuery("AAPL")
			if (err != nil) != tt.wantErr {
				t.Fatalf("ApiQuery() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Validate the entire struct
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApiQuery() got = %+v, want = %+v", got, tt.want)
			}
		})
	}

}

func TestReadApiKey(t *testing.T) {
	//  Tests readApiKey in API code

	t.Run("Testing reading APIKey", func(t *testing.T) {

		// Backup the original .env file
		originalEnv := ".env"
		backupEnv := ".env.bak"
		if _, err := os.Stat(originalEnv); err == nil {
			// The original env file has a status

			if err := os.Rename(originalEnv, backupEnv); err != nil {
				t.Fatalf("Failed to rename the original .env file: %v", err)
			}
			defer os.Rename(backupEnv, originalEnv) // Restore to orignal filename

		} else {
			t.Fatalf("Failed to get status of .env file : %v", err)
		}

		// Creates a tempory .env file
		fn := ".env"
		k := "APIKey"
		body := "ALPHA_VANTAGE_API_KEY=" + k
		err := os.WriteFile(fn, []byte(body), 0644)
		if err != nil {
			t.Fatalf("Failed to write .env file : %v", err)
		}

		defer os.Remove(fn) // Deltes tempory file after code finishes

		got, err := readAPIKey()
		if err != nil {
			t.Fatalf("unexpected error in readAPIKey: %v", err)
		}

		want := k
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
}

func TestAlphaQueryGen(t *testing.T) {
	// Tests generating a query in the format alpha vantage requires

	t.Run("Testing query generation", func(t *testing.T) {

		want := "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=IBM&outputsize=full&apikey=demo"

		qs := []string{
			"function=TIME_SERIES_DAILY",
			"symbol=IBM",
			"outputsize=full",
			"apikey=demo",
		}
		got := alphaQueryGen(qs)

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})

}
func TestMetaDataUnmarshalJSON(t *testing.T) {
	/*
		Two different test scenarios
		Happy path = A normal JSON which is expected
		Unexpected Parameter = Alpha vantage has a new parameter added to it's json
	*/

	tests := []struct {
		name     string
		jsonData string
		want     MetaData
	}{
		{
			name: "Happy Path",
			jsonData: `{
				"1. Information": "Daily Prices",
				"2. Symbol": "AAPL",
				"3. Last Refreshed": "2024-12-21",
				"4. Output Size": "Compact",
				"5. Time Zone": "US/Eastern"
			}`,
			want: MetaData{
				Information:   "Daily Prices",
				Symbol:        "AAPL",
				LastRefreshed: time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
				OutputSize:    "Compact",
				TimeZone:      "US/Eastern",
			},
		},
		{
			name: "Unexpected Parameter",
			jsonData: `{
				"1. Information": "Daily Prices",
				"2. Symbol": "AAPL",
				"3. Last Refreshed": "2024-12-21",
				"4. Output Size": "Compact",
				"5. Time Zone": "US/Eastern",
				"6. Unexpected": "Scary"
			}`,
			want: MetaData{
				Information:   "Daily Prices",
				Symbol:        "AAPL",
				LastRefreshed: time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
				OutputSize:    "Compact",
				TimeZone:      "US/Eastern",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got MetaData
			err := json.Unmarshal([]byte(tt.jsonData), &got)
			if err != nil {
				t.Fatalf("json.Unmarshal() error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("json.Unmarshal() got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestDailyDataUnmarshalJSON(t *testing.T) {
	/*
		Two different test scenarios
		Happy path = A normal JSON which is expected
		Unexpected Parameter = Alpha vantage has a new parameter added to it's json
	*/
	tests := []struct {
		name     string
		jsonData string
		want     DailyData
	}{
		{
			name: "Happy Path",
			jsonData: `{
				"1. open": "100",
				"2. high": "200",
				"3. low": "50",
				"4. close": "124",
				"5. volume": "1000"
			}`,
			want: DailyData{
				Open:   100,
				High:   200,
				Low:    50,
				Close:  124,
				Volume: 1000,
			},
		},
		{
			name: "Unexpected Parameter",
			jsonData: `{
				"1. open": "100",
				"2. high": "200",
				"3. low": "50",
				"4. close": "124",
				"5. volume": "1000",
				"6. velocity": "1"
			}`,
			want: DailyData{
				Open:   100,
				High:   200,
				Low:    50,
				Close:  124,
				Volume: 1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got DailyData
			err := json.Unmarshal([]byte(tt.jsonData), &got)
			if err != nil {
				t.Fatalf("json.Unmarshal() error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("json.Unmarshal() got: %v, want: %v", got, tt.want)
			}
		})
	}
}
