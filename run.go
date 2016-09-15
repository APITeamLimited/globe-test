package main

import (
	"context"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/api"
	"github.com/loadimpact/speedboat/js"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/simple"
	"gopkg.in/guregu/null.v3"
	"gopkg.in/urfave/cli.v1"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	TypeAuto = "auto"
	TypeURL  = "url"
	TypeJS   = "js"
)

var (
	ErrUnknownType = errors.New("Unable to infer type from argument; specify with -t/--type")
	ErrInvalidType = errors.New("Invalid type specified, see --help")
)

var commandRun = cli.Command***REMOVED***
	Name:      "run",
	Usage:     "Starts running a load test",
	ArgsUsage: "url|filename",
	Flags: []cli.Flag***REMOVED***
		cli.Int64Flag***REMOVED***
			Name:  "vus, u",
			Usage: "virtual users to simulate",
			Value: 10,
		***REMOVED***,
		cli.Int64Flag***REMOVED***
			Name:  "max, m",
			Usage: "max number of virtual users, if more than --vus",
		***REMOVED***,
		cli.DurationFlag***REMOVED***
			Name:  "duration, d",
			Usage: "test duration, 0 to run until cancelled",
			Value: 10 * time.Second,
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "paused, p",
			Usage: "start test in a paused state",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "type, t",
			Usage: "input type, one of: auto, url, js",
			Value: "auto",
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "quit, q",
			Usage: "quit immediately on test completion",
		***REMOVED***,
	***REMOVED***,
	Action: actionRun,
	Description: `Run starts a load test.

   This is the main entry point to Speedboat, and will do two things:
   
   - Construct an Engine and provide it with a Runner, depending on the first
     argument and the --type flag, which is used to execute the test.
   
   - Start an a web server on the address specified by the global --address
     flag, which serves a web interface and a REST API for remote control.
   
   For ease of use, you may also pass initial status parameters (vus, max,
   duration) to 'run', which will be applied through a normal API call.`,
***REMOVED***

var commandInspect = cli.Command***REMOVED***
	Name:      "inspect",
	Aliases:   []string***REMOVED***"i"***REMOVED***,
	Usage:     "Merges and prints test configuration",
	ArgsUsage: "url|filename",
	Flags: []cli.Flag***REMOVED***
		cli.StringFlag***REMOVED***
			Name:  "type, t",
			Usage: "input type, one of: auto, url, js",
			Value: "auto",
		***REMOVED***,
		cli.Int64Flag***REMOVED***
			Name:  "vus, u",
			Usage: "override vus",
			Value: 10,
		***REMOVED***,
		cli.Int64Flag***REMOVED***
			Name:  "max, m",
			Usage: "override vus-max",
		***REMOVED***,
		cli.DurationFlag***REMOVED***
			Name:  "duration, d",
			Usage: "override duration",
		***REMOVED***,
	***REMOVED***,
	Action: actionInspect,
***REMOVED***

func guessType(filename string) string ***REMOVED***
	switch ***REMOVED***
	case strings.Contains(filename, "://"):
		return TypeURL
	case strings.HasSuffix(filename, ".js"):
		return TypeJS
	default:
		return ""
	***REMOVED***
***REMOVED***

func makeRunner(filename, t string, opts *lib.Options) (lib.Runner, error) ***REMOVED***
	if t == TypeAuto ***REMOVED***
		t = guessType(filename)
	***REMOVED***

	switch t ***REMOVED***
	case "":
		return nil, ErrUnknownType
	case TypeURL:
		return simple.New(filename)
	case TypeJS:
		rt, err := js.New()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		exports, err := rt.Load(filename)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if err := rt.ExtractOptions(exports, opts); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return &js.Runner***REMOVED***Runtime: rt, Exports: exports***REMOVED***, nil
	default:
		return nil, ErrInvalidType
	***REMOVED***
***REMOVED***

func actionRun(cc *cli.Context) error ***REMOVED***
	wg := sync.WaitGroup***REMOVED******REMOVED***

	args := cc.Args()
	if len(args) != 1 ***REMOVED***
		return cli.NewExitError("Wrong number of arguments!", 1)
	***REMOVED***

	// Make the Runner
	filename := args[0]
	runnerType := cc.String("type")
	opts := lib.Options***REMOVED******REMOVED***
	runner, err := makeRunner(filename, runnerType, &opts)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a runner")
		return err
	***REMOVED***

	// Collect arguments
	addr := cc.GlobalString("address")
	paused := cc.Bool("paused")
	quit := cc.Bool("quit")

	duration := opts.Duration
	if cc.IsSet("duration") ***REMOVED***
		duration = cc.Duration("duration")
	***REMOVED***

	vus := opts.VUs
	if cc.IsSet("vus") ***REMOVED***
		vus = cc.Int64("vus")
	***REMOVED***

	max := opts.VUsMax
	if cc.IsSet("max") ***REMOVED***
		max = cc.Int64("max")
	***REMOVED***
	if max == 0 ***REMOVED***
		max = vus
	***REMOVED***

	if vus > max ***REMOVED***
		return cli.NewExitError(lib.ErrTooManyVUs.Error(), 1)
	***REMOVED***

	// Make the Engine
	engine, err := lib.NewEngine(runner)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create the engine")
		return err
	***REMOVED***
	engineC, engineCancel := context.WithCancel(context.Background())

	// Make the API Server
	srv := &api.Server***REMOVED***
		Engine: engine,
		Info:   lib.Info***REMOVED***Version: cc.App.Version***REMOVED***,
	***REMOVED***
	srvC, srvCancel := context.WithCancel(context.Background())

	// Make the Client
	cl, err := api.NewClient(addr)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't make a client; is the address valid?")
		return err
	***REMOVED***

	// Run the engine and API server in the background
	wg.Add(2)
	go func() ***REMOVED***
		defer func() ***REMOVED***
			log.Debug("Engine terminated")
			wg.Done()
		***REMOVED***()
		log.Debug("Starting engine...")
		if err := engine.Run(engineC); err != nil ***REMOVED***
			log.WithError(err).Error("Engine Error")
		***REMOVED***
	***REMOVED***()
	go func() ***REMOVED***
		defer func() ***REMOVED***
			log.Debug("API Server terminated")
			wg.Done()
		***REMOVED***()
		log.WithField("addr", addr).Debug("API Server starting...")
		srv.Run(srvC, addr)
	***REMOVED***()

	// Wait for the API server to come online
	startTime := time.Now()
	for ***REMOVED***
		if err := cl.Ping(); err != nil ***REMOVED***
			if time.Since(startTime) < 1*time.Second ***REMOVED***
				log.WithError(err).Debug("Waiting for API server to start...")
				time.Sleep(1 * time.Millisecond)
			***REMOVED*** else ***REMOVED***
				log.WithError(err).Warn("Connection to API server failed; retrying...")
				time.Sleep(1 * time.Second)
			***REMOVED***
			continue
		***REMOVED***
		break
	***REMOVED***

	// Start the test with the desired state
	log.WithField("vus", vus).Debug("Starting test...")
	status := lib.Status***REMOVED***
		Running: null.BoolFrom(!paused),
		VUs:     null.IntFrom(vus),
		VUsMax:  null.IntFrom(max),
	***REMOVED***
	if _, err := cl.UpdateStatus(status); err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't scale test")
	***REMOVED***

	// Pause the test once the duration expires
	if duration > 0 ***REMOVED***
		log.WithField("duration", duration).Debug("Test will pause after...")
		go func() ***REMOVED***
			time.Sleep(duration)
			log.Debug("Duration expired, pausing...")
			status := lib.Status***REMOVED***Running: null.BoolFrom(false)***REMOVED***
			if _, err := cl.UpdateStatus(status); err != nil ***REMOVED***
				log.WithError(err).Error("Couldn't pause test")
			***REMOVED***

			if quit ***REMOVED***
				log.Debug("Quit requested, terminating...")
				srvCancel()
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// Wait for a signal or timeout before shutting down
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	log.Debug("Waiting for test to finish")
	select ***REMOVED***
	case <-srvC.Done():
		log.Debug("API server terminated; shutting down...")
	case sig := <-signals:
		log.WithField("signal", sig).Debug("Signal received; shutting down...")
	***REMOVED***

	// Shut down the API server and engine, wait for them to terminate before exiting
	srvCancel()
	engineCancel()
	wg.Wait()

	return nil
***REMOVED***

func actionInspect(cc *cli.Context) error ***REMOVED***
	args := cc.Args()
	if len(args) != 1 ***REMOVED***
		return cli.NewExitError("Wrong number of arguments!", 1)
	***REMOVED***
	filename := args[0]

	t := cc.String("type")
	if t == TypeAuto ***REMOVED***
		t = guessType(filename)
	***REMOVED***

	var opts lib.Options
	switch t ***REMOVED***
	case TypeJS:
		r, err := js.New()
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***

		exports, err := r.Load(filename)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		if err := r.ExtractOptions(exports, &opts); err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
	***REMOVED***

	return dumpYAML(opts)
***REMOVED***
