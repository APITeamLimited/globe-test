package v8js

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
)

func consoleLogFields(args []interface***REMOVED******REMOVED***) log.Fields ***REMOVED***
	fields := log.Fields***REMOVED******REMOVED***
	for i, arg := range args ***REMOVED***
		fields[strconv.Itoa(i+1)] = arg
	***REMOVED***
	return fields
***REMOVED***

// TODO: Match console.log()'s sprintf()-like formatting behavior
func (vu *VUContext) ConsoleLog(msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	log.WithFields(consoleLogFields(args)).Info(msg)
***REMOVED***

func (vu *VUContext) ConsoleWarn(msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	log.WithFields(consoleLogFields(args)).Warn(msg)
***REMOVED***

func (vu *VUContext) ConsoleError(msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	log.WithFields(consoleLogFields(args)).Error(msg)
***REMOVED***
