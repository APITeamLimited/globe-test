package grpcext

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// ReflectionClient wraps a grpc.ServerReflectionClient.
type ReflectionClient struct ***REMOVED***
	Conn grpc.ClientConnInterface
***REMOVED***

// Reflect will use the grpc reflection api to make the file descriptors available to request.
// It is called in the connect function the first time the Client.Connect function is called.
func (rc *ReflectionClient) Reflect(ctx context.Context) (*descriptorpb.FileDescriptorSet, error) ***REMOVED***
	client := reflectpb.NewServerReflectionClient(rc.Conn)
	methodClient, err := client.ServerReflectionInfo(ctx)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("can't get server info: %w", err)
	***REMOVED***
	req := &reflectpb.ServerReflectionRequest***REMOVED***
		MessageRequest: &reflectpb.ServerReflectionRequest_ListServices***REMOVED******REMOVED***,
	***REMOVED***
	resp, err := sendReceive(methodClient, req)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("can't list services: %w", err)
	***REMOVED***
	listResp := resp.GetListServicesResponse()
	if listResp == nil ***REMOVED***
		return nil, fmt.Errorf("can't list services, nil response")
	***REMOVED***
	fdset, err := rc.resolveServiceFileDescriptors(methodClient, listResp)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("can't resolve services' file descriptors: %w", err)
	***REMOVED***
	return fdset, nil
***REMOVED***

func (rc *ReflectionClient) resolveServiceFileDescriptors(
	client sendReceiver,
	res *reflectpb.ListServiceResponse,
) (*descriptorpb.FileDescriptorSet, error) ***REMOVED***
	services := res.GetService()
	seen := make(map[fileDescriptorLookupKey]bool, len(services))
	fdset := &descriptorpb.FileDescriptorSet***REMOVED***
		File: make([]*descriptorpb.FileDescriptorProto, 0, len(services)),
	***REMOVED***

	for _, service := range services ***REMOVED***
		req := &reflectpb.ServerReflectionRequest***REMOVED***
			MessageRequest: &reflectpb.ServerReflectionRequest_FileContainingSymbol***REMOVED***
				FileContainingSymbol: service.GetName(),
			***REMOVED***,
		***REMOVED***
		resp, err := sendReceive(client, req)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("can't get method on service %q: %w", service, err)
		***REMOVED***
		fdResp := resp.GetFileDescriptorResponse()
		for _, raw := range fdResp.GetFileDescriptorProto() ***REMOVED***
			var fdp descriptorpb.FileDescriptorProto
			if err = proto.Unmarshal(raw, &fdp); err != nil ***REMOVED***
				return nil, fmt.Errorf("can't unmarshal proto on service %q: %w", service, err)
			***REMOVED***
			fdkey := fileDescriptorLookupKey***REMOVED***
				Package: *fdp.Package,
				Name:    *fdp.Name,
			***REMOVED***
			if seen[fdkey] ***REMOVED***
				// When a proto file contains declarations for multiple services
				// then the same proto file is returned multiple times,
				// this prevents adding the returned proto file as a duplicate.
				continue
			***REMOVED***
			seen[fdkey] = true
			fdset.File = append(fdset.File, &fdp)
		***REMOVED***
	***REMOVED***
	return fdset, nil
***REMOVED***

// sendReceiver is a smaller interface for decoupling
// from `reflectpb.ServerReflection_ServerReflectionInfoClient`,
// that has the dependency from `grpc.ClientStream`,
// which is too much in the case the requirement is to just make a reflection's request.
// It makes the API more restricted and with a controlled surface,
// in this way the testing should be easier also.
type sendReceiver interface ***REMOVED***
	Send(*reflectpb.ServerReflectionRequest) error
	Recv() (*reflectpb.ServerReflectionResponse, error)
***REMOVED***

// sendReceive sends a request to a reflection client and,
// receives a response.
func sendReceive(
	client sendReceiver,
	req *reflectpb.ServerReflectionRequest,
) (*reflectpb.ServerReflectionResponse, error) ***REMOVED***
	if err := client.Send(req); err != nil ***REMOVED***
		return nil, fmt.Errorf("can't send request: %w", err)
	***REMOVED***
	resp, err := client.Recv()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("can't receive response: %w", err)
	***REMOVED***
	return resp, nil
***REMOVED***

type fileDescriptorLookupKey struct ***REMOVED***
	Package string
	Name    string
***REMOVED***
