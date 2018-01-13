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

	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
)

func TestBuildK6Headers(t *testing.T) ***REMOVED***
	var headers = []struct ***REMOVED***
		values   []Header
		expected []string
	***REMOVED******REMOVED***
		***REMOVED***[]Header***REMOVED******REMOVED***"name", "1"***REMOVED***, ***REMOVED***"name", "2"***REMOVED******REMOVED***, []string***REMOVED***`"name": "1"`***REMOVED******REMOVED***,
		***REMOVED***[]Header***REMOVED******REMOVED***"name", "1"***REMOVED***, ***REMOVED***"name2", "2"***REMOVED******REMOVED***, []string***REMOVED***`"name": "1"`, `"name2": "2"`***REMOVED******REMOVED***,
		***REMOVED***[]Header***REMOVED******REMOVED***":host", "localhost"***REMOVED******REMOVED***, []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	for _, pair := range headers ***REMOVED***
		v := buildK6Headers(pair.values)
		if len(v) != len(pair.expected) ***REMOVED***
			t.Errorf("buildK6Headers(%v): expected %v, actual %v", pair.values, pair.expected, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBuildK6RequestObject(t *testing.T) ***REMOVED***
	req := &Request***REMOVED***
		Method:  "get",
		URL:     "http://www.google.es",
		Headers: []Header***REMOVED******REMOVED***"accept-language", "es-ES,es;q=0.8"***REMOVED******REMOVED***,
		Cookies: []Cookie***REMOVED******REMOVED***Name: "a", Value: "b"***REMOVED******REMOVED***,
	***REMOVED***
	v, err := buildK6RequestObject(req)
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
