package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AnatolyKoltun/calculator-api/services"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/AnatolyKoltun/calculator-api/models"
)

// CreateCalculationWithNATS — публикует задачу в NATS вместо прямого вычисления
func CreateCalculationWithNATS(c *gin.Context, js nats.JetStreamContext) {
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

	// Публикуем сообщение в NATS
	data, err := json.Marshal(calculation)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка формирования сообщения"})
		return
	}

	_, err = js.Publish("calculations.create", data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки в очередь"})
		return
	}

	// Возвращаем 202 Accepted, задача принята в обработку
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "processing",
		"message": "Задача принята",
	})
}
