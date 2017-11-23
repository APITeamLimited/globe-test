package logrus

import (
	"encoding/json"
	"fmt"
)

type fieldKey string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[fieldKey]string

// Default key names for the default fields
const (
	FieldKeyMsg   = "msg"
	FieldKeyLevel = "level"
	FieldKeyTime  = "time"
)

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

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter***REMOVED***
	//   	FieldMap: FieldMap***REMOVED***
	// 		 FieldKeyTime: "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg: "@message",
	//    ***REMOVED***,
	// ***REMOVED***
	FieldMap FieldMap
***REMOVED***

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) ***REMOVED***
	data := make(Fields, len(entry.Data)+3)
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
	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" ***REMOVED***
		timestampFormat = defaultTimestampFormat
	***REMOVED***

	if !f.DisableTimestamp ***REMOVED***
		data[f.FieldMap.resolve(FieldKeyTime)] = entry.Time.Format(timestampFormat)
	***REMOVED***
	data[f.FieldMap.resolve(FieldKeyMsg)] = entry.Message
	data[f.FieldMap.resolve(FieldKeyLevel)] = entry.Level.String()

	serialized, err := json.Marshal(data)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	***REMOVED***
	return append(serialized, '\n'), nil
***REMOVED***
