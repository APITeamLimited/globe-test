package loadtest

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/actions/registry"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/message"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/runner/js"
	"github.com/loadimpact/speedboat/worker"
)

func init() ***REMOVED***
	registry.RegisterProcessor(func(*worker.Worker) master.Processor ***REMOVED***
		return &LoadTestProcessor***REMOVED******REMOVED***
	***REMOVED***)
***REMOVED***

type LoadTestProcessor struct ***REMOVED***
	// Close this channel to stop the currently running test
	stopChannel chan interface***REMOVED******REMOVED***
***REMOVED***

func (p *LoadTestProcessor) Process(msg message.Message) <-chan message.Message ***REMOVED***
	ch := make(chan message.Message)

	go func() ***REMOVED***
		defer close(ch)

		switch msg.Type ***REMOVED***
		case "test.run":
			p.stopChannel = make(chan interface***REMOVED******REMOVED***)

			// filename := msg.Fields["filename"].(string)
			// src := msg.Fields["src"].(string)
			// vus := int(msg.Fields["vus"].(float64))
			data := MessageTestRun***REMOVED******REMOVED***
			if err := msg.Take(&data); err != nil ***REMOVED***
				log.WithError(err).Error("Couldn't decode test.run")
				return
			***REMOVED***

			log.WithFields(log.Fields***REMOVED***
				"filename": data.Filename,
				"vus":      data.VUs,
			***REMOVED***).Debug("Running script")

			var r runner.Runner = nil

			r, err := js.New()
			if err != nil ***REMOVED***
				ch <- message.NewToClient("run.error", message.Fields***REMOVED***"error": err***REMOVED***)
				break
			***REMOVED***

			err = r.Load(data.Filename, data.Source)
			if err != nil ***REMOVED***
				ch <- message.NewToClient("run.error", message.Fields***REMOVED***"error": err***REMOVED***)
				break
			***REMOVED***

			for res := range runner.Run(r, data.VUs, p.stopChannel) ***REMOVED***
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
		case "run.stop":
			close(p.stopChannel)
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
