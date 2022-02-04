/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package js

import (
	"encoding/json"
	"os"
	"reflect"
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
	exptype := v.ExportType()
	if exptype == nil ***REMOVED***
		return v.String()
	***REMOVED***

	kind := exptype.Kind()
	if kind != reflect.Map && kind != reflect.Slice ***REMOVED***
		return v.String()
	***REMOVED***

	b, err := json.Marshal(v)
	if err != nil ***REMOVED***
		return v.String()
	***REMOVED***
	return string(b)
***REMOVED***
