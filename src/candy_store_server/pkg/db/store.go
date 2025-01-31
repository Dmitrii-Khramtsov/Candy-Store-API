// candy_store_server/pkg/db/store.go
package db

import "github.com/lonmouth/candy_store_server/internal/repositories"

// NewStore создает и возвращает новый экземпляр хранилища конфет, реализующего интерфейс candyRepository
func NewStore() repositories.CandyRepository {
	return repositories.NewStorege()
}