package js2

import (
	"context"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/spf13/afero"
)

type Runner struct ***REMOVED***
	Options lib.Options

	defaultGroup *lib.Group
***REMOVED***

func New(src *lib.SourceData, fs afero.Fs) (*Runner, error) ***REMOVED***
	defaultGroup, err := lib.NewGroup("", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Runner***REMOVED***
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
	return &VU***REMOVED******REMOVED***, nil
***REMOVED***

func (r *Runner) GetDefaultGroup() *lib.Group ***REMOVED***
	return r.defaultGroup
***REMOVED***

func (r *Runner) GetOptions() lib.Options ***REMOVED***
	return r.Options
***REMOVED***

func (r *Runner) ApplyOptions(opts lib.Options) ***REMOVED***
	r.Options = r.Options.Apply(opts)
***REMOVED***

type VU struct ***REMOVED***
	VM *goja.Runtime
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
	return []stats.Sample***REMOVED******REMOVED***, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	return nil
***REMOVED***
