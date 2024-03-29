package worker

import (
	"context"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	defaultConfigFileName   = "config.json"
	waitRemoteLoggerTimeout = time.Second * 5
)

// globalFlags contains global config values that apply for all k6 sub-commands.
type globalFlags struct {
	configFilePath string
	address        string
	logOutput      string
}

// globalState contains the globalFlags and accessors for most of the global
// process-external state like CLI arguments, env vars, standard input, output
// and error, etc. In practice, most of it is normally accessed through the `os`
// package from the Go stdlibWorker.
//
// We group them here so we can prevent direct access to them from the rest of
// the k6 codebase. This gives us the ability to mock them and have robust and
// easy-to-write integration-like tests to check the k6 end-to-end behavior in
// any simulated conditions.
//
// `newGlobalState()` returns a globalState object with the real `os`
// parameters, while `newGlobalTestState()` can be used in tests to create
// simulated environments.
type globalState struct {
	ctx context.Context

	fs      afero.Fs
	getwd   func() (string, error)
	args    []string
	envVars map[string]string

	defaultFlags, flags globalFlags

	stdOut, stdErr *consoleWriter
	stdIn          io.Reader

	osExit       func(int)
	signalNotify func(chan<- os.Signal, ...os.Signal)
	signalStop   func(chan<- os.Signal)

	logger         *logrus.Logger
	fallbackLogger *logrus.Logger

	conn           *websocket.Conn
	connWriteMutex *sync.Mutex
	connReadMutex  *sync.Mutex

	workerId   string
	jobId      string
	childJobId string
	status     string

	funcModeEnabled bool
	funcModeInfo    *lib.FuncModeInfo
	messageQueue    *libWorker.MessageQueue

	cancelFunc context.CancelFunc
}

var _ libWorker.BaseGlobalState = &globalState{}

// Ideally, this should be the only function in the whole codebase where we use
// global variables and functions from the os package. Anywhere else, things
// like os.Stdout, os.Stderr, os.Stdin, os.Getenv(), etc. should be removed and
// the respective properties of globalState used instead.

// Care is needed to prevent leaking system info to malicious actors.

func newGlobalState(ctx context.Context, conn *websocket.Conn, job *libOrch.ChildJob, workerId string, funcModeInfo *lib.FuncModeInfo, connReadMutex, connWriteMutex *sync.Mutex) *globalState {
	gs := &globalState{
		ctx:            ctx,
		fs:             afero.NewMemMapFs(),
		getwd:          os.Getwd,
		args:           []string{},
		envVars:        make(map[string]string),
		stdIn:          os.Stdin,
		osExit:         os.Exit,
		signalNotify:   signal.Notify,
		signalStop:     signal.Stop,
		workerId:       workerId,
		conn:           conn,
		connWriteMutex: connWriteMutex,
		connReadMutex:  connReadMutex,
		jobId:          job.Id,
		childJobId:     job.ChildJobId,
		funcModeInfo:   funcModeInfo,
		messageQueue: &libWorker.MessageQueue{
			Mutex:         sync.Mutex{},
			QueueCount:    0,
			NewQueueCount: make(chan int),
		},
		cancelFunc: func() {},
	}

	loggerChannel := make(chan map[string]interface{}, 100)

	cw := &consoleWriter{
		gs,
		loggerChannel,
	}

	gs.stdOut = cw
	gs.stdErr = cw

	gs.logger = &logrus.Logger{
		Out:       gs.stdOut,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	gs.fallbackLogger = gs.logger

	confDir, err := os.UserConfigDir()
	if err != nil {
		gs.logger.WithError(err).Warn("could not get config directory")
		confDir = ".config"
	}

	defaultFlags := getDefaultFlags(confDir)

	gs.defaultFlags = defaultFlags
	gs.flags = defaultFlags

	return gs
}

func getDefaultFlags(homeFolder string) globalFlags {
	return globalFlags{
		address:        "localhost:6565",
		configFilePath: filepath.Join(homeFolder, "loadimpact", "k6", defaultConfigFileName),
		logOutput:      "stderr",
	}
}

func (gs *globalState) Ctx() context.Context {
	return gs.ctx
}

func (gs *globalState) Conn() *websocket.Conn {
	return gs.conn
}

func (gs *globalState) ConnWriteMutex() *sync.Mutex {
	return gs.connWriteMutex
}

func (gs *globalState) ConnReadMutex() *sync.Mutex {
	return gs.connReadMutex
}

func (gs *globalState) JobId() string {
	return gs.jobId
}

func (gs *globalState) ChildJobId() string {
	return gs.childJobId
}

func (gs *globalState) WorkerId() string {
	return gs.workerId
}

func (gs *globalState) GetWorkerStatus() string {
	return gs.status
}

func (gs *globalState) SetWorkerStatus(status string) {
	gs.status = status
}

func (gs *globalState) FuncModeEnabled() bool {
	return gs.funcModeEnabled
}

func (gs *globalState) FuncModeInfo() *lib.FuncModeInfo {
	return gs.funcModeInfo
}

func (gs *globalState) MessageQueue() *libWorker.MessageQueue {
	return gs.messageQueue
}

func (gs *globalState) GetRunAbortFunc() context.CancelFunc {
	return gs.cancelFunc
}

func (gs *globalState) SetRunAbortFunc(cancelFunc context.CancelFunc) {
	gs.cancelFunc = cancelFunc
}

func (gs *globalState) GetLoggerChannel() chan map[string]interface{} {
	return gs.stdOut.loggerChannel
}
