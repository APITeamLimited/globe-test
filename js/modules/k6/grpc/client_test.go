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
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"

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
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/fsext"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/testutils/httpmultibin"
	"go.k6.io/k6/stats"
)

const isWindows = runtime.GOOS == "windows"

func assertMetricEmitted(t *testing.T, metric *stats.Metric, sampleContainers []stats.SampleContainer, url string) ***REMOVED***
	seenMetric := false

	for _, sampleContainer := range sampleContainers ***REMOVED***
		for _, sample := range sampleContainer.GetSamples() ***REMOVED***
			surl, ok := sample.Tags.Get("url")
			assert.True(t, ok)
			if surl == url ***REMOVED***
				if sample.Metric == metric ***REMOVED***
					seenMetric = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	assert.True(t, seenMetric, "url %s didn't emit %s", url, metric.Name)
***REMOVED***

func TestClient(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()
	sr := tb.Replacer.Replace

	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
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

	ctx := common.WithRuntime(context.Background(), rt)
	ctx = common.WithInitEnv(ctx, initEnv)

	rt.Set("grpc", common.Bind(rt, New(), &ctx))

	t.Run("New", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			var client = new grpc.Client();
			if (!client) throw new Error("no client created")
		`)
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("LoadNotFound", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.load([], "./does_not_exist.proto");
		`)
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***

		// (rogchap) this is a bit of a hack as windows reports a different system error than unix
		errStr := strings.Replace(err.Error(), "The system cannot find the file specified", "no such file or directory", 1)

		assert.Contains(t, errStr, "no such file or directory")
	***REMOVED***)

	t.Run("Load", func(t *testing.T) ***REMOVED***
		respV, err := rt.RunString(`
			client.load([], "../../../../vendor/google.golang.org/grpc/test/grpc_testing/test.proto");
		`)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***
		resp := respV.Export()
		assert.IsType(t, []MethodInfo***REMOVED******REMOVED***, resp)
		assert.Len(t, resp, 6)
	***REMOVED***)

	t.Run("ConnectInit", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.connect();
		`)
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "connecting to a gRPC server in the init context is not supported")
	***REMOVED***)

	t.Run("invokeInit", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			var err = client.invoke();
			throw new Error(err)
		`)
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "invoking RPC methods in the init context is not supported")
	***REMOVED***)

	ctx = lib.WithState(ctx, state)

	t.Run("NoConnect", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
		`)
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "no gRPC connection, you must call connect first")
	***REMOVED***)

	t.Run("UnknownConnectParam", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
			client.connect("GRPCBIN_ADDR", ***REMOVED*** name: "k6" ***REMOVED***);
		`))
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "unknown connect param: \"name\"")
	***REMOVED***)

	t.Run("ConnectInvalidTimeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
			client.connect("GRPCBIN_ADDR", ***REMOVED*** timeout: "k6" ***REMOVED***);
		`))
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "invalid duration")
	***REMOVED***)

	t.Run("ConnectStringTimeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
			client.connect("GRPCBIN_ADDR", ***REMOVED*** timeout: "1h3s" ***REMOVED***);
		`))
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("ConnectFloatTimeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
			client.connect("GRPCBIN_ADDR", ***REMOVED*** timeout: 3456.3 ***REMOVED***);
		`))
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("ConnectIntegerTimeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
			client.connect("GRPCBIN_ADDR", ***REMOVED*** timeout: 3000 ***REMOVED***);
		`))
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("Connect", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(sr(`
			client.connect("GRPCBIN_ADDR");
		`))
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("InvokeNotFound", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.invoke("foo/bar", ***REMOVED******REMOVED***)
		`)
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "method \"/foo/bar\" not found in file descriptors")
	***REMOVED***)

	t.Run("InvokeInvalidParam", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** void: true ***REMOVED***)
		`)
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "unknown param: \"void\"")
	***REMOVED***)

	t.Run("InvokeInvalidTimeoutType", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: true ***REMOVED***)
		`)
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "invalid timeout value: unable to use type bool as a duration value")
	***REMOVED***)

	t.Run("InvokeInvalidTimeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: "please" ***REMOVED***)
		`)
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "invalid duration")
	***REMOVED***)

	t.Run("InvokeStringTimeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: "1h42m" ***REMOVED***)
		`)
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("InvokeFloatTimeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: 400.50 ***REMOVED***)
		`)
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("InvokeIntegerTimeout", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** timeout: 2000 ***REMOVED***)
		`)
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("Invoke", func(t *testing.T) ***REMOVED***
		tb.GRPCStub.EmptyCallFunc = func(context.Context, *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
			return &grpc_testing.Empty***REMOVED******REMOVED***, nil
		***REMOVED***
		_, err := rt.RunString(`
			var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
			if (resp.status !== grpc.StatusOK) ***REMOVED***
				throw new Error("unexpected error status: " + resp.status)
			***REMOVED***
		`)
		assert.NoError(t, err)
		samplesBuf := stats.GetBufferedSamples(samples)
		assertMetricEmitted(t, metrics.GRPCReqDuration, samplesBuf, sr("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
	***REMOVED***)

	t.Run("RequestMessage", func(t *testing.T) ***REMOVED***
		tb.GRPCStub.UnaryCallFunc = func(_ context.Context, req *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error) ***REMOVED***
			if req.Payload == nil || string(req.Payload.Body) != "负载测试" ***REMOVED***
				return nil, status.Error(codes.InvalidArgument, "")
			***REMOVED***

			return &grpc_testing.SimpleResponse***REMOVED******REMOVED***, nil
		***REMOVED***
		_, err := rt.RunString(`
			var resp = client.invoke("grpc.testing.TestService/UnaryCall", ***REMOVED*** payload: ***REMOVED*** body: "6LSf6L295rWL6K+V"***REMOVED*** ***REMOVED***)
			if (resp.status !== grpc.StatusOK) ***REMOVED***
				throw new Error("server did not receive the correct request message")
			***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("RequestHeaders", func(t *testing.T) ***REMOVED***
		tb.GRPCStub.EmptyCallFunc = func(ctx context.Context, _ *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok || len(md["x-load-tester"]) == 0 || md["x-load-tester"][0] != "k6" ***REMOVED***
				return nil, status.Error(codes.FailedPrecondition, "")
			***REMOVED***

			return &grpc_testing.Empty***REMOVED******REMOVED***, nil
		***REMOVED***
		_, err := rt.RunString(`
			var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** headers: ***REMOVED*** "X-Load-Tester": "k6" ***REMOVED*** ***REMOVED***)
			if (resp.status !== grpc.StatusOK) ***REMOVED***
				throw new Error("failed to send correct headers in the request")
			***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("ResponseMessage", func(t *testing.T) ***REMOVED***
		tb.GRPCStub.UnaryCallFunc = func(context.Context, *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error) ***REMOVED***
			return &grpc_testing.SimpleResponse***REMOVED***
				OauthScope: "水",
			***REMOVED***, nil
		***REMOVED***
		_, err := rt.RunString(`
			var resp = client.invoke("grpc.testing.TestService/UnaryCall", ***REMOVED******REMOVED***)
			if (!resp.message || resp.message.username !== "" || resp.message.oauthScope !== "水") ***REMOVED***
				throw new Error("unexpected response message: " + JSON.stringify(resp.message))
			***REMOVED***
		`)
		assert.NoError(t, err)
		samplesBuf := stats.GetBufferedSamples(samples)
		assertMetricEmitted(t, metrics.GRPCReqDuration, samplesBuf, sr("GRPCBIN_ADDR/grpc.testing.TestService/UnaryCall"))
	***REMOVED***)

	t.Run("ResponseError", func(t *testing.T) ***REMOVED***
		tb.GRPCStub.EmptyCallFunc = func(context.Context, *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
			return nil, status.Error(codes.DataLoss, "foobar")
		***REMOVED***
		_, err := rt.RunString(`
			var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
			if (resp.status !== grpc.StatusDataLoss) ***REMOVED***
				throw new Error("unexpected error status: " + resp.status)
			***REMOVED***
			if (!resp.error || resp.error.message !== "foobar" || resp.error.code !== 15) ***REMOVED***
				throw new Error("unexpected error object: " + JSON.stringify(resp.error.code))
			***REMOVED***
		`)
		assert.NoError(t, err)
		samplesBuf := stats.GetBufferedSamples(samples)
		assertMetricEmitted(t, metrics.GRPCReqDuration, samplesBuf, sr("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
	***REMOVED***)

	t.Run("ResponseHeaders", func(t *testing.T) ***REMOVED***
		tb.GRPCStub.EmptyCallFunc = func(ctx context.Context, _ *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
			md := metadata.Pairs("foo", "bar")
			_ = grpc.SetHeader(ctx, md)

			return &grpc_testing.Empty***REMOVED******REMOVED***, nil
		***REMOVED***
		_, err := rt.RunString(`
			var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
			if (resp.status !== grpc.StatusOK) ***REMOVED***
				throw new Error("unexpected error status: " + resp.status)
			***REMOVED***
			if (!resp.headers || !resp.headers["foo"] || resp.headers["foo"][0] !== "bar") ***REMOVED***
				throw new Error("unexpected headers object: " + JSON.stringify(resp.trailers))
			***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("ResponseTrailers", func(t *testing.T) ***REMOVED***
		tb.GRPCStub.EmptyCallFunc = func(ctx context.Context, _ *grpc_testing.Empty) (*grpc_testing.Empty, error) ***REMOVED***
			md := metadata.Pairs("foo", "bar")
			_ = grpc.SetTrailer(ctx, md)

			return &grpc_testing.Empty***REMOVED******REMOVED***, nil
		***REMOVED***
		_, err := rt.RunString(`
			var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
			if (resp.status !== grpc.StatusOK) ***REMOVED***
				throw new Error("unexpected error status: " + resp.status)
			***REMOVED***
			if (!resp.trailers || !resp.trailers["foo"] || resp.trailers["foo"][0] !== "bar") ***REMOVED***
				throw new Error("unexpected trailers object: " + JSON.stringify(resp.trailers))
			***REMOVED***
		`)
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("LoadNotInit", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString("client.load()")
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "load must be called in the init context")
	***REMOVED***)

	t.Run("Close", func(t *testing.T) ***REMOVED***
		_, err := rt.RunString(`
			client.close();
			client.invoke();
		`)
		if !assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Contains(t, err.Error(), "no gRPC connection")
	***REMOVED***)
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
