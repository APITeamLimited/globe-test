package run

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/actions/registry"
	"github.com/loadimpact/speedboat/common"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/message"
	"github.com/loadimpact/speedboat/runner"
	"io/ioutil"
	"path"
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
			cli.StringFlag***REMOVED***
				Name:  "duration, d",
				Usage: "Duration of the test",
				Value: "10s",
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
	in, out, errors := client.Run()

	filename := c.Args()[0]
	conf := loadtest.NewConfig()
	if len(c.Args()) > 0 ***REMOVED***
		data, err := ioutil.ReadFile(filename)
		if err != nil ***REMOVED***
			log.WithError(err).Fatal("Couldn't read test file")
		***REMOVED***

		loadtest.ParseConfig(data, &conf)
	***REMOVED***

	if c.IsSet("script") ***REMOVED***
		conf.Script = c.String("script")
	***REMOVED***
	if c.IsSet("duration") ***REMOVED***
		conf.Duration = c.String("duration")
	***REMOVED***
	if c.IsSet("vus") ***REMOVED***
		conf.VUs = c.Int("vus")
	***REMOVED***

	log.WithField("conf", conf).Info("Config")
	test, err := conf.Compile()
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Invalid test")
	***REMOVED***
	log.WithField("test", test).Info("Test")

	if err = test.Load(path.Dir(filename)); err != nil ***REMOVED***
		log.WithError(err).Fatal("Couldn't load script")
	***REMOVED***

	in, out, errors = test.Run(in, out, errors)
	sequencer := runner.NewSequencer()
runLoop:
	for ***REMOVED***
		select ***REMOVED***
		case msg, ok := <-in:
			// ok is false if in is closed
			if !ok ***REMOVED***
				break runLoop
			***REMOVED***

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
			log.WithError(err).Error("Error")
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
