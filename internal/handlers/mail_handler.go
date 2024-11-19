package handlers

import (
	"archiva/internal/frameworks"
	"archiva/internal/models"
	"archiva/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type MailHandler struct {
	MailService *services.MailService
}

func NewMailHandler(mailService *services.MailService) *MailHandler {
	return &MailHandler{MailService: mailService}
}

// Обработчик отправки письма
func (h *MailHandler) SendMailHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Ограничиваем размер до 10 МБ
		http.Error(w, fmt.Sprintf("Ошибка при разборе формы: %v", err), http.StatusBadRequest)
		return
	}

	// Извлекаем файл из формы
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при извлечении файла: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Извлекаем список email из формы
	emails := r.FormValue("emails")
	if emails == "" {
		http.Error(w, "Список email-адресов не предоставлен", http.StatusBadRequest)
		return
	}

	// Разбиваем строки email в массив
	emailList := strings.Split(emails, ",")
	for i, email := range emailList {
		emailList[i] = strings.TrimSpace(email) // Убираем пробелы
	}

	// Создаем запрос на отправку email
	emailRequest := models.EmailRequest{
		FilePath:       fileHeader.Filename,
		Filename:       fileHeader.Filename,
		EmailAddresses: emailList,
	}

	// Записываем файл во временную директорию
	tempFile, err := os.CreateTemp("", "uploaded-*.tmp")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при сохранении файла: %v", err), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name()) // Удаляем временный файл после отправки

	// Копируем содержимое файла
	if _, err := io.Copy(tempFile, file); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при копировании файла: %v", err), http.StatusInternalServerError)
		return
	}

	emailRequest.FilePath = tempFile.Name() // Указываем путь к временной копии

	// Вызываем MailService для отправки файла
	if err := h.MailService.SendFile(emailRequest); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при отправке email: %v", err), http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Файл успешно отправлен на указанные email-адреса",
	})
}

// Регистрация маршрута
func (h *MailHandler) RegisterRoutes(r *frameworks.Router) {
	r.Handle("POST", "/api/mail/file", h.SendMailHandler)
}
