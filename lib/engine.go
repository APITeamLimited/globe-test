package lib

import (
	"context"
	"time"
)

type Status struct ***REMOVED***
	StartTime time.Time `json:"startTime" yaml:"startTime"`

	Running bool  `json:"running" yaml:"running"`
	VUs     int64 `json:"vus" yaml:"vus"`
***REMOVED***

type Engine struct ***REMOVED***
	Runner Runner
	Status Status
***REMOVED***

func (e *Engine) Run(ctx context.Context) error ***REMOVED***
	e.Status.StartTime = time.Now()
	e.Status.Running = true
	defer func() ***REMOVED***
		e.Status.Running = false
		e.Status.VUs = 0
	***REMOVED***()

	<-ctx.Done()
	time.Sleep(1 * time.Second)

	return nil
***REMOVED***

func (e *Engine) Scale(vus int64) error ***REMOVED***
	e.Status.VUs = vus
	return nil
***REMOVED***
