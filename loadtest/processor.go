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

			data := MessageTestRun***REMOVED******REMOVED***
			if err := msg.Take(&data); err != nil ***REMOVED***
				ch <- message.ToClient("error").WithError(err)
				return
			***REMOVED***

			log.WithFields(log.Fields***REMOVED***
				"filename": data.Filename,
				"vus":      data.VUs,
			***REMOVED***).Debug("Running script")

			var r runner.Runner = nil

			r, err := js.New()
			if err != nil ***REMOVED***
				ch <- message.ToClient("error").WithError(err)
				break
			***REMOVED***

			err = r.Load(data.Filename, data.Source)
			if err != nil ***REMOVED***
				ch <- message.ToClient("error").WithError(err)
				break
			***REMOVED***

			for res := range runner.Run(r, data.VUs, p.stopChannel) ***REMOVED***
				switch res := res.(type) ***REMOVED***
				case runner.LogEntry:
					ch <- message.ToClient("test.log").With(res)
				case runner.Metric:
					ch <- message.ToClient("test.metric").With(res)
				case error:
					ch <- message.ToClient("error").WithError(res)
				***REMOVED***
			***REMOVED***
		case "test.stop":
			close(p.stopChannel)
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
