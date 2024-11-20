//The structures used in the data models for the API

package api

type Index struct {
	Name       string     `json:"name"`
	Ticker     string     `json:"ticker"`
	Securities []Security `json:"securities"`
}

type Security struct {
	Name        string  `json:"name"`
	Ticker      string  `json:"ticker"`
	SystemicRho float64 `json:"systemicRho"`
}
