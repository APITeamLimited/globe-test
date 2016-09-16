package lib

import (
	"gopkg.in/guregu/null.v3"
	"sync"
	"time"
)

type Status struct ***REMOVED***
	Running null.Bool `json:"running"`
	VUs     null.Int  `json:"vus"`
	VUsMax  null.Int  `json:"vus-max"`
***REMOVED***

func (s Status) GetName() string ***REMOVED***
	return "status"
***REMOVED***

func (s Status) GetID() string ***REMOVED***
	return "default"
***REMOVED***

func (s Status) SetID(id string) error ***REMOVED***
	return nil
***REMOVED***

type Info struct ***REMOVED***
	Version string `json:"version"`
***REMOVED***

func (i Info) GetName() string ***REMOVED***
	return "info"
***REMOVED***

func (i Info) GetID() string ***REMOVED***
	return "default"
***REMOVED***

type Options struct ***REMOVED***
	VUs      int64         `json:"vus"`
	VUsMax   int64         `json:"vus-max"`
	Duration time.Duration `json:"duration"`
***REMOVED***

func (o Options) GetName() string ***REMOVED***
	return "options"
***REMOVED***

func (o Options) GetID() string ***REMOVED***
	return "default"
***REMOVED***

type Group struct ***REMOVED***
	Parent *Group
	Name   string
	Tests  map[string]*Test

	TestMutex sync.Mutex
***REMOVED***

type Test struct ***REMOVED***
	Group *Group
	Name  string

	Passes int64
	Fails  int64
***REMOVED***
