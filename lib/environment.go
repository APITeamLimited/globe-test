package lib

import (
	"encoding/hex"
	"fmt"
	"os"
)

func GetEnvVariableBool(key string, defaultValue bool) bool {
	value := GetEnvVariableRaw(key, "None", true)
	if value == "None" {
		return defaultValue
	}

	if value == "true" {
		return true
	}

	if value == "false" {
		return false
	}

	panic(fmt.Sprintf("%s is not a valid boolean value", value))
}

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

func GetHexEnvVariable(key, defaultValue string) []byte {
	// Ensure key ends in _HEX
	if len(key) < 4 || key[len(key)-4:] != "_HEX" {
		panic(fmt.Sprintf("GetHexEnvVariable: key %s must end in _HEX", key))
	}

	valueHex := GetEnvVariable(key, defaultValue)

	value, err := hex.DecodeString(valueHex)
	if err != nil {
		panic(err)
	}

	return value
}
