package orchestrator

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/metrics"
	"github.com/APITeamLimited/globe-test/metrics/engine"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

type (
	consoleWriter struct {
		gs *globalState
	}

	globalState struct {
		ctx            context.Context
		logger         *logrus.Logger
		client         *redis.Client
		jobId          string
		orchestratorId string
		metricsStore   metrics.Registry
		status         string
		childJobStates []libOrch.WorkerState
		creditsManager *lib.CreditsManager
		standalone     bool
		funcAuthClient libOrch.RunAuthClient
		messageQueue   *libOrch.MessageQueue
		loadZones      []string
		registry       *metrics.Registry
		metricsEngine  *engine.MetricsEngine
		startTime      null.Time
	}
)

var _ libOrch.BaseGlobalState = &globalState{}

func NewGlobalState(ctx context.Context, orchestratorClient *redis.Client, job *libOrch.Job,
	orchestratorId string, creditsClient *redis.Client, standalone bool,
	funcAuthClient libOrch.RunAuthClient, loadZones []string) *globalState {
	gs := &globalState{
		ctx:            ctx,
		client:         orchestratorClient,
		jobId:          job.Id,
		orchestratorId: orchestratorId,
		childJobStates: []libOrch.WorkerState{},
		standalone:     standalone,
		funcAuthClient: funcAuthClient,
		messageQueue: &libOrch.MessageQueue{
			Mutex:         sync.Mutex{},
			QueueCount:    0,
			NewQueueCount: make(chan int),
		},
		loadZones: loadZones,
		registry:  metrics.NewRegistry(),

		// Haven't started yet so set to false
		startTime: null.NewTime(time.Now(), false),
	}

	if creditsClient != nil && job.FuncModeInfo != nil {
		gs.creditsManager = lib.CreateCreditsManager(ctx, job.Scope.Variant, job.Scope.VariantTargetId, creditsClient, *job.FuncModeInfo)
	}

	gs.logger = &logrus.Logger{
		Out:       &consoleWriter{gs: gs},
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	return gs
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
			libOrch.HandleStringError(w.gs, parsed["error"].(string))
		} else {
			libOrch.HandleStringError(w.gs, parsed["msg"].(string))
		}
		return
	}

	libOrch.DispatchMessage(w.gs, string(p), "STDOUT")

	return origLen, err
}

func (g *globalState) Ctx() context.Context {
	return g.ctx
}

func (g *globalState) Logger() *logrus.Logger {
	return g.logger
}

func (g *globalState) OrchestratorClient() *redis.Client {
	return g.client
}

func (g *globalState) JobId() string {
	return g.jobId
}

func (g *globalState) OrchestratorId() string {
	return g.orchestratorId
}

func (g *globalState) GetStatus() string {
	return g.status
}

func (g *globalState) SetStatus(status string) {
	if status == "RUNNING" {
		g.startTime = null.NewTime(time.Now(), true)
	}

	g.status = status
}

func (g *globalState) CreditsManager() *lib.CreditsManager {
	return g.creditsManager
}

func (g *globalState) Standalone() bool {
	return g.standalone
}

func (g *globalState) FuncAuthClient() libOrch.RunAuthClient {
	return g.funcAuthClient
}

func (g *globalState) MessageQueue() *libOrch.MessageQueue {
	return g.messageQueue
}

func (g *globalState) LoadZones() []string {
	return g.loadZones
}

func (g *globalState) GetCurrentTestRunDuration() time.Duration {
	if g.startTime.Valid == false {
		return 0
	}

	return time.Since(g.startTime.Time)
}
