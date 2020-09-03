package logrus

import (
	"context"
	"io"
	"time"
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = New()
)

func StandardLogger() *Logger ***REMOVED***
	return std
***REMOVED***

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) ***REMOVED***
	std.SetOutput(out)
***REMOVED***

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) ***REMOVED***
	std.SetFormatter(formatter)
***REMOVED***

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func SetReportCaller(include bool) ***REMOVED***
	std.SetReportCaller(include)
***REMOVED***

// SetLevel sets the standard logger level.
func SetLevel(level Level) ***REMOVED***
	std.SetLevel(level)
***REMOVED***

// GetLevel returns the standard logger level.
func GetLevel() Level ***REMOVED***
	return std.GetLevel()
***REMOVED***

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func IsLevelEnabled(level Level) bool ***REMOVED***
	return std.IsLevelEnabled(level)
***REMOVED***

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook Hook) ***REMOVED***
	std.AddHook(hook)
***REMOVED***

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry ***REMOVED***
	return std.WithField(ErrorKey, err)
***REMOVED***

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *Entry ***REMOVED***
	return std.WithContext(ctx)
***REMOVED***

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface***REMOVED******REMOVED***) *Entry ***REMOVED***
	return std.WithField(key, value)
***REMOVED***

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry ***REMOVED***
	return std.WithFields(fields)
***REMOVED***

// WithTime creates an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithTime(t time.Time) *Entry ***REMOVED***
	return std.WithTime(t)
***REMOVED***

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Trace(args...)
***REMOVED***

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Debug(args...)
***REMOVED***

// Print logs a message at level Info on the standard logger.
func Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Print(args...)
***REMOVED***

// Info logs a message at level Info on the standard logger.
func Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Info(args...)
***REMOVED***

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Warn(args...)
***REMOVED***

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Warning(args...)
***REMOVED***

// Error logs a message at level Error on the standard logger.
func Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Error(args...)
***REMOVED***

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Panic(args...)
***REMOVED***

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Fatal(args...)
***REMOVED***

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Tracef(format, args...)
***REMOVED***

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Debugf(format, args...)
***REMOVED***

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Printf(format, args...)
***REMOVED***

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Infof(format, args...)
***REMOVED***

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Warnf(format, args...)
***REMOVED***

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Warningf(format, args...)
***REMOVED***

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Errorf(format, args...)
***REMOVED***

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Panicf(format, args...)
***REMOVED***

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Fatalf(format, args...)
***REMOVED***

// Traceln logs a message at level Trace on the standard logger.
func Traceln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Traceln(args...)
***REMOVED***

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Debugln(args...)
***REMOVED***

// Println logs a message at level Info on the standard logger.
func Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Println(args...)
***REMOVED***

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Infoln(args...)
***REMOVED***

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Warnln(args...)
***REMOVED***

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Warningln(args...)
***REMOVED***

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Errorln(args...)
***REMOVED***

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Panicln(args...)
***REMOVED***

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	std.Fatalln(args...)
***REMOVED***
