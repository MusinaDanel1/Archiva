package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

func ProcessFiles(files []*multipart.FileHeader) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("unable to open file %s: %v", fileHeader.Filename, err)
		}
		defer file.Close()

		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			return nil, fmt.Errorf("unable to create zip file %s: %v", fileHeader.Filename, err)
		}

		_, err = io.Copy(zipFile, file)
		if err != nil {
			return nil, fmt.Errorf("unable to copy file %s to zip: %v", fileHeader.Filename, err)
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close zip archive: %v", err)
	}

	return buf.Bytes(), nil
}
