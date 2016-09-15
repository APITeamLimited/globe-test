package js

import (
	"context"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
)

type Runner struct ***REMOVED***
	Runtime *Runtime
	Exports otto.Value
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	return &VU***REMOVED******REMOVED***, nil
***REMOVED***

type VU struct ***REMOVED***
	ID int64
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
	return nil, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	return nil
***REMOVED***
