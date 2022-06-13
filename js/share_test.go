/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package js

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/metrics"
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

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.RuntimeState***REMOVED***
			Logger:         logger,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive())
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
			samples := make(chan metrics.SampleContainer, 100)
			initVU, err := r.NewVU(1, 1, samples)
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)
			entries := hook.Drain()
			assert.Len(t, entries, 0)
		***REMOVED***)
	***REMOVED***
***REMOVED***
