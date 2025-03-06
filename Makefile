APP_NAME = backuper
BUILD_DIR = build

GOOS_MAC = darwin
GOOS_LINUX = linux
GOARCH = amd64

LDFLAGS = -s -w  # Убираем отладочную информацию (уменьшает размер бинарника)
GCFLAGS =        # Можно добавить доп. флаги для оптимизации

build-mac:
	@echo "🔨 Компиляция для macOS..."
	GOOS=$(GOOS_MAC) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-local-mac ./cmd/local/main.go
	GOOS=$(GOOS_MAC) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-s3-mac ./cmd/s3/main.go
	@echo "✅ Собрано: $(BUILD_DIR)/$(APP_NAME)-mac"

build-linux:
	@echo "🔨 Компиляция для Linux..."
	GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-local ./cmd/local/main.go
	GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-s3 ./cmd/s3/main.go
	@echo "✅ Собрано: $(BUILD_DIR)/$(APP_NAME)-linux"

build-all: build-mac build-linux

clean:
	@echo "🗑 Удаление бинарников..."
	rm -rf $(BUILD_DIR)
	@echo "✅ Очистка завершена."

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

rebuild: clean build-all