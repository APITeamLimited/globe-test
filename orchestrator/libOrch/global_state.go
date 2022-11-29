package libOrch

import (
	"context"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
)

type WorkerState struct ***REMOVED***
	Status     string `json:"status"`
	WorkerId   string `json:"workerId"`
	ChildJobId string `json:"childJobId"`
***REMOVED***

type BaseGlobalState interface ***REMOVED***
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
***REMOVED***
