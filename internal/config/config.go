package config

import "os"

type Config struct {
	Port        string
	AdminUser   string
	AdminPass   string
	StoragePath string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", ":8080"),
		AdminUser:   getEnv("ADMIN_USER", "admin"),
		AdminPass:   getEnv("ADMIN_PASS", "admin"),
		StoragePath: getEnv("STORAGE_PATH", "./data/posts.json"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}
