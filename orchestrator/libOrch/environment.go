package libOrch

import (
	"fmt"
	"os"
)

func GetEnvVariable(key, defaultValue string) string ***REMOVED***
	SYSTEM_ENV := os.Getenv("SYSTEM_ENV")

	if SYSTEM_ENV == "" ***REMOVED***
		// Assume development environment
		SYSTEM_ENV = "development"
	***REMOVED***

	if value := os.Getenv(key); value != "" ***REMOVED***
		return value
	***REMOVED***

	if SYSTEM_ENV == "development" ***REMOVED***
		return defaultValue
	***REMOVED***

	panic(fmt.Sprintf("%s is not set, and environment was %s, not development", key, SYSTEM_ENV))
***REMOVED***
