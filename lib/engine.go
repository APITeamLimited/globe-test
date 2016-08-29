package lib

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

type Engine struct ***REMOVED***
	Runner Runner
	Status Status

	ctx       context.Context
	cancelers []context.CancelFunc
	mutex     sync.Mutex
***REMOVED***

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	e.ctx = ctx

	e.Status.StartTime = time.Now()
	e.Status.Running = true
	defer func() ***REMOVED***
		e.Status.Running = false
		e.Status.VUs = 0
	***REMOVED***()

	<-ctx.Done()

	return nil
***REMOVED***

func (e *Engine) Scale(vus int64) error ***REMOVED***
	e.mutex.Lock()
	defer e.mutex.Unlock()

	l := int64(len(e.cancelers))
	switch ***REMOVED***
	case l < vus:
		for i := int64(len(e.cancelers)); i < vus; i++ ***REMOVED***
			vu, err := e.Runner.NewVU()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if err := vu.Reconfigure(i + 1); err != nil ***REMOVED***
				return err
			***REMOVED***

			ctx, cancel := context.WithCancel(e.ctx)
			e.cancelers = append(e.cancelers, cancel)
			go e.runVU(ctx, vu)
		***REMOVED***
	case l > vus:
		for _, cancel := range e.cancelers[vus+1:] ***REMOVED***
			cancel()
		***REMOVED***
		e.cancelers = e.cancelers[:vus]
	***REMOVED***

	e.Status.VUs = int64(len(e.cancelers))

	return nil
***REMOVED***

func (e *Engine) runVU(ctx context.Context, vu VU) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
			if err := vu.RunOnce(ctx); err != nil ***REMOVED***
				log.WithError(err).Error("Runtime Error")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
