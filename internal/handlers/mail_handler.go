package handlers

import (
	"archiva/internal/frameworks"
	"archiva/internal/models"
	"archiva/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type MailHandler struct {
	MailService *services.MailService
}

func NewMailHandler(mailService *services.MailService) *MailHandler {
	return &MailHandler{MailService: mailService}
}

func (h *MailHandler) SendMailHandler(w http.ResponseWriter, r *http.Request) {
	var emailRequest models.EmailRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&emailRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при разборе запроса: %v", err), http.StatusBadRequest)
		return
	}

	err = h.MailService.SendFile(emailRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при отправке файла: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Файл успешно отправлен"))
}

func (h *MailHandler) RegisterRoutes(r *frameworks.Router) {
	r.Handle("POST", "/api/mail/file", h.SendMailHandler)
}
