package lib

import (
	"fmt"
	"os"
)

func GetEnvVariable(key, defaultValue string) string ***REMOVED***
	return GetEnvVariableRaw(key, defaultValue, false)
***REMOVED***

func GetEnvVariableRaw(key, defaultValue string, hideError bool) string ***REMOVED***
	SYSTEM_ENV := os.Getenv("SYSTEM_ENV")

	if SYSTEM_ENV == "" ***REMOVED***
		// Assume development environment
		SYSTEM_ENV = "development"
	***REMOVED***

	if value := os.Getenv(key); value != "" ***REMOVED***
		return value
	***REMOVED***

	if defaultValue != "" ***REMOVED***
		if !hideError ***REMOVED***
			fmt.Printf("%s is not set, using default value %s\n", key, defaultValue)
		***REMOVED***

		return defaultValue
	***REMOVED***

	panic(fmt.Sprintf("%s is not set and no default value was provided\n", key))
***REMOVED***
