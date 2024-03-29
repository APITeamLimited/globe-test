package orchestrator

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/aggregator"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
)

type (
	consoleWriter struct {
		gs *globalState
	}

	globalState struct {
		ctx                  context.Context
		logger               *logrus.Logger
		client               *redis.Client
		jobId                string
		orchestratorId       string
		metricsStore         libOrch.BaseAggregator
		status               string
		childJobStates       []libOrch.WorkerState
		creditsManager       *lib.CreditsManager
		standalone           bool
		funcAuthClient       libOrch.RunAuthClient
		messageQueue         *libOrch.MessageQueue
		loadZones            []string
		statusUpdatesChannel chan string
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
		loadZones:            loadZones,
		statusUpdatesChannel: make(chan string, 10),
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

	gs.metricsStore = aggregator.NewCachedMetricsStore(gs)

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

func (g *globalState) MetricsStore() libOrch.BaseAggregator {
	return g.metricsStore
}

func (g *globalState) GetStatus() string {
	return g.status
}

func (g *globalState) SetStatus(status string) {
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

func (g *globalState) StatusUpdatesChannel() chan string {
	return g.statusUpdatesChannel
}
