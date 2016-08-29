package lib

import (
	"time"
)

type Status struct ***REMOVED***
	StartTime time.Time `json:"startTime" yaml:"startTime"`

	Running bool  `json:"running" yaml:"running"`
	VUs     int64 `json:"vus" yaml:"vus"`
	Pooled  int64 `json:"pooled" yaml:"pooled"`
***REMOVED***

type Info struct ***REMOVED***
	Version string `json:"version"`
***REMOVED***
