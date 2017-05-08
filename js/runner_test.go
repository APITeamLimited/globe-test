/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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
	"fmt"
	"testing"
	"time"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestRunnerNew(t *testing.T) ***REMOVED***
	t.Run("Valid", func(t *testing.T) ***REMOVED***
		r, err := New(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
			let counter = 0;
			export default function() ***REMOVED*** counter++; ***REMOVED***
		`),
		***REMOVED***, afero.NewMemMapFs())
		assert.NoError(t, err)

		t.Run("NewVU", func(t *testing.T) ***REMOVED***
			vu_, err := r.NewVU()
			assert.NoError(t, err)
			vu := vu_.(*VU)
			assert.Equal(t, int64(0), vu.Runtime.Get("counter").Export())

			t.Run("RunOnce", func(t *testing.T) ***REMOVED***
				_, err = vu.RunOnce(context.Background())
				assert.NoError(t, err)
				assert.Equal(t, int64(1), vu.Runtime.Get("counter").Export())
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)

	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		_, err := New(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(`blarg`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "ReferenceError: blarg is not defined at /script.js:1:14(0)")
	***REMOVED***)
***REMOVED***

func TestRunnerGetDefaultGroup(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
	***REMOVED***, afero.NewMemMapFs())
	assert.NoError(t, err)
	assert.NotNil(t, r.GetDefaultGroup())
***REMOVED***

func TestRunnerOptions(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
	***REMOVED***, afero.NewMemMapFs())
	assert.NoError(t, err)

	assert.Equal(t, r.Bundle.Options, r.GetOptions())
	assert.Equal(t, null.NewBool(false, false), r.Bundle.Options.Paused)
	r.ApplyOptions(lib.Options***REMOVED***Paused: null.BoolFrom(true)***REMOVED***)
	assert.Equal(t, r.Bundle.Options, r.GetOptions())
	assert.Equal(t, null.NewBool(true, true), r.Bundle.Options.Paused)
	r.ApplyOptions(lib.Options***REMOVED***Paused: null.BoolFrom(false)***REMOVED***)
	assert.Equal(t, r.Bundle.Options, r.GetOptions())
	assert.Equal(t, null.NewBool(false, true), r.Bundle.Options.Paused)
***REMOVED***

func TestRunnerIntegrationImports(t *testing.T) ***REMOVED***
	t.Run("Modules", func(t *testing.T) ***REMOVED***
		modules := []string***REMOVED***
			"k6",
			"k6/http",
			"k6/metrics",
			"k6/html",
		***REMOVED***
		for _, mod := range modules ***REMOVED***
			t.Run(mod, func(t *testing.T) ***REMOVED***
				_, err := New(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data:     []byte(fmt.Sprintf(`import "%s"; export default function() ***REMOVED******REMOVED***`, mod)),
				***REMOVED***, afero.NewMemMapFs())
				assert.NoError(t, err)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)

	t.Run("Files", func(t *testing.T) ***REMOVED***
		fs := afero.NewMemMapFs()
		assert.NoError(t, fs.MkdirAll("/path/to", 0755))
		assert.NoError(t, afero.WriteFile(fs, "/path/to/lib.js", []byte(`export default "hi!";`), 0644))

		testdata := map[string]struct***REMOVED*** filename, path string ***REMOVED******REMOVED***
			"Absolute":       ***REMOVED***"/path/script.js", "/path/to/lib.js"***REMOVED***,
			"Relative":       ***REMOVED***"/path/script.js", "./to/lib.js"***REMOVED***,
			"Adjacent":       ***REMOVED***"/path/to/script.js", "./lib.js"***REMOVED***,
			"STDIN-Absolute": ***REMOVED***"-", "/path/to/lib.js"***REMOVED***,
			"STDIN-Relative": ***REMOVED***"-", "./path/to/lib.js"***REMOVED***,
		***REMOVED***
		for name, data := range testdata ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				r, err := New(&lib.SourceData***REMOVED***
					Filename: data.filename,
					Data: []byte(fmt.Sprintf(`
					import hi from "%s";
					export default function() ***REMOVED***
						if (hi != "hi!") ***REMOVED*** throw new Error("incorrect value"); ***REMOVED***
					***REMOVED***`, data.path)),
				***REMOVED***, fs)
				if !assert.NoError(t, err) ***REMOVED***
					return
				***REMOVED***

				vu, err := r.NewVU()
				if !assert.NoError(t, err) ***REMOVED***
					return
				***REMOVED***
				_, err = vu.RunOnce(context.Background())
				assert.NoError(t, err)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestVURunContext(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED*** fn(); ***REMOVED***`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	vu, err := r.newVU()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	fnCalled := false
	vu.Runtime.Set("fn", func() ***REMOVED***
		fnCalled = true
		assert.Equal(t, vu.Runtime, common.GetRuntime(*vu.Context), "incorrect runtime in context")
		assert.Equal(t, r.GetDefaultGroup(), common.GetState(*vu.Context).Group, "incorrect group in context")
	***REMOVED***)
	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
	assert.True(t, fnCalled, "fn() not called")
***REMOVED***

func TestVURunSamples(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED*** fn(); ***REMOVED***`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	vu, err := r.newVU()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	metric := stats.New("my_metric", stats.Counter)
	sample := stats.Sample***REMOVED***Time: time.Now(), Metric: metric, Value: 1***REMOVED***
	vu.Runtime.Set("fn", func() ***REMOVED***
		state := common.GetState(*vu.Context)
		state.Samples = append(state.Samples, sample)
	***REMOVED***)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []stats.Sample***REMOVED***sample***REMOVED***, common.GetState(*vu.Context).Samples)
***REMOVED***

func TestVUIntegrationGroups(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
		import ***REMOVED*** group ***REMOVED*** from "k6";
		export default function() ***REMOVED***
			fnOuter();
			group("my group", function() ***REMOVED***
				fnInner();
				group("nested group", function() ***REMOVED***
					fnNested();
				***REMOVED***)
			***REMOVED***);
		***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	vu, err := r.newVU()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	fnOuterCalled := false
	fnInnerCalled := false
	fnNestedCalled := false
	vu.Runtime.Set("fnOuter", func() ***REMOVED***
		fnOuterCalled = true
		assert.Equal(t, r.GetDefaultGroup(), common.GetState(*vu.Context).Group)
	***REMOVED***)
	vu.Runtime.Set("fnInner", func() ***REMOVED***
		fnInnerCalled = true
		g := common.GetState(*vu.Context).Group
		assert.Equal(t, "my group", g.Name)
		assert.Equal(t, r.GetDefaultGroup(), g.Parent)
	***REMOVED***)
	vu.Runtime.Set("fnNested", func() ***REMOVED***
		fnNestedCalled = true
		g := common.GetState(*vu.Context).Group
		assert.Equal(t, "nested group", g.Name)
		assert.Equal(t, "my group", g.Parent.Name)
		assert.Equal(t, r.GetDefaultGroup(), g.Parent.Parent)
	***REMOVED***)
	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
	assert.True(t, fnOuterCalled, "fnOuter() not called")
	assert.True(t, fnInnerCalled, "fnInner() not called")
	assert.True(t, fnNestedCalled, "fnNested() not called")
***REMOVED***

func TestVUIntegrationMetrics(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
		import ***REMOVED*** group ***REMOVED*** from "k6";
		import ***REMOVED*** Trend ***REMOVED*** from "k6/metrics";
		let myMetric = new Trend("my_metric");
		export default function() ***REMOVED*** myMetric.add(5); ***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	vu, err := r.newVU()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	samples, err := vu.RunOnce(context.Background())
	if assert.NoError(t, err) && assert.Len(t, samples, 1) ***REMOVED***
		assert.Equal(t, 5.0, samples[0].Value)
		assert.Equal(t, "my_metric", samples[0].Metric.Name)
		assert.Equal(t, stats.Trend, samples[0].Metric.Type)
	***REMOVED***
***REMOVED***
