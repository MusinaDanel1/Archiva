package services

import (
	"archiva/internal/models"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"os"
)

type MailService struct {
	SMTPServer   string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
}

func NewMailService(smtpServer, smtpPort, smtoUser, smtpPassword string) *MailService {
	return &MailService{
		SMTPServer:   smtpServer,
		SMTPPort:     smtpPort,
		SMTPUser:     smtoUser,
		SMTPPassword: smtpPassword,
	}
}

func (m *MailService) IsValidMimeType(mimeType string) bool {
	validMimeTypes := []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/pdf",
	}
	for _, valid := range validMimeTypes {
		if mimeType == valid {
			return true
		}
	}
	return false
}

func (m *MailService) SendFile(emailRequest models.EmailRequest) error {
	// Проверяем MIME тип
	if !m.IsValidMimeType("application/vnd.openxmlformats-officedocument.wordprocessingml.document") &&
		!m.IsValidMimeType("application/pdf") {
		return fmt.Errorf("неподдерживаемый формат файла")
	}

	// Подготовка письма
	subject := "File Attached"
	body := "Please find the attached file."

	// Создание сообщения
	msg := &mail.Message{
		Header: map[string][]string{
			"From":    {m.SMTPUser},
			"To":      emailRequest.EmailAddresses,
			"Subject": {subject},
		},
		Body: bytes.NewBufferString(body),
	}

	// Открываем файл
	file, err := os.Open(emailRequest.FilePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	// Создаем multipart
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Создаём форму для файла
	part, err := writer.CreateFormFile("attachment", emailRequest.Filename)
	if err != nil {
		return fmt.Errorf("не удалось создать файл в форме: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("не удалось скопировать файл: %w", err)
	}

	// Закрываем writer
	writer.Close()

	// Устанавливаем Content-Type
	msg.Header["Content-Type"] = []string{fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary())}

	// Отправляем письмо
	auth := smtp.PlainAuth("", m.SMTPUser, m.SMTPPassword, m.SMTPServer)
	err = smtp.SendMail(m.SMTPServer+":"+m.SMTPPort, auth, m.SMTPUser, emailRequest.EmailAddresses, buffer.Bytes()) // исправлено
	if err != nil {
		return fmt.Errorf("не удалось отправить письмо: %w", err)
	}

	return nil
}
