package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ArchiveFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Разбираем multipart форму
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10MB max
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	// Получаем файлы
	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	// Создаём новый буфер для архива
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Обрабатываем файлы и добавляем их в архив
	for _, fileHeader := range files {
		// Открываем файл
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to open file %s: %v", fileHeader.Filename, err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Создаем новый файл в архиве
		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create zip entry: %v", err), http.StatusInternalServerError)
			return
		}

		// Копируем содержимое файла в архив
		_, err = io.Copy(zipFile, file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to write file to zip: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Закрываем архив
	err = zipWriter.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to close zip archive: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправляем архив в ответе
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes()) // Отправляем архив как бинарные данные
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending zip file: %v", err), http.StatusInternalServerError)
	}
}
