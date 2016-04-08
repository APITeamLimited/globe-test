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

// Parses commandline arguments.
//
// topic - The topic (master or worker) to ping
func Parse(c *cli.Context) (topic string, err error) ***REMOVED***
	topic = comm.MasterTopic
	if c.Bool("worker") ***REMOVED***
		topic = comm.WorkerTopic
	***REMOVED***

	return topic, nil
***REMOVED***

// Runs the command.
func Run(in <-chan comm.Message, topic string) <-chan comm.Message ***REMOVED***
	out := make(chan comm.Message)

	go func() ***REMOVED***
		defer close(out)

		// Send a ping
		out <- comm.To(topic, "ping.ping").With(PingMessage***REMOVED***
			Time: time.Now(),
		***REMOVED***)

		// Wait for a reply
		for msg := range in ***REMOVED***
			switch msg.Type ***REMOVED***
			case "ping.pong":
				data := PingMessage***REMOVED******REMOVED***
				if err := msg.Take(&data); err != nil ***REMOVED***
					log.WithError(err).Error("Couldn't decode pong")
					break
				***REMOVED***
				log.WithField("time", data.Time).Info("Pong!")
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return out
***REMOVED***

// Pings a master or specified workers.
func actionPing(c *cli.Context) ***REMOVED***
	topic, err := Parse(c)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Couldn't parse arguments")
	***REMOVED***

	ct, local := util.MustGetClient(c)
	if local && !c.Bool("local") ***REMOVED***
		log.Fatal("You're about to ping an in-process system, which doesn't make a lot of sense. You probably want to specify --master=..., or use --local if this is actually what you want.")
	***REMOVED***

	in, out := ct.Connector.Run()
	for res := range Run(in, topic) ***REMOVED***
		out <- res
	***REMOVED***
***REMOVED***
