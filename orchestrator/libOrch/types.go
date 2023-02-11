package libOrch

import (
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/gorilla/websocket"
)

type (
	EnvironmentContext struct {
		Variables []libWorker.KeyValueItem `json:"variables"`
		Name      string                   `json:"name"`
	}

	CollectionContext struct {
		Variables []libWorker.KeyValueItem `json:"variables"`
		Name      string                   `json:"name"`
	}

	Job struct {
		Id string `json:"id"`

		EnvironmentContext *EnvironmentContext `json:"environmentContext"`
		CollectionContext  *CollectionContext  `json:"collectionContext"`

		AssignedOrchestrator string                 `json:"assignedOrchestrator"`
		Scope                Scope                  `json:"scope"`
		Options              *libWorker.Options     `json:"options"`
		VerifiedDomains      []string               `json:"verifiedDomains"`
		TestDataRaw          map[string]interface{} `json:"testData"`
		TestData             *libWorker.TestData    `json:"-"`

		FuncModeInfo           *lib.FuncModeInfo `json:"funcModeInfo"`
		PermittedLoadZones     []string          `json:"permittedLoadZones"`
		MaxTestDurationMinutes int64             `json:"maxTestDurationMinutes"`
		MaxSimulatedUsers      int64             `json:"maxSimulatedUsers"`
	}

	Scope struct {
		Variant         string `json:"variant"`
		VariantTargetId string `json:"variantTargetId"`
	}

	ChildJob struct {
		Job
		ChildJobId       string            `json:"childJobId"`
		ChildOptions     libWorker.Options `json:"childOptions"`
		SubFraction      float64           `json:"subFraction"`
		WorkerConnection *websocket.Conn
		ConnWriteMutex   *sync.Mutex
		ConnReadMutex    *sync.Mutex
	}

	ChildJobDistribution struct {
		ChildJobs []*ChildJob `json:"jobs"`
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

	LocalhostFile struct {
		FileName string `json:"fileName"`
		Contents string `json:"contents"`
		Kind     string `json:"kind"`
	}
)
