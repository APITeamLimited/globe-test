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

package grpclog

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc/internal/grpclog"
)

// LoggerV2 does underlying logging work for grpclog.
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

// SetLoggerV2 sets logger that is used in grpc to a V2 logger.
// Not mutex-protected, should be called before any gRPC functions.
func SetLoggerV2(l LoggerV2) ***REMOVED***
	if _, ok := l.(*componentData); ok ***REMOVED***
		panic("cannot use component logger as grpclog logger")
	***REMOVED***
	grpclog.Logger = l
	grpclog.DepthLogger, _ = l.(grpclog.DepthLoggerV2)
***REMOVED***

const (
	// infoLog indicates Info severity.
	infoLog int = iota
	// warningLog indicates Warning severity.
	warningLog
	// errorLog indicates Error severity.
	errorLog
	// fatalLog indicates Fatal severity.
	fatalLog
)

// severityName contains the string representation of each severity.
var severityName = []string***REMOVED***
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
***REMOVED***

// loggerT is the default logger used by grpclog.
type loggerT struct ***REMOVED***
	m          []*log.Logger
	v          int
	jsonFormat bool
***REMOVED***

// NewLoggerV2 creates a loggerV2 with the provided writers.
// Fatal logs will be written to errorW, warningW, infoW, followed by exit(1).
// Error logs will be written to errorW, warningW and infoW.
// Warning logs will be written to warningW and infoW.
// Info logs will be written to infoW.
func NewLoggerV2(infoW, warningW, errorW io.Writer) LoggerV2 ***REMOVED***
	return newLoggerV2WithConfig(infoW, warningW, errorW, loggerV2Config***REMOVED******REMOVED***)
***REMOVED***

// NewLoggerV2WithVerbosity creates a loggerV2 with the provided writers and
// verbosity level.
func NewLoggerV2WithVerbosity(infoW, warningW, errorW io.Writer, v int) LoggerV2 ***REMOVED***
	return newLoggerV2WithConfig(infoW, warningW, errorW, loggerV2Config***REMOVED***verbose: v***REMOVED***)
***REMOVED***

type loggerV2Config struct ***REMOVED***
	verbose    int
	jsonFormat bool
***REMOVED***

func newLoggerV2WithConfig(infoW, warningW, errorW io.Writer, c loggerV2Config) LoggerV2 ***REMOVED***
	var m []*log.Logger
	flag := log.LstdFlags
	if c.jsonFormat ***REMOVED***
		flag = 0
	***REMOVED***
	m = append(m, log.New(infoW, "", flag))
	m = append(m, log.New(io.MultiWriter(infoW, warningW), "", flag))
	ew := io.MultiWriter(infoW, warningW, errorW) // ew will be used for error and fatal.
	m = append(m, log.New(ew, "", flag))
	m = append(m, log.New(ew, "", flag))
	return &loggerT***REMOVED***m: m, v: c.verbose, jsonFormat: c.jsonFormat***REMOVED***
***REMOVED***

// newLoggerV2 creates a loggerV2 to be used as default logger.
// All logs are written to stderr.
func newLoggerV2() LoggerV2 ***REMOVED***
	errorW := ioutil.Discard
	warningW := ioutil.Discard
	infoW := ioutil.Discard

	logLevel := os.Getenv("GRPC_GO_LOG_SEVERITY_LEVEL")
	switch logLevel ***REMOVED***
	case "", "ERROR", "error": // If env is unset, set level to ERROR.
		errorW = os.Stderr
	case "WARNING", "warning":
		warningW = os.Stderr
	case "INFO", "info":
		infoW = os.Stderr
	***REMOVED***

	var v int
	vLevel := os.Getenv("GRPC_GO_LOG_VERBOSITY_LEVEL")
	if vl, err := strconv.Atoi(vLevel); err == nil ***REMOVED***
		v = vl
	***REMOVED***

	jsonFormat := strings.EqualFold(os.Getenv("GRPC_GO_LOG_FORMATTER"), "json")

	return newLoggerV2WithConfig(infoW, warningW, errorW, loggerV2Config***REMOVED***
		verbose:    v,
		jsonFormat: jsonFormat,
	***REMOVED***)
***REMOVED***

func (g *loggerT) output(severity int, s string) ***REMOVED***
	sevStr := severityName[severity]
	if !g.jsonFormat ***REMOVED***
		g.m[severity].Output(2, fmt.Sprintf("%v: %v", sevStr, s))
		return
	***REMOVED***
	// TODO: we can also include the logging component, but that needs more
	// (API) changes.
	b, _ := json.Marshal(map[string]string***REMOVED***
		"severity": sevStr,
		"message":  s,
	***REMOVED***)
	g.m[severity].Output(2, string(b))
***REMOVED***

func (g *loggerT) Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(infoLog, fmt.Sprint(args...))
***REMOVED***

func (g *loggerT) Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(infoLog, fmt.Sprintln(args...))
***REMOVED***

func (g *loggerT) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(infoLog, fmt.Sprintf(format, args...))
***REMOVED***

func (g *loggerT) Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(warningLog, fmt.Sprint(args...))
***REMOVED***

func (g *loggerT) Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(warningLog, fmt.Sprintln(args...))
***REMOVED***

func (g *loggerT) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(warningLog, fmt.Sprintf(format, args...))
***REMOVED***

func (g *loggerT) Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(errorLog, fmt.Sprint(args...))
***REMOVED***

func (g *loggerT) Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(errorLog, fmt.Sprintln(args...))
***REMOVED***

func (g *loggerT) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(errorLog, fmt.Sprintf(format, args...))
***REMOVED***

func (g *loggerT) Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(fatalLog, fmt.Sprint(args...))
	os.Exit(1)
***REMOVED***

func (g *loggerT) Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(fatalLog, fmt.Sprintln(args...))
	os.Exit(1)
***REMOVED***

func (g *loggerT) Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	g.output(fatalLog, fmt.Sprintf(format, args...))
	os.Exit(1)
***REMOVED***

func (g *loggerT) V(l int) bool ***REMOVED***
	return l <= g.v
***REMOVED***

// DepthLoggerV2 logs at a specified call frame. If a LoggerV2 also implements
// DepthLoggerV2, the below functions will be called with the appropriate stack
// depth set for trivial functions the logger may ignore.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type DepthLoggerV2 interface ***REMOVED***
	LoggerV2
	// InfoDepth logs to INFO log at the specified depth. Arguments are handled in the manner of fmt.Println.
	InfoDepth(depth int, args ...interface***REMOVED******REMOVED***)
	// WarningDepth logs to WARNING log at the specified depth. Arguments are handled in the manner of fmt.Println.
	WarningDepth(depth int, args ...interface***REMOVED******REMOVED***)
	// ErrorDepth logs to ERROR log at the specified depth. Arguments are handled in the manner of fmt.Println.
	ErrorDepth(depth int, args ...interface***REMOVED******REMOVED***)
	// FatalDepth logs to FATAL log at the specified depth. Arguments are handled in the manner of fmt.Println.
	FatalDepth(depth int, args ...interface***REMOVED******REMOVED***)
***REMOVED***
