// candy_store_server/internal/transport/handlers/model_order.go
package handlers

type Order struct {
	Money int `json:"money"` // колличество денег
	CandyType string `json:"candyType"` // тип конфеты
	CandyCount int `json:"candyCount"` // колличество конфет
}