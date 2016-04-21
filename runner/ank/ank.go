package ank

import (
	"errors"
	"fmt"
	"github.com/loadimpact/speedboat/runner"
	anko_core "github.com/mattn/anko/builtins"
	anko "github.com/mattn/anko/vm"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math"
	"time"
)

type Runner struct ***REMOVED***
	Filename string
	Source   string
	Client   *fasthttp.Client
***REMOVED***

type VUContext struct ***REMOVED***
	r   *Runner
	ctx context.Context
	ch  chan runner.Result
***REMOVED***

func New(filename, src string) *Runner ***REMOVED***
	return &Runner***REMOVED***
		Filename: filename,
		Source:   src,
		Client: &fasthttp.Client***REMOVED***
			Dial:                fasthttp.Dial,
			MaxIdleConnDuration: time.Duration(0),
			MaxConnsPerHost:     math.MaxInt64,
		***REMOVED***,
	***REMOVED***
***REMOVED***
func (r *Runner) Run(ctx context.Context, id int64) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		vu := VUContext***REMOVED***r: r, ctx: ctx, ch: ch***REMOVED***

		vm := anko.NewEnv()
		anko_core.Import(vm)

		vm.Set("__id", id)
		vm.Define("sleep", vu.Sleep)

		pkgs := map[string]func(env *anko.Env) *anko.Env***REMOVED***
			"http": vu.HTTPLoader,
		***REMOVED***
		vm.Define("import", func(s string) interface***REMOVED******REMOVED*** ***REMOVED***
			if loader, ok := pkgs[s]; ok ***REMOVED***
				m := loader(vm)
				return m
			***REMOVED***
			ch <- runner.Result***REMOVED***Error: errors.New(fmt.Sprintf("Package not found: %s", s))***REMOVED***
			return nil
		***REMOVED***)

		for ***REMOVED***
			if _, err := vm.Execute(r.Source); err != nil ***REMOVED***
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
