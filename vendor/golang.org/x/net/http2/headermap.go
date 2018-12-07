// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"net/http"
	"strings"
	"sync"
)

var (
	commonBuildOnce   sync.Once
	commonLowerHeader map[string]string // Go-Canonical-Case -> lower-case
	commonCanonHeader map[string]string // lower-case -> Go-Canonical-Case
)

func buildCommonHeaderMapsOnce() ***REMOVED***
	commonBuildOnce.Do(buildCommonHeaderMaps)
***REMOVED***

func buildCommonHeaderMaps() ***REMOVED***
	common := []string***REMOVED***
		"accept",
		"accept-charset",
		"accept-encoding",
		"accept-language",
		"accept-ranges",
		"age",
		"access-control-allow-origin",
		"allow",
		"authorization",
		"cache-control",
		"content-disposition",
		"content-encoding",
		"content-language",
		"content-length",
		"content-location",
		"content-range",
		"content-type",
		"cookie",
		"date",
		"etag",
		"expect",
		"expires",
		"from",
		"host",
		"if-match",
		"if-modified-since",
		"if-none-match",
		"if-unmodified-since",
		"last-modified",
		"link",
		"location",
		"max-forwards",
		"proxy-authenticate",
		"proxy-authorization",
		"range",
		"referer",
		"refresh",
		"retry-after",
		"server",
		"set-cookie",
		"strict-transport-security",
		"trailer",
		"transfer-encoding",
		"user-agent",
		"vary",
		"via",
		"www-authenticate",
	***REMOVED***
	commonLowerHeader = make(map[string]string, len(common))
	commonCanonHeader = make(map[string]string, len(common))
	for _, v := range common ***REMOVED***
		chk := http.CanonicalHeaderKey(v)
		commonLowerHeader[chk] = v
		commonCanonHeader[v] = chk
	***REMOVED***
***REMOVED***

func lowerHeader(v string) string ***REMOVED***
	buildCommonHeaderMapsOnce()
	if s, ok := commonLowerHeader[v]; ok ***REMOVED***
		return s
	***REMOVED***
	return strings.ToLower(v)
***REMOVED***
