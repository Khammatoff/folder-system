package handler

import (
	"folder-system/internal/service"
)

type Handler struct {
	auth     *AuthHandler
	document *DocumentHandler
	folder   *FolderHandler
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		auth:     NewAuthHandler(services.Auth),
		document: NewDocumentHandler(services.Document),
		folder:   NewFolderHandler(services.Folder),
	}
}

func (h *Handler) AuthHandler() *AuthHandler {
	return h.auth
}

func (h *Handler) DocumentHandler() *DocumentHandler {
	return h.document
}

func (h *Handler) FolderHandler() *FolderHandler {
	return h.folder
}
