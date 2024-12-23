//Entry point into Go applications

package main

import (
	"Backend/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getSec(c *gin.Context) {

	IBM, err := api.ApiQuery("IBM")
	if err != nil {
		// For now be fatal for the API. Correct later for error handling
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, IBM)
}

func main() {

	router := gin.Default()
	router.GET("/getSec", getSec)
	router.Run("localhost:8080")

}
