package js

import (
	log "github.com/Sirupsen/logrus"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
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

func TestAPIVUSleep1s(t *testing.T) ***REMOVED***
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

func TestAPIHTTPRequestReportsStats(t *testing.T) ***REMOVED***
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
			assert.Contains(t, p.Tags, "url")
			assert.Contains(t, p.Tags, "method")
			assert.Contains(t, p.Tags, "status")
			assert.Contains(t, p.Values, "duration")
		case &mErrors:
			mErrorsFound = true
			assert.Contains(t, p.Tags, "url")
			assert.Contains(t, p.Tags, "method")
			assert.Contains(t, p.Tags, "status")
			assert.Contains(t, p.Values, "value")
		***REMOVED***
	***REMOVED***
	assert.True(t, mRequestsFound)
	assert.True(t, mErrorsFound)
***REMOVED***
