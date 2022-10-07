package libOrch

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
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
	}

	CollectionContext struct {
		Variables []libWorker.KeyValueItem `json:"variables"`
	}

	Job struct {
		Id                   string                 `json:"id"`
		Source               string                 `json:"source"`
		SourceName           string                 `json:"sourceName"`
		ScopeId              string                 `json:"scopeId"`
		EnvironmentContext   *EnvironmentContext    `json:"environmentContext"`
		CollectionContext    *CollectionContext     `json:"collectionContext"`
		RestRequest          map[string]interface{} `json:"restRequest"`
		AssignedOrchestrator string                 `json:"assignedOrchestrator"`
	}

	ChildJob struct {
		Job
		ChildJobId string            `json:"childJobId"`
		Options    libWorker.Options `json:"options"`
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
		Time        time.Time `json:"time"`
		WorkerId    string    `json:"workerId"`
		Message     string    `json:"message"`
		MessageType string    `json:"messageType"`
	}

	OrchestratorOrWorkerMessage struct {
		JobId          string    `json:"jobId"`
		Time           time.Time `json:"time"`
		OrchestratorId string    `json:"orchestratorId"`
		WorkerId       string    `json:"workerId"`
		Message        string    `json:"message"`
		MessageType    string    `json:"messageType"`
	}

	MarkMessage struct {
		Mark    string      `json:"mark"`
		Message interface{} `json:"message"`
	}
)
