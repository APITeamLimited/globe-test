package v8js

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/runner"
	"github.com/ry/v8worker"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math"
	"os"
	"time"
)

type libFile struct ***REMOVED***
	Filename string
	Source   string
***REMOVED***

type Runner struct ***REMOVED***
	Filename string
	Source   string
	Client   *fasthttp.Client

	stdlib []libFile
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
	***REMOVED***

	// Load the standard library as a rice box; panic if any part of this fails
	// (The only possible cause is a programming/developer error, not user error)
	box := rice.MustFindBox("lib")
	box.Walk("/", func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		if !info.IsDir() ***REMOVED***
			r.stdlib = append(r.stdlib, libFile***REMOVED***
				Filename: path,
				Source:   box.MustString(path),
			***REMOVED***)
		***REMOVED***
		return nil
	***REMOVED***)

	return r
***REMOVED***

func (r *Runner) Run(ctx context.Context, id int64) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		vu := VUContext***REMOVED***r: r, ctx: ctx, ch: ch***REMOVED***

		w := v8worker.New(vu.Recv, vu.RecvSync)

		for _, f := range r.stdlib ***REMOVED***
			if err := w.Load(f.Filename, f.Source); err != nil ***REMOVED***
				log.WithError(err).WithField("file", f.Filename).Error("Couldn't load lib")
			***REMOVED***
		***REMOVED***

		src := fmt.Sprintf("function __run__() ***REMOVED***%s***REMOVED***; undefined", r.Source)
		if err := w.Load(r.Filename, src); err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't load JS")
			return
		***REMOVED***

		for ***REMOVED***
			log.Info("-> run")
			w.Send("run")
			log.Info("<- run")

			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***

func (vu *VUContext) Recv(raw string) ***REMOVED***
***REMOVED***

func (vu *VUContext) RecvSync(raw string) string ***REMOVED***
	return ""
***REMOVED***
