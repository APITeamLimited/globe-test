// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ptypes

import (
	"errors"
	"fmt"
	"time"

	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
)

// Range of google.protobuf.Duration as specified in timestamp.proto.
const (
	// Seconds field of the earliest valid Timestamp.
	// This is time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	minValidSeconds = -62135596800
	// Seconds field just after the latest valid Timestamp.
	// This is time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	maxValidSeconds = 253402300800
)

// Timestamp converts a timestamppb.Timestamp to a time.Time.
// It returns an error if the argument is invalid.
//
// Unlike most Go functions, if Timestamp returns an error, the first return
// value is not the zero time.Time. Instead, it is the value obtained from the
// time.Unix function when passed the contents of the Timestamp, in the UTC
// locale. This may or may not be a meaningful time; many invalid Timestamps
// do map to valid time.Times.
//
// A nil Timestamp returns an error. The first return value in that case is
// undefined.
func Timestamp(ts *timestamppb.Timestamp) (time.Time, error) ***REMOVED***
	// Don't return the zero value on error, because corresponds to a valid
	// timestamp. Instead return whatever time.Unix gives us.
	var t time.Time
	if ts == nil ***REMOVED***
		t = time.Unix(0, 0).UTC() // treat nil like the empty Timestamp
	***REMOVED*** else ***REMOVED***
		t = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	***REMOVED***
	return t, validateTimestamp(ts)
***REMOVED***

// TimestampNow returns a google.protobuf.Timestamp for the current time.
func TimestampNow() *timestamppb.Timestamp ***REMOVED***
	ts, err := TimestampProto(time.Now())
	if err != nil ***REMOVED***
		panic("ptypes: time.Now() out of Timestamp range")
	***REMOVED***
	return ts
***REMOVED***

// TimestampProto converts the time.Time to a google.protobuf.Timestamp proto.
// It returns an error if the resulting Timestamp is invalid.
func TimestampProto(t time.Time) (*timestamppb.Timestamp, error) ***REMOVED***
	ts := &timestamppb.Timestamp***REMOVED***
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	***REMOVED***
	if err := validateTimestamp(ts); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return ts, nil
***REMOVED***

// TimestampString returns the RFC 3339 string for valid Timestamps.
// For invalid Timestamps, it returns an error message in parentheses.
func TimestampString(ts *timestamppb.Timestamp) string ***REMOVED***
	t, err := Timestamp(ts)
	if err != nil ***REMOVED***
		return fmt.Sprintf("(%v)", err)
	***REMOVED***
	return t.Format(time.RFC3339Nano)
***REMOVED***

// validateTimestamp determines whether a Timestamp is valid.
// A valid timestamp represents a time in the range [0001-01-01, 10000-01-01)
// and has a Nanos field in the range [0, 1e9).
//
// If the Timestamp is valid, validateTimestamp returns nil.
// Otherwise, it returns an error that describes the problem.
//
// Every valid Timestamp can be represented by a time.Time,
// but the converse is not true.
func validateTimestamp(ts *timestamppb.Timestamp) error ***REMOVED***
	if ts == nil ***REMOVED***
		return errors.New("timestamp: nil Timestamp")
	***REMOVED***
	if ts.Seconds < minValidSeconds ***REMOVED***
		return fmt.Errorf("timestamp: %v before 0001-01-01", ts)
	***REMOVED***
	if ts.Seconds >= maxValidSeconds ***REMOVED***
		return fmt.Errorf("timestamp: %v after 10000-01-01", ts)
	***REMOVED***
	if ts.Nanos < 0 || ts.Nanos >= 1e9 ***REMOVED***
		return fmt.Errorf("timestamp: %v: nanos not in range [0, 1e9)", ts)
	***REMOVED***
	return nil
***REMOVED***
