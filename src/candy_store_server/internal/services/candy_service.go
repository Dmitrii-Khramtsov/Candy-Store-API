// candy_store_server/internal/services
package services

import (
	"fmt"

	"github.com/lonmouth/candy_store_server/internal/repositories"
)

type CandyService struct {
	CandyRepo repositories.CandyRepository
}

func NewCandyService(CandyRepo repositories.CandyRepository) *CandyService {
	return &CandyService{CandyRepo: CandyRepo}
}

// BuyCandy покупает конфеты и возвращает сдачу или ошибку
func (cs *CandyService) BuyCandy(candyType string, money, candyCount int) (int, error) {
	// если клиент предоставил отрицательное значение candyCount или неправильный candyType (помните: все пять типов конфет кодируются двумя буквами, так что это либо "CE", "AA", "NT", "DE" или "YR", все остальные случаи считаются неверными), то сервер должен ответить 400 и ошибкой в JSON, описывающей, что пошло не так
	if candyCount <= 0 {
		return 0, fmt.Errorf("количество конфет не может быть отрицательным или равным нулю")
	}

	// если клиент предоставил отрицательное значение money
	if money < 0 {
		return 0, fmt.Errorf("количество денег не может быть отрицательным")
	}

	price, err := cs.CandyRepo.GetCandyPrice(candyType)
	if err != nil {
		return 0, err
	}

	totalPrice := price * candyCount

	// если сумма цен конфет меньше или равна сумме денег, которую покупатель дал автомату, сервер отвечает HTTP 201 и возвращает JSON с двумя полями — "thanks", которое говорит "thanks!", и "change", которое является суммой сдачи, которую машина должна вернуть клиенту
	if money >= totalPrice {
		change := money - totalPrice
		return change, nil
	}

	// если сумма больше предоставленной, сервер отвечает HTTP 402 и сообщением об ошибке в JSON, которое говорит "Вам нужно {amount} больше денег!", где {amount} — разница между предоставленной суммой и ожидаемой
	if money < totalPrice {
		neededAmount := totalPrice - money
		return 0, fmt.Errorf("вам нужно на %d больше денег", neededAmount)
	}

	return 0, fmt.Errorf("неизвестная ошибка")
}