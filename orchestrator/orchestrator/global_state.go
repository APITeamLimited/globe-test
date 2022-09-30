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
	consoleWriter struct ***REMOVED***
		ctx            context.Context
		client         *redis.Client
		jobId          string
		orchestratorId string
	***REMOVED***

	globalState struct ***REMOVED***
		ctx            context.Context
		logger         *logrus.Logger
		client         *redis.Client
		jobId          string
		orchestratorId string
		metricsStore   libOrch.BaseMetricsStore
	***REMOVED***
)

var _ libOrch.BaseGlobalState = &globalState***REMOVED******REMOVED***

func (g *globalState) Ctx() context.Context ***REMOVED***
	return g.ctx
***REMOVED***

func (g *globalState) Logger() *logrus.Logger ***REMOVED***
	return g.logger
***REMOVED***

func (g *globalState) Client() *redis.Client ***REMOVED***
	return g.client
***REMOVED***

func (g *globalState) JobId() string ***REMOVED***
	return g.jobId
***REMOVED***

func (g *globalState) OrchestratorId() string ***REMOVED***
	return g.orchestratorId
***REMOVED***

func (g *globalState) MetricsStore() *libOrch.BaseMetricsStore ***REMOVED***
	return &g.metricsStore
***REMOVED***

func NewGlobalState(ctx context.Context, client *redis.Client, jobId string, orchestratorId string) *globalState ***REMOVED***
	redisStdOut := &consoleWriter***REMOVED***ctx, client, jobId, orchestratorId***REMOVED***

	logger := &logrus.Logger***REMOVED***
		Out:       redisStdOut,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	***REMOVED***

	metricsStore := orchMetrics.NewCachedMetricsStore(ctx, client, orchestratorId, jobId)

	return &globalState***REMOVED***
		ctx:            ctx,
		logger:         logger,
		client:         client,
		jobId:          jobId,
		orchestratorId: orchestratorId,
		metricsStore:   metricsStore,
	***REMOVED***
***REMOVED***

var _ io.Writer = &consoleWriter***REMOVED******REMOVED***

func (w *consoleWriter) Write(p []byte) (n int, err error) ***REMOVED***
	origLen := len(p)

	// Intercept the write message so can assess log errors parse json
	parsed := make(map[string]interface***REMOVED******REMOVED***)
	if err := json.Unmarshal(p, &parsed); err != nil ***REMOVED***

		return origLen, err
	***REMOVED***

	// Check message level, if error then log error
	if parsed["level"] == "error" ***REMOVED***
		if parsed["error"] != nil ***REMOVED***
			libOrch.HandleStringError(w.ctx, w.client, w.jobId, w.orchestratorId, parsed["error"].(string))
		***REMOVED*** else ***REMOVED***
			libOrch.HandleStringError(w.ctx, w.client, w.jobId, w.orchestratorId, parsed["msg"].(string))
		***REMOVED***
		return
	***REMOVED***

	libOrch.DispatchMessage(w.ctx, w.client, w.jobId, w.orchestratorId, string(p), "STDOUT")

	return origLen, err
***REMOVED***
