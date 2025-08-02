package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"src/utils"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CheckEmailRequest struct {
	Email string `json:"email"`
}

type CheckEmailResponse struct {
	Email        string `json:"email"`
	Compromised  bool   `json:"compromised"`
	Message      string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	count, err := app.EmailService.GetCompromisedEmailCount()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	response := HealthResponse{
		Status:  "healthy",
		Message: "Breach checker API is running with " + strconv.Itoa(count) + " compromised emails in database",
	}

	json.NewEncoder(w).Encode(response)
}

func (app *App) checkEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CheckEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON format"})
		return
	}

	if req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Email is required"})
		return
	}

	// Validate and hash the email
	emailHash, valid := utils.ValidateAndHashEmail(req.Email)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid email format"})
		return
	}

	// Check cache first
	if compromised, found := app.Cache.Get(emailHash); found {
		response := CheckEmailResponse{
			Email:       req.Email,
			Compromised: compromised,
			Message:     getResponseMessage(compromised, true),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Cache miss - check database
	compromised, err := app.EmailService.IsEmailCompromised(emailHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Database error occurred"})
		return
	}

	// Update cache with result
	app.Cache.Set(emailHash, compromised)

	response := CheckEmailResponse{
		Email:       req.Email,
		Compromised: compromised,
		Message:     getResponseMessage(compromised, false),
	}

	json.NewEncoder(w).Encode(response)
}

func getResponseMessage(compromised, fromCache bool) string {
	source := ""
	if fromCache {
		source = " (cached)"
	}

	if compromised {
		return "This email address has been found in known data breaches" + source
	}
	return "This email address was not found in known data breaches" + source
}