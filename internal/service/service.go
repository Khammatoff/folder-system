package service

// Service holds all the service interfaces.
type Service struct {
	Auth     AuthService
	Document DocumentService
	Folder   FolderService
}
