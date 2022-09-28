package libOrch

import (
	"context"
	"encoding/json"

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

	GlobalState struct ***REMOVED***
		Ctx            context.Context
		Logger         *logrus.Logger
		Client         *redis.Client
		JobId          string
		OrchestratorId string
	***REMOVED***
)

func NewGlobalState(ctx context.Context, client *redis.Client, jobId string, orchestratorId string) *GlobalState ***REMOVED***
	redisStdOut := &consoleWriter***REMOVED***ctx, client, jobId, orchestratorId***REMOVED***

	logger := &logrus.Logger***REMOVED***
		Out:       redisStdOut,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	***REMOVED***

	return &GlobalState***REMOVED***
		Ctx:            ctx,
		Logger:         logger,
		Client:         client,
		JobId:          jobId,
		OrchestratorId: orchestratorId,
	***REMOVED***
***REMOVED***

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
			HandleStringError(w.ctx, w.client, w.jobId, w.orchestratorId, parsed["error"].(string))
		***REMOVED*** else ***REMOVED***
			HandleStringError(w.ctx, w.client, w.jobId, w.orchestratorId, parsed["msg"].(string))
		***REMOVED***
		return
	***REMOVED***

	DispatchMessage(w.ctx, w.client, w.jobId, w.orchestratorId, string(p), "STDOUT")

	return origLen, err
***REMOVED***
