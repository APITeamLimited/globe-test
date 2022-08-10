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

	"go.k6.io/k6/api/common"
	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/minirunner"
	"go.k6.io/k6/metrics"
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
					newLogger(l, http.HandlerFunc(testHTTPHandler))(rw, r)

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
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	registry := metrics.NewRegistry()
	testState := &lib.TestRunState***REMOVED***
		TestPreInitState: &lib.TestPreInitState***REMOVED***
			Logger:         logger,
			Registry:       registry,
			BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
		***REMOVED***,
		Options: lib.Options***REMOVED******REMOVED***,
		Runner:  &minirunner.MiniRunner***REMOVED******REMOVED***,
	***REMOVED***

	execScheduler, err := local.NewExecutionScheduler(testState)
	require.NoError(t, err)
	engine, err := core.NewEngine(testState, execScheduler, nil)
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com/", nil)
	withEngine(engine, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		assert.Equal(t, engine, common.GetEngine(r.Context()))
	***REMOVED***))(rw, r)
***REMOVED***

func TestPing(t *testing.T) ***REMOVED***
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	mux := newHandler(logger)

	rw := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ping", nil)
	mux.ServeHTTP(rw, r)

	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, []byte***REMOVED***'o', 'k'***REMOVED***, rw.Body.Bytes())
***REMOVED***
