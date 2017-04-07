package metrics

import (
	"context"
	"errors"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js2/common"
	"github.com/loadimpact/k6/stats"
)

type Metric struct ***REMOVED***
	metric *stats.Metric
***REMOVED***

func newMetric(ctxPtr *context.Context, name string, t stats.MetricType, isTime []bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if common.GetState(*ctxPtr) != nil ***REMOVED***
		return nil, errors.New("Metrics must be declared in the init context")
	***REMOVED***

	valueType := stats.Default
	if len(isTime) > 0 && isTime[0] ***REMOVED***
		valueType = stats.Time
	***REMOVED***

	rt := common.GetRuntime(*ctxPtr)
	return common.Bind(rt, Metric***REMOVED***stats.New(name, t, valueType)***REMOVED***, ctxPtr), nil
***REMOVED***

func (m Metric) Add(ctx context.Context, v goja.Value, addTags ...map[string]string) ***REMOVED***
	state := common.GetState(ctx)

	tags := map[string]string***REMOVED***
		"group": state.Group.Path,
	***REMOVED***
	for _, ts := range addTags ***REMOVED***
		for k, v := range ts ***REMOVED***
			tags[k] = v
		***REMOVED***
	***REMOVED***

	vfloat := v.ToFloat()
	if vfloat == 0 && v.ToBoolean() ***REMOVED***
		vfloat = 1.0
	***REMOVED***

	state.Samples = append(state.Samples,
		stats.Sample***REMOVED***Time: time.Now(), Metric: m.metric, Value: vfloat, Tags: tags***REMOVED***,
	)
***REMOVED***

type Metrics struct***REMOVED******REMOVED***

func (*Metrics) XCounter(ctx *context.Context, name string, isTime ...bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return newMetric(ctx, name, stats.Counter, isTime)
***REMOVED***

func (*Metrics) XGauge(ctx *context.Context, name string, isTime ...bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return newMetric(ctx, name, stats.Gauge, isTime)
***REMOVED***

func (*Metrics) XTrend(ctx *context.Context, name string, isTime ...bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return newMetric(ctx, name, stats.Trend, isTime)
***REMOVED***

func (*Metrics) XRate(ctx *context.Context, name string, isTime ...bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return newMetric(ctx, name, stats.Rate, isTime)
***REMOVED***
