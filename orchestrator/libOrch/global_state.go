package libOrch

import (
	"context"

	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
)

type BaseGlobalState interface {
	Ctx() context.Context
	Logger() *logrus.Logger
	Client() *redis.Client
	JobId() string
	OrchestratorId() string
	MetricsStore() *BaseMetricsStore
}
