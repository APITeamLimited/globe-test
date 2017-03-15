/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/loadimpact/k6/api"
	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/js2"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/simple"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/influxdb"
	"github.com/loadimpact/k6/stats/json"
	"github.com/loadimpact/k6/ui"
	"github.com/spf13/afero"
	"gopkg.in/guregu/null.v3"
	"gopkg.in/urfave/cli.v1"
)

const (
	TypeAuto = "auto"
	TypeURL  = "url"
	TypeJS   = "js"
	TypeJS2  = "js2"
)

var urlRegex = regexp.MustCompile(`(?i)^https?://`)

var commandRun = cli.Command***REMOVED***
	Name:      "run",
	Usage:     "Starts running a load test",
	ArgsUsage: "url|filename",
	Flags: []cli.Flag***REMOVED***
		cli.BoolFlag***REMOVED***
			Name:  "quiet, q",
			Usage: "hide the progress bar",
		***REMOVED***,
		cli.Int64Flag***REMOVED***
			Name:  "vus, u",
			Usage: "virtual users to simulate",
			Value: 1,
		***REMOVED***,
		cli.Int64Flag***REMOVED***
			Name:  "max, m",
			Usage: "max number of virtual users, if more than --vus",
		***REMOVED***,
		cli.DurationFlag***REMOVED***
			Name:  "duration, d",
			Usage: "test duration, 0 to run until cancelled",
		***REMOVED***,
		cli.Int64Flag***REMOVED***
			Name:  "iterations, i",
			Usage: "run a set number of iterations, multiplied by VU count",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "stage, s",
			Usage: "define a test stage, in the format time[:vus] (10s:100)",
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
		cli.Int64Flag***REMOVED***
			Name:  "max-redirects",
			Usage: "follow at most n redirects",
			Value: 10,
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "insecure-skip-tls-verify",
			Usage: "INSECURE: skip verification of TLS certificates",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "out, o",
			Usage: "output metrics to an external data store",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "config, c",
			Usage: "read additional config files",
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:   "no-usage-report",
			Usage:  "don't send heartbeat to k6 project on test execution",
			EnvVar: "K6_NO_USAGE_REPORT",
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

func guessType(data []byte) string ***REMOVED***
	if urlRegex.Match(data) ***REMOVED***
		return TypeURL
	***REMOVED***
	return TypeJS
***REMOVED***

func getSrcData(filename string, fs afero.Fs) (*lib.SourceData, error) ***REMOVED***
	reader := io.Reader(os.Stdin)
	if filename != "-" ***REMOVED***
		f, err := fs.Open(filename)
		if err != nil ***REMOVED***
			// If the file doesn't exist, but it looks like a URL, try using it as one.
			if os.IsNotExist(err) && urlRegex.MatchString(filename) ***REMOVED***
				return &lib.SourceData***REMOVED***
					Data:     []byte(filename),
					Filename: filename,
				***REMOVED***, nil
			***REMOVED***

			return nil, err
		***REMOVED***
		defer func() ***REMOVED*** _ = f.Close() ***REMOVED***()
		reader = f
	***REMOVED***

	data, err := ioutil.ReadAll(reader)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &lib.SourceData***REMOVED***
		Data:     data,
		Filename: filename,
	***REMOVED***, nil
***REMOVED***

func makeRunner(runnerType string, src *lib.SourceData, fs afero.Fs) (lib.Runner, error) ***REMOVED***
	switch runnerType ***REMOVED***
	case TypeAuto:
		return makeRunner(guessType(src.Data), src, fs)
	case TypeURL:
		u, err := url.Parse(strings.TrimSpace(string(src.Data)))
		if err != nil || u.Scheme == "" ***REMOVED***
			return nil, errors.New("Failed to parse URL")
		***REMOVED***
		r, err := simple.New(u)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return r, err
	case TypeJS:
		rt, err := js.New()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		exports, err := rt.Load(src, fs)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		r, err := js.NewRunner(rt, exports)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return r, nil
	case TypeJS2:
		return js2.New(src, fs)
	default:
		return nil, errors.New("Invalid type specified, see --help")
	***REMOVED***
***REMOVED***

func parseCollectorString(s string) (t, p string, err error) ***REMOVED***
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 ***REMOVED***
		return "", "", errors.New("Malformed output; must be in the form 'type=url'")
	***REMOVED***

	return parts[0], parts[1], nil
***REMOVED***

func makeCollector(s string, opts lib.Options) (lib.Collector, error) ***REMOVED***
	t, p, err := parseCollectorString(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch t ***REMOVED***
	case "influxdb":
		return influxdb.New(p, opts)
	case "json":
		return json.New(p, afero.NewOsFs(), opts)
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
	quiet := cc.Bool("quiet")
	cliOpts := lib.Options***REMOVED***
		Paused:                cliBool(cc, "paused"),
		VUs:                   cliInt64(cc, "vus"),
		VUsMax:                cliInt64(cc, "max"),
		Duration:              cliDuration(cc, "duration"),
		Iterations:            cliInt64(cc, "iterations"),
		Linger:                cliBool(cc, "linger"),
		MaxRedirects:          cliInt64(cc, "max-redirects"),
		InsecureSkipTLSVerify: cliBool(cc, "insecure-skip-tls-verify"),
		NoUsageReport:         cliBool(cc, "no-usage-report"),
	***REMOVED***
	for _, s := range cc.StringSlice("stage") ***REMOVED***
		stage, err := ParseStage(s)
		if err != nil ***REMOVED***
			log.WithError(err).Error("Invalid stage specified")
			return err
		***REMOVED***
		cliOpts.Stages = append(cliOpts.Stages, stage)
	***REMOVED***
	opts := cliOpts

	// Make the Runner, extract script-defined options.
	arg := args[0]
	fs := afero.NewOsFs()
	src, err := getSrcData(arg, fs)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Failed to parse input data")
		return err
	***REMOVED***
	runnerType := cc.String("type")
	if runnerType == TypeAuto ***REMOVED***
		runnerType = guessType(src.Data)
	***REMOVED***
	runner, err := makeRunner(runnerType, src, fs)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a runner")
		return err
	***REMOVED***
	opts = opts.Apply(runner.GetOptions())

	// Read config files.
	for _, filename := range cc.StringSlice("config") ***REMOVED***
		data, err := afero.ReadFile(fs, filename)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***

		var configOpts lib.Options
		if err := yaml.Unmarshal(data, &configOpts); err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		opts = opts.Apply(configOpts)
	***REMOVED***

	// CLI options override everything.
	opts = opts.Apply(cliOpts)

	// Default to 1 iteration if duration and stages are unspecified.
	if !opts.Duration.Valid && !opts.Iterations.Valid && len(opts.Stages) == 0 ***REMOVED***
		opts.Iterations = null.IntFrom(1)
	***REMOVED***

	// Apply defaults.
	opts = opts.SetAllValid(true)

	// Make sure VUsMax defaults to VUs if not specified.
	if opts.VUsMax.Int64 == 0 ***REMOVED***
		opts.VUsMax.Int64 = opts.VUs.Int64
		if len(opts.Stages) > 0 ***REMOVED***
			for _, stage := range opts.Stages ***REMOVED***
				if stage.Target.Valid && stage.Target.Int64 > opts.VUsMax.Int64 ***REMOVED***
					opts.VUsMax = stage.Target
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Update the runner's options.
	runner.ApplyOptions(opts)

	// Make the metric collector, if requested.
	var collector lib.Collector
	collectorString := "-"
	if out != "" ***REMOVED***
		c, err := makeCollector(out, opts)
		if err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't create output")
			return err
		***REMOVED***
		collector = c
		collectorString = fmt.Sprint(collector)
	***REMOVED***

	fmt.Fprintln(color.Output, "")

	color.Green(`          /\      |‾‾|  /‾‾/  /‾/   `)
	color.Green(`     /\  /  \     |  |_/  /  / /   `)
	color.Green(`    /  \/    \    |      |  /  ‾‾\  `)
	color.Green(`   /          \   |  |‾\  \ | (_) | `)
	color.Green(`  / __________ \  |__|  \__\ \___/  Welcome to k6 v%s!`, cc.App.Version)

	fmt.Fprintln(color.Output, "")

	fmt.Fprintf(color.Output, "  execution: %s\n", color.CyanString("local"))
	fmt.Fprintf(color.Output, "     output: %s\n", color.CyanString(collectorString))
	fmt.Fprintf(color.Output, "     script: %s (%s)\n", color.CyanString(src.Filename), color.CyanString(runnerType))
	fmt.Fprintf(color.Output, "\n")
	fmt.Fprintf(color.Output, "   duration: %s, iterations: %s\n", color.CyanString(opts.Duration.String), color.CyanString("%d", opts.Iterations.Int64))
	fmt.Fprintf(color.Output, "        vus: %s, max: %s\n", color.CyanString("%d", opts.VUs.Int64), color.CyanString("%d", opts.VUsMax.Int64))
	fmt.Fprintf(color.Output, "\n")
	fmt.Fprintf(color.Output, "    web ui: %s\n", color.CyanString("http://%s/", addr))
	fmt.Fprintf(color.Output, "\n")

	// Make the Engine
	engine, err := lib.NewEngine(runner, opts)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create the engine")
		return err
	***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	engine.Collector = collector

	// Send usage report, if we're allowed to
	if opts.NoUsageReport.Valid && !opts.NoUsageReport.Bool ***REMOVED***
		go func() ***REMOVED***
			conn, err := net.Dial("udp", "k6reports.loadimpact.com:6565")
			if err == nil ***REMOVED***
				// This is a best-effort attempt to send a usage report. We don't want
				// to inconvenience users if this doesn't work, for whatever reason
				_, _ = conn.Write([]byte("nyoom"))
				_ = conn.Close()
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// Run the engine.
	wg.Add(1)
	go func() ***REMOVED***
		defer func() ***REMOVED***
			log.Debug("Engine terminated")
			wg.Done()
		***REMOVED***()
		log.Debug("Starting engine...")
		if err := engine.Run(ctx); err != nil ***REMOVED***
			log.WithError(err).Error("Engine Error")
		***REMOVED***
		cancel()
	***REMOVED***()

	// Start the API server in the background.
	go func() ***REMOVED***
		if err := api.ListenAndServe(addr, engine); err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't start API server!")
		***REMOVED***
	***REMOVED***()

	// Progress bar for TTYs.
	progressBar := ui.ProgressBar***REMOVED***Width: 60***REMOVED***
	if isTTY && !quiet ***REMOVED***
		fmt.Fprintf(color.Output, " starting %s -- / --\r", progressBar.String())
	***REMOVED***

	// Wait for a signal or timeout before shutting down
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Print status at a set interval; less frequently on non-TTYs.
	tickInterval := 10 * time.Millisecond
	if !isTTY || quiet ***REMOVED***
		tickInterval = 1 * time.Second
	***REMOVED***
	ticker := time.NewTicker(tickInterval)

loop:
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			if !engine.IsRunning() ***REMOVED***
				break loop
			***REMOVED***

			statusString := "running"
			if engine.IsPaused() ***REMOVED***
				statusString = "paused"
			***REMOVED***

			atTime := engine.AtTime()
			totalTime := engine.TotalTime()
			progress := 0.0
			if totalTime > 0 ***REMOVED***
				progress = float64(atTime) / float64(totalTime)
			***REMOVED***

			if isTTY && !quiet ***REMOVED***
				progressBar.Progress = progress
				fmt.Fprintf(color.Output, "%10s %s %10s / %s\r",
					statusString,
					progressBar.String(),
					roundDuration(atTime, 100*time.Millisecond),
					roundDuration(totalTime, 100*time.Millisecond),
				)
			***REMOVED*** else ***REMOVED***
				fmt.Fprintf(color.Output, "[%-10s] %s / %s\n",
					statusString,
					roundDuration(atTime, 100*time.Millisecond),
					roundDuration(totalTime, 100*time.Millisecond),
				)
			***REMOVED***
		case <-ctx.Done():
			log.Debug("Engine terminated; shutting down...")
			break loop
		case sig := <-signals:
			log.WithField("signal", sig).Debug("Signal received; shutting down...")
			break loop
		***REMOVED***
	***REMOVED***

	// Shut down the API server and engine.
	cancel()
	wg.Wait()

	// Test done, leave that status as the final progress bar!
	atTime := engine.AtTime()
	if isTTY && !quiet ***REMOVED***
		progressBar.Progress = 1.0
		fmt.Fprintf(color.Output, "      done %s %10s / %s\n",
			progressBar.String(),
			roundDuration(atTime, 100*time.Millisecond),
			roundDuration(atTime, 100*time.Millisecond),
		)
	***REMOVED*** else ***REMOVED***
		fmt.Fprintf(color.Output, "[%-10s] %s / %s\n",
			"done",
			roundDuration(atTime, 100*time.Millisecond),
			roundDuration(atTime, 100*time.Millisecond),
		)
	***REMOVED***
	fmt.Fprintf(color.Output, "\n")

	// Print groups.
	var printGroup func(g *lib.Group, level int)
	printGroup = func(g *lib.Group, level int) ***REMOVED***
		indent := strings.Repeat("  ", level)

		if g.Name != "" && g.Parent != nil ***REMOVED***
			fmt.Fprintf(color.Output, "%s█ %s\n", indent, g.Name)
		***REMOVED***

		if len(g.Checks) > 0 ***REMOVED***
			if g.Name != "" && g.Parent != nil ***REMOVED***
				fmt.Fprintf(color.Output, "\n")
			***REMOVED***
			for _, check := range g.Checks ***REMOVED***
				icon := "✓"
				statusColor := color.GreenString
				if check.Fails > 0 ***REMOVED***
					icon = "✗"
					statusColor = color.RedString
				***REMOVED***
				fmt.Fprint(color.Output, statusColor("%s  %s %2.2f%% - %s\n",
					indent,
					icon,
					100*(float64(check.Passes)/float64(check.Passes+check.Fails)),
					check.Name,
				))
			***REMOVED***
			fmt.Fprintf(color.Output, "\n")
		***REMOVED***
		if len(g.Groups) > 0 ***REMOVED***
			if g.Name != "" && g.Parent != nil && len(g.Checks) > 0 ***REMOVED***
				fmt.Fprintf(color.Output, "\n")
			***REMOVED***
			for _, g := range g.Groups ***REMOVED***
				printGroup(g, level+1)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	printGroup(engine.Runner.GetDefaultGroup(), 1)

	// Sort and print metrics.
	metrics := make(map[string]*stats.Metric, len(engine.Metrics))
	metricNames := make([]string, 0, len(engine.Metrics))
	metricNameWidth := 0
	for m := range engine.Metrics ***REMOVED***
		metrics[m.Name] = m
		metricNames = append(metricNames, m.Name)
		if l := len(m.Name); l > metricNameWidth ***REMOVED***
			metricNameWidth = l
		***REMOVED***
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
		if m.Tainted.Valid ***REMOVED***
			if !m.Tainted.Bool ***REMOVED***
				icon = color.GreenString("✓")
			***REMOVED*** else ***REMOVED***
				icon = color.RedString("✗")
			***REMOVED***
		***REMOVED***

		// Hack some color in there.
		parts := strings.Split(val, ", ")
		newParts := make([]string, len(parts))
		for i, part := range parts ***REMOVED***
			kv := strings.SplitN(part, "=", 2)
			switch len(kv) ***REMOVED***
			case 1:
				newParts[i] = color.CyanString(kv[0])
			case 2:
				newParts[i] = fmt.Sprintf(
					"%s%s",
					color.New(color.Reset).Sprint(kv[0]+"="),
					color.CyanString(kv[1]),
				)
			***REMOVED***
		***REMOVED***
		val = strings.Join(newParts, ", ")

		namePadding := strings.Repeat(".", metricNameWidth-len(name)+3)
		fmt.Fprintf(color.Output, "  %s %s%s %s\n",
			icon,
			name,
			color.New(color.Faint).Sprint(namePadding+":"),
			color.CyanString(val),
		)
	***REMOVED***

	if opts.Linger.Bool ***REMOVED***
		<-signals
	***REMOVED***

	if engine.IsTainted() ***REMOVED***
		return cli.NewExitError("", 99)
	***REMOVED***
	return nil
***REMOVED***

func actionInspect(cc *cli.Context) error ***REMOVED***
	args := cc.Args()
	if len(args) != 1 ***REMOVED***
		return cli.NewExitError("Wrong number of arguments!", 1)
	***REMOVED***
	arg := args[0]

	fs := afero.NewOsFs()
	src, err := getSrcData(arg, fs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	runnerType := cc.String("type")
	if runnerType == TypeAuto ***REMOVED***
		runnerType = guessType(src.Data)
	***REMOVED***

	var opts lib.Options

	switch runnerType ***REMOVED***
	case TypeJS:
		r, err := js.New()
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***

		if _, err := r.Load(src, fs); err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		opts = opts.Apply(r.Options)
	***REMOVED***

	for _, filename := range cc.StringSlice("config") ***REMOVED***
		data, err := afero.ReadFile(fs, filename)
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
