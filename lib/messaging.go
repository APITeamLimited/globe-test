package lib

type JobUserUpdate struct ***REMOVED***
	UpdateType string `json:"updateType"`
***REMOVED***

type WrappedJobUserUpdate struct ***REMOVED***
	Update JobUserUpdate `json:"update"`
	JobId  string        `json:"jobId"`
***REMOVED***
