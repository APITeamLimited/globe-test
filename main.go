package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/aggregate"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/report"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/runner/simple"
	"github.com/loadimpact/speedboat/runner/v8js"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"time"
)

func makeTest(c *cli.Context) (test loadtest.LoadTest, err error) ***REMOVED***
	base := ""
	conf := loadtest.NewConfig()
	if len(c.Args()) > 0 ***REMOVED***
		filename := c.Args()[0]
		base = path.Dir(filename)
		data, err := ioutil.ReadFile(filename)
		if err != nil ***REMOVED***
			return test, err
		***REMOVED***

		loadtest.ParseConfig(data, &conf)
	***REMOVED***

	if c.IsSet("script") ***REMOVED***
		conf.Script = c.String("script")
		base = ""
	***REMOVED***
	if c.IsSet("url") ***REMOVED***
		conf.URL = c.String("url")
	***REMOVED***
	if c.IsSet("duration") ***REMOVED***
		conf.Duration = c.Duration("duration").String()
	***REMOVED***
	if c.IsSet("vus") ***REMOVED***
		conf.VUs = c.Int("vus")
	***REMOVED***

	test, err = conf.Compile()
	if err != nil ***REMOVED***
		return test, err
	***REMOVED***

	if test.Script != "" ***REMOVED***
		srcb, err := ioutil.ReadFile(path.Join(base, test.Script))
		if err != nil ***REMOVED***
			return test, err
		***REMOVED***
		test.Source = string(srcb)
	***REMOVED***

	return test, nil
***REMOVED***

func run(test loadtest.LoadTest, r runner.Runner) (<-chan runner.Result, chan int) ***REMOVED***
	ch := make(chan runner.Result)
	scale := make(chan int, 1)

	go func() ***REMOVED***
		defer close(ch)

		timeout := time.Duration(0)
		for _, stage := range test.Stages ***REMOVED***
			timeout += stage.Duration
		***REMOVED***

		ctx, _ := context.WithTimeout(context.Background(), timeout)
		scale <- test.Stages[0].VUs.Start

		for res := range runner.Run(ctx, r, test, scale) ***REMOVED***
			ch <- res
		***REMOVED***
	***REMOVED***()

	return ch, scale
***REMOVED***

func action(c *cli.Context) error ***REMOVED***
	test, err := makeTest(c)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Configuration error")
	***REMOVED***

	r := runner.Runner(nil)

	// Start the pipeline by just running requests
	if test.Script != "" ***REMOVED***
		ext := path.Ext(test.Script)
		switch ext ***REMOVED***
		case ".js":
			r = v8js.New(test.Script, test.Source)
		default:
			log.WithField("ext", ext).Fatal("No runner found")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r = simple.New()
	***REMOVED***
	pipeline, scale := run(test, r)

	// Ramp VUs according to the test definition
	pipeline = runner.Ramp(&test, scale, pipeline)

	// Stick result aggregation onto it
	stats := aggregate.Stats***REMOVED******REMOVED***
	stats.Time.Values = make([]time.Duration, 30000000)[:0]
	pipeline = aggregate.Aggregate(&stats, pipeline)

	// Log results to a file
	outFilename := c.String("out-file")
	if outFilename != "" ***REMOVED***
		reporter := report.CSVReporter***REMOVED******REMOVED***
		if outFilename != "-" ***REMOVED***
			f, err := os.Create("results.csv")
			if err != nil ***REMOVED***
				log.WithError(err).Fatal("Couldn't open log file")
			***REMOVED***
			pipeline = report.Report(reporter, f, pipeline)
		***REMOVED*** else ***REMOVED***
			pipeline = report.Report(reporter, os.Stdout, pipeline)
		***REMOVED***
	***REMOVED***

	// Listen for SIGINT (Ctrl+C)
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

runLoop:
	for ***REMOVED***
		select ***REMOVED***
		case res, ok := <-pipeline:
			if !ok ***REMOVED***
				break runLoop
			***REMOVED***

			switch ***REMOVED***
			case res.Error != nil:
				l := log.WithError(res.Error)
				if res.Time != time.Duration(0) ***REMOVED***
					l = l.WithField("t", res.Time)
				***REMOVED***
				l.Error("Error")
			case res.Text != "":
				l := log.WithField("text", res.Text)
				if res.Time != time.Duration(0) ***REMOVED***
					l = l.WithField("t", res.Time)
				***REMOVED***
				l.Info("Log")
			default:
				// log.WithField("t", res.Time).Debug("Metric")
			***REMOVED***
		case <-stop:
			break runLoop
		***REMOVED***
	***REMOVED***

	log.WithField("results", stats.Results).Info("Finished")
	log.WithFields(log.Fields***REMOVED***
		"min": stats.Time.Min,
		"max": stats.Time.Max,
		"med": stats.Time.Med,
		"avg": stats.Time.Avg,
	***REMOVED***).Info("Time")

	return nil
***REMOVED***

// Configure the global logger.
func configureLogging(c *cli.Context) ***REMOVED***
	if c.GlobalBool("verbose") ***REMOVED***
		log.SetLevel(log.DebugLevel)
	***REMOVED***
***REMOVED***

func main() ***REMOVED***
	// Up the thread limit (default: 10.000)
	debug.SetMaxThreads(100000)
	// Up the stack size limit (default: 1GB)
	debug.SetMaxStack(3 * 1000000000)

	// Free up -v and -h for our own flags
	cli.VersionFlag.Name = "version"
	cli.HelpFlag.Name = "help, ?"

	// Bootstrap using action-registered commandline flags
	app := cli.NewApp()
	app.Name = "speedboat"
	app.Usage = "A next-generation load generator"
	app.Version = "0.0.1a1"
	app.Flags = []cli.Flag***REMOVED***
		cli.BoolFlag***REMOVED***
			Name:  "verbose, v",
			Usage: "More verbose output",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "script, s",
			Usage: "Script to run",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "url",
			Usage: "URL to test",
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
		cli.StringFlag***REMOVED***
			Name:  "out-file, o",
			Usage: "Output raw metrics to a file",
		***REMOVED***,
	***REMOVED***
	app.Before = func(c *cli.Context) error ***REMOVED***
		configureLogging(c)
		return nil
	***REMOVED***
	app.Action = action
	app.Run(os.Args)
***REMOVED***
