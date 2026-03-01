.PHONY: build run test clean lint help

# Переменные
BINARY_NAME=drone-simulation
CMD_PATH=./cmd/app
GO=go
GOFLAGS=-v

# Сборка бинарного файла
build:
	@echo "Сборка $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Сборка завершена: ./$(BINARY_NAME)"

# Сборка с оптимизациями
build-release:
	@echo "Сборка релизной версии..."
	$(GO) build -ldflags="-s -w" -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Релизная сборка завершена: ./$(BINARY_NAME)"

# Кросс-компиляция для macOS (Apple Silicon M1/M2/M3)
# Примечание: требует установленных Xcode Command Line Tools и библиотек GLFW на целевой системе
build-darwin-arm64:
	@echo "Сборка для macOS (ARM64)..."
	@echo "ВНИМАНИЕ: Для сборки требуется macOS с установленными Xcode CLI и GLFW"
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 CC=o64-clang CXX=o64-clang++ $(GO) build -ldflags="-s -w" -o $(BINARY_NAME)-darwin-arm64 $(CMD_PATH)
	@echo "Сборка завершена: ./$(BINARY_NAME)-darwin-arm64"

# Кросс-компиляция для macOS (Intel)
build-darwin-amd64:
	@echo "Сборка для macOS (AMD64)..."
	@echo "ВНИМАНИЕ: Для сборки требуется macOS с установленными Xcode CLI и GLFW"
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-apple-darwin20.2-clang CXX=x86_64-apple-darwin20.2-clang++ $(GO) build -ldflags="-s -w" -o $(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	@echo "Сборка завершена: ./$(BINARY_NAME)-darwin-amd64"

# Кросс-компиляция для macOS (универсальный бинарник)
build-darwin-universal: build-darwin-arm64 build-darwin-amd64
	@echo "Создание универсального бинарника для macOS..."
	lipo -create -output $(BINARY_NAME)-darwin-universal $(BINARY_NAME)-darwin-arm64 $(BINARY_NAME)-darwin-amd64
	@echo "Универсальный бинарник: ./$(BINARY_NAME)-darwin-universal"

# Сборка для macOS БЕЗ визуализации (только CLI режим)
build-darwin-arm64-cli:
	@echo "Сборка для macOS (ARM64) без визуализации..."
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GO) build -ldflags="-s -w" -tags=no_vis -o $(BINARY_NAME)-darwin-arm64-cli $(CMD_PATH)
	@echo "Сборка завершена: ./$(BINARY_NAME)-darwin-arm64-cli"

build-darwin-amd64-cli:
	@echo "Сборка для macOS (AMD64) без визуализации..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags="-s -w" -tags=no_vis -o $(BINARY_NAME)-darwin-amd64-cli $(CMD_PATH)
	@echo "Сборка завершена: ./$(BINARY_NAME)-darwin-amd64-cli"

# Кросс-компиляция для Windows
build-windows:
	@echo "Сборка для Windows..."
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BINARY_NAME).exe $(CMD_PATH)
	@echo "Сборка завершена: ./$(BINARY_NAME).exe"

# Кросс-компиляция для Linux
build-linux:
	@echo "Сборка для Linux..."
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BINARY_NAME)-linux $(CMD_PATH)
	@echo "Сборка завершена: ./$(BINARY_NAME)-linux"

# Сборка для всех платформ
build-all: build-release build-darwin-arm64 build-darwin-amd64 build-windows build-linux
	@echo "Сборка для всех платформ завершена"

# Запуск симуляции
run:
	@echo "Запуск симуляции..."
	$(GO) run $(CMD_PATH)

# Запуск без визуализации
run-cli:
	@echo "Запуск симуляции без визуализации..."
	$(GO) run $(CMD_PATH) -no-vis

# Запуск в режиме отладки
run-debug:
	@echo "Запуск симуляции в режиме отладки..."
	$(GO) run $(CMD_PATH) -debug

# Запуск с конфигурацией
run-config:
	@echo "Запуск симуляции с конфигурацией..."
	$(GO) run $(CMD_PATH) -config=config.json

# Запуск тестов
test:
	@echo "Запуск тестов..."
	$(GO) test $(GOFLAGS) -race ./...

# Запуск тестов с покрытием
test-coverage:
	@echo "Запуск тестов с покрытием..."
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Отчёт о покрытии: coverage.html"

# Очистка
clean:
	@echo "Очистка..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "Очистка завершена"

# Форматирование кода
fmt:
	@echo "Форматирование кода..."
	$(GO) fmt ./...

# Линтинг (требуется golangci-lint)
lint:
	@echo "Запуск линтера..."
	@which golangci-lint > /dev/null || (echo "Установите golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

# Обновление зависимостей
deps:
	@echo "Обновление зависимостей..."
	$(GO) get -u ./...
	$(GO) mod tidy

# Генерация конфигурации по умолчанию
gen-config:
	@echo "Генерация конфигурации по умолчанию..."
	@printf '%s\n' '{' \
		'  "drone_count": 100,' \
		'  "min_distance": 5.0,' \
		'  "max_distance": 10.0,' \
		'  "update_interval": 2000000000,' \
		'  "formation_radius": 7.5,' \
		'  "movement_variation": 0.5,' \
		'  "smoothing_factor": 0.2,' \
		'  "debug": false,' \
		'  "visualization_enabled": true,' \
		'  "window_width": 800,' \
		'  "window_height": 600' \
		'}' > config.json
	@echo "Конфигурация сохранена в config.json"

# Установка
install:
	@echo "Установка..."
	$(GO) install $(GOFLAGS) $(CMD_PATH)
	@echo "Установка завершена"

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make build              - Сборка бинарного файла"
	@echo "  make build-release      - Сборка релизной версии"
	@echo "  make build-darwin-arm64 - Сборка для macOS (Apple Silicon)"
	@echo "  make build-darwin-amd64 - Сборка для macOS (Intel)"
	@echo "  make build-darwin-universal - Универсальный бинарник для macOS"
	@echo "  make build-darwin-arm64-cli - macOS (ARM64) без визуализации (CGO=0)"
	@echo "  make build-darwin-amd64-cli - macOS (AMD64) без визуализации (CGO=0)"
	@echo "  make build-windows      - Сборка для Windows"
	@echo "  make build-linux        - Сборка для Linux"
	@echo "  make build-all          - Сборка для всех платформ"
	@echo "  make run                - Запуск симуляции"
	@echo "  make run-cli            - Запуск без визуализации"
	@echo "  make run-debug          - Запуск в режиме отладки"
	@echo "  make run-config         - Запуск с конфигурацией"
	@echo "  make test               - Запуск тестов"
	@echo "  make test-coverage      - Запуск тестов с отчётом"
	@echo "  make clean              - Очистка"
	@echo "  make fmt                - Форматирование кода"
	@echo "  make lint               - Линтинг"
	@echo "  make deps               - Обновление зависимостей"
	@echo "  make gen-config         - Генерация config.json"
	@echo "  make install            - Установка"
	@echo "  make help               - Эта справка"
