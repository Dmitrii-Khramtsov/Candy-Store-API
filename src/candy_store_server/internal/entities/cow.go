// candy_store_server/internal/entities/cow.go
package entities

/*
#cgo CFLAGS: -I/Users/dmitrii/Go_Day04-1/src/candy_store_server/internal/entities
#cgo LDFLAGS: -L/Users/dmitrii/Go_Day04-1/src/candy_store_server/internal/entities -lcow
#include "cow.h"
*/
import "C"
import "unsafe"

// AskCow вызывает ask_cow для генерации ASCII-коровы
func AskCow(phrase string) string {
	cPhrase := C.CString(phrase) // преобразуем строку Go в C-строку
	defer C.free(unsafe.Pointer(cPhrase)) // освобождаем память после использования

	cCow := C.ask_cow(cPhrase) // вызываем C-функцию ask_cow
	defer C.free(unsafe.Pointer(cCow)) // освобождаем память после использования

	return C.GoString(cCow) // преобразуем C-строку обратно в строку Go
}
