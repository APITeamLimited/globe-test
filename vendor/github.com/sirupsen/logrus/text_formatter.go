package logrus

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

var baseTimestamp time.Time

func init() ***REMOVED***
	baseTimestamp = time.Now()
***REMOVED***

// TextFormatter formats logs into text
type TextFormatter struct ***REMOVED***
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Force quoting of all values
	ForceQuote bool

	// DisableQuote disables quoting for all values.
	// DisableQuote will have a lower priority than ForceQuote.
	// If both of them are set to true, quote will be forced on all values.
	DisableQuote bool

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// PadLevelText Adds padding the level text so that all the levels output at the same length
	// PadLevelText is a superset of the DisableLevelTruncation option
	PadLevelText bool

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

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	terminalInitOnce sync.Once

	// The max length of the level text, generated dynamically on init
	levelTextMaxLength int
***REMOVED***

func (f *TextFormatter) init(entry *Entry) ***REMOVED***
	if entry.Logger != nil ***REMOVED***
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	***REMOVED***
	// Get the max length of the level text
	for _, level := range AllLevels ***REMOVED***
		levelTextLength := utf8.RuneCount([]byte(level.String()))
		if levelTextLength > f.levelTextMaxLength ***REMOVED***
			f.levelTextMaxLength = levelTextLength
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *TextFormatter) isColored() bool ***REMOVED***
	isColored := f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows"))

	if f.EnvironmentOverrideColors ***REMOVED***
		switch force, ok := os.LookupEnv("CLICOLOR_FORCE"); ***REMOVED***
		case ok && force != "0":
			isColored = true
		case ok && force == "0", os.Getenv("CLICOLOR") == "0":
			isColored = false
		***REMOVED***
	***REMOVED***

	return isColored && !f.DisableColors
***REMOVED***

// Format renders a single log entry
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) ***REMOVED***
	data := make(Fields)
	for k, v := range entry.Data ***REMOVED***
		data[k] = v
	***REMOVED***
	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())
	keys := make([]string, 0, len(data))
	for k := range data ***REMOVED***
		keys = append(keys, k)
	***REMOVED***

	var funcVal, fileVal string

	fixedKeys := make([]string, 0, 4+len(data))
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
		if f.CallerPrettyfier != nil ***REMOVED***
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		***REMOVED*** else ***REMOVED***
			funcVal = entry.Caller.Function
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		***REMOVED***

		if funcVal != "" ***REMOVED***
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyFunc))
		***REMOVED***
		if fileVal != "" ***REMOVED***
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyFile))
		***REMOVED***
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
		f.printColored(b, entry, keys, data, timestampFormat)
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
				value = funcVal
			case key == f.FieldMap.resolve(FieldKeyFile) && entry.HasCaller():
				value = fileVal
			default:
				value = data[key]
			***REMOVED***
			f.appendKeyValue(b, key, value)
		***REMOVED***
	***REMOVED***

	b.WriteByte('\n')
	return b.Bytes(), nil
***REMOVED***

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *Entry, keys []string, data Fields, timestampFormat string) ***REMOVED***
	var levelColor int
	switch entry.Level ***REMOVED***
	case DebugLevel, TraceLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	case InfoLevel:
		levelColor = blue
	default:
		levelColor = blue
	***REMOVED***

	levelText := strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation && !f.PadLevelText ***REMOVED***
		levelText = levelText[0:4]
	***REMOVED***
	if f.PadLevelText ***REMOVED***
		// Generates the format string used in the next line, for example "%-6s" or "%-7s".
		// Based on the max level text length.
		formatString := "%-" + strconv.Itoa(f.levelTextMaxLength) + "s"
		// Formats the level text by appending spaces up to the max length, for example:
		// 	- "INFO   "
		//	- "WARNING"
		levelText = fmt.Sprintf(formatString, levelText)
	***REMOVED***

	// Remove a single newline if it already exists in the message to keep
	// the behavior of logrus text_formatter the same as the stdlib log package
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	caller := ""
	if entry.HasCaller() ***REMOVED***
		funcVal := fmt.Sprintf("%s()", entry.Caller.Function)
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)

		if f.CallerPrettyfier != nil ***REMOVED***
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		***REMOVED***

		if fileVal == "" ***REMOVED***
			caller = funcVal
		***REMOVED*** else if funcVal == "" ***REMOVED***
			caller = fileVal
		***REMOVED*** else ***REMOVED***
			caller = fileVal + " " + funcVal
		***REMOVED***
	***REMOVED***

	switch ***REMOVED***
	case f.DisableTimestamp:
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m%s %-44s ", levelColor, levelText, caller, entry.Message)
	case !f.FullTimestamp:
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d]%s %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), caller, entry.Message)
	default:
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s]%s %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), caller, entry.Message)
	***REMOVED***
	for _, k := range keys ***REMOVED***
		v := data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	***REMOVED***
***REMOVED***

func (f *TextFormatter) needsQuoting(text string) bool ***REMOVED***
	if f.ForceQuote ***REMOVED***
		return true
	***REMOVED***
	if f.QuoteEmptyFields && len(text) == 0 ***REMOVED***
		return true
	***REMOVED***
	if f.DisableQuote ***REMOVED***
		return false
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
