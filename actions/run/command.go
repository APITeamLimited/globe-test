package run

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/actions/registry"
	"github.com/loadimpact/speedboat/common"
	"github.com/loadimpact/speedboat/message"
	"github.com/loadimpact/speedboat/runner"
	"io/ioutil"
	"time"
)

func init() ***REMOVED***
	registry.RegisterCommand(cli.Command***REMOVED***
		Name:   "run",
		Usage:  "Runs a load test",
		Action: actionRun,
		Flags: []cli.Flag***REMOVED***
			common.MasterHostFlag,
			common.MasterPortFlag,
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

func parseMetric(msg message.Message) (m runner.Metric, err error) ***REMOVED***
	duration, ok := msg.Fields["duration"].(float64)
	if !ok ***REMOVED***
		return m, errors.New("Duration is not a float64")
	***REMOVED***

	m.Duration = time.Duration(int64(duration))
	return m, nil
***REMOVED***

func actionRun(c *cli.Context) ***REMOVED***
	client, _ := common.MustGetClient(c)
	in, out, errors := client.Connector.Run()

	if !c.IsSet("script") ***REMOVED***
		log.Fatal("No script file specified!")
	***REMOVED***

	duration := c.Duration("duration")
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
	***REMOVED***)

	timeout := time.After(duration)
	sequencer := runner.NewSequencer()
runLoop:
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-in:
			switch msg.Type ***REMOVED***
			case "run.log":
				log.WithFields(log.Fields***REMOVED***
					"text": msg.Fields["text"],
				***REMOVED***).Info("Test Log")
			case "run.metric":
				m, err := parseMetric(msg)
				if err != nil ***REMOVED***
					log.WithError(err).Error("Couldn't parse metric")
					break
				***REMOVED***

				log.WithFields(log.Fields***REMOVED***
					"start":    m.Start,
					"duration": m.Duration,
				***REMOVED***).Debug("Test Metric")

				sequencer.Add(m)
			case "run.error":
				log.WithFields(log.Fields***REMOVED***
					"error": msg.Fields["error"],
				***REMOVED***).Error("Script Error")
			***REMOVED***
		case err := <-errors:
			log.WithError(err).Error("Ping failed")
		case <-timeout:
			out <- message.NewToWorker("run.stop", message.Fields***REMOVED******REMOVED***)
			log.Info("Test Ended")
			break runLoop
		***REMOVED***
	***REMOVED***

	stats := sequencer.Stats()
	log.WithField("count", sequencer.Count()).Info("Results")
	log.WithFields(log.Fields***REMOVED***
		"min": stats.Duration.Min,
		"max": stats.Duration.Max,
		"avg": stats.Duration.Avg,
		"med": stats.Duration.Med,
	***REMOVED***).Info("Duration")
***REMOVED***
