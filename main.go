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

func main() {
	http.HandleFunc("/", handlePaste)
	http.HandleFunc("/data/", viewDataHandler)

	port := "8080"
	fmt.Printf("Starting server on port: %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func handlePaste(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	// random ID
	id := generateRandomID()
	err = os.MkdirAll("data", 0755) // Save data in the "data" folder
	if err != nil {
		http.Error(w, "Error creating folder", http.StatusInternalServerError)
		return
	}

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

	fmt.Fprintf(w, "Paste received and stored with ID: %s\n", id)
}

func viewDataHandler(w http.ResponseWriter, r *http.Request) {
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
