package main

import (
	"backuper/config"
	"backuper/internal/pkg/app"
	"log"
)

func main() {
	config.LoadConfig()
	dumperS3, err := app.New(&config.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = dumperS3.Run()
	if err != nil {
		log.Fatal(err)
	}
}
