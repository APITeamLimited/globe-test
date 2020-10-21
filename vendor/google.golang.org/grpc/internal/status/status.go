/*
 *
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package status implements errors returned by gRPC.  These errors are
// serialized and transmitted on the wire between server and client, and allow
// for additional data to be transmitted via the Details field in the status
// proto.  gRPC service handlers should return an error created by this
// package, and gRPC clients should expect a corresponding error to be
// returned from the RPC call.
//
// This package upholds the invariants that a non-nil error may not
// contain an OK code, and an OK code must result in a nil error.
package status

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// Status represents an RPC status code, message, and details.  It is immutable
// and should be created with New, Newf, or FromProto.
type Status struct ***REMOVED***
	s *spb.Status
***REMOVED***

// New returns a Status representing c and msg.
func New(c codes.Code, msg string) *Status ***REMOVED***
	return &Status***REMOVED***s: &spb.Status***REMOVED***Code: int32(c), Message: msg***REMOVED******REMOVED***
***REMOVED***

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(c codes.Code, format string, a ...interface***REMOVED******REMOVED***) *Status ***REMOVED***
	return New(c, fmt.Sprintf(format, a...))
***REMOVED***

// FromProto returns a Status representing s.
func FromProto(s *spb.Status) *Status ***REMOVED***
	return &Status***REMOVED***s: proto.Clone(s).(*spb.Status)***REMOVED***
***REMOVED***

// Err returns an error representing c and msg.  If c is OK, returns nil.
func Err(c codes.Code, msg string) error ***REMOVED***
	return New(c, msg).Err()
***REMOVED***

// Errorf returns Error(c, fmt.Sprintf(format, a...)).
func Errorf(c codes.Code, format string, a ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return Err(c, fmt.Sprintf(format, a...))
***REMOVED***

// Code returns the status code contained in s.
func (s *Status) Code() codes.Code ***REMOVED***
	if s == nil || s.s == nil ***REMOVED***
		return codes.OK
	***REMOVED***
	return codes.Code(s.s.Code)
***REMOVED***

// Message returns the message contained in s.
func (s *Status) Message() string ***REMOVED***
	if s == nil || s.s == nil ***REMOVED***
		return ""
	***REMOVED***
	return s.s.Message
***REMOVED***

// Proto returns s's status as an spb.Status proto message.
func (s *Status) Proto() *spb.Status ***REMOVED***
	if s == nil ***REMOVED***
		return nil
	***REMOVED***
	return proto.Clone(s.s).(*spb.Status)
***REMOVED***

// Err returns an immutable error representing s; returns nil if s.Code() is OK.
func (s *Status) Err() error ***REMOVED***
	if s.Code() == codes.OK ***REMOVED***
		return nil
	***REMOVED***
	return &Error***REMOVED***e: s.Proto()***REMOVED***
***REMOVED***

// WithDetails returns a new status with the provided details messages appended to the status.
// If any errors are encountered, it returns nil and the first error encountered.
func (s *Status) WithDetails(details ...proto.Message) (*Status, error) ***REMOVED***
	if s.Code() == codes.OK ***REMOVED***
		return nil, errors.New("no error details for status with code OK")
	***REMOVED***
	// s.Code() != OK implies that s.Proto() != nil.
	p := s.Proto()
	for _, detail := range details ***REMOVED***
		any, err := ptypes.MarshalAny(detail)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		p.Details = append(p.Details, any)
	***REMOVED***
	return &Status***REMOVED***s: p***REMOVED***, nil
***REMOVED***

// Details returns a slice of details messages attached to the status.
// If a detail cannot be decoded, the error is returned in place of the detail.
func (s *Status) Details() []interface***REMOVED******REMOVED*** ***REMOVED***
	if s == nil || s.s == nil ***REMOVED***
		return nil
	***REMOVED***
	details := make([]interface***REMOVED******REMOVED***, 0, len(s.s.Details))
	for _, any := range s.s.Details ***REMOVED***
		detail := &ptypes.DynamicAny***REMOVED******REMOVED***
		if err := ptypes.UnmarshalAny(any, detail); err != nil ***REMOVED***
			details = append(details, err)
			continue
		***REMOVED***
		details = append(details, detail.Message)
	***REMOVED***
	return details
***REMOVED***

// Error wraps a pointer of a status proto. It implements error and Status,
// and a nil *Error should never be returned by this package.
type Error struct ***REMOVED***
	e *spb.Status
***REMOVED***

func (e *Error) Error() string ***REMOVED***
	return fmt.Sprintf("rpc error: code = %s desc = %s", codes.Code(e.e.GetCode()), e.e.GetMessage())
***REMOVED***

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *Status ***REMOVED***
	return FromProto(e.e)
***REMOVED***

// Is implements future error.Is functionality.
// A Error is equivalent if the code and message are identical.
func (e *Error) Is(target error) bool ***REMOVED***
	tse, ok := target.(*Error)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return proto.Equal(e.e, tse.e)
***REMOVED***
