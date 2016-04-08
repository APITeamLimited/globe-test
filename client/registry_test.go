package client

import (
	"github.com/codegangsta/cli"
	"testing"
)

func TestRegisterCommand(t *testing.T) ***REMOVED***
	oldGlobalCommands := GlobalCommands
	GlobalCommands = nil
	defer func() ***REMOVED*** GlobalCommands = oldGlobalCommands ***REMOVED***()

	RegisterCommand(cli.Command***REMOVED***Name: "test"***REMOVED***)
	if len(GlobalCommands) != 1 ***REMOVED***
		t.Error("Command not registered")
	***REMOVED***
***REMOVED***
