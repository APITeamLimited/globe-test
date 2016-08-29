package main

import (
	"context"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/simple"
	"gopkg.in/urfave/cli.v1"
	"math"
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
		cli.DurationFlag***REMOVED***
			Name:  "duration, d",
			Usage: "test duration, 0 to run until cancelled",
			Value: 10 * time.Second,
		***REMOVED***,
		cli.Int64Flag***REMOVED***
			Name:  "prepare, p",
			Usage: "VUs to prepare (but not start)",
			Value: 0,
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "type, t",
			Usage: "input type, one of: auto, url, js",
			Value: "auto",
		***REMOVED***,
	***REMOVED***,
	Action: actionRun,
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

func makeRunner(filename, t string) (lib.Runner, error) ***REMOVED***
	if t == TypeAuto ***REMOVED***
		t = guessType(filename)
	***REMOVED***

	switch t ***REMOVED***
	case TypeAuto:
		return makeRunner(filename, t)
	case "":
		return nil, ErrUnknownType
	case TypeURL:
		return simple.New(filename)
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

	// Collect arguments
	addr := cc.GlobalString("address")

	duration := cc.Duration("duration")
	if duration == 0 ***REMOVED***
		duration = time.Duration(math.MaxInt64)
	***REMOVED***

	vus := cc.Int64("vus")

	prepared := cc.Int64("prepare")
	if prepared == 0 ***REMOVED***
		prepared = vus
	***REMOVED***

	// Make the Runner
	filename := args[0]
	runnerType := cc.String("type")
	runner, err := makeRunner(filename, runnerType)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a runner")
	***REMOVED***

	// Make the Engine
	engine := &lib.Engine***REMOVED***
		Runner: runner,
	***REMOVED***
	engineC, cancelEngine := context.WithCancel(context.Background())

	// Make the API Server
	api := &APIServer***REMOVED***
		Engine: engine,
		Cancel: cancelEngine,
		Info: lib.Info***REMOVED***
			Version: cc.App.Version,
		***REMOVED***,
	***REMOVED***
	apiC, cancelAPI := context.WithCancel(context.Background())

	// Make the Client
	cl, err := client.New(addr)
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
		log.WithField("prepared", prepared).Debug("Starting engine...")
		if err := engine.Run(engineC, prepared); err != nil ***REMOVED***
			log.WithError(err).Error("Engine Error")
		***REMOVED***
	***REMOVED***()
	go func() ***REMOVED***
		defer func() ***REMOVED***
			log.Debug("API Server terminated")
			wg.Done()
		***REMOVED***()
		log.WithField("addr", addr).Debug("API Server starting...")
		api.Run(apiC, addr)
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

	// Scale the test up to the desired VU count
	if vus > 0 ***REMOVED***
		log.WithField("vus", vus).Debug("Scaling test...")
		if err := cl.Scale(vus); err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't scale test")
		***REMOVED***
	***REMOVED***

	// Wait for a signal or timeout before shutting down
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	log.Debug("Waiting for test to finish")
	select ***REMOVED***
	case <-time.After(duration):
		log.Debug("Duration expired; shutting down...")
	case sig := <-quit:
		log.WithField("signal", sig).Debug("Signal received; shutting down...")
	***REMOVED***

	// Shut down the API server and engine, wait for them to terminate before exiting
	cancelAPI()
	cancelEngine()
	wg.Wait()

	return nil
***REMOVED***
