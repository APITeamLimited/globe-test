package run

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/actions/registry"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/message"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/runner/js"
	"github.com/loadimpact/speedboat/worker"
	"time"
)

func init() ***REMOVED***
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
			vus := int(msg.Fields["vus"].(float64))
			duration := time.Duration(msg.Fields["duration"].(float64)) * time.Millisecond

			log.WithFields(log.Fields***REMOVED***
				"filename": filename,
				"vus":      vus,
				"duration": duration,
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

			for res := range runner.Run(r, vus, duration) ***REMOVED***
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
