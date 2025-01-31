// candy_store_server/cmd/server/main.go

// candy_store_server/
// ├── cmd/
// │   └── server/
// │       └── main.go
// ├── internal/
// │   ├── entities/
// │   │   └── cow.c
// │   │   └── cow.h
// │   │   └── cow.go
// │   ├── repositories/
// │   │   └── candy_data.go
// │   └── service/
// │   │   └── candy_service.go
// │   └── transport/
// │       └── handlers/
// │           └── api_default.go
// │           └── model_inline_response_201.go
// │           └── model_inline_response_400.go
// │           └── model_order.go
// ├── pkg/
// │   └── db/
// │       └── store.go
// └── cert/
// |   └── server/
// |   │   └── cert.pem
// |   │   └── key.pem
// |   └── minica.pem
// ├── go.mod
// └── go.sum

package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux" // добавление библиотеки (go get -u github.com/gorilla/mux)
	"github.com/lonmouth/candy_store_server/internal/services"
	"github.com/lonmouth/candy_store_server/internal/transports/handlers"
	"github.com/lonmouth/candy_store_server/pkg/db"
)

var (
	partCert   = "cert/server/cert.pem" // путь к сертификату сервера
	partKey    = "cert/server/key.pem"  // путь к ключу сервера
	partMinica = "cert/minica.pem"      // путь к CA-сертификату
)

func main() {
	candyRepo := db.NewStore()
	candyService := services.NewCandyService(candyRepo)
	handlers.NewCandyHandler(candyService)

	router := mux.NewRouter() // создаем роутер и настраиваем маршруты
	router.HandleFunc("/buy_candy", handlers.BuyCandyHandler).Methods("POST")

	log.Println("The server is running on port 3333")

	// log.Fatal(http.ListenAndServe(":3333", router))
	// log.Fatal(http.ListenAndServeTLS(":3333", partCert, partKey, router))

	server := creatingServerHTTPS(router)

	// запускаем сервер с поддержкой HTTPS
	log.Fatal(server.ListenAndServeTLS(partCert, partKey))
}

func creatingServerHTTPS(router *mux.Router) *http.Server {
	// читаем CA-сертификат
	caCert, err := os.ReadFile(partMinica)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	// создаём пул сертификатов и добавляем в него CA-сертификат
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// настраиваем TLS-сертификацию
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool, // устанавливаем пул CA-сертификатов
		ClientAuth: tls.RequireAndVerifyClientCert, // требуем и проверяем сертификат у клиента
		Certificates: []tls.Certificate{loadServerCertificate()}, // загружаем сертификат сервера
	}

	// создаём HTTPS-server с настроенной TLS-сертификацией
	server := &http.Server{
		Addr: ":3333",
		TLSConfig: tlsConfig,
		Handler: router, // устанавливаем роутер в качестве обработчика запросов
	}

	return server
}

// функция для загрузки сертификата сервера
func loadServerCertificate() tls.Certificate {
	cert, err := tls.LoadX509KeyPair(partCert, partKey)
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}
	return cert
}
