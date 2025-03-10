APP_NAME = backuper
BUILD_DIR = build

GOOS_MAC = darwin
GOOS_LINUX = linux
GOARCH = amd64

LDFLAGS = -s -w  # Убираем отладочную информацию (уменьшает размер бинарника)
GCFLAGS =        # Можно добавить доп. флаги для оптимизации

build-linux:
	@echo "🔨 Компиляция для Linux..."
	GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/backup/main.go
	@echo "✅ Собрано: $(BUILD_DIR)/$(APP_NAME)-linux"

clean:
	@echo "🗑 Удаление бинарников..."
	rm -rf $(BUILD_DIR)
	@echo "✅ Очистка завершена."

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

rebuild: clean build-linux