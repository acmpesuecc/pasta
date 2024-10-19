package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	sitename    = "http://localhost:8080" // skip port number if pushing to prod
	passphrase  = os.Getenv("PASSPHRASE")
	cooldown    = os.Getenv("COOLDOWN")
	rateLimiter = NewRateLimiter(cooldown, time.Minute)
)

const maxFileSize int64 = 1 * 1024 * 1024 in bytes

type RateLimiter struct {
	visitors map[string]int
	mu       sync.Mutex
	limit    int
	reset    time.Duration
}

func NewRateLimiter(limit int, reset time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]int),
		limit:    limit,
		reset:    reset,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.visitors[ip] >= rl.limit {
		return false
	}
	rl.visitors[ip]++
	time.AfterFunc(rl.reset, func() {
		rl.mu.Lock()
		rl.visitors[ip] = 0
		rl.mu.Unlock()
	})
	return true
}

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

	// Check for passphrase in headers
	if r.Header.Get("X-Auth-Passphrase") != passphrase {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Rate limiting check
	ip := getIP(r)
	if !rateLimiter.Allow(ip) {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
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

	body, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file data", http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Empty file", http.StatusBadRequest)
		return
	}

	id := generateRandomID()
	err = os.MkdirAll("data", 0755)
	if err != nil {
		http.Error(w, "Error creating folder", http.StatusInternalServerError)
		return
	}

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
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path[len("/data/"):]
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

func getIP(r *http.Request) string {
	addr := r.RemoteAddr

	// Check if the address is an IPv6 address (contains brackets)
	if strings.Contains(addr, "[") && strings.Contains(addr, "]") {
		// Extract the part between the brackets
		start := strings.Index(addr, "[") + 1
		end := strings.Index(addr, "]")
		return addr[start:end]
	}

	// For IPv4 addresses, split on ':' and take the first part
	ip := strings.Split(addr, ":")[0]
	return ip
}
