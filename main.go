package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/js"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/simple"
	"github.com/loadimpact/speedboat/stats"
	"github.com/loadimpact/speedboat/stats/accumulate"
	"github.com/loadimpact/speedboat/stats/influxdb"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const (
	typeURL = "url"
	typeJS  = "js"
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
	case out == "-":
		return stats.NewJSONBackend(os.Stdout), nil
	case strings.HasPrefix(out, "influxdb+"):
		url := strings.TrimPrefix(out, "influxdb+")
		return influxdb.NewFromURL(url)
	default:
		f, err := os.Create(out)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return stats.NewJSONBackend(f), nil
	***REMOVED***
***REMOVED***

func parseStages(vus []string, total time.Duration) (stages []lib.TestStage, err error) ***REMOVED***
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

func guessType(arg string) string ***REMOVED***
	switch ***REMOVED***
	case strings.Contains(arg, "://"):
		return typeURL
	case strings.HasSuffix(arg, ".js"):
		return typeJS
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
	default:
		return nil, errors.New("Type ambiguous, please specify -t/--type")
	***REMOVED***
***REMOVED***

func action(cc *cli.Context) error ***REMOVED***
	if cc.IsSet("verbose") ***REMOVED***
		log.SetLevel(log.DebugLevel)
	***REMOVED***

	for _, out := range cc.StringSlice("metrics") ***REMOVED***
		backend, err := parseBackend(out)
		if err != nil ***REMOVED***
			return cli.NewExitError(err.Error(), 1)
		***REMOVED***
		stats.DefaultRegistry.Backends = append(stats.DefaultRegistry.Backends, backend)
	***REMOVED***

	var accumulator *accumulate.Backend
	if !cc.Bool("quiet") ***REMOVED***
		accumulator = accumulate.New()
		for _, stat := range cc.StringSlice("exclude") ***REMOVED***
			accumulator.Exclude[stat] = true
		***REMOVED***
		stats.DefaultRegistry.Backends = append(stats.DefaultRegistry.Backends, accumulator)
	***REMOVED***

	stages, err := parseStages(cc.StringSlice("vus"), cc.Duration("duration"))
	if err != nil ***REMOVED***
		return cli.NewExitError(err.Error(), 1)
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

	quit := make(chan os.Signal)
	signal.Notify(quit)

	vus.Start(ctx)
	scaleTo := pollVURamping(ctx, t)
	mVUs := stats.Stat***REMOVED***Name: "vus", Type: stats.GaugeType***REMOVED***
mainLoop:
	for ***REMOVED***
		select ***REMOVED***
		case num := <-scaleTo:
			vus.Scale(num)
			stats.Add(stats.Point***REMOVED***
				Stat:   &mVUs,
				Values: stats.Value(float64(num)),
			***REMOVED***)
		case <-quit:
			cancel()
		case <-ctx.Done():
			break mainLoop
		***REMOVED***
	***REMOVED***

	vus.Stop()

	stats.Add(stats.Point***REMOVED***Stat: &mVUs, Values: stats.Value(0)***REMOVED***)
	stats.Submit()

	if accumulator != nil ***REMOVED***
		for stat, dimensions := range accumulator.Data ***REMOVED***
			switch stat.Type ***REMOVED***
			case stats.CounterType:
				for dname, dim := range dimensions ***REMOVED***
					e := log.WithField("count", stats.ApplyIntent(dim.Sum(), stat.Intent))
					if len(dimensions) == 1 ***REMOVED***
						e.Infof("Metric: %s", stat.Name)
					***REMOVED*** else ***REMOVED***
						e.Infof("Metric: %s.%s", stat.Name, *dname)
					***REMOVED***
				***REMOVED***
			case stats.GaugeType:
				for dname, dim := range dimensions ***REMOVED***
					last := dim.Last
					if last == 0 ***REMOVED***
						continue
					***REMOVED***

					e := log.WithField("val", stats.ApplyIntent(last, stat.Intent))
					if len(dimensions) == 1 ***REMOVED***
						e.Infof("Metric: %s", stat.Name)
					***REMOVED*** else ***REMOVED***
						e.Infof("Metric: %s.%s", stat.Name, *dname)
					***REMOVED***
				***REMOVED***
			case stats.HistogramType:
				first := true
				for dname, dim := range dimensions ***REMOVED***
					if first ***REMOVED***
						log.WithField("count", len(dim.Values)).Infof("Metric: %s", stat.Name)
						first = false
					***REMOVED***
					log.WithFields(log.Fields***REMOVED***
						"min": stats.ApplyIntent(dim.Min(), stat.Intent),
						"max": stats.ApplyIntent(dim.Max(), stat.Intent),
						"avg": stats.ApplyIntent(dim.Avg(), stat.Intent),
						"med": stats.ApplyIntent(dim.Med(), stat.Intent),
					***REMOVED***).Infof("  - %s", *dname)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func main() ***REMOVED***
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
		cli.StringFlag***REMOVED***
			Name:  "type, t",
			Usage: "Input file type, if not evident (url or js)",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "vus, u",
			Usage: "Number of VUs to simulate",
			Value: &cli.StringSlice***REMOVED***"10"***REMOVED***,
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
		cli.StringSliceFlag***REMOVED***
			Name:  "metrics, m",
			Usage: "Write metrics to a file or database",
		***REMOVED***,
		cli.StringSliceFlag***REMOVED***
			Name:  "exclude, e",
			Usage: "Exclude named metrics",
		***REMOVED***,
	***REMOVED***
	app.Action = action
	app.Run(os.Args)
***REMOVED***
