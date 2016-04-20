package duktapejs

import (
	"errors"
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"gopkg.in/olebedev/go-duktape.v2"
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

		c := duktape.New()

		// Get a reference to the global object
		c.PushGlobalObject()
		globalIndex := c.RequireTopIndex()

		// Set __id (key + val -> -)
		c.PushString("__id")
		c.PushInt(int(id))
		if !c.PutProp(globalIndex) ***REMOVED***
			ch <- runner.Result***REMOVED***Error: errors.New("Couldn't push __id")***REMOVED***
			return
		***REMOVED***

		// Bridge functions (no stack change)
		if _, err := c.PushGlobalGoFunction("get", vu.HTTPGet); err != nil ***REMOVED***
			ch <- runner.Result***REMOVED***Error: err***REMOVED***
			return
		***REMOVED***
		if _, err := c.PushGlobalGoFunction("sleep", vu.Sleep); err != nil ***REMOVED***
			ch <- runner.Result***REMOVED***Error: err***REMOVED***
		***REMOVED***

		// Compile the script (source + filename -> func)
		c.PushString(r.Source)
		c.PushString(r.Filename)
		if err := c.Pcompile(0); err != nil ***REMOVED***
			ch <- runner.Result***REMOVED***Error: err***REMOVED***
			return
		***REMOVED***

		// Set it as the global __code__ (key + val -> -)
		c.PushString("__code__")
		c.Insert(-2)
		if !c.PutProp(globalIndex) ***REMOVED***
			ch <- runner.Result***REMOVED***Error: errors.New("Couldn't push __code__")***REMOVED***
			return
		***REMOVED***

		for ***REMOVED***
			c.PushGlobalObject()
			c.PushString("__code__")
			if code := c.PcallProp(-2, 0); code != duktape.ExecSuccess ***REMOVED***
				e := c.SafeToString(-1)
				c.Pop()
				ch <- runner.Result***REMOVED***Error: errors.New(e)***REMOVED***
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
