package handlers

import (
	"archive/zip"
	"download-api/internal/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadZipHandler handles the /download endpoint and creates a zip archive of the specified folder.
func DownloadZipHandler(w http.ResponseWriter, r *http.Request) {
	// Get the folder path from the configuration
	folderPath := config.K.String("FOLDER_PATH")
	if folderPath == "" {
		http.Error(w, "FOLDER_PATH environment variable not set", http.StatusInternalServerError)
		return
	}

	// Check if the folder path exists and is a directory
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) || !info.IsDir() {
		http.Error(w, "Invalid folder path", http.StatusBadRequest)
		return
	}

	// Get the zip file name from the configuration
	zipFileName := config.K.String("ZIP_FILE_NAME")
	if zipFileName == "" {
		zipFileName = "archive.zip"
	}

	// Set the response headers for the zip file
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+zipFileName)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Create a new zip writer
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// Walk through the folder and add files to the zip archive
	err = filepath.Walk(folderPath, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Get the relative path of the file
		relPath, err := filepath.Rel(folderPath, file)
		if err != nil {
			return err
		}

		// Create a new file in the zip archive
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Open the source file
		sourceFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		// Copy the contents of the source file to the zip file
		_, err = io.Copy(zipFile, sourceFile)
		return err
	})

	// Handle any errors that occurred during the creation of the zip file
	if err != nil {
		http.Error(w, "Failed to create zip file", http.StatusInternalServerError)
		return
	}
}
