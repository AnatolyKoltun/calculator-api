package handlers

import (
	"net/http"

	"github.com/AnatolyKoltun/calculator-api/services"
	"github.com/gin-gonic/gin"

	"github.com/AnatolyKoltun/calculator-api/models"
)

func CreateCalculation(c *gin.Context) {
	var req models.RequestBody

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных: " + err.Error()})
		return
	}

	calculation, errCount := services.Calculate(req)

	if errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errCount.Error()})
		return
	}

	c.JSON(http.StatusOK, calculation)
}
