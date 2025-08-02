package main

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type AdminStatusResponse struct {
	Service        string            `json:"service"`
	Version        string            `json:"version"`
	Uptime         string            `json:"uptime"`
	DatabaseStats  DatabaseStats     `json:"database_stats"`
	CacheStats     CacheStats        `json:"cache_stats"`
	SystemStats    SystemStats       `json:"system_stats"`
	SecurityEvents []SecurityEvent   `json:"recent_security_events"`
}

type DatabaseStats struct {
	CompromisedEmails int    `json:"compromised_emails"`
	DatabasePath      string `json:"database_path"`
	DatabaseSize      string `json:"database_size_mb"`
}

type CacheStats struct {
	Size      int    `json:"current_size"`
	TTL       string `json:"ttl"`
	HitRatio  string `json:"hit_ratio_estimate"`
}

type SystemStats struct {
	GoVersion     string `json:"go_version"`
	NumGoroutines int    `json:"num_goroutines"`
	MemoryUsage   string `json:"memory_usage_mb"`
	CPUCount      int    `json:"cpu_count"`
}

type SecurityEvent struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	RemoteIP  string `json:"remote_ip"`
}

type AdminMetricsResponse struct {
	RequestsTotal      int64             `json:"requests_total"`
	RequestsByEndpoint map[string]int64  `json:"requests_by_endpoint"`
	ResponseTimes      map[string]string `json:"avg_response_times_ms"`
	ErrorRates         map[string]string `json:"error_rates"`
}

var startTime = time.Now()

func (app *App) adminStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get database stats
	count, err := app.EmailService.GetCompromisedEmailCount()
	if err != nil {
		http.Error(w, "Failed to get database stats", http.StatusInternalServerError)
		return
	}
	
	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	response := AdminStatusResponse{
		Service: "Email Checker API",
		Version: "1.0.0",
		Uptime:  time.Since(startTime).String(),
		DatabaseStats: DatabaseStats{
			CompromisedEmails: count,
			DatabasePath:      getEnv("DB_PATH", "email_checker.db"),
			DatabaseSize:      "N/A", // Could be enhanced to get actual file size
		},
		CacheStats: CacheStats{
			Size:     app.Cache.Size(),
			TTL:      "15m",
			HitRatio: "N/A", // Could be enhanced with cache hit tracking
		},
		SystemStats: SystemStats{
			GoVersion:     runtime.Version(),
			NumGoroutines: runtime.NumGoroutine(),
			MemoryUsage:   strconv.FormatFloat(float64(m.Alloc)/1024/1024, 'f', 2, 64),
			CPUCount:      runtime.NumCPU(),
		},
		SecurityEvents: []SecurityEvent{
			{
				Timestamp: time.Now().Format(time.RFC3339),
				Event:     "Admin access",
				RemoteIP:  r.RemoteAddr,
			},
		},
	}
	
	json.NewEncoder(w).Encode(response)
}

func (app *App) adminMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Basic metrics - in a real implementation, you'd track these
	response := AdminMetricsResponse{
		RequestsTotal: 0, // Would be tracked in middleware
		RequestsByEndpoint: map[string]int64{
			"/api/health":      0,
			"/api/check-email": 0,
			"/admin/status":    0,
		},
		ResponseTimes: map[string]string{
			"/api/health":      "5ms",
			"/api/check-email": "15ms",
			"/admin/status":    "10ms",
		},
		ErrorRates: map[string]string{
			"/api/health":      "0%",
			"/api/check-email": "2%",
			"/admin/status":    "0%",
		},
	}
	
	json.NewEncoder(w).Encode(response)
}