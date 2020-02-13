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

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"

	"github.com/loadimpact/k6/api/common"
	"github.com/loadimpact/k6/core"
	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils/minirunner"
)

func testHTTPHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if _, err := fmt.Fprint(rw, "ok"); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

func TestLogger(t *testing.T) ***REMOVED***
	for _, method := range []string***REMOVED***"GET", "POST", "PUT", "PATCH"***REMOVED*** ***REMOVED***
		t.Run("method="+method, func(t *testing.T) ***REMOVED***
			for _, path := range []string***REMOVED***"/", "/test", "/test/path"***REMOVED*** ***REMOVED***
				t.Run("path="+path, func(t *testing.T) ***REMOVED***
					rw := httptest.NewRecorder()
					r := httptest.NewRequest(method, "http://example.com"+path, nil)

					l, hook := logtest.NewNullLogger()
					l.Level = logrus.DebugLevel
					NewLogger(l)(negroni.NewResponseWriter(rw), r, testHTTPHandler)

					res := rw.Result()
					assert.Equal(t, http.StatusOK, res.StatusCode)
					assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))

					if !assert.Len(t, hook.Entries, 1) ***REMOVED***
						return
					***REMOVED***

					e := hook.LastEntry()
					assert.Equal(t, logrus.DebugLevel, e.Level)
					assert.Equal(t, fmt.Sprintf("%s %s", method, path), e.Message)
					assert.Equal(t, http.StatusOK, e.Data["status"])
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestWithEngine(t *testing.T) ***REMOVED***
	execScheduler, err := local.NewExecutionScheduler(&minirunner.MiniRunner***REMOVED******REMOVED***, logrus.StandardLogger())
	require.NoError(t, err)
	engine, err := core.NewEngine(execScheduler, lib.Options***REMOVED******REMOVED***, logrus.StandardLogger())
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com/", nil)
	WithEngine(engine)(rw, r, func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		assert.Equal(t, engine, common.GetEngine(r.Context()))
	***REMOVED***)
***REMOVED***

func TestPing(t *testing.T) ***REMOVED***
	mux := NewHandler()

	rw := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ping", nil)
	mux.ServeHTTP(rw, r)

	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, []byte***REMOVED***'o', 'k'***REMOVED***, rw.Body.Bytes())
***REMOVED***
