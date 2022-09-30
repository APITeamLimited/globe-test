package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"

	"google.golang.org/grpc/reflection"

	"github.com/dop251/goja"
	"github.com/golang/protobuf/ptypes/any"
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

	"github.com/APITeamLimited/globe-test/worker/js/modulestest"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/fsext"
	"github.com/APITeamLimited/globe-test/worker/libWorker/netext/grpcext"
	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils"
	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils/httpmultibin"
	grpcanytesting "github.com/APITeamLimited/globe-test/worker/libWorker/testutils/httpmultibin/grpc_any_testing"
)

const isWindows = runtime.GOOS == "windows"

// codeBlock represents an execution of a k6 script.
type codeBlock struct ***REMOVED***
	code       string
	val        interface***REMOVED******REMOVED***
	err        string
	windowsErr string
	asserts    func(*testing.T, *httpmultibin.HTTPMultiBin, chan workerMetrics.SampleContainer, error)
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
		*modulestest.Runtime
		httpBin *httpmultibin.HTTPMultiBin
		samples chan workerMetrics.SampleContainer
	***REMOVED***
	setup := func(t *testing.T) testState ***REMOVED***
		t.Helper()

		tb := httpmultibin.NewHTTPMultiBin(t)
		samples := make(chan workerMetrics.SampleContainer, 1000)
		testRuntime := modulestest.NewRuntime(t)

		cwd, err := os.Getwd()
		require.NoError(t, err)
		fs := afero.NewOsFs()
		if isWindows ***REMOVED***
			fs = fsext.NewTrimFilePathSeparatorFs(fs)
		***REMOVED***
		testRuntime.VU.InitEnvField.CWD = &url.URL***REMOVED***Path: cwd***REMOVED***
		testRuntime.VU.InitEnvField.FileSystems = map[string]afero.Fs***REMOVED***"file": fs***REMOVED***

		return testState***REMOVED***
			Runtime: testRuntime,
			httpBin: tb,
			samples: samples,
		***REMOVED***
	***REMOVED***

	assertMetricEmitted := func(
		t *testing.T,
		metricName string,
		sampleContainers []workerMetrics.SampleContainer,
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
				// for affecting the libWorker.State
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
					throw new Error("unexpected error: " + JSON.stringify(resp.error) + "or status: " + resp.status)
				***REMOVED***`,
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan workerMetrics.SampleContainer, _ error) ***REMOVED***
					samplesBuf := workerMetrics.GetBufferedSamples(samples)
					assertMetricEmitted(t, workerMetrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "InvokeAnyProto",
			initString: codeBlock***REMOVED***code: `
				var client = new grpc.Client();
				client.load([], "../../../../lib/testutils/httpmultibin/grpc_any_testing/any_test.proto");`***REMOVED***,
			setup: func(tb *httpmultibin.HTTPMultiBin) ***REMOVED***
				tb.GRPCAnyStub.SumFunc = func(ctx context.Context, req *grpcanytesting.SumRequest) (*grpcanytesting.SumReply, error) ***REMOVED***
					var sumRequestData grpcanytesting.SumRequestData
					if err := req.Data.UnmarshalTo(&sumRequestData); err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					sumReplyData := &grpcanytesting.SumReplyData***REMOVED***
						V:   sumRequestData.A + sumRequestData.B,
						Err: "",
					***REMOVED***
					sumReply := &grpcanytesting.SumReply***REMOVED***
						Data: &any.Any***REMOVED******REMOVED***,
					***REMOVED***
					if err := sumReply.Data.MarshalFrom(sumReplyData); err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					return sumReply, nil
				***REMOVED***
			***REMOVED***,
			vuString: codeBlock***REMOVED***
				code: `
				client.connect("GRPCBIN_ADDR");
				var resp = client.invoke("grpc.any.testing.AnyTestService/Sum",  ***REMOVED***
					data: ***REMOVED***
						"@type": "type.googleapis.com/grpc.any.testing.SumRequestData",
						"a": 1,
						"b": 2,
					***REMOVED***,
				***REMOVED***)
				if (resp.status !== grpc.StatusOK) ***REMOVED***
					throw new Error("unexpected error: " + JSON.stringify(resp.error) + "or status: " + resp.status)
				***REMOVED***
				if (resp.message.data.v !== "3") ***REMOVED***
					throw new Error("unexpected resp message data")
				***REMOVED***`,
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan workerMetrics.SampleContainer, _ error) ***REMOVED***
					samplesBuf := workerMetrics.GetBufferedSamples(samples)
					assertMetricEmitted(t, workerMetrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.any.testing.AnyTestService/Sum"))
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
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan workerMetrics.SampleContainer, _ error) ***REMOVED***
					samplesBuf := workerMetrics.GetBufferedSamples(samples)
					assertMetricEmitted(t, workerMetrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/UnaryCall"))
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
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan workerMetrics.SampleContainer, _ error) ***REMOVED***
					samplesBuf := workerMetrics.GetBufferedSamples(samples)
					assertMetricEmitted(t, workerMetrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
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
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan workerMetrics.SampleContainer, _ error) ***REMOVED***
					samplesBuf := workerMetrics.GetBufferedSamples(samples)
					assertMetricEmitted(t, workerMetrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
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
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan workerMetrics.SampleContainer, _ error) ***REMOVED***
					samplesBuf := workerMetrics.GetBufferedSamples(samples)
					assertMetricEmitted(t, workerMetrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/EmptyCall"))
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

			m, ok := New().NewModuleInstance(ts.VU).(*ModuleInstance)
			require.True(t, ok)
			require.NoError(t, ts.VU.Runtime().Set("grpc", m.Exports().Named))

			// setup necessary environment if needed by a test
			if tt.setup != nil ***REMOVED***
				tt.setup(ts.httpBin)
			***REMOVED***

			replace := func(code string) (goja.Value, error) ***REMOVED***
				return ts.VU.Runtime().RunString(ts.httpBin.Replacer.Replace(code))
			***REMOVED***

			val, err := replace(tt.initString.code)
			assertResponse(t, tt.initString, err, val, ts)

			root, err := libWorker.NewGroup("", nil)
			require.NoError(t, err)
			state := &libWorker.State***REMOVED***
				Group:     root,
				Dialer:    ts.httpBin.Dialer,
				TLSConfig: ts.httpBin.TLSClientConfig,
				Samples:   ts.samples,
				Options: libWorker.Options***REMOVED***
					SystemTags: workerMetrics.NewSystemTagSet(
						workerMetrics.TagName,
						workerMetrics.TagURL,
					),
					UserAgent: null.StringFrom("k6-test"),
				***REMOVED***,
				BuiltinMetrics: workerMetrics.RegisterBuiltinMetrics(
					workerMetrics.NewRegistry(),
				),
				Tags: libWorker.NewTagMap(nil),
			***REMOVED***
			ts.MoveToVUContext(state)
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

			grpcext.DebugStat(logger.WithField("source", "test"), tt.stat, "full")
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
			StateField: &libWorker.State***REMOVED***
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
