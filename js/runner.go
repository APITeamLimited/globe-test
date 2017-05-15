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
	"net"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
	"github.com/spf13/afero"
)

type Runner struct ***REMOVED***
	Bundle       *Bundle
	Logger       *log.Logger
	defaultGroup *lib.Group

	Dialer *netext.Dialer
***REMOVED***

func New(src *lib.SourceData, fs afero.Fs) (*Runner, error) ***REMOVED***
	bundle, err := NewBundle(src, fs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewFromBundle(bundle)
***REMOVED***

func NewFromArchive(arc *lib.Archive) (*Runner, error) ***REMOVED***
	bundle, err := NewBundleFromArchive(arc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewFromBundle(bundle)
***REMOVED***

func NewFromBundle(b *Bundle) (*Runner, error) ***REMOVED***
	defaultGroup, err := lib.NewGroup("", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Runner***REMOVED***
		Bundle:       b,
		Logger:       log.StandardLogger(),
		defaultGroup: defaultGroup,
		Dialer: netext.NewDialer(net.Dialer***REMOVED***
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		***REMOVED***),
	***REMOVED***, nil
***REMOVED***

func (r *Runner) MakeArchive() *lib.Archive ***REMOVED***
	return r.Bundle.MakeArchive()
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
		HTTPTransport:  &http.Transport***REMOVED***DialContext: r.Dialer.DialContext***REMOVED***,
		VUContext:      NewVUContext(),
	***REMOVED***
	common.BindToGlobal(vu.Runtime, common.Bind(vu.Runtime, vu.VUContext, vu.Context))

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

	Runner        *Runner
	HTTPTransport *http.Transport
	ID            int64
	Iteration     int64

	VUContext *VUContext
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
	state := &common.State***REMOVED***
		Logger:        u.Runner.Logger,
		Options:       u.Runner.Bundle.Options,
		Group:         u.Runner.defaultGroup,
		HTTPTransport: u.HTTPTransport,
	***REMOVED***

	ctx = common.WithRuntime(ctx, u.Runtime)
	ctx = common.WithState(ctx, state)
	*u.Context = ctx

	u.Runtime.Set("__ITER", u.Iteration)
	u.Iteration++

	_, err := u.Default(goja.Undefined())

	if u.Runner.Bundle.Options.NoConnectionReuse.Bool ***REMOVED***
		u.HTTPTransport.CloseIdleConnections()
	***REMOVED***
	return state.Samples, err
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	u.Iteration = 0
	u.Runtime.Set("__VU", u.ID)
	return nil
***REMOVED***
