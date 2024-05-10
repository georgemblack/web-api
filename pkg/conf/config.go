package conf

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	GCloudProjectID       string `json:"gcloudProjectID"`
	BuildServiceEndpoint  string `json:"buildServiceEndpoint"`
	AllowedOriginHeader   string `json:"allowedOriginHeader"`
	FirestoreDatabasename string `json:"firestoreDatabaseName"`
	BackupBucketName      string `json:"backupBucketName"`
	APIUsername           string `json:"apiUsername"`
	APIPassword           string `json:"apiPassword"`
	TokenSecret           string `json:"tokenSecret"`
}

//go:embed config/*
var configFiles embed.FS

// LoadConfig attemps to load configuration first from an embedded JSON file,
// then via environment variables that can override.
func LoadConfig() (Config, error) {
	// Load config via static files
	var config Config
	bytes, err := configFiles.ReadFile("config/" + getEnv() + ".json")
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file; %w", err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config file; %w", err)
	}

	// Load config via env vars
	if config.APIUsername == "" {
		config.APIUsername = os.Getenv("API_USERNAME")
	}
	if config.APIPassword == "" {
		config.APIPassword = os.Getenv("API_PASSWORD")
	}
	if config.TokenSecret == "" {
		config.TokenSecret = os.Getenv("TOKEN_SECRET")
	}

	return config, nil
}

// Base64UserPass returns the base64 encoded API username and password.
func (c *Config) Base64UserPass() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.APIUsername, c.APIPassword)))
}

func getEnv() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "staging"
	}
	return env
}
