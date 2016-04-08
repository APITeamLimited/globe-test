package ping

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/comm"
	"github.com/loadimpact/speedboat/util"
	"time"
)

func init() ***REMOVED***
	client.RegisterCommand(cli.Command***REMOVED***
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
			util.MasterHostFlag,
			util.MasterPortFlag,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// Pings a master or specified workers.
func actionPing(c *cli.Context) ***REMOVED***
	ct, local := util.MustGetClient(c)
	if local && !c.Bool("local") ***REMOVED***
		log.Fatal("You're about to ping an in-process system, which doesn't make a lot of sense. You probably want to specify --master=..., or use --local if this is actually what you want.")
	***REMOVED***

	in, out := ct.Connector.Run()

	topic := comm.MasterTopic
	if c.Bool("worker") ***REMOVED***
		topic = comm.WorkerTopic
	***REMOVED***
	out <- comm.To(topic, "ping.ping").With(PingMessage***REMOVED***
		Time: time.Now(),
	***REMOVED***)

readLoop:
	for msg := range in ***REMOVED***
		switch msg.Type ***REMOVED***
		case "ping.pong":
			data := PingMessage***REMOVED******REMOVED***
			if err := msg.Take(&data); err != nil ***REMOVED***
				log.WithError(err).Error("Couldn't decode pong")
				break
			***REMOVED***
			log.WithField("time", data.Time).Info("Pong!")
			break readLoop
		***REMOVED***
	***REMOVED***
***REMOVED***
