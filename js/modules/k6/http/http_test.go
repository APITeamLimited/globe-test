/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package http

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/stretchr/testify/assert"
)

func TestTagURL(t *testing.T) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	rt.Set("http", common.Bind(rt, New(), nil))

	testdata := map[string]URLTag***REMOVED***
		`http://httpbin.org/anything/`:               ***REMOVED***URL: "http://httpbin.org/anything/", Name: "http://httpbin.org/anything/"***REMOVED***,
		`http://httpbin.org/anything/$***REMOVED***1+1***REMOVED***`:         ***REMOVED***URL: "http://httpbin.org/anything/2", Name: "http://httpbin.org/anything/$***REMOVED******REMOVED***"***REMOVED***,
		`http://httpbin.org/anything/$***REMOVED***1+1***REMOVED***/`:        ***REMOVED***URL: "http://httpbin.org/anything/2/", Name: "http://httpbin.org/anything/$***REMOVED******REMOVED***/"***REMOVED***,
		`http://httpbin.org/anything/$***REMOVED***1+1***REMOVED***/$***REMOVED***1+2***REMOVED***`:  ***REMOVED***URL: "http://httpbin.org/anything/2/3", Name: "http://httpbin.org/anything/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***"***REMOVED***,
		`http://httpbin.org/anything/$***REMOVED***1+1***REMOVED***/$***REMOVED***1+2***REMOVED***/`: ***REMOVED***URL: "http://httpbin.org/anything/2/3/", Name: "http://httpbin.org/anything/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***/"***REMOVED***,
	***REMOVED***
	for expr, tag := range testdata ***REMOVED***
		t.Run("expr="+expr, func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, "http.url`"+expr+"`")
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, tag, v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
