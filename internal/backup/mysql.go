package backup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"backuper/config"
	"backuper/internal/util"

	_ "github.com/go-sql-driver/mysql"
)

// Глобальный мьютекс для защиты файла от одновременной записи
var fileMutex sync.Mutex

// Количество одновременно работающих горутин
const maxGoroutines = 4

func WithGO() (string, error) {
	dbConfig := config.Cfg.Database

	// Подключаемся к базе
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", fmt.Errorf("ошибка подключения к базе: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	dumpFile, err := util.GetDumpFileName(config.Cfg.Dir, false)
	if err != nil {
		return "", err
	}

	// Создаем файл
	file, err := os.Create(dumpFile)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(fmt.Sprintf("-- Дамп базы данных %s\n", dbConfig.Name)))
	if err != nil {
		return "", err
	}

	// Дамп структуры таблиц
	tables, err := getTables(db)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения имени таблицы: %v", err)
	}

	// счетчик потоков, чтобы дождаться остановки всех
	var wg sync.WaitGroup
	// канал для ошибок
	errCh := make(chan error, len(tables))

	// Создаём семафор с буфером 4
	sem := make(chan struct{}, maxGoroutines)
	for _, tableName := range tables {
		fmt.Printf("Дамп таблицы: %s\n", tableName)
		wg.Add(1)
		go dumpTable(tableName, db, file, dbConfig.ExcludeTables, &wg, errCh, sem)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		fmt.Println("Ошибка:", err)
	}

	log.Printf("Дамп базы данных создан: %s", dumpFile)
	return dumpFile, nil
}

func dumpTable(tableName string, db *sql.DB, file *os.File, excludeTables []string, wg *sync.WaitGroup, errCh chan error, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{}        // Захватываем слот в семафоре
	defer func() { <-sem }() // Освобождаем слот

	schema, err := getTableStructure(tableName, db)
	if err != nil {
		errCh <- fmt.Errorf("ошибка при получении структуры %s: %v", tableName, err)
		return
	}
	var tableData string = ""
	if contains(excludeTables, tableName) == false {
		err = getTableData(tableName, db, &tableData)
		if err != nil {
			errCh <- err
		}
	}

	// Блокируем файл перед записью
	fileMutex.Lock()
	_, err = file.Write([]byte(fmt.Sprintf("%s;\n\n%s", schema, tableData)))
	fileMutex.Unlock() // Разблокируем файл
	if err != nil {
		errCh <- err
	}
}

// getTableData - получаем данные
func getTableData(tableName string, db *sql.DB, tableData *string) error {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM `%s`", tableName))
	if err != nil {
		return fmt.Errorf("ошибка получения данных из таблицы %s: %v", tableName, err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	//var tableData string = ""
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("ошибка чтения строки из %s: %v", tableName, err)
		}

		var tableValues string
		tableValues = fmt.Sprintf("INSERT INTO `%s` VALUES (", tableName)

		for i, val := range values {
			if val == nil {
				tableValues += "NULL"
			} else {
				tableValues += fmt.Sprintf("'%v'", val)
			}
			if i < len(values)-1 {
				tableValues += ", "
			}
		}
		tableValues += ");\n"

		*tableData += tableValues
	}

	return nil
}

// getTableStructure - получаем структуру таблицы
func getTableStructure(tableName string, db *sql.DB) (string, error) {
	var schema string
	err := db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)).Scan(&tableName, &schema)
	if err != nil {
		return "", fmt.Errorf("ошибка получения структуры таблицы %s: %v", tableName, err)
	}

	return fmt.Sprintf("\n-- Структура таблицы `%s`\n%s", tableName, schema), nil
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// TODO: handle error
		}
	}(rows)

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// contains проверяет, содержится ли строка в массиве
func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// formatExcludeTables преобразует список исключаемых таблиц в строку аргументов для mysqldump
func formatExcludeTables(tables []string, dbName string) string {
	var excludeArgs string
	for _, table := range tables {
		excludeArgs += fmt.Sprintf(" --ignore-table=%s.%s", dbName, table)
	}
	return excludeArgs
}
