package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetric(t *testing.T) ***REMOVED***
	t.Parallel()
	testdata := map[string]struct ***REMOVED***
		Type     MetricType
		SinkType Sink
	***REMOVED******REMOVED***
		"Counter": ***REMOVED***Counter, &CounterSink***REMOVED******REMOVED******REMOVED***,
		"Gauge":   ***REMOVED***Gauge, &GaugeSink***REMOVED******REMOVED******REMOVED***,
		"Trend":   ***REMOVED***Trend, &TrendSink***REMOVED******REMOVED******REMOVED***,
		"Rate":    ***REMOVED***Rate, &RateSink***REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	for name, data := range testdata ***REMOVED***
		name, data := name, data
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			m := newMetric("my_metric", data.Type)
			assert.Equal(t, "my_metric", m.Name)
			assert.IsType(t, data.SinkType, m.Sink)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestAddSubmetric(t *testing.T) ***REMOVED***
	t.Parallel()
	testdata := map[string]struct ***REMOVED***
		err  bool
		tags map[string]string
	***REMOVED******REMOVED***
		"":                        ***REMOVED***true, nil***REMOVED***,
		"  ":                      ***REMOVED***true, nil***REMOVED***,
		"a":                       ***REMOVED***false, map[string]string***REMOVED***"a": ""***REMOVED******REMOVED***,
		"a:1":                     ***REMOVED***false, map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		" a : 1 ":                 ***REMOVED***false, map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		"a,b":                     ***REMOVED***false, map[string]string***REMOVED***"a": "", "b": ""***REMOVED******REMOVED***,
		` a:"",b: ''`:             ***REMOVED***false, map[string]string***REMOVED***"a": "", "b": ""***REMOVED******REMOVED***,
		`a:1,b:2`:                 ***REMOVED***false, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		` a : 1, b : 2 `:          ***REMOVED***false, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		`a : '1' , b : "2"`:       ***REMOVED***false, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		`" a" : ' 1' , b : "2 " `: ***REMOVED***false, map[string]string***REMOVED***" a": " 1", "b": "2 "***REMOVED******REMOVED***, //nolint:gocritic
	***REMOVED***

	for name, expected := range testdata ***REMOVED***
		name, expected := name, expected
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			m := newMetric("metric", Trend)
			sm, err := m.AddSubmetric(name)
			if expected.err ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)
			require.NotNil(t, sm)
			assert.EqualValues(t, expected.tags, sm.Tags.tags)
		***REMOVED***)
	***REMOVED***
***REMOVED***
