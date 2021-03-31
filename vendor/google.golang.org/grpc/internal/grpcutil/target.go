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

// Package grpcutil provides a bunch of utility functions to be used across the
// gRPC codebase.
package grpcutil

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

// split2 returns the values from strings.SplitN(s, sep, 2).
// If sep is not found, it returns ("", "", false) instead.
func split2(s, sep string) (string, string, bool) ***REMOVED***
	spl := strings.SplitN(s, sep, 2)
	if len(spl) < 2 ***REMOVED***
		return "", "", false
	***REMOVED***
	return spl[0], spl[1], true
***REMOVED***

// ParseTarget splits target into a resolver.Target struct containing scheme,
// authority and endpoint. skipUnixColonParsing indicates that the parse should
// not parse "unix:[path]" cases. This should be true in cases where a custom
// dialer is present, to prevent a behavior change.
//
// If target is not a valid scheme://authority/endpoint as specified in
// https://github.com/grpc/grpc/blob/master/doc/naming.md,
// it returns ***REMOVED***Endpoint: target***REMOVED***.
func ParseTarget(target string, skipUnixColonParsing bool) (ret resolver.Target) ***REMOVED***
	var ok bool
	if strings.HasPrefix(target, "unix-abstract:") ***REMOVED***
		if strings.HasPrefix(target, "unix-abstract://") ***REMOVED***
			// Maybe, with Authority specified, try to parse it
			var remain string
			ret.Scheme, remain, _ = split2(target, "://")
			ret.Authority, ret.Endpoint, ok = split2(remain, "/")
			if !ok ***REMOVED***
				// No Authority, add the "//" back
				ret.Endpoint = "//" + remain
			***REMOVED*** else ***REMOVED***
				// Found Authority, add the "/" back
				ret.Endpoint = "/" + ret.Endpoint
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Without Authority specified, split target on ":"
			ret.Scheme, ret.Endpoint, _ = split2(target, ":")
		***REMOVED***
		return ret
	***REMOVED***
	ret.Scheme, ret.Endpoint, ok = split2(target, "://")
	if !ok ***REMOVED***
		if strings.HasPrefix(target, "unix:") && !skipUnixColonParsing ***REMOVED***
			// Handle the "unix:[local/path]" and "unix:[/absolute/path]" cases,
			// because splitting on :// only handles the
			// "unix://[/absolute/path]" case. Only handle if the dialer is nil,
			// to avoid a behavior change with custom dialers.
			return resolver.Target***REMOVED***Scheme: "unix", Endpoint: target[len("unix:"):]***REMOVED***
		***REMOVED***
		return resolver.Target***REMOVED***Endpoint: target***REMOVED***
	***REMOVED***
	ret.Authority, ret.Endpoint, ok = split2(ret.Endpoint, "/")
	if !ok ***REMOVED***
		return resolver.Target***REMOVED***Endpoint: target***REMOVED***
	***REMOVED***
	if ret.Scheme == "unix" ***REMOVED***
		// Add the "/" back in the unix case, so the unix resolver receives the
		// actual endpoint in the "unix://[/absolute/path]" case.
		ret.Endpoint = "/" + ret.Endpoint
	***REMOVED***
	return ret
***REMOVED***
