package environment

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

	if SYSTEM_ENV == "development" {
		return defaultValue
	}

	panic(fmt.Sprintf("%s is not set, and environment was %s, not development", key, SYSTEM_ENV))
}
