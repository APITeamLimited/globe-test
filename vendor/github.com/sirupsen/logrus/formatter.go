package logrus

import "time"

// Default key names for the default fields
const (
	defaultTimestampFormat = time.RFC3339
	FieldKeyMsg            = "msg"
	FieldKeyLevel          = "level"
	FieldKeyTime           = "time"
	FieldKeyLogrusError    = "logrus_error"
	FieldKeyFunc           = "func"
	FieldKeyFile           = "file"
)

// The Formatter interface is used to implement a custom Formatter. It takes an
// `Entry`. It exposes all the fields, including the default ones:
//
// * `entry.Data["msg"]`. The message passed from Info, Warn, Error ..
// * `entry.Data["time"]`. The timestamp.
// * `entry.Data["level"]. The level the entry was logged at.
//
// Any additional fields added with `WithField` or `WithFields` are also in
// `entry.Data`. Format is expected to return an array of bytes which are then
// logged to `logger.Out`.
type Formatter interface ***REMOVED***
	Format(*Entry) ([]byte, error)
***REMOVED***

// This is to not silently overwrite `time`, `msg`, `func` and `level` fields when
// dumping it. If this code wasn't there doing:
//
//  logrus.WithField("level", 1).Info("hello")
//
// Would just silently drop the user provided level. Instead with this code
// it'll logged as:
//
//  ***REMOVED***"level": "info", "fields.level": 1, "msg": "hello", "time": "..."***REMOVED***
//
// It's not exported because it's still using Data in an opinionated way. It's to
// avoid code duplication between the two default formatters.
func prefixFieldClashes(data Fields, fieldMap FieldMap, reportCaller bool) ***REMOVED***
	timeKey := fieldMap.resolve(FieldKeyTime)
	if t, ok := data[timeKey]; ok ***REMOVED***
		data["fields."+timeKey] = t
		delete(data, timeKey)
	***REMOVED***

	msgKey := fieldMap.resolve(FieldKeyMsg)
	if m, ok := data[msgKey]; ok ***REMOVED***
		data["fields."+msgKey] = m
		delete(data, msgKey)
	***REMOVED***

	levelKey := fieldMap.resolve(FieldKeyLevel)
	if l, ok := data[levelKey]; ok ***REMOVED***
		data["fields."+levelKey] = l
		delete(data, levelKey)
	***REMOVED***

	logrusErrKey := fieldMap.resolve(FieldKeyLogrusError)
	if l, ok := data[logrusErrKey]; ok ***REMOVED***
		data["fields."+logrusErrKey] = l
		delete(data, logrusErrKey)
	***REMOVED***

	// If reportCaller is not set, 'func' will not conflict.
	if reportCaller ***REMOVED***
		funcKey := fieldMap.resolve(FieldKeyFunc)
		if l, ok := data[funcKey]; ok ***REMOVED***
			data["fields."+funcKey] = l
		***REMOVED***
		fileKey := fieldMap.resolve(FieldKeyFile)
		if l, ok := data[fileKey]; ok ***REMOVED***
			data["fields."+fileKey] = l
		***REMOVED***
	***REMOVED***
***REMOVED***
