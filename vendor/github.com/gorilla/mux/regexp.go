// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type routeRegexpOptions struct ***REMOVED***
	strictSlash    bool
	useEncodedPath bool
***REMOVED***

type regexpType int

const (
	regexpTypePath   regexpType = 0
	regexpTypeHost   regexpType = 1
	regexpTypePrefix regexpType = 2
	regexpTypeQuery  regexpType = 3
)

// newRouteRegexp parses a route template and returns a routeRegexp,
// used to match a host, a path or a query string.
//
// It will extract named variables, assemble a regexp to be matched, create
// a "reverse" template to build URLs and compile regexps to validate variable
// values used in URL building.
//
// Previously we accepted only Python-like identifiers for variable
// names ([a-zA-Z_][a-zA-Z0-9_]*), but currently the only restriction is that
// name and pattern can't be empty, and names can't contain a colon.
func newRouteRegexp(tpl string, typ regexpType, options routeRegexpOptions) (*routeRegexp, error) ***REMOVED***
	// Check if it is well-formed.
	idxs, errBraces := braceIndices(tpl)
	if errBraces != nil ***REMOVED***
		return nil, errBraces
	***REMOVED***
	// Backup the original.
	template := tpl
	// Now let's parse it.
	defaultPattern := "[^/]+"
	if typ == regexpTypeQuery ***REMOVED***
		defaultPattern = ".*"
	***REMOVED*** else if typ == regexpTypeHost ***REMOVED***
		defaultPattern = "[^.]+"
	***REMOVED***
	// Only match strict slash if not matching
	if typ != regexpTypePath ***REMOVED***
		options.strictSlash = false
	***REMOVED***
	// Set a flag for strictSlash.
	endSlash := false
	if options.strictSlash && strings.HasSuffix(tpl, "/") ***REMOVED***
		tpl = tpl[:len(tpl)-1]
		endSlash = true
	***REMOVED***
	varsN := make([]string, len(idxs)/2)
	varsR := make([]*regexp.Regexp, len(idxs)/2)
	pattern := bytes.NewBufferString("")
	pattern.WriteByte('^')
	reverse := bytes.NewBufferString("")
	var end int
	var err error
	for i := 0; i < len(idxs); i += 2 ***REMOVED***
		// Set all values we are interested in.
		raw := tpl[end:idxs[i]]
		end = idxs[i+1]
		parts := strings.SplitN(tpl[idxs[i]+1:end-1], ":", 2)
		name := parts[0]
		patt := defaultPattern
		if len(parts) == 2 ***REMOVED***
			patt = parts[1]
		***REMOVED***
		// Name or pattern can't be empty.
		if name == "" || patt == "" ***REMOVED***
			return nil, fmt.Errorf("mux: missing name or pattern in %q",
				tpl[idxs[i]:end])
		***REMOVED***
		// Build the regexp pattern.
		fmt.Fprintf(pattern, "%s(?P<%s>%s)", regexp.QuoteMeta(raw), varGroupName(i/2), patt)

		// Build the reverse template.
		fmt.Fprintf(reverse, "%s%%s", raw)

		// Append variable name and compiled pattern.
		varsN[i/2] = name
		varsR[i/2], err = regexp.Compile(fmt.Sprintf("^%s$", patt))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	// Add the remaining.
	raw := tpl[end:]
	pattern.WriteString(regexp.QuoteMeta(raw))
	if options.strictSlash ***REMOVED***
		pattern.WriteString("[/]?")
	***REMOVED***
	if typ == regexpTypeQuery ***REMOVED***
		// Add the default pattern if the query value is empty
		if queryVal := strings.SplitN(template, "=", 2)[1]; queryVal == "" ***REMOVED***
			pattern.WriteString(defaultPattern)
		***REMOVED***
	***REMOVED***
	if typ != regexpTypePrefix ***REMOVED***
		pattern.WriteByte('$')
	***REMOVED***
	reverse.WriteString(raw)
	if endSlash ***REMOVED***
		reverse.WriteByte('/')
	***REMOVED***
	// Compile full regexp.
	reg, errCompile := regexp.Compile(pattern.String())
	if errCompile != nil ***REMOVED***
		return nil, errCompile
	***REMOVED***

	// Check for capturing groups which used to work in older versions
	if reg.NumSubexp() != len(idxs)/2 ***REMOVED***
		panic(fmt.Sprintf("route %s contains capture groups in its regexp. ", template) +
			"Only non-capturing groups are accepted: e.g. (?:pattern) instead of (pattern)")
	***REMOVED***

	// Done!
	return &routeRegexp***REMOVED***
		template:   template,
		regexpType: typ,
		options:    options,
		regexp:     reg,
		reverse:    reverse.String(),
		varsN:      varsN,
		varsR:      varsR,
	***REMOVED***, nil
***REMOVED***

// routeRegexp stores a regexp to match a host or path and information to
// collect and validate route variables.
type routeRegexp struct ***REMOVED***
	// The unmodified template.
	template string
	// The type of match
	regexpType regexpType
	// Options for matching
	options routeRegexpOptions
	// Expanded regexp.
	regexp *regexp.Regexp
	// Reverse template.
	reverse string
	// Variable names.
	varsN []string
	// Variable regexps (validators).
	varsR []*regexp.Regexp
***REMOVED***

// Match matches the regexp against the URL host or path.
func (r *routeRegexp) Match(req *http.Request, match *RouteMatch) bool ***REMOVED***
	if r.regexpType != regexpTypeHost ***REMOVED***
		if r.regexpType == regexpTypeQuery ***REMOVED***
			return r.matchQueryString(req)
		***REMOVED***
		path := req.URL.Path
		if r.options.useEncodedPath ***REMOVED***
			path = req.URL.EscapedPath()
		***REMOVED***
		return r.regexp.MatchString(path)
	***REMOVED***

	return r.regexp.MatchString(getHost(req))
***REMOVED***

// url builds a URL part using the given values.
func (r *routeRegexp) url(values map[string]string) (string, error) ***REMOVED***
	urlValues := make([]interface***REMOVED******REMOVED***, len(r.varsN))
	for k, v := range r.varsN ***REMOVED***
		value, ok := values[v]
		if !ok ***REMOVED***
			return "", fmt.Errorf("mux: missing route variable %q", v)
		***REMOVED***
		if r.regexpType == regexpTypeQuery ***REMOVED***
			value = url.QueryEscape(value)
		***REMOVED***
		urlValues[k] = value
	***REMOVED***
	rv := fmt.Sprintf(r.reverse, urlValues...)
	if !r.regexp.MatchString(rv) ***REMOVED***
		// The URL is checked against the full regexp, instead of checking
		// individual variables. This is faster but to provide a good error
		// message, we check individual regexps if the URL doesn't match.
		for k, v := range r.varsN ***REMOVED***
			if !r.varsR[k].MatchString(values[v]) ***REMOVED***
				return "", fmt.Errorf(
					"mux: variable %q doesn't match, expected %q", values[v],
					r.varsR[k].String())
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return rv, nil
***REMOVED***

// getURLQuery returns a single query parameter from a request URL.
// For a URL with foo=bar&baz=ding, we return only the relevant key
// value pair for the routeRegexp.
func (r *routeRegexp) getURLQuery(req *http.Request) string ***REMOVED***
	if r.regexpType != regexpTypeQuery ***REMOVED***
		return ""
	***REMOVED***
	templateKey := strings.SplitN(r.template, "=", 2)[0]
	for key, vals := range req.URL.Query() ***REMOVED***
		if key == templateKey && len(vals) > 0 ***REMOVED***
			return key + "=" + vals[0]
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func (r *routeRegexp) matchQueryString(req *http.Request) bool ***REMOVED***
	return r.regexp.MatchString(r.getURLQuery(req))
***REMOVED***

// braceIndices returns the first level curly brace indices from a string.
// It returns an error in case of unbalanced braces.
func braceIndices(s string) ([]int, error) ***REMOVED***
	var level, idx int
	var idxs []int
	for i := 0; i < len(s); i++ ***REMOVED***
		switch s[i] ***REMOVED***
		case '***REMOVED***':
			if level++; level == 1 ***REMOVED***
				idx = i
			***REMOVED***
		case '***REMOVED***':
			if level--; level == 0 ***REMOVED***
				idxs = append(idxs, idx, i+1)
			***REMOVED*** else if level < 0 ***REMOVED***
				return nil, fmt.Errorf("mux: unbalanced braces in %q", s)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if level != 0 ***REMOVED***
		return nil, fmt.Errorf("mux: unbalanced braces in %q", s)
	***REMOVED***
	return idxs, nil
***REMOVED***

// varGroupName builds a capturing group name for the indexed variable.
func varGroupName(idx int) string ***REMOVED***
	return "v" + strconv.Itoa(idx)
***REMOVED***

// ----------------------------------------------------------------------------
// routeRegexpGroup
// ----------------------------------------------------------------------------

// routeRegexpGroup groups the route matchers that carry variables.
type routeRegexpGroup struct ***REMOVED***
	host    *routeRegexp
	path    *routeRegexp
	queries []*routeRegexp
***REMOVED***

// setMatch extracts the variables from the URL once a route matches.
func (v *routeRegexpGroup) setMatch(req *http.Request, m *RouteMatch, r *Route) ***REMOVED***
	// Store host variables.
	if v.host != nil ***REMOVED***
		host := getHost(req)
		matches := v.host.regexp.FindStringSubmatchIndex(host)
		if len(matches) > 0 ***REMOVED***
			extractVars(host, matches, v.host.varsN, m.Vars)
		***REMOVED***
	***REMOVED***
	path := req.URL.Path
	if r.useEncodedPath ***REMOVED***
		path = req.URL.EscapedPath()
	***REMOVED***
	// Store path variables.
	if v.path != nil ***REMOVED***
		matches := v.path.regexp.FindStringSubmatchIndex(path)
		if len(matches) > 0 ***REMOVED***
			extractVars(path, matches, v.path.varsN, m.Vars)
			// Check if we should redirect.
			if v.path.options.strictSlash ***REMOVED***
				p1 := strings.HasSuffix(path, "/")
				p2 := strings.HasSuffix(v.path.template, "/")
				if p1 != p2 ***REMOVED***
					u, _ := url.Parse(req.URL.String())
					if p1 ***REMOVED***
						u.Path = u.Path[:len(u.Path)-1]
					***REMOVED*** else ***REMOVED***
						u.Path += "/"
					***REMOVED***
					m.Handler = http.RedirectHandler(u.String(), 301)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Store query string variables.
	for _, q := range v.queries ***REMOVED***
		queryURL := q.getURLQuery(req)
		matches := q.regexp.FindStringSubmatchIndex(queryURL)
		if len(matches) > 0 ***REMOVED***
			extractVars(queryURL, matches, q.varsN, m.Vars)
		***REMOVED***
	***REMOVED***
***REMOVED***

// getHost tries its best to return the request host.
func getHost(r *http.Request) string ***REMOVED***
	if r.URL.IsAbs() ***REMOVED***
		return r.URL.Host
	***REMOVED***
	host := r.Host
	// Slice off any port information.
	if i := strings.Index(host, ":"); i != -1 ***REMOVED***
		host = host[:i]
	***REMOVED***
	return host

***REMOVED***

func extractVars(input string, matches []int, names []string, output map[string]string) ***REMOVED***
	for i, name := range names ***REMOVED***
		output[name] = input[matches[2*i+2]:matches[2*i+3]]
	***REMOVED***
***REMOVED***
