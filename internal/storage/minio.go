package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"backuper/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// UploadToMinio загружает файл в MinIO
func UploadToMinio(filePath string) error {
	minioConfig := config.Cfg.Minio

	// Создаём клиент MinIO
	minioClient, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure: minioConfig.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("ошибка подключения к MinIO: %v", err)
	}

	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			_ = fmt.Errorf("ошибка закрытия файла: %v", err)
		}
	}(file)

	// Получаем размер файла
	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("ошибка получения информации о файле: %v", err)
	}

	// Имя файла в MinIO (без локального пути)
	fileName := filepath.Base(filePath)

	// Загружаем файл в MinIO
	_, err = minioClient.PutObject(context.Background(), minioConfig.Bucket, fileName, file, fileStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("ошибка загрузки в MinIO: %v", err)
	}

	log.Printf("Файл успешно загружен в MinIO: %s/%s", minioConfig.Bucket, fileName)
	return nil
}
