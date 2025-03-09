package minio

import (
	"backuper/config"
	"backuper/internal/util"
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
	cfg := config.Cfg
	client, err := createClient(&cfg)
	if err != nil {
		return nil, err
	}

	return &Minio{Client: client, Bucket: cfg.Minio.Bucket, MaxFiles: cfg.MaxFiles}, nil
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

	return nil
}

func (M Minio) DownloadLastObject(localDir string) (string, error) {
	loading := make(chan bool)
	defer close(loading)

	go util.Spinner(loading, "Загружаю последний дамп")

	objects, err := M.GetObjects()
	if err != nil {
		return "", err
	}

	fileName := objects[0].Key
	localPath := filepath.Join(localDir, filepath.Base(fileName))
	err = M.Client.FGetObject(context.Background(), M.Bucket, fileName, localPath, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln("Ошибка скачивания файла:", err)
	}

	loading <- true
	fmt.Println("✅ Файл успешно скачан в:", localPath)

	return localPath, nil
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

func createClient(params *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(params.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(params.Minio.AccessKey, params.Minio.SecretKey, ""),
		Secure: params.Minio.UseSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MinIO: %v", err)
	}

	return minioClient, nil
}
