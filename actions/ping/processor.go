package ping

import (
	"github.com/loadimpact/speedboat/comm"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/worker"
)

func init() ***REMOVED***
	master.RegisterProcessor(func(*master.Master) comm.Processor ***REMOVED***
		return &PingProcessor***REMOVED******REMOVED***
	***REMOVED***)
	worker.RegisterProcessor(func(*worker.Worker) comm.Processor ***REMOVED***
		return &PingProcessor***REMOVED******REMOVED***
	***REMOVED***)
***REMOVED***

// Processes pings, on both master and worker.
type PingProcessor struct***REMOVED******REMOVED***

func (p *PingProcessor) Process(msg comm.Message) <-chan comm.Message ***REMOVED***
	out := make(chan comm.Message)

	go func() ***REMOVED***
		defer close(out)
		switch msg.Type ***REMOVED***
		case "ping.ping":
			data := PingMessage***REMOVED******REMOVED***
			if err := msg.Take(&data); err != nil ***REMOVED***
				out <- comm.ToClient("error").WithError(err)
				break
			***REMOVED***
			for res := range p.ProcessPing(data) ***REMOVED***
				out <- res
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return out
***REMOVED***

func (p *PingProcessor) ProcessPing(data PingMessage) <-chan comm.Message ***REMOVED***
	ch := make(chan comm.Message)

	go func() ***REMOVED***
		defer close(ch)

		ch <- comm.ToClient("ping.pong").With(data)
	***REMOVED***()

	return ch
***REMOVED***
