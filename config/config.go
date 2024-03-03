package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	NodeEnv           string
	MongoDBURI        string
	MongoDBCollection string
	CohereAPIKey      string
	CohereAPIPrompt   string
	GPTModel          string
	NYTimesBaseURL    string
	NYTimesAPIKey     string
	NewsAPIBaseURL    string
	NewsAPIKey        string
	GhostKey          string
	GhostURL          string
}

// NewConfig creates a new Config instance with values loaded from environment variables
func NewConfig() *Config {
	return &Config{
		NodeEnv:           getEnv("NODE_ENV", "development"),
		MongoDBURI:        getEnv("MONGO_URL", "mongodb://localhost/"),
		MongoDBCollection: getEnv("MONGO_DB_COLLECTION", ""),
		CohereAPIKey:      getEnv("COHERE_API_KEY", ""),
		CohereAPIPrompt:   getEnv("COHERE_API_PROMPT", ""),
		GPTModel:          getEnv("GPT_MODEL", ""),
		NYTimesBaseURL:    getEnv("NYTIMES_BASE_URL", ""),
		NYTimesAPIKey:     getEnv("NYTIMES_API_KEY", ""),
		NewsAPIBaseURL:    getEnv("NEWSAPI_BASE_URL", ""),
		NewsAPIKey:        getEnv("NEWSAPI_KEY", ""),
		GhostKey:          getEnv("GHOST_KEY", ""),
		GhostURL:          getEnv("GHOST_URL", "http://localhost:8080"),
	}
}

// getEnv retrieves the value of the environment variable with the given key,
// falling back to the provided default value if the variable is not set.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
