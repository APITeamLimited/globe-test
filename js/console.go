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
	"context"
	"strconv"

	"github.com/dop251/goja"
	log "github.com/Sirupsen/logrus"
)

type Console struct ***REMOVED***
	Logger *log.Logger
***REMOVED***

func NewConsole() *Console ***REMOVED***
	return &Console***REMOVED***log.StandardLogger()***REMOVED***
***REMOVED***

func (c Console) log(ctx *context.Context, level log.Level, msgobj goja.Value, args ...goja.Value) ***REMOVED***
	if ctx != nil && *ctx != nil ***REMOVED***
		select ***REMOVED***
		case <-(*ctx).Done():
			return
		default:
		***REMOVED***
	***REMOVED***

	fields := make(log.Fields)
	for i, arg := range args ***REMOVED***
		fields[strconv.Itoa(i)] = arg.String()
	***REMOVED***
	msg := msgobj.ToString()
	e := c.Logger.WithFields(fields)
	switch level ***REMOVED***
	case log.DebugLevel:
		e.Debug(msg)
	case log.InfoLevel:
		e.Info(msg)
	case log.WarnLevel:
		e.Warn(msg)
	case log.ErrorLevel:
		e.Error(msg)
	***REMOVED***
***REMOVED***

func (c Console) Log(ctx *context.Context, msg goja.Value, args ...goja.Value) ***REMOVED***
	c.Info(ctx, msg, args...)
***REMOVED***

func (c Console) Debug(ctx *context.Context, msg goja.Value, args ...goja.Value) ***REMOVED***
	c.log(ctx, log.DebugLevel, msg, args...)
***REMOVED***

func (c Console) Info(ctx *context.Context, msg goja.Value, args ...goja.Value) ***REMOVED***
	c.log(ctx, log.InfoLevel, msg, args...)
***REMOVED***

func (c Console) Warn(ctx *context.Context, msg goja.Value, args ...goja.Value) ***REMOVED***
	c.log(ctx, log.WarnLevel, msg, args...)
***REMOVED***

func (c Console) Error(ctx *context.Context, msg goja.Value, args ...goja.Value) ***REMOVED***
	c.log(ctx, log.ErrorLevel, msg, args...)
***REMOVED***
