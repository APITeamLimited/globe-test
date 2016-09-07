package lib

type Status struct ***REMOVED***
	ID string `jsonapi:"primary,status"`

	Running     bool  `jsonapi:"attr,running"`
	ActiveVUs   int64 `jsonapi:"attr,active-vus"`
	InactiveVUs int64 `jsonapi:"attr,inactive-vus"`
***REMOVED***

type Info struct ***REMOVED***
	ID      string `jsonapi:"primary,info"`
	Version string `jsonapi:"attr,version"`
***REMOVED***
