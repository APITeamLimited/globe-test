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
	GetMethodLogger(methodName string) MethodLogger
***REMOVED***

// binLogger is the global binary logger for the binary. One of this should be
// built at init time from the configuration (environment variable or flags).
//
// It is used to get a methodLogger for each individual method.
var binLogger Logger

var grpclogLogger = grpclog.Component("binarylog")

// SetLogger sets the binary logger.
//
// Only call this at init time.
func SetLogger(l Logger) ***REMOVED***
	binLogger = l
***REMOVED***

// GetLogger gets the binary logger.
//
// Only call this at init time.
func GetLogger() Logger ***REMOVED***
	return binLogger
***REMOVED***

// GetMethodLogger returns the methodLogger for the given methodName.
//
// methodName should be in the format of "/service/method".
//
// Each methodLogger returned by this method is a new instance. This is to
// generate sequence id within the call.
func GetMethodLogger(methodName string) MethodLogger ***REMOVED***
	if binLogger == nil ***REMOVED***
		return nil
	***REMOVED***
	return binLogger.GetMethodLogger(methodName)
***REMOVED***

func init() ***REMOVED***
	const envStr = "GRPC_BINARY_LOG_FILTER"
	configStr := os.Getenv(envStr)
	binLogger = NewLoggerFromConfigString(configStr)
***REMOVED***

// MethodLoggerConfig contains the setting for logging behavior of a method
// logger. Currently, it contains the max length of header and message.
type MethodLoggerConfig struct ***REMOVED***
	// Max length of header and message.
	Header, Message uint64
***REMOVED***

// LoggerConfig contains the config for loggers to create method loggers.
type LoggerConfig struct ***REMOVED***
	All      *MethodLoggerConfig
	Services map[string]*MethodLoggerConfig
	Methods  map[string]*MethodLoggerConfig

	Blacklist map[string]struct***REMOVED******REMOVED***
***REMOVED***

type logger struct ***REMOVED***
	config LoggerConfig
***REMOVED***

// NewLoggerFromConfig builds a logger with the given LoggerConfig.
func NewLoggerFromConfig(config LoggerConfig) Logger ***REMOVED***
	return &logger***REMOVED***config: config***REMOVED***
***REMOVED***

// newEmptyLogger creates an empty logger. The map fields need to be filled in
// using the set* functions.
func newEmptyLogger() *logger ***REMOVED***
	return &logger***REMOVED******REMOVED***
***REMOVED***

// Set method logger for "*".
func (l *logger) setDefaultMethodLogger(ml *MethodLoggerConfig) error ***REMOVED***
	if l.config.All != nil ***REMOVED***
		return fmt.Errorf("conflicting global rules found")
	***REMOVED***
	l.config.All = ml
	return nil
***REMOVED***

// Set method logger for "service/*".
//
// New methodLogger with same service overrides the old one.
func (l *logger) setServiceMethodLogger(service string, ml *MethodLoggerConfig) error ***REMOVED***
	if _, ok := l.config.Services[service]; ok ***REMOVED***
		return fmt.Errorf("conflicting service rules for service %v found", service)
	***REMOVED***
	if l.config.Services == nil ***REMOVED***
		l.config.Services = make(map[string]*MethodLoggerConfig)
	***REMOVED***
	l.config.Services[service] = ml
	return nil
***REMOVED***

// Set method logger for "service/method".
//
// New methodLogger with same method overrides the old one.
func (l *logger) setMethodMethodLogger(method string, ml *MethodLoggerConfig) error ***REMOVED***
	if _, ok := l.config.Blacklist[method]; ok ***REMOVED***
		return fmt.Errorf("conflicting blacklist rules for method %v found", method)
	***REMOVED***
	if _, ok := l.config.Methods[method]; ok ***REMOVED***
		return fmt.Errorf("conflicting method rules for method %v found", method)
	***REMOVED***
	if l.config.Methods == nil ***REMOVED***
		l.config.Methods = make(map[string]*MethodLoggerConfig)
	***REMOVED***
	l.config.Methods[method] = ml
	return nil
***REMOVED***

// Set blacklist method for "-service/method".
func (l *logger) setBlacklist(method string) error ***REMOVED***
	if _, ok := l.config.Blacklist[method]; ok ***REMOVED***
		return fmt.Errorf("conflicting blacklist rules for method %v found", method)
	***REMOVED***
	if _, ok := l.config.Methods[method]; ok ***REMOVED***
		return fmt.Errorf("conflicting method rules for method %v found", method)
	***REMOVED***
	if l.config.Blacklist == nil ***REMOVED***
		l.config.Blacklist = make(map[string]struct***REMOVED******REMOVED***)
	***REMOVED***
	l.config.Blacklist[method] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return nil
***REMOVED***

// getMethodLogger returns the methodLogger for the given methodName.
//
// methodName should be in the format of "/service/method".
//
// Each methodLogger returned by this method is a new instance. This is to
// generate sequence id within the call.
func (l *logger) GetMethodLogger(methodName string) MethodLogger ***REMOVED***
	s, m, err := grpcutil.ParseMethod(methodName)
	if err != nil ***REMOVED***
		grpclogLogger.Infof("binarylogging: failed to parse %q: %v", methodName, err)
		return nil
	***REMOVED***
	if ml, ok := l.config.Methods[s+"/"+m]; ok ***REMOVED***
		return newMethodLogger(ml.Header, ml.Message)
	***REMOVED***
	if _, ok := l.config.Blacklist[s+"/"+m]; ok ***REMOVED***
		return nil
	***REMOVED***
	if ml, ok := l.config.Services[s]; ok ***REMOVED***
		return newMethodLogger(ml.Header, ml.Message)
	***REMOVED***
	if l.config.All == nil ***REMOVED***
		return nil
	***REMOVED***
	return newMethodLogger(l.config.All.Header, l.config.All.Message)
***REMOVED***
