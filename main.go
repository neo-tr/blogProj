package main

import (
	"blogProj/db"

	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Ошибка загрузки .env файла")
	}

	// Инициализация бд
	db.Init()

	// Отложим закрытие при нормальном завершении
	defer db.Close()

	// Ловим сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Получен сигнал завершения, закрываем базу данных...")
		db.Close()
		os.Exit(0)
	}()
}
