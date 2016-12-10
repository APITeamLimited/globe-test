package v2

import (
	"encoding/json"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNullMetricTypeJSON(t *testing.T) ***REMOVED***
	values := map[NullMetricType]string***REMOVED***
		NullMetricType***REMOVED******REMOVED***:                    `null`,
		NullMetricType***REMOVED***stats.Counter, true***REMOVED***: `"counter"`,
		NullMetricType***REMOVED***stats.Gauge, true***REMOVED***:   `"gauge"`,
		NullMetricType***REMOVED***stats.Trend, true***REMOVED***:   `"trend"`,
		NullMetricType***REMOVED***stats.Rate, true***REMOVED***:    `"rate"`,
	***REMOVED***
	t.Run("Marshal", func(t *testing.T) ***REMOVED***
		for mt, val := range values ***REMOVED***
			t.Run(val, func(t *testing.T) ***REMOVED***
				data, err := json.Marshal(mt)
				assert.NoError(t, err)
				assert.Equal(t, val, string(data))
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Unmarshal", func(t *testing.T) ***REMOVED***
		for mt, val := range values ***REMOVED***
			t.Run(val, func(t *testing.T) ***REMOVED***
				var value NullMetricType
				assert.NoError(t, json.Unmarshal([]byte(val), &value))
				assert.Equal(t, mt, value)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestNullValueTypeJSON(t *testing.T) ***REMOVED***
	values := map[NullValueType]string***REMOVED***
		NullValueType***REMOVED******REMOVED***:                    `null`,
		NullValueType***REMOVED***stats.Default, true***REMOVED***: `"default"`,
		NullValueType***REMOVED***stats.Time, true***REMOVED***:    `"time"`,
	***REMOVED***
	t.Run("Marshal", func(t *testing.T) ***REMOVED***
		for mt, val := range values ***REMOVED***
			t.Run(val, func(t *testing.T) ***REMOVED***
				data, err := json.Marshal(mt)
				assert.NoError(t, err)
				assert.Equal(t, val, string(data))
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Unmarshal", func(t *testing.T) ***REMOVED***
		for mt, val := range values ***REMOVED***
			t.Run(val, func(t *testing.T) ***REMOVED***
				var value NullValueType
				assert.NoError(t, json.Unmarshal([]byte(val), &value))
				assert.Equal(t, mt, value)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***
