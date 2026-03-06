# Variables
APP_NAME=prueba-go
DOCKER_COMPOSE=docker compose
GO=go

.PHONY: all build run test clean docker-up docker-down docker-build swagger help

all: build

## Construcción:
build: ## Compilar el binario de la aplicación
	@echo "Compilando binario..."
	CGO_ENABLED=0 $(GO) build -o bin/$(APP_NAME) cmd/main.go

run: ## Ejecutar la aplicación localmente
	@echo "Ejecutando localmente..."
	$(GO) run cmd/main.go

test: ## Ejecutar las pruebas (tests)
	@echo "Ejecutando pruebas..."
	$(GO) test -v ./...

clean: ## Eliminar artefactos de construcción
	@echo "Limpiando..."
	rm -rf bin/
	rm -rf vendor/

## Docker:
docker-up: ## Levantar contenedores de Docker
	@echo "Iniciando contenedores..."
	$(DOCKER_COMPOSE) up -d

docker-down: ## Detener y eliminar contenedores y volúmenes de Docker
	@echo "Deteniendo contenedores..."
	$(DOCKER_COMPOSE) down -v

docker-build: ## Construir imágenes de Docker
	@echo "Construyendo imágenes..."
	$(DOCKER_COMPOSE) build

docker-logs: ## Ver logs de los contenedores
	$(DOCKER_COMPOSE) logs -f

## Herramientas:
swagger: ## Generar documentación de Swagger
	@echo "Generando documentación de Swagger..."
	swag init -g cmd/main.go -o internal/docs

help: ## Mostrar esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0pt %s\n", $$1, $$2}'
