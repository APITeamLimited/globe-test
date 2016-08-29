package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"gopkg.in/urfave/cli.v1"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var commandRun = cli.Command***REMOVED***
	Name:      "run",
	Usage:     "Starts running a load test",
	ArgsUsage: "url|filename",
	Flags: []cli.Flag***REMOVED***
		cli.IntFlag***REMOVED***
			Name:  "vus, u",
			Usage: "virtual users to simulate",
			Value: 10,
		***REMOVED***,
		cli.DurationFlag***REMOVED***
			Name:  "duration, d",
			Usage: "test duration, 0 to run until cancelled",
			Value: 10 * time.Second,
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
		return "url"
	case strings.HasSuffix(filename, ".js"):
		return "js"
	default:
		return ""
	***REMOVED***
***REMOVED***

func makeRunner(filename, t string) (lib.Runner, error) ***REMOVED***
	if t == "auto" ***REMOVED***
		t = guessType(filename)
	***REMOVED***
	return nil, nil
***REMOVED***

func actionRun(cc *cli.Context) error ***REMOVED***
	wg := sync.WaitGroup***REMOVED******REMOVED***

	args := cc.Args()
	if len(args) != 1 ***REMOVED***
		return cli.NewExitError("Wrong number of arguments!", 1)
	***REMOVED***

	filename := args[0]
	runnerType := cc.String("type")
	runner, err := makeRunner(filename, runnerType)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a runner")
	***REMOVED***

	engine := &lib.Engine***REMOVED***
		Runner: runner,
	***REMOVED***
	engineC, cancelEngine := context.WithCancel(context.Background())

	api := &APIServer***REMOVED***
		Engine: engine,
		Cancel: cancelEngine,
		Info: lib.Info***REMOVED***
			Version: cc.App.Version,
		***REMOVED***,
	***REMOVED***
	apiC, cancelAPI := context.WithCancel(context.Background())

	timeout := cc.Duration("duration")
	if timeout > 0 ***REMOVED***
		engineC, _ = context.WithTimeout(engineC, timeout)
	***REMOVED***

	wg.Add(2)
	go func() ***REMOVED***
		defer func() ***REMOVED***
			log.Debug("Engine terminated")
			wg.Done()
		***REMOVED***()
		log.Debug("Starting engine...")
		if err := engine.Run(engineC); err != nil ***REMOVED***
			log.WithError(err).Error("Runtime Error")
		***REMOVED***
	***REMOVED***()
	go func() ***REMOVED***
		defer func() ***REMOVED***
			log.Debug("API Server terminated")
			wg.Done()
		***REMOVED***()

		addr := cc.GlobalString("address")
		log.WithField("addr", addr).Debug("API Server starting...")
		api.Run(apiC, addr)
	***REMOVED***()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	log.WithField("signal", sig).Debug("Signal received; shutting down...")

	cancelAPI()
	cancelEngine()
	wg.Wait()

	return nil
***REMOVED***
