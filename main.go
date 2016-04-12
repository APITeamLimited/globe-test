package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/util"
	"io/ioutil"
	"os"
	"path"
	"runtime/pprof"
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

	if err = test.Load(base); err != nil ***REMOVED***
		return test, err
	***REMOVED***

	return test, nil
***REMOVED***

func action(c *cli.Context) ***REMOVED***
	test, err := makeTest(c)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Configuration error")
	***REMOVED***

	r, err := util.GetRunner(test.Script)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Couldn't get a runner")
	***REMOVED***
	log.WithField("r", r).Info("Runner")

	err = r.Load(test.Script, test.Source)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Couldn't load script")
	***REMOVED***

	controlChannel := make(chan int, 1)
	currentVUs := test.Stages[0].VUs.Start
	controlChannel <- currentVUs

	startTime := time.Now()
	intervene := time.Tick(time.Duration(1) * time.Second)
	sequencer := runner.NewSequencer()
	results := runner.Run(r, controlChannel)
runLoop:
	for ***REMOVED***
		select ***REMOVED***
		case res := <-results:
			switch res := res.(type) ***REMOVED***
			case runner.LogEntry:
				log.WithField("text", res.Text).Info("Test Log")
			case runner.Metric:
				log.WithField("d", res.Duration).Debug("Test Metric")
				sequencer.Add(res)
			case error:
				log.WithError(res).Error("Test Error")
			***REMOVED***
		case <-intervene:
			vus, stop := test.VUsAt(time.Since(startTime))
			if stop ***REMOVED***
				break runLoop
			***REMOVED***
			if vus != currentVUs ***REMOVED***
				delta := vus - currentVUs
				controlChannel <- delta
				currentVUs = vus
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

// Set up a CPU profile, if requested.
func startCPUProfile(c *cli.Context) ***REMOVED***
	cpuProfile := c.String("cpuprofile")
	if cpuProfile != "" ***REMOVED***
		f, err := os.Create(cpuProfile)
		if err != nil ***REMOVED***
			log.WithError(err).Fatal("Couldn't create CPU profile file")
		***REMOVED***

		pprof.StartCPUProfile(f)
	***REMOVED***
***REMOVED***

// End an ongoing CPU profile.
func endCPUProfile(c *cli.Context) ***REMOVED***
	pprof.StopCPUProfile()
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
		cli.StringFlag***REMOVED***
			Name:  "cpuprofile",
			Usage: "Write a CPU profile to this file",
		***REMOVED***,
	***REMOVED***
	app.Before = func(c *cli.Context) error ***REMOVED***
		configureLogging(c)
		startCPUProfile(c)
		return nil
	***REMOVED***
	app.After = func(c *cli.Context) error ***REMOVED***
		endCPUProfile(c)
		return nil
	***REMOVED***
	app.Action = action
	app.Run(os.Args)
***REMOVED***
