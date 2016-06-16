package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/js"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/loadimpact/speedboat/sampler/influxdb"
	"github.com/loadimpact/speedboat/sampler/stream"
	"github.com/loadimpact/speedboat/simple"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	typeURL = "url"
	typeYML = "yml"
	typeJS  = "js"
)

// Configure the global logger.
func configureLogging(c *cli.Context) ***REMOVED***
	log.SetLevel(log.InfoLevel)
	if c.GlobalBool("verbose") ***REMOVED***
		log.SetLevel(log.DebugLevel)
	***REMOVED***
***REMOVED***

// Configure the global sampler.
func configureSampler(c *cli.Context) ***REMOVED***
	sampler.DefaultSampler.OnError = func(err error) ***REMOVED***
		log.WithError(err).Error("[Sampler error]")
	***REMOVED***

	for _, output := range c.GlobalStringSlice("metrics") ***REMOVED***
		parts := strings.SplitN(output, "+", 2)
		switch parts[0] ***REMOVED***
		case "influxdb":
			out, err := influxdb.NewFromURL(parts[1])
			if err != nil ***REMOVED***
				log.WithError(err).Fatal("Couldn't create InfluxDB client")
			***REMOVED***
			sampler.DefaultSampler.Outputs = append(sampler.DefaultSampler.Outputs, out)
		default:
			var writer io.WriteCloser
			switch output ***REMOVED***
			case "stdout", "-":
				writer = os.Stdout
			default:
				file, err := os.Create(output)
				if err != nil ***REMOVED***
					log.WithError(err).Fatal("Couldn't create output file")
				***REMOVED***
				writer = file
			***REMOVED***

			var out sampler.Output
			switch c.GlobalString("format") ***REMOVED***
			case "json":
				out = &stream.JSONOutput***REMOVED***Output: writer***REMOVED***
			case "csv":
				out = &stream.CSVOutput***REMOVED***Output: writer***REMOVED***
			default:
				log.Fatal("Unknown output format")
			***REMOVED***
			sampler.DefaultSampler.Outputs = append(sampler.DefaultSampler.Outputs, out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func guessType(arg string) string ***REMOVED***
	switch ***REMOVED***
	case strings.Contains(arg, "://"):
		return typeURL
	case strings.HasSuffix(arg, ".js"):
		return typeJS
	case strings.HasSuffix(arg, ".yml"):
		return typeYML
	***REMOVED***
	return ""
***REMOVED***

func parse(cc *cli.Context) (conf Config, err error) ***REMOVED***
	if len(cc.Args()) == 0 ***REMOVED***
		return conf, errors.New("Nothing to do!")
	***REMOVED***

	conf.VUs = cc.Int("vus")
	conf.Duration = cc.Duration("duration").String()

	arg := cc.Args()[0]
	argType := cc.String("type")
	if argType == "" ***REMOVED***
		argType = guessType(arg)
	***REMOVED***

	switch argType ***REMOVED***
	case typeYML:
		bytes, err := ioutil.ReadFile(cc.Args()[0])
		if err != nil ***REMOVED***
			return conf, errors.New("Couldn't read config file")
		***REMOVED***
		if err := yaml.Unmarshal(bytes, &conf); err != nil ***REMOVED***
			return conf, errors.New("Couldn't parse config file")
		***REMOVED***
	case typeURL:
		conf.URL = arg
	case typeJS:
		conf.Script = arg
	default:
		return conf, errors.New("Unsure of what to do, try specifying --type")
	***REMOVED***

	return conf, nil
***REMOVED***

func headlessController(c context.Context, t *speedboat.Test) <-chan int ***REMOVED***
	ch := make(chan int)

	go func() ***REMOVED***
		defer close(ch)

		select ***REMOVED***
		case ch <- t.VUsAt(0):
		case <-c.Done():
			return
		***REMOVED***

		startTime := time.Now()
		ticker := time.NewTicker(100 * time.Millisecond)
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				ch <- t.VUsAt(time.Since(startTime))
			case <-c.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***

func action(cc *cli.Context) error ***REMOVED***
	if len(cc.Args()) == 0 ***REMOVED***
		cli.ShowAppHelp(cc)
		return nil
	***REMOVED***

	conf, err := parse(cc)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Invalid arguments; see --help")
	***REMOVED***

	t, err := conf.MakeTest()
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Configuration error")
	***REMOVED***

	// Inspect the test to find a suitable runner; additional ones can easily be added
	var runner speedboat.Runner
	switch ***REMOVED***
	case t.Script == "":
		runner = simple.New(t)
	case strings.HasSuffix(t.Script, ".js"):
		src, err := ioutil.ReadFile(t.Script)
		if err != nil ***REMOVED***
			log.WithError(err).Fatal("Couldn't read script")
		***REMOVED***
		runner = js.New(t, t.Script, string(src))
	default:
		log.Fatal("No suitable runner found!")
	***REMOVED***

	// Context that expires at the end of the test
	ctx, cancel := context.WithTimeout(context.Background(), t.TotalDuration())

	// Configure the VU logger
	logger := &log.Logger***REMOVED***
		Out:       os.Stderr,
		Level:     log.DebugLevel,
		Formatter: &log.TextFormatter***REMOVED******REMOVED***,
	***REMOVED***
	ctx = speedboat.WithLogger(ctx, logger)

	// Store metrics unless the --quiet flag is specified
	quiet := cc.Bool("quiet")
	sampler.DefaultSampler.Accumulate = !quiet

	// Commit metrics to any configured backends once per second
	go func() ***REMOVED***
		ticker := time.NewTicker(1 * time.Second)
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				commitMetrics()
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Use a "headless controller" to scale VUs by polling the test ramp
	mVUs := sampler.Gauge("vus")
	vus := []context.CancelFunc***REMOVED******REMOVED***
	for scale := range headlessController(ctx, &t) ***REMOVED***
		for i := len(vus); i < scale; i++ ***REMOVED***
			log.WithField("id", i).Debug("Spawning VU")
			vuCtx, vuCancel := context.WithCancel(ctx)
			vus = append(vus, vuCancel)

			go func(ctx context.Context) ***REMOVED***
				defer func() ***REMOVED***
					if v := recover(); v != nil ***REMOVED***
						switch err := v.(type) ***REMOVED***
						case speedboat.FlowControl:
							switch err ***REMOVED***
							case speedboat.AbortTest:
								log.Error("Test aborted")
								cancel()
							***REMOVED***
						default:
							panic(err)
						***REMOVED***
					***REMOVED***
				***REMOVED***()

				vu, err := runner.NewVU()
				if err != nil ***REMOVED***
					log.WithError(err).Error("Couldn't spawn VU")
					return
				***REMOVED***

				vu.Reconfigure(int64(i))
				for ***REMOVED***
					select ***REMOVED***
					case <-ctx.Done():
						return
					default:
						vu.RunOnce(ctx)
					***REMOVED***
				***REMOVED***
			***REMOVED***(vuCtx)
		***REMOVED***
		for i := len(vus); i > scale; i-- ***REMOVED***
			log.WithField("id", i-1).Debug("Dropping VU")
			vus[i-1]()
			vus = vus[:i-1]
		***REMOVED***
		mVUs.Int(len(vus))
	***REMOVED***

	// Wait until the end of the test
	<-ctx.Done()

	// Print and commit final metrics
	if !quiet ***REMOVED***
		printMetrics()
	***REMOVED***
	commitMetrics()
	closeMetrics()

	return nil
***REMOVED***

func main() ***REMOVED***
	// Free up -v and -h for our own flags
	cli.VersionFlag.Name = "version"
	cli.HelpFlag.Name = "help, ?"

	// Bootstrap using action-registered commandline flags
	app := cli.NewApp()
	app.Name = "speedboat"
	app.Usage = "A next-generation load generator"
	app.Version = "1.0.0-mvp1"
	app.Flags = []cli.Flag***REMOVED***
		cli.StringFlag***REMOVED***
			Name:  "type, t",
			Usage: "Input file type, if not evident (url, yml or js)",
		***REMOVED***,
		cli.IntFlag***REMOVED***
			Name:  "vus, u",
			Usage: "Number of VUs to simulate",
			Value: 10,
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
		cli.StringFlag***REMOVED***
			Name:  "format, f",
			Usage: "Metric output format (json or csv)",
			Value: "json",
		***REMOVED***,
	***REMOVED***
	app.Before = func(c *cli.Context) error ***REMOVED***
		configureLogging(c)
		configureSampler(c)
		return nil
	***REMOVED***
	app.Action = action
	app.Run(os.Args)
***REMOVED***
