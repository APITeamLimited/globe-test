// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

// ErrUnacknowledgedWrite is returned by operations that have an unacknowledged write concern.
var ErrUnacknowledgedWrite = errors.New("unacknowledged write")

// ErrClientDisconnected is returned when disconnected Client is used to run an operation.
var ErrClientDisconnected = errors.New("client is disconnected")

// ErrNilDocument is returned when a nil document is passed to a CRUD method.
var ErrNilDocument = errors.New("document is nil")

// ErrNilValue is returned when a nil value is passed to a CRUD method.
var ErrNilValue = errors.New("value is nil")

// ErrEmptySlice is returned when an empty slice is passed to a CRUD method that requires a non-empty slice.
var ErrEmptySlice = errors.New("must provide at least one element in input slice")

// ErrMapForOrderedArgument is returned when a map with multiple keys is passed to a CRUD method for an ordered parameter
type ErrMapForOrderedArgument struct ***REMOVED***
	ParamName string
***REMOVED***

// Error implements the error interface.
func (e ErrMapForOrderedArgument) Error() string ***REMOVED***
	return fmt.Sprintf("multi-key map passed in for ordered parameter %v", e.ParamName)
***REMOVED***

func replaceErrors(err error) error ***REMOVED***
	if err == topology.ErrTopologyClosed ***REMOVED***
		return ErrClientDisconnected
	***REMOVED***
	if de, ok := err.(driver.Error); ok ***REMOVED***
		return CommandError***REMOVED***
			Code:    de.Code,
			Message: de.Message,
			Labels:  de.Labels,
			Name:    de.Name,
			Wrapped: de.Wrapped,
			Raw:     bson.Raw(de.Raw),
		***REMOVED***
	***REMOVED***
	if qe, ok := err.(driver.QueryFailureError); ok ***REMOVED***
		// qe.Message is "command failure"
		ce := CommandError***REMOVED***
			Name:    qe.Message,
			Wrapped: qe.Wrapped,
			Raw:     bson.Raw(qe.Response),
		***REMOVED***

		dollarErr, err := qe.Response.LookupErr("$err")
		if err == nil ***REMOVED***
			ce.Message, _ = dollarErr.StringValueOK()
		***REMOVED***
		code, err := qe.Response.LookupErr("code")
		if err == nil ***REMOVED***
			ce.Code, _ = code.Int32OK()
		***REMOVED***

		return ce
	***REMOVED***
	if me, ok := err.(mongocrypt.Error); ok ***REMOVED***
		return MongocryptError***REMOVED***Code: me.Code, Message: me.Message***REMOVED***
	***REMOVED***

	return err
***REMOVED***

// IsDuplicateKeyError returns true if err is a duplicate key error
func IsDuplicateKeyError(err error) bool ***REMOVED***
	// handles SERVER-7164 and SERVER-11493
	for ; err != nil; err = unwrap(err) ***REMOVED***
		if e, ok := err.(ServerError); ok ***REMOVED***
			return e.HasErrorCode(11000) || e.HasErrorCode(11001) || e.HasErrorCode(12582) ||
				e.HasErrorCodeWithMessage(16460, " E11000 ")
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// IsTimeout returns true if err is from a timeout
func IsTimeout(err error) bool ***REMOVED***
	for ; err != nil; err = unwrap(err) ***REMOVED***
		// check unwrappable errors together
		if err == context.DeadlineExceeded ***REMOVED***
			return true
		***REMOVED***
		if err == driver.ErrDeadlineWouldBeExceeded ***REMOVED***
			return true
		***REMOVED***
		if ne, ok := err.(net.Error); ok ***REMOVED***
			return ne.Timeout()
		***REMOVED***
		//timeout error labels
		if le, ok := err.(labeledError); ok ***REMOVED***
			if le.HasErrorLabel("NetworkTimeoutError") || le.HasErrorLabel("ExceededTimeLimitError") ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// unwrap returns the inner error if err implements Unwrap(), otherwise it returns nil.
func unwrap(err error) error ***REMOVED***
	u, ok := err.(interface ***REMOVED***
		Unwrap() error
	***REMOVED***)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return u.Unwrap()
***REMOVED***

// errorHasLabel returns true if err contains the specified label
func errorHasLabel(err error, label string) bool ***REMOVED***
	for ; err != nil; err = unwrap(err) ***REMOVED***
		if le, ok := err.(labeledError); ok && le.HasErrorLabel(label) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// IsNetworkError returns true if err is a network error
func IsNetworkError(err error) bool ***REMOVED***
	return errorHasLabel(err, "NetworkError")
***REMOVED***

// MongocryptError represents an libmongocrypt error during client-side encryption.
type MongocryptError struct ***REMOVED***
	Code    int32
	Message string
***REMOVED***

// Error implements the error interface.
func (m MongocryptError) Error() string ***REMOVED***
	return fmt.Sprintf("mongocrypt error %d: %v", m.Code, m.Message)
***REMOVED***

// EncryptionKeyVaultError represents an error while communicating with the key vault collection during client-side
// encryption.
type EncryptionKeyVaultError struct ***REMOVED***
	Wrapped error
***REMOVED***

// Error implements the error interface.
func (ekve EncryptionKeyVaultError) Error() string ***REMOVED***
	return fmt.Sprintf("key vault communication error: %v", ekve.Wrapped)
***REMOVED***

// Unwrap returns the underlying error.
func (ekve EncryptionKeyVaultError) Unwrap() error ***REMOVED***
	return ekve.Wrapped
***REMOVED***

// MongocryptdError represents an error while communicating with mongocryptd during client-side encryption.
type MongocryptdError struct ***REMOVED***
	Wrapped error
***REMOVED***

// Error implements the error interface.
func (e MongocryptdError) Error() string ***REMOVED***
	return fmt.Sprintf("mongocryptd communication error: %v", e.Wrapped)
***REMOVED***

// Unwrap returns the underlying error.
func (e MongocryptdError) Unwrap() error ***REMOVED***
	return e.Wrapped
***REMOVED***

type labeledError interface ***REMOVED***
	error
	// HasErrorLabel returns true if the error contains the specified label.
	HasErrorLabel(string) bool
***REMOVED***

// ServerError is the interface implemented by errors returned from the server. Custom implementations of this
// interface should not be used in production.
type ServerError interface ***REMOVED***
	error
	// HasErrorCode returns true if the error has the specified code.
	HasErrorCode(int) bool
	// HasErrorLabel returns true if the error contains the specified label.
	HasErrorLabel(string) bool
	// HasErrorMessage returns true if the error contains the specified message.
	HasErrorMessage(string) bool
	// HasErrorCodeWithMessage returns true if any of the contained errors have the specified code and message.
	HasErrorCodeWithMessage(int, string) bool

	serverError()
***REMOVED***

var _ ServerError = CommandError***REMOVED******REMOVED***
var _ ServerError = WriteError***REMOVED******REMOVED***
var _ ServerError = WriteException***REMOVED******REMOVED***
var _ ServerError = BulkWriteException***REMOVED******REMOVED***

// CommandError represents a server error during execution of a command. This can be returned by any operation.
type CommandError struct ***REMOVED***
	Code    int32
	Message string
	Labels  []string // Categories to which the error belongs
	Name    string   // A human-readable name corresponding to the error code
	Wrapped error    // The underlying error, if one exists.
	Raw     bson.Raw // The original server response containing the error.
***REMOVED***

// Error implements the error interface.
func (e CommandError) Error() string ***REMOVED***
	if e.Name != "" ***REMOVED***
		return fmt.Sprintf("(%v) %v", e.Name, e.Message)
	***REMOVED***
	return e.Message
***REMOVED***

// Unwrap returns the underlying error.
func (e CommandError) Unwrap() error ***REMOVED***
	return e.Wrapped
***REMOVED***

// HasErrorCode returns true if the error has the specified code.
func (e CommandError) HasErrorCode(code int) bool ***REMOVED***
	return int(e.Code) == code
***REMOVED***

// HasErrorLabel returns true if the error contains the specified label.
func (e CommandError) HasErrorLabel(label string) bool ***REMOVED***
	if e.Labels != nil ***REMOVED***
		for _, l := range e.Labels ***REMOVED***
			if l == label ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// HasErrorMessage returns true if the error contains the specified message.
func (e CommandError) HasErrorMessage(message string) bool ***REMOVED***
	return strings.Contains(e.Message, message)
***REMOVED***

// HasErrorCodeWithMessage returns true if the error has the specified code and Message contains the specified message.
func (e CommandError) HasErrorCodeWithMessage(code int, message string) bool ***REMOVED***
	return int(e.Code) == code && strings.Contains(e.Message, message)
***REMOVED***

// IsMaxTimeMSExpiredError returns true if the error is a MaxTimeMSExpired error.
func (e CommandError) IsMaxTimeMSExpiredError() bool ***REMOVED***
	return e.Code == 50 || e.Name == "MaxTimeMSExpired"
***REMOVED***

// serverError implements the ServerError interface.
func (e CommandError) serverError() ***REMOVED******REMOVED***

// WriteError is an error that occurred during execution of a write operation. This error type is only returned as part
// of a WriteException or BulkWriteException.
type WriteError struct ***REMOVED***
	// The index of the write in the slice passed to an InsertMany or BulkWrite operation that caused this error.
	Index int

	Code    int
	Message string
	Details bson.Raw

	// The original write error from the server response.
	Raw bson.Raw
***REMOVED***

func (we WriteError) Error() string ***REMOVED***
	msg := we.Message
	if len(we.Details) > 0 ***REMOVED***
		msg = fmt.Sprintf("%s: %s", msg, we.Details.String())
	***REMOVED***
	return msg
***REMOVED***

// HasErrorCode returns true if the error has the specified code.
func (we WriteError) HasErrorCode(code int) bool ***REMOVED***
	return we.Code == code
***REMOVED***

// HasErrorLabel returns true if the error contains the specified label. WriteErrors do not contain labels,
// so we always return false.
func (we WriteError) HasErrorLabel(label string) bool ***REMOVED***
	return false
***REMOVED***

// HasErrorMessage returns true if the error contains the specified message.
func (we WriteError) HasErrorMessage(message string) bool ***REMOVED***
	return strings.Contains(we.Message, message)
***REMOVED***

// HasErrorCodeWithMessage returns true if the error has the specified code and Message contains the specified message.
func (we WriteError) HasErrorCodeWithMessage(code int, message string) bool ***REMOVED***
	return we.Code == code && strings.Contains(we.Message, message)
***REMOVED***

// serverError implements the ServerError interface.
func (we WriteError) serverError() ***REMOVED******REMOVED***

// WriteErrors is a group of write errors that occurred during execution of a write operation.
type WriteErrors []WriteError

// Error implements the error interface.
func (we WriteErrors) Error() string ***REMOVED***
	errs := make([]error, len(we))
	for i := 0; i < len(we); i++ ***REMOVED***
		errs[i] = we[i]
	***REMOVED***
	// WriteErrors isn't returned from batch operations, but we can still use the same formatter.
	return "write errors: " + joinBatchErrors(errs)
***REMOVED***

func writeErrorsFromDriverWriteErrors(errs driver.WriteErrors) WriteErrors ***REMOVED***
	wes := make(WriteErrors, 0, len(errs))
	for _, err := range errs ***REMOVED***
		wes = append(wes, WriteError***REMOVED***
			Index:   int(err.Index),
			Code:    int(err.Code),
			Message: err.Message,
			Details: bson.Raw(err.Details),
			Raw:     bson.Raw(err.Raw),
		***REMOVED***)
	***REMOVED***
	return wes
***REMOVED***

// WriteConcernError represents a write concern failure during execution of a write operation. This error type is only
// returned as part of a WriteException or a BulkWriteException.
type WriteConcernError struct ***REMOVED***
	Name    string
	Code    int
	Message string
	Details bson.Raw
	Raw     bson.Raw // The original write concern error from the server response.
***REMOVED***

// Error implements the error interface.
func (wce WriteConcernError) Error() string ***REMOVED***
	if wce.Name != "" ***REMOVED***
		return fmt.Sprintf("(%v) %v", wce.Name, wce.Message)
	***REMOVED***
	return wce.Message
***REMOVED***

// WriteException is the error type returned by the InsertOne, DeleteOne, DeleteMany, UpdateOne, UpdateMany, and
// ReplaceOne operations.
type WriteException struct ***REMOVED***
	// The write concern error that occurred, or nil if there was none.
	WriteConcernError *WriteConcernError

	// The write errors that occurred during operation execution.
	WriteErrors WriteErrors

	// The categories to which the exception belongs.
	Labels []string

	// The original server response containing the error.
	Raw bson.Raw
***REMOVED***

// Error implements the error interface.
func (mwe WriteException) Error() string ***REMOVED***
	causes := make([]string, 0, 2)
	if mwe.WriteConcernError != nil ***REMOVED***
		causes = append(causes, "write concern error: "+mwe.WriteConcernError.Error())
	***REMOVED***
	if len(mwe.WriteErrors) > 0 ***REMOVED***
		// The WriteErrors error message already starts with "write errors:", so don't add it to the
		// error message again.
		causes = append(causes, mwe.WriteErrors.Error())
	***REMOVED***

	message := "write exception: "
	if len(causes) == 0 ***REMOVED***
		return message + "no causes"
	***REMOVED***
	return message + strings.Join(causes, ", ")
***REMOVED***

// HasErrorCode returns true if the error has the specified code.
func (mwe WriteException) HasErrorCode(code int) bool ***REMOVED***
	if mwe.WriteConcernError != nil && mwe.WriteConcernError.Code == code ***REMOVED***
		return true
	***REMOVED***
	for _, we := range mwe.WriteErrors ***REMOVED***
		if we.Code == code ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// HasErrorLabel returns true if the error contains the specified label.
func (mwe WriteException) HasErrorLabel(label string) bool ***REMOVED***
	if mwe.Labels != nil ***REMOVED***
		for _, l := range mwe.Labels ***REMOVED***
			if l == label ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// HasErrorMessage returns true if the error contains the specified message.
func (mwe WriteException) HasErrorMessage(message string) bool ***REMOVED***
	if mwe.WriteConcernError != nil && strings.Contains(mwe.WriteConcernError.Message, message) ***REMOVED***
		return true
	***REMOVED***
	for _, we := range mwe.WriteErrors ***REMOVED***
		if strings.Contains(we.Message, message) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// HasErrorCodeWithMessage returns true if any of the contained errors have the specified code and message.
func (mwe WriteException) HasErrorCodeWithMessage(code int, message string) bool ***REMOVED***
	if mwe.WriteConcernError != nil &&
		mwe.WriteConcernError.Code == code && strings.Contains(mwe.WriteConcernError.Message, message) ***REMOVED***
		return true
	***REMOVED***
	for _, we := range mwe.WriteErrors ***REMOVED***
		if we.Code == code && strings.Contains(we.Message, message) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// serverError implements the ServerError interface.
func (mwe WriteException) serverError() ***REMOVED******REMOVED***

func convertDriverWriteConcernError(wce *driver.WriteConcernError) *WriteConcernError ***REMOVED***
	if wce == nil ***REMOVED***
		return nil
	***REMOVED***

	return &WriteConcernError***REMOVED***
		Name:    wce.Name,
		Code:    int(wce.Code),
		Message: wce.Message,
		Details: bson.Raw(wce.Details),
		Raw:     bson.Raw(wce.Raw),
	***REMOVED***
***REMOVED***

// BulkWriteError is an error that occurred during execution of one operation in a BulkWrite. This error type is only
// returned as part of a BulkWriteException.
type BulkWriteError struct ***REMOVED***
	WriteError            // The WriteError that occurred.
	Request    WriteModel // The WriteModel that caused this error.
***REMOVED***

// Error implements the error interface.
func (bwe BulkWriteError) Error() string ***REMOVED***
	return bwe.WriteError.Error()
***REMOVED***

// BulkWriteException is the error type returned by BulkWrite and InsertMany operations.
type BulkWriteException struct ***REMOVED***
	// The write concern error that occurred, or nil if there was none.
	WriteConcernError *WriteConcernError

	// The write errors that occurred during operation execution.
	WriteErrors []BulkWriteError

	// The categories to which the exception belongs.
	Labels []string
***REMOVED***

// Error implements the error interface.
func (bwe BulkWriteException) Error() string ***REMOVED***
	causes := make([]string, 0, 2)
	if bwe.WriteConcernError != nil ***REMOVED***
		causes = append(causes, "write concern error: "+bwe.WriteConcernError.Error())
	***REMOVED***
	if len(bwe.WriteErrors) > 0 ***REMOVED***
		errs := make([]error, len(bwe.WriteErrors))
		for i := 0; i < len(bwe.WriteErrors); i++ ***REMOVED***
			errs[i] = &bwe.WriteErrors[i]
		***REMOVED***
		causes = append(causes, "write errors: "+joinBatchErrors(errs))
	***REMOVED***

	message := "bulk write exception: "
	if len(causes) == 0 ***REMOVED***
		return message + "no causes"
	***REMOVED***
	return "bulk write exception: " + strings.Join(causes, ", ")
***REMOVED***

// HasErrorCode returns true if any of the errors have the specified code.
func (bwe BulkWriteException) HasErrorCode(code int) bool ***REMOVED***
	if bwe.WriteConcernError != nil && bwe.WriteConcernError.Code == code ***REMOVED***
		return true
	***REMOVED***
	for _, we := range bwe.WriteErrors ***REMOVED***
		if we.Code == code ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// HasErrorLabel returns true if the error contains the specified label.
func (bwe BulkWriteException) HasErrorLabel(label string) bool ***REMOVED***
	if bwe.Labels != nil ***REMOVED***
		for _, l := range bwe.Labels ***REMOVED***
			if l == label ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// HasErrorMessage returns true if the error contains the specified message.
func (bwe BulkWriteException) HasErrorMessage(message string) bool ***REMOVED***
	if bwe.WriteConcernError != nil && strings.Contains(bwe.WriteConcernError.Message, message) ***REMOVED***
		return true
	***REMOVED***
	for _, we := range bwe.WriteErrors ***REMOVED***
		if strings.Contains(we.Message, message) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// HasErrorCodeWithMessage returns true if any of the contained errors have the specified code and message.
func (bwe BulkWriteException) HasErrorCodeWithMessage(code int, message string) bool ***REMOVED***
	if bwe.WriteConcernError != nil &&
		bwe.WriteConcernError.Code == code && strings.Contains(bwe.WriteConcernError.Message, message) ***REMOVED***
		return true
	***REMOVED***
	for _, we := range bwe.WriteErrors ***REMOVED***
		if we.Code == code && strings.Contains(we.Message, message) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// serverError implements the ServerError interface.
func (bwe BulkWriteException) serverError() ***REMOVED******REMOVED***

// returnResult is used to determine if a function calling processWriteError should return
// the result or return nil. Since the processWriteError function is used by many different
// methods, both *One and *Many, we need a way to differentiate if the method should return
// the result and the error.
type returnResult int

const (
	rrNone returnResult = 1 << iota // None means do not return the result ever.
	rrOne                           // One means return the result if this was called by a *One method.
	rrMany                          // Many means return the result is this was called by a *Many method.

	rrAll returnResult = rrOne | rrMany // All means always return the result.
)

// processWriteError handles processing the result of a write operation. If the retrunResult matches
// the calling method's type, it should return the result object in addition to the error.
// This function will wrap the errors from other packages and return them as errors from this package.
//
// WriteConcernError will be returned over WriteErrors if both are present.
func processWriteError(err error) (returnResult, error) ***REMOVED***
	switch ***REMOVED***
	case err == driver.ErrUnacknowledgedWrite:
		return rrAll, ErrUnacknowledgedWrite
	case err != nil:
		switch tt := err.(type) ***REMOVED***
		case driver.WriteCommandError:
			return rrMany, WriteException***REMOVED***
				WriteConcernError: convertDriverWriteConcernError(tt.WriteConcernError),
				WriteErrors:       writeErrorsFromDriverWriteErrors(tt.WriteErrors),
				Labels:            tt.Labels,
				Raw:               bson.Raw(tt.Raw),
			***REMOVED***
		default:
			return rrNone, replaceErrors(err)
		***REMOVED***
	default:
		return rrAll, nil
	***REMOVED***
***REMOVED***

// batchErrorsTargetLength is the target length of error messages returned by batch operation
// error types. Try to limit batch error messages to 2kb to prevent problems when printing error
// messages from large batch operations.
const batchErrorsTargetLength = 2000

// joinBatchErrors appends messages from the given errors to a comma-separated string. If the
// string exceeds 2kb, it stops appending error messages and appends the message "+N more errors..."
// to the end.
//
// Example format:
//     "[message 1, message 2, +8 more errors...]"
func joinBatchErrors(errs []error) string ***REMOVED***
	var buf bytes.Buffer
	fmt.Fprint(&buf, "[")
	for idx, err := range errs ***REMOVED***
		if idx != 0 ***REMOVED***
			fmt.Fprint(&buf, ", ")
		***REMOVED***
		// If the error message has exceeded the target error message length, stop appending errors
		// to the message and append the number of remaining errors instead.
		if buf.Len() > batchErrorsTargetLength ***REMOVED***
			fmt.Fprintf(&buf, "+%d more errors...", len(errs)-idx)
			break
		***REMOVED***
		fmt.Fprint(&buf, err.Error())
	***REMOVED***
	fmt.Fprint(&buf, "]")

	return buf.String()
***REMOVED***
