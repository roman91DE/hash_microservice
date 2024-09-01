package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/hash", hashHandler)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Text      string `json:"text"`
		Algorithm string `json:"algorithm"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashChan := make(chan string)
	errChan := make(chan error)

	go func() {
		hash, err := computeHash(req.Text, req.Algorithm)
		if err != nil {
			errChan <- err
			return
		}
		hashChan <- hash
	}()

	select {
	case hash := <-hashChan:
		json.NewEncoder(w).Encode(map[string]string{"hash": hash})
	case err := <-errChan:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func computeHash(text, algorithm string) (string, error) {
	var hashBytes []byte

	switch strings.ToLower(algorithm) { // Updated line
	case "md5":
		hash := md5.Sum([]byte(text))
		hashBytes = hash[:]
	case "sha1":
		hash := sha1.Sum([]byte(text))
		hashBytes = hash[:]
	case "sha256":
		hash := sha256.Sum256([]byte(text))
		hashBytes = hash[:]
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}

	return hex.EncodeToString(hashBytes), nil
}
