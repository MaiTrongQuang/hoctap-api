package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"hoctap-api/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Response represents a standard API response
type Response struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// Global user repository
var userRepo *database.UserRepository

// Middleware for CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware for logging requests
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// Helper function to send JSON response
func sendJSONResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	dbStatus := "healthy"
	if database.DB == nil {
		dbStatus = "disconnected"
	} else if err := database.DB.Ping(); err != nil {
		dbStatus = "error: " + err.Error()
	}

	sendJSONResponse(w, http.StatusOK, "API is running successfully", map[string]interface{}{
		"status":    "healthy",
		"version":   "1.0.0",
		"database":  dbStatus,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Get all users
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := userRepo.GetAllUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		sendJSONResponse(w, http.StatusInternalServerError, "Failed to retrieve users", nil)
		return
	}

	sendJSONResponse(w, http.StatusOK, "Users retrieved successfully", users)
}

// Get user by ID
func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID %d: %v", userID, err)
		sendJSONResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}

	sendJSONResponse(w, http.StatusOK, "User found", user)
}

// Create new user
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	// Validation
	if userData.Name == "" || userData.Email == "" {
		sendJSONResponse(w, http.StatusBadRequest, "Name and email are required", nil)
		return
	}

	user, err := userRepo.CreateUser(userData.Name, userData.Email)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		if err.Error() == fmt.Sprintf("user with email '%s' already exists", userData.Email) {
			sendJSONResponse(w, http.StatusConflict, err.Error(), nil)
		} else {
			sendJSONResponse(w, http.StatusInternalServerError, "Failed to create user", nil)
		}
		return
	}

	sendJSONResponse(w, http.StatusCreated, "User created successfully", user)
}

// Update user
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var userData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	// Validation
	if userData.Name == "" || userData.Email == "" {
		sendJSONResponse(w, http.StatusBadRequest, "Name and email are required", nil)
		return
	}

	user, err := userRepo.UpdateUser(userID, userData.Name, userData.Email)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		if err.Error() == fmt.Sprintf("user with ID %d not found", userID) {
			sendJSONResponse(w, http.StatusNotFound, err.Error(), nil)
		} else if err.Error() == fmt.Sprintf("user with email '%s' already exists", userData.Email) {
			sendJSONResponse(w, http.StatusConflict, err.Error(), nil)
		} else {
			sendJSONResponse(w, http.StatusInternalServerError, "Failed to update user", nil)
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, "User updated successfully", user)
}

// Delete user
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	err = userRepo.DeleteUser(userID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		if err.Error() == fmt.Sprintf("user with ID %d not found", userID) {
			sendJSONResponse(w, http.StatusNotFound, err.Error(), nil)
		} else {
			sendJSONResponse(w, http.StatusInternalServerError, "Failed to delete user", nil)
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, "User deleted successfully", nil)
}

// Get users statistics
func getUsersStatsHandler(w http.ResponseWriter, r *http.Request) {
	count, err := userRepo.GetUsersCount()
	if err != nil {
		log.Printf("Error getting users count: %v", err)
		sendJSONResponse(w, http.StatusInternalServerError, "Failed to get users statistics", nil)
		return
	}

	stats := map[string]interface{}{
		"total_users": count,
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	sendJSONResponse(w, http.StatusOK, "Users statistics retrieved successfully", stats)
}

// Serve the main HTML page
func serveIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Welcome endpoint (moved to /welcome)
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	sendJSONResponse(w, http.StatusOK, "Welcome to HocTap API!", map[string]interface{}{
		"endpoints": map[string]string{
			"health":      "GET /health",
			"users":       "GET /api/users",
			"user_by_id":  "GET /api/users/{id}",
			"create_user": "POST /api/users",
			"update_user": "PUT /api/users/{id}",
			"delete_user": "DELETE /api/users/{id}",
			"users_stats": "GET /api/users/stats",
			"dashboard":   "GET / (HTML Dashboard)",
		},
		"database":      "MySQL with environment configuration",
		"documentation": "Use the endpoints above to interact with the API, or visit / for the web dashboard",
	})
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	// Load environment variables
	if err := godotenv.Load("config.env"); err != nil {
		log.Printf("Warning: Could not load config.env file: %v", err)
		log.Println("Using system environment variables or defaults")
	}

	// Initialize database
	log.Println("🔧 Initializing database connection...")
	if err := database.InitDB(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize user repository
	userRepo = database.NewUserRepository()

	// Seed initial users
	log.Println("🌱 Seeding initial users...")
	if err := userRepo.SeedUsers(); err != nil {
		log.Printf("⚠️ Warning: Failed to seed users: %v", err)
	} else {
		log.Println("✅ Initial users seeded successfully")
	}

	// Create a new router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(enableCORS)
	router.Use(logRequest)

	// Serve static files (CSS, JS)
	router.HandleFunc("/static/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "styles.css")
	}).Methods("GET")

	router.HandleFunc("/static/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "script.js")
	}).Methods("GET")

	// Serve the main HTML page at root
	router.HandleFunc("/", serveIndexHandler).Methods("GET")

	// API routes
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/welcome", welcomeHandler).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/users", getUsersHandler).Methods("GET")
	api.HandleFunc("/users/stats", getUsersStatsHandler).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", getUserByIDHandler).Methods("GET")
	api.HandleFunc("/users", createUserHandler).Methods("POST")
	api.HandleFunc("/users/{id:[0-9]+}", updateUserHandler).Methods("PUT")
	api.HandleFunc("/users/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")

	// Server configuration
	port := getEnv("SERVER_PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("🛑 Shutting down server...")
		database.CloseDB()
		os.Exit(0)
	}()

	fmt.Printf("🚀 HocTap API Server starting on port %s\n", port)
	fmt.Printf("📍 Available endpoints:\n")
	fmt.Printf("   • http://localhost:%s/ (HTML Dashboard)\n", port)
	fmt.Printf("   • http://localhost:%s/health (Health check)\n", port)
	fmt.Printf("   • http://localhost:%s/welcome (API welcome)\n", port)
	fmt.Printf("   • http://localhost:%s/api/users (Users API)\n", port)
	fmt.Printf("   • http://localhost:%s/api/users/stats (Users statistics)\n", port)
	fmt.Printf("   • http://localhost:%s/static/* (Static files)\n", port)
	fmt.Printf("\n💾 Database: MySQL with environment configuration\n")
	fmt.Printf("💡 Press Ctrl+C to stop the server\n")
	fmt.Printf("🌐 Open http://localhost:%s in your browser to use the dashboard\n\n", port)

	// Start server
	log.Fatal(server.ListenAndServe())
}
