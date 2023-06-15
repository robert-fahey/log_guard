package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestWebhookHandlerMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("GET", "/webhook", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webhookHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestWebhookHandlerBadRequest(t *testing.T) {
	var jsonStr = []byte(`{"wrong":"json"}`)
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webhookHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestWebhookHandler(t *testing.T) {
	var jsonStr = []byte(`{"message":"test"}`)
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webhookHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Open log.txt and check if the message is written
	fileContent, err := os.ReadFile("log.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(fileContent), "test") {
		t.Errorf("log.txt does not contain the expected message: got %v want %v",
			string(fileContent), "test")
	}
}
