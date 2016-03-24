package actions

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/actions/registry"
	"github.com/loadimpact/speedboat/common"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/message"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/runner/js"
	"github.com/loadimpact/speedboat/worker"
	"io/ioutil"
)

func init() ***REMOVED***
	registry.RegisterCommand(cli.Command***REMOVED***
		Name:   "run",
		Usage:  "Runs a load test",
		Action: actionRun,
		Flags: []cli.Flag***REMOVED***
			cli.StringFlag***REMOVED***
				Name:  "script",
				Usage: "Script file to run",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	registry.RegisterProcessor(func(*worker.Worker) master.Processor ***REMOVED***
		return &RunProcessor***REMOVED******REMOVED***
	***REMOVED***)
***REMOVED***

type RunProcessor struct***REMOVED******REMOVED***

func (p *RunProcessor) Process(msg message.Message) <-chan message.Message ***REMOVED***
	ch := make(chan message.Message)

	go func() ***REMOVED***
		defer func() ***REMOVED***
			ch <- message.NewToClient("run.end", message.Fields***REMOVED******REMOVED***)
			close(ch)
		***REMOVED***()

		switch msg.Type ***REMOVED***
		case "run.run":
			filename := msg.Fields["filename"].(string)
			src := msg.Fields["src"].(string)

			log.WithFields(log.Fields***REMOVED***
				"filename": filename,
			***REMOVED***).Debug("Running script")

			var r runner.Runner = nil

			r, err := js.New()
			if err != nil ***REMOVED***
				ch <- message.NewToClient("run.error", message.Fields***REMOVED***"error": err***REMOVED***)
				break
			***REMOVED***

			err = r.Load(filename, src)
			if err != nil ***REMOVED***
				ch <- message.NewToClient("run.error", message.Fields***REMOVED***"error": err***REMOVED***)
				break
			***REMOVED***

			for res := range runner.Run(r, 1) ***REMOVED***
				switch res := res.(type) ***REMOVED***
				case runner.LogEntry:
					ch <- message.NewToClient("run.log", message.Fields***REMOVED***
						"text": res.Text,
					***REMOVED***)
				case runner.Metric:
					ch <- message.NewToClient("run.metric", message.Fields***REMOVED***
						"start":    res.Start,
						"duration": res.Duration,
					***REMOVED***)
				case error:
					ch <- message.NewToClient("run.error", message.Fields***REMOVED***
						"error": res.Error(),
					***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
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
