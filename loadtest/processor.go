package loadtest

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/message"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/runner/js"
	"github.com/loadimpact/speedboat/worker"
)

func init() ***REMOVED***
	worker.RegisterProcessor(func(*worker.Worker) master.Processor ***REMOVED***
		return &LoadTestProcessor***REMOVED******REMOVED***
	***REMOVED***)
***REMOVED***

type LoadTestProcessor struct ***REMOVED***
	// Write a positive number to this to spawn so many VUs, negative to kill
	// that many. Close it to kill all VUs and end the running test.
	controlChannel chan int

	// Counter for how many VUs we currently have running.
	currentVUs int
***REMOVED***

func (p *LoadTestProcessor) Process(msg message.Message) <-chan message.Message ***REMOVED***
	ch := make(chan message.Message)

	go func() ***REMOVED***
		defer close(ch)

		switch msg.Type ***REMOVED***
		case "test.run":
			data := MessageTestRun***REMOVED******REMOVED***
			if err := msg.Take(&data); err != nil ***REMOVED***
				ch <- message.ToClient("error").WithError(err)
				return
			***REMOVED***

			p.controlChannel = make(chan int, 1)
			p.currentVUs = data.VUs

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

			p.controlChannel <- data.VUs
			for res := range runner.Run(r, p.controlChannel) ***REMOVED***
				switch res := res.(type) ***REMOVED***
				case runner.LogEntry:
					ch <- message.ToClient("test.log").With(res)
				case runner.Metric:
					ch <- message.ToClient("test.metric").With(res)
				case error:
					ch <- message.ToClient("error").WithError(res)
				***REMOVED***
			***REMOVED***
		case "test.scale":
			data := MessageTestScale***REMOVED******REMOVED***
			if err := msg.Take(&data); err != nil ***REMOVED***
				ch <- message.ToClient("error").WithError(err)
				return
			***REMOVED***

			delta := data.VUs - p.currentVUs
			log.WithFields(log.Fields***REMOVED***
				"from":  p.currentVUs,
				"to":    data.VUs,
				"delta": delta,
			***REMOVED***).Debug("Scaling")
			p.controlChannel <- delta
			p.currentVUs = data.VUs
		case "test.stop":
			close(p.controlChannel)
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
