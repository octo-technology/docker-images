package main

import (
	"download-api/internal/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func setupServer() {
	// Start the server in a goroutine
	go func() {
		main()
	}()

	// Allow some time for the server to start
	time.Sleep(1 * time.Second)
}

func TestMain(m *testing.M) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testfolder")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary file in the directory
	tempFile, err := os.CreateTemp(tempDir, "testfile")
	if err != nil {
		panic(err)
	}
	tempFile.WriteString("This is a test file")
	tempFile.Close()

	// Set the FOLDER_PATH environment variable to the temporary directory
	os.Setenv("FOLDER_PATH", tempDir)
	os.Setenv("ZIP_FILE_NAME", "testarchive.zip")

	// Load configuration
	config.LoadConfig()

	// Start the server
	setupServer()

	// Run the tests
	code := m.Run()

	os.Exit(code)
}

func TestLivezHandler(t *testing.T) {
	// Create a request to the /livez endpoint
	req, err := http.NewRequest("GET", "/livez", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request
	http.DefaultServeMux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("/livez handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	if rr.Body.String() != "OK" {
		t.Errorf("/livez handler returned unexpected body: got %v want %v", rr.Body.String(), "OK")
	}
}

func TestReadyzHandler(t *testing.T) {
	// Create a request to the /readyz endpoint
	req, err := http.NewRequest("GET", "/readyz", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request
	http.DefaultServeMux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("/readyz handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	if rr.Body.String() != "OK" {
		t.Errorf("/readyz handler returned unexpected body: got %v want %v", rr.Body.String(), "OK")
	}
}
