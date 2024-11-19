package services

import (
	"archiva/internal/models"
	"archive/zip"
	"errors"
	"mime/multipart"
)

func ProcessArchive(file multipart.File, filename string) (*models.Archive, error) {
	// Определяем размер файла
	fileSize, err := file.Seek(0, 2) // Переходим в конец файла
	if err != nil {
		return nil, errors.New("unable to determine file size")
	}

	_, err = file.Seek(0, 0) // Возвращаемся в начало файла
	if err != nil {
		return nil, errors.New("unable to reset file pointer")
	}
	zipReader, err := zip.NewReader(file, fileSize)
	if err != nil {
		return nil, errors.New("file is not a valid ZIP archive")
	}

	archive := &models.Archive{
		Filename: filename,
	}

	var totalSize float64
	for _, f := range zipReader.File {
		totalSize += float64(f.UncompressedSize64)
		archive.Files = append(archive.Files, models.File{
			FilePath: f.Name,
			Size:     f.UncompressedSize64,
		})
	}

	archive.TotalSize = totalSize
	archive.TotalFiles = float64(len(archive.Files))

	return archive, nil
}
