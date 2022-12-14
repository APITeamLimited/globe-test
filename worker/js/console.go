package js

import (
	"encoding/json"
	"strings"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
)

// console represents a JS console implemented as a logrus.Logger.
type console struct {
	logger logrus.FieldLogger
}

// Creates a console with the standard logrus logger.
func newConsole(logger logrus.FieldLogger) *console {
	return &console{logger.WithField("source", "console")}
}

func (c console) log(level logrus.Level, args ...goja.Value) {
	var strs strings.Builder
	for i := 0; i < len(args); i++ {
		if i > 0 {
			strs.WriteString(" ")
		}
		strs.WriteString(c.valueString(args[i]))
	}
	msg := strs.String()

	switch level { //nolint:exhaustive
	case logrus.DebugLevel:
		c.logger.Debug(msg)
	case logrus.InfoLevel:
		c.logger.Info(msg)
	case logrus.WarnLevel:
		c.logger.Warn(msg)
	case logrus.ErrorLevel:
		c.logger.Error(msg)
	}
}

func (c console) Log(args ...goja.Value) {
	c.Info(args...)
}

func (c console) Debug(args ...goja.Value) {
	c.log(logrus.DebugLevel, args...)
}

func (c console) Info(args ...goja.Value) {
	c.log(logrus.InfoLevel, args...)
}

func (c console) Warn(args ...goja.Value) {
	c.log(logrus.WarnLevel, args...)
}

func (c console) Error(args ...goja.Value) {
	c.log(logrus.ErrorLevel, args...)
}

func (c console) valueString(v goja.Value) string {
	mv, ok := v.(json.Marshaler)
	if !ok {
		return v.String()
	}

	b, err := json.Marshal(mv)
	if err != nil {
		return v.String()
	}
	return string(b)
}
