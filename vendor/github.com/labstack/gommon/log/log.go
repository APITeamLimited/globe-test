package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"strconv"

	"github.com/mattn/go-isatty"
	"github.com/valyala/fasttemplate"

	"github.com/labstack/gommon/color"
)

type (
	Logger struct ***REMOVED***
		prefix     string
		level      Lvl
		output     io.Writer
		template   *fasttemplate.Template
		levels     []string
		color      *color.Color
		bufferPool sync.Pool
		mutex      sync.Mutex
	***REMOVED***

	Lvl uint8

	JSON map[string]interface***REMOVED******REMOVED***
)

const (
	DEBUG Lvl = iota + 1
	INFO
	WARN
	ERROR
	OFF
)

var (
	global        = New("-")
	defaultHeader = `***REMOVED***"time":"$***REMOVED***time_rfc3339_nano***REMOVED***","level":"$***REMOVED***level***REMOVED***","prefix":"$***REMOVED***prefix***REMOVED***",` +
		`"file":"$***REMOVED***short_file***REMOVED***","line":"$***REMOVED***line***REMOVED***"***REMOVED***`
)

func New(prefix string) (l *Logger) ***REMOVED***
	l = &Logger***REMOVED***
		level:    INFO,
		prefix:   prefix,
		template: l.newTemplate(defaultHeader),
		color:    color.New(),
		bufferPool: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED***
				return bytes.NewBuffer(make([]byte, 256))
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	l.initLevels()
	l.SetOutput(output())
	return
***REMOVED***

func (l *Logger) initLevels() ***REMOVED***
	l.levels = []string***REMOVED***
		"-",
		l.color.Blue("DEBUG"),
		l.color.Green("INFO"),
		l.color.Yellow("WARN"),
		l.color.Red("ERROR"),
	***REMOVED***
***REMOVED***

func (l *Logger) newTemplate(format string) *fasttemplate.Template ***REMOVED***
	return fasttemplate.New(format, "$***REMOVED***", "***REMOVED***")
***REMOVED***

func (l *Logger) DisableColor() ***REMOVED***
	l.color.Disable()
	l.initLevels()
***REMOVED***

func (l *Logger) EnableColor() ***REMOVED***
	l.color.Enable()
	l.initLevels()
***REMOVED***

func (l *Logger) Prefix() string ***REMOVED***
	return l.prefix
***REMOVED***

func (l *Logger) SetPrefix(p string) ***REMOVED***
	l.prefix = p
***REMOVED***

func (l *Logger) Level() Lvl ***REMOVED***
	return l.level
***REMOVED***

func (l *Logger) SetLevel(v Lvl) ***REMOVED***
	l.level = v
***REMOVED***

func (l *Logger) Output() io.Writer ***REMOVED***
	return l.output
***REMOVED***

func (l *Logger) SetOutput(w io.Writer) ***REMOVED***
	l.output = w
	if w, ok := w.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) ***REMOVED***
		l.DisableColor()
	***REMOVED***
***REMOVED***

func (l *Logger) Color() *color.Color ***REMOVED***
	return l.color
***REMOVED***

func (l *Logger) SetHeader(h string) ***REMOVED***
	l.template = l.newTemplate(h)
***REMOVED***

func (l *Logger) Print(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(0, "", i...)
	// fmt.Fprintln(l.output, i...)
***REMOVED***

func (l *Logger) Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(0, format, args...)
***REMOVED***

func (l *Logger) Printj(j JSON) ***REMOVED***
	l.log(0, "json", j)
***REMOVED***

func (l *Logger) Debug(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(DEBUG, "", i...)
***REMOVED***

func (l *Logger) Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(DEBUG, format, args...)
***REMOVED***

func (l *Logger) Debugj(j JSON) ***REMOVED***
	l.log(DEBUG, "json", j)
***REMOVED***

func (l *Logger) Info(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(INFO, "", i...)
***REMOVED***

func (l *Logger) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(INFO, format, args...)
***REMOVED***

func (l *Logger) Infoj(j JSON) ***REMOVED***
	l.log(INFO, "json", j)
***REMOVED***

func (l *Logger) Warn(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(WARN, "", i...)
***REMOVED***

func (l *Logger) Warnf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(WARN, format, args...)
***REMOVED***

func (l *Logger) Warnj(j JSON) ***REMOVED***
	l.log(WARN, "json", j)
***REMOVED***

func (l *Logger) Error(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(ERROR, "", i...)
***REMOVED***

func (l *Logger) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.log(ERROR, format, args...)
***REMOVED***

func (l *Logger) Errorj(j JSON) ***REMOVED***
	l.log(ERROR, "json", j)
***REMOVED***

func (l *Logger) Fatal(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Print(i...)
	os.Exit(1)
***REMOVED***

func (l *Logger) Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Printf(format, args...)
	os.Exit(1)
***REMOVED***

func (l *Logger) Fatalj(j JSON) ***REMOVED***
	l.Printj(j)
	os.Exit(1)
***REMOVED***

func (l *Logger) Panic(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Print(i...)
	panic(fmt.Sprint(i...))
***REMOVED***

func (l *Logger) Panicf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Printf(format, args...)
	panic(fmt.Sprintf(format, args))
***REMOVED***

func (l *Logger) Panicj(j JSON) ***REMOVED***
	l.Printj(j)
	panic(j)
***REMOVED***

func DisableColor() ***REMOVED***
	global.DisableColor()
***REMOVED***

func EnableColor() ***REMOVED***
	global.EnableColor()
***REMOVED***

func Prefix() string ***REMOVED***
	return global.Prefix()
***REMOVED***

func SetPrefix(p string) ***REMOVED***
	global.SetPrefix(p)
***REMOVED***

func Level() Lvl ***REMOVED***
	return global.Level()
***REMOVED***

func SetLevel(v Lvl) ***REMOVED***
	global.SetLevel(v)
***REMOVED***

func Output() io.Writer ***REMOVED***
	return global.Output()
***REMOVED***

func SetOutput(w io.Writer) ***REMOVED***
	global.SetOutput(w)
***REMOVED***

func SetHeader(h string) ***REMOVED***
	global.SetHeader(h)
***REMOVED***

func Print(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Print(i...)
***REMOVED***

func Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Printf(format, args...)
***REMOVED***

func Printj(j JSON) ***REMOVED***
	global.Printj(j)
***REMOVED***

func Debug(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Debug(i...)
***REMOVED***

func Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Debugf(format, args...)
***REMOVED***

func Debugj(j JSON) ***REMOVED***
	global.Debugj(j)
***REMOVED***

func Info(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Info(i...)
***REMOVED***

func Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Infof(format, args...)
***REMOVED***

func Infoj(j JSON) ***REMOVED***
	global.Infoj(j)
***REMOVED***

func Warn(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Warn(i...)
***REMOVED***

func Warnf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Warnf(format, args...)
***REMOVED***

func Warnj(j JSON) ***REMOVED***
	global.Warnj(j)
***REMOVED***

func Error(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Error(i...)
***REMOVED***

func Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Errorf(format, args...)
***REMOVED***

func Errorj(j JSON) ***REMOVED***
	global.Errorj(j)
***REMOVED***

func Fatal(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Fatal(i...)
***REMOVED***

func Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Fatalf(format, args...)
***REMOVED***

func Fatalj(j JSON) ***REMOVED***
	global.Fatalj(j)
***REMOVED***

func Panic(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Panic(i...)
***REMOVED***

func Panicf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	global.Panicf(format, args...)
***REMOVED***

func Panicj(j JSON) ***REMOVED***
	global.Panicj(j)
***REMOVED***

func (l *Logger) log(v Lvl, format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.mutex.Lock()
	defer l.mutex.Unlock()
	buf := l.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer l.bufferPool.Put(buf)
	_, file, line, _ := runtime.Caller(3)

	if v >= l.level || v == 0 ***REMOVED***
		message := ""
		if format == "" ***REMOVED***
			message = fmt.Sprint(args...)
		***REMOVED*** else if format == "json" ***REMOVED***
			b, err := json.Marshal(args[0])
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***
			message = string(b)
		***REMOVED*** else ***REMOVED***
			message = fmt.Sprintf(format, args...)
		***REMOVED***

		_, err := l.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) ***REMOVED***
			switch tag ***REMOVED***
			case "time_rfc3339":
				return w.Write([]byte(time.Now().Format(time.RFC3339)))
			case "time_rfc3339_nano":
				return w.Write([]byte(time.Now().Format(time.RFC3339Nano)))
			case "level":
				return w.Write([]byte(l.levels[v]))
			case "prefix":
				return w.Write([]byte(l.prefix))
			case "long_file":
				return w.Write([]byte(file))
			case "short_file":
				return w.Write([]byte(path.Base(file)))
			case "line":
				return w.Write([]byte(strconv.Itoa(line)))
			***REMOVED***
			return 0, nil
		***REMOVED***)

		if err == nil ***REMOVED***
			s := buf.String()
			i := buf.Len() - 1
			if s[i] == '***REMOVED***' ***REMOVED***
				// JSON header
				buf.Truncate(i)
				buf.WriteByte(',')
				if format == "json" ***REMOVED***
					buf.WriteString(message[1:])
				***REMOVED*** else ***REMOVED***
					buf.WriteString(`"message":`)
					buf.WriteString(strconv.Quote(message))
					buf.WriteString(`***REMOVED***`)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// Text header
				buf.WriteByte(' ')
				buf.WriteString(message)
			***REMOVED***
			buf.WriteByte('\n')
			l.output.Write(buf.Bytes())
		***REMOVED***
	***REMOVED***
***REMOVED***
