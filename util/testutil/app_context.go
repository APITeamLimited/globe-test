package testutil

import (
	"github.com/codegangsta/cli"
	"strings"
)

func WithAppContext(args string, cmd cli.Command, fn func(c *cli.Context)) error ***REMOVED***
	app := cli.NewApp()

	// Make a surrogate command, with the same invocation as cmd
	app.Commands = []cli.Command***REMOVED***
		cli.Command***REMOVED***
			Name:   "cmd",
			Flags:  cmd.Flags,
			Action: func(c *cli.Context) ***REMOVED*** fn(c) ***REMOVED***,
		***REMOVED***,
	***REMOVED***

	// Don't print an usage message to stdout for invalid arguments
	app.OnUsageError = func(_ *cli.Context, err error, _ bool) error ***REMOVED*** return err ***REMOVED***

	// Imitate os.Args by prepending a program name, and invoke the surrogate
	appArgs := []string***REMOVED***"program", "cmd"***REMOVED***
	for _, arg := range strings.Split(args, " ") ***REMOVED***
		appArgs = append(appArgs, arg)
	***REMOVED***

	// Returns an error for invalid arguments
	return app.Run(appArgs)
***REMOVED***
