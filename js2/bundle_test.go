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

package js2

import (
	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"testing"
	"time"
)

func TestNewBundle(t *testing.T) ***REMOVED***
	t.Run("Blank", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(``),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultUndefined", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default undefined;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultNull", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default null;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultWrongType", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default 12345;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "default export must be a function")
	***REMOVED***)
	t.Run("Minimal", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default function() ***REMOVED******REMOVED***;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.NoError(t, err)
	***REMOVED***)
	t.Run("Options", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
					export let options = ***REMOVED******REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
		***REMOVED***, afero.NewMemMapFs())
		assert.NoError(t, err)

		t.Run("Paused", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						paused: true,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.Paused)
			***REMOVED***
		***REMOVED***)
		t.Run("VUs", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						vus: 100,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.VUs)
			***REMOVED***
		***REMOVED***)
		t.Run("VUsMax", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						vusMax: 100,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.VUsMax)
			***REMOVED***
		***REMOVED***)
		t.Run("Duration", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						duration: "10s",
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.StringFrom("10s"), b.Options.Duration)
			***REMOVED***
		***REMOVED***)
		t.Run("Iterations", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						iterations: 100,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.Iterations)
			***REMOVED***
		***REMOVED***)
		t.Run("Stages", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						stages: [],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Len(t, b.Options.Stages, 0)
			***REMOVED***

			t.Run("Empty", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED******REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED******REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("Target", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED***target: 10***REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("Duration", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED***duration: "10s"***REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: 10 * time.Second***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("DurationAndTarget", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED***duration: "10s", target: 10***REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: 10 * time.Second, Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("RampUpAndPlateau", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED***duration: "10s", target: 10***REMOVED***,
								***REMOVED***duration: "5s"***REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 2) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: 10 * time.Second, Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
						assert.Equal(t, lib.Stage***REMOVED***Duration: 5 * time.Second***REMOVED***, b.Options.Stages[1])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
		t.Run("Linger", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						linger: true,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.Linger)
			***REMOVED***
		***REMOVED***)
		t.Run("NoUsageReport", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						noUsageReport: true,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.NoUsageReport)
			***REMOVED***
		***REMOVED***)
		t.Run("MaxRedirects", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						maxRedirects: 10,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(10), b.Options.MaxRedirects)
			***REMOVED***
		***REMOVED***)
		t.Run("InsecureSkipTLSVerify", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						insecureSkipTLSVerify: true,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.InsecureSkipTLSVerify)
			***REMOVED***
		***REMOVED***)
		t.Run("Thresholds", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						thresholds: ***REMOVED***
							http_req_duration: ["avg<100"],
						***REMOVED***,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				if assert.Len(t, b.Options.Thresholds["http_req_duration"].Thresholds, 1) ***REMOVED***
					assert.Equal(t, "avg<100", b.Options.Thresholds["http_req_duration"].Thresholds[0].Source)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
***REMOVED***
