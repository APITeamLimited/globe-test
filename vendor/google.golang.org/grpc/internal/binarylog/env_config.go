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

package binarylog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// NewLoggerFromConfigString reads the string and build a logger. It can be used
// to build a new logger and assign it to binarylog.Logger.
//
// Example filter config strings:
//  - "" Nothing will be logged
//  - "*" All headers and messages will be fully logged.
//  - "****REMOVED***h***REMOVED***" Only headers will be logged.
//  - "****REMOVED***m:256***REMOVED***" Only the first 256 bytes of each message will be logged.
//  - "Foo/*" Logs every method in service Foo
//  - "Foo/*,-Foo/Bar" Logs every method in service Foo except method /Foo/Bar
//  - "Foo/*,Foo/Bar***REMOVED***m:256***REMOVED***" Logs the first 256 bytes of each message in method
//    /Foo/Bar, logs all headers and messages in every other method in service
//    Foo.
//
// If two configs exist for one certain method or service, the one specified
// later overrides the previous config.
func NewLoggerFromConfigString(s string) Logger ***REMOVED***
	if s == "" ***REMOVED***
		return nil
	***REMOVED***
	l := newEmptyLogger()
	methods := strings.Split(s, ",")
	for _, method := range methods ***REMOVED***
		if err := l.fillMethodLoggerWithConfigString(method); err != nil ***REMOVED***
			grpclogLogger.Warningf("failed to parse binary log config: %v", err)
			return nil
		***REMOVED***
	***REMOVED***
	return l
***REMOVED***

// fillMethodLoggerWithConfigString parses config, creates methodLogger and adds
// it to the right map in the logger.
func (l *logger) fillMethodLoggerWithConfigString(config string) error ***REMOVED***
	// "" is invalid.
	if config == "" ***REMOVED***
		return errors.New("empty string is not a valid method binary logging config")
	***REMOVED***

	// "-service/method", blacklist, no * or ***REMOVED******REMOVED*** allowed.
	if config[0] == '-' ***REMOVED***
		s, m, suffix, err := parseMethodConfigAndSuffix(config[1:])
		if err != nil ***REMOVED***
			return fmt.Errorf("invalid config: %q, %v", config, err)
		***REMOVED***
		if m == "*" ***REMOVED***
			return fmt.Errorf("invalid config: %q, %v", config, "* not allowed in blacklist config")
		***REMOVED***
		if suffix != "" ***REMOVED***
			return fmt.Errorf("invalid config: %q, %v", config, "header/message limit not allowed in blacklist config")
		***REMOVED***
		if err := l.setBlacklist(s + "/" + m); err != nil ***REMOVED***
			return fmt.Errorf("invalid config: %v", err)
		***REMOVED***
		return nil
	***REMOVED***

	// "****REMOVED***h:256;m:256***REMOVED***"
	if config[0] == '*' ***REMOVED***
		hdr, msg, err := parseHeaderMessageLengthConfig(config[1:])
		if err != nil ***REMOVED***
			return fmt.Errorf("invalid config: %q, %v", config, err)
		***REMOVED***
		if err := l.setDefaultMethodLogger(&methodLoggerConfig***REMOVED***hdr: hdr, msg: msg***REMOVED***); err != nil ***REMOVED***
			return fmt.Errorf("invalid config: %v", err)
		***REMOVED***
		return nil
	***REMOVED***

	s, m, suffix, err := parseMethodConfigAndSuffix(config)
	if err != nil ***REMOVED***
		return fmt.Errorf("invalid config: %q, %v", config, err)
	***REMOVED***
	hdr, msg, err := parseHeaderMessageLengthConfig(suffix)
	if err != nil ***REMOVED***
		return fmt.Errorf("invalid header/message length config: %q, %v", suffix, err)
	***REMOVED***
	if m == "*" ***REMOVED***
		if err := l.setServiceMethodLogger(s, &methodLoggerConfig***REMOVED***hdr: hdr, msg: msg***REMOVED***); err != nil ***REMOVED***
			return fmt.Errorf("invalid config: %v", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := l.setMethodMethodLogger(s+"/"+m, &methodLoggerConfig***REMOVED***hdr: hdr, msg: msg***REMOVED***); err != nil ***REMOVED***
			return fmt.Errorf("invalid config: %v", err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

const (
	// TODO: this const is only used by env_config now. But could be useful for
	// other config. Move to binarylog.go if necessary.
	maxUInt = ^uint64(0)

	// For "p.s/m" plus any suffix. Suffix will be parsed again. See test for
	// expected output.
	longMethodConfigRegexpStr = `^([\w./]+)/((?:\w+)|[*])(.+)?$`

	// For suffix from above, "***REMOVED***h:123,m:123***REMOVED***". See test for expected output.
	optionalLengthRegexpStr      = `(?::(\d+))?` // Optional ":123".
	headerConfigRegexpStr        = `^***REMOVED***h` + optionalLengthRegexpStr + `***REMOVED***$`
	messageConfigRegexpStr       = `^***REMOVED***m` + optionalLengthRegexpStr + `***REMOVED***$`
	headerMessageConfigRegexpStr = `^***REMOVED***h` + optionalLengthRegexpStr + `;m` + optionalLengthRegexpStr + `***REMOVED***$`
)

var (
	longMethodConfigRegexp    = regexp.MustCompile(longMethodConfigRegexpStr)
	headerConfigRegexp        = regexp.MustCompile(headerConfigRegexpStr)
	messageConfigRegexp       = regexp.MustCompile(messageConfigRegexpStr)
	headerMessageConfigRegexp = regexp.MustCompile(headerMessageConfigRegexpStr)
)

// Turn "service/method***REMOVED***h;m***REMOVED***" into "service", "method", "***REMOVED***h;m***REMOVED***".
func parseMethodConfigAndSuffix(c string) (service, method, suffix string, _ error) ***REMOVED***
	// Regexp result:
	//
	// in:  "p.s/m***REMOVED***h:123,m:123***REMOVED***",
	// out: []string***REMOVED***"p.s/m***REMOVED***h:123,m:123***REMOVED***", "p.s", "m", "***REMOVED***h:123,m:123***REMOVED***"***REMOVED***,
	match := longMethodConfigRegexp.FindStringSubmatch(c)
	if match == nil ***REMOVED***
		return "", "", "", fmt.Errorf("%q contains invalid substring", c)
	***REMOVED***
	service = match[1]
	method = match[2]
	suffix = match[3]
	return
***REMOVED***

// Turn "***REMOVED***h:123;m:345***REMOVED***" into 123, 345.
//
// Return maxUInt if length is unspecified.
func parseHeaderMessageLengthConfig(c string) (hdrLenStr, msgLenStr uint64, err error) ***REMOVED***
	if c == "" ***REMOVED***
		return maxUInt, maxUInt, nil
	***REMOVED***
	// Header config only.
	if match := headerConfigRegexp.FindStringSubmatch(c); match != nil ***REMOVED***
		if s := match[1]; s != "" ***REMOVED***
			hdrLenStr, err = strconv.ParseUint(s, 10, 64)
			if err != nil ***REMOVED***
				return 0, 0, fmt.Errorf("failed to convert %q to uint", s)
			***REMOVED***
			return hdrLenStr, 0, nil
		***REMOVED***
		return maxUInt, 0, nil
	***REMOVED***

	// Message config only.
	if match := messageConfigRegexp.FindStringSubmatch(c); match != nil ***REMOVED***
		if s := match[1]; s != "" ***REMOVED***
			msgLenStr, err = strconv.ParseUint(s, 10, 64)
			if err != nil ***REMOVED***
				return 0, 0, fmt.Errorf("failed to convert %q to uint", s)
			***REMOVED***
			return 0, msgLenStr, nil
		***REMOVED***
		return 0, maxUInt, nil
	***REMOVED***

	// Header and message config both.
	if match := headerMessageConfigRegexp.FindStringSubmatch(c); match != nil ***REMOVED***
		// Both hdr and msg are specified, but one or two of them might be empty.
		hdrLenStr = maxUInt
		msgLenStr = maxUInt
		if s := match[1]; s != "" ***REMOVED***
			hdrLenStr, err = strconv.ParseUint(s, 10, 64)
			if err != nil ***REMOVED***
				return 0, 0, fmt.Errorf("failed to convert %q to uint", s)
			***REMOVED***
		***REMOVED***
		if s := match[2]; s != "" ***REMOVED***
			msgLenStr, err = strconv.ParseUint(s, 10, 64)
			if err != nil ***REMOVED***
				return 0, 0, fmt.Errorf("failed to convert %q to uint", s)
			***REMOVED***
		***REMOVED***
		return hdrLenStr, msgLenStr, nil
	***REMOVED***
	return 0, 0, fmt.Errorf("%q contains invalid substring", c)
***REMOVED***
