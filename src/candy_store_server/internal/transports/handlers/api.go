// candy_store_server/internal/transport/handlers/api.go
package handlers

import (
	"encoding/json"
	"net/http"
	"errors"
	"strings"

	"github.com/lonmouth/candy_store_server/internal/services"
	"github.com/lonmouth/candy_store_server/pkg/db"
	"github.com/lonmouth/candy_store_server/internal/entities"
)

type CandyHandler struct {
	service *services.CandyService
}

func NewCandyHandler(service *services.CandyService) *CandyHandler {
	return &CandyHandler{service: service}
}

func BuyCandyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	
	var (
		model Order
		unmarshalerr *json.UnmarshalTypeError
	)

	decoder := json.NewDecoder(r.Body) // создаём новый декодер для чтения тела запроса
	decoder.DisallowUnknownFields() // запрещаем декодору игнорировать неизвестные поля
	err := decoder.Decode(&model) // декодируем JSON из тела запроса в структуру model
	if err != nil {
		if errors.As(err, &unmarshalerr) { // является ли ошибка - ошибкой типа
			errorResponse(w, "указан неправильный тип поля" +unmarshalerr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "ошибка декодирования" +err.Error(), http.StatusBadRequest)
		}
		return
	}

	candyRepo := db.NewStore()
	candyService := services.NewCandyService(candyRepo)

	change, err := candyService.BuyCandy(model.CandyType, model.Money, model.CandyCount)
	if err != nil {
		if strings.Contains(err.Error(), "нужно") {
			errorResponse(w, err.Error(), http.StatusPaymentRequired)
		} else if strings.Contains(err.Error(), "тип") || strings.Contains(err.Error(), "количество конфет не может быть отрицательным или равным нулю") {
			errorResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Генерируем ASCII-корову с сообщением "Thank you!"
	thanksMessage := entities.AskCow("Thank you!")
	response := InlineResponse201{
		Thanks: thanksMessage, // Используем сгенерированное сообщение с ASCII-коровой
		Change: int(change),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := InlineResponse400{
		Error_: message,
	}
	json.NewEncoder(w).Encode(resp)
}
