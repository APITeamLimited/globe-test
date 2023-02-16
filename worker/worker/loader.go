package worker

import (
	"errors"
	"fmt"
	"strings"

	"github.com/APITeamLimited/globe-test/errext"
	"github.com/APITeamLimited/globe-test/errext/exitcodes"
	"github.com/APITeamLimited/globe-test/js"
	proxy_registry "github.com/APITeamLimited/globe-test/metrics/proxy"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/spf13/afero"
)

func loadAndConfigureTest(
	gs *globalState,
	job *libOrch.ChildJob,
	workerInfo *libWorker.WorkerInfo,
) (*workerLoadedAndConfiguredTest, error) {
	sourceData, err := loader.LoadTestData(job.TestData)
	if err != nil {
		return nil, err
	}

	filesystems := make(map[string]afero.Fs, 1)
	filesystems["file"] = afero.NewMemMapFs()

	registry := proxy_registry.NewProxyRegistry(job.Options.MetricSamplesBufferSize, gs)

	preInitState := &libWorker.TestPreInitState{
		// These gs will need to be changed as on the cloud
		Logger:         gs.logger,
		BuiltinMetrics: proxy_registry.RegisterBuiltinMetrics(registry),
	}

	test := &workerLoadedTest{
		fs:            gs.fs,
		pwd:           "",
		fileSystems:   filesystems,
		preInitState:  preInitState,
		sourceData:    sourceData,
		proxyRegistry: registry,
	}

	if err := test.initializeFirstRunner(gs, workerInfo, job); err != nil {
		return nil, fmt.Errorf("could not initialize first runner: %w", err)
	}

	return test.consolidateDeriveAndValidateConfig(gs, job)
}

func (lt *workerLoadedTest) initializeFirstRunner(gs *globalState, workerInfo *libWorker.WorkerInfo, job *libOrch.ChildJob) error {
	runner, err := js.New(lt.preInitState, lt.sourceData, lt.fileSystems, workerInfo, job.TestData, lt.proxyRegistry)
	// TODO: should we use common.UnwrapGojaInterruptedError() here?
	if err != nil {
		return fmt.Errorf("could not load JS test: %w", err)
	}
	lt.initRunner = runner
	return nil

}

func (lt *workerLoadedTest) consolidateDeriveAndValidateConfig(
	gs *globalState, job *libOrch.ChildJob,
) (*workerLoadedAndConfiguredTest, error) {
	// ChildOptions have already been determined by the orchestrator, just load them
	consolidatedConfig := Config{
		Options: job.ChildOptions,
	}

	derivedConfig, err := deriveAndValidateConfig(consolidatedConfig, lt.initRunner.IsExecutable)
	if err != nil {
		return nil, err
	}

	return &workerLoadedAndConfiguredTest{
		workerLoadedTest: lt,
		derivedConfig:    derivedConfig,
	}, nil
}

func deriveAndValidateConfig(
	conf Config, isExecutable func(string) bool,
) (result Config, err error) {
	// Don't modify this, need to write to the original config
	result = conf
	err = validateConfig(result, isExecutable)
	return result, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
}

func validateConfig(conf Config, isExecutable func(string) bool) error {
	errList := conf.Validate()

	for _, ec := range conf.Scenarios {
		if err := validateScenarioConfig(ec, isExecutable); err != nil {
			errList = append(errList, err)
		}
	}

	return consolidateErrorMessage(errList, "There were problems with the specified script configuration:")
}

func validateScenarioConfig(conf libWorker.ExecutorConfig, isExecutable func(string) bool) error {
	execFn := conf.GetExec()
	if !isExecutable(execFn) {
		return fmt.Errorf("executor %s: function '%s' not found in exports", conf.GetName(), execFn)
	}
	return nil
}

func consolidateErrorMessage(errList []error, title string) error {
	if len(errList) == 0 {
		return nil
	}

	errMsgParts := []string{title}
	for _, err := range errList {
		errMsgParts = append(errMsgParts, fmt.Sprintf("\t- %s", err.Error()))
	}

	return errors.New(strings.Join(errMsgParts, "\n"))
}
