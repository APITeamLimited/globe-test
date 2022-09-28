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

package grpclog

import (
	"fmt"
)

// PrefixLogger does logging with a prefix.
//
// Logging method on a nil logs without any prefix.
type PrefixLogger struct ***REMOVED***
	logger DepthLoggerV2
	prefix string
***REMOVED***

// Infof does info logging.
func (pl *PrefixLogger) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if pl != nil ***REMOVED***
		// Handle nil, so the tests can pass in a nil logger.
		format = pl.prefix + format
		pl.logger.InfoDepth(1, fmt.Sprintf(format, args...))
		return
	***REMOVED***
	InfoDepth(1, fmt.Sprintf(format, args...))
***REMOVED***

// Warningf does warning logging.
func (pl *PrefixLogger) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if pl != nil ***REMOVED***
		format = pl.prefix + format
		pl.logger.WarningDepth(1, fmt.Sprintf(format, args...))
		return
	***REMOVED***
	WarningDepth(1, fmt.Sprintf(format, args...))
***REMOVED***

// Errorf does error logging.
func (pl *PrefixLogger) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if pl != nil ***REMOVED***
		format = pl.prefix + format
		pl.logger.ErrorDepth(1, fmt.Sprintf(format, args...))
		return
	***REMOVED***
	ErrorDepth(1, fmt.Sprintf(format, args...))
***REMOVED***

// Debugf does info logging at verbose level 2.
func (pl *PrefixLogger) Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if !Logger.V(2) ***REMOVED***
		return
	***REMOVED***
	if pl != nil ***REMOVED***
		// Handle nil, so the tests can pass in a nil logger.
		format = pl.prefix + format
		pl.logger.InfoDepth(1, fmt.Sprintf(format, args...))
		return
	***REMOVED***
	InfoDepth(1, fmt.Sprintf(format, args...))
***REMOVED***

// NewPrefixLogger creates a prefix logger with the given prefix.
func NewPrefixLogger(logger DepthLoggerV2, prefix string) *PrefixLogger ***REMOVED***
	return &PrefixLogger***REMOVED***logger: logger, prefix: prefix***REMOVED***
***REMOVED***
