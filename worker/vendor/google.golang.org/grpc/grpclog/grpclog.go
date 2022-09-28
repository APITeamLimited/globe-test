/*
 *
 * Copyright 2017 gRPC authors.
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

// Package grpclog defines logging for grpc.
//
// All logs in transport and grpclb packages only go to verbose level 2.
// All logs in other packages in grpc are logged in spite of the verbosity level.
//
// In the default logger,
// severity level can be set by environment variable GRPC_GO_LOG_SEVERITY_LEVEL,
// verbosity level can be set by GRPC_GO_LOG_VERBOSITY_LEVEL.
package grpclog // import "google.golang.org/grpc/grpclog"

import (
	"os"

	"google.golang.org/grpc/internal/grpclog"
)

func init() ***REMOVED***
	SetLoggerV2(newLoggerV2())
***REMOVED***

// V reports whether verbosity level l is at least the requested verbose level.
func V(l int) bool ***REMOVED***
	return grpclog.Logger.V(l)
***REMOVED***

// Info logs to the INFO log.
func Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Info(args...)
***REMOVED***

// Infof logs to the INFO log. Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Infof(format, args...)
***REMOVED***

// Infoln logs to the INFO log. Arguments are handled in the manner of fmt.Println.
func Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Infoln(args...)
***REMOVED***

// Warning logs to the WARNING log.
func Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Warning(args...)
***REMOVED***

// Warningf logs to the WARNING log. Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Warningf(format, args...)
***REMOVED***

// Warningln logs to the WARNING log. Arguments are handled in the manner of fmt.Println.
func Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Warningln(args...)
***REMOVED***

// Error logs to the ERROR log.
func Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Error(args...)
***REMOVED***

// Errorf logs to the ERROR log. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Errorf(format, args...)
***REMOVED***

// Errorln logs to the ERROR log. Arguments are handled in the manner of fmt.Println.
func Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Errorln(args...)
***REMOVED***

// Fatal logs to the FATAL log. Arguments are handled in the manner of fmt.Print.
// It calls os.Exit() with exit code 1.
func Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Fatal(args...)
	// Make sure fatal logs will exit.
	os.Exit(1)
***REMOVED***

// Fatalf logs to the FATAL log. Arguments are handled in the manner of fmt.Printf.
// It calls os.Exit() with exit code 1.
func Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Fatalf(format, args...)
	// Make sure fatal logs will exit.
	os.Exit(1)
***REMOVED***

// Fatalln logs to the FATAL log. Arguments are handled in the manner of fmt.Println.
// It calle os.Exit()) with exit code 1.
func Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Fatalln(args...)
	// Make sure fatal logs will exit.
	os.Exit(1)
***REMOVED***

// Print prints to the logger. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: use Info.
func Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Info(args...)
***REMOVED***

// Printf prints to the logger. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: use Infof.
func Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Infof(format, args...)
***REMOVED***

// Println prints to the logger. Arguments are handled in the manner of fmt.Println.
//
// Deprecated: use Infoln.
func Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	grpclog.Logger.Infoln(args...)
***REMOVED***
