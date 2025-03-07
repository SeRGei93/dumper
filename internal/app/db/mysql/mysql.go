package mysql

import (
	"backuper/config"
	"fmt"
	"os"
	"os/exec"
)

type Mysql struct {
	Params *config.DatabaseParams
}

func New(params *config.DatabaseParams) *Mysql {
	return &Mysql{Params: params}
}

func (M Mysql) Dump(filePath string) error {

	// Формируем команду с учетом сжатия
	cmdStr := fmt.Sprintf(
		"mysqldump -h %s -u %s -p%s %s %s | gzip > %s",
		M.Params.Host, M.Params.User, M.Params.Password, M.Params.Name, formatExcludeTables(M.Params.ExcludeTables, M.Params.Name), filePath,
	)

	// Выполняем команду через sh -c (чтобы работал пайп)
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Env = os.Environ()

	// Запускаем команду
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка выполнения mysqldump: %v\nВывод: %s", err, string(output))
	}

	return nil
}

// formatExcludeTables преобразует список исключаемых таблиц в строку аргументов для mysqldump
func formatExcludeTables(tables []string, dbName string) string {
	var excludeArgs string
	for _, table := range tables {
		excludeArgs += fmt.Sprintf(" --ignore-table=%s.%s", dbName, table)
	}
	return excludeArgs
}
