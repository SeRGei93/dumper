package main

import (
	"backuper/config"
	"backuper/internal/backup"
	"fmt"
	"log"
)

func main() {
	config.LoadConfig()
	dumpFile, err := backup.WithGO()
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return
	}

	log.Printf("Дамп успешно создан: %s", dumpFile)
}
