package jstest

import (
	"fmt"
	"time"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/stats"
)

func init() ***REMOVED***
	modules.Register("k6/x/jsexttest", New())
***REMOVED***

type (
	RootModule struct***REMOVED******REMOVED***

	// JSTest is meant to test xk6 and the JS extension sub-system of k6.
	JSTest struct ***REMOVED***
		vu modules.VU

		foos *stats.Metric
	***REMOVED***
)

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &JSTest***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.Module interface and returns
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &JSTest***REMOVED***
		vu:   vu,
		foos: vu.InitEnv().Registry.MustNewMetric("foos", stats.Counter),
	***REMOVED***
***REMOVED***

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (j *JSTest) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***Default: j***REMOVED***
***REMOVED***

// Foo emits a foo metric
func (j *JSTest) Foo(arg float64) (bool, error) ***REMOVED***
	state := j.vu.State()
	if state == nil ***REMOVED***
		return false, fmt.Errorf("the VU State is not avaialble in the init context")
	***REMOVED***

	ctx := j.vu.Context()

	tags := state.CloneTags()
	tags["foo"] = "bar"
	stats.PushIfNotDone(ctx, state.Samples, stats.Sample***REMOVED***
		Time:   time.Now(),
		Metric: j.foos, Tags: stats.IntoSampleTags(&tags),
		Value: arg,
	***REMOVED***)

	return true, nil
***REMOVED***
