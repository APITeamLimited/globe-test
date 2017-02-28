/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"gopkg.in/urfave/cli.v1"
	"os"
)

var isTTY = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

func main() ***REMOVED***
	// This won't be needed in cli v2
	cli.VersionFlag.Name = "version"
	cli.HelpFlag.Name = "help"
	cli.HelpFlag.Hidden = true

	app := cli.NewApp()
	app.Name = "k6"
	app.Usage = "a next generation load generator"
	app.Version = "0.12.0"
	app.Commands = []cli.Command***REMOVED***
		commandRun,
		commandInspect,
		commandStatus,
		commandStats,
		commandScale,
		commandPause,
		commandResume,
	***REMOVED***
	app.Flags = []cli.Flag***REMOVED***
		cli.BoolFlag***REMOVED***
			Name:  "verbose, v",
			Usage: "show debug messages",
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:   "address, a",
			Usage:  "address for the API",
			Value:  "127.0.0.1:6565",
			EnvVar: "K6_ADDRESS",
		***REMOVED***,
		cli.BoolFlag***REMOVED***
			Name:   "no-color, n",
			Usage:  "disable colored output",
			EnvVar: "K6_NO_COLOR",
		***REMOVED***,
	***REMOVED***
	app.Before = func(cc *cli.Context) error ***REMOVED***
		if cc.Bool("verbose") ***REMOVED***
			log.SetLevel(log.DebugLevel)
		***REMOVED***
		if cc.Bool("no-color") ***REMOVED***
			color.NoColor = true
		***REMOVED***

		return nil
	***REMOVED***
	if err := app.Run(os.Args); err != nil ***REMOVED***
		os.Exit(1)
	***REMOVED***
***REMOVED***
