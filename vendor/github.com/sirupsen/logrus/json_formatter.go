package logrus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
)

type fieldKey string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[fieldKey]string

func (f FieldMap) resolve(key fieldKey) string ***REMOVED***
	if k, ok := f[key]; ok ***REMOVED***
		return k
	***REMOVED***

	return string(key)
***REMOVED***

// JSONFormatter formats logs into parsable json
type JSONFormatter struct ***REMOVED***
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DisableHTMLEscape allows disabling html escaping in output
	DisableHTMLEscape bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter***REMOVED***
	//   	FieldMap: FieldMap***REMOVED***
	// 		 FieldKeyTime:  "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg:   "@message",
	// 		 FieldKeyFunc:  "@caller",
	//    ***REMOVED***,
	// ***REMOVED***
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// PrettyPrint will indent all json logs
	PrettyPrint bool
***REMOVED***

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) ***REMOVED***
	data := make(Fields, len(entry.Data)+4)
	for k, v := range entry.Data ***REMOVED***
		switch v := v.(type) ***REMOVED***
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		***REMOVED***
	***REMOVED***

	if f.DataKey != "" ***REMOVED***
		newData := make(Fields, 4)
		newData[f.DataKey] = data
		data = newData
	***REMOVED***

	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" ***REMOVED***
		timestampFormat = defaultTimestampFormat
	***REMOVED***

	if entry.err != "" ***REMOVED***
		data[f.FieldMap.resolve(FieldKeyLogrusError)] = entry.err
	***REMOVED***
	if !f.DisableTimestamp ***REMOVED***
		data[f.FieldMap.resolve(FieldKeyTime)] = entry.Time.Format(timestampFormat)
	***REMOVED***
	data[f.FieldMap.resolve(FieldKeyMsg)] = entry.Message
	data[f.FieldMap.resolve(FieldKeyLevel)] = entry.Level.String()
	if entry.HasCaller() ***REMOVED***
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if f.CallerPrettyfier != nil ***REMOVED***
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		***REMOVED***
		if funcVal != "" ***REMOVED***
			data[f.FieldMap.resolve(FieldKeyFunc)] = funcVal
		***REMOVED***
		if fileVal != "" ***REMOVED***
			data[f.FieldMap.resolve(FieldKeyFile)] = fileVal
		***REMOVED***
	***REMOVED***

	var b *bytes.Buffer
	if entry.Buffer != nil ***REMOVED***
		b = entry.Buffer
	***REMOVED*** else ***REMOVED***
		b = &bytes.Buffer***REMOVED******REMOVED***
	***REMOVED***

	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!f.DisableHTMLEscape)
	if f.PrettyPrint ***REMOVED***
		encoder.SetIndent("", "  ")
	***REMOVED***
	if err := encoder.Encode(data); err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	***REMOVED***

	return b.Bytes(), nil
***REMOVED***
