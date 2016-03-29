package run

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/actions/registry"
	"github.com/loadimpact/speedboat/common"
	"github.com/loadimpact/speedboat/message"
	"io/ioutil"
	"time"
)

func init() ***REMOVED***
	registry.RegisterCommand(cli.Command***REMOVED***
		Name:   "run",
		Usage:  "Runs a load test",
		Action: actionRun,
		Flags: []cli.Flag***REMOVED***
			cli.StringFlag***REMOVED***
				Name:  "script, s",
				Usage: "Script file to run",
			***REMOVED***,
			cli.IntFlag***REMOVED***
				Name:  "vus, u",
				Usage: "Virtual Users to simulate",
				Value: 2,
			***REMOVED***,
			cli.DurationFlag***REMOVED***
				Name:  "duration, d",
				Usage: "Duration of the test",
				Value: time.Duration(10) * time.Second,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func actionRun(c *cli.Context) ***REMOVED***
	client, _ := common.MustGetClient(c)
	in, out, errors := client.Connector.Run()

	if !c.IsSet("script") ***REMOVED***
		log.Fatal("No script file specified!")
	***REMOVED***

	filename := c.String("script")
	srcb, err := ioutil.ReadFile(filename)
	src := string(srcb)
	if err != nil ***REMOVED***
		log.WithError(err).WithFields(log.Fields***REMOVED***
			"filename": filename,
		***REMOVED***).Fatal("Couldn't read script")
	***REMOVED***

	out <- message.NewToWorker("run.run", message.Fields***REMOVED***
		"filename": c.String("script"),
		"src":      src,
		"vus":      c.Int("vus"),
		"duration": int64(c.Duration("duration")) / int64(time.Millisecond),
	***REMOVED***)

readLoop:
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-in:
			switch msg.Type ***REMOVED***
			case "run.log":
				log.WithFields(log.Fields***REMOVED***
					"text": msg.Fields["text"],
				***REMOVED***).Info("Test Log")
			case "run.metric":
				log.WithFields(log.Fields***REMOVED***
					"start":    msg.Fields["start"],
					"duration": msg.Fields["duration"],
				***REMOVED***).Info("Test Metric")
			case "run.error":
				log.WithFields(log.Fields***REMOVED***
					"error": msg.Fields["error"],
				***REMOVED***).Error("Script Error")
			case "run.end":
				log.Info("-- Test End --")
				break readLoop
			***REMOVED***
		case err := <-errors:
			log.WithError(err).Error("Ping failed")
		***REMOVED***
	***REMOVED***
***REMOVED***
