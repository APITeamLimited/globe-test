// Custom loader for execution from redis

package worker

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
	"gopkg.in/guregu/null.v3"
)

func loadAndConfigureTest(
	gs *globalState,
	job map[string]string,
	workerInfo *lib.WorkerInfo,
) (*workerLoadedAndConfiguredTest, error) ***REMOVED***
	test, err := loadTest(gs, job, workerInfo)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return test.consolidateDeriveAndValidateConfig(gs, job)
***REMOVED***

func loadTest(gs *globalState, job map[string]string, workerInfo *lib.WorkerInfo) (*workerLoadedTest, error) ***REMOVED***
	sourceName := job["sourceName"]

	if sourceName == "" ***REMOVED***
		return nil, fmt.Errorf("sourceName not found on job, this is probably a bug")
	***REMOVED***

	stringSource := job["source"]

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
	runtimeOptions := lib.RuntimeOptions***REMOVED***
		TestType:             null.StringFrom(testTypeJS),
		IncludeSystemEnvVars: null.BoolFrom(false),
		CompatibilityMode:    null.StringFrom("extended"),
		NoThresholds:         null.BoolFrom(false),
		NoSummary:            null.BoolFrom(false),
		SummaryExport:        null.StringFrom(""),
		Env:                  make(map[string]string),
	***REMOVED***

	registry := metrics.NewRegistry()

	preInitState := &lib.TestPreInitState***REMOVED***
		// These gs will need to be changed as on the cloud
		Logger:         gs.logger,
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
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

	marshalled, _ := json.Marshal(test.initRunner.GetOptions())

	fmt.Println(string(marshalled))

	return test, nil
***REMOVED***

func detectTestType(data []byte) string ***REMOVED***
	if _, err := tar.NewReader(bytes.NewReader(data)).Next(); err == nil ***REMOVED***
		return testTypeArchive
	***REMOVED***
	return testTypeJS
***REMOVED***

func (lt *workerLoadedTest) initializeFirstRunner(gs *globalState, workerInfo *lib.WorkerInfo) error ***REMOVED***
	testPath := lt.source.URL.String()
	logger := gs.logger.WithField("test_path", testPath)

	testType := lt.preInitState.RuntimeOptions.TestType.String

	if testType == "" ***REMOVED***
		logger.Debug("Detecting test type for...")
		testType = detectTestType(lt.source.Data)
	***REMOVED***

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
		lt.preInitState.KeyLogger = &syncWriter***REMOVED***w: f***REMOVED***
	***REMOVED***

	switch testType ***REMOVED***
	case testTypeJS:
		runner, err := js.New(lt.preInitState, lt.source, lt.fileSystems, workerInfo)
		// TODO: should we use common.UnwrapGojaInterruptedError() here?
		if err != nil ***REMOVED***
			return fmt.Errorf("could not load JS test '%s': %w", testPath, err)
		***REMOVED***
		lt.initRunner = runner
		return nil
	default:
		return fmt.Errorf("unknown or unspecified test type '%s' for '%s'", testType, testPath)
	***REMOVED***
***REMOVED***

func (lt *workerLoadedTest) consolidateDeriveAndValidateConfig(
	gs *globalState, job map[string]string,
) (*workerLoadedAndConfiguredTest, error) ***REMOVED***

	// TODO: implement consolidateDeriveAndValidateConfig behavior

	var parsedOptions lib.Options
	err := json.Unmarshal([]byte(job["options"]), &parsedOptions)

	// Get from source data

	if err != nil ***REMOVED***
		return nil, fmt.Errorf("could not parse options: %w", err)
	***REMOVED***

	consolidatedConfig := getConsolidatedConfig(parsedOptions)

	// TODO: get other config sources eg

	// Parse the thresholds, only if the --no-threshold flag is not set.
	// If parsing the threshold expressions failed, consider it as an
	// invalid configuration error.
	if !lt.preInitState.RuntimeOptions.NoThresholds.Bool ***REMOVED***
		for metricName, thresholdsDefinition := range consolidatedConfig.Options.Thresholds ***REMOVED***
			err = thresholdsDefinition.Parse()
			if err != nil ***REMOVED***
				return nil, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			***REMOVED***

			err = thresholdsDefinition.Validate(metricName, lt.preInitState.Registry)
			if err != nil ***REMOVED***
				return nil, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	derivedConfig, err := deriveAndValidateConfig(consolidatedConfig, lt.initRunner.IsExecutable, gs.logger)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &workerLoadedAndConfiguredTest***REMOVED***
		workerLoadedTest:   lt,
		consolidatedConfig: consolidatedConfig,
		derivedConfig:      derivedConfig,
	***REMOVED***, nil
***REMOVED***

func getConsolidatedConfig(parsedOptions lib.Options) Config ***REMOVED***
	consolidatedConfig := Config***REMOVED***
		Options: parsedOptions,
	***REMOVED***

	consolidatedConfig = applyDefault(consolidatedConfig)

	return consolidatedConfig
***REMOVED***

func deriveAndValidateConfig(
	conf Config, isExecutable func(string) bool, logger logrus.FieldLogger,
) (result Config, err error) ***REMOVED***
	result = conf
	result.Options, err = executor.DeriveScenariosFromShortcuts(conf.Options, logger)
	if err == nil ***REMOVED***
		err = validateConfig(result, isExecutable)
	***REMOVED***
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

func validateScenarioConfig(conf lib.ExecutorConfig, isExecutable func(string) bool) error ***REMOVED***
	execFn := conf.GetExec()
	if !isExecutable(execFn) ***REMOVED***
		return fmt.Errorf("executor %s: function '%s' not found in exports", conf.GetName(), execFn)
	***REMOVED***
	return nil
***REMOVED***

func (lct *workerLoadedAndConfiguredTest) buildTestRunState(
	configToReinject lib.Options,
) (*lib.TestRunState, error) ***REMOVED***
	// This might be the full derived or just the consodlidated options
	if err := lct.initRunner.SetOptions(configToReinject); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// TODO: init atlas root worker, etc.

	return &lib.TestRunState***REMOVED***
		TestPreInitState: lct.preInitState,
		Runner:           lct.initRunner,
		Options:          lct.derivedConfig.Options, // we will always run with the derived options
	***REMOVED***, nil
***REMOVED***

func applyDefault(conf Config) Config ***REMOVED***
	if conf.SystemTags == nil ***REMOVED***
		conf.SystemTags = &metrics.DefaultSystemTagSet
	***REMOVED***
	if conf.SummaryTrendStats == nil ***REMOVED***
		conf.SummaryTrendStats = lib.DefaultSummaryTrendStats
	***REMOVED***
	defDNS := types.DefaultDNSConfig()
	if !conf.DNS.TTL.Valid ***REMOVED***
		conf.DNS.TTL = defDNS.TTL
	***REMOVED***
	if !conf.DNS.Select.Valid ***REMOVED***
		conf.DNS.Select = defDNS.Select
	***REMOVED***
	if !conf.DNS.Policy.Valid ***REMOVED***
		conf.DNS.Policy = defDNS.Policy
	***REMOVED***
	if !conf.SetupTimeout.Valid ***REMOVED***
		conf.SetupTimeout.Duration = types.Duration(60 * time.Second)
	***REMOVED***
	if !conf.TeardownTimeout.Valid ***REMOVED***
		conf.TeardownTimeout.Duration = types.Duration(60 * time.Second)
	***REMOVED***
	return conf
***REMOVED***
