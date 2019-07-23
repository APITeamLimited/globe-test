// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

package httputil

import (
	"github.com/golang/gddo/httputil/header"
	"net/http"
	"strings"
)

// NegotiateContentEncoding returns the best offered content encoding for the
// request's Accept-Encoding header. If two offers match with equal weight and
// then the offer earlier in the list is preferred. If no offers are
// acceptable, then "" is returned.
func NegotiateContentEncoding(r *http.Request, offers []string) string ***REMOVED***
	bestOffer := "identity"
	bestQ := -1.0
	specs := header.ParseAccept(r.Header, "Accept-Encoding")
	for _, offer := range offers ***REMOVED***
		for _, spec := range specs ***REMOVED***
			if spec.Q > bestQ &&
				(spec.Value == "*" || spec.Value == offer) ***REMOVED***
				bestQ = spec.Q
				bestOffer = offer
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if bestQ == 0 ***REMOVED***
		bestOffer = ""
	***REMOVED***
	return bestOffer
***REMOVED***

// NegotiateContentType returns the best offered content type for the request's
// Accept header. If two offers match with equal weight, then the more specific
// offer is preferred.  For example, text/* trumps */*. If two offers match
// with equal weight and specificity, then the offer earlier in the list is
// preferred. If no offers match, then defaultOffer is returned.
func NegotiateContentType(r *http.Request, offers []string, defaultOffer string) string ***REMOVED***
	bestOffer := defaultOffer
	bestQ := -1.0
	bestWild := 3
	specs := header.ParseAccept(r.Header, "Accept")
	for _, offer := range offers ***REMOVED***
		for _, spec := range specs ***REMOVED***
			switch ***REMOVED***
			case spec.Q == 0.0:
				// ignore
			case spec.Q < bestQ:
				// better match found
			case spec.Value == "*/*":
				if spec.Q > bestQ || bestWild > 2 ***REMOVED***
					bestQ = spec.Q
					bestWild = 2
					bestOffer = offer
				***REMOVED***
			case strings.HasSuffix(spec.Value, "/*"):
				if strings.HasPrefix(offer, spec.Value[:len(spec.Value)-1]) &&
					(spec.Q > bestQ || bestWild > 1) ***REMOVED***
					bestQ = spec.Q
					bestWild = 1
					bestOffer = offer
				***REMOVED***
			default:
				if spec.Value == offer &&
					(spec.Q > bestQ || bestWild > 0) ***REMOVED***
					bestQ = spec.Q
					bestWild = 0
					bestOffer = offer
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return bestOffer
***REMOVED***
