package config

import "os"

type Config struct {
	PostsBaseURL string
	ServerPort   string
}

func Load() *Config {
	return &Config{
		PostsBaseURL: getEnvVal("POSTS_BASE_URL", "https://json-placeholder.mock.beeceptor.com"),
		ServerPort:   getEnvVal("SERVER_PORT", "9091"),
	}
}

func getEnvVal(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}
