package js

import (
	"errors"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/runner"
	"golang.org/x/net/context"
	"gopkg.in/olebedev/go-duktape.v2"
	"strconv"
)

type Runner struct ***REMOVED***
***REMOVED***

func New() *Runner ***REMOVED***
	return &Runner***REMOVED******REMOVED***
***REMOVED***

func (r *Runner) Run(ctx context.Context, t loadtest.LoadTest, id int64) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		c := duktape.New()
		defer c.Destroy()

		c.PushGlobalObject()
		c.PushObject() // __internal__

		pushModules(c, t, id, ch)
		c.PutPropString(-2, "modules")

		c.PutPropString(-2, "__internal__")

		if top := c.GetTopIndex(); top != 0 ***REMOVED***
			panic("PROGRAMMING ERROR: Excess items on stack: " + strconv.Itoa(top+1))
		***REMOVED***

		if err := c.PcompileString(0, t.Source); err != nil ***REMOVED***
			ch <- runner.Result***REMOVED***Error: err***REMOVED***
			return
		***REMOVED***

		for ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
				c.DupTop()
				if c.Pcall(0) != duktape.ErrNone ***REMOVED***
					err := errors.New(c.SafeToString(-1))
					ch <- runner.Result***REMOVED***Error: err***REMOVED***
				***REMOVED***
				c.Pop() // Pop return value
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***

func pushModules(c *duktape.Context, t loadtest.LoadTest, id int64, ch <-chan runner.Result) ***REMOVED***
	api := map[string]map[string]apiFunc***REMOVED***
		"http": map[string]apiFunc***REMOVED***
			"get": apiHTTPGet,
		***REMOVED***,
	***REMOVED***

	c.PushObject() // __internal__.modules
	for name, mod := range api ***REMOVED***
		pushModule(c, ch, mod)
		c.PutPropString(-2, name)
	***REMOVED***
***REMOVED***

func pushModule(c *duktape.Context, ch <-chan runner.Result, members map[string]apiFunc) int ***REMOVED***
	idx := c.PushObject()

	for name, fn := range members ***REMOVED***
		c.PushGoFunction(func(lc *duktape.Context) int ***REMOVED***
			return fn(lc, ch)
		***REMOVED***)
		c.PutPropString(idx, name)
	***REMOVED***

	return idx
***REMOVED***
