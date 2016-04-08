package run

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/actions/registry"
	"github.com/loadimpact/speedboat/common"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/runner"
	"io/ioutil"
	"path"
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

func actionRun(c *cli.Context) ***REMOVED***
	client, _ := common.MustGetClient(c)
	in, out := client.Run()

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

	in, out = test.Run(in, out)
	sequencer := runner.NewSequencer()

runLoop:
	for msg := range in ***REMOVED***
		switch msg.Type ***REMOVED***
		case "test.log":
			entry := runner.LogEntry***REMOVED******REMOVED***
			if err := msg.Take(&entry); err != nil ***REMOVED***
				log.WithError(err).Error("Couldn't decode log entry")
				break
			***REMOVED***
			log.WithFields(log.Fields***REMOVED***
				"text": entry.Text,
			***REMOVED***).Info("Test Log")
		case "test.metric":
			metric := runner.Metric***REMOVED******REMOVED***
			if err := msg.Take(&metric); err != nil ***REMOVED***
				log.WithError(err).Error("Couldn't decode metric")
				break
			***REMOVED***

			log.WithFields(log.Fields***REMOVED***
				"start":    metric.Start,
				"duration": metric.Duration,
			***REMOVED***).Debug("Test Metric")

			sequencer.Add(metric)
		case "test.end":
			break runLoop
		case "error":
			log.WithError(msg.TakeError()).Error("Test Error")
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
