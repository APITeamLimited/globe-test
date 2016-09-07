package lib

import (
	"time"
)

type Status struct ***REMOVED***
	ID string `json:"-" jsonapi:"primary,status"`

	StartTime time.Time `json:"startTime" jsonapi:"attr,startTime" yaml:"startTime"`

	Running bool  `json:"running" jsonapi:"attr,running"`
	VUs     int64 `json:"vus" jsonapi:"attr,vus"`
	Pooled  int64 `json:"pooled" jsonapi:"attr,pooled"`
***REMOVED***

type Info struct ***REMOVED***
	ID      string `json:"-" jsonapi:"primary,info"`
	Version string `json:"version" jsonapi:"attr,version"`
***REMOVED***
