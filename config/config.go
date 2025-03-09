package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type DatabaseParams struct {
	Name     string
	User     string
	Password string
	Host     string
}

type S3Params struct {
	Endpoint  string
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

// структура для чтения файла .env
type envCfg struct {
	Name      string `env:"DB_NAME,required"`
	User      string `env:"DB_USER,required"`
	Password  string `env:"DB_PASSWORD,required"`
	Host      string `env:"DB_HOST,required"`
	Endpoint  string `env:"S3_HOST,required"`
	Bucket    string `env:"S3_BUCKET,required"`
	Region    string `env:"S3_REGION,required"`
	AccessKey string `env:"S3_ACCESS_KEY"`
	SecretKey string `env:"S3_SECRET_KEY"`
	UseSSL    bool   `env:"S3_USE_SSL"`
}

// структура для чтения файла yaml
type yamlConfig struct {
	App struct {
		Dir      string `yaml:"dir"`
		MaxFiles int64  `yaml:"max_files"`
	}
	Database struct {
		ExcludeTables []string `yaml:"exclude_tables"`
	}
}

type Config struct {
	Database      DatabaseParams
	ExcludeTables []string
	Minio         S3Params
	Dir           string
	MaxFiles      int64
}

var Cfg Config

func LoadConfig() {
	c, err := loadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	y, err := loadYaml()
	if err != nil {
		log.Fatal(err)
	}

	Cfg = Config{
		Database: DatabaseParams{
			Name:     c.Name,
			User:     c.User,
			Password: c.Password,
			Host:     c.Host,
		},
		ExcludeTables: y.Database.ExcludeTables,
		Minio: S3Params{
			Endpoint:  c.Endpoint,
			Bucket:    c.Bucket,
			Region:    c.Region,
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		},
		Dir:      y.App.Dir,
		MaxFiles: y.App.MaxFiles,
	}
}

func loadYaml() (*yamlConfig, error) {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	c := yamlConfig{}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func loadEnvironment() (*envCfg, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	// parse
	var c envCfg
	err = env.Parse(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
