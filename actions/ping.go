package actions

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/actions/registry"
	"github.com/loadimpact/speedboat/common"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/message"
	"github.com/loadimpact/speedboat/worker"
	"time"
)

func init() ***REMOVED***
	registry.RegisterCommand(cli.Command***REMOVED***
		Name:   "ping",
		Usage:  "Tests master connectivity",
		Action: actionPing,
		Flags: []cli.Flag***REMOVED***
			cli.BoolFlag***REMOVED***
				Name:  "worker",
				Usage: "Pings a worker instead of the master",
			***REMOVED***,
			cli.BoolFlag***REMOVED***
				Name:  "local",
				Usage: "Allow pinging an inproc master/worker",
			***REMOVED***,
			common.MasterHostFlag,
			common.MasterPortFlag,
		***REMOVED***,
	***REMOVED***)
	registry.RegisterMasterProcessor(func(*master.Master) master.Processor ***REMOVED***
		return &PingMasterProcessor***REMOVED******REMOVED***
	***REMOVED***)
	registry.RegisterProcessor(func(*worker.Worker) master.Processor ***REMOVED***
		return &PingProcessor***REMOVED******REMOVED***
	***REMOVED***)
***REMOVED***

// Processes worker pings.
type PingProcessor struct***REMOVED******REMOVED***

func (*PingProcessor) Process(msg message.Message) <-chan message.Message ***REMOVED***
	out := make(chan message.Message)

	go func() ***REMOVED***
		defer close(out)
		switch msg.Type ***REMOVED***
		case "ping.ping":
			out <- message.NewToClient("ping.pong", msg.Fields)
		***REMOVED***
	***REMOVED***()

	return out
***REMOVED***

// Processes master pings.
type PingMasterProcessor struct***REMOVED******REMOVED***

func (*PingMasterProcessor) Process(msg message.Message) <-chan message.Message ***REMOVED***
	out := make(chan message.Message)

	go func() ***REMOVED***
		defer close(out)
		switch msg.Type ***REMOVED***
		case "ping.ping":
			out <- message.NewToClient("ping.pong", msg.Fields)
		***REMOVED***
	***REMOVED***()

	return out
***REMOVED***

// Pings a master or specified workers.
func actionPing(c *cli.Context) ***REMOVED***
	client, local := common.MustGetClient(c)
	if local && !c.Bool("local") ***REMOVED***
		log.Fatal("You're about to ping an in-process system, which doesn't make a lot of sense. You probably want to specify --master=..., or use --local if this is actually what you want.")
	***REMOVED***

	in, out, errors := client.Connector.Run()

	msgTopic := message.MasterTopic
	if c.Bool("worker") ***REMOVED***
		msgTopic = message.WorkerTopic
	***REMOVED***
	out <- message.Message***REMOVED***
		Topic: msgTopic,
		Type:  "ping.ping",
		Fields: message.Fields***REMOVED***
			"time": time.Now().Format("15:04:05 2006-01-02 MST"),
		***REMOVED***,
	***REMOVED***

readLoop:
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-in:
			switch msg.Type ***REMOVED***
			case "ping.pong":
				log.WithFields(log.Fields***REMOVED***
					"time": msg.Fields["time"],
				***REMOVED***).Info("Pong!")
				break readLoop
			***REMOVED***
		case err := <-errors:
			log.WithError(err).Error("Ping failed")
		***REMOVED***
	***REMOVED***
***REMOVED***
