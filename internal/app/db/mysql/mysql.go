package mysql

import (
	"backuper/config"
	"backuper/internal/util"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Mysql struct {
	Params        *config.DatabaseParams
	ExcludeTables []string
}

func New(params *config.Config) *Mysql {
	return &Mysql{Params: &params.Database, ExcludeTables: params.ExcludeTables}
}

func (M Mysql) Dump(filePath string) error {
	// Формируем команду с учетом сжатия
	cmdStr := fmt.Sprintf(
		"mysqldump -h %s -u %s -p%s %s %s | gzip > %s",
		M.Params.Host, M.Params.User, M.Params.Password, M.Params.Name, formatExcludeTables(M.ExcludeTables, M.Params.Name), filePath,
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

func (M Mysql) Restore(dumpFile string) error {
	loading := make(chan bool)
	defer close(loading)
	go util.Spinner(loading, fmt.Sprintf("Восстановление: %s", dumpFile))

	// Проверяем, существует ли файл дампа
	if _, err := os.Stat(dumpFile); os.IsNotExist(err) {
		return fmt.Errorf("файл дампа не найден: %s", dumpFile)
	}

	restoreCmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("gunzip -c %s | mysql -h %s -u %s -p%s %s",
			dumpFile, M.Params.Host, M.Params.User, M.Params.Password, M.Params.Name))

	// Выполняем команду
	var stdout, stderr bytes.Buffer
	restoreCmd.Stdout = &stdout
	restoreCmd.Stderr = &stderr

	if err := restoreCmd.Run(); err != nil {
		return fmt.Errorf("ошибка при восстановлении дампа: %s", stderr.String())
	}

	loading <- true
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
