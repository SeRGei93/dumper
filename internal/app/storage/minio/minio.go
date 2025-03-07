package minio

import (
	"backuper/config"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	Client   *minio.Client
	Bucket   string
	MaxFiles int64
}

func New() (*Minio, error) {
	cfg := config.Cfg.Minio
	client, err := createClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Minio{Client: client, Bucket: cfg.Bucket, MaxFiles: cfg.MaxFiles}, nil
}

// Upload загружает файл в MinIO
func (M Minio) Upload(filePath string) error {
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
	_, err = M.Client.PutObject(context.Background(), M.Bucket, fileName, file, fileStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("ошибка загрузки в MinIO: %v", err)
	}

	log.Printf("Файл успешно загружен в MinIO: %s/%s", M.Bucket, fileName)
	return nil
}

func (M Minio) DownloadLastObject(localDir string) (string, error) {
	objects, err := M.GetObjects()
	if err != nil {
		return "", err
	}

	err = M.Client.FGetObject(context.Background(), M.Bucket, objects[0].Key, localDir, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln("Ошибка скачивания файла:", err)
	}

	path := fmt.Sprintf("./%s/%s", localDir, objects[0].Key)

	fmt.Println("✅ Файл успешно скачан в:", path)

	return path, nil
}

func (M Minio) GetObjects() ([]minio.ObjectInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := M.Client.ListObjects(ctx, M.Bucket, minio.ListObjectsOptions{
		Prefix:    "backup_",
		Recursive: true,
	})

	files := make([]minio.ObjectInfo, 0, M.MaxFiles+5)
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		files = append(files, object)
	}

	// сортируем от нового к старому
	sort.Slice(files, func(i, j int) bool {
		return files[i].LastModified.After(files[j].LastModified)
	})

	return files, nil
}

func (M Minio) RemoveObjects(objects []minio.ObjectInfo) error {
	objectCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectCh)
		for _, obj := range objects {
			objectCh <- obj
		}
	}()

	for rErr := range M.Client.RemoveObjects(context.Background(), M.Bucket, objectCh, minio.RemoveObjectsOptions{}) {
		if rErr.Err != nil {
			fmt.Println("Ошибка удаления:", rErr.Err)
		}
	}

	return nil
}

func createClient(params config.S3Params) (*minio.Client, error) {
	minioClient, err := minio.New(params.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(params.AccessKey, params.SecretKey, ""),
		Secure: params.UseSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MinIO: %v", err)
	}

	return minioClient, nil
}
