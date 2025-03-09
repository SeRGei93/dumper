package main

import (
	"backuper/config"
	"backuper/internal/pkg/app"
	"fmt"
	"log"
)

func main() {
	config.LoadConfig()
	dumperS3, err := app.New(&config.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = dumperS3.RunRestore()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Дамп успешно восстановлен")
}
