package libOrch

import (
	"context"

	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
)

type BaseGlobalState interface ***REMOVED***
	Ctx() context.Context
	Logger() *logrus.Logger
	// The orchestrator client
	Client() *redis.Client
	JobId() string
	OrchestratorId() string
	MetricsStore() *BaseMetricsStore
***REMOVED***
