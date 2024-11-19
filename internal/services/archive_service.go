package services

import (
	"archiva/internal/models"
	"archive/zip"
	"errors"
	"mime"
	"mime/multipart"
	"path/filepath"
)

func ProcessArchive(file multipart.File, filename string) (*models.Archive, error) {
	fileSize, err := file.Seek(0, 2)
	if err != nil {
		return nil, errors.New("unable to determine file size")
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, errors.New("unable to reset file pointer")
	}
	zipReader, err := zip.NewReader(file, fileSize)
	if err != nil {
		return nil, errors.New("file is not a valid ZIP archive")
	}

	archive := &models.Archive{
		Filename:    filename,
		ArchiveSize: float64(fileSize) / (1024 * 1024), // Размер архива в MB
	}

	var totalSize float64
	for _, f := range zipReader.File {
		totalSize += float64(f.UncompressedSize64)

		// Извлекаем MIME-тип из расширения файла
		ext := filepath.Ext(f.Name)           // Получаем расширение файла
		mimeType := mime.TypeByExtension(ext) // Определяем MIME-тип по расширению

		// Если MIME-тип не определен, можно задать значение по умолчанию
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		// Добавление файла в структуру архива с MIME-типом
		archive.Files = append(archive.Files, models.File{
			FilePath: f.Name,
			Size:     float64(f.UncompressedSize64) / (1024 * 1024), // Размер файла в MB
			Mimetype: mimeType,
		})
	}

	archive.TotalSize = totalSize / (1024 * 1024) // Общий размер всех файлов в MB
	archive.TotalFiles = float64(len(archive.Files))

	return archive, nil
}
