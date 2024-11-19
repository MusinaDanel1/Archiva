package handlers

import (
	"archiva/internal/services"
	"fmt"
	"net/http"
)

func ArchiveFilesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	archiveData, err := services.ProcessFiles(files)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing files: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(archiveData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending zip file: %v", err), http.StatusInternalServerError)
	}
}
