package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Notification struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// We expect a POST request with a body
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON body into our Notification struct
	var notif Notification
	if err := json.Unmarshal(body, &notif); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Check if 'message' field is empty
	if notif.Message == "" {
		http.Error(w, "Missing 'message' in request body", http.StatusBadRequest)
		return
	}

	// Append to file
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err = f.WriteString(notif.Message + "\n"); err != nil {
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	// Send a response back to the requester
	w.Write([]byte("Successfully wrote to file\n"))
}
