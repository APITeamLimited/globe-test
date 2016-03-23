package registry

import (
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/worker"
)

// All registered cli commands.
var GlobalCommands []cli.Command

// All registered master processors.
var GlobalMasterProcessors []func(*master.Master) master.Processor

// All registered worker processors.
var GlobalProcessors []func(*worker.Worker) master.Processor

// Register an application subcommand.
func RegisterCommand(cmd cli.Command) ***REMOVED***
	GlobalCommands = append(GlobalCommands, cmd)
***REMOVED***

// Register a master handler.
func RegisterMasterProcessor(factory func(*master.Master) master.Processor) ***REMOVED***
	GlobalMasterProcessors = append(GlobalMasterProcessors, factory)
***REMOVED***

// Register a worker processor.
func RegisterProcessor(factory func(*worker.Worker) master.Processor) ***REMOVED***
	GlobalProcessors = append(GlobalProcessors, factory)
***REMOVED***
