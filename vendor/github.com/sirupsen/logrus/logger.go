package logrus

import (
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type Logger struct ***REMOVED***
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventorous, such as logging to Kafka.
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
	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged.
	Level Level
	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap
	// Reusable empty entry
	entryPool sync.Pool
***REMOVED***

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
//    var log = &Logger***REMOVED***
//      Out: os.Stderr,
//      Formatter: new(JSONFormatter),
//      Hooks: make(LevelHooks),
//      Level: logrus.DebugLevel,
//    ***REMOVED***
//
// It's recommended to make this a global instance called `log`.
func New() *Logger ***REMOVED***
	return &Logger***REMOVED***
		Out:       os.Stderr,
		Formatter: new(TextFormatter),
		Hooks:     make(LevelHooks),
		Level:     InfoLevel,
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
	logger.entryPool.Put(entry)
***REMOVED***

// Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Fatal or Panic. It only creates a log entry.
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

func (logger *Logger) Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= DebugLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Debugf(format, args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= InfoLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Infof(format, args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry := logger.newEntry()
	entry.Printf(format, args...)
	logger.releaseEntry(entry)
***REMOVED***

func (logger *Logger) Warnf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= WarnLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Warnf(format, args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= WarnLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Warnf(format, args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= ErrorLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Errorf(format, args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= FatalLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Fatalf(format, args...)
		logger.releaseEntry(entry)
	***REMOVED***
	Exit(1)
***REMOVED***

func (logger *Logger) Panicf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= PanicLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Panicf(format, args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Debug(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= DebugLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Debug(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= InfoLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Info(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry := logger.newEntry()
	entry.Info(args...)
	logger.releaseEntry(entry)
***REMOVED***

func (logger *Logger) Warn(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= WarnLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Warn(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= WarnLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Warn(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= ErrorLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Error(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= FatalLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Fatal(args...)
		logger.releaseEntry(entry)
	***REMOVED***
	Exit(1)
***REMOVED***

func (logger *Logger) Panic(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= PanicLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Panic(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Debugln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= DebugLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Debugln(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= InfoLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Infoln(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry := logger.newEntry()
	entry.Println(args...)
	logger.releaseEntry(entry)
***REMOVED***

func (logger *Logger) Warnln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= WarnLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Warnln(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= WarnLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Warnln(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= ErrorLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Errorln(args...)
		logger.releaseEntry(entry)
	***REMOVED***
***REMOVED***

func (logger *Logger) Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= FatalLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Fatalln(args...)
		logger.releaseEntry(entry)
	***REMOVED***
	Exit(1)
***REMOVED***

func (logger *Logger) Panicln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger.level() >= PanicLevel ***REMOVED***
		entry := logger.newEntry()
		entry.Panicln(args...)
		logger.releaseEntry(entry)
	***REMOVED***
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

func (logger *Logger) SetLevel(level Level) ***REMOVED***
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
***REMOVED***

func (logger *Logger) AddHook(hook Hook) ***REMOVED***
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Hooks.Add(hook)
***REMOVED***
