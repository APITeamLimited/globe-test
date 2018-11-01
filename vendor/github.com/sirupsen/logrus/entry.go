package logrus

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	bufferPool *sync.Pool

	// qualified package name, cached at first use
	logrusPackage string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

func init() ***REMOVED***
	bufferPool = &sync.Pool***REMOVED***
		New: func() interface***REMOVED******REMOVED*** ***REMOVED***
			return new(bytes.Buffer)
		***REMOVED***,
	***REMOVED***

	// start at the bottom of the stack before the package-name cache is primed
	minimumCallerDepth = 1
***REMOVED***

// Defines the key when adding errors using WithError.
var ErrorKey = "error"

// An entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField***REMOVED***,s***REMOVED***. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.
type Entry struct ***REMOVED***
	Logger *Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level Level

	// Calling method, with package name
	Caller *runtime.Frame

	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// When formatter is called in entry.log(), a Buffer may be set to entry
	Buffer *bytes.Buffer

	// err may contain a field formatting error
	err string
***REMOVED***

func NewEntry(logger *Logger) *Entry ***REMOVED***
	return &Entry***REMOVED***
		Logger: logger,
		// Default is three fields, plus one optional.  Give a little extra room.
		Data: make(Fields, 6),
	***REMOVED***
***REMOVED***

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) ***REMOVED***
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	str := string(serialized)
	return str, nil
***REMOVED***

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry ***REMOVED***
	return entry.WithField(ErrorKey, err)
***REMOVED***

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface***REMOVED******REMOVED***) *Entry ***REMOVED***
	return entry.WithFields(Fields***REMOVED***key: value***REMOVED***)
***REMOVED***

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry ***REMOVED***
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data ***REMOVED***
		data[k] = v
	***REMOVED***
	var field_err string
	for k, v := range fields ***REMOVED***
		if t := reflect.TypeOf(v); t != nil && t.Kind() == reflect.Func ***REMOVED***
			field_err = fmt.Sprintf("can not add field %q", k)
			if entry.err != "" ***REMOVED***
				field_err = entry.err + ", " + field_err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			data[k] = v
		***REMOVED***
	***REMOVED***
	return &Entry***REMOVED***Logger: entry.Logger, Data: data, Time: entry.Time, err: field_err***REMOVED***
***REMOVED***

// Overrides the time of the Entry.
func (entry *Entry) WithTime(t time.Time) *Entry ***REMOVED***
	return &Entry***REMOVED***Logger: entry.Logger, Data: entry.Data, Time: t***REMOVED***
***REMOVED***

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string ***REMOVED***
	for ***REMOVED***
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash ***REMOVED***
			f = f[:lastPeriod]
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return f
***REMOVED***

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame ***REMOVED***
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() ***REMOVED***
		logrusPackage = getPackageName(runtime.FuncForPC(pcs[0]).Name())

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		// XXX this is dubious, the number of frames may vary store an entry in a logger interface
		minimumCallerDepth = knownLogrusFrames
	***REMOVED***)

	for f, again := frames.Next(); again; f, again = frames.Next() ***REMOVED***
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage ***REMOVED***
			return &f
		***REMOVED***
	***REMOVED***

	// if we got here, we failed to find the caller's context
	return nil
***REMOVED***

func (entry Entry) HasCaller() (has bool) ***REMOVED***
	return entry.Logger != nil &&
		entry.Logger.ReportCaller &&
		entry.Caller != nil
***REMOVED***

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) log(level Level, msg string) ***REMOVED***
	var buffer *bytes.Buffer

	// Default to now, but allow users to override if they want.
	//
	// We don't have to worry about polluting future calls to Entry#log()
	// with this assignment because this function is declared with a
	// non-pointer receiver.
	if entry.Time.IsZero() ***REMOVED***
		entry.Time = time.Now()
	***REMOVED***

	entry.Level = level
	entry.Message = msg
	if entry.Logger.ReportCaller ***REMOVED***
		entry.Caller = getCaller()
	***REMOVED***

	entry.fireHooks()

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer

	entry.write()

	entry.Buffer = nil

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel ***REMOVED***
		panic(&entry)
	***REMOVED***
***REMOVED***

func (entry *Entry) fireHooks() ***REMOVED***
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	err := entry.Logger.Hooks.Fire(entry.Level, entry)
	if err != nil ***REMOVED***
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
	***REMOVED***
***REMOVED***

func (entry *Entry) write() ***REMOVED***
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil ***REMOVED***
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
	***REMOVED*** else ***REMOVED***
		_, err = entry.Logger.Out.Write(serialized)
		if err != nil ***REMOVED***
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (entry *Entry) Trace(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(TraceLevel) ***REMOVED***
		entry.log(TraceLevel, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Debug(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(DebugLevel) ***REMOVED***
		entry.log(DebugLevel, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Info(args...)
***REMOVED***

func (entry *Entry) Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(InfoLevel) ***REMOVED***
		entry.log(InfoLevel, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Warn(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(WarnLevel) ***REMOVED***
		entry.log(WarnLevel, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Warn(args...)
***REMOVED***

func (entry *Entry) Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(ErrorLevel) ***REMOVED***
		entry.log(ErrorLevel, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(FatalLevel) ***REMOVED***
		entry.log(FatalLevel, fmt.Sprint(args...))
	***REMOVED***
	entry.Logger.Exit(1)
***REMOVED***

func (entry *Entry) Panic(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(PanicLevel) ***REMOVED***
		entry.log(PanicLevel, fmt.Sprint(args...))
	***REMOVED***
	panic(fmt.Sprint(args...))
***REMOVED***

// Entry Printf family functions

func (entry *Entry) Tracef(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(TraceLevel) ***REMOVED***
		entry.Trace(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(DebugLevel) ***REMOVED***
		entry.Debug(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(InfoLevel) ***REMOVED***
		entry.Info(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Infof(format, args...)
***REMOVED***

func (entry *Entry) Warnf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(WarnLevel) ***REMOVED***
		entry.Warn(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Warnf(format, args...)
***REMOVED***

func (entry *Entry) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(ErrorLevel) ***REMOVED***
		entry.Error(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(FatalLevel) ***REMOVED***
		entry.Fatal(fmt.Sprintf(format, args...))
	***REMOVED***
	entry.Logger.Exit(1)
***REMOVED***

func (entry *Entry) Panicf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(PanicLevel) ***REMOVED***
		entry.Panic(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

// Entry Println family functions

func (entry *Entry) Traceln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(TraceLevel) ***REMOVED***
		entry.Trace(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Debugln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(DebugLevel) ***REMOVED***
		entry.Debug(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(InfoLevel) ***REMOVED***
		entry.Info(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Infoln(args...)
***REMOVED***

func (entry *Entry) Warnln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(WarnLevel) ***REMOVED***
		entry.Warn(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Warnln(args...)
***REMOVED***

func (entry *Entry) Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(ErrorLevel) ***REMOVED***
		entry.Error(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(FatalLevel) ***REMOVED***
		entry.Fatal(entry.sprintlnn(args...))
	***REMOVED***
	entry.Logger.Exit(1)
***REMOVED***

func (entry *Entry) Panicln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(PanicLevel) ***REMOVED***
		entry.Panic(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...interface***REMOVED******REMOVED***) string ***REMOVED***
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
***REMOVED***
