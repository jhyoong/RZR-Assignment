package main

import (
	"log"
	"net/http"
	"os"
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
	dbPath := getEnv("DB_PATH", "email_checker.db")
	
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
	
	// API routes
	api := app.Router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", app.healthHandler).Methods("GET")
	api.HandleFunc("/check-email", app.checkEmailHandler).Methods("POST")

	// Enable CORS for all routes
	app.Router.Use(corsMiddleware)

	// Serve static files
	app.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/")))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow specific domains including Cloudflare tunnel
		allowedOrigins := []string{
			"http://localhost",
			"http://localhost:80",
			"https://*.yoongjiahui.com",
			"https://razerassignment.yoongjiahui.com",
		}
		
		origin := r.Header.Get("Origin")
		host := r.Header.Get("Host")
		
		// Allow requests from tunnel domains
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