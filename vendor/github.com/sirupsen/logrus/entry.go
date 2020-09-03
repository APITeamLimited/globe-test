package logrus

import (
	"bytes"
	"context"
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

	// Contains the context set by the user. Useful for hook processing etc.
	Context context.Context

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

// Returns the bytes representation of this entry from the formatter.
func (entry *Entry) Bytes() ([]byte, error) ***REMOVED***
	return entry.Logger.Formatter.Format(entry)
***REMOVED***

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) ***REMOVED***
	serialized, err := entry.Bytes()
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

// Add a context to the Entry.
func (entry *Entry) WithContext(ctx context.Context) *Entry ***REMOVED***
	dataCopy := make(Fields, len(entry.Data))
	for k, v := range entry.Data ***REMOVED***
		dataCopy[k] = v
	***REMOVED***
	return &Entry***REMOVED***Logger: entry.Logger, Data: dataCopy, Time: entry.Time, err: entry.err, Context: ctx***REMOVED***
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
	fieldErr := entry.err
	for k, v := range fields ***REMOVED***
		isErrField := false
		if t := reflect.TypeOf(v); t != nil ***REMOVED***
			switch t.Kind() ***REMOVED***
			case reflect.Func:
				isErrField = true
			case reflect.Ptr:
				isErrField = t.Elem().Kind() == reflect.Func
			***REMOVED***
		***REMOVED***
		if isErrField ***REMOVED***
			tmp := fmt.Sprintf("can not add field %q", k)
			if fieldErr != "" ***REMOVED***
				fieldErr = entry.err + ", " + tmp
			***REMOVED*** else ***REMOVED***
				fieldErr = tmp
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			data[k] = v
		***REMOVED***
	***REMOVED***
	return &Entry***REMOVED***Logger: entry.Logger, Data: data, Time: entry.Time, err: fieldErr, Context: entry.Context***REMOVED***
***REMOVED***

// Overrides the time of the Entry.
func (entry *Entry) WithTime(t time.Time) *Entry ***REMOVED***
	dataCopy := make(Fields, len(entry.Data))
	for k, v := range entry.Data ***REMOVED***
		dataCopy[k] = v
	***REMOVED***
	return &Entry***REMOVED***Logger: entry.Logger, Data: dataCopy, Time: t, err: entry.err, Context: entry.Context***REMOVED***
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
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() ***REMOVED***
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ ***REMOVED***
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") ***REMOVED***
				logrusPackage = getPackageName(funcName)
				break
			***REMOVED***
		***REMOVED***

		minimumCallerDepth = knownLogrusFrames
	***REMOVED***)

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() ***REMOVED***
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage ***REMOVED***
			return &f //nolint:scopelint
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
	entry.Logger.mu.Lock()
	if entry.Logger.ReportCaller ***REMOVED***
		entry.Caller = getCaller()
	***REMOVED***
	entry.Logger.mu.Unlock()

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
		return
	***REMOVED***
	if _, err = entry.Logger.Out.Write(serialized); err != nil ***REMOVED***
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	***REMOVED***
***REMOVED***

func (entry *Entry) Log(level Level, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(level) ***REMOVED***
		entry.log(level, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Trace(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Log(TraceLevel, args...)
***REMOVED***

func (entry *Entry) Debug(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Log(DebugLevel, args...)
***REMOVED***

func (entry *Entry) Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Info(args...)
***REMOVED***

func (entry *Entry) Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Log(InfoLevel, args...)
***REMOVED***

func (entry *Entry) Warn(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Log(WarnLevel, args...)
***REMOVED***

func (entry *Entry) Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Warn(args...)
***REMOVED***

func (entry *Entry) Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Log(ErrorLevel, args...)
***REMOVED***

func (entry *Entry) Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Log(FatalLevel, args...)
	entry.Logger.Exit(1)
***REMOVED***

func (entry *Entry) Panic(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Log(PanicLevel, args...)
	panic(fmt.Sprint(args...))
***REMOVED***

// Entry Printf family functions

func (entry *Entry) Logf(level Level, format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(level) ***REMOVED***
		entry.Log(level, fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Tracef(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logf(TraceLevel, format, args...)
***REMOVED***

func (entry *Entry) Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logf(DebugLevel, format, args...)
***REMOVED***

func (entry *Entry) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logf(InfoLevel, format, args...)
***REMOVED***

func (entry *Entry) Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Infof(format, args...)
***REMOVED***

func (entry *Entry) Warnf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logf(WarnLevel, format, args...)
***REMOVED***

func (entry *Entry) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Warnf(format, args...)
***REMOVED***

func (entry *Entry) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logf(ErrorLevel, format, args...)
***REMOVED***

func (entry *Entry) Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logf(FatalLevel, format, args...)
	entry.Logger.Exit(1)
***REMOVED***

func (entry *Entry) Panicf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logf(PanicLevel, format, args...)
***REMOVED***

// Entry Println family functions

func (entry *Entry) Logln(level Level, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.IsLevelEnabled(level) ***REMOVED***
		entry.Log(level, entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Traceln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logln(TraceLevel, args...)
***REMOVED***

func (entry *Entry) Debugln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logln(DebugLevel, args...)
***REMOVED***

func (entry *Entry) Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logln(InfoLevel, args...)
***REMOVED***

func (entry *Entry) Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Infoln(args...)
***REMOVED***

func (entry *Entry) Warnln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logln(WarnLevel, args...)
***REMOVED***

func (entry *Entry) Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Warnln(args...)
***REMOVED***

func (entry *Entry) Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logln(ErrorLevel, args...)
***REMOVED***

func (entry *Entry) Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logln(FatalLevel, args...)
	entry.Logger.Exit(1)
***REMOVED***

func (entry *Entry) Panicln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Logln(PanicLevel, args...)
***REMOVED***

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...interface***REMOVED******REMOVED***) string ***REMOVED***
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
***REMOVED***
