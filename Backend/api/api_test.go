package api

import (
	"os"
	"testing"
)

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
