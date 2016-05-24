package js

import (
	"errors"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"gopkg.in/olebedev/go-duktape.v2"
	"strconv"
	"time"
)

const srcRequire = `
function require(name) ***REMOVED***
	var mod = __internal__.modules[name];
	if (!mod) ***REMOVED***
		throw new Error("Unknown module: " + name);
	***REMOVED***
	return mod;
***REMOVED***
`

type Runner struct ***REMOVED***
	Client *fasthttp.Client
***REMOVED***

func New() *Runner ***REMOVED***
	return &Runner***REMOVED***
		Client: &fasthttp.Client***REMOVED***
			MaxIdleConnDuration: time.Duration(0),
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (r *Runner) Run(ctx context.Context, t loadtest.LoadTest, id int64) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		c := duktape.New()
		defer c.Destroy()

		c.PushGlobalObject()

		c.PushObject()
		***REMOVED***
			pushModules(c, r, ch)
			c.PutPropString(-2, "modules")

			pushData(c, t, id)
			c.PutPropString(-2, "data")
		***REMOVED***
		c.PutPropString(-2, "__internal__")

		if err := c.PcompileString(duktape.CompileFunction|duktape.CompileStrict, srcRequire); err != nil ***REMOVED***
			ch <- runner.Result***REMOVED***Error: err***REMOVED***
			return
		***REMOVED***
		c.PutPropString(-2, "require")

		// This will probably be moved to a module; global for compatibility
		c.PushGoFunction(func(c *duktape.Context) int ***REMOVED***
			t := argNumber(c, 0)
			time.Sleep(time.Duration(t) * time.Second)
			return 0
		***REMOVED***)
		c.PutPropString(-2, "sleep")

		if top := c.GetTopIndex(); top != 0 ***REMOVED***
			panic("PROGRAMMING ERROR: Excess items on stack: " + strconv.Itoa(top+1))
		***REMOVED***

		// It should be cheaper memory-wise to keep the code on the global object (where it's
		// safe from GC shenanigans) and duplicate the key every iteration, than to keep it on
		// the stack and duplicate the whole function every iteration; just make it ABUNDANTLY
		// CLEAR that the script should NEVER touch that property or dumb things will happen
		codeProp := "__!!__seriously__don't__touch__from__script__!!__"
		if err := c.PcompileString(0, t.Source); err != nil ***REMOVED***
			ch <- runner.Result***REMOVED***Error: err***REMOVED***
			return
		***REMOVED***
		c.PutPropString(-2, codeProp)

		c.PushString(codeProp)
		for ***REMOVED***
			select ***REMOVED***
			default:
				c.DupTop()
				if c.PcallProp(-3, 0) != duktape.ErrNone ***REMOVED***
					err := errors.New(c.SafeToString(-1))
					ch <- runner.Result***REMOVED***Error: err***REMOVED***
				***REMOVED***
				c.Pop() // Pop return value
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***

func pushData(c *duktape.Context, t loadtest.LoadTest, id int64) ***REMOVED***
	c.PushObject()

	c.PushInt(int(id))
	c.PutPropString(-2, "id")

	c.PushObject()
	***REMOVED***
		c.PushString(t.URL)
		c.PutPropString(-2, "url")
	***REMOVED***
	c.PutPropString(-2, "test")
***REMOVED***

func pushModules(c *duktape.Context, r *Runner, ch chan<- runner.Result) ***REMOVED***
	c.PushObject()

	api := map[string]map[string]apiFunc***REMOVED***
		"http": map[string]apiFunc***REMOVED***
			"get": apiHTTPGet,
		***REMOVED***,
	***REMOVED***
	for name, mod := range api ***REMOVED***
		pushModule(c, r, ch, mod)
		c.PutPropString(-2, name)
	***REMOVED***
***REMOVED***

func pushModule(c *duktape.Context, r *Runner, ch chan<- runner.Result, members map[string]apiFunc) ***REMOVED***
	c.PushObject()

	for name, fn := range members ***REMOVED***
		c.PushGoFunction(func(lc *duktape.Context) int ***REMOVED***
			return fn(r, lc, ch)
		***REMOVED***)
		c.PutPropString(-2, name)
	***REMOVED***
***REMOVED***
