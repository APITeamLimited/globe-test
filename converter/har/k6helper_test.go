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
        "fmt"
        "testing"
        "time"

        "github.com/loadimpact/k6/js"
        "github.com/loadimpact/k6/lib"
        "github.com/spf13/afero"
)

func TestBuildK6Cookies(t *testing.T) ***REMOVED***
        var cookies = []struct ***REMOVED***
                values   []Cookie
                expected string
        ***REMOVED******REMOVED***
                ***REMOVED***[]Cookie***REMOVED******REMOVED***Name: "a", Value: "b"***REMOVED******REMOVED***, "a=b"***REMOVED***,
                ***REMOVED***[]Cookie***REMOVED******REMOVED***Name: "a", Value: "b"***REMOVED***, ***REMOVED***Name: "c", Value: "d"***REMOVED******REMOVED***, "a=b; c=d"***REMOVED***,
        ***REMOVED***

        for _, pair := range cookies ***REMOVED***
                v := BuildK6CookiesValues(pair.values)
                if v != pair.expected ***REMOVED***
                        t.Errorf("BuildK6Cookies(%v): expected %v, actual %v", pair.values, pair.expected, v)
                ***REMOVED***
        ***REMOVED***
***REMOVED***

func TestBuildK6Headers(t *testing.T) ***REMOVED***
        var headers = []struct ***REMOVED***
                values   []Header
                expected string
        ***REMOVED******REMOVED***
                ***REMOVED***[]Header***REMOVED******REMOVED***"name", "1"***REMOVED***, ***REMOVED***"name", "2"***REMOVED******REMOVED***, "\"headers\" : ***REMOVED*** \"name\" : \"1\" ***REMOVED***"***REMOVED***,
                ***REMOVED***[]Header***REMOVED******REMOVED***"name", "1"***REMOVED***, ***REMOVED***"Name", "2"***REMOVED******REMOVED***, "\"headers\" : ***REMOVED*** \"name\" : \"1\" ***REMOVED***"***REMOVED***,
                ***REMOVED***[]Header***REMOVED******REMOVED***"Name", "1"***REMOVED***, ***REMOVED***"name", "2"***REMOVED******REMOVED***, "\"headers\" : ***REMOVED*** \"Name\" : \"1\" ***REMOVED***"***REMOVED***,
                ***REMOVED***[]Header***REMOVED******REMOVED***"name", "value"***REMOVED***, ***REMOVED***"name2", "value2"***REMOVED******REMOVED***, "\"headers\" : ***REMOVED*** \"name\" : \"value\", \"name2\" : \"value2\" ***REMOVED***"***REMOVED***,
                ***REMOVED***[]Header***REMOVED******REMOVED***"accept-language", "es-ES,es;q=0.8"***REMOVED******REMOVED***, "\"headers\" : ***REMOVED*** \"accept-language\" : \"es-ES,es;q=0.8\" ***REMOVED***"***REMOVED***,
                ***REMOVED***[]Header***REMOVED******REMOVED***":host", "localhost"***REMOVED******REMOVED***, "\"headers\" : ***REMOVED***  ***REMOVED***"***REMOVED***, // avoid SPDY’s colon headers
        ***REMOVED***

        for _, pair := range headers ***REMOVED***
                v := BuildK6Headers(pair.values)
                if v != pair.expected ***REMOVED***
                        t.Errorf("BuildK6Headers(%v): expected %v, actual %v", pair.values, pair.expected, v)
                ***REMOVED***
        ***REMOVED***
***REMOVED***

func TestBuildK6Request(t *testing.T) ***REMOVED***
        v, err := BuildK6Request("get", "http://www.google.es", "", []Header***REMOVED******REMOVED***"accept-language", "es-ES,es;q=0.8"***REMOVED******REMOVED***, []Cookie***REMOVED******REMOVED***Name: "a", Value: "b"***REMOVED******REMOVED***)

        if err != nil ***REMOVED***
                t.Error(err)
        ***REMOVED***

        _, err = js.New(&lib.SourceData***REMOVED***
                Filename: "/script.js",
                Data:     []byte(fmt.Sprintf("export default function() ***REMOVED*** %v ***REMOVED***", v)),
        ***REMOVED***, afero.NewMemMapFs())

        if err != nil ***REMOVED***
                t.Error(err)
        ***REMOVED***
***REMOVED***

func TestBuildK6RequestObject(t *testing.T) ***REMOVED***
        v, err := BuildK6RequestObject("get", "http://www.google.es", "", []Header***REMOVED******REMOVED***"accept-language", "es-ES,es;q=0.8"***REMOVED******REMOVED***, []Cookie***REMOVED******REMOVED***Name: "a", Value: "b"***REMOVED******REMOVED***)
        if err != nil ***REMOVED***
                t.Error(err)
        ***REMOVED***
        _, err = js.New(&lib.SourceData***REMOVED***
                Filename: "/script.js",
                Data:     []byte(fmt.Sprintf("export default function() ***REMOVED*** res = http.batch([%v]); ***REMOVED***", v)),
        ***REMOVED***, afero.NewMemMapFs())

        if err != nil ***REMOVED***
                t.Error(err)
        ***REMOVED***
***REMOVED***

func TestIsAllowedURL(t *testing.T) ***REMOVED***
        var allowed = []struct ***REMOVED***
                url      string
                only     []string
                skip     []string
                expected bool
        ***REMOVED******REMOVED***
                ***REMOVED***"http://www.google.com/", []string***REMOVED******REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
                ***REMOVED***"http://www.google.com/", []string***REMOVED***"google.com"***REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
                ***REMOVED***"https://www.google.com/", []string***REMOVED***"google.com"***REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
                ***REMOVED***"https://www.google.com/", []string***REMOVED***"http://"***REMOVED***, []string***REMOVED******REMOVED***, false***REMOVED***,
                ***REMOVED***"http://www.google.com/?hl=en", []string***REMOVED***"http://www.google.com"***REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
                ***REMOVED***"http://www.google.com/?hl=en", []string***REMOVED***"google.com", "google.co.uk"***REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
                ***REMOVED***"http://www.google.com/?hl=en", []string***REMOVED******REMOVED***, []string***REMOVED***"google.com"***REMOVED***, false***REMOVED***,
                ***REMOVED***"http://www.google.com/?hl=en", []string***REMOVED******REMOVED***, []string***REMOVED***"google.co.uk"***REMOVED***, true***REMOVED***,
        ***REMOVED***

        for _, s := range allowed ***REMOVED***
                v := IsAllowedURL(s.url, s.only, s.skip)
                if v != s.expected ***REMOVED***
                        t.Errorf("IsAllowedURL(%v, %v, %v): expected %v, actual %v", s.url, s.only, s.skip, s.expected, v)
                ***REMOVED***
        ***REMOVED***
***REMOVED***

func TestGroupHarEntriesByIntervals(t *testing.T) ***REMOVED***
        // max number of requests in a batch statement
        const maxentries uint = 5

        t1 := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)

        entries := []*Entry***REMOVED******REMOVED***

        // 10 time entries with increments of 100ms (from 0 a 1000ms)
        for i := 0; i < 10; i++ ***REMOVED***
                entries = append(entries, &Entry***REMOVED***StartedDateTime: t1.Add(time.Duration(100*i) * time.Millisecond)***REMOVED***)
        ***REMOVED***

        splitValues := []struct ***REMOVED***
                diff, groups uint
        ***REMOVED******REMOVED***
                ***REMOVED***0, 0***REMOVED***,
                ***REMOVED***1000, 2***REMOVED***,
                ***REMOVED***100, 10***REMOVED***,
                ***REMOVED***500, 2***REMOVED***,
                ***REMOVED***800, 3***REMOVED***, // group with 5 entries, group with 3 entries, group with 2 entries
        ***REMOVED***

        for _, v := range splitValues ***REMOVED***
                result := groupHarEntriesByIntervals(entries, t1, v.diff, maxentries)
                if len(result) != int(v.groups) ***REMOVED***
                        t.Errorf("groupHarEntriesByIntervals(%v, %v, %v, %v) Expected %v, actual %v", entries, t1, v.diff, maxentries, v.groups, len(result))
                ***REMOVED***
        ***REMOVED***
***REMOVED***
