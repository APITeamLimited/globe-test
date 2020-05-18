/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package minirunner

import (
	"context"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
)

// Ensure mock implementations conform to the interfaces.
var (
	_ lib.Runner        = &MiniRunner***REMOVED******REMOVED***
	_ lib.InitializedVU = &VU***REMOVED******REMOVED***
	_ lib.ActiveVU      = &ActiveVU***REMOVED******REMOVED***
)

// MiniRunner partially implements the lib.Runner interface, but instead of
// using a real JS runtime, it allows us to directly specify the options and
// functions with Go code.
type MiniRunner struct ***REMOVED***
	Fn         func(ctx context.Context, out chan<- stats.SampleContainer) error
	SetupFn    func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error)
	TeardownFn func(ctx context.Context, out chan<- stats.SampleContainer) error

	SetupData []byte

	NextVUID int64
	Group    *lib.Group
	Options  lib.Options
***REMOVED***

// MakeArchive isn't implemented, it always returns nil and is just here to
// satisfy the lib.Runner interface.
func (r MiniRunner) MakeArchive() *lib.Archive ***REMOVED***
	return nil
***REMOVED***

// NewVU returns a new VU with an incremental ID.
func (r *MiniRunner) NewVU(id int64, out chan<- stats.SampleContainer) (lib.InitializedVU, error) ***REMOVED***
	return &VU***REMOVED***R: r, Out: out, ID: id***REMOVED***, nil
***REMOVED***

// Setup calls the supplied mock setup() function, if present.
func (r *MiniRunner) Setup(ctx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	if fn := r.SetupFn; fn != nil ***REMOVED***
		r.SetupData, err = fn(ctx, out)
	***REMOVED***
	return
***REMOVED***

// GetSetupData returns json representation of the setup data if setup() is
// specified and was ran, nil otherwise.
func (r MiniRunner) GetSetupData() []byte ***REMOVED***
	return r.SetupData
***REMOVED***

// SetSetupData saves the externally supplied setup data as JSON in the runner.
func (r *MiniRunner) SetSetupData(data []byte) ***REMOVED***
	r.SetupData = data
***REMOVED***

// Teardown calls the supplied mock teardown() function, if present.
func (r MiniRunner) Teardown(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
	if fn := r.TeardownFn; fn != nil ***REMOVED***
		return fn(ctx, out)
	***REMOVED***
	return nil
***REMOVED***

// GetDefaultGroup returns the default group.
func (r MiniRunner) GetDefaultGroup() *lib.Group ***REMOVED***
	if r.Group == nil ***REMOVED***
		r.Group = &lib.Group***REMOVED******REMOVED***
	***REMOVED***
	return r.Group
***REMOVED***

// IsExecutable satisfies lib.Runner, but is mocked for MiniRunner since
// it doesn't deal with JS.
func (r MiniRunner) IsExecutable(name string) bool ***REMOVED***
	return true
***REMOVED***

// GetOptions returns the supplied options struct.
func (r MiniRunner) GetOptions() lib.Options ***REMOVED***
	return r.Options
***REMOVED***

// SetOptions allows you to override the runner options.
func (r *MiniRunner) SetOptions(opts lib.Options) error ***REMOVED***
	r.Options = opts
	return nil
***REMOVED***

// VU is a mock VU, spawned by a MiniRunner.
type VU struct ***REMOVED***
	R         *MiniRunner
	Out       chan<- stats.SampleContainer
	ID        int64
	Iteration int64
***REMOVED***

// ActiveVU holds a VU and its activation parameters
type ActiveVU struct ***REMOVED***
	*VU
	*lib.VUActivationParams
	busy chan struct***REMOVED******REMOVED***
***REMOVED***

// Activate the VU so it will be able to run code.
func (vu *VU) Activate(params *lib.VUActivationParams) lib.ActiveVU ***REMOVED***
	avu := &ActiveVU***REMOVED***
		VU:                 vu,
		VUActivationParams: params,
		busy:               make(chan struct***REMOVED******REMOVED***, 1),
	***REMOVED***

	go func() ***REMOVED***
		<-params.RunContext.Done()

		// Wait for the VU to stop running, if it was, and prevent it from
		// running again for this activation
		avu.busy <- struct***REMOVED******REMOVED******REMOVED******REMOVED***

		if params.DeactivateCallback != nil ***REMOVED***
			params.DeactivateCallback(vu)
		***REMOVED***
	***REMOVED***()

	return avu
***REMOVED***

// RunOnce runs the mock default function once, incrementing its iteration.
func (vu *ActiveVU) RunOnce() error ***REMOVED***
	if vu.R.Fn == nil ***REMOVED***
		return nil
	***REMOVED***

	select ***REMOVED***
	case <-vu.RunContext.Done():
		return vu.RunContext.Err() // we are done, return
	case vu.busy <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		// nothing else can run now, and the VU cannot be deactivated
	***REMOVED***
	defer func() ***REMOVED***
		<-vu.busy // unlock deactivation again
	***REMOVED***()

	state := &lib.State***REMOVED***
		Vu:        vu.ID,
		Iteration: vu.Iteration,
	***REMOVED***
	newctx := lib.WithState(vu.RunContext, state)

	vu.Iteration++

	return vu.R.Fn(newctx, vu.Out)
***REMOVED***
