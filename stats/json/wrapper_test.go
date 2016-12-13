package json

import (
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrapWithNilArg(t *testing.T) ***REMOVED***
	out := Wrap(nil)
	assert.Equal(t, out, (*Envelope)(nil))
***REMOVED***

func TestWrapWithUnusedType(t *testing.T) ***REMOVED***
	out := Wrap(JSONSample***REMOVED******REMOVED***)
	assert.Equal(t, out, (*Envelope)(nil))
***REMOVED***

func TestWrapWithSample(t *testing.T) ***REMOVED***
	out := Wrap(stats.Sample***REMOVED***
		Metric: &stats.Metric***REMOVED******REMOVED***,
	***REMOVED***)
	assert.NotEqual(t, out, (*Envelope)(nil))
***REMOVED***

func TestWrapWithMetricPointer(t *testing.T) ***REMOVED***
	out := Wrap(&stats.Metric***REMOVED******REMOVED***)
	assert.NotEqual(t, out, (*Envelope)(nil))
***REMOVED***

func TestWrapWithMetric(t *testing.T) ***REMOVED***
	out := Wrap(stats.Metric***REMOVED******REMOVED***)
	assert.Equal(t, out, (*Envelope)(nil))
***REMOVED***
