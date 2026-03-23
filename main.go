package main

import (
	"github.com/gin-gonic/gin"

	"github.com/AnatolyKoltun/calculator-api/handlers"
)

func setupAndRunServer() {
	router := gin.Default()

	router.POST("/calculate", handlers.CreateCalculation)
	router.GET("/calculations", handlers.GetCalculations)

	router.Run(":8080")
}

func main() {
	setupAndRunServer()
}
