package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	// creates an HTTP handler from our main package
	hf := http.HandlerFunc(handler)

	// Serve the HTTP request to our recorder. This line executes the handler.
	hf.ServeHTTP(recorder, req)

	// Create error message for status code to check assertions

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `Hello, World!`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned wrong body: got %v want %v", actual, expected)
	}
}

// The recorder is like a mini-browser, it accepts the result of the request
