package services

import (
	"archiva/internal/models"
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

var allowedMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}

// ProcessFiles обрабатывает файлы, проверяет их типы и создаёт ZIP-архив.
func ProcessFiles(files []*multipart.FileHeader) (*models.Archive, error) {
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	var totalSize float64
	var archiveFiles []models.File

	for _, fileHeader := range files {
		// Открываем файл
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Проверяем MIME тип
		mimeType, err := validateFileMimeType(file)
		if err != nil {
			return nil, err
		}

		// Записываем файл в архив
		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			return nil, err
		}

		// Копируем содержимое файла в архив
		_, err = io.Copy(zipFile, file)
		if err != nil {
			return nil, err
		}

		// Добавляем информацию о файле в список
		archiveFiles = append(archiveFiles, models.File{
			FilePath: fileHeader.Filename,
			Size:     float64(fileHeader.Size),
			Mimetype: mimeType,
		})

		// Накапливаем общий размер архива
		totalSize += float64(fileHeader.Size)
	}

	// Закрываем архивный файл
	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	// Создаём объект архива
	archive := &models.Archive{
		Filename:    "archive.zip", // или можно передать имя в параметре
		ArchiveSize: float64(zipBuffer.Len()),
		TotalSize:   totalSize,
		TotalFiles:  float64(len(archiveFiles)),
		Files:       archiveFiles,
	}

	return archive, nil
}

// validateFileMimeType проверяет MIME тип файла
func validateFileMimeType(file multipart.File) (string, error) {
	// Read the first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", errors.New("unable to read file for MIME type detection")
	}

	// Use http.DetectContentType to get the MIME type
	mimeType := http.DetectContentType(buffer)

	// If the MIME type is not valid, return an error
	if !allowedMimeTypes[mimeType] {
		return "", errors.New("unsupported MIME type: " + mimeType)
	}

	// Reset the file pointer back to the start after reading the buffer
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", errors.New("unable to reset file pointer")
	}

	return mimeType, nil
}
