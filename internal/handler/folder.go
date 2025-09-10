package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"folder-system/internal/service"
)

type FolderHandler struct {
	folderService service.FolderService
}

func NewFolderHandler(folderService service.FolderService) *FolderHandler {
	return &FolderHandler{folderService: folderService}
}

func (h *FolderHandler) GetRecommendedFolder(w http.ResponseWriter, r *http.Request) {
	docTypeIDStr := r.URL.Query().Get("document_type_id")
	sheetsCountStr := r.URL.Query().Get("sheets_count")

	if docTypeIDStr == "" {
		http.Error(w, `{"error": "document_type_id query parameter is required"}`, http.StatusBadRequest)
		return
	}

	docTypeID, err := strconv.ParseUint(docTypeIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error": "Invalid document_type_id"}`, http.StatusBadRequest)
		return
	}

	sheetsCount := 0 // Default if not provided
	if sheetsCountStr != "" {
		sc, err := strconv.ParseInt(sheetsCountStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "Invalid sheets_count"}`, http.StatusBadRequest)
			return
		}
		sheetsCount = int(sc)
	}

	folder, err := h.folderService.GetRecommendedFolder(uint(docTypeID), sheetsCount)
	if err != nil {
		// "Not found" is an acceptable outcome, return null
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}
