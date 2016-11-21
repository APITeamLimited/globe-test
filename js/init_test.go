package js

import (
	"fmt"
	"github.com/loadimpact/speedboat/stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetric(t *testing.T) ***REMOVED***
	tpl := `
		import %s from "speedboat/metrics";
		let myMetric = new %s(%s"my_metric", %s);
		export default function() ***REMOVED******REMOVED***;
	`

	types := map[string]stats.MetricType***REMOVED***
		"Counter": stats.Counter,
		"Gauge":   stats.Gauge,
		"Trend":   stats.Trend,
	***REMOVED***

	for s, tp := range types ***REMOVED***
		t.Run("t="+s, func(t *testing.T) ***REMOVED***
			// name: [import, type, arg0]
			imports := map[string][]string***REMOVED***
				"wrapper,direct": []string***REMOVED***
					fmt.Sprintf("***REMOVED*** %s ***REMOVED***", s),
					s,
					"",
				***REMOVED***,
				"wrapper,module": []string***REMOVED***
					"metrics",
					fmt.Sprintf("metrics.%s", s),
					"",
				***REMOVED***,
				"const,direct": []string***REMOVED***
					fmt.Sprintf("***REMOVED*** Metric, %sType ***REMOVED***", s),
					"Metric",
					fmt.Sprintf("%sType, ", s),
				***REMOVED***,
				"const,module": []string***REMOVED***
					"metrics",
					"metrics.Metric",
					fmt.Sprintf("metrics.%sType, ", s),
				***REMOVED***,
			***REMOVED***

			for name, imp := range imports ***REMOVED***
				t.Run("import="+name, func(t *testing.T) ***REMOVED***
					isTimes := map[string]bool***REMOVED***
						"undefined": false,
						"false":     false,
						"true":      true,
					***REMOVED***

					for arg2, isTime := range isTimes ***REMOVED***
						t.Run("isTime="+arg2, func(t *testing.T) ***REMOVED***
							vt := stats.Default
							if isTime ***REMOVED***
								vt = stats.Time
							***REMOVED***

							src := fmt.Sprintf(tpl, imp[0], imp[1], imp[2], arg2)
							r, err := newSnippetRunner(src)
							if !assert.NoError(t, err) ***REMOVED***
								t.Log(src)
								return
							***REMOVED***

							assert.Contains(t, r.Runtime.Metrics, "my_metric")
							m := r.Runtime.Metrics["my_metric"]
							assert.Equal(t, tp, m.Type, "wrong metric type")
							assert.Equal(t, vt, m.Contains, "wrong value type")
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
