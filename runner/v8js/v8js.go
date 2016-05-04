package v8js

import (
	"encoding/json"
	"fmt"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/api"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/runner"
	"github.com/ry/v8worker"
	"golang.org/x/net/context"
	"os"
)

type libFile struct ***REMOVED***
	Filename string
	Source   string
***REMOVED***

type Runner struct ***REMOVED***
	Filename string
	Source   string

	stdlib []libFile
***REMOVED***

type VUContext struct ***REMOVED***
	r   *Runner
	ctx context.Context
	ch  chan runner.Result
	api map[string]map[string]interface***REMOVED******REMOVED***
***REMOVED***

type Module map[string]Member

type Member struct ***REMOVED***
	Func  interface***REMOVED******REMOVED***
	Async bool
***REMOVED***

type workerData struct ***REMOVED***
	ID        int64
	Test      loadtest.LoadTest
	Iteration int
***REMOVED***

func New(filename, src string) *Runner ***REMOVED***
	r := &Runner***REMOVED***
		Filename: filename,
		Source:   src,
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

func (r *Runner) Run(ctx context.Context, t loadtest.LoadTest, id int64) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		vu := VUContext***REMOVED***r: r, ctx: ctx, ch: ch, api: api.New()***REMOVED***
		w := v8worker.New(vu.Recv, vu.RecvSync)

		for _, f := range r.stdlib ***REMOVED***
			if err := w.Load(f.Filename, f.Source); err != nil ***REMOVED***
				log.WithError(err).WithField("file", f.Filename).Error("Couldn't load lib")
			***REMOVED***
		***REMOVED***

		wdata := workerData***REMOVED***
			ID:   id,
			Test: t,
		***REMOVED***
		wjson, err := json.Marshal(wdata)
		if err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't encode worker data")
			return
		***REMOVED***
		w.Load("internal:constants", fmt.Sprintf(`speedboat._data = %s;`, wjson))

		if err := vu.BridgeAPI(w); err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't register bridged functions")
			return
		***REMOVED***

		src := fmt.Sprintf(`
		function __run__() ***REMOVED***
			speedboat._data.Iteration++;
			try ***REMOVED***
		%s
			***REMOVED*** catch (e) ***REMOVED***
				console.error("Script Error", '' + e);
			***REMOVED***
		***REMOVED***
		`, r.Source)
		if err := w.Load(r.Filename, src); err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't load JS")
			return
		***REMOVED***

		done := make(chan interface***REMOVED******REMOVED***)
		for ***REMOVED***
			go func() ***REMOVED***
				w.SendSync("run")
				done <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***()

			select ***REMOVED***
			case <-done:
			case <-ctx.Done():
				w.TerminateExecution()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
