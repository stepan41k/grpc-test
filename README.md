# USDT Rates gRPC Service

Микросервис на Go для получения и обработки курсов USDT с биржи Grinex. Сервис предоставляет gRPC интерфейс для получения текущих котировок, расчета статистических показателей и дальнейшего автоматического сохранения данных в PostgreSQL.

## Основной функционал

- **gRPC API**: Метод GetRates для получения цен Ask/Bid и timestamp.
    
- **Статистика**: Расчет значения конкретной позиции (topN) и среднего значения в диапазоне (avgNM).
    
- **Хранение данных**: Автоматическое сохранение каждой полученной котировки в базу данных PostgreSQL.
    
- **Observability**:
    
    - Трассировка запросов через **OpenTelemetry** и **Jaeger**.
        
    - Сбор метрик (gRPC, Go Runtime, бизнес-метрики) через **Prometheus**.
        
- **Healthcheck**: Встроенный механизм проверки работоспособности сервиса.
    
- **Надежность**: Реализован **Graceful Shutdown** для корректного закрытия соединений с БД, Jaeger и завершения активных запросов.


## Технологический стек

- **Language**: Go 1.25.7
    
- **Protocol**: gRPC (Protobuf)
    
- **Database**: PostgreSQL (pgx/v5(pgxpool) driver)
    
- **HTTP Client**: Resty (для работы с API Grinex)
    
- **Monitoring**: Prometheus
    
- **Tracing**: OpenTelemetry (OTLP) + Jaeger
    
- **Logging**: uber zap
    
- **Containerization**: Docker / Docker Compose

## Быстрый запуск

Самый простой способ запустить весь стек (приложение, Postgres, Jaeger, Prometheus) — использовать docker-compose:

codeBash

```
make run
```

После запуска сервисы будут доступны по следующим адресам:

- **gRPC Server**: localhost:50051
    
- **Prometheus Metrics**: http://localhost:9091/metrics
    
- **Jaeger UI (Traces)**: http://localhost:16686
    
- **Prometheus UI**: http://localhost:9090



## Конфигурация

Приложение поддерживает гибкую настройку через аргументы командной строки (флаги) и переменные окружения.

**Приоритет настроек:**

1. Флаги запуска (высший приоритет).
    
2. Переменные окружения (.env файл).
    
3. Значения по умолчанию.
    

### Параметры подключения

|                               |               |                             |                                                    |
| ----------------------------- | ------------- | --------------------------- | -------------------------------------------------- |
| Параметр                      | Флаг запуска  | Переменная окружения        | По умолчанию                                       |
| Хост Postgres                 | -pg-host      | PG_HOST                     | db                                                 |
| Порт подключения Postgres     | -pg-port      | PG_PORT                     | 5432                                               |
| Пользователь Postgres         | -pg-user      | PG_USERNAME                 | user                                               |
| Имя базы данных Postgres      | -pg-dbname    | PG_DBNAME                   | rates                                              | 
| Пароль Postgres               | -pg-password  | PG_PASSWORD                 | https://grinex.io                                  |
| SSL mode Postgres             | -pg-sslmode   | PG_SSLMODE                  | disable                                            |
| Окружение                     |       -       | ENV                         | jaeger:4317                                        |
| URL Grinex API                |       -       | GRINEX_URL                  | https://grinex.io                                  |
| gRPC port                     |       -       | GRPC_PORT                   | 50051                                              |
| Хост и порта экспорта трейсов |       -       | OTEL_EXPORTER               | jaeger:4317                                        |
| Порт метрик                   |       -       | METRICS_PORT                | 9090                                               |



## Команды Makefile

В корне проекта доступен Makefile для автоматизации основных задач:

- make build — Сборка исполняемого файла в папку bin/.
    
- make test — Запуск всех unit-тестов (бизнес-логика, расчеты, клиенты)._

- make cover — Запуск проверки покрытия тестами._
    
- make docker-build — Сборка Docker-образа приложения.
    
- make run — Запуск всего стека через docker-compose.

- make stop — Выполняет docker-compose down 
    
- make lint — Запуск линтера golangci-lint (требуется предустановленный golangci-lint v1.64.8).
    
- make proto — Генерация Go-кода из .proto файлов.

- make migrate — Генерация файлов миграции, которые состоят из текущей даты и введенного названия.



## Описание gRPC API

### Метод GetRates

Принимает индексы для расчета и возвращает текущий курс.

**Request:**

- top_n_index (int32): Индекс конкретной позиции в массиве.
    
- avg_n (int32): Начальный индекс диапазона для среднего.
    
- avg_m (int32): Конечный индекс диапазона для среднего.
    

**Response:**

- top_n_price (double): Цена на позиции N.
    
- avg_nm_price (double): Средняя цена в диапазоне [N; M].
    
- timestamp (string): Таймстемп запроса к API (тип protobuf.timestamp).
    

### Метод Check

Стандартный Healthcheck. Возвращает статус SERVING, если сервис готов к работе.



## Разработка и тестирование

### Unit-тесты

Покрывают ключевой функционал:

- Математические расчеты диапазона и выборки.
    
- Мокирование внешних зависимостей (БД и API) для проверки бизнес-логики сервиса.
    
- Парсинг JSON-ответов от внешнего API.
    

### Линтер

Используется golangci-lint. Конфигурация находится в файле .golangci.yml. Проверяет стиль кода, потенциальные утечки памяти (bodyclose) и необработанные ошибки.

#### Линтеры:

- errcheck:      unhandled errors
- govet:         standard
- ineffassign:   finding unused assignments
- staticcheck:   checks for logical errors
- unused:        find unused constants, variables, and functions
- gocritic:      code quality tips
- bodyclose:     checks for closing of HTTP body
- sqlclosecheck: checks the closure of sql.Rows rows

---

## Дополнительная информация
    
- Для управления миграциями БД используется golang-migrate (запускается автоматически при старте через Docker).
    
- Реализован механизм **Graceful Shutdown**: при получении сигналов SIGINT/SIGTERM сервис дожидается завершения активных RPC, сбрасывает трейсы в Jaeger и корректно закрывает пул соединений с PostgreSQL.
