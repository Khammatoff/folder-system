package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"folder-system/internal/service"

	"github.com/go-chi/chi/v5"
)

type DocumentHandler struct {
	documentService service.DocumentService
}

func NewDocumentHandler(documentService service.DocumentService) *DocumentHandler {
	return &DocumentHandler{documentService: documentService}
}

type CreateDocumentRequest struct {
	Title          string `json:"title"`
	SheetsCount    int    `json:"sheets_count"`
	FolderID       *uint  `json:"folder_id"`
	DocumentTypeID uint   `json:"document_type_id"`
}

type UpdateDocumentRequest struct {
	Title       *string `json:"title,omitempty"`
	SheetsCount *int    `json:"sheets_count,omitempty"`
	FolderID    *uint   `json:"folder_id,omitempty"`
}

func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var req CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.SheetsCount <= 0 {
		WriteJSONError(w, http.StatusBadRequest, "Title and positive sheets_count are required")
		return
	}

	if req.DocumentTypeID == 0 {
		req.DocumentTypeID = 1 // default document type id
	}

	document, err := h.documentService.CreateDocument(req.Title, req.SheetsCount, req.FolderID, req.DocumentTypeID)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(document)
}

func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	document, err := h.documentService.GetDocument(uint(id))
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(document); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

func (h *DocumentHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var req UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	document, err := h.documentService.UpdateDocument(uint(id), req.Title, req.SheetsCount, req.FolderID)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(document)
}

func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	err = h.documentService.DeleteDocument(uint(id))
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
