package backup

import (
	"backuper/config"
	"backuper/internal/util"
	"fmt"
	"os"
	"os/exec"
)

func WithMysqlDump() (string, error) {
	params := config.Cfg.Database

	dumpFile, err := util.GetDumpFileName(config.Cfg.Dir, true)
	if err != nil {
		return "", err
	}

	// Формируем команду с учетом сжатия
	cmdStr := fmt.Sprintf(
		"mysqldump -h %s -u %s -p%s %s %s | gzip > %s",
		params.Host, params.User, params.Password, params.Name, formatExcludeTables(params.ExcludeTables, params.Name), dumpFile,
	)

	// Выполняем команду через sh -c (чтобы работал пайп)
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Env = os.Environ()

	// Запускаем команду
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения mysqldump: %v\nВывод: %s", err, string(output))
	}

	return dumpFile, nil
}
