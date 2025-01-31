// candy_store_server/internal/repositories/candy_data.go
package repositories

import (
	"errors"
	"log"
)

// интерфейс candyRepository определяет метод для получения цены конфеты по типу
type CandyRepository interface {
	GetCandyPrice(candy string) (int, error)
}

// структура хранит карту типов конфет и их цен
type CandyStorege struct {
	candies map[string]int
}

// функция newStorege инициализирует и отравляет новый экземпляр candyStorege с предопределёнными ценами на конфеты
func NewStorege() *CandyStorege {
	return &CandyStorege{
		candies: map[string]int{
			"CE": 10,
			"AA": 15,
			"NT": 17,
			"DE": 21,
			"YR": 23,
		},
	}
}

// метод GetCandyPrice возвращает цену указанного типа конфеты или ошибку, если тип конфеты не найден
func (cs *CandyStorege) GetCandyPrice(candyType string) (int, error){
	if candyPrice, exists := cs.candies[candyType]; exists {
		return candyPrice, nil
	}

	err := errors.New("тип конфет " + candyType + " не найден")
	log.Printf("ошибка: %v", err)
	return 0, err
}
