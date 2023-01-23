package libOrch

import (
	"context"
	"sync"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
)

type WorkerState struct {
	Status     string `json:"status"`
	WorkerId   string `json:"workerId"`
	ChildJobId string `json:"childJobId"`
}

type MessageQueue struct {
	Mutex sync.Mutex

	// The count of currently actively being sent messages
	QueueCount    int
	NewQueueCount chan int
}

type BaseGlobalState interface {
	Ctx() context.Context
	Logger() *logrus.Logger
	OrchestratorClient() *redis.Client
	JobId() string
	OrchestratorId() string
	MetricsStore() *BaseMetricsStore
	GetStatus() string
	SetStatus(string)

	GetChildJobStates() []WorkerState
	SetChildJobState(workerId string, childJobId string, status string)
	CreditsManager() *lib.CreditsManager
	Standalone() bool
	FuncMode() bool
	FuncAuthClient() FunctionAuthClient

	IndependentWorkerRedisHosts() bool

	MessageQueue() *MessageQueue
}
