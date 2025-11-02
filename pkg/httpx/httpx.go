package httpx

import (
	"encoding/json"
	"net/http"
)

type Status string

const (
	SUCCESS Status = "success"
	ERROR   Status = "error"
)

// SendSuccess sends a success JSON response
func SendSuccess(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}

// SendCreated sends a created JSON response
func SendCreated(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(data)
}

// SendBadRequest sends a bad request JSON response
func SendBadRequest(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(data)
}

// SendMethodNotAllowed sends a method not allowed JSON response
func SendMethodNotAllowed(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	return json.NewEncoder(w).Encode(data)
}

// SendInternalServerError sends an internal server error JSON response
func SendInternalServerError(w http.ResponseWriter, message string, details any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	return json.NewEncoder(w).Encode(map[string]any{
		"error":   message,
		"details": details,
	})
}

// SendNotFound sends a not found JSON response
func SendNotFound(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	return json.NewEncoder(w).Encode(data)
}

// SendUnauthorized sends an unauthorized JSON response
func SendUnauthorized(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	return json.NewEncoder(w).Encode(data)
}

// SendForbidden sends a forbidden JSON response
func SendForbidden(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	return json.NewEncoder(w).Encode(data)
}

// HTTP Method helpers - wrap handlers to only allow specific HTTP methods

// Get wraps a handler to only allow GET requests
func Get(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			SendMethodNotAllowed(w, "Only GET method allowed")
			return
		}
		handler(w, r)
	}
}

// Post wraps a handler to only allow POST requests
func Post(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			SendMethodNotAllowed(w, "Only POST method allowed")
			return
		}
		handler(w, r)
	}
}

// Put wraps a handler to only allow PUT requests
func Put(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			SendMethodNotAllowed(w, "Only PUT method allowed")
			return
		}
		handler(w, r)
	}
}

// Patch wraps a handler to only allow PATCH requests
func Patch(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			SendMethodNotAllowed(w, "Only PATCH method allowed")
			return
		}
		handler(w, r)
	}
}

// Delete wraps a handler to only allow DELETE requests
func Delete(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			SendMethodNotAllowed(w, "Only DELETE method allowed")
			return
		}
		handler(w, r)
	}
}

// PutOrPatch wraps a handler to allow both PUT and PATCH requests
func PutOrPatch(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			SendMethodNotAllowed(w, "Only PUT and PATCH methods allowed")
			return
		}
		handler(w, r)
	}
}
