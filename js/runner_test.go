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
	"crypto/x509"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/testutils"
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
		***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
		assert.NoError(t, err)

		t.Run("NewVU", func(t *testing.T) ***REMOVED***
			vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
			assert.NoError(t, err)
			vuc, ok := vu.(*VU)
			assert.True(t, ok)
			assert.Equal(t, int64(0), vuc.Runtime.Get("counter").Export())

			t.Run("RunOnce", func(t *testing.T) ***REMOVED***
				err = vu.RunOnce(context.Background())
				assert.NoError(t, err)
				assert.Equal(t, int64(1), vuc.Runtime.Get("counter").Export())
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)

	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		_, err := New(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(`blarg`),
		***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
		assert.EqualError(t, err, "ReferenceError: blarg is not defined at /script.js:1:1(0)")
	***REMOVED***)
***REMOVED***

func TestRunnerGetDefaultGroup(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if assert.NoError(t, err) ***REMOVED***
		assert.NotNil(t, r1.GetDefaultGroup())
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if assert.NoError(t, err) ***REMOVED***
		assert.NotNil(t, r2.GetDefaultGroup())
	***REMOVED***
***REMOVED***

func TestRunnerOptions(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, r.Bundle.Options, r.GetOptions())
			assert.Equal(t, null.NewBool(false, false), r.Bundle.Options.Paused)
			r.SetOptions(lib.Options***REMOVED***Paused: null.BoolFrom(true)***REMOVED***)
			assert.Equal(t, r.Bundle.Options, r.GetOptions())
			assert.Equal(t, null.NewBool(true, true), r.Bundle.Options.Paused)
			r.SetOptions(lib.Options***REMOVED***Paused: null.BoolFrom(false)***REMOVED***)
			assert.Equal(t, r.Bundle.Options, r.GetOptions())
			assert.Equal(t, null.NewBool(false, true), r.Bundle.Options.Paused)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSetupTeardown(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			export let options = ***REMOVED***
				setupTimeout: "1s",
				teardownTimeout: "1s"
			***REMOVED***;

			export function setup() ***REMOVED***
				return ***REMOVED*** v: 1 ***REMOVED***;
			***REMOVED***
			export function teardown(data) ***REMOVED***
				if (data.v != 1) ***REMOVED***
					throw new Error("teardown: wrong data: " + JSON.stringify(data))
				***REMOVED***
			***REMOVED***
			export default function(data) ***REMOVED***
				if (data.v != 1) ***REMOVED***
					throw new Error("default: wrong data: " + JSON.stringify(data))
				***REMOVED***
			***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		samples := make(chan stats.SampleContainer, 100)
		t.Run(name, func(t *testing.T) ***REMOVED***
			if !assert.NoError(t, r.Setup(context.Background(), samples)) ***REMOVED***
				return
			***REMOVED***

			vu, err := r.NewVU(samples)
			if assert.NoError(t, err) ***REMOVED***
				err := vu.RunOnce(context.Background())
				assert.NoError(t, err)
			***REMOVED***

			assert.NoError(t, r.Teardown(context.Background(), samples))
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
					***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
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
				***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
				if !assert.NoError(t, err) ***REMOVED***
					return
				***REMOVED***

				r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
				if !assert.NoError(t, err) ***REMOVED***
					return
				***REMOVED***

				testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
				for name, r := range testdata ***REMOVED***
					t.Run(name, func(t *testing.T) ***REMOVED***
						vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
						if !assert.NoError(t, err) ***REMOVED***
							return
						***REMOVED***
						err = vu.RunOnce(context.Background())
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
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(r1.GetOptions().Apply(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***))

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.newVU(make(chan stats.SampleContainer, 100))
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
			err = vu.RunOnce(context.Background())
			assert.NoError(t, err)
			assert.True(t, fnCalled, "fn() not called")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVURunInterrupt(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
		export default function() ***REMOVED*** while(true) ***REMOVED******REMOVED*** ***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.newVU(make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			err = vu.RunOnce(ctx)
			assert.EqualError(t, err, "context cancelled at /script.js:1:1(1)")
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
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.newVU(make(chan stats.SampleContainer, 100))
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
			err = vu.RunOnce(context.Background())
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
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			samples := make(chan stats.SampleContainer, 100)
			vu, err := r.newVU(samples)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			err = vu.RunOnce(context.Background())
			assert.NoError(t, err)
			sampleCount := 0
			for i, sampleC := range stats.GetBufferedSamples(samples) ***REMOVED***
				for j, s := range sampleC.GetSamples() ***REMOVED***
					sampleCount++
					switch i + j ***REMOVED***
					case 0:
						assert.Equal(t, 5.0, s.Value)
						assert.Equal(t, "my_metric", s.Metric.Name)
						assert.Equal(t, stats.Trend, s.Metric.Type)
					case 1:
						assert.Equal(t, 0.0, s.Value)
						assert.Equal(t, metrics.DataSent, s.Metric, "`data_sent` sample is before `data_received` and `iteration_duration`")
					case 2:
						assert.Equal(t, 0.0, s.Value)
						assert.Equal(t, metrics.DataReceived, s.Metric, "`data_received` sample is after `data_received`")
					case 3:
						assert.Equal(t, metrics.IterationDuration, s.Metric, "`iteration-duration` sample is after `data_received`")
					***REMOVED***
				***REMOVED***
			***REMOVED***
			assert.Equal(t, sampleCount, 4)
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
			***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***.Apply(data.opts))

			r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				t.Run(name, func(t *testing.T) ***REMOVED***
					r.Logger, _ = logtest.NewNullLogger()

					vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***
					err = vu.RunOnce(context.Background())
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

func TestVUIntegrationBlacklist(t *testing.T) ***REMOVED***
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
					import http from "k6/http";
					export default function() ***REMOVED*** http.get("http://10.1.2.3/"); ***REMOVED***
				`),
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	_, cidr, err := net.ParseCIDR("10.0.0.0/8")
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(lib.Options***REMOVED***
		Throw:        null.BoolFrom(true),
		BlacklistIPs: []*net.IPNet***REMOVED***cidr***REMOVED***,
	***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			err = vu.RunOnce(context.Background())
			assert.EqualError(t, err, "GoError: Get http://10.1.2.3/: IP (10.1.2.3) is in a blacklisted range (10.0.0.0/8)")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationHosts(t *testing.T) ***REMOVED***
	tb := testutils.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(tb.Replacer.Replace(`
					import ***REMOVED*** check, fail ***REMOVED*** from "k6";
					import http from "k6/http";
					export default function() ***REMOVED***
						let res = http.get("http://test.loadimpact.com:HTTPBIN_PORT/");
						check(res, ***REMOVED***
							"is correct IP": (r) => r.remote_ip === "127.0.0.1"
						***REMOVED***) || fail("failed to override dns");
					***REMOVED***
				`)),
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	r1.SetOptions(lib.Options***REMOVED***
		Throw: null.BoolFrom(true),
		Hosts: map[string]net.IP***REMOVED***
			"test.loadimpact.com": net.ParseIP("127.0.0.1"),
		***REMOVED***,
	***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			err = vu.RunOnce(context.Background())
			if !assert.NoError(t, err) ***REMOVED***
				return
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
			lib.Options***REMOVED***TLSVersion: &lib.TLSVersions***REMOVED***Min: tls.VersionTLS12, Max: tls.VersionTLS12***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		"UnsupportedVersion": ***REMOVED***
			lib.Options***REMOVED***TLSVersion: &lib.TLSVersions***REMOVED***Min: tls.VersionSSL30, Max: tls.VersionSSL30***REMOVED******REMOVED***,
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
			***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***.Apply(data.opts))

			r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				t.Run(name, func(t *testing.T) ***REMOVED***
					r.Logger, _ = logtest.NewNullLogger()

					vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***
					err = vu.RunOnce(context.Background())
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
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(lib.Options***REMOVED***
		Throw:      null.BoolFrom(true),
		SystemTags: lib.GetTagSet("proto"),
	***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			samples := make(chan stats.SampleContainer, 100)
			vu, err := r.NewVU(samples)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			err = vu.RunOnce(context.Background())
			assert.NoError(t, err)

			protoFound := false
			for _, sampleC := range stats.GetBufferedSamples(samples) ***REMOVED***
				for _, sample := range sampleC.GetSamples() ***REMOVED***
					if proto, ok := sample.Tags.Get("proto"); ok ***REMOVED***
						protoFound = true
						assert.Equal(t, "HTTP/2.0", proto)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			assert.True(t, protoFound)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationOpenFunctionError(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			export default function() ***REMOVED*** open("/tmp/foo") ***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	assert.NoError(t, err)

	vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
	assert.NoError(t, err)
	err = vu.RunOnce(context.Background())
	assert.EqualError(t, err, "GoError: \"open\" function is only available to the init code (aka global scope), see https://docs.k6.io/docs/test-life-cycle for more information")
***REMOVED***

func TestVUIntegrationCookies(t *testing.T) ***REMOVED***
	tb := testutils.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(tb.Replacer.Replace(`
			import http from "k6/http";
			export default function() ***REMOVED***
				let url = "HTTPBIN_URL";
				let preRes = http.get(url + "/cookies");
				if (preRes.status != 200) ***REMOVED*** throw new Error("wrong status (pre): " + preRes.status); ***REMOVED***
				if (preRes.json().k1 || preRes.json().k2) ***REMOVED***
					throw new Error("cookies persisted: " + preRes.body);
				***REMOVED***

				let res = http.get(url + "/cookies/set?k2=v2&k1=v1");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
				if (res.json().k1 != "v1" || res.json().k2 != "v2") ***REMOVED***
					throw new Error("wrong cookies: " + res.body);
				***REMOVED***
			***REMOVED***
		`)),
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(lib.Options***REMOVED***
		Throw:        null.BoolFrom(true),
		MaxRedirects: null.IntFrom(10),
		Hosts:        tb.Dialer.Hosts,
	***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			for i := 0; i < 2; i++ ***REMOVED***
				err = vu.RunOnce(context.Background())
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
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			assert.NoError(t, vu.Reconfigure(1234))
			err = vu.RunOnce(context.Background())
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationClientCerts(t *testing.T) ***REMOVED***
	clientCAPool := x509.NewCertPool()
	assert.True(t, clientCAPool.AppendCertsFromPEM(
		[]byte("-----BEGIN CERTIFICATE-----\n"+
			"MIIBYzCCAQqgAwIBAgIUMYw1pqZ1XhXdFG0S2ITXhfHBsWgwCgYIKoZIzj0EAwIw\n"+
			"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE1MTYxODAwWhcNMjIwODE0MTYxODAw\n"+
			"WjAQMQ4wDAYDVQQDEwVNeSBDQTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABFWO\n"+
			"fg4dgL8cdvjoSWDQFLBJxlbQFlZfOSyUR277a4g91BD07KWX+9ny+Q8WuUODog06\n"+
			"xH1g8fc6zuaejllfzM6jQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTAD\n"+
			"AQH/MB0GA1UdDgQWBBTeoSFylGCmyqj1X4sWez1r6hkhjDAKBggqhkjOPQQDAgNH\n"+
			"ADBEAiAfuKi6u/BVXenCkgnU2sfXsYjel6rACuXEcx01yaaWuQIgXAtjrDisdlf4\n"+
			"0ZdoIoYjNhDAXUtnyRBt+V6+rIklv/8=\n"+
			"-----END CERTIFICATE-----"),
	))
	serverCert, err := tls.X509KeyPair(
		[]byte("-----BEGIN CERTIFICATE-----\n"+
			"MIIBxjCCAW2gAwIBAgIUICcYHG1bI28NZm676wHlMPxL+CEwCgYIKoZIzj0EAwIw\n"+
			"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE3MTQwNjAwWhcNMTgwODE3MTQwNjAw\n"+
			"WjAZMRcwFQYDVQQDEw4xMjcuMC4wLjE6Njk2OTBZMBMGByqGSM49AgEGCCqGSM49\n"+
			"AwEHA0IABCdD1IqowucJ5oUjGYCZZnXvgi7EMD4jD1osbOkzOFFnHSLRvdm6fcJu\n"+
			"vPUcl4g8zUs466sC0AVUNpk21XbA/QajgZswgZgwDgYDVR0PAQH/BAQDAgWgMB0G\n"+
			"A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1Ud\n"+
			"DgQWBBTeAc8HY3sgGIV+fu/lY0OKr2Ho0jAfBgNVHSMEGDAWgBTeoSFylGCmyqj1\n"+
			"X4sWez1r6hkhjDAZBgNVHREEEjAQgg4xMjcuMC4wLjE6Njk2OTAKBggqhkjOPQQD\n"+
			"AgNHADBEAiAt3gC5FGQfSJXQ5DloXAOeJDFnKIL7d6xhftgPS5O08QIgRuAyysB8\n"+
			"5JXHvvze5DMN/clHYptos9idVFc+weUZAUQ=\n"+
			"-----END CERTIFICATE-----\n"+
			"-----BEGIN CERTIFICATE-----\n"+
			"MIIBYzCCAQqgAwIBAgIUMYw1pqZ1XhXdFG0S2ITXhfHBsWgwCgYIKoZIzj0EAwIw\n"+
			"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE1MTYxODAwWhcNMjIwODE0MTYxODAw\n"+
			"WjAQMQ4wDAYDVQQDEwVNeSBDQTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABFWO\n"+
			"fg4dgL8cdvjoSWDQFLBJxlbQFlZfOSyUR277a4g91BD07KWX+9ny+Q8WuUODog06\n"+
			"xH1g8fc6zuaejllfzM6jQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTAD\n"+
			"AQH/MB0GA1UdDgQWBBTeoSFylGCmyqj1X4sWez1r6hkhjDAKBggqhkjOPQQDAgNH\n"+
			"ADBEAiAfuKi6u/BVXenCkgnU2sfXsYjel6rACuXEcx01yaaWuQIgXAtjrDisdlf4\n"+
			"0ZdoIoYjNhDAXUtnyRBt+V6+rIklv/8=\n"+
			"-----END CERTIFICATE-----"),
		[]byte("-----BEGIN EC PRIVATE KEY-----\n"+
			"MHcCAQEEIKYptA4VtQ8UOKL+d1wkhl+51aPpvO+ppY62nLF9Z1w5oAoGCCqGSM49\n"+
			"AwEHoUQDQgAEJ0PUiqjC5wnmhSMZgJlmde+CLsQwPiMPWixs6TM4UWcdItG92bp9\n"+
			"wm689RyXiDzNSzjrqwLQBVQ2mTbVdsD9Bg==\n"+
			"-----END EC PRIVATE KEY-----"),
	)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config***REMOVED***
		Certificates: []tls.Certificate***REMOVED***serverCert***REMOVED***,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAPool,
	***REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	defer func() ***REMOVED*** _ = listener.Close() ***REMOVED***()
	srv := &http.Server***REMOVED***
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
			_, _ = fmt.Fprintf(w, "ok")
		***REMOVED***),
		ErrorLog: stdlog.New(ioutil.Discard, "", 0),
	***REMOVED***
	go func() ***REMOVED*** _ = srv.Serve(listener) ***REMOVED***()

	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(fmt.Sprintf(`
			import http from "k6/http";
			export default function() ***REMOVED*** http.get("https://%s")***REMOVED***
		`, listener.Addr().String())),
	***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(lib.Options***REMOVED***
		Throw: null.BoolFrom(true),
		InsecureSkipTLSVerify: null.BoolFrom(true),
	***REMOVED***)

	t.Run("Unauthenticated", func(t *testing.T) ***REMOVED***
		r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***

		runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
		for name, r := range runners ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				r.Logger, _ = logtest.NewNullLogger()
				vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
				if assert.NoError(t, err) ***REMOVED***
					err := vu.RunOnce(context.Background())
					assert.EqualError(t, err, fmt.Sprintf("GoError: Get https://%s: remote error: tls: bad certificate", listener.Addr().String()))
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)

	r1.SetOptions(lib.Options***REMOVED***
		TLSAuth: []*lib.TLSAuth***REMOVED***
			***REMOVED***
				TLSAuthFields: lib.TLSAuthFields***REMOVED***
					Domains: []string***REMOVED***"127.0.0.1"***REMOVED***,
					Cert: "-----BEGIN CERTIFICATE-----\n" +
						"MIIBoTCCAUigAwIBAgIUd6XedDxP+rGo+kq0APqHElGZzs4wCgYIKoZIzj0EAwIw\n" +
						"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE3MTUwNjAwWhcNMTgwODE3MTUwNjAw\n" +
						"WjARMQ8wDQYDVQQDEwZjbGllbnQwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATL\n" +
						"mi/a1RVvk05FyrYmartbo/9cW+53DrQLW1twurII2q5ZfimdMX05A32uB3Ycoy/J\n" +
						"x+w7Ifyd/YRw0zEc3NHQo38wfTAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYI\n" +
						"KwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFN2SR/TD\n" +
						"yNW5DQWxZSkoXHQWsLY+MB8GA1UdIwQYMBaAFN6hIXKUYKbKqPVfixZ7PWvqGSGM\n" +
						"MAoGCCqGSM49BAMCA0cAMEQCICtETmyOmupmg4w3tw59VYJyOBqRTxg6SK+rOQmq\n" +
						"kE1VAiAUvsflDfmWBZ8EMPu46OhX6RX6MbvJ9NNvRco2G5ek1w==\n" +
						"-----END CERTIFICATE-----",
					Key: "-----BEGIN EC PRIVATE KEY-----\n" +
						"MHcCAQEEIOrnhT05alCeQEX66HgnSHah/m5LazjJHLDawYRnhUtZoAoGCCqGSM49\n" +
						"AwEHoUQDQgAEy5ov2tUVb5NORcq2Jmq7W6P/XFvudw60C1tbcLqyCNquWX4pnTF9\n" +
						"OQN9rgd2HKMvycfsOyH8nf2EcNMxHNzR0A==\n" +
						"-----END EC PRIVATE KEY-----",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)

	t.Run("Authenticated", func(t *testing.T) ***REMOVED***
		r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***

		runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
		for name, r := range runners ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				vu, err := r.NewVU(make(chan stats.SampleContainer, 100))
				if assert.NoError(t, err) ***REMOVED***
					err := vu.RunOnce(context.Background())
					assert.NoError(t, err)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***
