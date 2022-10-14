package libOrch

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

type (
	Distribution struct ***REMOVED***
		LoadZone string `json:"loadZone"`
		Percent  int    `json:"percent"`
	***REMOVED***

	APITeamOptions struct ***REMOVED***
		Distribution `json:"distribution"`
	***REMOVED***

	EnvironmentContext struct ***REMOVED***
		Variables []libWorker.KeyValueItem `json:"variables"`
	***REMOVED***

	CollectionContext struct ***REMOVED***
		Variables []libWorker.KeyValueItem `json:"variables"`
	***REMOVED***

	Job struct ***REMOVED***
		Id                   string                 `json:"id"`
		Source               string                 `json:"source"`
		SourceName           string                 `json:"sourceName"`
		ScopeId              string                 `json:"scopeId"`
		EnvironmentContext   *EnvironmentContext    `json:"environmentContext"`
		CollectionContext    *CollectionContext     `json:"collectionContext"`
		RestRequest          map[string]interface***REMOVED******REMOVED*** `json:"restRequest"`
		AssignedOrchestrator string                 `json:"assignedOrchestrator"`
		Scope                Scope                  `json:"scope"`
	***REMOVED***

	Scope struct ***REMOVED***
		Variant         string `json:"variant"`
		VariantTargetId string `json:"variantTargetId"`
	***REMOVED***

	ChildJob struct ***REMOVED***
		Job
		ChildJobId string            `json:"childJobId"`
		Options    libWorker.Options `json:"options"`
	***REMOVED***

	OrchestratorMessage struct ***REMOVED***
		JobId          string    `json:"jobId"`
		Time           time.Time `json:"time"`
		OrchestratorId string    `json:"orchestratorId"`
		Message        string    `json:"message"`
		MessageType    string    `json:"messageType"`
	***REMOVED***

	WorkerMessage struct ***REMOVED***
		JobId       string    `json:"jobId"`
		Time        time.Time `json:"time"`
		WorkerId    string    `json:"workerId"`
		Message     string    `json:"message"`
		MessageType string    `json:"messageType"`
	***REMOVED***

	OrchestratorOrWorkerMessage struct ***REMOVED***
		JobId          string    `json:"jobId"`
		Time           time.Time `json:"time"`
		OrchestratorId string    `json:"orchestratorId"`
		WorkerId       string    `json:"workerId"`
		Message        string    `json:"message"`
		MessageType    string    `json:"messageType"`
	***REMOVED***

	MarkMessage struct ***REMOVED***
		Mark    string      `json:"mark"`
		Message interface***REMOVED******REMOVED*** `json:"message"`
	***REMOVED***
)
