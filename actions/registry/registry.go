package registry

import (
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/message"
	"github.com/loadimpact/speedboat/worker"
)

// All registered cli commands.
var GlobalCommands []cli.Command

// All registered master handlers.
var GlobalHandlers []func(*master.Master, message.Message, chan message.Message) bool

// All registered worker processors.
var GlobalProcessors []func(*worker.Worker, message.Message, chan message.Message) bool

// Register an application subcommand.
func RegisterCommand(cmd cli.Command) ***REMOVED***
	GlobalCommands = append(GlobalCommands, cmd)
***REMOVED***

// Register a master handler.
func RegisterHandler(handler func(*master.Master, message.Message, chan message.Message) bool) ***REMOVED***
	GlobalHandlers = append(GlobalHandlers, handler)
***REMOVED***

// Register a worker processor.
func RegisterProcessor(proc func(*worker.Worker, message.Message, chan message.Message) bool) ***REMOVED***
	GlobalProcessors = append(GlobalProcessors, proc)
***REMOVED***
