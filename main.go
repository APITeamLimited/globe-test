package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/js"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/postman"
	"github.com/loadimpact/speedboat/simple"
	"github.com/loadimpact/speedboat/stats"
	"github.com/loadimpact/speedboat/stats/accumulate"
	"github.com/loadimpact/speedboat/stats/influxdb"
	"github.com/loadimpact/speedboat/stats/writer"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const (
	typeURL     = "url"
	typeJS      = "js"
	typePostman = "postman"
)

// Help text template
const helpTemplate = `NAME:
   ***REMOVED******REMOVED***.Name***REMOVED******REMOVED*** - ***REMOVED******REMOVED***.Usage***REMOVED******REMOVED***

USAGE:
   ***REMOVED******REMOVED***if .UsageText***REMOVED******REMOVED******REMOVED******REMOVED***.UsageText***REMOVED******REMOVED******REMOVED******REMOVED***else***REMOVED******REMOVED******REMOVED******REMOVED***.HelpName***REMOVED******REMOVED*** ***REMOVED******REMOVED***if .VisibleFlags***REMOVED******REMOVED***[options] ***REMOVED******REMOVED***end***REMOVED******REMOVED***filename|url***REMOVED******REMOVED***end***REMOVED******REMOVED***
   ***REMOVED******REMOVED***if .Version***REMOVED******REMOVED******REMOVED******REMOVED***if not .HideVersion***REMOVED******REMOVED***
VERSION:
   ***REMOVED******REMOVED***.Version***REMOVED******REMOVED***
   ***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .VisibleFlags***REMOVED******REMOVED***
OPTIONS:
   ***REMOVED******REMOVED***range .VisibleFlags***REMOVED******REMOVED******REMOVED******REMOVED***.***REMOVED******REMOVED***
   ***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
`

var mVUs = stats.Stat***REMOVED***Name: "vus", Type: stats.GaugeType***REMOVED***

func pollVURamping(ctx context.Context, t lib.Test) <-chan int ***REMOVED***
	ch := make(chan int)
	startTime := time.Now()

	go func() ***REMOVED***
		defer close(ch)

		ticker := time.NewTicker(1 * time.Second)
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				ch <- t.VUsAt(time.Since(startTime))
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***

func parseBackend(out string) (stats.Backend, error) ***REMOVED***
	switch ***REMOVED***
	case strings.HasPrefix(out, "influxdb+"):
		url := strings.TrimPrefix(out, "influxdb+")
		return influxdb.NewFromURL(url)
	default:
		return nil, errors.New("Unknown output destination")
	***REMOVED***
***REMOVED***

func parseStages(vus []string, total time.Duration) (stages []lib.TestStage, err error) ***REMOVED***
	if len(vus) == 0 ***REMOVED***
		return []lib.TestStage***REMOVED***
			lib.TestStage***REMOVED***Duration: total, StartVUs: 10, EndVUs: 10***REMOVED***,
		***REMOVED***, nil
	***REMOVED***

	accountedTime := time.Duration(0)
	fluidStages := []int***REMOVED******REMOVED***
	for i, spec := range vus ***REMOVED***
		parts := strings.SplitN(spec, ":", 2)
		countParts := strings.SplitN(parts[0], "-", 2)

		stage := lib.TestStage***REMOVED******REMOVED***

		// An absent first part means keep going from the last stage's end
		// If it's the first stage, just start with 0
		if countParts[0] == "" ***REMOVED***
			if i > 0 ***REMOVED***
				stage.StartVUs = stages[i-1].EndVUs
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			start, err := strconv.ParseInt(countParts[0], 10, 64)
			if err != nil ***REMOVED***
				return stages, err
			***REMOVED***
			stage.StartVUs = int(start)
		***REMOVED***

		// If an end is specified, use that, otherwise keep the VU level constant
		if len(countParts) > 1 && countParts[1] != "" ***REMOVED***
			end, err := strconv.ParseInt(countParts[1], 10, 64)
			if err != nil ***REMOVED***
				return stages, err
			***REMOVED***
			stage.EndVUs = int(end)
		***REMOVED*** else ***REMOVED***
			stage.EndVUs = stage.StartVUs
		***REMOVED***

		// If a time is specified, use that, otherwise mark the stage as "fluid", allotting it an
		// even slice of what time remains after all fixed stages are accounted for (may be 0)
		if len(parts) > 1 ***REMOVED***
			duration, err := time.ParseDuration(parts[1])
			if err != nil ***REMOVED***
				return stages, err
			***REMOVED***
			stage.Duration = duration
			accountedTime += duration
		***REMOVED*** else ***REMOVED***
			fluidStages = append(fluidStages, i)
		***REMOVED***

		stages = append(stages, stage)
	***REMOVED***

	// We're ignoring fluid stages if the fixed stages already take up all the allotted time
	// Otherwise, evenly divide the test's remaining time between all fluid stages
	if len(fluidStages) > 0 && accountedTime < total ***REMOVED***
		fluidDuration := (total - accountedTime) / time.Duration(len(fluidStages))
		for _, i := range fluidStages ***REMOVED***
			stage := stages[i]
			stage.Duration = fluidDuration
			stages[i] = stage
		***REMOVED***
	***REMOVED***

	return stages, nil
***REMOVED***

func parseTags(lines []string) stats.Tags ***REMOVED***
	tags := make(stats.Tags)
	for _, line := range lines ***REMOVED***
		idx := strings.IndexAny(line, ":=")
		if idx == -1 ***REMOVED***
			tags[line] = line
			continue
		***REMOVED***

		key := line[:idx]
		val := line[idx+1:]
		if key == "" ***REMOVED***
			key = val
		***REMOVED***
		tags[key] = val
	***REMOVED***
	return tags
***REMOVED***

func guessType(arg string) string ***REMOVED***
	switch ***REMOVED***
	case strings.Contains(arg, "://"):
		return typeURL
	case strings.HasSuffix(arg, ".js"):
		return typeJS
	case strings.HasSuffix(arg, ".postman_collection.json"):
		return typePostman
	***REMOVED***
	return ""
***REMOVED***

func readAll(filename string) ([]byte, error) ***REMOVED***
	if filename == "-" ***REMOVED***
		return ioutil.ReadAll(os.Stdin)
	***REMOVED***

	return ioutil.ReadFile(filename)
***REMOVED***

func makeRunner(t lib.Test, filename, typ string) (lib.Runner, error) ***REMOVED***
	if typ == typeURL ***REMOVED***
		return simple.New(filename), nil
	***REMOVED***

	bytes, err := readAll(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch typ ***REMOVED***
	case typeJS:
		return js.New(filename, string(bytes)), nil
	case typePostman:
		return postman.New(bytes)
	default:
		return nil, errors.New("Type ambiguous, please specify -t/--type")
	***REMOVED***
***REMOVED***

func action(cc *cli.Context) error ***REMOVED***
	once := cc.Bool("once")

	if cc.IsSet("verbose") ***REMOVED***
		log.SetLevel(log.DebugLevel)
	***REMOVED***

	for _, out := range cc.StringSlice("out") ***REMOVED***
		backend, err := parseBackend(out)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		stats.DefaultRegistry.Backends = append(stats.DefaultRegistry.Backends, backend)
	***REMOVED***

	var formatter writer.Formatter
	switch cc.String("format") ***REMOVED***
	case "":
	case "json":
		formatter = writer.JSONFormatter***REMOVED******REMOVED***
	case "prettyjson":
		formatter = writer.PrettyJSONFormatter***REMOVED******REMOVED***
	case "yaml":
		formatter = writer.YAMLFormatter***REMOVED******REMOVED***
	default:
		return cli.NewExitError("Unknown output format", 1)
	***REMOVED***

	stats.DefaultRegistry.ExtraTags = parseTags(cc.StringSlice("tag"))

	var summarizer *Summarizer
	if formatter != nil ***REMOVED***
		filter := stats.MakeFilter(cc.StringSlice("exclude"), cc.StringSlice("select"))
		if cc.Bool("raw") ***REMOVED***
			backend := &writer.Backend***REMOVED***
				Writer:    os.Stdout,
				Formatter: formatter,
			***REMOVED***
			backend.Filter = filter
			stats.DefaultRegistry.Backends = append(stats.DefaultRegistry.Backends, backend)
		***REMOVED*** else ***REMOVED***
			accumulator := accumulate.New()
			accumulator.Filter = filter
			accumulator.GroupBy = cc.StringSlice("group-by")
			stats.DefaultRegistry.Backends = append(stats.DefaultRegistry.Backends, accumulator)

			summarizer = &Summarizer***REMOVED***
				Accumulator: accumulator,
				Formatter:   formatter,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	stages, err := parseStages(cc.StringSlice("vus"), cc.Duration("duration"))
	if err != nil ***REMOVED***
		return cli.NewExitError(err.Error(), 1)
	***REMOVED***
	if once ***REMOVED***
		stages = []lib.TestStage***REMOVED***
			lib.TestStage***REMOVED***Duration: 0, StartVUs: 1, EndVUs: 1***REMOVED***,
		***REMOVED***
	***REMOVED***
	t := lib.Test***REMOVED***Stages: stages***REMOVED***

	var r lib.Runner
	switch len(cc.Args()) ***REMOVED***
	case 0:
		cli.ShowAppHelp(cc)
		return nil
	case 1:
		filename := cc.Args()[0]
		typ := cc.String("type")
		if typ == "" ***REMOVED***
			typ = guessType(filename)
		***REMOVED***

		if filename == "-" && typ == "" ***REMOVED***
			typ = typeJS
		***REMOVED***

		runner, err := makeRunner(t, filename, typ)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		r = runner
	default:
		return cli.NewExitError("Too many arguments!", 1)
	***REMOVED***

	if cc.Bool("plan") ***REMOVED***
		data, err := yaml.Marshal(map[string]interface***REMOVED******REMOVED******REMOVED***
			"stages": stages,
		***REMOVED***)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		os.Stdout.Write(data)
		return nil
	***REMOVED***

	vus := lib.VUGroup***REMOVED***
		Pool: lib.VUPool***REMOVED***
			New: r.NewVU,
		***REMOVED***,
		RunOnce: func(ctx context.Context, vu lib.VU) ***REMOVED***
			if err := vu.RunOnce(ctx); err != nil ***REMOVED***
				log.WithError(err).Error("Uncaught Error")
			***REMOVED***
		***REMOVED***,
	***REMOVED***

	for i := 0; i < t.MaxVUs(); i++ ***REMOVED***
		vu, err := vus.Pool.New()
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		vus.Pool.Put(vu)
	***REMOVED***

	ctx, cancel := context.WithTimeout(context.Background(), t.TotalDuration())
	if once ***REMOVED***
		ctx, cancel = context.WithCancel(context.Background())
	***REMOVED***

	go func() ***REMOVED***
		ticker := time.NewTicker(1 * time.Second)
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				if err := stats.Submit(); err != nil ***REMOVED***
					log.WithError(err).Error("[Couldn't submit stats]")
				***REMOVED***
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	interval := cc.Duration("interval")
	if interval > 0 && summarizer != nil ***REMOVED***
		go func() ***REMOVED***
			ticker := time.NewTicker(interval)
			for ***REMOVED***
				select ***REMOVED***
				case <-ticker.C:
					if err := summarizer.Print(os.Stdout); err != nil ***REMOVED***
						log.WithError(err).Error("Couldn't print statistics!")
					***REMOVED***
				case <-ctx.Done():
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	go func() ***REMOVED***
		quit := make(chan os.Signal)
		signal.Notify(quit)

		select ***REMOVED***
		case <-quit:
			cancel()
		case <-ctx.Done():
		***REMOVED***
	***REMOVED***()

	if !cc.Bool("quiet") ***REMOVED***
		log.WithFields(log.Fields***REMOVED***
			"at":     time.Now(),
			"length": t.TotalDuration(),
		***REMOVED***).Info("Starting test...")
	***REMOVED***

	if once ***REMOVED***
		stats.Add(stats.Sample***REMOVED***Stat: &mVUs, Values: stats.Value(1)***REMOVED***)

		vu, _ := vus.Pool.Get()
		if err := vu.RunOnce(ctx); err != nil ***REMOVED***
			log.WithError(err).Error("Uncaught Error")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		vus.Start(ctx)
		scaleTo := pollVURamping(ctx, t)
	mainLoop:
		for ***REMOVED***
			select ***REMOVED***
			case num := <-scaleTo:
				vus.Scale(num)
				stats.Add(stats.Sample***REMOVED***
					Stat:   &mVUs,
					Values: stats.Value(float64(num)),
				***REMOVED***)
			case <-ctx.Done():
				break mainLoop
			***REMOVED***
		***REMOVED***

		vus.Stop()
	***REMOVED***

	stats.Add(stats.Sample***REMOVED***Stat: &mVUs, Values: stats.Value(0)***REMOVED***)
	stats.Submit()

	if summarizer != nil ***REMOVED***
		if err := summarizer.Print(os.Stdout); err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't print statistics!")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func main() ***REMOVED***
	// Submit usage statistics for the closed beta
	invocation := Invocation***REMOVED******REMOVED***
	invocationError := make(chan error, 1)

	go func() ***REMOVED***
		// Set SUBMIT=false to prevent stat collection
		submitURL := os.Getenv("SB_SUBMIT")
		switch submitURL ***REMOVED***
		case "false", "no":
			return
		case "":
			submitURL = "http://52.209.216.227:8080"
		***REMOVED***

		// Wait at most 2s for an invocation error to be reported
		select ***REMOVED***
		case err := <-invocationError:
			invocation.Error = err.Error()
		case <-time.After(2 * time.Second):
		***REMOVED***

		// Submit stats to a specified server
		if err := invocation.Submit(submitURL); err != nil ***REMOVED***
			log.WithError(err).Debug("Couldn't submit statistics")
		***REMOVED***
	***REMOVED***()

	// Free up -v and -h for our own flags
	cli.VersionFlag.Name = "version"
	cli.HelpFlag.Name = "help, ?"
	cli.AppHelpTemplate = helpTemplate

	// Bootstrap the app from commandline flags
	app := cli.NewApp()
	app.Name = "speedboat"
	app.Usage = "A next-generation load generator"
	app.Version = "1.0.0-mvp1"
	app.Flags = []cli.Flag***REMOVED***
		cli.BoolFlag***REMOVED***
			Name:  "plan",
			Usage: "Don't run anything, just show the test plan",
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "once",
			Usage: "Run only a single test iteration, with one VU",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "type, t",
			Usage: "Input file type, if not evident (url, js or postman)",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "vus, u",
			Usage: "Number of VUs to simulate",
		***REMOVED***,
		cli.DurationFlag***REMOVED***
			Name:  "duration, d",
			Usage: "Test duration",
			Value: time.Duration(10) * time.Second,
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "verbose, v",
			Usage: "More verbose output",
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "quiet, q",
			Usage: "Suppress the summary at the end of a test",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "format, f",
			Usage: "Format for printed metrics (yaml, json, prettyjson)",
			Value: "yaml",
		***REMOVED***,
		cli.DurationFlag***REMOVED***
			Name:  "interval, i",
			Usage: "Write periodic summaries",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "out, o",
			Usage: "Write metrics to a database",
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:  "raw",
			Usage: "Instead of summaries, dump raw samples to stdout",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "select, s",
			Usage: "Include only named metrics",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "exclude, e",
			Usage: "Exclude named metrics",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "group-by, g",
			Usage: "Group metrics by tags",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "tag",
			Usage: "Additional metric tags",
		***REMOVED***,
	***REMOVED***
	app.Before = func(cc *cli.Context) error ***REMOVED***
		invocation.PopulateWithContext(cc)
		return nil
	***REMOVED***
	app.Action = action
	invocationError <- app.Run(os.Args)
***REMOVED***
