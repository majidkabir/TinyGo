package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/majidkabir/TinyGo/internal/database"
	"github.com/majidkabir/TinyGo/internal/handlers"
	"github.com/majidkabir/TinyGo/internal/repository"
	"github.com/majidkabir/TinyGo/internal/service"
	"github.com/majidkabir/TinyGo/pkg/middleware"
)

func main() {
	// Database configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "userdb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to database
	db, err := database.New(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize layers (Dependency Injection)
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Setup routes
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// User routes
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.Logger(middleware.CORS(userHandler.ListUsers))(w, r)
		case http.MethodPost:
			middleware.Logger(middleware.CORS(userHandler.CreateUser))(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.Logger(middleware.CORS(userHandler.GetUser))(w, r)
		case http.MethodPut:
			middleware.Logger(middleware.CORS(userHandler.UpdateUser))(w, r)
		case http.MethodDelete:
			middleware.Logger(middleware.CORS(userHandler.DeleteUser))(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start server
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on http://localhost:%s", port)
		log.Println("Available endpoints:")
		log.Println("  GET    /health           - Health check")
		log.Println("  GET    /api/users        - List users")
		log.Println("  POST   /api/users        - Create user")
		log.Println("  GET    /api/users/{id}   - Get user")
		log.Println("  PUT    /api/users/{id}   - Update user")
		log.Println("  DELETE /api/users/{id}   - Delete user")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	// Could add server.Shutdown(ctx) here for graceful shutdown
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
