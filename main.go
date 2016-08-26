package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"os"
	"time"
)

func main() ***REMOVED***
	// This won't be needed in cli v2
	cli.VersionFlag.Name = "version"
	cli.HelpFlag.Name = "help"
	cli.HelpFlag.Hidden = true

	app := cli.NewApp()
	app.Name = "speedboat"
	app.Usage = "a next generation load generator"
	app.Version = "0.2.0"
	app.Commands = []cli.Command***REMOVED***
		cli.Command***REMOVED***
			Name:      "run",
			Usage:     "Starts running a load test",
			ArgsUsage: "url|filename",
			Flags: []cli.Flag***REMOVED***
				cli.IntFlag***REMOVED***
					Name:  "vus, u",
					Usage: "virtual users to simulate",
					Value: 10,
				***REMOVED***,
				cli.DurationFlag***REMOVED***
					Name:  "duration, d",
					Usage: "test duration, 0 to run until cancelled",
					Value: 10 * time.Second,
				***REMOVED***,
				cli.StringFlag***REMOVED***
					Name:  "type, t",
					Usage: "input type, one of: auto, url, js",
					Value: "auto",
				***REMOVED***,
			***REMOVED***,
			Action: actionRun,
		***REMOVED***,
		cli.Command***REMOVED***
			Name:      "status",
			Usage:     "Looks up the status of a running test",
			ArgsUsage: " ",
			Action:    actionStatus,
		***REMOVED***,
		cli.Command***REMOVED***
			Name:      "scale",
			Usage:     "Scales a running test",
			ArgsUsage: "vus",
			Action:    actionScale,
		***REMOVED***,
		cli.Command***REMOVED***
			Name:      "abort",
			Usage:     "Aborts a running test",
			ArgsUsage: " ",
			Action:    actionAbort,
		***REMOVED***,
	***REMOVED***
	app.Flags = []cli.Flag***REMOVED***
		cli.BoolFlag***REMOVED***
			Name:  "verbose, v",
			Usage: "show debug messages",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "address, a",
			Usage: "address for the API",
			Value: "127.0.0.1:6565",
		***REMOVED***,
	***REMOVED***
	app.Before = func(cc *cli.Context) error ***REMOVED***
		if cc.Bool("verbose") ***REMOVED***
			log.SetLevel(log.DebugLevel)
		***REMOVED***

		return nil
	***REMOVED***
	if err := app.Run(os.Args); err != nil ***REMOVED***
		os.Exit(1)
	***REMOVED***
***REMOVED***
