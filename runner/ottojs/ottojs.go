package ottojs

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math"
	"sync"
	"time"
)

type Runner struct ***REMOVED***
	Filename string
	Source   string
	Client   *fasthttp.Client
	VMs      sync.Pool
***REMOVED***

type VUContext struct ***REMOVED***
	r   *Runner
	ctx context.Context
	ch  chan runner.Result
***REMOVED***

func New(filename, src string) *Runner ***REMOVED***
	r := &Runner***REMOVED***
		Filename: filename,
		Source:   src,
		Client: &fasthttp.Client***REMOVED***
			Dial:                fasthttp.Dial,
			MaxIdleConnDuration: time.Duration(0),
			MaxConnsPerHost:     math.MaxInt64,
		***REMOVED***,
		VMs: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED***
				return otto.New()
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for i := 0; i < 10000; i++ ***REMOVED***
		r.VMs.Put(r.VMs.New())
	***REMOVED***
	return r
***REMOVED***

func (r *Runner) Run(ctx context.Context, id int64) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		vu := VUContext***REMOVED***r: r, ctx: ctx, ch: ch***REMOVED***

		vm := r.VMs.Get().(*otto.Otto)
		defer r.VMs.Put(vm)

		vm.Set("__id", id)
		vm.Set("get", vu.HTTPGet)
		vm.Set("sleep", vu.Sleep)

		script, err := vm.Compile(r.Filename, r.Source)
		if err != nil ***REMOVED***
			ch <- runner.Result***REMOVED***Error: err***REMOVED***
			return
		***REMOVED***

		for ***REMOVED***
			if _, err := vm.Run(script); err != nil ***REMOVED***
				ch <- runner.Result***REMOVED***Error: err***REMOVED***
			***REMOVED***

			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
