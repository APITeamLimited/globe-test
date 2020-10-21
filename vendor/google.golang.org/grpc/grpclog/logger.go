/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpclog

import "google.golang.org/grpc/internal/grpclog"

// Logger mimics golang's standard Logger as an interface.
//
// Deprecated: use LoggerV2.
type Logger interface ***REMOVED***
	Fatal(args ...interface***REMOVED******REMOVED***)
	Fatalf(format string, args ...interface***REMOVED******REMOVED***)
	Fatalln(args ...interface***REMOVED******REMOVED***)
	Print(args ...interface***REMOVED******REMOVED***)
	Printf(format string, args ...interface***REMOVED******REMOVED***)
	Println(args ...interface***REMOVED******REMOVED***)
***REMOVED***

// SetLogger sets the logger that is used in grpc. Call only from
// init() functions.
//
// Deprecated: use SetLoggerV2.
func SetLogger(l Logger) ***REMOVED***
	grpclog.Logger = &loggerWrapper***REMOVED***Logger: l***REMOVED***
***REMOVED***

// loggerWrapper wraps Logger into a LoggerV2.
type loggerWrapper struct ***REMOVED***
	Logger
***REMOVED***

func (g *loggerWrapper) Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.Logger.Print(args...)
***REMOVED***

func (g *loggerWrapper) Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.Logger.Println(args...)
***REMOVED***

func (g *loggerWrapper) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.Logger.Printf(format, args...)
***REMOVED***

func (g *loggerWrapper) Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.Logger.Print(args...)
***REMOVED***

func (g *loggerWrapper) Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.Logger.Println(args...)
***REMOVED***

func (g *loggerWrapper) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.Logger.Printf(format, args...)
***REMOVED***

func (g *loggerWrapper) Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.Logger.Print(args...)
***REMOVED***

func (g *loggerWrapper) Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.Logger.Println(args...)
***REMOVED***

func (g *loggerWrapper) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.Logger.Printf(format, args...)
***REMOVED***

func (g *loggerWrapper) V(l int) bool ***REMOVED***
	// Returns true for all verbose level.
	return true
***REMOVED***
