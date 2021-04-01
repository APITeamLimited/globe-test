package logrus

import (
	"context"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// LogFunction For big messages, it can be more efficient to pass a function
// and only call it if the log level is actually enables rather than
// generating the log message and then checking if the level is enabled
type LogFunction func() []interface***REMOVED******REMOVED***

type Logger struct ***REMOVED***
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventurous, such as logging to Kafka.
	Out io.Writer
	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	Hooks LevelHooks
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter Formatter

	// Flag for whether to log caller info (off by default)
	ReportCaller bool

	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged.
	Level Level
	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap
	// Reusable empty entry
	entryPool sync.Pool
	// Function to exit the application, defaults to `os.Exit()`
	ExitFunc exitFunc
***REMOVED***

type exitFunc func(int)

type MutexWrap struct ***REMOVED***
	lock     sync.Mutex
	disabled bool
***REMOVED***

func (mw *MutexWrap) Lock() ***REMOVED***
	if !mw.disabled ***REMOVED***
		mw.lock.Lock()
	***REMOVED***
***REMOVED***

func (mw *MutexWrap) Unlock() ***REMOVED***
	if !mw.disabled ***REMOVED***
		mw.lock.Unlock()
	***REMOVED***
***REMOVED***

func (mw *MutexWrap) Disable() ***REMOVED***
	mw.disabled = true
***REMOVED***

// Creates a new logger. Configuration should be set by changing `Formatter`,
// `Out` and `Hooks` directly on the default logger instance. You can also just
// instantiate your own:
//
//    var log = &logrus.Logger***REMOVED***
//      Out: os.Stderr,
//      Formatter: new(logrus.TextFormatter),
//      Hooks: make(logrus.LevelHooks),
//      Level: logrus.DebugLevel,
//    ***REMOVED***
//
// It's recommended to make this a global instance called `log`.
func New() *Logger ***REMOVED***
	return &Logger***REMOVED***
		Out:          os.Stderr,
		Formatter:    new(TextFormatter),
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	***REMOVED***
***REMOVED***

func (logger *Logger) newEntry() *Entry ***REMOVED***
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok ***REMOVED***
		return entry
	***REMOVED***
	return NewEntry(logger)
***REMOVED***

func (logger *Logger) releaseEntry(entry *Entry) ***REMOVED***
	entry.Data = map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	logger.entryPool.Put(entry)
***REMOVED***

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface***REMOVED******REMOVED***) *Entry ***REMOVED***
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithField(key, value)
***REMOVED***

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields Fields) *Entry ***REMOVED***
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields)
***REMOVED***

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *Entry ***REMOVED***
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
***REMOVED***

// Add a context to the log entry.
func (logger *Logger) WithContext(ctx context.Context) *Entry ***REMOVED***
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithContext(ctx)
***REMOVED***

// Overrides the time of the log entry.
func (logger *Logger) WithTime(t time.Time) *Entry ***REMOVED***
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithTime(t)
***REMOVED***

func (logger *Logger) Logf(level Level, format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.IsLevelEnabled(level) ***REMOVED***
		entry := logger.newEntry()
		entry.Logf(level, format, args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Tracef(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logf(TraceLevel, format, args...)
***REMOVED***

func (logger *Logger) Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logf(DebugLevel, format, args...)
***REMOVED***

func (logger *Logger) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logf(InfoLevel, format, args...)
***REMOVED***

func (logger *Logger) Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry := logger.newEntry()
	entry.Printf(format, args...)
	logger.releaseEntry(entry)
***REMOVED***

func (logger *Logger) Warnf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logf(WarnLevel, format, args...)
***REMOVED***

func (logger *Logger) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Warnf(format, args...)
***REMOVED***

func (logger *Logger) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logf(ErrorLevel, format, args...)
***REMOVED***

func (logger *Logger) Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logf(FatalLevel, format, args...)
	logger.Exit(1)
***REMOVED***

func (logger *Logger) Panicf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logf(PanicLevel, format, args...)
***REMOVED***

func (logger *Logger) Log(level Level, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.IsLevelEnabled(level) ***REMOVED***
		entry := logger.newEntry()
		entry.Log(level, args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) LogFn(level Level, fn LogFunction) ***REMOVED***
	if logger.IsLevelEnabled(level) ***REMOVED***
		entry := logger.newEntry()
		entry.Log(level, fn()...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Trace(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Log(TraceLevel, args...)
***REMOVED***

func (logger *Logger) Debug(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Log(DebugLevel, args...)
***REMOVED***

func (logger *Logger) Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Log(InfoLevel, args...)
***REMOVED***

func (logger *Logger) Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry := logger.newEntry()
	entry.Print(args...)
	logger.releaseEntry(entry)
***REMOVED***

func (logger *Logger) Warn(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Log(WarnLevel, args...)
***REMOVED***

func (logger *Logger) Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Warn(args...)
***REMOVED***

func (logger *Logger) Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Log(ErrorLevel, args...)
***REMOVED***

func (logger *Logger) Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Log(FatalLevel, args...)
	logger.Exit(1)
***REMOVED***

func (logger *Logger) Panic(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Log(PanicLevel, args...)
***REMOVED***

func (logger *Logger) TraceFn(fn LogFunction) ***REMOVED***
	logger.LogFn(TraceLevel, fn)
***REMOVED***

func (logger *Logger) DebugFn(fn LogFunction) ***REMOVED***
	logger.LogFn(DebugLevel, fn)
***REMOVED***

func (logger *Logger) InfoFn(fn LogFunction) ***REMOVED***
	logger.LogFn(InfoLevel, fn)
***REMOVED***

func (logger *Logger) PrintFn(fn LogFunction) ***REMOVED***
	entry := logger.newEntry()
	entry.Print(fn()...)
	logger.releaseEntry(entry)
***REMOVED***

func (logger *Logger) WarnFn(fn LogFunction) ***REMOVED***
	logger.LogFn(WarnLevel, fn)
***REMOVED***

func (logger *Logger) WarningFn(fn LogFunction) ***REMOVED***
	logger.WarnFn(fn)
***REMOVED***

func (logger *Logger) ErrorFn(fn LogFunction) ***REMOVED***
	logger.LogFn(ErrorLevel, fn)
***REMOVED***

func (logger *Logger) FatalFn(fn LogFunction) ***REMOVED***
	logger.LogFn(FatalLevel, fn)
	logger.Exit(1)
***REMOVED***

func (logger *Logger) PanicFn(fn LogFunction) ***REMOVED***
	logger.LogFn(PanicLevel, fn)
***REMOVED***

func (logger *Logger) Logln(level Level, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.IsLevelEnabled(level) ***REMOVED***
		entry := logger.newEntry()
		entry.Logln(level, args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Traceln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logln(TraceLevel, args...)
***REMOVED***

func (logger *Logger) Debugln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logln(DebugLevel, args...)
***REMOVED***

func (logger *Logger) Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logln(InfoLevel, args...)
***REMOVED***

func (logger *Logger) Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry := logger.newEntry()
	entry.Println(args...)
	logger.releaseEntry(entry)
***REMOVED***

func (logger *Logger) Warnln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logln(WarnLevel, args...)
***REMOVED***

func (logger *Logger) Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Warnln(args...)
***REMOVED***

func (logger *Logger) Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logln(ErrorLevel, args...)
***REMOVED***

func (logger *Logger) Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logln(FatalLevel, args...)
	logger.Exit(1)
***REMOVED***

func (logger *Logger) Panicln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Logln(PanicLevel, args...)
***REMOVED***

func (logger *Logger) Exit(code int) ***REMOVED***
	runHandlers()
	if logger.ExitFunc == nil ***REMOVED***
		logger.ExitFunc = os.Exit
	***REMOVED***
	logger.ExitFunc(code)
***REMOVED***

//When file is opened with appending mode, it's safe to
//write concurrently to a file (within 4k message on Linux).
//In these cases user can choose to disable the lock.
func (logger *Logger) SetNoLock() ***REMOVED***
	logger.mu.Disable()
***REMOVED***

func (logger *Logger) level() Level ***REMOVED***
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
***REMOVED***

// SetLevel sets the logger level.
func (logger *Logger) SetLevel(level Level) ***REMOVED***
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
***REMOVED***

// GetLevel returns the logger level.
func (logger *Logger) GetLevel() Level ***REMOVED***
	return logger.level()
***REMOVED***

// AddHook adds a hook to the logger hooks.
func (logger *Logger) AddHook(hook Hook) ***REMOVED***
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Hooks.Add(hook)
***REMOVED***

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (logger *Logger) IsLevelEnabled(level Level) bool ***REMOVED***
	return logger.level() >= level
***REMOVED***

// SetFormatter sets the logger formatter.
func (logger *Logger) SetFormatter(formatter Formatter) ***REMOVED***
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Formatter = formatter
***REMOVED***

// SetOutput sets the logger output.
func (logger *Logger) SetOutput(output io.Writer) ***REMOVED***
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Out = output
***REMOVED***

func (logger *Logger) SetReportCaller(reportCaller bool) ***REMOVED***
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.ReportCaller = reportCaller
***REMOVED***

// ReplaceHooks replaces the logger hooks and returns the old ones
func (logger *Logger) ReplaceHooks(hooks LevelHooks) LevelHooks ***REMOVED***
	logger.mu.Lock()
	oldHooks := logger.Hooks
	logger.Hooks = hooks
	logger.mu.Unlock()
	return oldHooks
***REMOVED***
