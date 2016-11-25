package main

import (
	"context"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/api"
	"github.com/loadimpact/speedboat/js"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/simple"
	"github.com/loadimpact/speedboat/stats"
	"github.com/loadimpact/speedboat/stats/influxdb"
	"gopkg.in/urfave/cli.v1"
	"net/url"
	"os"
	"os/signal"
	"sort"
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
			Name:  "run, r",
			Usage: "start test immediately",
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
		cli.BoolFlag***REMOVED***
			Name:  "quit-on-taint",
			Usage: "quit immediately if the test gets tainted",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "out, o",
			Usage: "output metrics to an external data store",
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
		return js.NewRunner(rt, exports)
	default:
		return nil, ErrInvalidType
	***REMOVED***
***REMOVED***

func parseCollectorString(s string) (t string, u *url.URL, err error) ***REMOVED***
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 ***REMOVED***
		return "", nil, errors.New("Malformed output; must be in the form 'type=url'")
	***REMOVED***

	u, err = url.Parse(parts[1])
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***

	return parts[0], u, nil
***REMOVED***

func makeCollector(s string) (stats.Collector, error) ***REMOVED***
	t, u, err := parseCollectorString(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch t ***REMOVED***
	case "influxdb":
		return influxdb.New(u)
	default:
		return nil, errors.New("Unknown output type: " + t)
	***REMOVED***
***REMOVED***

func actionRun(cc *cli.Context) error ***REMOVED***
	wg := sync.WaitGroup***REMOVED******REMOVED***

	args := cc.Args()
	if len(args) != 1 ***REMOVED***
		return cli.NewExitError("Wrong number of arguments!", 1)
	***REMOVED***

	// Collect CLI arguments, most (not all) relating to options.
	addr := cc.GlobalString("address")
	out := cc.String("out")
	cliOpts := lib.Options***REMOVED***
		Run:         cliBool(cc, "run"),
		VUs:         cliInt64(cc, "vus"),
		VUsMax:      cliInt64(cc, "vus-max"),
		Duration:    cliDuration(cc, "duration"),
		Quit:        cliBool(cc, "quit"),
		QuitOnTaint: cliBool(cc, "quit-on-taint"),
	***REMOVED***

	// Make the Runner, extract script-defined options.
	filename := args[0]
	runnerType := cc.String("type")
	runnerOpts := lib.Options***REMOVED******REMOVED***
	runner, err := makeRunner(filename, runnerType, &runnerOpts)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a runner")
		return err
	***REMOVED***

	// Make the metric collector, if requested.
	var collector stats.Collector
	if out != "" ***REMOVED***
		c, err := makeCollector(out)
		if err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't create output")
			return err
		***REMOVED***
		collector = c
	***REMOVED***

	// Option predecence: CLI > Script.
	// CLI has defaults, which are set as invalid, but have potentially nonzero values.
	opts := cliOpts.Apply(runnerOpts).SetAllValid(true)

	// Make the Engine
	engine, err := lib.NewEngine(runner)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create the engine")
		return err
	***REMOVED***
	engineC, engineCancel := context.WithCancel(context.Background())
	engine.Collector = collector

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
		if err := engine.Run(engineC, opts); err != nil ***REMOVED***
			log.WithError(err).Error("Engine Error")
		***REMOVED***
		engineCancel()
	***REMOVED***()
	go func() ***REMOVED***
		defer func() ***REMOVED***
			log.Debug("API Server terminated")
			wg.Done()
		***REMOVED***()
		log.WithField("addr", addr).Debug("API Server starting...")
		srv.Run(srvC, addr)
		srvCancel()
	***REMOVED***()

	log.Infof("Web UI available at: http://%s/", addr)

	// Wait for a signal or timeout before shutting down
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	log.Debug("Waiting for test to finish")
	select ***REMOVED***
	case <-srvC.Done():
		log.Debug("API server terminated; shutting down...")
	case <-engineC.Done():
		log.Debug("Engine terminated; shutting down...")
	case sig := <-signals:
		log.WithField("signal", sig).Debug("Signal received; shutting down...")
	***REMOVED***

	// If API server is still available, write final metrics to stdout.
	// (An unavailable API server most likely means a port binding failure.)
	select ***REMOVED***
	case <-srvC.Done():
	default:
		metricList, err := cl.Metrics()
		if err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't get metrics!")
			break
		***REMOVED***

		// Poor man's object sort.
		metrics := make(map[string]stats.Metric, len(metricList))
		keys := make([]string, len(metricList))
		for i, metric := range metricList ***REMOVED***
			metrics[metric.Name] = metric
			keys[i] = metric.Name
		***REMOVED***
		sort.Strings(keys)

		for _, key := range keys ***REMOVED***
			val := metrics[key].Humanize()
			if val == "0" ***REMOVED***
				continue
			***REMOVED***
			fmt.Printf("%s: %s\n", key, val)
		***REMOVED***
	***REMOVED***

	// Shut down the API server and engine, wait for them to terminate before exiting
	srvCancel()
	engineCancel()
	wg.Wait()

	if engine.Status.Tainted.Bool ***REMOVED***
		return cli.NewExitError("", 99)
	***REMOVED***
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
