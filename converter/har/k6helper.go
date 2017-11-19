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
        "math/rand"
        "net/url"
        "sort"
        "strings"
        "time"
)

const (
        Get   = "get"
        Post  = "post"
        Del   = "del"
        Put   = "put"
        Patch = "patch"
)

// build a K6 request
func BuildK6Request(method, uri, data string, headers []Header, cookies []Cookie) (string, error) ***REMOVED***
        var b bytes.Buffer
        w := bufio.NewWriter(&b)

        // method and url
        method = strings.ToLower(method)
        switch method ***REMOVED***
        case Get, Post, Put, Patch:
                fmt.Fprintf(w, "http.%v(\"%v\"", method, uri)
        case "delete":
                fmt.Fprintf(w, "http.del(\"%v\"", uri)
        default:
                fmt.Fprintf(w, "http.request(\"%v\",\"%v\"", method, uri)
        ***REMOVED***

        // data
        if data != "" ***REMOVED***
                fmt.Fprintf(w, ",\n\t\"%s\"", url.QueryEscape(data))
        ***REMOVED*** else if method != Get ***REMOVED***
                fmt.Fprint(w, ",\n\tnull")
        ***REMOVED***

        // Add cookie as header
        c := BuildK6CookiesValues(cookies)
        if c != "" ***REMOVED***
                headers = append(headers, Header***REMOVED***Name: "Cookie", Value: c***REMOVED***)
        ***REMOVED***

        if header := BuildK6Headers(headers); len(header) > 0 ***REMOVED***
                fmt.Fprintf(w, ",\n\t***REMOVED*** %v ***REMOVED***", header)
        ***REMOVED***

        fmt.Fprint(w, "\n);\n")

        if err := w.Flush(); err != nil ***REMOVED***
                return "", err
        ***REMOVED***
        return b.String(), nil
***REMOVED***

// Build a K6 request object for batch requests
func BuildK6RequestObject(method, uri, data string, headers []Header, cookies []Cookie) (string, error) ***REMOVED***
        var b bytes.Buffer
        w := bufio.NewWriter(&b)

        fmt.Fprint(w, "***REMOVED***\n")

        method = strings.ToLower(method)
        if method == "delete" ***REMOVED***
                method = Del
        ***REMOVED***
        fmt.Fprintf(w, "\"method\": %q,\n\"url\": %q", method, uri)

        // data
        if data != "" && method != Get ***REMOVED***
                fmt.Fprintf(w, ", \"body\": \n%q\n", data)
        ***REMOVED***

        // Add cookie as header
        if c := BuildK6CookiesValues(cookies); c != "" ***REMOVED***
                headers = append(headers, Header***REMOVED***Name: "Cookie", Value: c***REMOVED***)
        ***REMOVED***

        if header := BuildK6Headers(headers); len(header) > 0 ***REMOVED***
                fmt.Fprintf(w, ", \"params\": ***REMOVED***\n%s\n***REMOVED***\n", header)
        ***REMOVED***

        fmt.Fprint(w, "***REMOVED***")
        if err := w.Flush(); err != nil ***REMOVED***
                return "", err
        ***REMOVED***

        // json indentation
        var buffer bytes.Buffer
        err := json.Indent(&buffer, b.Bytes(), "", "    ")
        if err != nil ***REMOVED***
                return "", err
        ***REMOVED***

        return buffer.String(), nil
***REMOVED***

// Build the string representation of a K6 headers object from a given HAR.NVP array
func BuildK6Headers(headers []Header) string ***REMOVED***
        if len(headers) == 0 ***REMOVED***
                return ""
        ***REMOVED***
        m := make(map[string]Header)

        var h []string
        for _, header := range headers ***REMOVED***
                if header.Name[0] != ':' ***REMOVED*** // avoid SPDY's colon headers
                        // avoid header duplicity
                        _, exists := m[strings.ToLower(header.Name)]
                        if !exists ***REMOVED***
                                m[strings.ToLower(header.Name)] = header
                                h = append(h, fmt.Sprintf("%q : %q", header.Name, header.Value))
                        ***REMOVED***
                ***REMOVED***
        ***REMOVED***

        return fmt.Sprintf("\"headers\" : ***REMOVED*** %v ***REMOVED***", strings.Join(h, ", "))
***REMOVED***

// Build the string representation of K6 cookie values from a given HAR.Cookie array
func BuildK6CookiesValues(cookies []Cookie) string ***REMOVED***
        if len(cookies) == 0 ***REMOVED***
                return ""
        ***REMOVED***

        var c []string
        for _, cookie := range cookies ***REMOVED***
                c = append(c, fmt.Sprintf("%v=%v", cookie.Name, cookie.Value))
        ***REMOVED***

        return strings.Join(c, "; ")
***REMOVED***

// Returns true if the given url is allowed from the only (only domains) and skip (skip domains) values, otherwise false
func IsAllowedURL(url string, only, skip []string) bool ***REMOVED***
        if len(only) != 0 ***REMOVED***
                for _, v := range only ***REMOVED***
                        v = strings.Trim(v, " ")
                        if v != "" && strings.Contains(url, v) ***REMOVED***
                                return true
                        ***REMOVED***
                ***REMOVED***
                return false
        ***REMOVED***
        if len(skip) != 0 ***REMOVED***
                for _, v := range skip ***REMOVED***
                        v = strings.Trim(v, " ")
                        if v != "" && strings.Contains(url, v) ***REMOVED***
                                return false
                        ***REMOVED***
                ***REMOVED***
        ***REMOVED***
        return true
***REMOVED***

func WriteK6Script(w *bufio.Writer, h *HAR, includeCodeCheck bool, batchTime uint, only, skip []string, maxRequestsBatch uint) error ***REMOVED***
        if includeCodeCheck ***REMOVED***
                fmt.Fprint(w, "import ***REMOVED*** group, check, sleep ***REMOVED*** from 'k6';\n")
        ***REMOVED*** else ***REMOVED***
                fmt.Fprint(w, "import ***REMOVED*** group, sleep ***REMOVED*** from 'k6';\n")
        ***REMOVED***
        fmt.Fprint(w, "import http from 'k6/http';\n\n")

        fmt.Fprintf(w, "// Version: %v\n", h.Log.Version)
        fmt.Fprintf(w, "// Creator: %v\n", h.Log.Creator.Name)

        if h.Log.Browser != nil ***REMOVED*** // browser is optional
                fmt.Fprintf(w, "// Browser: %v\n", h.Log.Browser.Name)
        ***REMOVED***
        if h.Log.Comment != "" ***REMOVED***
                fmt.Fprintf(w, "// %v\n", h.Log.Comment)
        ***REMOVED***

        fmt.Fprint(w, "\nexport default function() ***REMOVED***\n\nlet res;\n\n")

        // name used by group entries
        pagenames := make(map[string]string)
        for _, e := range h.Log.Pages ***REMOVED***
                pagenames[e.ID] = fmt.Sprintf("%s - %s", e.ID, e.Title)
        ***REMOVED***

        // grouping by page and URL filtering
        groups := make(map[string][]*Entry)
        var nameGroups []string
        for _, e := range h.Log.Entries ***REMOVED***

                // URL filtering
                u, err := url.Parse(e.Request.URL)
                if err != nil ***REMOVED***
                        return err
                ***REMOVED***
                if !IsAllowedURL(u.Host, only, skip) ***REMOVED***
                        continue
                ***REMOVED***

                // avoid multipart/form-data requests until k6 scripts can support binary data
                if e.Request.PostData != nil && strings.HasPrefix(e.Request.PostData.MimeType, "multipart/form-data") ***REMOVED***
                        continue
                ***REMOVED***

                // create new group o adding page to a existing one
                if _, ok := groups[e.Pageref]; !ok ***REMOVED***
                        groups[e.Pageref] = append([]*Entry***REMOVED******REMOVED***, e)
                        nameGroups = append(nameGroups, e.Pageref)
                ***REMOVED*** else ***REMOVED***
                        groups[e.Pageref] = append(groups[e.Pageref], e)
                ***REMOVED***
        ***REMOVED***

        for _, n := range nameGroups ***REMOVED***

                // sort entries by requests started date time
                sort.Sort(byRequestDate(groups[n]))

                fmt.Fprintf(w, "group(\"%v\", function() ***REMOVED***\n", pagenames[n])

                if batchTime > 0 ***REMOVED***
                        // batch mode, multiple HTTP requests together
                        entries := groupHarEntriesByIntervals(groups[n], groups[n][0].StartedDateTime, batchTime, maxRequestsBatch)

                        fmt.Fprint(w, "\tlet req\n")

                        for _, e := range entries ***REMOVED***
                                var statuses []int

                                fmt.Fprint(w, "\treq = [\n")
                                for _, r := range e ***REMOVED***
                                        var data string
                                        if r.Request.PostData != nil ***REMOVED***
                                                data = r.Request.PostData.Text
                                        ***REMOVED***
                                        b, err := BuildK6RequestObject(r.Request.Method, r.Request.URL, data, r.Request.Headers, r.Request.Cookies)
                                        if err != nil ***REMOVED***
                                                return err
                                        ***REMOVED***
                                        fmt.Fprintf(w, "%v,\n", b)
                                        statuses = append(statuses, r.Response.Status)
                                ***REMOVED***
                                fmt.Fprint(w, "];\n")

                                fmt.Fprint(w, "\tres = http.batch(req);\n")
                                if includeCodeCheck ***REMOVED***
                                        for i, s := range statuses ***REMOVED***
                                                if s > 0 ***REMOVED*** // avoid empty responses, browsers with adblockers, antivirus.. can block HTTP requests
                                                        fmt.Fprintf(w, "\tcheck(res[%v], ***REMOVED***\n\t\t\"status is %v\": (r) => r.status === %v,\n\t***REMOVED***);\n", i, s, s)
                                                ***REMOVED***
                                        ***REMOVED***
                                ***REMOVED***
                        ***REMOVED***
                ***REMOVED*** else ***REMOVED***
                        // no batch mode
                        for _, e := range groups[n] ***REMOVED***
                                var data string
                                if e.Request.PostData != nil ***REMOVED***
                                        data = e.Request.PostData.Text
                                ***REMOVED***
                                b, err := BuildK6Request(e.Request.Method, e.Request.URL, data, e.Request.Headers, e.Request.Cookies)
                                if err != nil ***REMOVED***
                                        return err
                                ***REMOVED***
                                fmt.Fprintf(w, "res = %v", b)
                                if includeCodeCheck ***REMOVED***
                                        if e.Response.Status > 0 ***REMOVED*** // avoid empty responses, browsers with adblockers, antivirus.. can block HTTP requests
                                                fmt.Fprintf(w, "check(res, ***REMOVED***\n\"status is %v\": (r) => r.status === %v,\n***REMOVED***);\n", e.Response.Status, e.Response.Status)
                                        ***REMOVED***
                                ***REMOVED***
                        ***REMOVED***
                ***REMOVED***

                // random sleep from 100ms to 500ms
                fmt.Fprintf(w, "\tsleep(%.1f);\n", float64(rand.Intn(500)+100)/1000)

                fmt.Fprint(w, "***REMOVED***);\n")
        ***REMOVED***

        fmt.Fprint(w, "\n***REMOVED***\n")
        err := w.Flush()

        return err
***REMOVED***

func groupHarEntriesByIntervals(entries []*Entry, starttime time.Time, interval uint, maxentries uint) [][]*Entry ***REMOVED***
        var ordered [][]*Entry
        var j int

        if interval > 0 ***REMOVED***
                t := starttime
                d := time.Duration(interval) * time.Millisecond

                for _, e := range entries ***REMOVED***
                        // new interval by date
                        if e.StartedDateTime.Sub(t) >= d ***REMOVED***
                                t = t.Add(d)
                                j++
                        ***REMOVED***
                        if len(ordered) == j ***REMOVED***
                                ordered = append(ordered, []*Entry***REMOVED******REMOVED***)
                        ***REMOVED***
                        // new interval by maxentries value
                        if len(ordered[j]) == int(maxentries) ***REMOVED***
                                ordered = append(ordered, []*Entry***REMOVED******REMOVED***)
                                j++
                        ***REMOVED***
                        ordered[j] = append(ordered[j], e)
                ***REMOVED***
        ***REMOVED***

        return ordered
***REMOVED***
