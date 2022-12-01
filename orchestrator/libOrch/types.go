package libOrch

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
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
		Name      string                   `json:"name"`
	***REMOVED***

	CollectionContext struct ***REMOVED***
		Variables []libWorker.KeyValueItem `json:"variables"`
		Name      string                   `json:"name"`
	***REMOVED***

	Job struct ***REMOVED***
		Id                   string                 `json:"id"`
		Source               string                 `json:"source"`
		SourceName           string                 `json:"sourceName"`
		EnvironmentContext   *EnvironmentContext    `json:"environmentContext"`
		CollectionContext    *CollectionContext     `json:"collectionContext"`
		UnderlyingRequest    map[string]interface***REMOVED******REMOVED*** `json:"underlyingRequest"`
		FinalRequest         map[string]interface***REMOVED******REMOVED*** `json:"finalRequest"`
		AssignedOrchestrator string                 `json:"assignedOrchestrator"`
		Scope                Scope                  `json:"scope"`
		Options              *libWorker.Options     `json:"options"`
		VerifiedDomains      []string               `json:"verifiedDomains"`
	***REMOVED***

	Scope struct ***REMOVED***
		Variant         string `json:"variant"`
		VariantTargetId string `json:"variantTargetId"`
	***REMOVED***

	ChildJob struct ***REMOVED***
		Job
		ChildJobId        string                 `json:"childJobId"`
		Options           libWorker.Options      `json:"options"`
		UnderlyingRequest map[string]interface***REMOVED******REMOVED*** `json:"underlyingRequest"`
		FinalRequest      map[string]interface***REMOVED******REMOVED*** `json:"finalRequest"`
		SubFraction       float64                `json:"subFraction"`
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
		ChildJobId  string    `json:"childJobId"`
		Time        time.Time `json:"time"`
		WorkerId    string    `json:"workerId"`
		Message     string    `json:"message"`
		MessageType string    `json:"messageType"`
	***REMOVED***

	OrchestratorOrWorkerMessage struct ***REMOVED***
		JobId          string    `json:"jobId"`
		Time           time.Time `json:"time"`
		ChildJobId     string    `json:"childJobId"`
		OrchestratorId string    `json:"orchestratorId"`
		WorkerId       string    `json:"workerId"`
		Message        string    `json:"message"`
		MessageType    string    `json:"messageType"`
	***REMOVED***

	MarkMessage struct ***REMOVED***
		Mark    string      `json:"mark"`
		Message interface***REMOVED******REMOVED*** `json:"message"`
	***REMOVED***

	WorkerClients struct ***REMOVED***
		Clients       map[string]*NamedClient
		DefaultClient *NamedClient
	***REMOVED***

	NamedClient struct ***REMOVED***
		Name   string
		Client *redis.Client
	***REMOVED***

	LocalhostFile struct ***REMOVED***
		FileName string `json:"fileName"`
		Contents string `json:"contents"`
		Kind     string `json:"kind"`
	***REMOVED***
)
