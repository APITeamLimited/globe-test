package lib

import (
	"context"
	"time"
)

type State struct ***REMOVED***
	StartTime time.Time `json:"startTime"`

	Running bool  `json:"running"`
	VUs     int64 `json:"vus"`
***REMOVED***

type Engine struct ***REMOVED***
	Runner Runner
	State  State
***REMOVED***

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	e.State.StartTime = time.Now()
	e.State.Running = true
	defer func() ***REMOVED***
		e.State.Running = false
		e.State.VUs = 0
	***REMOVED***()

	<-ctx.Done()
	time.Sleep(1 * time.Second)

	return nil
***REMOVED***

func (e *Engine) Scale(vus int64) error ***REMOVED***
	e.State.VUs = vus
	return nil
***REMOVED***
