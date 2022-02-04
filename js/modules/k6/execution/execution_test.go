/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package execution

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/stats"
	"gopkg.in/guregu/null.v3"
)

type execEnv struct ***REMOVED***
	Runtime *goja.Runtime
	Module  *ModuleInstance
	LogHook *testutils.SimpleLogrusHook
***REMOVED***

func setupTagsExecEnv(t *testing.T) execEnv ***REMOVED***
	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(ioutil.Discard)

	state := &lib.State***REMOVED***
		Options: lib.Options***REMOVED***
			SystemTags: stats.NewSystemTagSet(stats.TagVU),
		***REMOVED***,
		Tags: lib.NewTagMap(map[string]string***REMOVED***
			"vu": "42",
		***REMOVED***),
		Logger: testLog,
	***REMOVED***

	var (
		rt  = goja.New()
		ctx = context.Background()
	)

	m, ok := New().NewModuleInstance(
		&modulestest.VU***REMOVED***
			RuntimeField: rt,
			InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
			CtxField:     ctx,
			StateField:   state,
		***REMOVED***,
	).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("exec", m.Exports().Default))

	return execEnv***REMOVED***
		Module:  m,
		Runtime: rt,
		LogHook: logHook,
	***REMOVED***
***REMOVED***

func TestVUTags(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("Get", func(t *testing.T) ***REMOVED***
		t.Parallel()

		tenv := setupTagsExecEnv(t)
		tag, err := tenv.Runtime.RunString(`exec.vu.tags["vu"]`)
		require.NoError(t, err)
		assert.Equal(t, "42", tag.String())

		// not found
		tag, err = tenv.Runtime.RunString(`exec.vu.tags["not-existing-tag"]`)
		require.NoError(t, err)
		assert.Equal(t, "undefined", tag.String())
	***REMOVED***)

	t.Run("JSONEncoding", func(t *testing.T) ***REMOVED***
		t.Parallel()

		tenv := setupTagsExecEnv(t)
		state := tenv.Module.vu.State()
		state.Tags.Set("custom-tag", "mytag1")

		encoded, err := tenv.Runtime.RunString(`JSON.stringify(exec.vu.tags)`)
		require.NoError(t, err)
		assert.JSONEq(t, `***REMOVED***"vu":"42","custom-tag":"mytag1"***REMOVED***`, encoded.String())
	***REMOVED***)

	t.Run("Set", func(t *testing.T) ***REMOVED***
		t.Parallel()

		t.Run("SuccessAccetedTypes", func(t *testing.T) ***REMOVED***
			t.Parallel()

			// bool and numbers are implicitly converted into string

			tests := map[string]struct ***REMOVED***
				v   interface***REMOVED******REMOVED***
				exp string
			***REMOVED******REMOVED***
				"string": ***REMOVED***v: `"tag1"`, exp: "tag1"***REMOVED***,
				"bool":   ***REMOVED***v: true, exp: "true"***REMOVED***,
				"int":    ***REMOVED***v: 101, exp: "101"***REMOVED***,
				"float":  ***REMOVED***v: 3.14, exp: "3.14"***REMOVED***,
			***REMOVED***

			tenv := setupTagsExecEnv(t)

			for _, tc := range tests ***REMOVED***
				_, err := tenv.Runtime.RunString(fmt.Sprintf(`exec.vu.tags["mytag"] = %v`, tc.v))
				require.NoError(t, err)

				val, err := tenv.Runtime.RunString(`exec.vu.tags["mytag"]`)
				require.NoError(t, err)

				assert.Equal(t, tc.exp, val.String())
			***REMOVED***
		***REMOVED***)

		t.Run("SuccessOverwriteSystemTag", func(t *testing.T) ***REMOVED***
			t.Parallel()

			tenv := setupTagsExecEnv(t)

			_, err := tenv.Runtime.RunString(`exec.vu.tags["vu"] = "vu101"`)
			require.NoError(t, err)
			val, err := tenv.Runtime.RunString(`exec.vu.tags["vu"]`)
			require.NoError(t, err)
			assert.Equal(t, "vu101", val.String())
		***REMOVED***)

		t.Run("DiscardWrongTypeRaisingError", func(t *testing.T) ***REMOVED***
			t.Parallel()

			tenv := setupTagsExecEnv(t)
			state := tenv.Module.vu.State()
			state.Options.Throw = null.BoolFrom(true)
			require.NotNil(t, state)

			// array
			_, err := tenv.Runtime.RunString(`exec.vu.tags["custom-tag"] = [1, 3, 5]`)
			require.Contains(t, err.Error(), "only String, Boolean and Number")

			// object
			_, err = tenv.Runtime.RunString(`exec.vu.tags["custom-tag"] = ***REMOVED***f1: "value1", f2: 4***REMOVED***`)
			require.Contains(t, err.Error(), "only String, Boolean and Number")
		***REMOVED***)

		t.Run("DiscardWrongTypeOnlyWarning", func(t *testing.T) ***REMOVED***
			t.Parallel()

			tenv := setupTagsExecEnv(t)
			_, err := tenv.Runtime.RunString(`exec.vu.tags["custom-tag"] = [1, 3, 5]`)
			require.NoError(t, err)

			entries := tenv.LogHook.Drain()
			require.Len(t, entries, 1)
			assert.Contains(t, entries[0].Message, "discarded")
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestAbortTest(t *testing.T) ***REMOVED*** //nolint: tparallel
	t.Parallel()

	var (
		rt    = goja.New()
		state = &lib.State***REMOVED******REMOVED***
		ctx   = context.Background()
	)

	m, ok := New().NewModuleInstance(
		&modulestest.VU***REMOVED***
			RuntimeField: rt,
			InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
			CtxField:     ctx,
			StateField:   state,
		***REMOVED***,
	).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("exec", m.Exports().Default))

	prove := func(t *testing.T, script, reason string) ***REMOVED***
		_, err := rt.RunString(script)
		require.NotNil(t, err)
		var x *goja.InterruptedError
		assert.ErrorAs(t, err, &x)
		v, ok := x.Value().(*common.InterruptError)
		require.True(t, ok)
		require.Equal(t, v.Reason, reason)
	***REMOVED***

	t.Run("default reason", func(t *testing.T) ***REMOVED*** //nolint: paralleltest
		prove(t, "exec.test.abort()", common.AbortTest)
	***REMOVED***)
	t.Run("custom reason", func(t *testing.T) ***REMOVED*** //nolint: paralleltest
		prove(t, `exec.test.abort("mayday")`, fmt.Sprintf("%s: mayday", common.AbortTest))
	***REMOVED***)
***REMOVED***
