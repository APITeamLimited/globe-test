/*
 *
 * Copyright 2018 gRPC authors.
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

// Package binarylog implementation binary logging as defined in
// https://github.com/grpc/proposal/blob/master/A16-binary-logging.md.
package binarylog

import (
	"fmt"
	"os"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/grpcutil"
)

// Logger is the global binary logger. It can be used to get binary logger for
// each method.
type Logger interface ***REMOVED***
	getMethodLogger(methodName string) *MethodLogger
***REMOVED***

// binLogger is the global binary logger for the binary. One of this should be
// built at init time from the configuration (environment variable or flags).
//
// It is used to get a methodLogger for each individual method.
var binLogger Logger

var grpclogLogger = grpclog.Component("binarylog")

// SetLogger sets the binarg logger.
//
// Only call this at init time.
func SetLogger(l Logger) ***REMOVED***
	binLogger = l
***REMOVED***

// GetMethodLogger returns the methodLogger for the given methodName.
//
// methodName should be in the format of "/service/method".
//
// Each methodLogger returned by this method is a new instance. This is to
// generate sequence id within the call.
func GetMethodLogger(methodName string) *MethodLogger ***REMOVED***
	if binLogger == nil ***REMOVED***
		return nil
	***REMOVED***
	return binLogger.getMethodLogger(methodName)
***REMOVED***

func init() ***REMOVED***
	const envStr = "GRPC_BINARY_LOG_FILTER"
	configStr := os.Getenv(envStr)
	binLogger = NewLoggerFromConfigString(configStr)
***REMOVED***

type methodLoggerConfig struct ***REMOVED***
	// Max length of header and message.
	hdr, msg uint64
***REMOVED***

type logger struct ***REMOVED***
	all      *methodLoggerConfig
	services map[string]*methodLoggerConfig
	methods  map[string]*methodLoggerConfig

	blacklist map[string]struct***REMOVED******REMOVED***
***REMOVED***

// newEmptyLogger creates an empty logger. The map fields need to be filled in
// using the set* functions.
func newEmptyLogger() *logger ***REMOVED***
	return &logger***REMOVED******REMOVED***
***REMOVED***

// Set method logger for "*".
func (l *logger) setDefaultMethodLogger(ml *methodLoggerConfig) error ***REMOVED***
	if l.all != nil ***REMOVED***
		return fmt.Errorf("conflicting global rules found")
	***REMOVED***
	l.all = ml
	return nil
***REMOVED***

// Set method logger for "service/*".
//
// New methodLogger with same service overrides the old one.
func (l *logger) setServiceMethodLogger(service string, ml *methodLoggerConfig) error ***REMOVED***
	if _, ok := l.services[service]; ok ***REMOVED***
		return fmt.Errorf("conflicting service rules for service %v found", service)
	***REMOVED***
	if l.services == nil ***REMOVED***
		l.services = make(map[string]*methodLoggerConfig)
	***REMOVED***
	l.services[service] = ml
	return nil
***REMOVED***

// Set method logger for "service/method".
//
// New methodLogger with same method overrides the old one.
func (l *logger) setMethodMethodLogger(method string, ml *methodLoggerConfig) error ***REMOVED***
	if _, ok := l.blacklist[method]; ok ***REMOVED***
		return fmt.Errorf("conflicting blacklist rules for method %v found", method)
	***REMOVED***
	if _, ok := l.methods[method]; ok ***REMOVED***
		return fmt.Errorf("conflicting method rules for method %v found", method)
	***REMOVED***
	if l.methods == nil ***REMOVED***
		l.methods = make(map[string]*methodLoggerConfig)
	***REMOVED***
	l.methods[method] = ml
	return nil
***REMOVED***

// Set blacklist method for "-service/method".
func (l *logger) setBlacklist(method string) error ***REMOVED***
	if _, ok := l.blacklist[method]; ok ***REMOVED***
		return fmt.Errorf("conflicting blacklist rules for method %v found", method)
	***REMOVED***
	if _, ok := l.methods[method]; ok ***REMOVED***
		return fmt.Errorf("conflicting method rules for method %v found", method)
	***REMOVED***
	if l.blacklist == nil ***REMOVED***
		l.blacklist = make(map[string]struct***REMOVED******REMOVED***)
	***REMOVED***
	l.blacklist[method] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return nil
***REMOVED***

// getMethodLogger returns the methodLogger for the given methodName.
//
// methodName should be in the format of "/service/method".
//
// Each methodLogger returned by this method is a new instance. This is to
// generate sequence id within the call.
func (l *logger) getMethodLogger(methodName string) *MethodLogger ***REMOVED***
	s, m, err := grpcutil.ParseMethod(methodName)
	if err != nil ***REMOVED***
		grpclogLogger.Infof("binarylogging: failed to parse %q: %v", methodName, err)
		return nil
	***REMOVED***
	if ml, ok := l.methods[s+"/"+m]; ok ***REMOVED***
		return newMethodLogger(ml.hdr, ml.msg)
	***REMOVED***
	if _, ok := l.blacklist[s+"/"+m]; ok ***REMOVED***
		return nil
	***REMOVED***
	if ml, ok := l.services[s]; ok ***REMOVED***
		return newMethodLogger(ml.hdr, ml.msg)
	***REMOVED***
	if l.all == nil ***REMOVED***
		return nil
	***REMOVED***
	return newMethodLogger(l.all.hdr, l.all.msg)
***REMOVED***
