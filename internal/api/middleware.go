package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Time    string `json:"timestamp"`
}

// SuccessResponse wraps a successful response
type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Time      string      `json:"timestamp"`
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		handler.ServeHTTP(w, r)
		log.Printf("Request completed in %v", time.Since(start))
	})
}

// JSONResponseMiddleware sets JSON content type
func JSONResponseMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteError writes a JSON error response
func WriteError(w http.ResponseWriter, statusCode int, err string, message string) {
	response := ErrorResponse{
		Error:   err,
		Message: message,
		Code:    statusCode,
		Time:    time.Now().UTC().Format(time.RFC3339),
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// WriteSuccess writes a JSON success response
func WriteSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := SuccessResponse{
		Data:    data,
		Message: message,
		Time:    time.Now().UTC().Format(time.RFC3339),
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
