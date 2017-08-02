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
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
	logtest "github.com/sirupsen/logrus/hooks/test"
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
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
	***REMOVED***, afero.NewMemMapFs())
	if assert.NoError(t, err) ***REMOVED***
		assert.NotNil(t, r1.GetDefaultGroup())
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive())
	if assert.NoError(t, err) ***REMOVED***
		assert.NotNil(t, r2.GetDefaultGroup())
	***REMOVED***
***REMOVED***

func TestRunnerOptions(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, r.Bundle.Options, r.GetOptions())
			assert.Equal(t, null.NewBool(false, false), r.Bundle.Options.Paused)
			r.ApplyOptions(lib.Options***REMOVED***Paused: null.BoolFrom(true)***REMOVED***)
			assert.Equal(t, r.Bundle.Options, r.GetOptions())
			assert.Equal(t, null.NewBool(true, true), r.Bundle.Options.Paused)
			r.ApplyOptions(lib.Options***REMOVED***Paused: null.BoolFrom(false)***REMOVED***)
			assert.Equal(t, r.Bundle.Options, r.GetOptions())
			assert.Equal(t, null.NewBool(false, true), r.Bundle.Options.Paused)
		***REMOVED***)
	***REMOVED***
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
				t.Run("Source", func(t *testing.T) ***REMOVED***
					_, err := New(&lib.SourceData***REMOVED***
						Filename: "/script.js",
						Data:     []byte(fmt.Sprintf(`import "%s"; export default function() ***REMOVED******REMOVED***`, mod)),
					***REMOVED***, afero.NewMemMapFs())
					assert.NoError(t, err)
				***REMOVED***)
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
				r1, err := New(&lib.SourceData***REMOVED***
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

				r2, err := NewFromArchive(r1.MakeArchive())
				if !assert.NoError(t, err) ***REMOVED***
					return
				***REMOVED***

				testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
				for name, r := range testdata ***REMOVED***
					t.Run(name, func(t *testing.T) ***REMOVED***
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
	***REMOVED***)
***REMOVED***

func TestVURunContext(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
		export let options = ***REMOVED*** vus: 10 ***REMOVED***;
		export default function() ***REMOVED*** fn(); ***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.ApplyOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.newVU()
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			fnCalled := false
			vu.Runtime.Set("fn", func() ***REMOVED***
				fnCalled = true

				assert.Equal(t, vu.Runtime, common.GetRuntime(*vu.Context), "incorrect runtime in context")

				state := common.GetState(*vu.Context)
				if assert.NotNil(t, state) ***REMOVED***
					assert.Equal(t, null.IntFrom(10), state.Options.VUs)
					assert.Equal(t, null.BoolFrom(true), state.Options.Throw)
					assert.NotNil(t, state.Logger)
					assert.Equal(t, r.GetDefaultGroup(), state.Group)
					assert.Equal(t, vu.HTTPTransport, state.HTTPTransport)
				***REMOVED***
			***REMOVED***)
			_, err = vu.RunOnce(context.Background())
			assert.NoError(t, err)
			assert.True(t, fnCalled, "fn() not called")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationGroups(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
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

	r2, err := NewFromArchive(r1.MakeArchive())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
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
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationMetrics(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
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

	r2, err := NewFromArchive(r1.MakeArchive())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.newVU()
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			samples, err := vu.RunOnce(context.Background())
			assert.NoError(t, err)
			assert.Len(t, samples, 3)
			for i, s := range samples ***REMOVED***
				switch i ***REMOVED***
				case 0:
					assert.Equal(t, 5.0, s.Value)
					assert.Equal(t, "my_metric", s.Metric.Name)
					assert.Equal(t, stats.Trend, s.Metric.Type)
				case 1:
					assert.Equal(t, 0.0, s.Value)
					assert.Equal(t, metrics.DataSent, s.Metric)
				case 2:
					assert.Equal(t, 0.0, s.Value)
					assert.Equal(t, metrics.DataReceived, s.Metric)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationInsecureRequests(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		opts   lib.Options
		errMsg string
	***REMOVED******REMOVED***
		"Null": ***REMOVED***
			lib.Options***REMOVED******REMOVED***,
			"GoError: Get https://expired.badssl.com/: x509: certificate has expired or is not yet valid",
		***REMOVED***,
		"False": ***REMOVED***
			lib.Options***REMOVED***InsecureSkipTLSVerify: null.BoolFrom(false)***REMOVED***,
			"GoError: Get https://expired.badssl.com/: x509: certificate has expired or is not yet valid",
		***REMOVED***,
		"True": ***REMOVED***
			lib.Options***REMOVED***InsecureSkipTLSVerify: null.BoolFrom(true)***REMOVED***,
			"",
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			r1, err := New(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					import http from "k6/http";
					export default function() ***REMOVED*** http.get("https://expired.badssl.com/"); ***REMOVED***
				`),
			***REMOVED***, afero.NewMemMapFs())
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			r1.ApplyOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***)
			r1.ApplyOptions(data.opts)

			r2, err := NewFromArchive(r1.MakeArchive())
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				t.Run(name, func(t *testing.T) ***REMOVED***
					r.Logger, _ = logtest.NewNullLogger()

					vu, err := r.NewVU()
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***
					_, err = vu.RunOnce(context.Background())
					if data.errMsg != "" ***REMOVED***
						assert.EqualError(t, err, data.errMsg)
					***REMOVED*** else ***REMOVED***
						assert.NoError(t, err)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
func TestVUIntegrationTLSConfig(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		opts   lib.Options
		errMsg string
	***REMOVED******REMOVED***
		"NullCipherSuites": ***REMOVED***
			lib.Options***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		"SupportedCipherSuite": ***REMOVED***
			lib.Options***REMOVED***TLSCipherSuites: &lib.TLSCipherSuites***REMOVED***tls.TLS_RSA_WITH_AES_128_GCM_SHA256***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		"UnsupportedCipherSuite": ***REMOVED***
			lib.Options***REMOVED***TLSCipherSuites: &lib.TLSCipherSuites***REMOVED***tls.TLS_RSA_WITH_RC4_128_SHA***REMOVED******REMOVED***,
			"GoError: Get https://sha256.badssl.com/: remote error: tls: handshake failure",
		***REMOVED***,
		"NullVersion": ***REMOVED***
			lib.Options***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		"SupportedVersion": ***REMOVED***
			lib.Options***REMOVED***TLSVersion: &lib.TLSVersion***REMOVED***Min: tls.VersionTLS12, Max: tls.VersionTLS12***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		"UnsupportedVersion": ***REMOVED***
			lib.Options***REMOVED***TLSVersion: &lib.TLSVersion***REMOVED***Min: tls.VersionSSL30, Max: tls.VersionSSL30***REMOVED******REMOVED***,
			"GoError: Get https://sha256.badssl.com/: remote error: tls: handshake failure",
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			r1, err := New(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					import http from "k6/http";
					export default function() ***REMOVED*** http.get("https://sha256.badssl.com/"); ***REMOVED***
				`),
			***REMOVED***, afero.NewMemMapFs())
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			r1.ApplyOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***)
			r1.ApplyOptions(data.opts)

			r2, err := NewFromArchive(r1.MakeArchive())
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				t.Run(name, func(t *testing.T) ***REMOVED***
					r.Logger, _ = logtest.NewNullLogger()

					vu, err := r.NewVU()
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***
					_, err = vu.RunOnce(context.Background())
					if data.errMsg != "" ***REMOVED***
						assert.EqualError(t, err, data.errMsg)
					***REMOVED*** else ***REMOVED***
						assert.NoError(t, err)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationHTTP2(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			import http from "k6/http";
			export default function() ***REMOVED***
				let res = http.request("GET", "https://http2.akamai.com/demo");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
				if (res.proto != "HTTP/2.0") ***REMOVED*** throw new Error("wrong proto: " + res.proto) ***REMOVED***
			***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.ApplyOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.NewVU()
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			samples, err := vu.RunOnce(context.Background())
			assert.NoError(t, err)

			protoFound := false
			for _, sample := range samples ***REMOVED***
				if proto := sample.Tags["proto"]; proto != "" ***REMOVED***
					protoFound = true
					assert.Equal(t, "HTTP/2.0", proto)
				***REMOVED***
			***REMOVED***
			assert.True(t, protoFound)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationCookies(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			import http from "k6/http";
			export default function() ***REMOVED***
				let preRes = http.get("https://httpbin.org/cookies");
				if (preRes.status != 200) ***REMOVED*** throw new Error("wrong status (pre): " + preRes.status); ***REMOVED***
				if (preRes.json().cookies.k1 || preRes.json().cookies.k2) ***REMOVED***
					throw new Error("cookies persisted: " + preRes.body);
				***REMOVED***

				let res = http.get("https://httpbin.org/cookies/set?k2=v2&k1=v1");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
				if (res.json().cookies.k1 != "v1" || res.json().cookies.k2 != "v2") ***REMOVED***
					throw new Error("wrong cookies: " + res.body);
				***REMOVED***
			***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.ApplyOptions(lib.Options***REMOVED***
		Throw:        null.BoolFrom(true),
		MaxRedirects: null.IntFrom(10),
	***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.NewVU()
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			for i := 0; i < 2; i++ ***REMOVED***
				_, err = vu.RunOnce(context.Background())
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationVUID(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			export default function() ***REMOVED***
				if (__VU != 1234) ***REMOVED*** throw new Error("wrong __VU: " + __VU); ***REMOVED***
			***REMOVED***`,
		),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.ApplyOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.NewVU()
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			assert.NoError(t, vu.Reconfigure(1234))
			_, err = vu.RunOnce(context.Background())
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***
