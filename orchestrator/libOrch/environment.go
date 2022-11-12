package libOrch

import (
	"fmt"
	"os"
)

func GetEnvVariable(key, defaultValue string) string {
	SYSTEM_ENV := os.Getenv("SYSTEM_ENV")

	if SYSTEM_ENV == "" {
		// Assume development environment
		SYSTEM_ENV = "development"
	}

	if value := os.Getenv(key); value != "" {
		return value
	}

	if SYSTEM_ENV != "development" {
		fmt.Printf("%s is not set, and environment was %s, not development\n", key, SYSTEM_ENV)
	}

	return defaultValue
}
