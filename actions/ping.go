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
	registry.RegisterHandler(handlePing)
	registry.RegisterProcessor(processPing)
***REMOVED***

func processPing(w *worker.Worker, msg message.Message, out chan message.Message) bool ***REMOVED***
	switch msg.Type ***REMOVED***
	case "ping.ping":
		out <- message.NewToClient("ping.pong", msg.Body)
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func handlePing(m *master.Master, msg message.Message, out chan message.Message) bool ***REMOVED***
	switch msg.Type ***REMOVED***
	case "ping.ping":
		out <- message.NewToClient("ping.pong", msg.Body)
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

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
		Body:  time.Now().Format("15:04:05 2006-01-02 MST"),
	***REMOVED***

readLoop:
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-in:
			switch msg.Type ***REMOVED***
			case "ping.pong":
				log.WithFields(log.Fields***REMOVED***
					"body": msg.Body,
				***REMOVED***).Info("Pong!")
				break readLoop
			***REMOVED***
		case err := <-errors:
			log.WithError(err).Error("Ping failed")
		***REMOVED***
	***REMOVED***
***REMOVED***
