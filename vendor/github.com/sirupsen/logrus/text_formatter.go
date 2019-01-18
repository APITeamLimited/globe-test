package logrus

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
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
	emptyFieldMap FieldMap
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

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

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

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &TextFormatter***REMOVED***
	//     FieldMap: FieldMap***REMOVED***
	//         FieldKeyTime:  "@timestamp",
	//         FieldKeyLevel: "@level",
	//         FieldKeyMsg:   "@message"***REMOVED******REMOVED***
	FieldMap FieldMap

	terminalInitOnce sync.Once
***REMOVED***

func (f *TextFormatter) init(entry *Entry) ***REMOVED***
	if entry.Logger != nil ***REMOVED***
		f.isTerminal = checkIfTerminal(entry.Logger.Out)

		if f.isTerminal ***REMOVED***
			initTerminal(entry.Logger.Out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *TextFormatter) isColored() bool ***REMOVED***
	isColored := f.ForceColors || f.isTerminal

	if f.EnvironmentOverrideColors ***REMOVED***
		if force, ok := os.LookupEnv("CLICOLOR_FORCE"); ok && force != "0" ***REMOVED***
			isColored = true
		***REMOVED*** else if ok && force == "0" ***REMOVED***
			isColored = false
		***REMOVED*** else if os.Getenv("CLICOLOR") == "0" ***REMOVED***
			isColored = false
		***REMOVED***
	***REMOVED***

	return isColored && !f.DisableColors
***REMOVED***

// Format renders a single log entry
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) ***REMOVED***
	prefixFieldClashes(entry.Data, f.FieldMap, entry.HasCaller())

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data ***REMOVED***
		keys = append(keys, k)
	***REMOVED***

	fixedKeys := make([]string, 0, 4+len(entry.Data))
	if !f.DisableTimestamp ***REMOVED***
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyTime))
	***REMOVED***
	fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyLevel))
	if entry.Message != "" ***REMOVED***
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyMsg))
	***REMOVED***
	if entry.err != "" ***REMOVED***
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyLogrusError))
	***REMOVED***
	if entry.HasCaller() ***REMOVED***
		fixedKeys = append(fixedKeys,
			f.FieldMap.resolve(FieldKeyFunc), f.FieldMap.resolve(FieldKeyFile))
	***REMOVED***

	if !f.DisableSorting ***REMOVED***
		if f.SortingFunc == nil ***REMOVED***
			sort.Strings(keys)
			fixedKeys = append(fixedKeys, keys...)
		***REMOVED*** else ***REMOVED***
			if !f.isColored() ***REMOVED***
				fixedKeys = append(fixedKeys, keys...)
				f.SortingFunc(fixedKeys)
			***REMOVED*** else ***REMOVED***
				f.SortingFunc(keys)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fixedKeys = append(fixedKeys, keys...)
	***REMOVED***

	var b *bytes.Buffer
	if entry.Buffer != nil ***REMOVED***
		b = entry.Buffer
	***REMOVED*** else ***REMOVED***
		b = &bytes.Buffer***REMOVED******REMOVED***
	***REMOVED***

	f.terminalInitOnce.Do(func() ***REMOVED*** f.init(entry) ***REMOVED***)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" ***REMOVED***
		timestampFormat = defaultTimestampFormat
	***REMOVED***
	if f.isColored() ***REMOVED***
		f.printColored(b, entry, keys, timestampFormat)
	***REMOVED*** else ***REMOVED***
		for _, key := range fixedKeys ***REMOVED***
			var value interface***REMOVED******REMOVED***
			switch ***REMOVED***
			case key == f.FieldMap.resolve(FieldKeyTime):
				value = entry.Time.Format(timestampFormat)
			case key == f.FieldMap.resolve(FieldKeyLevel):
				value = entry.Level.String()
			case key == f.FieldMap.resolve(FieldKeyMsg):
				value = entry.Message
			case key == f.FieldMap.resolve(FieldKeyLogrusError):
				value = entry.err
			case key == f.FieldMap.resolve(FieldKeyFunc) && entry.HasCaller():
				value = entry.Caller.Function
			case key == f.FieldMap.resolve(FieldKeyFile) && entry.HasCaller():
				value = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
			default:
				value = entry.Data[key]
			***REMOVED***
			f.appendKeyValue(b, key, value)
		***REMOVED***
	***REMOVED***

	b.WriteByte('\n')
	return b.Bytes(), nil
***REMOVED***

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *Entry, keys []string, timestampFormat string) ***REMOVED***
	var levelColor int
	switch entry.Level ***REMOVED***
	case DebugLevel, TraceLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	***REMOVED***

	levelText := strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation ***REMOVED***
		levelText = levelText[0:4]
	***REMOVED***

	// Remove a single newline if it already exists in the message to keep
	// the behavior of logrus text_formatter the same as the stdlib log package
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	caller := ""

	if entry.HasCaller() ***REMOVED***
		caller = fmt.Sprintf("%s:%d %s()",
			entry.Caller.File, entry.Caller.Line, entry.Caller.Function)
	***REMOVED***

	if f.DisableTimestamp ***REMOVED***
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m%s %-44s ", levelColor, levelText, caller, entry.Message)
	***REMOVED*** else if !f.FullTimestamp ***REMOVED***
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d]%s %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), caller, entry.Message)
	***REMOVED*** else ***REMOVED***
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s]%s %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), caller, entry.Message)
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
