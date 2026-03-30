package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/AnatolyKoltun/calculator-api/models"
	pb "github.com/AnatolyKoltun/calculator-api/proto"
)

// GetCalculationsWithGRPC — получает список вычислений через gRPC из storage
func GetCalculationsWithGRPC(c *gin.Context, client pb.StorageServiceClient) {
	var filter models.FilterRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметров: " + err.Error()})
		return
	}

	// Формируем gRPC запрос
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	req := &pb.ListRequest{
		DateFrom: string(filter.DateFrom),
		DateTo:   string(filter.DateTo),
	}

	resp, err := client.ListCalculations(ctx, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных из storage: " + err.Error()})
		return
	}

	// Конвертируем gRPC ответ в формат, ожидаемый фронтендом
	calculations := make([]models.Calculation, len(resp.Calculations))

	for i, calc := range resp.Calculations {
		calculations[i] = models.Calculation{
			ID:        int(calc.Id),
			Argument1: calc.Argument1,
			Argument2: calc.Argument2,
			Operator:  calc.Operator,
			Result:    calc.Result,
			CreatedAt: calc.CreatedAt.AsTime(),
		}
	}

	c.JSON(http.StatusOK, calculations)
}
