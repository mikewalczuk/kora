# Nazwa paczki z wygenerowanym kodem
API_PACKAGE=api
# Ścieżka do specyfikacji (relatywna do folderu backend)
OPENAPI_SPEC=../api-spec/openapi.yaml
# Plik wynikowy
API_OUT=internal/api/api.gen.go
GOOSE = $(shell go env GOPATH)/bin/goose
MIGRATIONS_DIR = ./migrations
GOOSE_FLAGS = -dir $(MIGRATIONS_DIR) postgres

.PHONY: gen tidy build dev install-tools

# Główna komenda generująca kod
gen:
	@echo "Generating OpenAPI code..."
	@mkdir -p internal/api
	$(shell go env GOPATH)/bin/oapi-codegen --config oapi-codegen.yaml $(OPENAPI_SPEC)
	@echo "Generating SQL queries..."
	sqlc generate
	@$(MAKE) tidy

# Porządkuje zależności po wygenerowaniu kodu
tidy:
	go mod tidy

# Buduje aplikację (sprawdza czy wszystko się kompiluje)
build: gen
	go build ./cmd/api

test:
	go test ./... -v

dev: gen migrate-up
	@set -a && . ./.env && set +a && $(shell go env GOPATH)/bin/air

migrate-up:
	@set -a && . ./.env && set +a && $(GOOSE) $(GOOSE_FLAGS) "$$DATABASE_URL" up

migrate-down:
	@set -a && . ./.env && set +a && $(GOOSE) $(GOOSE_FLAGS) "$$DATABASE_URL" down

install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest