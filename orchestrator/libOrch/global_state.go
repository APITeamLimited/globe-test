package libOrch

import (
	"context"

	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
)

type WorkerState struct {
	Status     string `json:"status"`
	WorkerId   string `json:"workerId"`
	ChildJobId string `json:"childJobId"`
}

type BaseGlobalState interface {
	Ctx() context.Context
	Logger() *logrus.Logger
	// The orchestrator client
	Client() *redis.Client
	JobId() string
	OrchestratorId() string
	MetricsStore() *BaseMetricsStore
	GetStatus() string
	SetStatus(string)

	GetChildJobStates() []WorkerState
	SetChildJobState(workerId string, childJobId string, status string)
}
