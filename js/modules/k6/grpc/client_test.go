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

package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	"google.golang.org/grpc/reflection"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstats "google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/grpc_testing"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/fsext"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/httpmultibin"
	"go.k6.io/k6/stats"
)

const isWindows = runtime.GOOS == "windows"

// codeBlock represents an execution of a k6 script.
type codeBlock struct ***REMOVED***
	code       string
	val        interface***REMOVED******REMOVED***
	err        string
	windowsErr string
	asserts    func(*testing.T, *httpmultibin.HTTPMultiBin, chan stats.SampleContainer, error)
***REMOVED***

type testcase struct ***REMOVED***
	name       string
	setup      func(*httpmultibin.HTTPMultiBin)
	initString codeBlock // runs in the init context
	vuString   codeBlock // runs in the vu context
***REMOVED***

func TestClient(t *testing.T) ***REMOVED***
	t.Parallel()

	type testState struct ***REMOVED***
		rt      *goja.Runtime
		vuState *lib.State
		env     *common.InitEnvironment
		httpBin *httpmultibin.HTTPMultiBin
		samples chan stats.SampleContainer
	***REMOVED***
	setup := func(t *testing.T) testState ***REMOVED***
		t.Helper()

		root, err := lib.NewGroup("", nil)
		require.NoError(t, err)
		tb := httpmultibin.NewHTTPMultiBin(t)
		samples := make(chan stats.SampleContainer, 1000)
		state := &lib.State***REMOVED***
			Group:     root,
			Dialer:    tb.Dialer,
			TLSConfig: tb.TLSClientConfig,
			Samples:   samples,
			Options: lib.Options***REMOVED***
				SystemTags: stats.NewSystemTagSet(
					stats.TagName,
					stats.TagURL,
				),
				UserAgent: null.StringFrom("k6-test"),
			***REMOVED***,
			BuiltinMetrics: metrics.RegisterBuiltinMetrics(
				metrics.NewRegistry(),
			),
			Tags: lib.NewTagMap(nil),
		***REMOVED***

		cwd, err := os.Getwd()
		require.NoError(t, err)
		fs := afero.NewOsFs()
		if isWindows ***REMOVED***
			fs = fsext.NewTrimFilePathSeparatorFs(fs)
		***REMOVED***
		initEnv := &common.InitEnvironment***REMOVED***
			Logger: logrus.New(),
			CWD:    &url.URL***REMOVED***Path: cwd***REMOVED***,
			FileSystems: map[string]afero.Fs***REMOVED***
				"file": fs,
			***REMOVED***,
		***REMOVED***

		rt := goja.New()
		rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

		return testState***REMOVED***
			rt:      rt,
			httpBin: tb,
			vuState: state,
			env:     initEnv,
			samples: samples,
		***REMOVED***
	***REMOVED***

	assertMetricEmitted := func(
		t *testing.T,
		metricName string,
		sampleContainers []stats.SampleContainer,
		url string,
	) ***REMOVED***
		seenMetric := false

		for _, sampleContainer := range sampleContainers ***REMOVED***
			for _, sample := range sampleContainer.GetSamples() ***REMOVED***
				surl, ok := sample.Tags.Get("url")
				assert.True(t, ok)
				if surl == url ***REMOVED***
					if sample.Metric.Name == metricName ***REMOVED***
						seenMetric = true
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		assert.True(t, seenMetric, "url %s didn't emit %s", url, metricName)
	***REMOVED***

	tests := []testcase***REMOVED***
		***REMOVED***
			name: "BadTLS",
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				// changing the pointer's value
				// for affecting the lib.State
				// that uses the same pointer
				*tb.TLSClientConfig = tls.Config***REMOVED***
					MinVersion: tls.VersionTLS13,
				***REMOVED***
			***REMOVED***,
			initString: codeBlock***REMOVED***
				code: `var client = new grpc.Client();`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `client.connect("GRPCBIN_ADDR", ***REMOVED***timeout: '1s'***REMOVED***)`,
				err:  "certificate signed by unknown authority",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "New",
			initString: codeBlock***REMOVED***
				code: `
			var client = new grpc.Client();
			if (!client) throw new Error("no client created")`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "LoadNotFound",
			initString: codeBlock***REMOVED***
				code: `
			var client = new grpc.Client();
			client.load([], "./does_not_exist.proto");`,
				err: "no such file or directory",
				// (rogchap) this is a bit of a hack as windows reports a different system error than unix.
				windowsErr: "The system cannot find the file specified",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "Load",
			initString: codeBlock***REMOVED***
				code: `
			var client = new grpc.Client();
			client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
				val: []MethodInfo***REMOVED******REMOVED***MethodInfo: grpc.MethodInfo***REMOVED***Name: "EmptyCall", IsClientStream: false, IsServerStream: false***REMOVED***, Package: "grpc.testing", Service: "TestService", FullMethod: "/grpc.testing.TestService/EmptyCall"***REMOVED***, ***REMOVED***MethodInfo: grpc.MethodInfo***REMOVED***Name: "UnaryCall", IsClientStream: false, IsServerStream: false***REMOVED***, Package: "grpc.testing", Service: "TestService", FullMethod: "/grpc.testing.TestService/UnaryCall"***REMOVED***, ***REMOVED***MethodInfo: grpc.MethodInfo***REMOVED***Name: "StreamingOutputCall", IsClientStream: false, IsServerStream: true***REMOVED***, Package: "grpc.testing", Service: "TestService", FullMethod: "/grpc.testing.TestService/StreamingOutputCall"***REMOVED***, ***REMOVED***MethodInfo: grpc.MethodInfo***REMOVED***Name: "StreamingInputCall", IsClientStream: true, IsServerStream: false***REMOVED***, Package: "grpc.testing", Service: "TestService", FullMethod: "/grpc.testing.TestService/StreamingInputCall"***REMOVED***, ***REMOVED***MethodInfo: grpc.MethodInfo***REMOVED***Name: "FullDuplexCall", IsClientStream: true, IsServerStream: true***REMOVED***, Package: "grpc.testing", Service: "TestService", FullMethod: "/grpc.testing.TestService/FullDuplexCall"***REMOVED***, ***REMOVED***MethodInfo: grpc.MethodInfo***REMOVED***Name: "HalfDuplexCall", IsClientStream: true, IsServerStream: true***REMOVED***, Package: "grpc.testing", Service: "TestService", FullMethod: "/grpc.testing.TestService/HalfDuplexCall"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ConnectInit",
			initString: codeBlock***REMOVED***
				code: `
			var client = new grpc.Client();
			client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");
			client.connect();`,
				err: "connecting to a gRPC server in the init context is not supported",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "InvokeInit",
			initString: codeBlock***REMOVED***
				code: `
			var client = new grpc.Client();
			client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");
			var err = client.invoke();
			throw new Error(err)`,
				err: "invoking RPC methods in the init context is not supported",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "NoConnect",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");
				client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)`,
				err: "invoking RPC methods in the init context is not supported",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "UnknownConnectParam",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `client.connect("GRPCBIN_ADDR", ***REMOVED*** name: "k6" ***REMOVED***);`,
				err:  `unknown connect param: "name"`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ConnectInvalidTimeout",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `client.connect("GRPCBIN_ADDR", ***REMOVED*** timeout: "k6" ***REMOVED***);`,
				err:  "invalid duration",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ConnectStringTimeout",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***code: `client.connect("GRPCBIN_ADDR", ***REMOVED*** timeout: "1h3s" ***REMOVED***);`***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ConnectIntegerTimeout",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***code: `client.connect("GRPCBIN_ADDR", ***REMOVED*** timeout: 3000 ***REMOVED***);`***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ConnectFloatTimeout",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***code: `client.connect("GRPCBIN_ADDR", ***REMOVED*** timeout: 3456.3 ***REMOVED***);`***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "Connect",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***code: `client.connect("GRPCBIN_ADDR");`***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "InvokeNotFound",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				client.invoke("foo/bar", ***REMOVED******REMOVED***)`,
				err: `method "/foo/bar" not found in file descriptors`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "InvokeInvalidParam",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** void: true ***REMOVED***)`,
				err: `unknown param: "void"`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "InvokeInvalidTimeoutType",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: true ***REMOVED***)`,
				err: "invalid timeout value: unable to use type bool as a duration value",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "InvokeInvalidTimeout",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: "please" ***REMOVED***)`,
				err: "invalid duration",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "InvokeStringTimeout",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: "1h42m" ***REMOVED***)`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "InvokeFloatTimeout",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: 400.50 ***REMOVED***)`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "InvokeIntegerTimeout",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: 2000 ***REMOVED***)`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "Invoke",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`***REMOVED***,
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				tb.GRPCStub.EmptyCallFunc = func(context.Context, *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
					return &grpc_testing.Empty***REMOVED******REMOVED***, nil
				***REMOVED***
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
				if (resp.status !== grpc.StatusOK) ***REMOVED***
					throw new Error("unexpected error status: " + resp.status)
				***REMOVED***`,
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan stats.SampleContainer, _ error) ***REMOVED***
					samplesBuf := stats.GetBufferedSamples(samples)
					assertMetricEmitted(t, metrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "RequestMessage",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				tb.GRPCStub.UnaryCallFunc = func(_ context.Context, req *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error) ***REMOVED***
					if req.Payload == nil || string(req.Payload.Body) != "负载测试" ***REMOVED***
						return nil, status.Error(codes.InvalidArgument, "")
					***REMOVED***
					return &grpc_testing.SimpleResponse***REMOVED******REMOVED***, nil
				***REMOVED***
			***REMOVED***,
			vuString: codeBlock***REMOVED***code: `
				client.connect("GRPCBIN_ADDR");
				var resp = client.invoke("grpc.testing.TestService/UnaryCall", ***REMOVED*** payload: ***REMOVED*** body: "6LSf6L295rWL6K+V"***REMOVED*** ***REMOVED***)
				if (resp.status !== grpc.StatusOK) ***REMOVED***
					throw new Error("server did not receive the correct request message")
				***REMOVED***`***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "RequestHeaders",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				tb.GRPCStub.EmptyCallFunc = func(ctx context.Context, _ *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
					md, ok := metadata.FromIncomingContext(ctx)
					if !ok || len(md["x-load-tester"]) == 0 || md["x-load-tester"][0] != "k6" ***REMOVED***
						return nil, status.Error(codes.FailedPrecondition, "")
					***REMOVED***

					return &grpc_testing.Empty***REMOVED******REMOVED***, nil
				***REMOVED***
			***REMOVED***,
			vuString: codeBlock***REMOVED***code: `
				client.connect("GRPCBIN_ADDR");
				var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** metadata: ***REMOVED*** "X-Load-Tester": "k6" ***REMOVED*** ***REMOVED***)
				if (resp.status !== grpc.StatusOK) ***REMOVED***
					throw new Error("failed to send correct headers in the request")
				***REMOVED***
			`***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ResponseMessage",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				tb.GRPCStub.UnaryCallFunc = func(context.Context, *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error) ***REMOVED***
					return &grpc_testing.SimpleResponse***REMOVED***
						OauthScope: "水",
					***REMOVED***, nil
				***REMOVED***
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				var resp = client.invoke("grpc.testing.TestService/UnaryCall", ***REMOVED******REMOVED***)
				if (!resp.message || resp.message.username !== "" || resp.message.oauthScope !== "水") ***REMOVED***
					throw new Error("unexpected response message: " + JSON.stringify(resp.message))
				***REMOVED***`,
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan stats.SampleContainer, _ error) ***REMOVED***
					samplesBuf := stats.GetBufferedSamples(samples)
					assertMetricEmitted(t, metrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/UnaryCall"))
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ResponseError",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				tb.GRPCStub.EmptyCallFunc = func(context.Context, *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
					return nil, status.Error(codes.DataLoss, "foobar")
				***REMOVED***
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
				if (resp.status !== grpc.StatusDataLoss) ***REMOVED***
					throw new Error("unexpected error status: " + resp.status)
				***REMOVED***
				if (!resp.error || resp.error.message !== "foobar" || resp.error.code !== 15) ***REMOVED***
					throw new Error("unexpected error object: " + JSON.stringify(resp.error.code))
				***REMOVED***`,
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan stats.SampleContainer, _ error) ***REMOVED***
					samplesBuf := stats.GetBufferedSamples(samples)
					assertMetricEmitted(t, metrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ResponseHeaders",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				tb.GRPCStub.EmptyCallFunc = func(ctx context.Context, _ *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
					md := metadata.Pairs("foo", "bar")
					_ = grpc.SetHeader(ctx, md)
					return &grpc_testing.Empty***REMOVED******REMOVED***, nil
				***REMOVED***
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
				if (resp.status !== grpc.StatusOK) ***REMOVED***
					throw new Error("unexpected error status: " + resp.status)
				***REMOVED***
				if (!resp.headers || !resp.headers["foo"] || resp.headers["foo"][0] !== "bar") ***REMOVED***
					throw new Error("unexpected headers object: " + JSON.stringify(resp.trailers))
				***REMOVED***`,
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan stats.SampleContainer, _ error) ***REMOVED***
					samplesBuf := stats.GetBufferedSamples(samples)
					assertMetricEmitted(t, metrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ResponseTrailers",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				tb.GRPCStub.EmptyCallFunc = func(ctx context.Context, _ *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
					md := metadata.Pairs("foo", "bar")
					_ = grpc.SetTrailer(ctx, md)
					return &grpc_testing.Empty***REMOVED******REMOVED***, nil
				***REMOVED***
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
				if (resp.status !== grpc.StatusOK) ***REMOVED***
					throw new Error("unexpected error status: " + resp.status)
				***REMOVED***
				if (!resp.trailers || !resp.trailers["foo"] || resp.trailers["foo"][0] !== "bar") ***REMOVED***
					throw new Error("unexpected trailers object: " + JSON.stringify(resp.trailers))
				***REMOVED***`,
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan stats.SampleContainer, _ error) ***REMOVED***
					samplesBuf := stats.GetBufferedSamples(samples)
					assertMetricEmitted(t, metrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "LoadNotInit",
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				tb.GRPCStub.EmptyCallFunc = func(ctx context.Context, _ *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
					md := metadata.Pairs("foo", "bar")
					_ = grpc.SetTrailer(ctx, md)
					return &grpc_testing.Empty***REMOVED******REMOVED***, nil
				***REMOVED***
			***REMOVED***,
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `client.load()`,
				err:  "load must be called in the init context",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ReflectUnregistered",
			initString: codeBlock***REMOVED***
				code: `var client = new grpc.Client();`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `client.connect("GRPCBIN_ADDR", ***REMOVED***reflect: true***REMOVED***)`,
				err:  "rpc error: code = Unimplemented desc = unknown service grpc.reflection.v1alpha.ServerReflection",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "Reflect",
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				reflection.Register(tb.ServerGRPC)
			***REMOVED***,
			initString: codeBlock***REMOVED***
				code: `var client = new grpc.Client();`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `client.connect("GRPCBIN_ADDR", ***REMOVED***reflect: true***REMOVED***)`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ReflectBadParam",
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				reflection.Register(tb.ServerGRPC)
			***REMOVED***,
			initString: codeBlock***REMOVED***
				code: `var client = new grpc.Client();`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `client.connect("GRPCBIN_ADDR", ***REMOVED***reflect: "true"***REMOVED***)`,
				err:  `invalid reflect value`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ReflectInvokeNoExist",
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				reflection.Register(tb.ServerGRPC)
				tb.GRPCStub.EmptyCallFunc = func(ctx context.Context, _ *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
					return &grpc_testing.Empty***REMOVED******REMOVED***, nil
				***REMOVED***
			***REMOVED***,
			initString: codeBlock***REMOVED***
				code: `var client = new grpc.Client();`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
					client.connect("GRPCBIN_ADDR", ***REMOVED***reflect: true***REMOVED***)
					client.invoke("foo/bar", ***REMOVED******REMOVED***)
				`,
				err: `method "/foo/bar" not found in file descriptors`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ReflectInvoke",
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				reflection.Register(tb.ServerGRPC)
				tb.GRPCStub.EmptyCallFunc = func(ctx context.Context, _ *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
					return &grpc_testing.Empty***REMOVED******REMOVED***, nil
				***REMOVED***
			***REMOVED***,
			initString: codeBlock***REMOVED***
				code: `var client = new grpc.Client();`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
					client.connect("GRPCBIN_ADDR", ***REMOVED***reflect: true***REMOVED***)
					client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
				`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "Close",
			initString: codeBlock***REMOVED***
				code: `
				var client = new grpc.Client();
				client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");`,
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
			client.close();
			client.invoke();`,
				err: "no gRPC connection",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	assertResponse := func(t *testing.T, cb codeBlock, err error, val goja.Value, ts testState) ***REMOVED***
		if isWindows && cb.windowsErr != "" && err != nil ***REMOVED***
			err = errors.New(strings.ReplaceAll(err.Error(), cb.windowsErr, cb.err))
		***REMOVED***
		if cb.err == "" ***REMOVED***
			assert.NoError(t, err)
		***REMOVED*** else ***REMOVED***
			require.Error(t, err)
			assert.Contains(t, err.Error(), cb.err)
		***REMOVED***
		if cb.val != nil ***REMOVED***
			require.NotNil(t, val)
			assert.Equal(t, cb.val, val.Export())
		***REMOVED***
		if cb.asserts != nil ***REMOVED***
			cb.asserts(t, ts.httpBin, ts.samples, err)
		***REMOVED***
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		tt := tt
		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			ts := setup(t)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mvu := &modulestest.VU***REMOVED***
				RuntimeField: ts.rt,
				InitEnvField: ts.env,
				CtxField:     ctx,
			***REMOVED***

			m, ok := New().NewModuleInstance(mvu).(*ModuleInstance)
			require.True(t, ok)
			require.NoError(t, ts.rt.Set("grpc", m.Exports().Named))

			// setup necessary environment if needed by a test
			if tt.setup != nil ***REMOVED***
				tt.setup(ts.httpBin)
			***REMOVED***

			replace := func(code string) (goja.Value, error) ***REMOVED***
				return ts.rt.RunString(ts.httpBin.Replacer.Replace(code))
			***REMOVED***

			val, err := replace(tt.initString.code)
			assertResponse(t, tt.initString, err, val, ts)

			mvu.StateField = ts.vuState
			val, err = replace(tt.vuString.code)
			assertResponse(t, tt.vuString, err, val, ts)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestDebugStat(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := [...]struct ***REMOVED***
		name     string
		stat     grpcstats.RPCStats
		expected string
	***REMOVED******REMOVED***
		***REMOVED***
			"OutHeader",
			&grpcstats.OutHeader***REMOVED******REMOVED***,
			"Out Header:",
		***REMOVED***,
		***REMOVED***
			"OutTrailer",
			&grpcstats.OutTrailer***REMOVED***
				Trailer: metadata.MD***REMOVED***
					"x-trail": []string***REMOVED***"out"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			"Out Trailer:",
		***REMOVED***,
		***REMOVED***
			"OutPayload",
			&grpcstats.OutPayload***REMOVED***
				Payload: &grpc_testing.SimpleRequest***REMOVED***
					FillUsername: true,
				***REMOVED***,
			***REMOVED***,
			"fill_username:",
		***REMOVED***,
		***REMOVED***
			"InHeader",
			&grpcstats.InHeader***REMOVED***
				Header: metadata.MD***REMOVED***
					"x-head": []string***REMOVED***"in"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			"x-head: in",
		***REMOVED***,
		***REMOVED***
			"InTrailer",
			&grpcstats.InTrailer***REMOVED***
				Trailer: metadata.MD***REMOVED***
					"x-trail": []string***REMOVED***"in"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			"x-trail: in",
		***REMOVED***,
		***REMOVED***
			"InPayload",
			&grpcstats.InPayload***REMOVED***
				Payload: &grpc_testing.SimpleResponse***REMOVED***
					Username: "k6-user",
				***REMOVED***,
			***REMOVED***,
			"username:",
		***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		tt := tt
		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			var b bytes.Buffer
			logger := logrus.New()
			logger.Out = &b

			debugStat(tt.stat, logger.WithField("source", "test"), "full")
			assert.Contains(t, b.String(), tt.expected)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestClientInvokeHeadersDeprecated(t *testing.T) ***REMOVED***
	t.Parallel()

	logHook := &testutils.SimpleLogrusHook***REMOVED***
		HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED***,
	***REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(ioutil.Discard)

	c := Client***REMOVED***
		vu: &modulestest.VU***REMOVED***
			StateField: &lib.State***REMOVED***
				Logger: testLog,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	params := map[string]interface***REMOVED******REMOVED******REMOVED***
		"headers": map[string]interface***REMOVED******REMOVED******REMOVED***
			"X-HEADER-FOO": "bar",
		***REMOVED***,
	***REMOVED***
	_, err := c.parseParams(params)
	require.NoError(t, err)

	entries := logHook.Drain()
	require.Len(t, entries, 1)
	require.Contains(t, entries[0].Message, "headers property is deprecated")
***REMOVED***

func TestResolveFileDescriptors(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		name                string
		pkgs                []string
		services            []string
		expectedDescriptors int
	***REMOVED******REMOVED***
		***REMOVED***
			name:                "SuccessSamePackage",
			pkgs:                []string***REMOVED***"mypkg"***REMOVED***,
			services:            []string***REMOVED***"Service1", "Service2", "Service3"***REMOVED***,
			expectedDescriptors: 3,
		***REMOVED***,
		***REMOVED***
			name:                "SuccessMultiPackages",
			pkgs:                []string***REMOVED***"mypkg1", "mypkg2", "mypkg3"***REMOVED***,
			services:            []string***REMOVED***"Service", "Service", "Service"***REMOVED***,
			expectedDescriptors: 3,
		***REMOVED***,
		***REMOVED***
			name:                "DeduplicateServices",
			pkgs:                []string***REMOVED***"mypkg1"***REMOVED***,
			services:            []string***REMOVED***"Service1", "Service2", "Service1"***REMOVED***,
			expectedDescriptors: 2,
		***REMOVED***,
		***REMOVED***
			name:                "NoServices",
			services:            []string***REMOVED******REMOVED***,
			expectedDescriptors: 0,
		***REMOVED***,
	***REMOVED***

	for _, tt := range tests ***REMOVED***
		tt := tt
		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			var (
				lsr  = &reflectpb.ListServiceResponse***REMOVED******REMOVED***
				mock = &getServiceFileDescriptorMock***REMOVED******REMOVED***
			)
			for i, service := range tt.services ***REMOVED***
				// if only one package is defined then
				// the package is the same for every service
				pkg := tt.pkgs[0]
				if len(tt.pkgs) > 1 ***REMOVED***
					pkg = tt.pkgs[i]
				***REMOVED***

				lsr.Service = append(lsr.Service, &reflectpb.ServiceResponse***REMOVED***
					Name: fmt.Sprintf("%s.%s", pkg, service),
				***REMOVED***)
				mock.pkgs = append(mock.pkgs, pkg)
				mock.names = append(mock.names, service)
			***REMOVED***

			fdset, err := resolveServiceFileDescriptors(mock, lsr)
			require.NoError(t, err)
			assert.Len(t, fdset.File, tt.expectedDescriptors)
		***REMOVED***)
	***REMOVED***
***REMOVED***

type getServiceFileDescriptorMock struct ***REMOVED***
	nreqs int64
	pkgs  []string
	names []string
***REMOVED***

func (m *getServiceFileDescriptorMock) Send(req *reflectpb.ServerReflectionRequest) error ***REMOVED***
	// TODO: check that the sent message is expected,
	// otherwise return an error
	return nil
***REMOVED***

func (m *getServiceFileDescriptorMock) Recv() (*reflectpb.ServerReflectionResponse, error) ***REMOVED***
	n := atomic.AddInt64(&m.nreqs, 1)
	ptr := func(s string) (sptr *string) ***REMOVED***
		return &s
	***REMOVED***
	index := n - 1
	fdp := &descriptorpb.FileDescriptorProto***REMOVED***
		Package: ptr(m.pkgs[index]),
		Name:    ptr(m.names[index]),
	***REMOVED***
	b, err := proto.Marshal(fdp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	srr := &reflectpb.ServerReflectionResponse***REMOVED***
		MessageResponse: &reflectpb.ServerReflectionResponse_FileDescriptorResponse***REMOVED***
			FileDescriptorResponse: &reflectpb.FileDescriptorResponse***REMOVED***
				FileDescriptorProto: [][]byte***REMOVED***b***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	return srr, nil
***REMOVED***
