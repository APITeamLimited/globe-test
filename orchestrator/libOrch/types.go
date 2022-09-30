package libOrch

import "time"

type Distribution struct ***REMOVED***
	LoadZone string `json:"loadZone"`
	Percent  int    `json:"percent"`
***REMOVED***

type APITeamOptions struct ***REMOVED***
	Distribution `json:"distribution"`
***REMOVED***

type Job struct ***REMOVED***
	Id         string `json:"id"`
	Source     string `json:"source"`
	SourceName string `json:"sourceName"`
	Options    string `json:"options"`
***REMOVED***

type OrchestratorMessage struct ***REMOVED***
	JobId          string    `json:"jobId"`
	Time           time.Time `json:"time"`
	OrchestratorId string    `json:"orchestratorId"`
	Message        string    `json:"message"`
	MessageType    string    `json:"messageType"`
***REMOVED***

type WorkerMessage struct ***REMOVED***
	JobId       string    `json:"jobId"`
	Time        time.Time `json:"time"`
	WorkerId    string    `json:"workerId"`
	Message     string    `json:"message"`
	MessageType string    `json:"messageType"`
***REMOVED***

type OrchestratorOrWorkerMessage struct ***REMOVED***
	JobId          string    `json:"jobId"`
	Time           time.Time `json:"time"`
	OrchestratorId string    `json:"orchestratorId"`
	WorkerId       string    `json:"workerId"`
	Message        string    `json:"message"`
	MessageType    string    `json:"messageType"`
***REMOVED***

type MarkMessage struct ***REMOVED***
	Mark    string      `json:"mark"`
	Message interface***REMOVED******REMOVED*** `json:"message"`
***REMOVED***
