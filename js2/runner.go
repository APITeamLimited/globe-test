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

package js2

import (
	"context"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js2/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/spf13/afero"
)

type Runner struct ***REMOVED***
	Bundle *Bundle

	defaultGroup *lib.Group
***REMOVED***

func New(src *lib.SourceData, fs afero.Fs) (*Runner, error) ***REMOVED***
	bundle, err := NewBundle(src, fs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defaultGroup, err := lib.NewGroup("", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Runner***REMOVED***
		Bundle:       bundle,
		defaultGroup: defaultGroup,
	***REMOVED***, nil
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	vu, err := r.newVU()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return lib.VU(vu), nil
***REMOVED***

func (r *Runner) newVU() (*VU, error) ***REMOVED***
	// Instantiate a new bundle, make a VU out of it.
	bi, err := r.Bundle.Instantiate()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Make a VU, apply the VU context.
	vu := &VU***REMOVED***
		BundleInstance: *bi,
		Runner:         r,
		VUContext:      NewVUContext(),
	***REMOVED***
	common.BindToGlobal(vu.Runtime, vu.VUContext)

	// Give the VU an initial sense of identity.
	if err := vu.Reconfigure(0); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return vu, nil
***REMOVED***

func (r *Runner) GetDefaultGroup() *lib.Group ***REMOVED***
	return r.defaultGroup
***REMOVED***

func (r *Runner) GetOptions() lib.Options ***REMOVED***
	return r.Bundle.Options
***REMOVED***

func (r *Runner) ApplyOptions(opts lib.Options) ***REMOVED***
	r.Bundle.Options = r.Bundle.Options.Apply(opts)
***REMOVED***

type VU struct ***REMOVED***
	BundleInstance

	Runner    *Runner
	ID        int64
	Iteration int64

	Samples   []stats.Sample
	VUContext *VUContext
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
	ctx = common.WithState(ctx, &common.State***REMOVED***
		Group: u.Runner.defaultGroup,
	***REMOVED***)

	u.Runtime.Set("__ITER", u.Iteration)
	u.Iteration++

	_, err := u.Default(goja.Undefined())
	samples := u.Samples
	u.Samples = nil
	return samples, err
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	u.Iteration = 0
	u.Runtime.Set("__VU", u.ID)
	return nil
***REMOVED***
