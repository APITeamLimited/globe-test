package logrus

import (
	"bufio"
	"io"
	"runtime"
)

func (logger *Logger) Writer() *io.PipeWriter ***REMOVED***
	return logger.WriterLevel(InfoLevel)
***REMOVED***

func (logger *Logger) WriterLevel(level Level) *io.PipeWriter ***REMOVED***
	return NewEntry(logger).WriterLevel(level)
***REMOVED***

func (entry *Entry) Writer() *io.PipeWriter ***REMOVED***
	return entry.WriterLevel(InfoLevel)
***REMOVED***

func (entry *Entry) WriterLevel(level Level) *io.PipeWriter ***REMOVED***
	reader, writer := io.Pipe()

	var printFunc func(args ...interface***REMOVED******REMOVED***)

	switch level ***REMOVED***
	case TraceLevel:
		printFunc = entry.Trace
	case DebugLevel:
		printFunc = entry.Debug
	case InfoLevel:
		printFunc = entry.Info
	case WarnLevel:
		printFunc = entry.Warn
	case ErrorLevel:
		printFunc = entry.Error
	case FatalLevel:
		printFunc = entry.Fatal
	case PanicLevel:
		printFunc = entry.Panic
	default:
		printFunc = entry.Print
	***REMOVED***

	go entry.writerScanner(reader, printFunc)
	runtime.SetFinalizer(writer, writerFinalizer)

	return writer
***REMOVED***

func (entry *Entry) writerScanner(reader *io.PipeReader, printFunc func(args ...interface***REMOVED******REMOVED***)) ***REMOVED***
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() ***REMOVED***
		printFunc(scanner.Text())
	***REMOVED***
	if err := scanner.Err(); err != nil ***REMOVED***
		entry.Errorf("Error while reading from Writer: %s", err)
	***REMOVED***
	reader.Close()
***REMOVED***

func writerFinalizer(writer *io.PipeWriter) ***REMOVED***
	writer.Close()
***REMOVED***
