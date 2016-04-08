package ping

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/comm"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/util"
	"github.com/loadimpact/speedboat/worker"
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
	master.RegisterProcessor(func(*master.Master) comm.Processor ***REMOVED***
		return &PingProcessor***REMOVED******REMOVED***
	***REMOVED***)
	worker.RegisterProcessor(func(*worker.Worker) comm.Processor ***REMOVED***
		return &PingProcessor***REMOVED******REMOVED***
	***REMOVED***)
***REMOVED***

// Processes pings, on both master and worker.
type PingProcessor struct***REMOVED******REMOVED***

type PingMessage struct ***REMOVED***
	Time time.Time
***REMOVED***

func (*PingProcessor) Process(msg comm.Message) <-chan comm.Message ***REMOVED***
	out := make(chan comm.Message)

	go func() ***REMOVED***
		defer close(out)
		switch msg.Type ***REMOVED***
		case "ping.ping":
			data := PingMessage***REMOVED******REMOVED***
			if err := msg.Take(&data); err != nil ***REMOVED***
				out <- comm.ToClient("error").WithError(err)
				break
			***REMOVED***
			out <- comm.ToClient("ping.pong").With(data)
		***REMOVED***
	***REMOVED***()

	return out
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
