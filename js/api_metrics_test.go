package js

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricAdd(t *testing.T) ***REMOVED***
	values := map[string]float64***REMOVED***
		`1234`:   1234.0,
		`1234.5`: 1234.5,
		`true`:   1.0,
		`false`:  0.0,
	***REMOVED***

	for jsV, v := range values ***REMOVED***
		t.Run("v="+jsV, func(t *testing.T) ***REMOVED***
			tags := map[string]map[string]string***REMOVED***
				`undefined`:     map[string]string***REMOVED******REMOVED***,
				`***REMOVED***tag:"value"***REMOVED***`: map[string]string***REMOVED***"tag": "value"***REMOVED***,
				`***REMOVED***tag:1234***REMOVED***`:    map[string]string***REMOVED***"tag": "1234"***REMOVED***,
				`***REMOVED***tag:1234.5***REMOVED***`:  map[string]string***REMOVED***"tag": "1234.5"***REMOVED***,
			***REMOVED***

			for jsT, t_ := range tags ***REMOVED***
				t.Run("t="+jsT, func(t *testing.T) ***REMOVED***
					r, err := newSnippetRunner(fmt.Sprintf(`
						import ***REMOVED*** _assert ***REMOVED*** from "speedboat";
						import ***REMOVED*** Counter ***REMOVED*** from "speedboat/metrics";
						let myMetric = new Counter("my_metric");
						export default function() ***REMOVED***
							let v = %s;
							let t = %s;
							_assert(myMetric.add(v, t) === v);
						***REMOVED***
					`, jsV, jsT))

					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***

					vu, err := r.NewVU()
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***

					samples, err := vu.RunOnce(context.Background())
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***

					assert.Len(t, samples, 1)
					s := samples[0]
					assert.Equal(t, r.Runtime.Metrics["my_metric"], s.Metric)
					assert.Equal(t, v, s.Value)
					assert.EqualValues(t, t_, s.Tags)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
