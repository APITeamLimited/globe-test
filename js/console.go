package js

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
)

// console represents a JS console implemented as a logrus.Logger.
type console struct ***REMOVED***
	logger logrus.FieldLogger
***REMOVED***

// Creates a console with the standard logrus logger.
func newConsole(logger logrus.FieldLogger) *console ***REMOVED***
	return &console***REMOVED***logger.WithField("source", "console")***REMOVED***
***REMOVED***

// Creates a console logger with its output set to the file at the provided `filepath`.
func newFileConsole(filepath string, formatter logrus.Formatter) (*console, error) ***REMOVED***
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644) //nolint:gosec
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	l := logrus.New()
	l.SetOutput(f)
	l.SetFormatter(formatter)

	return &console***REMOVED***l***REMOVED***, nil
***REMOVED***

func (c console) log(level logrus.Level, args ...goja.Value) ***REMOVED***
	var strs strings.Builder
	for i := 0; i < len(args); i++ ***REMOVED***
		if i > 0 ***REMOVED***
			strs.WriteString(" ")
		***REMOVED***
		strs.WriteString(c.valueString(args[i]))
	***REMOVED***
	msg := strs.String()

	switch level ***REMOVED*** //nolint:exhaustive
	case logrus.DebugLevel:
		c.logger.Debug(msg)
	case logrus.InfoLevel:
		c.logger.Info(msg)
	case logrus.WarnLevel:
		c.logger.Warn(msg)
	case logrus.ErrorLevel:
		c.logger.Error(msg)
	***REMOVED***
***REMOVED***

func (c console) Log(args ...goja.Value) ***REMOVED***
	c.Info(args...)
***REMOVED***

func (c console) Debug(args ...goja.Value) ***REMOVED***
	c.log(logrus.DebugLevel, args...)
***REMOVED***

func (c console) Info(args ...goja.Value) ***REMOVED***
	c.log(logrus.InfoLevel, args...)
***REMOVED***

func (c console) Warn(args ...goja.Value) ***REMOVED***
	c.log(logrus.WarnLevel, args...)
***REMOVED***

func (c console) Error(args ...goja.Value) ***REMOVED***
	c.log(logrus.ErrorLevel, args...)
***REMOVED***

func (c console) valueString(v goja.Value) string ***REMOVED***
	mv, ok := v.(json.Marshaler)
	if !ok ***REMOVED***
		return v.String()
	***REMOVED***

	b, err := json.Marshal(mv)
	if err != nil ***REMOVED***
		return v.String()
	***REMOVED***
	return string(b)
***REMOVED***
