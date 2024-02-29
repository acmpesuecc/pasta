package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
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

// usage
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

	// check file's data
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

	// gen random ID
	id := generateRandomID()
	err = os.MkdirAll("data", 0755)
	if err != nil {
		http.Error(w, "Error creating folder", http.StatusInternalServerError)
		return
	}

	// check if empty
	if !isEmptyFile(body) {
		newFilePath := "data/" + id
		newFile, err := os.Create(newFilePath)
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
	} else {
		http.Error(w, "Received file is empty", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s/data/%s\n", sitename, id)
}

// check if empty
func isEmptyFile(data []byte) bool {
	for _, b := range data {
		if b != 0 {
			return false
		}
	}
	return true
}

func viewDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // curl only
		return
	}
	id := r.URL.Path[len("/data/"):]
	// get paste from ID
	file, err := os.Open("data/" + id)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Paste not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error accessing paste", http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()
	// send paste to user
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving paste", http.StatusInternalServerError)
		return
	}
}

func generateRandomID() string {
	id := make([]byte, 4) // id length
	rand.Read(id)
	return hex.EncodeToString(id)
}
