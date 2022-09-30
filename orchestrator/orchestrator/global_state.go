package orchestrator

import (
	"context"
	"encoding/json"
	"io"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/orchMetrics"
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

	globalState struct {
		ctx            context.Context
		logger         *logrus.Logger
		client         *redis.Client
		jobId          string
		orchestratorId string
		metricsStore   libOrch.BaseMetricsStore
	}
)

var _ libOrch.BaseGlobalState = &globalState{}

func (g *globalState) Ctx() context.Context {
	return g.ctx
}

func (g *globalState) Logger() *logrus.Logger {
	return g.logger
}

func (g *globalState) Client() *redis.Client {
	return g.client
}

func (g *globalState) JobId() string {
	return g.jobId
}

func (g *globalState) OrchestratorId() string {
	return g.orchestratorId
}

func (g *globalState) MetricsStore() *libOrch.BaseMetricsStore {
	return &g.metricsStore
}

func NewGlobalState(ctx context.Context, client *redis.Client, jobId string, orchestratorId string) *globalState {
	redisStdOut := &consoleWriter{ctx, client, jobId, orchestratorId}

	logger := &logrus.Logger{
		Out:       redisStdOut,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	metricsStore := orchMetrics.NewCachedMetricsStore(ctx, client, orchestratorId, jobId)

	return &globalState{
		ctx:            ctx,
		logger:         logger,
		client:         client,
		jobId:          jobId,
		orchestratorId: orchestratorId,
		metricsStore:   metricsStore,
	}
}

var _ io.Writer = &consoleWriter{}

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
			libOrch.HandleStringError(w.ctx, w.client, w.jobId, w.orchestratorId, parsed["error"].(string))
		} else {
			libOrch.HandleStringError(w.ctx, w.client, w.jobId, w.orchestratorId, parsed["msg"].(string))
		}
		return
	}

	libOrch.DispatchMessage(w.ctx, w.client, w.jobId, w.orchestratorId, string(p), "STDOUT")

	return origLen, err
}
