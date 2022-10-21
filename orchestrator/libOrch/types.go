package libOrch

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

type (
	Distribution struct {
		LoadZone string `json:"loadZone"`
		Percent  int    `json:"percent"`
	}

	APITeamOptions struct {
		Distribution `json:"distribution"`
	}

	EnvironmentContext struct {
		Variables []libWorker.KeyValueItem `json:"variables"`
		Name      string                   `json:"name"`
	}

	CollectionContext struct {
		Variables []libWorker.KeyValueItem `json:"variables"`
		Name      string                   `json:"name"`
	}

	Job struct {
		Id                   string                 `json:"id"`
		Source               string                 `json:"source"`
		SourceName           string                 `json:"sourceName"`
		ScopeId              string                 `json:"scopeId"`
		EnvironmentContext   *EnvironmentContext    `json:"environmentContext"`
		CollectionContext    *CollectionContext     `json:"collectionContext"`
		UnderlyingRequest    map[string]interface{} `json:"underlyingRequest"`
		FinalRequest         map[string]interface{} `json:"finalRequest"`
		AssignedOrchestrator string                 `json:"assignedOrchestrator"`
		Scope                Scope                  `json:"scope"`
		Options              *libWorker.Options     `json:"options"`
	}

	Scope struct {
		Variant         string `json:"variant"`
		VariantTargetId string `json:"variantTargetId"`
	}

	ChildJob struct {
		Job
		ChildJobId        string                 `json:"childJobId"`
		Options           libWorker.Options      `json:"options"`
		UnderlyingRequest map[string]interface{} `json:"underlyingRequest"`
		FinalRequest      map[string]interface{} `json:"finalRequest"`
		SubFraction       float64                `json:"subFraction"`
	}

	OrchestratorMessage struct {
		JobId          string    `json:"jobId"`
		Time           time.Time `json:"time"`
		OrchestratorId string    `json:"orchestratorId"`
		Message        string    `json:"message"`
		MessageType    string    `json:"messageType"`
	}

	WorkerMessage struct {
		JobId       string    `json:"jobId"`
		ChildJobId  string    `json:"childJobId"`
		Time        time.Time `json:"time"`
		WorkerId    string    `json:"workerId"`
		Message     string    `json:"message"`
		MessageType string    `json:"messageType"`
	}

	OrchestratorOrWorkerMessage struct {
		JobId          string    `json:"jobId"`
		Time           time.Time `json:"time"`
		ChildJobId     string    `json:"childJobId"`
		OrchestratorId string    `json:"orchestratorId"`
		WorkerId       string    `json:"workerId"`
		Message        string    `json:"message"`
		MessageType    string    `json:"messageType"`
	}

	MarkMessage struct {
		Mark    string      `json:"mark"`
		Message interface{} `json:"message"`
	}

	WorkerClients struct {
		Clients       map[string]*NamedClient
		DefaultClient *NamedClient
	}

	NamedClient struct {
		Name   string
		Client *redis.Client
	}
)
