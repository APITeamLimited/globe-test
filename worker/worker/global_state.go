package worker

import (
	"context"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	defaultConfigFileName   = "config.json"
	waitRemoteLoggerTimeout = time.Second * 5
)

// globalFlags contains global config values that apply for all k6 sub-commands.
type globalFlags struct ***REMOVED***
	configFilePath string
	address        string
	logOutput      string
***REMOVED***

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
type globalState struct ***REMOVED***
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

	client *redis.Client

	workerId   string
	jobId      string
	childJobId string
	status     string
***REMOVED***

var _ libWorker.BaseGlobalState = &globalState***REMOVED******REMOVED***

// Ideally, this should be the only function in the whole codebase where we use
// global variables and functions from the os package. Anywhere else, things
// like os.Stdout, os.Stderr, os.Stdin, os.Getenv(), etc. should be removed and
// the respective properties of globalState used instead.

// Care is needed to prevent leaking system info to malicious actors.

func newGlobalState(ctx context.Context, client *redis.Client, job libOrch.ChildJob, workerId string) *globalState ***REMOVED***
	gs := &globalState***REMOVED***
		ctx:          ctx,
		fs:           afero.NewMemMapFs(),
		getwd:        os.Getwd,
		args:         []string***REMOVED******REMOVED***,
		envVars:      make(map[string]string),
		stdIn:        os.Stdin,
		osExit:       os.Exit,
		signalNotify: signal.Notify,
		signalStop:   signal.Stop,
		workerId:     workerId,
		client:       client,
		jobId:        job.Id,
		childJobId:   job.ChildJobId,
	***REMOVED***

	gs.stdOut = &consoleWriter***REMOVED***gs***REMOVED***
	gs.stdErr = &consoleWriter***REMOVED***gs***REMOVED***

	gs.logger = &logrus.Logger***REMOVED***
		Out:       gs.stdOut,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	***REMOVED***

	gs.fallbackLogger = gs.logger

	confDir, err := os.UserConfigDir()
	if err != nil ***REMOVED***
		gs.logger.WithError(err).Warn("could not get config directory")
		confDir = ".config"
	***REMOVED***

	defaultFlags := getDefaultFlags(confDir)

	gs.defaultFlags = defaultFlags
	gs.flags = defaultFlags

	return gs
***REMOVED***

func getDefaultFlags(homeFolder string) globalFlags ***REMOVED***
	return globalFlags***REMOVED***
		address:        "localhost:6565",
		configFilePath: filepath.Join(homeFolder, "loadimpact", "k6", defaultConfigFileName),
		logOutput:      "stderr",
	***REMOVED***
***REMOVED***

func (gs *globalState) Ctx() context.Context ***REMOVED***
	return gs.ctx
***REMOVED***

func (gs *globalState) Client() *redis.Client ***REMOVED***
	return gs.client
***REMOVED***

func (gs *globalState) JobId() string ***REMOVED***
	return gs.jobId
***REMOVED***

func (gs *globalState) ChildJobId() string ***REMOVED***
	return gs.childJobId
***REMOVED***

func (gs *globalState) WorkerId() string ***REMOVED***
	return gs.workerId
***REMOVED***

func (gs *globalState) GetWorkerStatus() string ***REMOVED***
	return gs.status
***REMOVED***

func (gs *globalState) SetWorkerStatus(status string) ***REMOVED***
	gs.status = status
***REMOVED***
