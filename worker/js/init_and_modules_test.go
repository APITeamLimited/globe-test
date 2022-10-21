package js_test

import (
	"context"
	"fmt"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/worker/js"
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

type CheckModule struct ***REMOVED***
	t             testing.TB
	initCtxCalled int
	vuCtxCalled   int
***REMOVED***

func (cm *CheckModule) InitCtx(ctx context.Context) ***REMOVED***
	cm.initCtxCalled++
***REMOVED***

func (cm *CheckModule) VuCtx(ctx context.Context) ***REMOVED***
	cm.vuCtxCalled++
***REMOVED***

func TestNewJSRunnerWithCustomModule(t *testing.T) ***REMOVED***
	t.Parallel()

	var uniqueModuleNumber int64
	checkModule := &CheckModule***REMOVED***t: t***REMOVED***
	moduleName := fmt.Sprintf("k6/x/check-%d", atomic.AddInt64(&uniqueModuleNumber, 1))
	modules.Register(moduleName, checkModule)

	script := fmt.Sprintf(`
		var check = require("%s");
		check.initCtx();

		module.exports.options = ***REMOVED*** vus: 1, iterations: 1 ***REMOVED***;
		module.exports.default = function() ***REMOVED***
			check.vuCtx();
		***REMOVED***;
	`, moduleName)

	logger := testutils.NewLogger(t)
	rtOptions := libWorker.RuntimeOptions***REMOVED***CompatibilityMode: null.StringFrom("base")***REMOVED***
	registry := workerMetrics.NewRegistry()
	builtinMetrics := workerMetrics.RegisterBuiltinMetrics(registry)
	runner, err := js.New(
		&libWorker.TestPreInitState***REMOVED***
			Logger:         logger,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
			RuntimeOptions: rtOptions,
		***REMOVED***,
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: "blah", Scheme: "file"***REMOVED***,
			Data: []byte(script),
		***REMOVED***,
		map[string]afero.Fs***REMOVED***"file": afero.NewMemMapFs(), "https": afero.NewMemMapFs()***REMOVED***, libWorker.GetTestWorkerInfo(),
	)
	require.NoError(t, err)
	assert.Equal(t, checkModule.initCtxCalled, 1)
	assert.Equal(t, checkModule.vuCtxCalled, 0)

	vu, err := runner.NewVU(1, 1, make(chan workerMetrics.SampleContainer, 100), libWorker.GetTestWorkerInfo())
	require.NoError(t, err)
	assert.Equal(t, checkModule.initCtxCalled, 2)
	assert.Equal(t, checkModule.vuCtxCalled, 0)

	vuCtx, vuCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer vuCancel()

	activeVU := vu.Activate(&libWorker.VUActivationParams***REMOVED***RunContext: vuCtx***REMOVED***)
	require.NoError(t, activeVU.RunOnce())
	assert.Equal(t, checkModule.initCtxCalled, 2)
	assert.Equal(t, checkModule.vuCtxCalled, 1)
	require.NoError(t, activeVU.RunOnce())
	assert.Equal(t, checkModule.initCtxCalled, 2)
	assert.Equal(t, checkModule.vuCtxCalled, 2)

	arc := runner.MakeArchive()
	assert.Equal(t, checkModule.initCtxCalled, 2) // shouldn't change, we're not executing the init context again
	assert.Equal(t, checkModule.vuCtxCalled, 2)

	runnerFromArc, err := js.NewFromArchive(
		&libWorker.TestPreInitState***REMOVED***
			Logger:         logger,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
			RuntimeOptions: rtOptions,
		***REMOVED***, arc, libWorker.GetTestWorkerInfo())
	require.NoError(t, err)
	assert.Equal(t, checkModule.initCtxCalled, 3) // changes because we need to get the exported functions
	assert.Equal(t, checkModule.vuCtxCalled, 2)
	vuFromArc, err := runnerFromArc.NewVU(2, 2, make(chan workerMetrics.SampleContainer, 100), libWorker.GetTestWorkerInfo())
	require.NoError(t, err)
	assert.Equal(t, checkModule.initCtxCalled, 4)
	assert.Equal(t, checkModule.vuCtxCalled, 2)
	activeVUFromArc := vuFromArc.Activate(&libWorker.VUActivationParams***REMOVED***RunContext: vuCtx***REMOVED***)
	require.NoError(t, activeVUFromArc.RunOnce())
	assert.Equal(t, checkModule.initCtxCalled, 4)
	assert.Equal(t, checkModule.vuCtxCalled, 3)
	require.NoError(t, activeVUFromArc.RunOnce())
	assert.Equal(t, checkModule.initCtxCalled, 4)
	assert.Equal(t, checkModule.vuCtxCalled, 4)
***REMOVED***
