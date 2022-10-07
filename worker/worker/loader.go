// Custom loader for execution from redis

package worker

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/errext"
	"github.com/APITeamLimited/globe-test/worker/errext/exitcodes"
	"github.com/APITeamLimited/globe-test/worker/js"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/spf13/afero"
	"gopkg.in/guregu/null.v3"
)

func loadAndConfigureTest(
	gs *globalState,
	job libOrch.ChildJob,
	workerInfo *libWorker.WorkerInfo,
) (*workerLoadedAndConfiguredTest, error) ***REMOVED***
	sourceName := job.SourceName

	if sourceName == "" ***REMOVED***
		return nil, fmt.Errorf("sourceName not found on job, this is probably a bug")
	***REMOVED***

	stringSource := job.Source

	if stringSource == "" ***REMOVED***
		return nil, fmt.Errorf("source not found on job, this is probably a bug")
	***REMOVED***

	source := &loader.SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Path: sourceName***REMOVED***,
		Data: []byte(stringSource),
	***REMOVED***

	filesystems := map[string]afero.Fs***REMOVED***
		"file": afero.NewMemMapFs(),
	***REMOVED***

	f, err := afero.TempFile(filesystems["file"], "", sourceName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	_, err = f.Write([]byte(stringSource))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Store the source in the filesystem
	sourceRootPath := sourceName

	// For now runtime options are constant for all tests
	// TODO: make this configurable
	runtimeOptions := libWorker.RuntimeOptions***REMOVED***
		TestType:             null.StringFrom(testTypeJS),
		IncludeSystemEnvVars: null.BoolFrom(false),
		CompatibilityMode:    null.StringFrom("extended"),
		NoThresholds:         null.BoolFrom(false),
		SummaryExport:        null.StringFrom(""),
		Env:                  make(map[string]string),
	***REMOVED***

	registry := workerMetrics.NewRegistry()

	preInitState := &libWorker.TestPreInitState***REMOVED***
		// These gs will need to be changed as on the cloud
		Logger:         gs.logger,
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: workerMetrics.RegisterBuiltinMetrics(registry),
	***REMOVED***

	test := &workerLoadedTest***REMOVED***
		sourceRootPath: sourceRootPath,
		source:         source,
		fs:             gs.fs,
		pwd:            "",
		fileSystems:    filesystems,
		preInitState:   preInitState,
	***REMOVED***

	gs.logger.Debugf("Initializing k6 runner for '%s' (%s)...", sourceRootPath)
	if err := test.initializeFirstRunner(gs, workerInfo); err != nil ***REMOVED***
		return nil, fmt.Errorf("could not initialize '%s': %w", sourceRootPath, err)
	***REMOVED***
	gs.logger.Debug("Runner successfully initialized!")

	return test.consolidateDeriveAndValidateConfig(gs, job)
***REMOVED***

func (lt *workerLoadedTest) initializeFirstRunner(gs *globalState, workerInfo *libWorker.WorkerInfo) error ***REMOVED***
	testPath := lt.source.URL.String()
	logger := gs.logger.WithField("test_path", testPath)

	if lt.preInitState.RuntimeOptions.KeyWriter.Valid ***REMOVED***

		logger.Warnf("SSLKEYLOGFILE was specified, logging TLS connection keys to '%s'...",
			lt.preInitState.RuntimeOptions.KeyWriter.String)
		keylogFilename := lt.preInitState.RuntimeOptions.KeyWriter.String
		// if path is absolute - no point doing anything
		if !filepath.IsAbs(keylogFilename) ***REMOVED***
			// filepath.Abs could be used but it will get the pwd from `os` package instead of what is in lt.pwd
			// this is against our general approach of not using `os` directly and makes testing harder
			keylogFilename = filepath.Join(lt.pwd, keylogFilename)
		***REMOVED***
		f, err := lt.fs.OpenFile(keylogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
		if err != nil ***REMOVED***
			return fmt.Errorf("couldn't get absolute path for keylog file: %w", err)
		***REMOVED***
		lt.keyLogger = f
		lt.preInitState.KeyLogger = &consoleWriter***REMOVED***
			ctx:      gs.ctx,
			client:   workerInfo.Client,
			jobId:    workerInfo.JobId,
			workerId: workerInfo.WorkerId,
		***REMOVED***
	***REMOVED***

	runner, err := js.New(lt.preInitState, lt.source, lt.fileSystems, workerInfo)
	// TODO: should we use common.UnwrapGojaInterruptedError() here?
	if err != nil ***REMOVED***
		return fmt.Errorf("could not load JS test '%s': %w", testPath, err)
	***REMOVED***
	lt.initRunner = runner
	return nil

***REMOVED***

func (lt *workerLoadedTest) consolidateDeriveAndValidateConfig(
	gs *globalState, job libOrch.ChildJob,
) (*workerLoadedAndConfiguredTest, error) ***REMOVED***
	// Options have already been determined by the orchestrator, just load them
	consolidatedConfig := Config***REMOVED***
		Options: job.Options,
	***REMOVED***

	// Parse the thresholds, only if the --no-threshold flag is not set.
	// If parsing the threshold expressions failed, consider it as an
	// invalid configuration error.
	if !lt.preInitState.RuntimeOptions.NoThresholds.Bool ***REMOVED***
		for metricName, thresholdsDefinition := range consolidatedConfig.Options.Thresholds ***REMOVED***
			err := thresholdsDefinition.Parse()
			if err != nil ***REMOVED***
				return nil, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			***REMOVED***

			err = thresholdsDefinition.Validate(metricName, lt.preInitState.Registry)
			if err != nil ***REMOVED***
				return nil, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	derivedConfig, err := deriveAndValidateConfig(consolidatedConfig, lt.initRunner.IsExecutable)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &workerLoadedAndConfiguredTest***REMOVED***
		workerLoadedTest: lt,
		derivedConfig:    derivedConfig,
	***REMOVED***, nil
***REMOVED***

func deriveAndValidateConfig(
	conf Config, isExecutable func(string) bool,
) (result Config, err error) ***REMOVED***
	result = conf
	err = validateConfig(result, isExecutable)
	return result, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
***REMOVED***

func validateConfig(conf Config, isExecutable func(string) bool) error ***REMOVED***
	errList := conf.Validate()

	for _, ec := range conf.Scenarios ***REMOVED***
		if err := validateScenarioConfig(ec, isExecutable); err != nil ***REMOVED***
			errList = append(errList, err)
		***REMOVED***
	***REMOVED***

	return consolidateErrorMessage(errList, "There were problems with the specified script configuration:")
***REMOVED***

func consolidateErrorMessage(errList []error, title string) error ***REMOVED***
	if len(errList) == 0 ***REMOVED***
		return nil
	***REMOVED***

	errMsgParts := []string***REMOVED***title***REMOVED***
	for _, err := range errList ***REMOVED***
		errMsgParts = append(errMsgParts, fmt.Sprintf("\t- %s", err.Error()))
	***REMOVED***

	return errors.New(strings.Join(errMsgParts, "\n"))
***REMOVED***

func validateScenarioConfig(conf libWorker.ExecutorConfig, isExecutable func(string) bool) error ***REMOVED***
	execFn := conf.GetExec()
	if !isExecutable(execFn) ***REMOVED***
		return fmt.Errorf("executor %s: function '%s' not found in exports", conf.GetName(), execFn)
	***REMOVED***
	return nil
***REMOVED***
