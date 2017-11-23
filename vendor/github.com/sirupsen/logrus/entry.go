package logrus

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

var bufferPool *sync.Pool

func init() ***REMOVED***
	bufferPool = &sync.Pool***REMOVED***
		New: func() interface***REMOVED******REMOVED*** ***REMOVED***
			return new(bytes.Buffer)
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Defines the key when adding errors using WithError.
var ErrorKey = "error"

// An entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField***REMOVED***,s***REMOVED***. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
type Entry struct ***REMOVED***
	Logger *Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// When formatter is called in entry.log(), an Buffer may be set to entry
	Buffer *bytes.Buffer
***REMOVED***

func NewEntry(logger *Logger) *Entry ***REMOVED***
	return &Entry***REMOVED***
		Logger: logger,
		// Default is three fields, give a little extra room
		Data: make(Fields, 5),
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
	for k, v := range fields ***REMOVED***
		data[k] = v
	***REMOVED***
	return &Entry***REMOVED***Logger: entry.Logger, Data: data***REMOVED***
***REMOVED***

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) log(level Level, msg string) ***REMOVED***
	var buffer *bytes.Buffer
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg

	entry.Logger.mu.Lock()
	err := entry.Logger.Hooks.Fire(level, &entry)
	entry.Logger.mu.Unlock()
	if err != nil ***REMOVED***
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
		entry.Logger.mu.Unlock()
	***REMOVED***
	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer
	serialized, err := entry.Logger.Formatter.Format(&entry)
	entry.Buffer = nil
	if err != nil ***REMOVED***
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		entry.Logger.mu.Unlock()
	***REMOVED*** else ***REMOVED***
		entry.Logger.mu.Lock()
		_, err = entry.Logger.Out.Write(serialized)
		if err != nil ***REMOVED***
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		***REMOVED***
		entry.Logger.mu.Unlock()
	***REMOVED***

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel ***REMOVED***
		panic(&entry)
	***REMOVED***
***REMOVED***

func (entry *Entry) Debug(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= DebugLevel ***REMOVED***
		entry.log(DebugLevel, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Info(args...)
***REMOVED***

func (entry *Entry) Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= InfoLevel ***REMOVED***
		entry.log(InfoLevel, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Warn(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= WarnLevel ***REMOVED***
		entry.log(WarnLevel, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Warn(args...)
***REMOVED***

func (entry *Entry) Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= ErrorLevel ***REMOVED***
		entry.log(ErrorLevel, fmt.Sprint(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= FatalLevel ***REMOVED***
		entry.log(FatalLevel, fmt.Sprint(args...))
	***REMOVED***
	Exit(1)
***REMOVED***

func (entry *Entry) Panic(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= PanicLevel ***REMOVED***
		entry.log(PanicLevel, fmt.Sprint(args...))
	***REMOVED***
	panic(fmt.Sprint(args...))
***REMOVED***

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= DebugLevel ***REMOVED***
		entry.Debug(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= InfoLevel ***REMOVED***
		entry.Info(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Infof(format, args...)
***REMOVED***

func (entry *Entry) Warnf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= WarnLevel ***REMOVED***
		entry.Warn(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Warnf(format, args...)
***REMOVED***

func (entry *Entry) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= ErrorLevel ***REMOVED***
		entry.Error(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= FatalLevel ***REMOVED***
		entry.Fatal(fmt.Sprintf(format, args...))
	***REMOVED***
	Exit(1)
***REMOVED***

func (entry *Entry) Panicf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= PanicLevel ***REMOVED***
		entry.Panic(fmt.Sprintf(format, args...))
	***REMOVED***
***REMOVED***

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= DebugLevel ***REMOVED***
		entry.Debug(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= InfoLevel ***REMOVED***
		entry.Info(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Infoln(args...)
***REMOVED***

func (entry *Entry) Warnln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= WarnLevel ***REMOVED***
		entry.Warn(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	entry.Warnln(args...)
***REMOVED***

func (entry *Entry) Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= ErrorLevel ***REMOVED***
		entry.Error(entry.sprintlnn(args...))
	***REMOVED***
***REMOVED***

func (entry *Entry) Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= FatalLevel ***REMOVED***
		entry.Fatal(entry.sprintlnn(args...))
	***REMOVED***
	Exit(1)
***REMOVED***

func (entry *Entry) Panicln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if entry.Logger.level() >= PanicLevel ***REMOVED***
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
