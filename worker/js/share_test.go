package js

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils"
)

func TestNewSharedArrayIntegration(t *testing.T) ***REMOVED***
	t.Parallel()
	data := `'use strict';
var SharedArray = require("k6/data").SharedArray;
function generateArray() ***REMOVED***
    console.log("once");
    var n = 50;
    var arr = new Array(n);
    for (var i = 0 ; i <n; i++) ***REMOVED***
        arr[i] = ***REMOVED***value: "something" +i***REMOVED***;
    ***REMOVED***
    return arr;
***REMOVED***

var s = new SharedArray("something", generateArray);

exports.default = function() ***REMOVED***
	if (s[2].value !== "something2") ***REMOVED***
		throw new Error("bad s[2]="+s[2].value);
	***REMOVED***
	if (s.length != 50) ***REMOVED***
		throw new Error("bad length " +_s.length);
	***REMOVED***

	var i = 0;
	for (var v of s) ***REMOVED***
		if (v.value !== "something"+i) ***REMOVED***
			throw new Error("bad v.value="+v.value+" for i="+i);
		***REMOVED***
		i++;
	***REMOVED***
***REMOVED***`

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.Out = ioutil.Discard
	hook := testutils.SimpleLogrusHook***REMOVED***
		HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel***REMOVED***,
	***REMOVED***
	logger.AddHook(&hook)

	r1, err := getSimpleRunner(t, "/script.js", data, logger)
	require.NoError(t, err)
	entries := hook.Drain()
	require.Len(t, entries, 1)
	assert.Equal(t, logrus.InfoLevel, entries[0].Level)
	assert.Equal(t, "once", entries[0].Message)

	registry := workerMetrics.NewRegistry()
	builtinMetrics := workerMetrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&libWorker.TestPreInitState***REMOVED***
			Logger:         logger,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), libWorker.GetTestWorkerInfo())
	require.NoError(t, err)
	entries = hook.Drain()
	require.Len(t, entries, 1)
	assert.Equal(t, logrus.InfoLevel, entries[0].Level)
	assert.Equal(t, "once", entries[0].Message)

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			samples := make(chan workerMetrics.SampleContainer, 100)
			initVU, err := r.NewVU(1, 1, samples, libWorker.GetTestWorkerInfo())
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&libWorker.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)
			entries := hook.Drain()
			assert.Len(t, entries, 0)
		***REMOVED***)
	***REMOVED***
***REMOVED***
