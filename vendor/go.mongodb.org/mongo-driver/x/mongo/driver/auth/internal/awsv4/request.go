// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on github.com/aws/aws-sdk-go by Amazon.com, Inc. with code from:
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/request/request.go
// See THIRD-PARTY-NOTICES for original license terms

package awsv4

import (
	"net/http"
	"strings"
)

// Returns host from request
func getHost(r *http.Request) string ***REMOVED***
	if r.Host != "" ***REMOVED***
		return r.Host
	***REMOVED***

	if r.URL == nil ***REMOVED***
		return ""
	***REMOVED***

	return r.URL.Host
***REMOVED***

// Hostname returns u.Host, without any port number.
//
// If Host is an IPv6 literal with a port number, Hostname returns the
// IPv6 literal without the square brackets. IPv6 literals may include
// a zone identifier.
//
// Copied from the Go 1.8 standard library (net/url)
func stripPort(hostport string) string ***REMOVED***
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 ***REMOVED***
		return hostport
	***REMOVED***
	if i := strings.IndexByte(hostport, ']'); i != -1 ***REMOVED***
		return strings.TrimPrefix(hostport[:i], "[")
	***REMOVED***
	return hostport[:colon]
***REMOVED***

// Port returns the port part of u.Host, without the leading colon.
// If u.Host doesn't contain a port, Port returns an empty string.
//
// Copied from the Go 1.8 standard library (net/url)
func portOnly(hostport string) string ***REMOVED***
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 ***REMOVED***
		return ""
	***REMOVED***
	if i := strings.Index(hostport, "]:"); i != -1 ***REMOVED***
		return hostport[i+len("]:"):]
	***REMOVED***
	if strings.Contains(hostport, "]") ***REMOVED***
		return ""
	***REMOVED***
	return hostport[colon+len(":"):]
***REMOVED***

// Returns true if the specified URI is using the standard port
// (i.e. port 80 for HTTP URIs or 443 for HTTPS URIs)
func isDefaultPort(scheme, port string) bool ***REMOVED***
	if port == "" ***REMOVED***
		return true
	***REMOVED***

	lowerCaseScheme := strings.ToLower(scheme)
	if (lowerCaseScheme == "http" && port == "80") || (lowerCaseScheme == "https" && port == "443") ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***
