package config

import (
	"os"
)

type Config struct {
	PublicHost  string
	Port        string
	Environment string
}

var (
	Env *Config
	MaxAttempts map[string]int
)

func initConfig() *Config {
	cfg := &Config{
		PublicHost:  getEnv("PUBLIC_HOST", "localhost"),
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("environment", "devlopment"),
	}
	return cfg
}

func init() { // Init is called automatically when the package is imported
	Env = initConfig()
	MaxAttempts = make(map[string]int)
	MaxAttempts["sending email"] = 5
}

func getEnv(key, fallback string) string {
	// A fallback is used if the environment variable is not set
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
