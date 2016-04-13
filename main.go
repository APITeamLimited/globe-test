package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/runner/js"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func getRunner(filename string) (runner.Runner, error) ***REMOVED***
	switch path.Ext(filename) ***REMOVED***
	case ".js":
		return js.New()
	default:
		return nil, errors.New("No runner found")
	***REMOVED***
***REMOVED***

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

	srcb, err := ioutil.ReadFile(path.Join(base, test.Script))
	if err != nil ***REMOVED***
		return test, err
	***REMOVED***
	test.Source = string(srcb)

	return test, nil
***REMOVED***

func action(c *cli.Context) ***REMOVED***
	test, err := makeTest(c)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Configuration error")
	***REMOVED***

	r, err := getRunner(test.Script)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Couldn't get a runner")
	***REMOVED***

	err = r.Load(test.Script, test.Source)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Couldn't load script")
	***REMOVED***

	// Write a number to the control channel to make the test scale to that many
	// VUs; close it to make the test terminate.
	controlChannel := make(chan int, 1)
	controlChannel <- test.Stages[0].VUs.Start

	sequencer := runner.NewSequencer()
	startTime := time.Now()

	intervene := time.NewTicker(time.Duration(1) * time.Second)
	results := runner.Run(r, controlChannel)
runLoop:
	for ***REMOVED***
		select ***REMOVED***
		case res, ok := <-results:
			// The results channel will be closed once all VUs are done.
			if !ok ***REMOVED***
				break runLoop
			***REMOVED***
			switch res := res.(type) ***REMOVED***
			case runner.LogEntry:
				log.WithField("text", res.Text).Info("Test Log")
			case runner.Metric:
				log.WithField("d", res.Duration).Debug("Test Metric")
				sequencer.Add(res)
			case error:
				log.WithError(res).Error("Test Error")
			***REMOVED***
		case <-intervene.C:
			vus, stop := test.VUsAt(time.Since(startTime))
			if stop ***REMOVED***
				// Stop the timer, and let VUs gracefully terminate.
				intervene.Stop()
				close(controlChannel)
			***REMOVED*** else ***REMOVED***
				controlChannel <- vus
			***REMOVED***
		***REMOVED***
	***REMOVED***

	stats := sequencer.Stats()
	log.WithField("count", sequencer.Count()).Info("Results")
	log.WithFields(log.Fields***REMOVED***
		"min": stats.Duration.Min,
		"max": stats.Duration.Max,
		"avg": stats.Duration.Avg,
		"med": stats.Duration.Med,
	***REMOVED***).Info("Duration")
***REMOVED***

// Configure the global logger.
func configureLogging(c *cli.Context) ***REMOVED***
	if c.GlobalBool("verbose") ***REMOVED***
		log.SetLevel(log.DebugLevel)
	***REMOVED***
***REMOVED***

func main() ***REMOVED***
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
	***REMOVED***
	app.Before = func(c *cli.Context) error ***REMOVED***
		configureLogging(c)
		return nil
	***REMOVED***
	app.Action = action
	app.Run(os.Args)
***REMOVED***
