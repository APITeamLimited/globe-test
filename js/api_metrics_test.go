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
						import ***REMOVED*** _assert ***REMOVED*** from "k6";
						import ***REMOVED*** Counter ***REMOVED*** from "k6/metrics";
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
