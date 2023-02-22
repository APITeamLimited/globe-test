package aggregator

import (
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/metrics"
)

type intervalWithSubfraction struct {
	subFraction float64
	location    string
	interval    *Interval
}

type periodIntervals struct {
	period    int32
	intervals []*intervalWithSubfraction
}

// Cached metrics are stored before being collated and sent
type aggregator struct {
	gs libOrch.BaseGlobalState

	childJobs     map[string]libOrch.ChildJobDistribution
	childJobCount int

	intervals             []*periodIntervals
	lastIntervalCleanupAt time.Time
	intervalsMutex        sync.Mutex

	consoleMessageHashes  []string
	consoleMessages       []*ConsoleMessage
	consoleMessagesTicker *time.Ticker
	sentMaxLogsMessage    bool
	consoleMutex          sync.Mutex

	previousIntervals []*Interval

	thresholds               map[string]metrics.Thresholds
	evaluateThresholdsTicker *time.Ticker
	thresholdsMutex          sync.Mutex
	thresholdStartTime       *time.Time
}

var (
	_ libOrch.BaseAggregator = &aggregator{}
)

func NewCachedMetricsStore(gs libOrch.BaseGlobalState) *aggregator {
	return &aggregator{
		gs: gs,

		intervals:             make([]*periodIntervals, 0),
		lastIntervalCleanupAt: time.Now(),
		intervalsMutex:        sync.Mutex{},

		consoleMessageHashes: make([]string, 0),
		consoleMessages:      make([]*ConsoleMessage, 0),
		consoleMutex:         sync.Mutex{},

		sentMaxLogsMessage: false,

		previousIntervals: nil,

		thresholdsMutex:    sync.Mutex{},
		thresholdStartTime: nil,
	}
}

func (aggregator *aggregator) InitAggregator(childJobs map[string]libOrch.ChildJobDistribution, thresholds map[string]metrics.Thresholds) {
	aggregator.lockAllMutexes()
	defer aggregator.unlockAllMutexes()

	aggregator.childJobs = childJobs
	aggregator.thresholds = thresholds

	// Determine total number of child jobs
	for _, childJob := range childJobs {
		aggregator.childJobCount += len(childJob.ChildJobs)
	}
}

func (aggregator *aggregator) StartConsoleLogging() {
	aggregator.consoleMutex.Lock()
	defer aggregator.consoleMutex.Unlock()

	aggregator.consoleMessagesTicker = time.NewTicker(1 * time.Second)

	go func() {
		for range aggregator.consoleMessagesTicker.C {
			err := aggregator.flushConsoleMessages()
			if err != nil {
				libOrch.HandleError(aggregator.gs, fmt.Errorf("error flushing console messages: %v", err))
			}
		}
	}()
}

func (aggregator *aggregator) StartThresholdEvaluation() {
	aggregator.intervalsMutex.Lock()
	aggregator.thresholdsMutex.Lock()
	defer func() {
		aggregator.intervalsMutex.Unlock()
		aggregator.thresholdsMutex.Unlock()
	}()

	timeStart := time.Now()

	aggregator.thresholdStartTime = &timeStart
	aggregator.evaluateThresholdsTicker = time.NewTicker(2 * time.Second)

	go func() {
		for range aggregator.evaluateThresholdsTicker.C {
			aggregator.evaluateThresholds()
		}
	}()
}

func (aggregator *aggregator) StopAndCleanup() {
	if aggregator.consoleMessagesTicker != nil {
		aggregator.consoleMessagesTicker.Stop()
	}

	if aggregator.evaluateThresholdsTicker != nil {
		aggregator.evaluateThresholdsTicker.Stop()
	}

	aggregator = nil
}

func (aggregator *aggregator) lockAllMutexes() {
	aggregator.intervalsMutex.Lock()
	aggregator.thresholdsMutex.Lock()
	aggregator.consoleMutex.Lock()
}

func (aggregator *aggregator) unlockAllMutexes() {
	aggregator.intervalsMutex.Unlock()
	aggregator.thresholdsMutex.Unlock()
	aggregator.consoleMutex.Unlock()
}
