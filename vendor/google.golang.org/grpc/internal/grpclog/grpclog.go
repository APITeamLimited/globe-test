/*
 *
 * Copyright 2020 gRPC authors.
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

// Package grpclog (internal) defines depth logging for grpc.
package grpclog

import (
	"os"
)

// Logger is the logger used for the non-depth log functions.
var Logger LoggerV2

// DepthLogger is the logger used for the depth log functions.
var DepthLogger DepthLoggerV2

// InfoDepth logs to the INFO log at the specified depth.
func InfoDepth(depth int, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if DepthLogger != nil ***REMOVED***
		DepthLogger.InfoDepth(depth, args...)
	***REMOVED*** else ***REMOVED***
		Logger.Infoln(args...)
	***REMOVED***
***REMOVED***

// WarningDepth logs to the WARNING log at the specified depth.
func WarningDepth(depth int, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if DepthLogger != nil ***REMOVED***
		DepthLogger.WarningDepth(depth, args...)
	***REMOVED*** else ***REMOVED***
		Logger.Warningln(args...)
	***REMOVED***
***REMOVED***

// ErrorDepth logs to the ERROR log at the specified depth.
func ErrorDepth(depth int, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if DepthLogger != nil ***REMOVED***
		DepthLogger.ErrorDepth(depth, args...)
	***REMOVED*** else ***REMOVED***
		Logger.Errorln(args...)
	***REMOVED***
***REMOVED***

// FatalDepth logs to the FATAL log at the specified depth.
func FatalDepth(depth int, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if DepthLogger != nil ***REMOVED***
		DepthLogger.FatalDepth(depth, args...)
	***REMOVED*** else ***REMOVED***
		Logger.Fatalln(args...)
	***REMOVED***
	os.Exit(1)
***REMOVED***

// LoggerV2 does underlying logging work for grpclog.
// This is a copy of the LoggerV2 defined in the external grpclog package. It
// is defined here to avoid a circular dependency.
type LoggerV2 interface ***REMOVED***
	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(args ...interface***REMOVED******REMOVED***)
	// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
	Infoln(args ...interface***REMOVED******REMOVED***)
	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface***REMOVED******REMOVED***)
	// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
	Warning(args ...interface***REMOVED******REMOVED***)
	// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
	Warningln(args ...interface***REMOVED******REMOVED***)
	// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
	Warningf(format string, args ...interface***REMOVED******REMOVED***)
	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(args ...interface***REMOVED******REMOVED***)
	// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	Errorln(args ...interface***REMOVED******REMOVED***)
	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, args ...interface***REMOVED******REMOVED***)
	// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatal(args ...interface***REMOVED******REMOVED***)
	// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalln(args ...interface***REMOVED******REMOVED***)
	// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalf(format string, args ...interface***REMOVED******REMOVED***)
	// V reports whether verbosity level l is at least the requested verbose level.
	V(l int) bool
***REMOVED***

// DepthLoggerV2 logs at a specified call frame. If a LoggerV2 also implements
// DepthLoggerV2, the below functions will be called with the appropriate stack
// depth set for trivial functions the logger may ignore.
// This is a copy of the DepthLoggerV2 defined in the external grpclog package.
// It is defined here to avoid a circular dependency.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type DepthLoggerV2 interface ***REMOVED***
	// InfoDepth logs to INFO log at the specified depth. Arguments are handled in the manner of fmt.Print.
	InfoDepth(depth int, args ...interface***REMOVED******REMOVED***)
	// WarningDepth logs to WARNING log at the specified depth. Arguments are handled in the manner of fmt.Print.
	WarningDepth(depth int, args ...interface***REMOVED******REMOVED***)
	// ErrorDetph logs to ERROR log at the specified depth. Arguments are handled in the manner of fmt.Print.
	ErrorDepth(depth int, args ...interface***REMOVED******REMOVED***)
	// FatalDepth logs to FATAL log at the specified depth. Arguments are handled in the manner of fmt.Print.
	FatalDepth(depth int, args ...interface***REMOVED******REMOVED***)
***REMOVED***
