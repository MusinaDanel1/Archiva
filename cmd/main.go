package main

import (
	"archiva/internal/frameworks"
	"archiva/internal/handlers"
	"archiva/internal/services"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := "8080"

	// Создаем новый роутер
	router := frameworks.NewRouter()
	mailService := services.NewMailService("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")

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
