package services

import (
	"context"
	"errors"
	"time"

	"github.com/AnatolyKoltun/calculator-api/models"
	//"github.com/AnatolyKoltun/calculator-api/repositories"
)

type CalculationService struct{}

//func (service CalculationService) Count(data models.RequestBody) (models.Calculation, error) {
//	var result float64
//
//	switch data.Operator {
//	case "+":
//		result = data.Argument1 + data.Argument2
//	case "-":
//		result = data.Argument1 - data.Argument2
//	case "*":
//		result = data.Argument1 * data.Argument2
//	case "/":
//		if data.Argument2 == 0 {
//			return models.Calculation{}, errors.New("деление на ноль запрещено")
//		}
//		result = data.Argument1 / data.Argument2
//	default:
//		return models.Calculation{}, errors.New("неподдерживаемая операция")
//	}
//
//	calculation := models.Calculation{
//		Argument1: data.Argument1,
//		Argument2: data.Argument2,
//		Operator:  data.Operator,
//		Result:    result,
//		CreatedAt: time.Now(),
//	}
//
//	//errSave := calcRepository.Save(context.Background(), &calculation)
//
//	//if errSave != nil {
//	//	return models.Calculation{}, errors.New("Ошибка сохранения в БД: " + errSave.Error())
//	//}
//
//	return calculation, nil
//}

func (service CalculationService) GetCalculationsList(filter models.FilterRequest) ([]models.Calculation, error) {

	//calculations, err := calcRepository.GetList(context.Background(), filter)

	return calculations, err
}
