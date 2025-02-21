package handlers

import (
	"archive/zip"
	"bytes"
	"download-api/internal/config"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func createTempDirWithFiles(t *testing.T, files map[string]string) string {
	tempDir, err := os.MkdirTemp("", "testfolder")
	if err != nil {
		t.Fatal(err)
	}

	for name, content := range files {
		filePath := filepath.Join(tempDir, name)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return tempDir
}

func verifyZipContent(t *testing.T, zipData []byte, expectedFiles map[string]string) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		t.Fatal(err)
	}

	for name, expectedContent := range expectedFiles {
		found := false
		for _, file := range zipReader.File {
			if file.Name == name {
				found = true
				rc, err := file.Open()
				if err != nil {
					t.Fatal(err)
				}
				defer rc.Close()

				var buf bytes.Buffer
				_, err = io.Copy(&buf, rc)
				if err != nil {
					t.Fatal(err)
				}

				if buf.String() != expectedContent {
					t.Errorf("file content mismatch for %s: got %v want %v", name, buf.String(), expectedContent)
				}
			}
		}
		if !found {
			t.Errorf("file %v not found in zip archive", name)
		}
	}
}

func TestDownloadZipHandler(t *testing.T) {
	files := map[string]string{
		"testfile": "This is a test file",
	}
	tempDir := createTempDirWithFiles(t, files)
	defer os.RemoveAll(tempDir)

	os.Setenv("FOLDER_PATH", tempDir)
	os.Setenv("ZIP_FILE_NAME", "testarchive.zip")
	config.LoadConfig()

	req, err := http.NewRequest("GET", "/download", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DownloadZipHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "application/zip" {
		t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/zip")
	}

	if rr.Header().Get("Content-Disposition") != `attachment; filename=testarchive.zip` {
		t.Errorf("handler returned wrong content disposition: got %v want %v", rr.Header().Get("Content-Disposition"), `attachment; filename=testarchive.zip`)
	}

	verifyZipContent(t, rr.Body.Bytes(), files)
}

func TestDownloadZipHandlerWithSubfolder(t *testing.T) {
	files := map[string]string{
		"testfile":            "This is a test file",
		"subfolder/testfile2": "This is another test file",
	}
	tempDir := createTempDirWithFiles(t, files)
	defer os.RemoveAll(tempDir)

	os.Setenv("FOLDER_PATH", tempDir)
	os.Setenv("ZIP_FILE_NAME", "subfolderarchive.zip")
	config.LoadConfig()

	req, err := http.NewRequest("GET", "/download", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DownloadZipHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "application/zip" {
		t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/zip")
	}

	if rr.Header().Get("Content-Disposition") != `attachment; filename=subfolderarchive.zip` {
		t.Errorf("handler returned wrong content disposition: got %v want %v", rr.Header().Get("Content-Disposition"), `attachment; filename=subfolderarchive.zip`)
	}

	verifyZipContent(t, rr.Body.Bytes(), files)
}
