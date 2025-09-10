package main

import (
	"fmt"
	"log"
	"net/http"

	"folder-system/internal/config"
	"folder-system/internal/handler"
	custommiddleware "folder-system/internal/middleware"
	"folder-system/internal/repository/postgresql"
	"folder-system/internal/service"
	"folder-system/pkg/lib"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := lib.NewLogger(cfg.Logging.Level, cfg.Logging.File)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Build DSN for PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	// Initialize repository (PostgreSQL)
	repo, err := postgresql.NewRepository(dsn)
	if err != nil {
		logger.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize services
	authService := service.NewAuthService(repo, cfg)
	documentService := service.NewDocumentService(repo, repo)
	folderService := service.NewFolderService(repo)

	services := &service.Service{
		Auth:     authService,
		Document: documentService,
		Folder:   folderService,
	}

	// Initialize handlers
	handlers := handler.NewHandler(services)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Logger)
	r.Use(custommiddleware.LoggerMiddleware(logger))

	// Public routes
	r.Post("/api/register", handlers.AuthHandler().Register)
	r.Post("/api/login", handlers.AuthHandler().Login)

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(custommiddleware.AuthMiddleware(cfg.JWT.AccessSecret))

		// Document routes
		r.Route("/documents", func(r chi.Router) {
			r.Post("/", handlers.DocumentHandler().CreateDocument)
			r.Get("/{id}", handlers.DocumentHandler().GetDocument)
			r.Put("/{id}", handlers.DocumentHandler().UpdateDocument)
			r.Delete("/{id}", handlers.DocumentHandler().DeleteDocument)
		})

		// Folder routes
		r.Route("/folders", func(r chi.Router) {
			r.Get("/recommended", handlers.FolderHandler().GetRecommendedFolder)
		})
	})

	// Serve frontend static files
	feFS := http.FileServer(http.Dir("./frontend"))
	r.Handle("/*", feFS)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})

	// Start server
	host := cfg.Server.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	serverAddr := fmt.Sprintf("%s:%s", host, port)
	logger.Infof("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
