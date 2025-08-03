package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"src/cache"
	"src/database"

	"github.com/gorilla/mux"
)

type App struct {
	Router       *mux.Router
	EmailService *database.EmailService
	Cache        *cache.Cache
}

func main() {
	app := &App{}
	app.Initialize()

	port := getEnv("PORT", "8082")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}

func (app *App) Initialize() {
	dbPath := getEnv("DB_PATH", "/data/email_checker.db")

	db, err := database.InitDatabase(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	app.EmailService = database.NewEmailService(db)
	app.Cache = cache.NewCache(15 * time.Minute) // Cache TTL of 15 minutes

	// Seed database with sample data for testing
	if err := database.SeedDatabase(app.EmailService); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	app.setupRoutes()
}

func (app *App) setupRoutes() {
	app.Router = mux.NewRouter()

	// Public API routes
	api := app.Router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", app.healthHandler).Methods("GET")
	api.HandleFunc("/check-email", app.checkEmailHandler).Methods("POST")

	// Protected admin routes (require authentication)
	admin := app.Router.PathPrefix("/admin").Subrouter()
	admin.Use(app.basicAuthMiddleware)
	admin.HandleFunc("/status", app.adminStatusHandler).Methods("GET")
	admin.HandleFunc("/metrics", app.adminMetricsHandler).Methods("GET")

	// Enable CORS for all routes
	app.Router.Use(corsMiddleware)

	// Add security middleware
	app.Router.Use(app.securityMiddleware)

	// Serve static files securely
	app.Router.PathPrefix("/").Handler(app.secureFileHandler())
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow specific domains including Cloudflare tunnel
		allowedOrigins := []string{
			"http://localhost",
			"http://localhost:80",
			"https://*.yoongjiahui.com",
			"https://razerassignmentapp.yoongjiahui.com",
			"https://razerassignment.yoongjiahui.com",
		}

		origin := r.Header.Get("Origin")
		host := r.Header.Get("Host")

		// Allow requests from tunnel domains - specific for linking to my Cloudflare
		if origin != "" {
			for _, allowed := range allowedOrigins {
				if origin == allowed ||
					(strings.Contains(allowed, "*.yoongjiahui.com") && strings.HasSuffix(origin, ".yoongjiahui.com")) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		} else if host != "" && strings.HasSuffix(host, ".yoongjiahui.com") {
			w.Header().Set("Access-Control-Allow-Origin", "https://"+host)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Real-IP, X-Forwarded-For, CF-Connecting-IP, CF-Ray, CF-Visitor")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// secureFileHandler creates a secure file handler that prevents directory traversal
func (app *App) secureFileHandler() http.Handler {
	frontendDir := "../frontend/"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean the path to prevent directory traversal
		cleanPath := filepath.Clean(r.URL.Path)

		// Remove leading slash for filepath.Join
		cleanPath = strings.TrimPrefix(cleanPath, "/")

		// If empty path, serve index.html
		if cleanPath == "" || cleanPath == "." {
			cleanPath = "index.html"
		}

		// Join with frontend directory
		fullPath := filepath.Join(frontendDir, cleanPath)

		// Ensure the final path is within the frontend directory
		absBasePath, err := filepath.Abs(frontendDir)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		absFullPath, err := filepath.Abs(fullPath)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Check if the resolved path is within the frontend directory
		if !strings.HasPrefix(absFullPath, absBasePath) {
			log.Printf("Security: Directory traversal attempt blocked: %s", r.URL.Path)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Check if file exists
		info, err := os.Stat(absFullPath)
		if err != nil {
			// If file doesn't exist, serve index.html for SPA routing
			indexPath := filepath.Join(frontendDir, "index.html")
			http.ServeFile(w, r, indexPath)
			return
		}

		// Don't serve directories
		if info.IsDir() {
			log.Printf("Security: Directory access attempt blocked: %s", r.URL.Path)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Serve the file
		http.ServeFile(w, r, absFullPath)
	})
}

// basicAuthMiddleware provides basic authentication for admin endpoints
func (app *App) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := getEnv("ADMIN_USERNAME", "")
		password := getEnv("ADMIN_PASSWORD", "")

		// Require both username and password to be set
		if username == "" || password == "" {
			log.Fatal("ADMIN_USERNAME and ADMIN_PASSWORD must be set in environment variables")
		}

		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Use constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			log.Printf("Security: Failed authentication attempt from %s", r.RemoteAddr)
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// securityMiddleware adds security headers and logging
func (app *App) securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Log suspicious requests
		if strings.Contains(r.URL.Path, "..") ||
			strings.Contains(r.URL.Path, "//") ||
			strings.Contains(r.URL.Path, "/etc/") ||
			strings.Contains(r.URL.Path, "/proc/") ||
			strings.Contains(r.URL.Path, "/home/") {
			log.Printf("Security: Suspicious request from %s: %s", r.RemoteAddr, r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})
}
