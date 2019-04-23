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

package lib

import (
	"context"

	"github.com/loadimpact/k6/stats"
)

// Ensure mock implementations conform to the interfaces.
var _ Runner = &MiniRunner***REMOVED******REMOVED***
var _ VU = &MiniRunnerVU***REMOVED******REMOVED***

// A Runner is a factory for VUs. It should precompute as much as possible upon creation (parse
// ASTs, load files into memory, etc.), so that spawning VUs becomes as fast as possible.
// The Runner doesn't actually *do* anything in itself, the Executor is responsible for wrapping
// and scheduling these VUs for execution.
//
// TODO: Rename this to something more obvious? This name made sense a very long time ago.
type Runner interface ***REMOVED***
	// Creates an Archive of the runner. There should be a corresponding NewFromArchive() function
	// that will restore the runner from the archive.
	MakeArchive() *Archive

	// Spawns a new VU. It's fine to make this function rather heavy, if it means a performance
	// improvement at runtime. Remember, this is called once per VU and normally only at the start
	// of a test - RunOnce() may be called hundreds of thousands of times, and must be fast.
	//TODO: pass context.Context, so initialization can be killed properly...
	NewVU(out chan<- stats.SampleContainer) (VU, error)

	// Runs pre-test setup, if applicable.
	Setup(ctx context.Context, out chan<- stats.SampleContainer) error

	// Returns json representation of the setup data if setup() is specified and run, nil otherwise
	GetSetupData() []byte

	// Saves the externally supplied setup data as json in the runner
	SetSetupData([]byte)

	// Runs post-test teardown, if applicable.
	Teardown(ctx context.Context, out chan<- stats.SampleContainer) error

	// Returns the default (root) Group.
	GetDefaultGroup() *Group

	// Get and set options. The initial value will be whatever the script specifies (for JS,
	// `export let options = ***REMOVED******REMOVED***`); cmd/run.go will mix this in with CLI-, config- and env-provided
	// values and write it back to the runner.
	GetOptions() Options
	SetOptions(opts Options) error
***REMOVED***

// A VU is a Virtual User, that can be scheduled by an Executor.
type VU interface ***REMOVED***
	// Runs the VU once. The VU is responsible for handling the Halting Problem, eg. making sure
	// that execution actually stops when the context is cancelled.
	RunOnce(ctx context.Context) error

	// Assign the VU a new ID. Called by the Executor upon creation, but may be called multiple
	// times if the VU is recycled because the test was scaled down and then back up.
	//TODO: support reconfiguring of env vars, tags and exec
	Reconfigure(id int64) error
***REMOVED***

// MiniRunner wraps a function in a runner whose VUs will simply call that function.
//TODO: move to testutils, or somewhere else that's not lib...
type MiniRunner struct ***REMOVED***
	Fn         func(ctx context.Context, out chan<- stats.SampleContainer) error
	SetupFn    func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error)
	TeardownFn func(ctx context.Context, out chan<- stats.SampleContainer) error

	setupData []byte

	Group   *Group
	Options Options
***REMOVED***

func (r MiniRunner) VU(out chan<- stats.SampleContainer) *MiniRunnerVU ***REMOVED***
	return &MiniRunnerVU***REMOVED***R: r, Out: out***REMOVED***
***REMOVED***

func (r MiniRunner) MakeArchive() *Archive ***REMOVED***
	return nil
***REMOVED***

func (r MiniRunner) NewVU(out chan<- stats.SampleContainer) (VU, error) ***REMOVED***
	return r.VU(out), nil
***REMOVED***

func (r *MiniRunner) Setup(ctx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	if fn := r.SetupFn; fn != nil ***REMOVED***
		r.setupData, err = fn(ctx, out)
	***REMOVED***
	return
***REMOVED***

// GetSetupData returns json representation of the setup data if setup() is specified and run, nil otherwise
func (r MiniRunner) GetSetupData() []byte ***REMOVED***
	return r.setupData
***REMOVED***

// SetSetupData saves the externally supplied setup data as json in the runner
func (r *MiniRunner) SetSetupData(data []byte) ***REMOVED***
	r.setupData = data
***REMOVED***

func (r MiniRunner) Teardown(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
	if fn := r.TeardownFn; fn != nil ***REMOVED***
		return fn(ctx, out)
	***REMOVED***
	return nil
***REMOVED***

func (r MiniRunner) GetDefaultGroup() *Group ***REMOVED***
	if r.Group == nil ***REMOVED***
		r.Group = &Group***REMOVED******REMOVED***
	***REMOVED***
	return r.Group
***REMOVED***

func (r MiniRunner) GetOptions() Options ***REMOVED***
	return r.Options
***REMOVED***

func (r *MiniRunner) SetOptions(opts Options) error ***REMOVED***
	r.Options = opts
	return nil
***REMOVED***

// A VU spawned by a MiniRunner.
type MiniRunnerVU struct ***REMOVED***
	R   MiniRunner
	Out chan<- stats.SampleContainer
	ID  int64
***REMOVED***

func (vu MiniRunnerVU) RunOnce(ctx context.Context) error ***REMOVED***
	if vu.R.Fn == nil ***REMOVED***
		return nil
	***REMOVED***
	return vu.R.Fn(ctx, vu.Out)
***REMOVED***

func (vu *MiniRunnerVU) Reconfigure(id int64) error ***REMOVED***
	vu.ID = id
	return nil
***REMOVED***
