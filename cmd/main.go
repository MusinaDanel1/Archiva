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

	router := frameworks.NewRouter()

	mailService := services.NewMailService(
		"smtp.gmail.com",         // SMTP-сервер
		"587",                    // Порт SMTP
		"danelmusina1@gmail.com", // Логин
		"Md310804@",              // Пароль
	)

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
