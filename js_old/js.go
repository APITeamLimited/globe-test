package js

import (
	"encoding/json"
	"errors"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"gopkg.in/olebedev/go-duktape.v2"
	"strconv"
	"time"
)

type Runner struct ***REMOVED***
	Client *fasthttp.Client

	source string
	lib    *rice.Box
	vendor *rice.Box

	mDuration *sampler.Metric
	mErrors   *sampler.Metric
***REMOVED***

func New(src string) *Runner ***REMOVED***
	return &Runner***REMOVED***
		Client: &fasthttp.Client***REMOVED***
			MaxIdleConnDuration: time.Duration(0),
		***REMOVED***,
		source:    src,
		lib:       rice.MustFindBox("lib"),
		vendor:    rice.MustFindBox("vendor"),
		mDuration: sampler.Stats("duration"),
		mErrors:   sampler.Counter("errors"),
	***REMOVED***
***REMOVED***

func (r *Runner) RunVU(ctx context.Context, t speedboat.Test, id int) ***REMOVED***
	c, err := r.newJSContext(t, id)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create JS context")
		return
	***REMOVED***
	defer c.Destroy()

	c.PushGlobalObject()

	// It should be cheaper memory-wise to keep the code on the global object (where it's
	// safe from GC shenanigans) and duplicate the key every iteration, than to keep it on
	// the stack and duplicate the whole function every iteration; just make it ABUNDANTLY
	// CLEAR that the script should NEVER touch that property or dumb things will happen
	codeProp := "__!!__seriously__don't__touch__from__script__!!__"
	c.PushString("script")
	if err := c.PcompileStringFilename(0, r.source); err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't compile script")
		return
	***REMOVED***
	c.PutPropString(-2, codeProp)

	iteration := 1
	c.PushString(codeProp)
	for ***REMOVED***
		select ***REMOVED***
		default:
			setIterationCounter(c, iteration)
			c.DupTop()
			if c.PcallProp(-3, 0) != duktape.ErrNone ***REMOVED***
				c.GetPropString(-1, "fileName")
				filename := c.SafeToString(-1)
				c.Pop()

				c.GetPropString(-1, "lineNumber")
				line := c.ToNumber(-1)
				c.Pop()

				err := errors.New(c.SafeToString(-1))
				log.WithError(err).WithFields(log.Fields***REMOVED***
					"file": filename,
					"line": line,
				***REMOVED***).Error("Script Error")
			***REMOVED***
			c.Pop() // Pop return value
		case <-ctx.Done():
			return
		***REMOVED***

		iteration++
	***REMOVED***
***REMOVED***

func (r *Runner) newJSContext(t speedboat.Test, id int) (*duktape.Context, error) ***REMOVED***
	c := duktape.New()
	c.PushGlobalObject()

	c.PushObject()
	***REMOVED***
		pushModules(c, r)
		c.PutPropString(-2, "modules")

		c.PushObject()
		c.PutPropString(-2, "types")

		pushData(c, t, id)
		c.PutPropString(-2, "data")
	***REMOVED***
	c.PutPropString(-2, "__internal__")

	load := map[*rice.Box][]string***REMOVED***
		r.lib:    []string***REMOVED***"require.js", "http.js", "log.js", "vu.js", "test.js"***REMOVED***,
		r.vendor: []string***REMOVED***"lodash/dist/lodash.min.js"***REMOVED***,
	***REMOVED***
	for box, files := range load ***REMOVED***
		for _, name := range files ***REMOVED***
			src, err := box.String(name)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if err = loadFile(c, name, src); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

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

	c.Pop() // Global object
	return c, nil
***REMOVED***

func setIterationCounter(c *duktape.Context, i int) ***REMOVED***
	c.PushGlobalObject()
	***REMOVED***
		c.GetPropString(-1, "__internal__")
		***REMOVED***
			c.PushInt(i)
			c.PutPropString(-2, "iteration")
		***REMOVED***
		c.Pop()
	***REMOVED***
	c.Pop()
***REMOVED***

func pushData(c *duktape.Context, t speedboat.Test, id int) ***REMOVED***
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

func pushModules(c *duktape.Context, r *Runner) ***REMOVED***
	c.PushObject()

	api := map[string]map[string]apiFunc***REMOVED***
		"http": map[string]apiFunc***REMOVED***
			"do": apiHTTPDo,
			"setMaxConnectionsPerHost": apiHTTPSetMaxConnectionsPerHost,
		***REMOVED***,
		"log": map[string]apiFunc***REMOVED***
			"type": apiLogType,
		***REMOVED***,
		"test": map[string]apiFunc***REMOVED***
			"abort": apiTestAbort,
		***REMOVED***,
		"vu": map[string]apiFunc***REMOVED******REMOVED***,
	***REMOVED***
	for name, mod := range api ***REMOVED***
		pushModule(c, r, mod)
		c.PutPropString(-2, name)
	***REMOVED***
***REMOVED***

func pushModule(c *duktape.Context, r *Runner, members map[string]apiFunc) ***REMOVED***
	c.PushObject()

	for name, fn := range members ***REMOVED***
		fn := fn
		c.PushGoFunction(func(lc *duktape.Context) int ***REMOVED***
			return fn(r, lc)
		***REMOVED***)
		c.PutPropString(-2, name)
	***REMOVED***
***REMOVED***

func loadFile(c *duktape.Context, name, src string) error ***REMOVED***
	c.PushString(name)
	if err := c.PcompileStringFilename(0, src); err != nil ***REMOVED***
		return err
	***REMOVED***
	c.Pcall(0)
	c.Pop()
	return nil
***REMOVED***

func pushInstance(c *duktape.Context, obj interface***REMOVED******REMOVED***, t string) error ***REMOVED***
	s, err := json.Marshal(obj)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.PushString(string(s))
	c.JsonDecode(-1)

	if t != "" ***REMOVED***
		c.PushGlobalObject()
		***REMOVED***
			c.GetPropString(-1, t)
			c.SetPrototype(-3)
		***REMOVED***
		c.Pop()
	***REMOVED***

	return nil
***REMOVED***
