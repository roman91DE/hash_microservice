package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestComputeHash(t *testing.T) {
	tests := []struct {
		text      string
		algorithm string
		expected  string
		wantErr   bool
	}{
		{"hello", "md5", "5d41402abc4b2a76b9719d911017c592", false},
		{"hello", "sha1", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d", false},
		{"hello", "sha256", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", false},
		{"hello", "unknown", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.algorithm, func(t *testing.T) {
			got, err := computeHash(tt.text, tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("computeHash() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHashHandler(t *testing.T) {
	tests := []struct {
		body       map[string]string
		statusCode int
		response   map[string]string
	}{
		{map[string]string{"text": "hello", "algorithm": "md5"}, http.StatusOK, map[string]string{"hash": "5d41402abc4b2a76b9719d911017c592"}},
		{map[string]string{"text": "hello", "algorithm": "sha1"}, http.StatusOK, map[string]string{"hash": "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"}},
		{map[string]string{"text": "hello", "algorithm": "sha256"}, http.StatusOK, map[string]string{"hash": "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"}},
		{map[string]string{"text": "hello", "algorithm": "unknown"}, http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.body["algorithm"], func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, err := http.NewRequest(http.MethodPost, "/hash", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(hashHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.statusCode)
			}

			if tt.statusCode == http.StatusOK {
				var response map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}
				if response["hash"] != tt.response["hash"] {
					t.Errorf("handler returned unexpected body: got %v want %v", response["hash"], tt.response["hash"])
				}
			}
		})
	}
}
