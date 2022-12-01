package lib

import (
	"fmt"
	"os"
)

func GetEnvVariable(key, defaultValue string) string ***REMOVED***
	return GetEnvVariableHideError(key, defaultValue, false)
***REMOVED***

func GetEnvVariableHideError(key, defaultValue string, hideError bool) string ***REMOVED***
	SYSTEM_ENV := os.Getenv("SYSTEM_ENV")

	if SYSTEM_ENV == "" ***REMOVED***
		// Assume development environment
		SYSTEM_ENV = "development"
	***REMOVED***

	if value := os.Getenv(key); value != "" ***REMOVED***
		return value
	***REMOVED***

	if SYSTEM_ENV != "development" ***REMOVED***
		if defaultValue != "" ***REMOVED***
			if !hideError ***REMOVED***
				fmt.Printf("%s is not set, and environment was %s, not development, resorting to default value %s\n", key, SYSTEM_ENV, defaultValue)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			panic(fmt.Sprintf("%s is not set, and environment was %s, not development, and no default value was provided", key, SYSTEM_ENV))
		***REMOVED***
	***REMOVED***

	return defaultValue
***REMOVED***