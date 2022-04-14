// Package grpcext allows gRPC requests collecting stats info.
package grpcext

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"

	protov1 "github.com/golang/protobuf/proto" //nolint:staticcheck,nolintlint // this is the old v1 version
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstats "google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Request represents a gRPC request.
type Request struct ***REMOVED***
	MethodDescriptor protoreflect.MethodDescriptor
	Tags             map[string]string
	Message          []byte
***REMOVED***

// Response represents a gRPC response.
type Response struct ***REMOVED***
	Message  interface***REMOVED******REMOVED***
	Error    interface***REMOVED******REMOVED***
	Headers  map[string][]string
	Trailers map[string][]string
	Status   codes.Code
***REMOVED***

type clientConnCloser interface ***REMOVED***
	grpc.ClientConnInterface
	Close() error
***REMOVED***

// Conn is a gRPC client connection.
type Conn struct ***REMOVED***
	raw clientConnCloser
***REMOVED***

// DefaultOptions generates an option set
// with common options for requests from a VU.
func DefaultOptions(vu modules.VU) []grpc.DialOption ***REMOVED***
	dialer := func(ctx context.Context, addr string) (net.Conn, error) ***REMOVED***
		return vu.State().Dialer.DialContext(ctx, "tcp", addr)
	***REMOVED***

	return []grpc.DialOption***REMOVED***
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithReturnConnectionError(),
		grpc.WithStatsHandler(statsHandler***REMOVED***vu: vu***REMOVED***),
		grpc.WithContextDialer(dialer),
	***REMOVED***
***REMOVED***

// Dial establish a gRPC connection.
func Dial(ctx context.Context, addr string, options ...grpc.DialOption) (*Conn, error) ***REMOVED***
	conn, err := grpc.DialContext(ctx, addr, options...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Conn***REMOVED***
		raw: conn,
	***REMOVED***, nil
***REMOVED***

// ReflectionClient returns a reflection client based on the current connection.
func (c *Conn) ReflectionClient() (*ReflectionClient, error) ***REMOVED***
	return &ReflectionClient***REMOVED***Conn: c.raw***REMOVED***, nil
***REMOVED***

// Invoke executes a unary gRPC request.
func (c *Conn) Invoke(
	ctx context.Context,
	url string,
	md metadata.MD,
	req Request,
	opts ...grpc.CallOption) (*Response, error) ***REMOVED***
	if url == "" ***REMOVED***
		return nil, fmt.Errorf("url is required")
	***REMOVED***
	if req.MethodDescriptor == nil ***REMOVED***
		return nil, fmt.Errorf("request method descriptor is required")
	***REMOVED***
	if len(req.Message) == 0 ***REMOVED***
		return nil, fmt.Errorf("request message is required")
	***REMOVED***

	ctx = metadata.NewOutgoingContext(ctx, md)

	reqdm := dynamicpb.NewMessage(req.MethodDescriptor.Input())
	if err := protojson.Unmarshal(req.Message, reqdm); err != nil ***REMOVED***
		return nil, fmt.Errorf("unable to serialise request object to protocol buffer: %w", err)
	***REMOVED***

	ctx = withTags(ctx, req.Tags)

	resp := dynamicpb.NewMessage(req.MethodDescriptor.Output())
	header, trailer := metadata.New(nil), metadata.New(nil)

	copts := make([]grpc.CallOption, 0, len(opts)+2)
	copts = append(copts, opts...)
	copts = append(copts, grpc.Header(&header), grpc.Trailer(&trailer))

	err := c.raw.Invoke(ctx, url, reqdm, resp, copts...)

	response := Response***REMOVED***
		Headers:  header,
		Trailers: trailer,
	***REMOVED***

	marshaler := protojson.MarshalOptions***REMOVED***EmitUnpopulated: true***REMOVED***

	if err != nil ***REMOVED***
		sterr := status.Convert(err)
		response.Status = sterr.Code()

		// (rogchap) when you access a JSON property in goja, you are actually accessing the underling
		// Go type (struct, map, slice etc); because these are dynamic messages the Unmarshaled JSON does
		// not map back to a "real" field or value (as a normal Go type would). If we don't marshal and then
		// unmarshal back to a map, you will get "undefined" when accessing JSON properties, even when
		// JSON.Stringify() shows the object to be correctly present.

		raw, _ := marshaler.Marshal(sterr.Proto())
		errMsg := make(map[string]interface***REMOVED******REMOVED***)
		_ = json.Unmarshal(raw, &errMsg)
		response.Error = errMsg
	***REMOVED***

	if resp != nil ***REMOVED***
		// (rogchap) there is a lot of marshaling/unmarshaling here, but if we just pass the dynamic message
		// the default Marshaller would be used, which would strip any zero/default values from the JSON.
		// eg. given this message:
		// message Point ***REMOVED***
		//    double x = 1;
		// 	  double y = 2;
		// 	  double z = 3;
		// ***REMOVED***
		// and a value like this:
		// msg := Point***REMOVED***X: 6, Y: 4, Z: 0***REMOVED***
		// would result in JSON output:
		// ***REMOVED***"x":6,"y":4***REMOVED***
		// rather than the desired:
		// ***REMOVED***"x":6,"y":4,"z":0***REMOVED***
		raw, _ := marshaler.Marshal(resp)
		msg := make(map[string]interface***REMOVED******REMOVED***)
		_ = json.Unmarshal(raw, &msg)
		response.Message = msg
	***REMOVED***
	return &response, nil
***REMOVED***

// Close closes the underhood connection.
func (c *Conn) Close() error ***REMOVED***
	return c.raw.Close()
***REMOVED***

type statsHandler struct ***REMOVED***
	vu modules.VU
***REMOVED***

// TagConn implements the grpcstats.Handler interface
func (statsHandler) TagConn(ctx context.Context, _ *grpcstats.ConnTagInfo) context.Context ***REMOVED*** // noop
	return ctx
***REMOVED***

// HandleConn implements the grpcstats.Handler interface
func (statsHandler) HandleConn(context.Context, grpcstats.ConnStats) ***REMOVED***
	// noop
***REMOVED***

// TagRPC implements the grpcstats.Handler interface
func (statsHandler) TagRPC(ctx context.Context, _ *grpcstats.RPCTagInfo) context.Context ***REMOVED***
	// noop
	return ctx
***REMOVED***

// HandleRPC implements the grpcstats.Handler interface
func (h statsHandler) HandleRPC(ctx context.Context, stat grpcstats.RPCStats) ***REMOVED***
	state := h.vu.State()
	tags := getTags(ctx)
	switch s := stat.(type) ***REMOVED***
	case *grpcstats.OutHeader:
		if state.Options.SystemTags.Has(metrics.TagIP) && s.RemoteAddr != nil ***REMOVED***
			if ip, _, err := net.SplitHostPort(s.RemoteAddr.String()); err == nil ***REMOVED***
				tags["ip"] = ip
			***REMOVED***
		***REMOVED***
	case *grpcstats.End:
		if state.Options.SystemTags.Has(metrics.TagStatus) ***REMOVED***
			tags["status"] = strconv.Itoa(int(status.Code(s.Error)))
		***REMOVED***

		mTags := map[string]string(tags)
		sampleTags := metrics.IntoSampleTags(&mTags)
		metrics.PushIfNotDone(ctx, state.Samples, metrics.ConnectedSamples***REMOVED***
			Samples: []metrics.Sample***REMOVED***
				***REMOVED***
					Metric: state.BuiltinMetrics.GRPCReqDuration,
					Tags:   sampleTags,
					Value:  metrics.D(s.EndTime.Sub(s.BeginTime)),
					Time:   s.EndTime,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***)
	***REMOVED***

	// (rogchap) Re-using --http-debug flag as gRPC is technically still HTTP
	if state.Options.HTTPDebug.String != "" ***REMOVED***
		logger := state.Logger.WithField("source", "http-debug")
		httpDebugOption := state.Options.HTTPDebug.String
		DebugStat(logger, stat, httpDebugOption)
	***REMOVED***
***REMOVED***

// DebugStat prints debugging information based on RPCStats.
func DebugStat(logger logrus.FieldLogger, stat grpcstats.RPCStats, httpDebugOption string) ***REMOVED***
	switch s := stat.(type) ***REMOVED***
	case *grpcstats.OutHeader:
		logger.Infof("Out Header:\nFull Method: %s\nRemote Address: %s\n%s\n",
			s.FullMethod, s.RemoteAddr, formatMetadata(s.Header))
	case *grpcstats.OutTrailer:
		if len(s.Trailer) > 0 ***REMOVED***
			logger.Infof("Out Trailer:\n%s\n", formatMetadata(s.Trailer))
		***REMOVED***
	case *grpcstats.OutPayload:
		if httpDebugOption == "full" ***REMOVED***
			logger.Infof("Out Payload:\nWire Length: %d\nSent Time: %s\n%s\n\n",
				s.WireLength, s.SentTime, formatPayload(s.Payload))
		***REMOVED***
	case *grpcstats.InHeader:
		if len(s.Header) > 0 ***REMOVED***
			logger.Infof("In Header:\nWire Length: %d\n%s\n", s.WireLength, formatMetadata(s.Header))
		***REMOVED***
	case *grpcstats.InTrailer:
		if len(s.Trailer) > 0 ***REMOVED***
			logger.Infof("In Trailer:\nWire Length: %d\n%s\n", s.WireLength, formatMetadata(s.Trailer))
		***REMOVED***
	case *grpcstats.InPayload:
		if httpDebugOption == "full" ***REMOVED***
			logger.Infof("In Payload:\nWire Length: %d\nReceived Time: %s\n%s\n\n",
				s.WireLength, s.RecvTime, formatPayload(s.Payload))
		***REMOVED***
	***REMOVED***
***REMOVED***

func formatMetadata(md metadata.MD) string ***REMOVED***
	var sb strings.Builder
	for k, v := range md ***REMOVED***
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(strings.Join(v, ", "))
		sb.WriteRune('\n')
	***REMOVED***

	return sb.String()
***REMOVED***

func formatPayload(payload interface***REMOVED******REMOVED***) string ***REMOVED***
	msg, ok := payload.(proto.Message)
	if !ok ***REMOVED***
		// check to see if we are dealing with a APIv1 message
		msgV1, ok := payload.(protov1.Message)
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		msg = protov1.MessageV2(msgV1)
	***REMOVED***

	marshaler := prototext.MarshalOptions***REMOVED***
		Multiline: true,
		Indent:    "  ",
	***REMOVED***
	b, err := marshaler.Marshal(msg)
	if err != nil ***REMOVED***
		return ""
	***REMOVED***
	return string(b)
***REMOVED***

type ctxKeyTags struct***REMOVED******REMOVED***

type reqtags map[string]string

func withTags(ctx context.Context, tags reqtags) context.Context ***REMOVED***
	if tags == nil ***REMOVED***
		tags = make(map[string]string)
	***REMOVED***
	return context.WithValue(ctx, ctxKeyTags***REMOVED******REMOVED***, tags)
***REMOVED***

func getTags(ctx context.Context) reqtags ***REMOVED***
	v := ctx.Value(ctxKeyTags***REMOVED******REMOVED***)
	if v == nil ***REMOVED***
		return make(map[string]string)
	***REMOVED***
	return v.(reqtags)
***REMOVED***
