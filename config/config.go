package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"regexp"
)

type DatabaseParams struct {
	Name          string   `mapstructure:"name"`
	User          string   `mapstructure:"user"`
	Password      string   `mapstructure:"password"`
	Host          string   `mapstructure:"host"`
	ExcludeTables []string `mapstructure:"exclude_tables"`
}

type S3Params struct {
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

type Config struct {
	Database DatabaseParams `mapstructure:"database"`
	Minio    S3Params       `mapstructure:"minio"`
	Dir      string         `mapstructure:"dir"`
}

var Cfg Config

// LoadConfig загружает YAML и заменяет ${ENV_VAR}
func LoadConfig() {
	_ = godotenv.Load()
	viper.SetConfigFile("config.yaml")

	// Загружаем конфигурацию
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка загрузки config.yaml: %v", err)
	}

	// Получаем данные в map
	configMap := viper.AllSettings()
	replaceEnvVariables(configMap)

	// Пересобираем в Viper и Unmarshal в структуру
	err := viper.MergeConfigMap(configMap)
	if err != nil {
		log.Fatalf("Ошибка объединения конфигурации: %v", err)
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("Ошибка парсинга конфигурации: %v", err)
	}
}

// replaceEnvVariables заменяет ${ENV_VAR} значениями из переменных окружения
func replaceEnvVariables(configMap map[string]interface{}) {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	for key, value := range configMap {
		switch v := value.(type) {
		case string:
			configMap[key] = re.ReplaceAllStringFunc(v, func(match string) string {
				envVar := match[2 : len(match)-1] // Убираем `${}`
				if envValue, exists := os.LookupEnv(envVar); exists {
					return envValue
				}
				return match // Если ENV нет, оставляем как есть
			})
		case map[string]interface{}:
			replaceEnvVariables(v) // Рекурсивно обрабатываем вложенные структуры
		}
	}
}
