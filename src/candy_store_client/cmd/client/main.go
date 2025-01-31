// candy_store_client/cmd/client/main.go

// candy_store_client/
// ├── cmd/
// │   └── client/
// │       └── main.go
// └── cert/
//     └── client/
//     │   └── cert.pem
//     │   └── key.pem
//     └── minica.pem

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	partCert   = "../../cert/client/cert.pem"  // путь к сертификату клиента
	partKey    = "../../cert/client/key.pem"   // путь к ключу клиента
	partMinica = "../../cert/minica.pem"       // путь к CA-сертификату
)

// Order - структура, которая определяет тело запроса для покупки конфет
type Order struct {
	Money      int    `json:"money"`      // колличество денег
	CandyType  string `json:"candyType"`  // тип конфет
	CandyCount int    `json:"candyCount"` // колличество конфет
}

type CommonResponse struct {
	Change int    `json:"change"`
	Error_ string `json:"error"`
	Thanks string `json:"thanks"`
}

// читаем и декодируем ответ от сервера
func handleResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close() // гарантируем закрытие тела ответа после чтения
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var commonResponse CommonResponse
	if err := json.Unmarshal(body, &commonResponse); err != nil {
		return "", err
	}

	if commonResponse.Error_ != "" {
		return commonResponse.Error_, nil
	}

	if commonResponse.Thanks != "" {
		return fmt.Sprintf("Спасибо! Ваша сдача составляет %v\n%s", strconv.Itoa(commonResponse.Change), commonResponse.Thanks), nil
	}

	return "", err
}

func createRequest(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body)) // создаём новый POST-запрос
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json") // устанавливаем загловок Content-Type
	return req, nil                                   // возвращаем созданный запрос
}

func createClient() (*http.Client, error) {
	// чтение клиентского сертификата
	cert, err := os.ReadFile(partCert)
	if err != nil {
		log.Fatalf("Failed to read cert: %v\n", err)
		return nil, err
	}

	// чтение клиентского ключа
	key, err := os.ReadFile(partKey)
	if err != nil {
		log.Fatalf("Failed to read key: %v\n", err)
		return nil, err
	}

	// загружаем сертификат и ключ клиента
	clientCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v\n", err)
		return nil, err
	}

	// чтение CA-сертификата
	caCert, err := os.ReadFile(partMinica)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v\n", err)
		return nil, err
	}

	// создаем пул сертификатов и добавляем в него CA-сертификат
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// настраиваем TLS-сертификацию клиента
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}

	// создаем HTTP-клиент с поддержкой HTTPS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return client, nil
}

// отправляем запрос и обрабатываем ответ от сервера
func sendOrder(url string, order Order) (string, error) {
	requestBody, err := json.Marshal(order) // кодируем структуру в JSON
	if err != nil {
		return "", err
	}

	req, err := createRequest(url, requestBody) // создаём новый HTTP запрос
	if err != nil {
		return "", err
	}

	client, err := createClient() // создаём новый HTTPS-клиент
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req) // выполняем HTTP-запрос
	if err != nil {
		return "", err
	}

	return handleResponse(resp)
}

func parseFlags(candyType *string, candyCount, money *int) {
	flag.StringVar(candyType, "k", "", "type of candy")
	flag.IntVar(candyCount, "c", 0, "number of candies")
	flag.IntVar(money, "m", 0, "sum of money")

	flag.Parse()

	if *candyType == "" || *candyCount == 0 || *money == 0 {
		log.Fatal("Пожалуйста, укажите все необходимые флаги: -k, -c, -m")
	}
}

func main() {
	// URL сервера к которому будет отправлен запрос
	url := "https://candy.tld:3333/buy_candy"
	var candyType string
	var candyCount, money int
	parseFlags(&candyType, &candyCount, &money)

	// данные для запроса
	order := Order{
		Money:      money,
		CandyType:  candyType,
		CandyCount: candyCount,
	}

	// отправляем запрос и обрабатываем ответ
	change, err := sendOrder(url, order)
	if err != nil {
		log.Fatalf("Failed to send order: %v", err)
	}

	// логируем ответ от сервера
	log.Println(change)
}
