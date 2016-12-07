package main

import (
	"context"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/loadimpact/k6/api"
	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/simple"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/influxdb"
	"github.com/loadimpact/k6/ui"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
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
		cli.Float64Flag***REMOVED***
			Name:  "acceptance, a",
			Usage: "acceptable margin of error before failing the test",
			Value: 0.0,
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
			Name:  "linger, l",
			Usage: "linger after test completion",
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "abort-on-taint",
			Usage: "abort immediately if the test gets tainted",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "out, o",
			Usage: "output metrics to an external data store",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "config, c",
			Usage: "read additional config files",
		***REMOVED***,
	***REMOVED***,
	Action: actionRun,
	Description: `Run starts a load test.

   This is the main entry point to k6, and will do two things:
   
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
		cli.StringSliceFlag***REMOVED***
			Name:  "config, c",
			Usage: "read additional config files",
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

func makeRunner(filename, t string) (lib.Runner, error) ***REMOVED***
	if t == TypeAuto ***REMOVED***
		t = guessType(filename)
	***REMOVED***

	switch t ***REMOVED***
	case "":
		return nil, ErrUnknownType
	case TypeURL:
		r, err := simple.New(filename)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return r, err
	case TypeJS:
		rt, err := js.New()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		exports, err := rt.Load(filename)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		r, err := js.NewRunner(rt, exports)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return r, nil
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
	opts := lib.Options***REMOVED***
		Paused:       cliBool(cc, "paused"),
		VUs:          cliInt64(cc, "vus"),
		VUsMax:       cliInt64(cc, "max"),
		Duration:     cliDuration(cc, "duration"),
		Linger:       cliBool(cc, "linger"),
		AbortOnTaint: cliBool(cc, "abort-on-taint"),
		Acceptance:   cliFloat64(cc, "acceptance"),
	***REMOVED***

	// Make the Runner, extract script-defined options.
	filename := args[0]
	runnerType := cc.String("type")
	runner, err := makeRunner(filename, runnerType)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a runner")
		return err
	***REMOVED***
	opts = opts.Apply(runner.GetOptions())

	// Read config files.
	for _, filename := range cc.StringSlice("config") ***REMOVED***
		data, err := ioutil.ReadFile(filename)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***

		var configOpts lib.Options
		if err := yaml.Unmarshal(data, &configOpts); err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		opts = opts.Apply(configOpts)
	***REMOVED***

	// CLI options have defaults, which are set as invalid, but have potentially nonzero values.
	// Flipping the Valid flag for all invalid options thus applies all defaults.
	if !opts.VUsMax.Valid ***REMOVED***
		opts.VUsMax.Int64 = opts.VUs.Int64
	***REMOVED***
	opts = opts.SetAllValid(true)
	runner.ApplyOptions(opts)

	// Make the metric collector, if requested.
	var collector stats.Collector
	collectorString := "-"
	if out != "" ***REMOVED***
		c, err := makeCollector(out)
		if err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't create output")
			return err
		***REMOVED***
		collector = c
		collectorString = fmt.Sprint(collector)
	***REMOVED***

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

	// Print the banner!
	fmt.Printf("Welcome to k6 v%s!\n", cc.App.Version)
	fmt.Printf("\n")
	fmt.Printf("  execution: local\n")
	fmt.Printf("     output: %s\n", collectorString)
	fmt.Printf("     script: %s\n", filename)
	fmt.Printf("             ↳ duration: %s\n", opts.Duration.String)
	fmt.Printf("             ↳ vus: %d, max: %d\n", opts.VUs.Int64, opts.VUsMax.Int64)
	fmt.Printf("\n")
	fmt.Printf("  web ui: http://%s/\n", addr)
	fmt.Printf("\n")

	progressBar := ui.ProgressBar***REMOVED***Width: 60***REMOVED***
	fmt.Printf(" starting %s -- / --\r", progressBar.String())

	// Wait for a signal or timeout before shutting down
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(10 * time.Millisecond)

loop:
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			statusString := "running"
			if !engine.Status.Running.Bool ***REMOVED***
				if engine.IsRunning() ***REMOVED***
					statusString = "paused"
				***REMOVED*** else ***REMOVED***
					statusString = "stopping"
				***REMOVED***
			***REMOVED***

			atTime := time.Duration(engine.Status.AtTime.Int64)
			totalTime, finite := engine.TotalTime()
			progress := 0.0
			if finite ***REMOVED***
				progress = float64(atTime) / float64(totalTime)
			***REMOVED***

			progressBar.Progress = progress
			fmt.Printf("%10s %s %10s / %s\r",
				statusString,
				progressBar.String(),
				roundDuration(atTime, 100*time.Millisecond),
				roundDuration(totalTime, 100*time.Millisecond),
			)
		case <-srvC.Done():
			log.Debug("API server terminated; shutting down...")
			break loop
		case <-engineC.Done():
			log.Debug("Engine terminated; shutting down...")
			break loop
		case sig := <-signals:
			log.WithField("signal", sig).Debug("Signal received; shutting down...")
			break loop
		***REMOVED***
	***REMOVED***

	// Shut down the API server and engine.
	srvCancel()
	engineCancel()
	wg.Wait()

	// Test done, leave that status as the final progress bar!
	atTime := time.Duration(engine.Status.AtTime.Int64)
	progressBar.Progress = 1.0
	fmt.Printf("      done %s %10s / %s\n",
		progressBar.String(),
		roundDuration(atTime, 100*time.Millisecond),
		roundDuration(atTime, 100*time.Millisecond),
	)
	fmt.Printf("\n")

	// Print groups.
	var printGroup func(g *lib.Group, level int)
	printGroup = func(g *lib.Group, level int) ***REMOVED***
		indent := strings.Repeat("  ", level)

		if g.Name != "" && g.Parent != nil ***REMOVED***
			fmt.Printf("%s█ %s\n", indent, g.Name)
		***REMOVED***

		if len(g.Checks) > 0 ***REMOVED***
			if g.Name != "" && g.Parent != nil ***REMOVED***
				fmt.Printf("\n")
			***REMOVED***
			for _, check := range g.Checks ***REMOVED***
				icon := "✓"
				if check.Fails > 0 ***REMOVED***
					icon = "✗"
				***REMOVED***
				fmt.Printf("%s  %s %2.2f%% - %s\n",
					indent,
					icon,
					100*(float64(check.Passes)/float64(check.Passes+check.Fails)),
					check.Name,
				)
			***REMOVED***
			fmt.Printf("\n")
		***REMOVED***
		if len(g.Groups) > 0 ***REMOVED***
			if g.Name != "" && g.Parent != nil && len(g.Checks) > 0 ***REMOVED***
				fmt.Printf("\n")
			***REMOVED***
			for _, g := range g.Groups ***REMOVED***
				printGroup(g, level+1)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	groups := engine.Runner.GetGroups()
	for _, g := range groups ***REMOVED***
		if g.Parent != nil ***REMOVED***
			continue
		***REMOVED***
		printGroup(g, 1)
	***REMOVED***

	// Sort and print metrics.
	metrics := make(map[string]*stats.Metric, len(engine.Metrics))
	metricNames := make([]string, 0, len(engine.Metrics))
	for m, _ := range engine.Metrics ***REMOVED***
		metrics[m.Name] = m
		metricNames = append(metricNames, m.Name)
	***REMOVED***
	sort.Strings(metricNames)

	for _, name := range metricNames ***REMOVED***
		m := metrics[name]
		m.Sample = engine.Metrics[m].Format()
		val := metrics[name].Humanize()
		if val == "0" ***REMOVED***
			continue
		***REMOVED***
		icon := " "
		for _, threshold := range engine.Thresholds[name] ***REMOVED***
			icon = "✓"
			if threshold.Failed ***REMOVED***
				icon = "✗"
				break
			***REMOVED***
		***REMOVED***
		fmt.Printf("  %s %s: %s\n", icon, name, val)
	***REMOVED***

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

		if _, err := r.Load(filename); err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		opts = opts.Apply(r.Options)
	***REMOVED***

	for _, filename := range cc.StringSlice("config") ***REMOVED***
		data, err := ioutil.ReadFile(filename)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***

		var configOpts lib.Options
		if err := yaml.Unmarshal(data, &configOpts); err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		opts = opts.Apply(configOpts)
	***REMOVED***

	return dumpYAML(opts)
***REMOVED***
