// polarhive.net/pasta
package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var sitename = "https://x.polarhive.net"
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

func serveRobotsTxt(w http.ResponseWriter, r *http.Request) {
	robotsTxt := `User-agent: *
Disallow: /`
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(robotsTxt))
}

// usage
func handlePaste(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf(`$ curl -d "@file.txt" "%s"`, sitename), http.StatusMethodNotAllowed) // curl only
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
                http.Error(w, "File size > 1 MB", http.StatusBadRequest)
                return
	}

	if len(body) == 0 {
                http.Error(w, "Empty file", http.StatusBadRequest)
                return
	}

	// random ID
	id := generateRandomID()
	err = os.MkdirAll("data", 0755) // Save data in the "data" folder
	if err != nil {
		http.Error(w, "Error creating folder", http.StatusInternalServerError)
		return
	}

	// check if empty
	if !isEmptyFile(body) {
		// write to disk
		file, err := os.Create("data/" + id) // Save data in the "data" folder
		if err != nil {
			http.Error(w, "Error saving data", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		_, err = file.Write(body)
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

	// get ID
	id := r.URL.Path[len("/data/"):]

	// get paste from the corresponding ID
	file, err := os.Open("data/" + id)
	if err != nil {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// send paste to the client
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

