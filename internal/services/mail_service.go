package services

import (
	"archiva/internal/models"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService struct {
	APIKey string
}

func NewMailService(apiKey string) *MailService {
	return &MailService{
		APIKey: apiKey,
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

	file, err := os.Open(emailRequest.FilePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	fileBytes := make([]byte, 512)
	if _, err := file.Read(fileBytes); err != nil {
		return fmt.Errorf("не удалось прочитать файл: %w", err)
	}
	mimeType := http.DetectContentType(fileBytes)
	if !m.IsValidMimeType(mimeType) {
		return fmt.Errorf("неподдерживаемый формат файла: %s", mimeType)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("ошибка при сбросе указателя файла: %w", err)
	}

	from := mail.NewEmail("user1", "musinadanel1@gmail.com")
	subject := "File Attached"
	to := mail.NewEmail("Recipient", emailRequest.EmailAddresses[0])
	plainTextContent := "Please find the attached file."
	htmlContent := "<strong>Please find the attached file.</strong>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл в память: %w", err)
	}

	encodedContent := base64.StdEncoding.EncodeToString(fileContent)

	attachment := mail.NewAttachment()
	attachment.SetContent(encodedContent)
	attachment.SetType(mimeType)
	attachment.SetFilename(emailRequest.Filename)
	attachment.SetDisposition("attachment")
	message.AddAttachment(attachment)

	client := sendgrid.NewSendClient(m.APIKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("ошибка при отправке email через SendGrid: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("ошибка при отправке email через SendGrid: %s", response.Body)
	}

	return nil
}
