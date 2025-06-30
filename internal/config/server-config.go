package config

import (
	"os"
)

type Config struct {
	Port        string
	Environment string
	JobQueueKey   string
	JobHashKeyPrefix string
}

var (
	Env         *Config
	MaxAttempts map[string]int
)

func initConfig() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "devlopment"),
		JobQueueKey:   getEnv("JOB_QUEUE_KEY", "jobQueue"),
		JobHashKeyPrefix: getEnv("JOB_HASH_KEY_PREFIX","job:data:"),
	}
	return cfg
}

func initMaxAttempts() {
	MaxAttempts = make(map[string]int)
	MaxAttempts["sending email"] = 2
}

func init() { // Init is called automatically when the package is imported
	Env = initConfig()
	initMaxAttempts()
}

func getEnv(key, fallback string) string {
	// A fallback is used if the environment variable is not set
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
