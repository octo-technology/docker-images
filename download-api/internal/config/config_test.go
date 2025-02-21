package config

import (
    "os"
    "testing"
)

func TestLoadConfig(t *testing.T) {
    os.Setenv("FOLDER_PATH", "/test/folder")
    os.Setenv("HOST", "127.0.0.1")
    os.Setenv("PORT", "9090")

    LoadConfig()

    if K.String("FOLDER_PATH") != "/test/folder" {
        t.Errorf("Expected FOLDER_PATH to be '/test/folder', got '%s'", K.String("FOLDER_PATH"))
    }

    if K.String("HOST") != "127.0.0.1" {
        t.Errorf("Expected HOST to be '127.0.0.1', got '%s'", K.String("HOST"))
    }

    if K.String("PORT") != "9090" {
        t.Errorf("Expected PORT to be '9090', got '%s'", K.String("PORT"))
    }
}
