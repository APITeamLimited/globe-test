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
	"testing"

	"google.golang.org/grpc/reflection"

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

func assertMetricEmitted(t *testing.T, metricName string, sampleContainers []stats.SampleContainer, url string) ***REMOVED***
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

// codeBlock represents an execution of a k6 script.
type codeBlock struct ***REMOVED***
	code         string
	val          interface***REMOVED******REMOVED***
	err          string
	windowsError string
	asserts      func(*testing.T, *httpmultibin.HTTPMultiBin, chan stats.SampleContainer, error)
***REMOVED***

type testcase struct ***REMOVED***
	name       string
	setup      func(*httpmultibin.HTTPMultiBin)
	initString codeBlock // runs in the init context
	vuString   codeBlock // runs in the vu context
***REMOVED***

func TestClient(t *testing.T) ***REMOVED***
	t.Parallel()
	tests := []testcase***REMOVED***
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
				windowsError: "The system cannot find the file specified",
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
			vuString: codeBlock***REMOVED***code: `
				client.connect("GRPCBIN_ADDR");
				var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***)
				if (resp.status !== grpc.StatusOK) ***REMOVED***
					throw new Error("unexpected error status: " + resp.status)
				***REMOVED***`***REMOVED***,
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
				var resp = client.invoke("grpc.testing.TestService/EmptyCall", ***REMOVED******REMOVED***, ***REMOVED*** headers: ***REMOVED*** "X-Load-Tester": "k6" ***REMOVED*** ***REMOVED***)
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
				err:  "error invoking reflect API: rpc error: code = Unimplemented desc = unknown service grpc.reflection.v1alpha.ServerReflection",
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
				err:  `invalid value for 'reflect': '"true"', it needs to be boolean`,
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
	for _, test := range tests ***REMOVED***
		test := test
		t.Run(test.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			rt, state, initEnv, tb, samples, err := setup(t)
			require.NoError(t, err)
			ctx := common.WithInitEnv(common.WithRuntime(context.Background(), rt), initEnv)
			require.NoError(t, rt.Set("grpc", common.Bind(rt, New(), &ctx)))
			if test.setup != nil ***REMOVED***
				test.setup(tb)
			***REMOVED***
			val, err := rt.RunString(tb.Replacer.Replace(test.initString.code))
			assertResponse(t, test.initString, err, val, tb, samples)
			ctx = lib.WithState(ctx, state)
			val, err = rt.RunString(tb.Replacer.Replace(test.vuString.code))
			assertResponse(t, test.vuString, err, val, tb, samples)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func assertResponse(t *testing.T, r codeBlock, err error, val goja.Value, tb *httpmultibin.HTTPMultiBin, samples chan stats.SampleContainer) ***REMOVED***
	if r.err == "" && r.windowsError == "" ***REMOVED***
		assert.NoError(t, err)
	***REMOVED*** else if r.err != "" ***REMOVED***
		require.Error(t, err)
		assert.Contains(t, err.Error(), r.err)
	***REMOVED*** else if r.windowsError != "" ***REMOVED***
		require.Error(t, err)
		assert.Contains(t, err.Error(), r.windowsError)
	***REMOVED***
	if r.val != nil ***REMOVED***
		require.NotNil(t, val)
		assert.Equal(t, r.val, val.Export())
	***REMOVED***
	if r.asserts != nil ***REMOVED***
		r.asserts(t, tb, samples, err)
	***REMOVED***
***REMOVED***

func setup(t *testing.T) (*goja.Runtime, *lib.State, *common.InitEnvironment, *httpmultibin.HTTPMultiBin, chan stats.SampleContainer, error) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)
	runtime := goja.New()
	runtime.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
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
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(metrics.NewRegistry()),
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
	return runtime, state, initEnv, tb, samples, err
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
