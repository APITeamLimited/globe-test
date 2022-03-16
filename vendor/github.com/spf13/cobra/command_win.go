//go:build windows
// +build windows

package cobra

import (
	"fmt"
	"os"
	"time"

	"github.com/inconshreveable/mousetrap"
)

var preExecHookFn = preExecHook

func preExecHook(c *Command) ***REMOVED***
	if MousetrapHelpText != "" && mousetrap.StartedByExplorer() ***REMOVED***
		c.Print(MousetrapHelpText)
		if MousetrapDisplayDuration > 0 ***REMOVED***
			time.Sleep(MousetrapDisplayDuration)
		***REMOVED*** else ***REMOVED***
			c.Println("Press return to continue...")
			fmt.Scanln()
		***REMOVED***
		os.Exit(1)
	***REMOVED***
***REMOVED***
