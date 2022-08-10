// Package k6 implements the module imported as 'k6' from inside k6.
package k6

import (
	"errors"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/dop251/goja"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

var (
	// ErrGroupInInitContext is returned when group() are using in the init context.
	ErrGroupInInitContext = common.NewInitContextError("Using group() in the init context is not supported")

	// ErrCheckInInitContext is returned when check() are using in the init context.
	ErrCheckInInitContext = common.NewInitContextError("Using check() in the init context is not supported")
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct***REMOVED******REMOVED***

	// K6 represents an instance of the k6 module.
	K6 struct ***REMOVED***
		vu modules.VU
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &K6***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &K6***REMOVED***vu: vu***REMOVED***
***REMOVED***

// Exports returns the exports of the k6 module.
func (mi *K6) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"check":      mi.Check,
			"fail":       mi.Fail,
			"group":      mi.Group,
			"randomSeed": mi.RandomSeed,
			"sleep":      mi.Sleep,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Fail is a fancy way of saying `throw "something"`.
func (*K6) Fail(msg string) (goja.Value, error) ***REMOVED***
	return goja.Undefined(), errors.New(msg)
***REMOVED***

// Sleep waits the provided seconds before continuing the execution.
func (mi *K6) Sleep(secs float64) ***REMOVED***
	ctx := mi.vu.Context()
	timer := time.NewTimer(time.Duration(secs * float64(time.Second)))
	select ***REMOVED***
	case <-timer.C:
	case <-ctx.Done():
		timer.Stop()
	***REMOVED***
***REMOVED***

// RandomSeed sets the seed to the random generator used for this VU.
func (mi *K6) RandomSeed(seed int64) ***REMOVED***
	randSource := rand.New(rand.NewSource(seed)).Float64 //nolint:gosec
	mi.vu.Runtime().SetRandSource(randSource)
***REMOVED***

// Group wraps a function call and executes it within the provided group name.
func (mi *K6) Group(name string, fn goja.Callable) (goja.Value, error) ***REMOVED***
	state := mi.vu.State()
	if state == nil ***REMOVED***
		return nil, ErrGroupInInitContext
	***REMOVED***

	if fn == nil ***REMOVED***
		return nil, errors.New("group() requires a callback as a second argument")
	***REMOVED***

	g, err := state.Group.Group(name)
	if err != nil ***REMOVED***
		return goja.Undefined(), err
	***REMOVED***

	old := state.Group
	state.Group = g

	shouldUpdateTag := state.Options.SystemTags.Has(metrics.TagGroup)
	if shouldUpdateTag ***REMOVED***
		state.Tags.Set("group", g.Path)
	***REMOVED***
	defer func() ***REMOVED***
		state.Group = old
		if shouldUpdateTag ***REMOVED***
			state.Tags.Set("group", old.Path)
		***REMOVED***
	***REMOVED***()

	startTime := time.Now()
	ret, err := fn(goja.Undefined())
	t := time.Now()

	tags := state.CloneTags()

	ctx := mi.vu.Context()
	metrics.PushIfNotDone(ctx, state.Samples, metrics.Sample***REMOVED***
		Time:   t,
		Metric: state.BuiltinMetrics.GroupDuration,
		Tags:   metrics.IntoSampleTags(&tags),
		Value:  metrics.D(t.Sub(startTime)),
	***REMOVED***)

	return ret, err
***REMOVED***

// Check will emit check metrics for the provided checks.
//nolint:cyclop
func (mi *K6) Check(arg0, checks goja.Value, extras ...goja.Value) (bool, error) ***REMOVED***
	state := mi.vu.State()
	if state == nil ***REMOVED***
		return false, ErrCheckInInitContext
	***REMOVED***
	if checks == nil ***REMOVED***
		return false, errors.New("no checks provided to `check`")
	***REMOVED***
	ctx := mi.vu.Context()
	rt := mi.vu.Runtime()
	t := time.Now()

	// Prepare the metric tags
	commonTags := state.CloneTags()
	if len(extras) > 0 ***REMOVED***
		obj := extras[0].ToObject(rt)
		for _, k := range obj.Keys() ***REMOVED***
			commonTags[k] = obj.Get(k).String()
		***REMOVED***
	***REMOVED***

	succ := true
	var exc error
	obj := checks.ToObject(rt)
	for _, name := range obj.Keys() ***REMOVED***
		val := obj.Get(name)

		tags := make(map[string]string, len(commonTags))
		for k, v := range commonTags ***REMOVED***
			tags[k] = v
		***REMOVED***

		// Resolve the check record.
		check, err := state.Group.Check(name)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		if state.Options.SystemTags.Has(metrics.TagCheck) ***REMOVED***
			tags["check"] = check.Name
		***REMOVED***

		// Resolve callables into values.
		fn, ok := goja.AssertFunction(val)
		if ok ***REMOVED***
			tmpVal, err := fn(goja.Undefined(), arg0)
			val = tmpVal
			if err != nil ***REMOVED***
				val = rt.ToValue(false)
				exc = err
			***REMOVED***
		***REMOVED***

		sampleTags := metrics.IntoSampleTags(&tags)

		// Emit! (But only if we have a valid context.)
		select ***REMOVED***
		case <-ctx.Done():
		default:
			if val.ToBoolean() ***REMOVED***
				atomic.AddInt64(&check.Passes, 1)
				metrics.PushIfNotDone(ctx, state.Samples,
					metrics.Sample***REMOVED***Time: t, Metric: state.BuiltinMetrics.Checks, Tags: sampleTags, Value: 1***REMOVED***)
			***REMOVED*** else ***REMOVED***
				atomic.AddInt64(&check.Fails, 1)
				metrics.PushIfNotDone(ctx, state.Samples,
					metrics.Sample***REMOVED***Time: t, Metric: state.BuiltinMetrics.Checks, Tags: sampleTags, Value: 0***REMOVED***)
				// A single failure makes the return value false.
				succ = false
			***REMOVED***
		***REMOVED***

		if exc != nil ***REMOVED***
			return succ, exc
		***REMOVED***
	***REMOVED***

	return succ, nil
***REMOVED***
