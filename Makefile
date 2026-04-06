.PHONY: build test cover docker-build run stop lint proto migrate-create

APP_NAME = usdt-rate-service
PROTO_DIR = api/proto/exchange/v1
OUT_DIR = .
MIGRATIONS_DIR = migrations

build:
	go build -o bin/$(APP_NAME) cmd/server/main.go

test:
	go test -v ./internal/...

cover:
	go test -cover ./...

docker-build:
	docker build -t $(APP_NAME):latest .

run:
	docker-compose up --build

stop:
	docker-compose down -v

lint:
	golangci-lint run ./...

proto:
	@mkdir -p internal/grpc/pb
	protoc --proto_path=api/proto/exchange/v1 \
			--go_out=internal/grpc/pb --go_opt=paths=source_relative \
			--go-grpc_out=internal/grpc/pb --go-grpc_opt=paths=source_relative \
			api/proto/exchange/v1/exchange.proto
	@echo "Proto files generated successfully."

migrate-create:
	@read -p "Enter migration name: " name; \
	TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	touch $(MIGRATIONS_DIR)/$${TIMESTAMP}_$${name}.up.sql; \
	touch $(MIGRATIONS_DIR)/$${TIMESTAMP}_$${name}.down.sql; \
	echo "Created: $(MIGRATIONS_DIR)/$${TIMESTAMP}_$${name}.up.sql"; \
	echo "Created: $(MIGRATIONS_DIR)/$${TIMESTAMP}_$${name}.down.sql"
