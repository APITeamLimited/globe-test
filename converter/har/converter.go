/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package har

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func Convert(h HAR, includeCodeCheck bool, batchTime uint, only, skip []string) (string, error) ***REMOVED***
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if includeCodeCheck ***REMOVED***
		fmt.Fprint(w, "import ***REMOVED*** group, check, sleep ***REMOVED*** from 'k6';\n")
	***REMOVED*** else ***REMOVED***
		fmt.Fprint(w, "import ***REMOVED*** group, sleep ***REMOVED*** from 'k6';\n")
	***REMOVED***
	fmt.Fprint(w, "import http from 'k6/http';\n\n")

	fmt.Fprintf(w, "// Version: %v\n", h.Log.Version)
	fmt.Fprintf(w, "// Creator: %v\n", h.Log.Creator.Name)
	if h.Log.Browser != nil ***REMOVED***
		fmt.Fprintf(w, "// Browser: %v\n", h.Log.Browser.Name)
	***REMOVED***
	if h.Log.Comment != "" ***REMOVED***
		fmt.Fprintf(w, "// %v\n", h.Log.Comment)
	***REMOVED***

	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "export default function() ***REMOVED***\n\n")

	pages := h.Log.Pages
	sort.Sort(PageByStarted(pages))

	// Grouping by page and URL filtering
	pageEntries := make(map[string][]*Entry)
	for _, e := range h.Log.Entries ***REMOVED***

		// URL filtering
		u, err := url.Parse(e.Request.URL)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if !IsAllowedURL(u.Host, only, skip) ***REMOVED***
			continue
		***REMOVED***

		// Avoid multipart/form-data requests until k6 scripts can support binary data
		if e.Request.PostData != nil && strings.HasPrefix(e.Request.PostData.MimeType, "multipart/form-data") ***REMOVED***
			continue
		***REMOVED***

		// Create new group o adding page to a existing one
		if _, ok := pageEntries[e.Pageref]; !ok ***REMOVED***
			pageEntries[e.Pageref] = append([]*Entry***REMOVED******REMOVED***, e)
		***REMOVED*** else ***REMOVED***
			pageEntries[e.Pageref] = append(pageEntries[e.Pageref], e)
		***REMOVED***
	***REMOVED***

	for i, page := range pages ***REMOVED***

		entries := pageEntries[page.ID]
		fmt.Fprintf(w, "\tgroup(\"%s - %s\", function() ***REMOVED***\n", page.ID, page.Title)

		sort.Sort(EntryByStarted(entries))
		batches := SplitEntriesInBatches(entries, batchTime)

		fmt.Fprint(w, "\t\tlet req, res;\n")

		for j, batchEntries := range batches ***REMOVED***

			fmt.Fprint(w, "\t\treq = [")
			for k, e := range batchEntries ***REMOVED***
				r, err := buildK6RequestObject(e.Request)
				if err != nil ***REMOVED***
					return "", err
				***REMOVED***
				fmt.Fprintf(w, "%v", r)
				if k != len(batchEntries)-1 ***REMOVED***
					fmt.Fprint(w, ",")
				***REMOVED***
			***REMOVED***
			fmt.Fprint(w, "];\n")
			fmt.Fprint(w, "\t\tres = http.batch(req);\n")

			if includeCodeCheck ***REMOVED***
				for k, e := range batchEntries ***REMOVED***
					if e.Response.Status > 0 ***REMOVED***
						fmt.Fprintf(w, "\t\tcheck(res[%v], ***REMOVED***\n\t\t\"status is %v\": (r) => r.status === %v,\n\t***REMOVED***);\n", k, e.Response.Status, e.Response.Status)
					***REMOVED***
				***REMOVED***
			***REMOVED***

			if j != len(batches)-1 ***REMOVED***
				lastBatchEntry := batchEntries[len(batchEntries)-1]
				firstBatchEntry := batches[j+1][0]
				t := firstBatchEntry.StartedDateTime.Sub(lastBatchEntry.StartedDateTime).Seconds()
				fmt.Fprintf(w, "\t\tsleep(%.2f);\n", t)
			***REMOVED***
		***REMOVED***

		if i == len(pages)-1 ***REMOVED***
			// Last page; add random sleep time at the group completion
			fmt.Fprint(w, "\t\t// Random sleep between 2s and 4s\n")
			fmt.Fprint(w, "\t\tsleep(Math.floor(Math.random()*3+2));\n")
		***REMOVED*** else ***REMOVED***
			// Add sleep time at the end of the group
			nextPage := pages[i+1]
			lastEntry := entries[len(entries)-1]
			t := nextPage.StartedDateTime.Sub(lastEntry.StartedDateTime).Seconds()
			if t < 0.01 ***REMOVED***
				t = 0.5
			***REMOVED***
			fmt.Fprintf(w, "\t\tsleep(%.2f);\n", t)
		***REMOVED***

		fmt.Fprint(w, "\t***REMOVED***);\n")
	***REMOVED***

	fmt.Fprint(w, "\n***REMOVED***\n")
	if err := w.Flush(); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return b.String(), nil
***REMOVED***

func buildK6RequestObject(req *Request) (string, error) ***REMOVED***
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	fmt.Fprint(w, "***REMOVED***\n")

	method := strings.ToLower(req.Method)
	if method == "delete" ***REMOVED***
		method = "del"
	***REMOVED***
	fmt.Fprintf(w, `"method": %q, "url": %q`, method, req.URL)

	if req.PostData != nil && method != "get" ***REMOVED***
		postParams, plainText, err := buildK6Body(req)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED*** else if len(postParams) > 0 ***REMOVED***
			fmt.Fprintf(w, `, "body": ***REMOVED*** %s ***REMOVED***`, strings.Join(postParams, ", "))
		***REMOVED*** else if plainText != "" ***REMOVED***
			fmt.Fprintf(w, `, "body": %q`, plainText)
		***REMOVED***
	***REMOVED***

	var params []string
	var cookies []string
	for _, c := range req.Cookies ***REMOVED***
		cookies = append(cookies, fmt.Sprintf(`%q: %q`, c.Name, c.Value))
	***REMOVED***
	if len(cookies) > 0 ***REMOVED***
		params = append(params, fmt.Sprintf(`"cookies": ***REMOVED*** %s ***REMOVED***`, strings.Join(cookies, ", ")))
	***REMOVED***

	if headers := buildK6Headers(req.Headers); len(headers) > 0 ***REMOVED***
		params = append(params, fmt.Sprintf(`"headers": ***REMOVED*** %s ***REMOVED***`, strings.Join(headers, ", ")))
	***REMOVED***

	if len(params) > 0 ***REMOVED***
		fmt.Fprintf(w, `, "params": ***REMOVED*** %s ***REMOVED***`, strings.Join(params, ", "))
	***REMOVED***

	fmt.Fprint(w, "***REMOVED***")
	if err := w.Flush(); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var buffer bytes.Buffer
	err := json.Indent(&buffer, b.Bytes(), "\t\t", "\t")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return buffer.String(), nil
***REMOVED***

func buildK6Headers(headers []Header) []string ***REMOVED***
	var h []string
	if len(headers) > 0 ***REMOVED***
		m := make(map[string]Header)
		for _, header := range headers ***REMOVED***
			name := strings.ToLower(header.Name)
			_, exists := m[name]
			// Avoid SPDY's, duplicated or cookie headers
			if !exists && name[0] != ':' && name != "cookie" ***REMOVED***
				m[strings.ToLower(header.Name)] = header
				h = append(h, fmt.Sprintf("%q: %q", header.Name, header.Value))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return h
***REMOVED***

func buildK6Body(req *Request) ([]string, string, error) ***REMOVED***
	var postParams []string
	if req.PostData.MimeType == "application/x-www-form-urlencoded" && len(req.PostData.Params) > 0 ***REMOVED***
		for _, p := range req.PostData.Params ***REMOVED***
			n, err := url.QueryUnescape(p.Name)
			if err != nil ***REMOVED***
				return postParams, "", err
			***REMOVED***
			v, err := url.QueryUnescape(p.Value)
			if err != nil ***REMOVED***
				return postParams, "", err
			***REMOVED***
			postParams = append(postParams, fmt.Sprintf(`%q: %q`, n, v))
		***REMOVED***
		return postParams, "", nil
	***REMOVED***
	return postParams, req.PostData.Text, nil
***REMOVED***
