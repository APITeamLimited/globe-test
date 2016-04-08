package run

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/common"
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/message"
	"github.com/loadimpact/speedboat/runner"
	"io/ioutil"
	"path"
	"time"
)

func init() ***REMOVED***
	client.RegisterCommand(cli.Command***REMOVED***
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
	ct, _ := common.MustGetClient(c)
	in, out := ct.Run()

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

	out <- message.ToWorker("test.run").With(MessageTestRun***REMOVED***
		Filename: test.Script,
		Source:   test.Source,
		VUs:      test.Stages[0].VUs.Start,
	***REMOVED***)

	startTime := time.Now()
	intervene := time.Tick(time.Duration(1) * time.Second)
	sequencer := runner.NewSequencer()
	currentVUs := 0
runLoop:
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-in:
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
			case "error":
				log.WithError(msg.TakeError()).Error("Test Error")
			***REMOVED***
		case <-intervene:
			vus, stop := test.VUsAt(time.Since(startTime))
			if stop ***REMOVED***
				out <- message.ToWorker("test.stop")
				break runLoop
			***REMOVED***
			if vus != currentVUs ***REMOVED***
				out <- message.ToWorker("test.scale").With(MessageTestScale***REMOVED***VUs: vus***REMOVED***)
				currentVUs = vus
			***REMOVED***
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
