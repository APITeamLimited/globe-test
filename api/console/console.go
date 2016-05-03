package console

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
)

var members = map[string]interface***REMOVED******REMOVED******REMOVED***
	"log":   Log,
	"warn":  Warn,
	"error": Error,
***REMOVED***

func New() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	return members
***REMOVED***

func consoleLogFields(args []interface***REMOVED******REMOVED***) log.Fields ***REMOVED***
	fields := log.Fields***REMOVED******REMOVED***
	for i, arg := range args ***REMOVED***
		fields[strconv.Itoa(i+1)] = arg
	***REMOVED***
	return fields
***REMOVED***

// TODO: Match console.log()'s sprintf()-like formatting behavior
func Log(msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	log.WithFields(consoleLogFields(args)).Info(msg)
***REMOVED***

func Warn(msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	log.WithFields(consoleLogFields(args)).Warn(msg)
***REMOVED***

func Error(msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	log.WithFields(consoleLogFields(args)).Error(msg)
***REMOVED***
