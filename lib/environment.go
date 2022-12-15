package lib

import (
	"fmt"
	"os"
)

func GetEnvVariable(key, defaultValue string) string {
	return GetEnvVariableRaw(key, defaultValue, false)
}

func GetEnvVariableRaw(key, defaultValue string, hideError bool) string {
	SYSTEM_ENV := os.Getenv("SYSTEM_ENV")

	if SYSTEM_ENV == "" {
		// Assume development environment
		SYSTEM_ENV = "development"
	}

	if value := os.Getenv(key); value != "" {
		return value
	}

	if defaultValue != "" {
		if !hideError {
			fmt.Printf("%s is not set, using default value %s\n", key, defaultValue)
		}

		return defaultValue
	}

	panic(fmt.Sprintf("%s is not set and no default value was provided\n", key))
}
