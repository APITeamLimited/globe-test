package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/util"
	"io/ioutil"
	"os"
	"path"
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
		conf.Duration = c.String("duration")
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
	***REMOVED***
	app.Before = func(c *cli.Context) error ***REMOVED***
		configureLogging(c)
		return nil
	***REMOVED***
	app.Action = action
	app.Run(os.Args)
***REMOVED***

// func main() ***REMOVED***
// 	// Free up -v and -h for our own flags
// 	cli.VersionFlag.Name = "version"
// 	cli.HelpFlag.Name = "help, ?"

// 	// Bootstrap using action-registered commandline flags
// 	app := cli.NewApp()
// 	app.Name = "speedboat"
// 	app.Usage = "A next-generation load generator"
// 	app.Version = "0.0.1a1"
// 	app.Flags = []cli.Flag***REMOVED***
// 		cli.BoolFlag***REMOVED***
// 			Name:  "verbose, v",
// 			Usage: "More verbose output",
// 		***REMOVED***,
// 	***REMOVED***
// 	app.Commands = client.GlobalCommands
// 	app.Before = func(c *cli.Context) error ***REMOVED***
// 		configureLogging(c)
// 		return nil
// 	***REMOVED***
// 	app.Run(os.Args)
// ***REMOVED***
