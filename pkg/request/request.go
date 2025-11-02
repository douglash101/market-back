package request

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type METHOD string

const (
	GET    METHOD = "GET"
	POST   METHOD = "POST"
	DELETE METHOD = "DELETE"
	PUT    METHOD = "PUT"
)

type request[T any] struct {
	Name    string            `json:"name"`
	Method  METHOD            `json:"method"`
	URL     string            `json:"path"`
	Retries int               `json:"retries"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    any               `json:"body,omitempty"`
}

type RequestParams struct {
	Name    string            `json:"name"`
	Method  METHOD            `json:"method"`
	URL     string            `json:"path"`
	Retries int               `json:"retries"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    any               `json:"body,omitempty"`
}

func NewRequest[T any](params RequestParams) *request[T] {
	return &request[T]{
		Name:    params.Name,
		Method:  params.Method,
		URL:     params.URL,
		Retries: params.Retries,
		Headers: params.Headers,
		Body:    params.Body,
	}
}

func (r *request[T]) Execute() (*T, error) {
	return r.ExecuteWithContext(context.Background())
}

func (r *request[T]) ExecuteWithContext(ctx context.Context) (*T, error) {
	var result T

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	var bodyReader io.Reader
	if r.Body != nil {
		switch body := r.Body.(type) {
		case string:
			bodyReader = strings.NewReader(body)
		default:
			bodyData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("error marshaling body: %v", err)
			}
			bodyReader = strings.NewReader(string(bodyData))
		}
	}

	req, err := http.NewRequestWithContext(ctx, string(r.Method), r.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Set custom headers
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	var resp *http.Response
	for attempt := 0; attempt <= int(r.Retries); attempt++ {
		// Check if context was cancelled before making the request
		if ctx.Err() != nil {
			return nil, fmt.Errorf("request cancelled: %v", ctx.Err())
		}

		resp, err = client.Do(req)

		if err == nil && resp.StatusCode != http.StatusTooManyRequests {
			defer resp.Body.Close()
			break
		}

		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			if attempt < int(r.Retries) {
				// Use context-aware sleep to allow cancellation during retry wait
				select {
				case <-time.After(time.Second * 2):
					// Continue after sleep
				case <-ctx.Done():
					return nil, fmt.Errorf("request cancelled during retry wait: %v", ctx.Err())
				}
				continue
			}
			defer resp.Body.Close()
			return nil, fmt.Errorf("too many retries, giving up: %s", resp.Status)
		}

		// If we have an error and it's not the last attempt, wait before retrying
		if err != nil && attempt < int(r.Retries) {
			select {
			case <-time.After(time.Second * 2):
				// Continue after sleep
			case <-ctx.Done():
				return nil, fmt.Errorf("request cancelled during retry wait: %v", ctx.Err())
			}
		}
	}

	// Final check for errors after all retries
	if err != nil {
		return nil, fmt.Errorf("request failed after %d retries: %v", r.Retries, err)
	}

	// Check if response is nil (shouldn't happen, but defensive programming)
	if resp == nil {
		return nil, fmt.Errorf("no response received")
	}

	// Check if context was cancelled before reading response
	if ctx.Err() != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("request cancelled before reading response: %v", ctx.Err())
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(responseData, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}
