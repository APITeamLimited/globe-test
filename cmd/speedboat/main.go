package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/simple"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	// Schedule all configured VUs. Because we know the VU curves ahead of time, we:
	// - Make a context with the test's duration as timeout
	// - Loop through all the stages of the test
	// - Spawn VUs that:
	//     - Sleep until they're scheduled to start
	//     - Expire at the projected end of their lifecycles
	//
	// TODO: Account for VU ramping in lifecycle projections!
	ctx, _ := context.WithTimeout(context.Background(), t.TotalDuration())
	offset := time.Duration(0)
	for _, stage := range t.Stages ***REMOVED***
		startOffset := offset
		ctx, _ := context.WithTimeout(ctx, startOffset+stage.Duration)
		go func() ***REMOVED***
			select ***REMOVED***
			case <-time.After(startOffset):
				runner.RunVU(ctx, t)
			case <-ctx.Done():
			***REMOVED***
		***REMOVED***()
		offset += stage.Duration
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
