// +build windows

package cobra

import (
	"os"
	"time"

	"github.com/inconshreveable/mousetrap"
)

var preExecHookFn = preExecHook

func preExecHook(c *Command) ***REMOVED***
	if MousetrapHelpText != "" && mousetrap.StartedByExplorer() ***REMOVED***
		c.Print(MousetrapHelpText)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	***REMOVED***
***REMOVED***
