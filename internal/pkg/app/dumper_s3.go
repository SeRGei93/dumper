package app

import (
	"backuper/config"
	"backuper/internal/app/db/mysql"
	"backuper/internal/app/storage/minio"
	"backuper/internal/util"
	"fmt"
	"log"
	"os"
)

type DumperS3 struct {
	db   *mysql.Mysql
	s3   *minio.Minio
	dir  string
	file string
}

func New(cfg *config.Config) (*DumperS3, error) {
	db := mysql.New(cfg)
	s3, err := minio.New()
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к S3: %v", err)
	}

	dumpFile, err := util.GetDumpFileName(cfg.Dir, true)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &DumperS3{db, s3, cfg.Dir, dumpFile}, nil
}

func (d *DumperS3) RunCreate() error {
	err := d.Dump()
	if err != nil {
		return err
	}

	err = d.UploadToS3()
	if err != nil {
		return err
	}

	err = d.RemoveOldFiles()
	if err != nil {
		return err
	}

	// удаляем созданный дамп с диска
	err = d.Clean(d.file)
	if err != nil {
		return err
	}

	return nil
}

func (d *DumperS3) RunRestore() error {
	file, err := d.s3.DownloadLastObject(d.dir)
	if err != nil {
		return err
	}

	err = d.db.Restore(file)
	if err != nil {
		return err
	}

	// удаляем созданный дамп с диска
	err = d.Clean(file)
	if err != nil {
		return err
	}

	return nil
}

func (d *DumperS3) Dump() error {
	err := d.db.DumpWithSchemaOfExcludedTables(d.file)
	if err != nil {
		_ = os.Remove(d.file)
		return fmt.Errorf("ошибка при создании дампа: %v", err)
	}

	log.Printf("Дамп создан: %s", d.file)

	return nil
}

func (d *DumperS3) UploadToS3() error {
	err := d.s3.Upload(d.file)
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла на S3: %v", err)
	}
	log.Println("Файл успешно загружен на S3")

	return nil
}

// RemoveOldFiles удаляем старые дампы с S3
func (d *DumperS3) RemoveOldFiles() error {
	// удаляем старые бекапы
	objects, err := d.s3.GetObjects()
	if err != nil {
		return fmt.Errorf("ошибка при удалении старых файлов с S3: %v", err)
	}

	if int64(len(objects)) > d.s3.MaxFiles {
		err = d.s3.RemoveObjects(objects[d.s3.MaxFiles:])
		if err != nil {
			return fmt.Errorf("ошибка при удалении старых файлов с S3: %v", err)
		}

		fmt.Println("✅ Старые дампы удалены")
	}

	return nil
}

func (d *DumperS3) Restore(file string) error {
	err := d.db.Restore(file)
	if err != nil {
		return err
	}

	return nil
}

func (d *DumperS3) Clean(file string) error {
	err := os.Remove(file)
	if err != nil {
		return fmt.Errorf("ошибка при удалении файла с диска %v", err)
	}
	return nil
}
