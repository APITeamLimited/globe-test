package libOrch

import (
	"context"
	"sync"
	"time"

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
	GetStatus() string
	SetStatus(string)

	CreditsManager() *lib.CreditsManager
	Standalone() bool
	FuncAuthClient() RunAuthClient
	LoadZones() []string
	MessageQueue() *MessageQueue

	GetCurrentTestRunDuration() time.Duration
}
