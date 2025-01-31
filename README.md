# Candy Store API

## Описание проекта

Проект "Candy Store" представляет собой API для магазина сладостей, который позволяет пользователям покупать конфеты через вендинговые автоматы. API реализован на языке Go и использует HTTPS для обеспечения безопасности передачи данных. Проект включает в себя серверную часть, которая обрабатывает запросы на покупку конфет, а также клиентскую часть, которая отправляет эти запросы. Проект также включает генерацию смешного сообщения с ASCII-коровой в ответе на успешную покупку.


## Структура проекта

### Сервер
```
candy_store_server/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── entities/
│   │   └── cow.c
│   │   └── cow.h
│   │   └── cow.go
│   ├── repositories/
│   │   └── candy_data.go
│   └── service/
│   │   └── candy_service.go
│   └── transport/
│       └── handlers/
│           └── api_default.go
│           └── model_inline_response_201.go
│           └── model_inline_response_400.go
│           └── model_order.go
├── pkg/
│   └── db/
│       └── store.go
├── cert/
│   └── server/
│   │   └── cert.pem
│   │   └── key.pem
│   └── minica.pem
├── go.mod
└── go.sum
```

### Клиент
```
candy_store_client/
├── cmd/
│   └── client/
│       └── main.go
└── cert/
    └── client/
    │   └── cert.pem
    │   └── key.pem
    └── minica.pem
```


## Шаги для запуска

### 1. Установка зависимостей

#### Go:
- Убедитесь, что установлен Go 1.21.6+.

### 2. Компиляция и запуск C-кода

Для запуска проекта необходимо выполнить следующие команды:
```
gcc -arch arm64 -c /Users/dmitrii/Go_Day04-1/src/candy_store_server/internal/entities/cow.c -o /Users/dmitrii/Go_Day04-1/src/candy_store_server/internal/entities/cow.o

gcc -arch arm64 -shared -o /Users/dmitrii/Go_Day04-1/src/candy_store_server/internal/entities/libcow.so /Users/dmitrii/Go_Day04-1/src/candy_store_server/internal/entities/cow.o
```

### 3. Запуск проекта

#### 1. Создание сертификатов:
- Используйте minica для генерации сертификатов.
- Сгенерируйте сертификаты для сервера и клиента, а также CA-сертификат.

#### Установка и использование minica:
1. Установка minica:
- Скачайте и установите minica с официального репозитория.
- Следуйте инструкциям по установке для вашей операционной системы.

2. Генерация сертификатов:
- Создайте директорию для хранения сертификатов:
```
mkdir -p cert/server cert/client
```

- Сгенерируйте сертификаты для сервера:
```
minica --domains candy.tld
```

- Сгенерируйте сертификаты для клиента:
```
minica --domains client.tld

```

- Переместите сгенерированные сертификаты в соответствующие директории:
```
mv candy.tld/cert.pem cert/server/
mv candy.tld/key.pem cert/server/
mv client.tld/cert.pem cert/client/
mv client.tld/key.pem cert/client/
mv minica.pem cert/
```

#### 2. Сборка и запуск сервера:
```
GOARCH=arm64 go build -o candy_store_server cmd/server/main.go
./candy_store_server
```

#### 3. Сборка и запуск клиента:
```
GOARCH=arm64 go build -o candy_store_client cmd/client/main.go
./candy_store_client -k AA -c 2 -m 50
```


## Описание API

### Маршруты
Статус
Метод: POST
Путь: /buy_candy
Назначение: Обрабатывает запросы на покупку конфет.

### Список баз данных
- В проекте используется внутренняя структура для хранения данных о конфетах и их ценах.

### Создание маппинга
- Маппинги создаются в функции main с использованием роутера mux.

### Выполнение запроса
- Запросы выполняются с использованием HTTPS и требуют наличия сертификатов клиента и сервера.

### Примечания

- Все данные передаются в формате JSON.
- Сервер проверяет корректность ввода и рассчитывает сдачу.

### Ошибки

- 400: Неверные данные ввода.
- 402: Недостаточно средств.
- 500: Внутренняя ошибка сервера.
- 415: Неподдерживаемый тип данных.
- 201: Успешная покупка.


## Интерфейс

- Клиентский интерфейс реализован в виде командной строки и поддерживает флаги -k, -c, -m для указания типа конфет, их количества и суммы денег соответственно.

## Пример использования

### Использование curl:
```
curl -v --key cert/client/key.pem --cert cert/client/cert.pem --cacert cert/minica.pem -X POST -H "Content-Type: application/json" -d '{"candyType": "NT", "candyCount": 2, "money": 34}' "https://candy.tld:3333/buy_candy"

```

### Использование клиентского приложения:
```
go run main.go -k AA -c 2 -m 36
```

## Заключение
- Проект "Candy Store" обеспечивает безопасную и надежную платформу для покупки конфет через вендинговые автоматы. Использование HTTPS и взаимной аутентификации TLS гарантирует защиту данных и предотвращает атаки типа "man-in-the-middle". Для управления сертификатами используется minica, чтобы сгенерировать сертификаты для сервера и клиента, а также CA-сертификат.
