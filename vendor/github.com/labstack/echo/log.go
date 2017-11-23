package echo

import (
	"io"

	"github.com/labstack/gommon/log"
)

type (
	// Logger defines the logging interface.
	Logger interface ***REMOVED***
		Output() io.Writer
		SetOutput(w io.Writer)
		Prefix() string
		SetPrefix(p string)
		Level() log.Lvl
		SetLevel(v log.Lvl)
		Print(i ...interface***REMOVED******REMOVED***)
		Printf(format string, args ...interface***REMOVED******REMOVED***)
		Printj(j log.JSON)
		Debug(i ...interface***REMOVED******REMOVED***)
		Debugf(format string, args ...interface***REMOVED******REMOVED***)
		Debugj(j log.JSON)
		Info(i ...interface***REMOVED******REMOVED***)
		Infof(format string, args ...interface***REMOVED******REMOVED***)
		Infoj(j log.JSON)
		Warn(i ...interface***REMOVED******REMOVED***)
		Warnf(format string, args ...interface***REMOVED******REMOVED***)
		Warnj(j log.JSON)
		Error(i ...interface***REMOVED******REMOVED***)
		Errorf(format string, args ...interface***REMOVED******REMOVED***)
		Errorj(j log.JSON)
		Fatal(i ...interface***REMOVED******REMOVED***)
		Fatalj(j log.JSON)
		Fatalf(format string, args ...interface***REMOVED******REMOVED***)
		Panic(i ...interface***REMOVED******REMOVED***)
		Panicj(j log.JSON)
		Panicf(format string, args ...interface***REMOVED******REMOVED***)
	***REMOVED***
)
