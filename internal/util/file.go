package util

import (
	"fmt"
	"os"
	"time"
)

// GetDumpFileName формирует имя файла и создает папку для хранения
func GetDumpFileName(dir string, gzip bool) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("ошибка создания директории: %w", err)
	}

	timeStamp := time.Now().Format("20060102_150405")
	var extension string = "sql"
	if gzip {
		extension = "sql.gz"
	}

	return fmt.Sprintf("%s/backup_%s.%s", dir, timeStamp, extension), nil
}
