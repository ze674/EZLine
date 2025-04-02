// internal/database/db.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// DB глобальная переменная для доступа к базе данных
var DB *sql.DB

// Connect устанавливает соединение с базой данных и выполняет миграции
func Connect(dbPath string) error {
	var err error

	// Убедимся, что директория существует
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию базы данных: %w", err)
	}

	// Открытие соединения с SQLite
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть соединение с базой данных: %w", err)
	}

	// Проверка соединения
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	log.Println("Соединение с базой данных установлено")

	// Запуск миграций
	if err = runMigrations(DB); err != nil {
		return fmt.Errorf("ошибка при выполнении миграций: %w", err)
	}

	return nil
}

// Выполняет миграции базы данных
func runMigrations(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver)
	if err != nil {
		return err
	}

	// Выполняем миграцию до последней версии
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Миграции выполнены успешно")
	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
