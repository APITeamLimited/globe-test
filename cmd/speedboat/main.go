package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/simple"
	"github.com/rcrowley/go-metrics"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	stdlog "log"
	"os"
	"time"
)

// Configure the global logger.
func configureLogging(c *cli.Context) ***REMOVED***
	log.SetLevel(log.InfoLevel)
	if c.GlobalBool("verbose") ***REMOVED***
		log.SetLevel(log.DebugLevel)
	***REMOVED***
***REMOVED***

func parse(cc *cli.Context) (conf Config, err error) ***REMOVED***
	switch len(cc.Args()) ***REMOVED***
	case 0:
		if !cc.IsSet("script") && !cc.IsSet("url") ***REMOVED***
			return conf, errors.New("No config file, script or URL")
		***REMOVED***
	case 1:
		bytes, err := ioutil.ReadFile(cc.Args()[0])
		if err != nil ***REMOVED***
			return conf, errors.New("Couldn't read config file")
		***REMOVED***
		if err := yaml.Unmarshal(bytes, &conf); err != nil ***REMOVED***
			return conf, errors.New("Couldn't parse config file")
		***REMOVED***
	default:
		return conf, errors.New("Too many arguments!")
	***REMOVED***

	// Let commandline flags override config files
	if cc.IsSet("script") ***REMOVED***
		conf.Script = cc.String("script")
	***REMOVED***
	if cc.IsSet("url") ***REMOVED***
		conf.URL = cc.String("url")
	***REMOVED***
	if cc.IsSet("vus") ***REMOVED***
		conf.VUs = cc.Int("vus")
	***REMOVED***
	if cc.IsSet("duration") ***REMOVED***
		conf.Duration = cc.Duration("duration").String()
	***REMOVED***

	return conf, nil
***REMOVED***

func dumpTest(t *speedboat.Test) ***REMOVED***
	log.WithFields(log.Fields***REMOVED***
		"script": t.Script,
		"url":    t.URL,
	***REMOVED***).Info("General")
	for i, stage := range t.Stages ***REMOVED***
		log.WithFields(log.Fields***REMOVED***
			"#":        i,
			"duration": stage.Duration,
			"start":    stage.StartVUs,
			"end":      stage.EndVUs,
		***REMOVED***).Info("Stage")
	***REMOVED***
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
	conf, err := parse(cc)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Invalid arguments; see --help")
	***REMOVED***

	t, err := conf.MakeTest()
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Configuration error")
	***REMOVED***

	if cc.Bool("dump") ***REMOVED***
		dumpTest(&t)
		return nil
	***REMOVED***

	// Inspect the test to find a suitable runner; additional ones can easily be added
	var runner speedboat.Runner
	switch ***REMOVED***
	case t.URL != "":
		runner = simple.New()
	default:
		log.Fatal("No suitable runner found!")
	***REMOVED***

	// Global metrics
	mVUs := metrics.NewRegisteredGauge("vus", speedboat.Registry)

	// Output metrics appropriately
	go metrics.Log(speedboat.Registry, time.Second, stdlog.New(os.Stderr, "metrics: ", stdlog.Lmicroseconds))

	// Use a "headless controller" to scale VUs by polling the test ramp
	ctx, _ := context.WithTimeout(context.Background(), t.TotalDuration())
	vus := []context.CancelFunc***REMOVED******REMOVED***
	for scale := range headlessController(ctx, &t) ***REMOVED***
		for i := len(vus); i < scale; i++ ***REMOVED***
			log.WithField("id", i).Debug("Spawning VU")
			vuCtx, cancel := context.WithCancel(ctx)
			vus = append(vus, cancel)
			go runner.RunVU(vuCtx, t, len(vus))
		***REMOVED***
		for i := len(vus); i > scale; i-- ***REMOVED***
			log.WithField("id", i-1).Debug("Dropping VU")
			vus[i-1]()
			vus = vus[:i-1]
		***REMOVED***
		mVUs.Update(int64(len(vus)))
	***REMOVED***

	// Wait until the end of the test
	<-ctx.Done()

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
		cli.BoolFlag***REMOVED***
			Name:  "dump",
			Usage: "Dump parsed test and exit",
		***REMOVED***,
	***REMOVED***
	app.Before = func(c *cli.Context) error ***REMOVED***
		configureLogging(c)
		return nil
	***REMOVED***
	app.Action = action
	app.Run(os.Args)
***REMOVED***
