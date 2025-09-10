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
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.SheetsCount <= 0 {
		http.Error(w, `{"error": "Title and positive sheets_count are required"}`, http.StatusBadRequest)
		return
	}

	document, err := h.documentService.CreateDocument(req.Title, req.SheetsCount, req.FolderID, req.DocumentTypeID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(document)
}

func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error": "Invalid document ID"}`, http.StatusBadRequest)
		return
	}

	document, err := h.documentService.GetDocument(uint(id))
	if err != nil {
		http.Error(w, `{"error": "Document not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *DocumentHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error": "Invalid document ID"}`, http.StatusBadRequest)
		return
	}

	var req UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	document, err := h.documentService.UpdateDocument(uint(id), req.Title, req.SheetsCount, req.FolderID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error": "Invalid document ID"}`, http.StatusBadRequest)
		return
	}

	err = h.documentService.DeleteDocument(uint(id))
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
