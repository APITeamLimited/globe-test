package logrus

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 36
	gray    = 37
)

var (
	baseTimestamp time.Time
)

func init() ***REMOVED***
	baseTimestamp = time.Now()
***REMOVED***

// TextFormatter formats logs into text
type TextFormatter struct ***REMOVED***
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	sync.Once
***REMOVED***

func (f *TextFormatter) init(entry *Entry) ***REMOVED***
	if entry.Logger != nil ***REMOVED***
		f.isTerminal = f.checkIfTerminal(entry.Logger.Out)
	***REMOVED***
***REMOVED***

func (f *TextFormatter) checkIfTerminal(w io.Writer) bool ***REMOVED***
	switch v := w.(type) ***REMOVED***
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	***REMOVED***
***REMOVED***

// Format renders a single log entry
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) ***REMOVED***
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data ***REMOVED***
		keys = append(keys, k)
	***REMOVED***

	if !f.DisableSorting ***REMOVED***
		sort.Strings(keys)
	***REMOVED***
	if entry.Buffer != nil ***REMOVED***
		b = entry.Buffer
	***REMOVED*** else ***REMOVED***
		b = &bytes.Buffer***REMOVED******REMOVED***
	***REMOVED***

	prefixFieldClashes(entry.Data)

	f.Do(func() ***REMOVED*** f.init(entry) ***REMOVED***)

	isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" ***REMOVED***
		timestampFormat = defaultTimestampFormat
	***REMOVED***
	if isColored ***REMOVED***
		f.printColored(b, entry, keys, timestampFormat)
	***REMOVED*** else ***REMOVED***
		if !f.DisableTimestamp ***REMOVED***
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat))
		***REMOVED***
		f.appendKeyValue(b, "level", entry.Level.String())
		if entry.Message != "" ***REMOVED***
			f.appendKeyValue(b, "msg", entry.Message)
		***REMOVED***
		for _, key := range keys ***REMOVED***
			f.appendKeyValue(b, key, entry.Data[key])
		***REMOVED***
	***REMOVED***

	b.WriteByte('\n')
	return b.Bytes(), nil
***REMOVED***

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *Entry, keys []string, timestampFormat string) ***REMOVED***
	var levelColor int
	switch entry.Level ***REMOVED***
	case DebugLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	***REMOVED***

	levelText := strings.ToUpper(entry.Level.String())[0:4]

	if f.DisableTimestamp ***REMOVED***
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m %-44s ", levelColor, levelText, entry.Message)
	***REMOVED*** else if !f.FullTimestamp ***REMOVED***
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d] %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message)
	***REMOVED*** else ***REMOVED***
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s] %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), entry.Message)
	***REMOVED***
	for _, k := range keys ***REMOVED***
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	***REMOVED***
***REMOVED***

func (f *TextFormatter) needsQuoting(text string) bool ***REMOVED***
	if f.QuoteEmptyFields && len(text) == 0 ***REMOVED***
		return true
	***REMOVED***
	for _, ch := range text ***REMOVED***
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface***REMOVED******REMOVED***) ***REMOVED***
	if b.Len() > 0 ***REMOVED***
		b.WriteByte(' ')
	***REMOVED***
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
***REMOVED***

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface***REMOVED******REMOVED***) ***REMOVED***
	stringVal, ok := value.(string)
	if !ok ***REMOVED***
		stringVal = fmt.Sprint(value)
	***REMOVED***

	if !f.needsQuoting(stringVal) ***REMOVED***
		b.WriteString(stringVal)
	***REMOVED*** else ***REMOVED***
		b.WriteString(fmt.Sprintf("%q", stringVal))
	***REMOVED***
***REMOVED***
