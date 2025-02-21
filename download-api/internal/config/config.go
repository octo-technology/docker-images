package config

import (
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

var K = koanf.New(".")

var defaultConfig = map[string]interface{}{
	"HOST":          "0.0.0.0",
	"PORT":          "8080",
	"FOLDER_PATH":   "/default/path",
	"ZIP_FILE_NAME": "archive.zip",
}

func LoadConfig() {
	// Load default configuration from the defaultConfig map
	K.Load(confmap.Provider(defaultConfig, "."), nil)

	// Override with environment variables
	K.Load(env.Provider("", ".", func(s string) string {
		return s
	}), nil)
}
