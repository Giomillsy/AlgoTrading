//Entry point into Go applications

package main

import (
	"AlgorithmicTrading/Backend/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

var FTSE100 = api.Index{
	Name:   "The FTSE 100",
	Ticker: "FTSE100",
	Securities: []api.Security{
		{
			Name:        "Unilever",
			Ticker:      "UNI",
			SystemicRho: 0.25,
		},
	},
}

func getIndex(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, FTSE100)
}

func main() {

	api.ApiQuery()

	/*
		router := gin.Default()
		router.GET("/getIndex", getIndex)
		router.Run("localhost:8080")
	*/
}
