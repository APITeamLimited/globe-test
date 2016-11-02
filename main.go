package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() ***REMOVED***
	// This won't be needed in cli v2
	cli.VersionFlag.Name = "version"
	cli.HelpFlag.Name = "help"
	cli.HelpFlag.Hidden = true

	app := cli.NewApp()
	app.Name = "speedboat"
	app.Usage = "a next generation load generator"
	app.Version = "0.2.1"
	app.Commands = []cli.Command***REMOVED***
		commandRun,
		commandInspect,
		commandStatus,
		commandScale,
		commandStart,
		commandPause,
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
		gin.SetMode(gin.ReleaseMode)

		if cc.Bool("verbose") ***REMOVED***
			log.SetLevel(log.DebugLevel)
		***REMOVED***

		return nil
	***REMOVED***
	if err := app.Run(os.Args); err != nil ***REMOVED***
		os.Exit(1)
	***REMOVED***
***REMOVED***
