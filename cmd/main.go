package main

import (
	"archiva/internal/frameworks"
	"archiva/internal/handlers"
	"archiva/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	port := "8080"

	// Загружаем .env файл
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Теперь получаем API ключ из переменной окружения
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Fatal("API ключ не найден")
	}

	// Создаем новый роутер
	router := frameworks.NewRouter()
	mailService := services.NewMailService(apiKey)

	// Регистрируем маршруты
	handlers.RegisterRoutes(router) // Для архивации
	mailHandler := handlers.NewMailHandler(mailService)
	router.Handle("POST", "/api/mail/file", mailHandler.SendMailHandler)

	// Запускаем сервер
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(server.ListenAndServe())
}
