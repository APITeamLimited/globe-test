package libOrch

import (
	"fmt"
	"os"
)

func GetEnvVariable(key, defaultValue string) string {
	return GetEnvVariableHideError(key, defaultValue, false)
}

func GetEnvVariableHideError(key, defaultValue string, hideError bool) string {
	SYSTEM_ENV := os.Getenv("SYSTEM_ENV")

	if SYSTEM_ENV == "" {
		// Assume development environment
		SYSTEM_ENV = "development"
	}

	if value := os.Getenv(key); value != "" {
		return value
	}

	if SYSTEM_ENV != "development" {
		if defaultValue != "" {
			if !hideError {
				fmt.Printf("%s is not set, and environment was %s, not development, resorting to default value %s\n", key, SYSTEM_ENV, defaultValue)
			}
		} else {
			panic(fmt.Sprintf("%s is not set, and environment was %s, not development, and no default value was provided", key, SYSTEM_ENV))
		}
	}

	return defaultValue
}
