package js

import (
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
		setupInternals(c, t, id, ch)
	***REMOVED***()

	return ch
***REMOVED***

func setupInternals(c *duktape.Context, t loadtest.LoadTest, id int64, ch <-chan runner.Result) ***REMOVED***
	api := map[string]map[string]apiFunc***REMOVED***
		"http": map[string]apiFunc***REMOVED***
			"get": apiHTTPGet,
		***REMOVED***,
	***REMOVED***

	c.PushGlobalObject()

	c.PushObject() // __internal__

	c.PushObject() // __internal__.modules
	for name, mod := range api ***REMOVED***
		pushModule(c, ch, mod)
		c.PutPropString(-2, name)
	***REMOVED***
	c.PutPropString(-2, "modules")

	c.PutPropString(-2, "__internal__")

	c.Pop() // global object

	if top := c.GetTopIndex(); !(top < 0) ***REMOVED*** // < 0 = invalid index = empty stack
		panic("PROGRAMMING ERROR: Stack depth must be 0, is: " + strconv.Itoa(top+1))
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
