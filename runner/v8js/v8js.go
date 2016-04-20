package v8js

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/runner"
	"github.com/ry/v8worker"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math"
	"strings"
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

		w := v8worker.New(func(msg string) ***REMOVED******REMOVED***, func(msg string) string ***REMOVED***
			parts := strings.SplitN(msg, ";", 2)
			switch parts[0] ***REMOVED***
			case "get":
				vu.HTTPGet(parts[1])
			case "sleep":
				vu.Sleep(parts[1])
			default:
				log.WithField("call", parts[0]).Fatal("Unknown JS call")
			***REMOVED***
			return ""
		***REMOVED***)
		w.Load(r.Filename, fmt.Sprintf(`
		$recvSync(function(msg) ***REMOVED***
			if(msg == 'run') ***REMOVED***
				run()
			***REMOVED***
		***REMOVED***)
		function get(url) ***REMOVED***
			$sendSync('get;' + url)
		***REMOVED***
		function sleep(t) ***REMOVED***
			$sendSync('sleep;' + t)
		***REMOVED***
		function run() ***REMOVED***
			%s
		***REMOVED***
		`, r.Source))

		// vm := otto.New()
		// vm.Set("__id", id)
		// vm.Set("get", vu.HTTPGet)
		// vm.Set("sleep", vu.Sleep)

		// script, err := vm.Compile(r.Filename, r.Source)
		// if err != nil ***REMOVED***
		// 	ch <- runner.Result***REMOVED***Error: err***REMOVED***
		// 	return
		// ***REMOVED***

		for ***REMOVED***
			// if _, err := vm.Run(script); err != nil ***REMOVED***
			// 	ch <- runner.Result***REMOVED***Error: err***REMOVED***
			// ***REMOVED***
			w.SendSync("run")

			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
