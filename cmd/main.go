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

	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Fatal("API ключ не найден")
	}

	router := frameworks.NewRouter()
	mailService := services.NewMailService(apiKey)

	handlers.RegisterRoutes(router)
	mailHandler := handlers.NewMailHandler(mailService)
	router.Handle("POST", "/api/mail/file", mailHandler.SendMailHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(server.ListenAndServe())
}
