package lib

import (
	"gopkg.in/guregu/null.v3"
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
	VUs      int           `json:"vus"`
	VUsMax   int           `json:"vus-max"`
	Duration time.Duration `json:"duration"`

	Ext map[string]interface***REMOVED******REMOVED*** `json:"ext"`
***REMOVED***

func (o Options) GetName() string ***REMOVED***
	return "options"
***REMOVED***

func (o Options) GetID() string ***REMOVED***
	return "default"
***REMOVED***
