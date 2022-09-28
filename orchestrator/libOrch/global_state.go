package libOrch

import (
	"context"
	"encoding/json"

	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
)

type (
	consoleWriter struct {
		ctx            context.Context
		client         *redis.Client
		jobId          string
		orchestratorId string
	}

	GlobalState struct {
		Ctx            context.Context
		Logger         *logrus.Logger
		Client         *redis.Client
		JobId          string
		OrchestratorId string
	}
)

func NewGlobalState(ctx context.Context, client *redis.Client, jobId string, orchestratorId string) *GlobalState {
	redisStdOut := &consoleWriter{ctx, client, jobId, orchestratorId}

	logger := &logrus.Logger{
		Out:       redisStdOut,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	return &GlobalState{
		Ctx:            ctx,
		Logger:         logger,
		Client:         client,
		JobId:          jobId,
		OrchestratorId: orchestratorId,
	}
}

func (w *consoleWriter) Write(p []byte) (n int, err error) {
	origLen := len(p)

	// Intercept the write message so can assess log errors parse json
	parsed := make(map[string]interface{})
	if err := json.Unmarshal(p, &parsed); err != nil {

		return origLen, err
	}

	// Check message level, if error then log error
	if parsed["level"] == "error" {
		if parsed["error"] != nil {
			HandleStringError(w.ctx, w.client, w.jobId, w.orchestratorId, parsed["error"].(string))
		} else {
			HandleStringError(w.ctx, w.client, w.jobId, w.orchestratorId, parsed["msg"].(string))
		}
		return
	}

	DispatchMessage(w.ctx, w.client, w.jobId, w.orchestratorId, string(p), "STDOUT")

	return origLen, err
}
