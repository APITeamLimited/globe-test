package worker

import (
	"context"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/APITeamLimited/redis/v9"
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
	noColor        bool
	address        string
	logOutput      string
	logFormat      string
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
}

// Ideally, this should be the only function in the whole codebase where we use
// global variables and functions from the os package. Anywhere else, things
// like os.Stdout, os.Stderr, os.Stdin, os.Getenv(), etc. should be removed and
// the respective properties of globalState used instead.

// Care is needed to prevent leaking system info to malicious actors.

func newGlobalState(ctx context.Context, client *redis.Client, jobId string, workerId string) *globalState {
	redisStdOut := &consoleWriter{ctx, client, jobId, workerId}
	redisStdErr := &consoleWriter{ctx, client, jobId, workerId}

	envVars := make(map[string]string)

	logger := &logrus.Logger{
		Out:       redisStdOut,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	confDir, err := os.UserConfigDir()
	if err != nil {
		logger.WithError(err).Warn("could not get config directory")
		confDir = ".config"
	}

	defaultFlags := getDefaultFlags(confDir)

	return &globalState{
		ctx:            ctx,
		fs:             afero.NewMemMapFs(),
		getwd:          os.Getwd,
		args:           []string{},
		envVars:        envVars,
		defaultFlags:   defaultFlags,
		flags:          getFlags(defaultFlags, envVars),
		stdOut:         redisStdOut,
		stdErr:         redisStdErr,
		stdIn:          os.Stdin,
		osExit:         os.Exit,
		signalNotify:   signal.Notify,
		signalStop:     signal.Stop,
		logger:         logger,
		fallbackLogger: logger,
	}
}

func getDefaultFlags(homeFolder string) globalFlags {
	return globalFlags{
		address:        "localhost:6565",
		configFilePath: filepath.Join(homeFolder, "loadimpact", "k6", defaultConfigFileName),
		logOutput:      "stderr",
	}
}

func getFlags(defaultFlags globalFlags, env map[string]string) globalFlags {
	result := defaultFlags

	// TODO: add env vars for the rest of the values (after adjusting
	// rootCmdPersistentFlagSet(), of course)

	if val, ok := env["K6_CONFIG"]; ok {
		result.configFilePath = val
	}
	if val, ok := env["K6_LOG_OUTPUT"]; ok {
		result.logOutput = val
	}
	if val, ok := env["K6_LOG_FORMAT"]; ok {
		result.logFormat = val
	}
	if env["K6_NO_COLOR"] != "" {
		result.noColor = true
	}
	// Support https://no-color.org/, even an empty value should disable the
	// color output from k6.
	if _, ok := env["NO_COLOR"]; ok {
		result.noColor = true
	}
	return result
}
