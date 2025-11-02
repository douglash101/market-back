package request

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Test response struct for testing JSON unmarshaling
type TestResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Test request body struct
type TestRequestBody struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name     string
		params   RequestParams
		expected *request[TestResponse]
	}{
		{
			name: "Basic request creation",
			params: RequestParams{
				Name:    "test-request",
				Method:  GET,
				URL:     "https://api.example.com/test",
				Retries: 3,
			},
			expected: &request[TestResponse]{
				Name:    "test-request",
				Method:  GET,
				URL:     "https://api.example.com/test",
				Retries: 3,
				Headers: nil,
				Body:    nil,
			},
		},
		{
			name: "Request with headers and body",
			params: RequestParams{
				Name:   "post-request",
				Method: POST,
				URL:    "https://api.example.com/create",
				Headers: map[string]string{
					"Authorization": "Bearer token123",
					"Custom-Header": "custom-value",
				},
				Body: TestRequestBody{Name: "test", Value: 42},
			},
			expected: &request[TestResponse]{
				Name:   "post-request",
				Method: POST,
				URL:    "https://api.example.com/create",
				Headers: map[string]string{
					"Authorization": "Bearer token123",
					"Custom-Header": "custom-value",
				},
				Body: TestRequestBody{Name: "test", Value: 42},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewRequest[TestResponse](tt.params)

			if result.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", result.Name, tt.expected.Name)
			}
			if result.Method != tt.expected.Method {
				t.Errorf("Method = %v, want %v", result.Method, tt.expected.Method)
			}
			if result.URL != tt.expected.URL {
				t.Errorf("URL = %v, want %v", result.URL, tt.expected.URL)
			}
			if result.Retries != tt.expected.Retries {
				t.Errorf("Retries = %v, want %v", result.Retries, tt.expected.Retries)
			}
		})
	}
}

func TestExecuteSuccess(t *testing.T) {
	// Create a test server that returns a successful response
	testResponse := TestResponse{
		ID:      1,
		Message: "success",
		Status:  "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Name:   "test-success",
		Method: GET,
		URL:    server.URL,
	})

	result, err := req.Execute()

	if err != nil {
		t.Fatalf("Execute() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Execute() returned nil result")
	}

	if result.ID != testResponse.ID {
		t.Errorf("Result ID = %v, want %v", result.ID, testResponse.ID)
	}
	if result.Message != testResponse.Message {
		t.Errorf("Result Message = %v, want %v", result.Message, testResponse.Message)
	}
	if result.Status != testResponse.Status {
		t.Errorf("Result Status = %v, want %v", result.Status, testResponse.Status)
	}
}

func TestExecuteWithDifferentMethods(t *testing.T) {
	methods := []METHOD{GET, POST, PUT, DELETE}

	for _, method := range methods {
		t.Run(string(method), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != string(method) {
					t.Errorf("Expected method %s, got %s", method, r.Method)
				}

				testResponse := TestResponse{ID: 1, Message: fmt.Sprintf("%s success", method), Status: "ok"}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(testResponse)
			}))
			defer server.Close()

			req := NewRequest[TestResponse](RequestParams{
				Method: method,
				URL:    server.URL,
			})

			result, err := req.Execute()
			if err != nil {
				t.Fatalf("Execute() failed: %v", err)
			}
			if result.Message != fmt.Sprintf("%s success", method) {
				t.Errorf("Unexpected response message: %v", result.Message)
			}
		})
	}
}

func TestExecuteWithHeaders(t *testing.T) {
	expectedHeaders := map[string]string{
		"Authorization": "Bearer test-token",
		"Custom-Header": "custom-value",
		"X-API-Key":     "api-key-123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify custom headers are present
		for key, expectedValue := range expectedHeaders {
			if actualValue := r.Header.Get(key); actualValue != expectedValue {
				t.Errorf("Header %s = %v, want %v", key, actualValue, expectedValue)
			}
		}

		// Verify default headers
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Content-Type = %v, want application/json", contentType)
		}
		if accept := r.Header.Get("Accept"); accept != "application/json" {
			t.Errorf("Accept = %v, want application/json", accept)
		}

		testResponse := TestResponse{ID: 1, Message: "headers test", Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method:  POST,
		URL:     server.URL,
		Headers: expectedHeaders,
	})

	_, err := req.Execute()
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}
}

func TestExecuteWithStringBody(t *testing.T) {
	expectedBody := `{"name":"test","value":42}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		if string(body) != expectedBody {
			t.Errorf("Body = %v, want %v", string(body), expectedBody)
		}

		testResponse := TestResponse{ID: 1, Message: "string body test", Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method: POST,
		URL:    server.URL,
		Body:   expectedBody,
	})

	_, err := req.Execute()
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}
}

func TestExecuteWithStructBody(t *testing.T) {
	testBody := TestRequestBody{Name: "test", Value: 42}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var receivedBody TestRequestBody
		json.NewDecoder(r.Body).Decode(&receivedBody)

		if receivedBody.Name != testBody.Name || receivedBody.Value != testBody.Value {
			t.Errorf("Body = %+v, want %+v", receivedBody, testBody)
		}

		testResponse := TestResponse{ID: 1, Message: "struct body test", Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method: POST,
		URL:    server.URL,
		Body:   testBody,
	})

	_, err := req.Execute()
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}
}

func TestExecuteWithRetries(t *testing.T) {
	callCount := 0
	maxRetries := 3

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		// Fail the first few attempts, succeed on the last
		if callCount <= maxRetries {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		testResponse := TestResponse{ID: 1, Message: "retry success", Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method:  GET,
		URL:     server.URL,
		Retries: maxRetries,
	})

	result, err := req.Execute()
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if result.Message != "retry success" {
		t.Errorf("Expected retry success, got %v", result.Message)
	}

	if callCount != maxRetries+1 {
		t.Errorf("Expected %d calls, got %d", maxRetries+1, callCount)
	}
}

func TestExecuteWithRateLimitRetry(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		// Return rate limit on first call, success on second
		if callCount == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		testResponse := TestResponse{ID: 1, Message: "rate limit retry success", Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method:  GET,
		URL:     server.URL,
		Retries: 2,
	})

	result, err := req.Execute()
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if result.Message != "rate limit retry success" {
		t.Errorf("Expected rate limit retry success, got %v", result.Message)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestExecuteWithMaxRetriesExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method:  GET,
		URL:     server.URL,
		Retries: 1,
	})

	_, err := req.Execute()
	if err == nil {
		t.Fatal("Expected error for max retries exceeded")
	}

	if !strings.Contains(err.Error(), "too many retries") {
		t.Errorf("Expected 'too many retries' error, got: %v", err)
	}
}

func TestExecuteWithContext(t *testing.T) {
	testResponse := TestResponse{ID: 1, Message: "context test", Status: "ok"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method: GET,
		URL:    server.URL,
	})

	ctx := context.Background()
	result, err := req.ExecuteWithContext(ctx)

	if err != nil {
		t.Fatalf("ExecuteWithContext() failed: %v", err)
	}

	if result.Message != testResponse.Message {
		t.Errorf("Expected %v, got %v", testResponse.Message, result.Message)
	}
}

func TestExecuteWithContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method: GET,
		URL:    server.URL,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := req.ExecuteWithContext(ctx)
	if err == nil {
		t.Fatal("Expected context cancellation error")
	}

	if !strings.Contains(err.Error(), "cancelled") && !strings.Contains(err.Error(), "deadline exceeded") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}
}

func TestExecuteWithInvalidURL(t *testing.T) {
	req := NewRequest[TestResponse](RequestParams{
		Method: GET,
		URL:    "invalid-url",
	})

	_, err := req.Execute()
	if err == nil {
		t.Fatal("Expected error for invalid URL")
	}
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method: GET,
		URL:    server.URL,
	})

	_, err := req.Execute()
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "unmarshaling") {
		t.Errorf("Expected unmarshaling error, got: %v", err)
	}
}

func TestExecuteWithInvalidBodyMarshaling(t *testing.T) {
	// Create a body that can't be marshaled to JSON (function type)
	invalidBody := func() {}

	req := NewRequest[TestResponse](RequestParams{
		Method: POST,
		URL:    "http://example.com",
		Body:   invalidBody,
	})

	_, err := req.Execute()
	if err == nil {
		t.Fatal("Expected error for invalid body marshaling")
	}

	if !strings.Contains(err.Error(), "marshaling") {
		t.Errorf("Expected marshaling error, got: %v", err)
	}
}

func TestExecuteCallsExecuteWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testResponse := TestResponse{ID: 1, Message: "execute calls context", Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse)
	}))
	defer server.Close()

	req := NewRequest[TestResponse](RequestParams{
		Method: GET,
		URL:    server.URL,
	})

	result, err := req.Execute()
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if result.Message != "execute calls context" {
		t.Errorf("Expected correct response from Execute() method")
	}
}
