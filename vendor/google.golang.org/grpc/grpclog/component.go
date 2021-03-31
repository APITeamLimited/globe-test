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

	"google.golang.org/grpc/internal/grpclog"
)

// componentData records the settings for a component.
type componentData struct ***REMOVED***
	name string
***REMOVED***

var cache = map[string]*componentData***REMOVED******REMOVED***

func (c *componentData) InfoDepth(depth int, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	args = append([]interface***REMOVED******REMOVED******REMOVED***"[" + string(c.name) + "]"***REMOVED***, args...)
	grpclog.InfoDepth(depth+1, args...)
***REMOVED***

func (c *componentData) WarningDepth(depth int, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	args = append([]interface***REMOVED******REMOVED******REMOVED***"[" + string(c.name) + "]"***REMOVED***, args...)
	grpclog.WarningDepth(depth+1, args...)
***REMOVED***

func (c *componentData) ErrorDepth(depth int, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	args = append([]interface***REMOVED******REMOVED******REMOVED***"[" + string(c.name) + "]"***REMOVED***, args...)
	grpclog.ErrorDepth(depth+1, args...)
***REMOVED***

func (c *componentData) FatalDepth(depth int, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	args = append([]interface***REMOVED******REMOVED******REMOVED***"[" + string(c.name) + "]"***REMOVED***, args...)
	grpclog.FatalDepth(depth+1, args...)
***REMOVED***

func (c *componentData) Info(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.InfoDepth(1, args...)
***REMOVED***

func (c *componentData) Warning(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.WarningDepth(1, args...)
***REMOVED***

func (c *componentData) Error(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.ErrorDepth(1, args...)
***REMOVED***

func (c *componentData) Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.FatalDepth(1, args...)
***REMOVED***

func (c *componentData) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.InfoDepth(1, fmt.Sprintf(format, args...))
***REMOVED***

func (c *componentData) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.WarningDepth(1, fmt.Sprintf(format, args...))
***REMOVED***

func (c *componentData) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.ErrorDepth(1, fmt.Sprintf(format, args...))
***REMOVED***

func (c *componentData) Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.FatalDepth(1, fmt.Sprintf(format, args...))
***REMOVED***

func (c *componentData) Infoln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.InfoDepth(1, args...)
***REMOVED***

func (c *componentData) Warningln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.WarningDepth(1, args...)
***REMOVED***

func (c *componentData) Errorln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.ErrorDepth(1, args...)
***REMOVED***

func (c *componentData) Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.FatalDepth(1, args...)
***REMOVED***

func (c *componentData) V(l int) bool ***REMOVED***
	return V(l)
***REMOVED***

// Component creates a new component and returns it for logging. If a component
// with the name already exists, nothing will be created and it will be
// returned. SetLoggerV2 will panic if it is called with a logger created by
// Component.
func Component(componentName string) DepthLoggerV2 ***REMOVED***
	if cData, ok := cache[componentName]; ok ***REMOVED***
		return cData
	***REMOVED***
	c := &componentData***REMOVED***componentName***REMOVED***
	cache[componentName] = c
	return c
***REMOVED***
