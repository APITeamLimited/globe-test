package logrus

import (
	"fmt"
	"log"
	"strings"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface***REMOVED******REMOVED***

// Level type
type Level uint32

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string ***REMOVED***
	switch level ***REMOVED***
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	***REMOVED***

	return "unknown"
***REMOVED***

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) ***REMOVED***
	switch strings.ToLower(lvl) ***REMOVED***
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	***REMOVED***

	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
***REMOVED***

// A constant exposing all logging levels
var AllLevels = []Level***REMOVED***
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
***REMOVED***

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

// Won't compile if StdLogger can't be realized by a log.Logger
var (
	_ StdLogger = &log.Logger***REMOVED******REMOVED***
	_ StdLogger = &Entry***REMOVED******REMOVED***
	_ StdLogger = &Logger***REMOVED******REMOVED***
)

// StdLogger is what your logrus-enabled library should take, that way
// it'll accept a stdlib logger and a logrus logger. There's no standard
// interface, this is the closest we get, unfortunately.
type StdLogger interface ***REMOVED***
	Print(...interface***REMOVED******REMOVED***)
	Printf(string, ...interface***REMOVED******REMOVED***)
	Println(...interface***REMOVED******REMOVED***)

	Fatal(...interface***REMOVED******REMOVED***)
	Fatalf(string, ...interface***REMOVED******REMOVED***)
	Fatalln(...interface***REMOVED******REMOVED***)

	Panic(...interface***REMOVED******REMOVED***)
	Panicf(string, ...interface***REMOVED******REMOVED***)
	Panicln(...interface***REMOVED******REMOVED***)
***REMOVED***

// The FieldLogger interface generalizes the Entry and Logger types
type FieldLogger interface ***REMOVED***
	WithField(key string, value interface***REMOVED******REMOVED***) *Entry
	WithFields(fields Fields) *Entry
	WithError(err error) *Entry

	Debugf(format string, args ...interface***REMOVED******REMOVED***)
	Infof(format string, args ...interface***REMOVED******REMOVED***)
	Printf(format string, args ...interface***REMOVED******REMOVED***)
	Warnf(format string, args ...interface***REMOVED******REMOVED***)
	Warningf(format string, args ...interface***REMOVED******REMOVED***)
	Errorf(format string, args ...interface***REMOVED******REMOVED***)
	Fatalf(format string, args ...interface***REMOVED******REMOVED***)
	Panicf(format string, args ...interface***REMOVED******REMOVED***)

	Debug(args ...interface***REMOVED******REMOVED***)
	Info(args ...interface***REMOVED******REMOVED***)
	Print(args ...interface***REMOVED******REMOVED***)
	Warn(args ...interface***REMOVED******REMOVED***)
	Warning(args ...interface***REMOVED******REMOVED***)
	Error(args ...interface***REMOVED******REMOVED***)
	Fatal(args ...interface***REMOVED******REMOVED***)
	Panic(args ...interface***REMOVED******REMOVED***)

	Debugln(args ...interface***REMOVED******REMOVED***)
	Infoln(args ...interface***REMOVED******REMOVED***)
	Println(args ...interface***REMOVED******REMOVED***)
	Warnln(args ...interface***REMOVED******REMOVED***)
	Warningln(args ...interface***REMOVED******REMOVED***)
	Errorln(args ...interface***REMOVED******REMOVED***)
	Fatalln(args ...interface***REMOVED******REMOVED***)
	Panicln(args ...interface***REMOVED******REMOVED***)
***REMOVED***
