package handlers

import (
	"archiva/internal/frameworks"
	"archiva/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
)

func ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	info, err := services.ProcessArchive(file, header.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing archive: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func RegisterRoutes(r *frameworks.Router) {
	r.Handle("POST", "/api/archive", ArchiveHandler)
}
