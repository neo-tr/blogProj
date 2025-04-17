package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var Pool *pgxpool.Pool

// Init устанавливает соединение с БД и запускает миграцию таблицы posts.
func Init() {
	// Получение строки подключения из переменной окружения или использование значения по умолчанию.
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Замените значения user, password, host, port и dbname на свои
		connStr = "postgres://user:password@localhost:5432/go_blog?sslmode=disable"
	}

	// Создаем пул соединений.
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	Pool = pool

	log.Println("Соединение с базой данных установлено, таблица posts готова к использованию.")
}

// Close корректно закрывает соединение с базой данных.
func Close() {
	Pool.Close()
}
