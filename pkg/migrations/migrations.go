package migrations

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Run(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Читаем SQL файл
	sqlBytes, err := os.ReadFile("pkg/migrations/001_create_tables.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Выполняем миграцию
	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	fmt.Println("Migrations completed successfully")
	return nil
}
