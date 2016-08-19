package js

import (
	log "github.com/Sirupsen/logrus"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	// "math"
	"os"
	// "strconv"
	"testing"
	"time"
)

func TestNewVU(t *testing.T) ***REMOVED***
	r := New("script", "1+1")
	_, err := r.NewVU()
	assert.NoError(t, err)
***REMOVED***

func TestNewVUInvalidJS(t *testing.T) ***REMOVED***
	r := New("script", "aiugbauibeuifa")
	_, err := r.NewVU()
	assert.NoError(t, err)
***REMOVED***

func TestReconfigure(t *testing.T) ***REMOVED***
	r := New("script", "1+1")
	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	vu.ID = 100
	vu.Iteration = 100

	vu.Reconfigure(1)
	assert.Equal(t, int64(1), vu.ID)
	assert.Equal(t, int64(0), vu.Iteration)
***REMOVED***

func TestRunOnceIncreasesIterations(t *testing.T) ***REMOVED***
	r := New("script", "1+1")
	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	assert.Equal(t, int64(0), vu.Iteration)
	vu.RunOnce(context.Background())
	assert.Equal(t, int64(1), vu.Iteration)
***REMOVED***

func TestRunOnceInvalidJS(t *testing.T) ***REMOVED***
	r := New("script", "diyfsybfbub")
	vu, err := r.NewVU()
	assert.NoError(t, err)
	err = vu.RunOnce(context.Background())
	assert.Error(t, err)
***REMOVED***

func TestAPILogDebug(t *testing.T) ***REMOVED***
	r := New("script", `$log.debug("test");`)
	logger, hook := logtest.NewNullLogger()
	logger.Level = log.DebugLevel
	r.logger = logger

	vu, _ := r.NewVU()
	assert.NoError(t, vu.RunOnce(context.Background()))

	e := hook.LastEntry()
	assert.NotNil(t, e)
	assert.Equal(t, log.DebugLevel, e.Level)
	assert.Equal(t, "test", e.Message)
	assert.Len(t, e.Data, 0)
***REMOVED***

func TestAPILogInfo(t *testing.T) ***REMOVED***
	r := New("script", `$log.info("test");`)
	logger, hook := logtest.NewNullLogger()
	r.logger = logger

	vu, _ := r.NewVU()
	assert.NoError(t, vu.RunOnce(context.Background()))

	e := hook.LastEntry()
	assert.NotNil(t, e)
	assert.Equal(t, log.InfoLevel, e.Level)
	assert.Equal(t, "test", e.Message)
	assert.Len(t, e.Data, 0)
***REMOVED***

func TestAPILogWarn(t *testing.T) ***REMOVED***
	r := New("script", `$log.warn("test");`)
	logger, hook := logtest.NewNullLogger()
	r.logger = logger

	vu, _ := r.NewVU()
	assert.NoError(t, vu.RunOnce(context.Background()))

	e := hook.LastEntry()
	assert.NotNil(t, e)
	assert.Equal(t, log.WarnLevel, e.Level)
	assert.Equal(t, "test", e.Message)
	assert.Len(t, e.Data, 0)
***REMOVED***

func TestAPILogError(t *testing.T) ***REMOVED***
	r := New("script", `$log.error("test");`)
	logger, hook := logtest.NewNullLogger()
	r.logger = logger

	vu, _ := r.NewVU()
	assert.NoError(t, vu.RunOnce(context.Background()))

	e := hook.LastEntry()
	assert.NotNil(t, e)
	assert.Equal(t, log.ErrorLevel, e.Level)
	assert.Equal(t, "test", e.Message)
	assert.Len(t, e.Data, 0)
***REMOVED***

func TestAPILogWithData(t *testing.T) ***REMOVED***
	r := New("script", `$log.info("test", ***REMOVED*** a: 'hi', b: 123 ***REMOVED***);`)
	logger, hook := logtest.NewNullLogger()
	r.logger = logger

	vu, _ := r.NewVU()
	assert.NoError(t, vu.RunOnce(context.Background()))

	e := hook.LastEntry()
	assert.NotNil(t, e)
	assert.Equal(t, log.InfoLevel, e.Level)
	assert.Equal(t, "test", e.Message)
	assert.Equal(t, log.Fields***REMOVED***"a": "hi", "b": int64(123)***REMOVED***, e.Data)
***REMOVED***

func TestAPIVUSleep1s(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `$vu.sleep(1);`)
	vu, _ := r.NewVU()

	startTime := time.Now()
	err := vu.RunOnce(context.Background())
	duration := time.Since(startTime)

	assert.NoError(t, err)

	// Allow 50ms margin for call overhead
	target := 1 * time.Second
	if duration < target || duration > target+(50*time.Millisecond) ***REMOVED***
		t.Fatalf("Incorrect sleep duration: %s", duration)
	***REMOVED***
***REMOVED***

func TestAPIVUSleep01s(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `$vu.sleep(0.1);`)
	vu, _ := r.NewVU()

	startTime := time.Now()
	err := vu.RunOnce(context.Background())
	duration := time.Since(startTime)

	assert.NoError(t, err)

	// Allow 50ms margin for call overhead
	target := 100 * time.Millisecond
	if duration < target || duration > target+(50*time.Millisecond) ***REMOVED***
		t.Fatalf("Incorrect sleep duration: %s", duration)
	***REMOVED***
***REMOVED***

func TestAPIVUID(t *testing.T) ***REMOVED***
	r := New("script", `if ($vu.id() !== 100) ***REMOVED*** throw new Error("invalid ID"); ***REMOVED***`)
	vu, _ := r.NewVU()
	vu.Reconfigure(100)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIVUIteration(t *testing.T) ***REMOVED***
	r := New("script", `if ($vu.iteration() !== 1) ***REMOVED*** throw new Error("invalid iteration"); ***REMOVED***`)
	vu, _ := r.NewVU()
	vu.Reconfigure(100)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPITestEnv(t *testing.T) ***REMOVED***
	os.Setenv("TEST_VAR", "hi")
	r := New("script", `if ($test.env("TEST_VAR") !== "hi") ***REMOVED*** throw new Error("assertion failed"); ***REMOVED***`)
	vu, _ := r.NewVU()
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPITestEnvUndefined(t *testing.T) ***REMOVED***
	os.Unsetenv("NOT_SET_VAR") // Just in case...
	r := New("script", `if ($test.env("NOT_SET_VAR") !== undefined) ***REMOVED*** throw new Error("assertion failed"); ***REMOVED***`)
	vu, _ := r.NewVU()
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPITestAbort(t *testing.T) ***REMOVED***
	r := New("script", `$test.abort();`)
	vu, _ := r.NewVU()
	assert.Panics(t, func() ***REMOVED*** vu.RunOnce(context.Background()) ***REMOVED***)
***REMOVED***

// func TestAPIHTTPSetMaxConnsPerHost(t *testing.T) ***REMOVED***
// 	r := New("script", `$http.setMaxConnsPerHost(100);`)
// 	vu, _ := r.NewVU()
// 	assert.NoError(t, vu.RunOnce(context.Background()))
// 	assert.Equal(t, 100, vu.(*VU).Client.MaxConnsPerHost)
// ***REMOVED***

// func TestAPIHTTPSetMaxConnsPerHostOverflow(t *testing.T) ***REMOVED***
// 	r := New("script", `$http.setMaxConnsPerHost(`+strconv.FormatInt(math.MaxInt64, 10)+`);`)
// 	vu, _ := r.NewVU()
// 	assert.NoError(t, vu.RunOnce(context.Background()))
// 	assert.Equal(t, math.MaxInt32, vu.(*VU).Client.MaxConnsPerHost)
// ***REMOVED***

// func TestAPIHTTPSetMaxConnsPerHostZero(t *testing.T) ***REMOVED***
// 	r := New("script", `$http.setMaxConnsPerHost(0);`)
// 	vu, _ := r.NewVU()
// 	assert.Error(t, vu.RunOnce(context.Background()))
// ***REMOVED***

// func TestAPIHTTPSetMaxConnsPerHostNegative(t *testing.T) ***REMOVED***
// 	r := New("script", `$http.setMaxConnsPerHost(-1);`)
// 	vu, _ := r.NewVU()
// 	assert.Error(t, vu.RunOnce(context.Background()))
// ***REMOVED***

// func TestAPIHTTPSetMaxConnsPerHostInvalid(t *testing.T) ***REMOVED***
// 	r := New("script", `$http.setMaxConnsPerHost("qwerty");`)
// 	vu, _ := r.NewVU()
// 	assert.Error(t, vu.RunOnce(context.Background()))
// ***REMOVED***

func TestAPIHTTPRequestReportsStats(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", "$http.get('http://httpbin.org/get');")
	vu, err := r.NewVU()
	assert.NoError(t, err)
	err = vu.RunOnce(context.Background())
	assert.NoError(t, err)

	mRequestsFound := false
	for _, p := range vu.(*VU).Collector.Batch ***REMOVED***
		switch p.Stat ***REMOVED***
		case &mRequests:
			mRequestsFound = true
			assert.Contains(t, p.Tags, "url")
			assert.Contains(t, p.Tags, "method")
			assert.Contains(t, p.Tags, "status")
			assert.Contains(t, p.Values, "duration")
		case &mErrors:
			assert.Fail(t, "Errors found")
		***REMOVED***
	***REMOVED***
	assert.True(t, mRequestsFound)
***REMOVED***

func TestAPIHTTPRequestErrorReportsStats(t *testing.T) ***REMOVED***
	r := New("script", "$http.get('http://255.255.255.255/');")
	vu, err := r.NewVU()
	assert.NoError(t, err)
	err = vu.RunOnce(context.Background())
	assert.Error(t, err)

	mRequestsFound := false
	mErrorsFound := false
	for _, p := range vu.(*VU).Collector.Batch ***REMOVED***
		switch p.Stat ***REMOVED***
		case &mRequests:
			mRequestsFound = true
		case &mErrors:
			mErrorsFound = true
			assert.Contains(t, p.Tags, "url")
			assert.Contains(t, p.Tags, "method")
			assert.Contains(t, p.Values, "value")
		***REMOVED***
	***REMOVED***
	assert.False(t, mRequestsFound)
	assert.True(t, mErrorsFound)
***REMOVED***

func TestAPIHTTPRequestQuietReportsNoStats(t *testing.T) ***REMOVED***
	r := New("script", "$http.get('http://255.255.255.255/', null, ***REMOVED*** quiet: true ***REMOVED***);")
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.Error(t, vu.RunOnce(context.Background()))
	assert.Len(t, vu.(*VU).Collector.Batch, 0)
***REMOVED***

func TestAPIHTTPRequestGET(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.get("http://httpbin.org/get")
	if (res.status !== 200) ***REMOVED***
		throw new Error("invalid status: " + res.status);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestGETArgs(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	data = $http.get("http://httpbin.org/get", ***REMOVED***a: 'b', b: 2***REMOVED***).json()
	if (data.args.a !== 'b') ***REMOVED***
		throw new Error("invalid args.a: " + data.args.a);
	***REMOVED***
	if (data.args.b !== '2') ***REMOVED***
		throw new Error("invalid args.b: " + data.args.b);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestGETHeaders(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	data = $http.get("http://httpbin.org/get", null, ***REMOVED*** headers: ***REMOVED*** 'X-Test': 'hi' ***REMOVED*** ***REMOVED***).json()
	if (data.headers['X-Test'] !== 'hi') ***REMOVED***
		throw new Error("invalid X-Test header: " + data.headers['X-Test'])
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestGETRedirect(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.get("http://httpbin.org/redirect/6");
	if (res.status !== 302) ***REMOVED***
		throw new Error("invalid response code: " + res.status);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestGETRedirectFollow(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.get("http://httpbin.org/redirect/6", null, ***REMOVED*** follow: true ***REMOVED***);
	if (res.status !== 200) ***REMOVED***
		throw new Error("invalid response code: " + res.status);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestGETRedirectFollowTooMany(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	$http.get("http://httpbin.org/redirect/15", null, ***REMOVED*** follow: true ***REMOVED***);
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.Error(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestHEAD(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.head("http://httpbin.org/get")
	if (res.status !== 200) ***REMOVED***
		throw new Error("invalid status: " + res.status);
	***REMOVED***
	if (res.body !== "") ***REMOVED***
		throw new Error("body not empty")
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestHEADWithArgsDoesntStickThemInTheBodyAndFail(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.head("http://httpbin.org/get", ***REMOVED*** a: 'b' ***REMOVED***)
	if (res.status !== 200) ***REMOVED***
		throw new Error("invalid status: " + res.status);
	***REMOVED***
	if (res.body !== "") ***REMOVED***
		throw new Error("body not empty")
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestPOST(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.post("http://httpbin.org/post")
	if (res.status !== 200) ***REMOVED***
		throw new Error("invalid status: " + res.status);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestPOSTArgs(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	data = $http.post("http://httpbin.org/post", ***REMOVED*** a: 'b' ***REMOVED***).json()
	if (data.form.a !== 'b') ***REMOVED***
		throw new Error("invalid form.a: " + data.form.a);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestPOSTBody(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	data = $http.post("http://httpbin.org/post", 'a=b').json()
	if (data.data !== 'a=b') ***REMOVED***
		throw new Error("invalid data: " + data.data);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestPUT(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.put("http://httpbin.org/put")
	if (res.status !== 200) ***REMOVED***
		throw new Error("invalid status: " + res.status);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestPATCH(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.patch("http://httpbin.org/patch")
	if (res.status !== 200) ***REMOVED***
		throw new Error("invalid status: " + res.status);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestDELETE(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.delete("http://httpbin.org/delete")
	if (res.status !== 200) ***REMOVED***
		throw new Error("invalid status: " + res.status);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***

func TestAPIHTTPRequestOPTIONS(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	r := New("script", `
	res = $http.options("http://httpbin.org/")
	if (res.status !== 200) ***REMOVED***
		throw new Error("invalid status: " + res.status);
	***REMOVED***
	if (res.body !== "") ***REMOVED***
		throw new Error("non-empty body: " + res.body);
	***REMOVED***
	`)
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.NoError(t, vu.RunOnce(context.Background()))
***REMOVED***
