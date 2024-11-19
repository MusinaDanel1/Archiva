package main

import (
	"archiva/internal/frameworks"
	"archiva/internal/handlers"
	"archiva/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8080"

	// Создаем новый роутер
	router := frameworks.NewRouter()
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Fatal("API ключ не найден")
	}
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
