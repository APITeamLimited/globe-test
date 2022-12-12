package lib

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
				fmt.Printf("%s is not set, resorting to default value %s\n", key, defaultValue)
			}
		} else {
			panic(fmt.Sprintf("%s is not set and no default value was provided\n", key))
		}
	}

	return defaultValue
}
