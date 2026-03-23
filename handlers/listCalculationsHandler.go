package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AnatolyKoltun/calculator-api/models"
	"github.com/AnatolyKoltun/calculator-api/services"
)

func GetCalculations(c *gin.Context) {
	var filter models.FilterRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметров: " + err.Error()})
		return
	}

	calculations, errList := services.GetCalculationsList(filter)

	if errList != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errList.Error()})
		return
	}

	c.JSON(http.StatusOK, calculations)
}
