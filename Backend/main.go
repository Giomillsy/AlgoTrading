//Entry point into Go applications

package main

import (
	"AlgorithmicTrading/Backend/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getSec(c *gin.Context) {

	IBM := api.ApiQuery("IBM")
	c.IndentedJSON(http.StatusOK, IBM)
}

func main() {

	router := gin.Default()
	router.GET("/getSec", getSec)
	router.Run("localhost:8080")

}
