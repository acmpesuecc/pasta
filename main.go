package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var sitename = "http://localhost:8080"    // skip port number if pushing to prod
const maxFileSize int64 = 1 * 1024 * 1024 // 1 MB in bytes

func main() {
	http.HandleFunc("/", handlePaste)
	http.HandleFunc("/data/", viewDataHandler)
	http.HandleFunc("/robots.txt", serveRobotsTxt) // no robots

	port := "8080"
	fmt.Printf("Starting server on port: %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

// no robots
func serveRobotsTxt(w http.ResponseWriter, r *http.Request) {
	robotsTxt := `User-agent: *
Disallow: /`
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(robotsTxt))
}

// handle file upload
func handlePaste(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("usage:\ncurl -F \"file=@file.txt\" \"%s\"", sitename), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, "File size > 1 MB", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// read file data
	body, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file data", http.StatusInternalServerError)
		return
	}

	// no empty files
	if len(body) == 0 {
		http.Error(w, "Empty file", http.StatusBadRequest)
		return
	}

	// compute hash of file contents
	hash := sha256.Sum256(body)
	hashStr := hex.EncodeToString(hash[:])

	err = os.MkdirAll("data", 0755)
	if err != nil {
		http.Error(w, "Error creating folder", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join("data", hashStr)

	// check if file with the same hash already exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create new file if it doesn't exist
		newFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error saving data", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		_, err = newFile.Write(body)
		if err != nil {
			http.Error(w, "Error saving data", http.StatusInternalServerError)
			return
		}
	}

	// return URL to the file (reuse if already exists)
	fmt.Fprintf(w, "%s/data/%s\n", sitename, hashStr)
}

func viewDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // curl only
		return
	}
	id := r.URL.Path[len("/data/"):]

	// open file by hash
	filePath := filepath.Join("data", id)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error accessing file", http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	// send file contents to client
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving file", http.StatusInternalServerError)
		return
	}
}

// utility to generate random ID (can be kept if needed)
func generateRandomID() string {
	id := make([]byte, 4) // id length
	rand.Read(id)
	return hex.EncodeToString(id)
}
