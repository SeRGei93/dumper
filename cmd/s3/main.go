package main

import (
	"backuper/config"
	"backuper/internal/backup"
	"backuper/internal/storage"
	"fmt"
	"log"
	"os"
)

func main() {
	config.LoadConfig()
	dumpFile, err := backup.WithMysqlDump()
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		_ = os.Remove(dumpFile)
		return
	}

	log.Printf("Дамп успешно создан: %s", dumpFile)

	err = storage.UploadToMinio(dumpFile)
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
	}

	_ = os.Remove(dumpFile)
}
