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
		WriteJSONError(w, http.StatusBadRequest, "document_type_id query parameter is required")
		return
	}

	docTypeID, err := strconv.ParseUint(docTypeIDStr, 10, 32)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid document_type_id")
		return
	}

	sheetsCount := 0
	if sheetsCountStr != "" {
		sc, err := strconv.ParseInt(sheetsCountStr, 10, 32)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, "Invalid sheets_count")
			return
		}
		sheetsCount = int(sc)
	}

	folder, err := h.folderService.GetRecommendedFolder(uint(docTypeID), sheetsCount)
	if err != nil {
		// Not found is acceptable — return JSON null
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("null"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(folder); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}
