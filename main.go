package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/runner/simple"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"path"
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

func action(c *cli.Context) ***REMOVED***
	test, err := makeTest(c)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Configuration error")
	***REMOVED***

	r := simple.New()
	r.URL = test.URL

	timeout := time.Duration(0)
	for _, stage := range test.Stages ***REMOVED***
		timeout += stage.Duration
	***REMOVED***

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	scale := make(chan int, 1)
	scale <- test.Stages[0].VUs.Start
	for t := range runner.Run(ctx, r, scale) ***REMOVED***
		log.WithField("t", t).Info("Test Metric")
	***REMOVED***
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
			Usage: "Script to run (do not use with --url)",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "url",
			Usage: "URL to test (do not use with --script)",
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
