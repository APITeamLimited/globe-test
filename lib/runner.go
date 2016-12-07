package lib

import (
	"context"
	"github.com/loadimpact/k6/stats"
)

// A Runner is a factory for VUs.
type Runner interface ***REMOVED***
	// Creates a new VU. As much as possible should be precomputed here, to allow a pool
	// of prepared VUs to be used to quickly scale up and down.
	NewVU() (VU, error)

	// Returns all known groups.
	GetGroups() []*Group

	// Returns all known tests.
	GetChecks() []*Check

	// Returns the option set.
	GetOptions() Options

	// Applies a set of options.
	ApplyOptions(opts Options)
***REMOVED***

// A VU is a Virtual User.
type VU interface ***REMOVED***
	// Runs the VU once. An iteration should be completely self-contained, and no state
	// or open connections should carry over from one iteration to the next.
	RunOnce(ctx context.Context) ([]stats.Sample, error)

	// Called when the VU's identity changes.
	Reconfigure(id int64) error
***REMOVED***
