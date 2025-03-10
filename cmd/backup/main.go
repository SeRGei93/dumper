package main

import (
	"backuper/config"
	"backuper/internal/pkg/app"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	config.LoadConfig()
	dumperS3, err := app.New(&config.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	if config.BackupFlag == true {
		err = dumperS3.RunCreate()
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if config.RestoreFlag == true {
		err = dumperS3.RunRestore()
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	fmt.Println("Usage:")
	flag.PrintDefaults()
	os.Exit(1)
}
