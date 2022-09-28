// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var (
	retryableCodes          = []int32***REMOVED***11600, 11602, 10107, 13435, 13436, 189, 91, 7, 6, 89, 9001, 262***REMOVED***
	nodeIsRecoveringCodes   = []int32***REMOVED***11600, 11602, 13436, 189, 91***REMOVED***
	notPrimaryCodes         = []int32***REMOVED***10107, 13435, 10058***REMOVED***
	nodeIsShuttingDownCodes = []int32***REMOVED***11600, 91***REMOVED***

	unknownReplWriteConcernCode   = int32(79)
	unsatisfiableWriteConcernCode = int32(100)
)

var (
	// UnknownTransactionCommitResult is an error label for unknown transaction commit results.
	UnknownTransactionCommitResult = "UnknownTransactionCommitResult"
	// TransientTransactionError is an error label for transient errors with transactions.
	TransientTransactionError = "TransientTransactionError"
	// NetworkError is an error label for network errors.
	NetworkError = "NetworkError"
	// RetryableWriteError is an error lable for retryable write errors.
	RetryableWriteError = "RetryableWriteError"
	// ErrCursorNotFound is the cursor not found error for legacy find operations.
	ErrCursorNotFound = errors.New("cursor not found")
	// ErrUnacknowledgedWrite is returned from functions that have an unacknowledged
	// write concern.
	ErrUnacknowledgedWrite = errors.New("unacknowledged write")
	// ErrUnsupportedStorageEngine is returned when a retryable write is attempted against a server
	// that uses a storage engine that does not support retryable writes
	ErrUnsupportedStorageEngine = errors.New("this MongoDB deployment does not support retryable writes. Please add retryWrites=false to your connection string")
	// ErrDeadlineWouldBeExceeded is returned when a Timeout set on an operation would be exceeded
	// if the operation were sent to the server.
	ErrDeadlineWouldBeExceeded = errors.New("operation not sent to server, as Timeout would be exceeded")
)

// QueryFailureError is an error representing a command failure as a document.
type QueryFailureError struct ***REMOVED***
	Message  string
	Response bsoncore.Document
	Wrapped  error
***REMOVED***

// Error implements the error interface.
func (e QueryFailureError) Error() string ***REMOVED***
	return fmt.Sprintf("%s: %v", e.Message, e.Response)
***REMOVED***

// Unwrap returns the underlying error.
func (e QueryFailureError) Unwrap() error ***REMOVED***
	return e.Wrapped
***REMOVED***

// ResponseError is an error parsing the response to a command.
type ResponseError struct ***REMOVED***
	Message string
	Wrapped error
***REMOVED***

// NewCommandResponseError creates a CommandResponseError.
func NewCommandResponseError(msg string, err error) ResponseError ***REMOVED***
	return ResponseError***REMOVED***Message: msg, Wrapped: err***REMOVED***
***REMOVED***

// Error implements the error interface.
func (e ResponseError) Error() string ***REMOVED***
	if e.Wrapped != nil ***REMOVED***
		return fmt.Sprintf("%s: %s", e.Message, e.Wrapped)
	***REMOVED***
	return e.Message
***REMOVED***

// WriteCommandError is an error for a write command.
type WriteCommandError struct ***REMOVED***
	WriteConcernError *WriteConcernError
	WriteErrors       WriteErrors
	Labels            []string
	Raw               bsoncore.Document
***REMOVED***

// UnsupportedStorageEngine returns whether or not the WriteCommandError comes from a retryable write being attempted
// against a server that has a storage engine where they are not supported
func (wce WriteCommandError) UnsupportedStorageEngine() bool ***REMOVED***
	for _, writeError := range wce.WriteErrors ***REMOVED***
		if writeError.Code == 20 && strings.HasPrefix(strings.ToLower(writeError.Message), "transaction numbers") ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (wce WriteCommandError) Error() string ***REMOVED***
	var buf bytes.Buffer
	fmt.Fprint(&buf, "write command error: [")
	fmt.Fprintf(&buf, "***REMOVED***%s***REMOVED***, ", wce.WriteErrors)
	fmt.Fprintf(&buf, "***REMOVED***%s***REMOVED***]", wce.WriteConcernError)
	return buf.String()
***REMOVED***

// Retryable returns true if the error is retryable
func (wce WriteCommandError) Retryable(wireVersion *description.VersionRange) bool ***REMOVED***
	for _, label := range wce.Labels ***REMOVED***
		if label == RetryableWriteError ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	if wireVersion != nil && wireVersion.Max >= 9 ***REMOVED***
		return false
	***REMOVED***

	if wce.WriteConcernError == nil ***REMOVED***
		return false
	***REMOVED***
	return (*wce.WriteConcernError).Retryable()
***REMOVED***

// WriteConcernError is a write concern failure that occurred as a result of a
// write operation.
type WriteConcernError struct ***REMOVED***
	Name            string
	Code            int64
	Message         string
	Details         bsoncore.Document
	Labels          []string
	TopologyVersion *description.TopologyVersion
	Raw             bsoncore.Document
***REMOVED***

func (wce WriteConcernError) Error() string ***REMOVED***
	if wce.Name != "" ***REMOVED***
		return fmt.Sprintf("(%v) %v", wce.Name, wce.Message)
	***REMOVED***
	return wce.Message
***REMOVED***

// Retryable returns true if the error is retryable
func (wce WriteConcernError) Retryable() bool ***REMOVED***
	for _, code := range retryableCodes ***REMOVED***
		if wce.Code == int64(code) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// NodeIsRecovering returns true if this error is a node is recovering error.
func (wce WriteConcernError) NodeIsRecovering() bool ***REMOVED***
	for _, code := range nodeIsRecoveringCodes ***REMOVED***
		if wce.Code == int64(code) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	hasNoCode := wce.Code == 0
	return hasNoCode && strings.Contains(wce.Message, "node is recovering")
***REMOVED***

// NodeIsShuttingDown returns true if this error is a node is shutting down error.
func (wce WriteConcernError) NodeIsShuttingDown() bool ***REMOVED***
	for _, code := range nodeIsShuttingDownCodes ***REMOVED***
		if wce.Code == int64(code) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	hasNoCode := wce.Code == 0
	return hasNoCode && strings.Contains(wce.Message, "node is shutting down")
***REMOVED***

// NotPrimary returns true if this error is a not primary error.
func (wce WriteConcernError) NotPrimary() bool ***REMOVED***
	for _, code := range notPrimaryCodes ***REMOVED***
		if wce.Code == int64(code) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	hasNoCode := wce.Code == 0
	return hasNoCode && strings.Contains(wce.Message, internal.LegacyNotPrimary)
***REMOVED***

// WriteError is a non-write concern failure that occurred as a result of a write
// operation.
type WriteError struct ***REMOVED***
	Index   int64
	Code    int64
	Message string
	Details bsoncore.Document
	Raw     bsoncore.Document
***REMOVED***

func (we WriteError) Error() string ***REMOVED*** return we.Message ***REMOVED***

// WriteErrors is a group of non-write concern failures that occurred as a result
// of a write operation.
type WriteErrors []WriteError

func (we WriteErrors) Error() string ***REMOVED***
	var buf bytes.Buffer
	fmt.Fprint(&buf, "write errors: [")
	for idx, err := range we ***REMOVED***
		if idx != 0 ***REMOVED***
			fmt.Fprintf(&buf, ", ")
		***REMOVED***
		fmt.Fprintf(&buf, "***REMOVED***%s***REMOVED***", err)
	***REMOVED***
	fmt.Fprint(&buf, "]")
	return buf.String()
***REMOVED***

// Error is a command execution error from the database.
type Error struct ***REMOVED***
	Code            int32
	Message         string
	Labels          []string
	Name            string
	Wrapped         error
	TopologyVersion *description.TopologyVersion
	Raw             bsoncore.Document
***REMOVED***

// UnsupportedStorageEngine returns whether e came as a result of an unsupported storage engine
func (e Error) UnsupportedStorageEngine() bool ***REMOVED***
	return e.Code == 20 && strings.HasPrefix(strings.ToLower(e.Message), "transaction numbers")
***REMOVED***

// Error implements the error interface.
func (e Error) Error() string ***REMOVED***
	if e.Name != "" ***REMOVED***
		return fmt.Sprintf("(%v) %v", e.Name, e.Message)
	***REMOVED***
	return e.Message
***REMOVED***

// Unwrap returns the underlying error.
func (e Error) Unwrap() error ***REMOVED***
	return e.Wrapped
***REMOVED***

// HasErrorLabel returns true if the error contains the specified label.
func (e Error) HasErrorLabel(label string) bool ***REMOVED***
	if e.Labels != nil ***REMOVED***
		for _, l := range e.Labels ***REMOVED***
			if l == label ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// RetryableRead returns true if the error is retryable for a read operation
func (e Error) RetryableRead() bool ***REMOVED***
	for _, label := range e.Labels ***REMOVED***
		if label == NetworkError ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	for _, code := range retryableCodes ***REMOVED***
		if e.Code == code ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// RetryableWrite returns true if the error is retryable for a write operation
func (e Error) RetryableWrite(wireVersion *description.VersionRange) bool ***REMOVED***
	for _, label := range e.Labels ***REMOVED***
		if label == NetworkError || label == RetryableWriteError ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	if wireVersion != nil && wireVersion.Max >= 9 ***REMOVED***
		return false
	***REMOVED***
	for _, code := range retryableCodes ***REMOVED***
		if e.Code == code ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// NetworkError returns true if the error is a network error.
func (e Error) NetworkError() bool ***REMOVED***
	for _, label := range e.Labels ***REMOVED***
		if label == NetworkError ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// NodeIsRecovering returns true if this error is a node is recovering error.
func (e Error) NodeIsRecovering() bool ***REMOVED***
	for _, code := range nodeIsRecoveringCodes ***REMOVED***
		if e.Code == code ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	hasNoCode := e.Code == 0
	return hasNoCode && strings.Contains(e.Message, "node is recovering")
***REMOVED***

// NodeIsShuttingDown returns true if this error is a node is shutting down error.
func (e Error) NodeIsShuttingDown() bool ***REMOVED***
	for _, code := range nodeIsShuttingDownCodes ***REMOVED***
		if e.Code == code ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	hasNoCode := e.Code == 0
	return hasNoCode && strings.Contains(e.Message, "node is shutting down")
***REMOVED***

// NotPrimary returns true if this error is a not primary error.
func (e Error) NotPrimary() bool ***REMOVED***
	for _, code := range notPrimaryCodes ***REMOVED***
		if e.Code == code ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	hasNoCode := e.Code == 0
	return hasNoCode && strings.Contains(e.Message, internal.LegacyNotPrimary)
***REMOVED***

// NamespaceNotFound returns true if this errors is a NamespaceNotFound error.
func (e Error) NamespaceNotFound() bool ***REMOVED***
	return e.Code == 26 || e.Message == "ns not found"
***REMOVED***

// ExtractErrorFromServerResponse extracts an error from a server response bsoncore.Document
// if there is one. Also used in testing for SDAM.
func ExtractErrorFromServerResponse(doc bsoncore.Document) error ***REMOVED***
	var errmsg, codeName string
	var code int32
	var labels []string
	var ok bool
	var tv *description.TopologyVersion
	var wcError WriteCommandError
	elems, err := doc.Elements()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, elem := range elems ***REMOVED***
		switch elem.Key() ***REMOVED***
		case "ok":
			switch elem.Value().Type ***REMOVED***
			case bson.TypeInt32:
				if elem.Value().Int32() == 1 ***REMOVED***
					ok = true
				***REMOVED***
			case bson.TypeInt64:
				if elem.Value().Int64() == 1 ***REMOVED***
					ok = true
				***REMOVED***
			case bson.TypeDouble:
				if elem.Value().Double() == 1 ***REMOVED***
					ok = true
				***REMOVED***
			***REMOVED***
		case "errmsg":
			if str, okay := elem.Value().StringValueOK(); okay ***REMOVED***
				errmsg = str
			***REMOVED***
		case "codeName":
			if str, okay := elem.Value().StringValueOK(); okay ***REMOVED***
				codeName = str
			***REMOVED***
		case "code":
			if c, okay := elem.Value().Int32OK(); okay ***REMOVED***
				code = c
			***REMOVED***
		case "errorLabels":
			if arr, okay := elem.Value().ArrayOK(); okay ***REMOVED***
				vals, err := arr.Values()
				if err != nil ***REMOVED***
					continue
				***REMOVED***
				for _, val := range vals ***REMOVED***
					if str, ok := val.StringValueOK(); ok ***REMOVED***
						labels = append(labels, str)
					***REMOVED***
				***REMOVED***

			***REMOVED***
		case "writeErrors":
			arr, exists := elem.Value().ArrayOK()
			if !exists ***REMOVED***
				break
			***REMOVED***
			vals, err := arr.Values()
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			for _, val := range vals ***REMOVED***
				var we WriteError
				doc, exists := val.DocumentOK()
				if !exists ***REMOVED***
					continue
				***REMOVED***
				if index, exists := doc.Lookup("index").AsInt64OK(); exists ***REMOVED***
					we.Index = index
				***REMOVED***
				if code, exists := doc.Lookup("code").AsInt64OK(); exists ***REMOVED***
					we.Code = code
				***REMOVED***
				if msg, exists := doc.Lookup("errmsg").StringValueOK(); exists ***REMOVED***
					we.Message = msg
				***REMOVED***
				if info, exists := doc.Lookup("errInfo").DocumentOK(); exists ***REMOVED***
					we.Details = make([]byte, len(info))
					copy(we.Details, info)
				***REMOVED***
				we.Raw = doc
				wcError.WriteErrors = append(wcError.WriteErrors, we)
			***REMOVED***
		case "writeConcernError":
			doc, exists := elem.Value().DocumentOK()
			if !exists ***REMOVED***
				break
			***REMOVED***
			wcError.WriteConcernError = new(WriteConcernError)
			wcError.WriteConcernError.Raw = doc
			if code, exists := doc.Lookup("code").AsInt64OK(); exists ***REMOVED***
				wcError.WriteConcernError.Code = code
			***REMOVED***
			if name, exists := doc.Lookup("codeName").StringValueOK(); exists ***REMOVED***
				wcError.WriteConcernError.Name = name
			***REMOVED***
			if msg, exists := doc.Lookup("errmsg").StringValueOK(); exists ***REMOVED***
				wcError.WriteConcernError.Message = msg
			***REMOVED***
			if info, exists := doc.Lookup("errInfo").DocumentOK(); exists ***REMOVED***
				wcError.WriteConcernError.Details = make([]byte, len(info))
				copy(wcError.WriteConcernError.Details, info)
			***REMOVED***
			if errLabels, exists := doc.Lookup("errorLabels").ArrayOK(); exists ***REMOVED***
				vals, err := errLabels.Values()
				if err != nil ***REMOVED***
					continue
				***REMOVED***
				for _, val := range vals ***REMOVED***
					if str, ok := val.StringValueOK(); ok ***REMOVED***
						labels = append(labels, str)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case "topologyVersion":
			doc, ok := elem.Value().DocumentOK()
			if !ok ***REMOVED***
				break
			***REMOVED***
			version, err := description.NewTopologyVersion(bson.Raw(doc))
			if err == nil ***REMOVED***
				tv = version
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !ok ***REMOVED***
		if errmsg == "" ***REMOVED***
			errmsg = "command failed"
		***REMOVED***

		return Error***REMOVED***
			Code:            code,
			Message:         errmsg,
			Name:            codeName,
			Labels:          labels,
			TopologyVersion: tv,
			Raw:             doc,
		***REMOVED***
	***REMOVED***

	if len(wcError.WriteErrors) > 0 || wcError.WriteConcernError != nil ***REMOVED***
		wcError.Labels = labels
		if wcError.WriteConcernError != nil ***REMOVED***
			wcError.WriteConcernError.TopologyVersion = tv
		***REMOVED***
		wcError.Raw = doc
		return wcError
	***REMOVED***

	return nil
***REMOVED***
