package simple

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestNew(t *testing.T) ***REMOVED***
	r := New("http://example.com/")
	assert.Equal(t, "http://example.com/", r.URL)
***REMOVED***

func TestNewVU(t *testing.T) ***REMOVED***
	r := New("http://example.com/")
	vu, err := r.NewVU()
	assert.NoError(t, err)
	assert.IsType(t, &VU***REMOVED******REMOVED***, vu)
***REMOVED***

func TestReconfigure(t *testing.T) ***REMOVED***
	r := New("http://example.com/")

	vu, err := r.NewVU()
	assert.NoError(t, err)

	err = vu.Reconfigure(12345)
	assert.NoError(t, err)
***REMOVED***

func TestRunOnceReportsStats(t *testing.T) ***REMOVED***
	r := New("http://255.255.255.255/")
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
