// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

var (
	firstBatchIdentifier     = "firstBatch"
	nextBatchIdentifier      = "nextBatch"
	listCollectionsNamespace = "system.namespaces"
	listIndexesNamespace     = "system.indexes"

	// ErrFilterType is returned when the filter for a legacy list collections operation is of the wrong type.
	ErrFilterType = errors.New("filter for list collections operation must be a string")
)

func (op Operation) getFullCollectionName(coll string) string ***REMOVED***
	return op.Database + "." + coll
***REMOVED***

func (op Operation) legacyFind(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error ***REMOVED***
	wm, startedInfo, collName, err := op.createLegacyFindWireMessage(dst, desc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	finishedInfo := finishedInformation***REMOVED***
		cmdName:   startedInfo.cmdName,
		requestID: startedInfo.requestID,
		startTime: time.Now(),
		connID:    startedInfo.connID,
	***REMOVED***

	finishedInfo.response, finishedInfo.cmdErr = op.roundTripLegacyCursor(ctx, wm, srvr, conn, collName, firstBatchIdentifier)
	op.publishFinishedEvent(ctx, finishedInfo)

	if finishedInfo.cmdErr != nil ***REMOVED***
		return finishedInfo.cmdErr
	***REMOVED***

	if op.ProcessResponseFn != nil ***REMOVED***
		// CurrentIndex is always 0 in this mode.
		info := ResponseInfo***REMOVED***
			ServerResponse:        finishedInfo.response,
			Server:                srvr,
			Connection:            conn,
			ConnectionDescription: desc.Server,
		***REMOVED***
		return op.ProcessResponseFn(info)
	***REMOVED***
	return nil
***REMOVED***

// returns wire message, collection name, error
func (op Operation) createLegacyFindWireMessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) ***REMOVED***
	info := startedInformation***REMOVED***
		requestID: wiremessage.NextRequestID(),
		cmdName:   "find",
	***REMOVED***

	// call CommandFn on an empty slice rather than dst because the options will need to be converted to legacy
	var cmdDoc bsoncore.Document
	var cmdIndex int32
	var err error

	cmdIndex, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	cmdDoc, err = op.CommandFn(cmdDoc, desc)
	if err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIndex)
	// for monitoring legacy events, the upconverted document should be captured rather than the legacy one
	info.cmd = cmdDoc

	cmdElems, err := cmdDoc.Elements()
	if err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***

	// take each option from the non-legacy command and convert it
	// build options as a byte slice of elements rather than a bsoncore.Document because they will be appended
	// to another document with $query
	var optsElems []byte
	flags := op.secondaryOK(desc)
	var numToSkip, numToReturn, batchSize, limit int32 // numToReturn calculated from batchSize and limit
	var filter, returnFieldsSelector bsoncore.Document
	var collName string
	var singleBatch bool
	for _, elem := range cmdElems ***REMOVED***
		switch elem.Key() ***REMOVED***
		case "find":
			collName = elem.Value().StringValue()
		case "filter":
			filter = elem.Value().Data
		case "sort":
			optsElems = bsoncore.AppendValueElement(optsElems, "$orderby", elem.Value())
		case "hint":
			optsElems = bsoncore.AppendValueElement(optsElems, "$hint", elem.Value())
		case "comment":
			optsElems = bsoncore.AppendValueElement(optsElems, "$comment", elem.Value())
		case "max":
			optsElems = bsoncore.AppendValueElement(optsElems, "$max", elem.Value())
		case "min":
			optsElems = bsoncore.AppendValueElement(optsElems, "$min", elem.Value())
		case "returnKey":
			optsElems = bsoncore.AppendValueElement(optsElems, "$returnKey", elem.Value())
		case "showRecordId":
			optsElems = bsoncore.AppendValueElement(optsElems, "$showDiskLoc", elem.Value())
		case "maxTimeMS":
			optsElems = bsoncore.AppendValueElement(optsElems, "$maxTimeMS", elem.Value())
		case "snapshot":
			optsElems = bsoncore.AppendValueElement(optsElems, "$snapshot", elem.Value())
		case "projection":
			returnFieldsSelector = elem.Value().Data
		case "skip":
			// CRUD spec declares skip as int64 but numToSkip is int32 in OP_QUERY
			numToSkip = int32(elem.Value().Int64())
		case "batchSize":
			batchSize = elem.Value().Int32()
			// Not possible to use batchSize = 1 because cursor will be closed on first batch
			if batchSize == 1 ***REMOVED***
				batchSize = 2
			***REMOVED***
		case "limit":
			// CRUD spec declares limit as int64 but numToReturn is int32 in OP_QUERY
			limit = int32(elem.Value().Int64())
		case "singleBatch":
			singleBatch = elem.Value().Boolean()
		case "tailable":
			flags |= wiremessage.TailableCursor
		case "awaitData":
			flags |= wiremessage.AwaitData
		case "oplogReplay":
			flags |= wiremessage.OplogReplay
		case "noCursorTimeout":
			flags |= wiremessage.NoCursorTimeout
		case "allowPartialResults":
			flags |= wiremessage.Partial
		***REMOVED***
	***REMOVED***

	// for non-legacy servers, a negative limit is implemented as a positive limit + singleBatch = true
	if singleBatch ***REMOVED***
		limit = limit * -1
	***REMOVED***
	numToReturn = op.calculateNumberToReturn(limit, batchSize)

	// add read preference if needed
	rp, err := op.createReadPref(desc, true)
	if err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***
	if len(rp) > 0 ***REMOVED***
		optsElems = bsoncore.AppendDocumentElement(optsElems, "$readPreference", rp)
	***REMOVED***

	if len(filter) == 0 ***REMOVED***
		var fidx int32
		fidx, filter = bsoncore.AppendDocumentStart(filter)
		filter, _ = bsoncore.AppendDocumentEnd(filter, fidx)
	***REMOVED***

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, flags)
	dst = wiremessage.AppendQueryFullCollectionName(dst, op.getFullCollectionName(collName))
	dst = wiremessage.AppendQueryNumberToSkip(dst, numToSkip)
	dst = wiremessage.AppendQueryNumberToReturn(dst, numToReturn)
	dst = op.appendLegacyQueryDocument(dst, filter, optsElems)
	if len(returnFieldsSelector) != 0 ***REMOVED***
		// returnFieldsSelector is optional
		dst = append(dst, returnFieldsSelector...)
	***REMOVED***

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, collName, nil
***REMOVED***

func (op Operation) calculateNumberToReturn(limit, batchSize int32) int32 ***REMOVED***
	var numToReturn int32

	if limit < 0 ***REMOVED***
		numToReturn = limit
	***REMOVED*** else if limit == 0 ***REMOVED***
		numToReturn = batchSize
	***REMOVED*** else if batchSize == 0 ***REMOVED***
		numToReturn = limit
	***REMOVED*** else if limit < batchSize ***REMOVED***
		numToReturn = limit
	***REMOVED*** else ***REMOVED***
		numToReturn = batchSize
	***REMOVED***

	return numToReturn
***REMOVED***

func (op Operation) legacyGetMore(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error ***REMOVED***
	wm, startedInfo, collName, err := op.createLegacyGetMoreWiremessage(dst, desc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	finishedInfo := finishedInformation***REMOVED***
		cmdName:   startedInfo.cmdName,
		requestID: startedInfo.requestID,
		startTime: time.Now(),
		connID:    startedInfo.connID,
	***REMOVED***
	finishedInfo.response, finishedInfo.cmdErr = op.roundTripLegacyCursor(ctx, wm, srvr, conn, collName, nextBatchIdentifier)
	op.publishFinishedEvent(ctx, finishedInfo)

	if finishedInfo.cmdErr != nil ***REMOVED***
		return finishedInfo.cmdErr
	***REMOVED***

	if op.ProcessResponseFn != nil ***REMOVED***
		// CurrentIndex is always 0 in this mode.
		info := ResponseInfo***REMOVED***
			ServerResponse:        finishedInfo.response,
			Server:                srvr,
			Connection:            conn,
			ConnectionDescription: desc.Server,
		***REMOVED***
		return op.ProcessResponseFn(info)
	***REMOVED***
	return nil
***REMOVED***

func (op Operation) createLegacyGetMoreWiremessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) ***REMOVED***
	info := startedInformation***REMOVED***
		requestID: wiremessage.NextRequestID(),
		cmdName:   "getMore",
	***REMOVED***

	var cmdDoc bsoncore.Document
	var cmdIdx int32
	var err error

	cmdIdx, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	cmdDoc, err = op.CommandFn(cmdDoc, desc)
	if err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIdx)
	info.cmd = cmdDoc

	cmdElems, err := cmdDoc.Elements()
	if err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***

	var cursorID int64
	var numToReturn int32
	var collName string
	for _, elem := range cmdElems ***REMOVED***
		switch elem.Key() ***REMOVED***
		case "getMore":
			cursorID = elem.Value().Int64()
		case "collection":
			collName = elem.Value().StringValue()
		case "batchSize":
			numToReturn = elem.Value().Int32()
		***REMOVED***
	***REMOVED***

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpGetMore)
	dst = wiremessage.AppendGetMoreZero(dst)
	dst = wiremessage.AppendGetMoreFullCollectionName(dst, op.getFullCollectionName(collName))
	dst = wiremessage.AppendGetMoreNumberToReturn(dst, numToReturn)
	dst = wiremessage.AppendGetMoreCursorID(dst, cursorID)

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, collName, nil
***REMOVED***

func (op Operation) legacyKillCursors(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error ***REMOVED***
	wm, startedInfo, _, err := op.createLegacyKillCursorsWiremessage(dst, desc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	// skip startTime because OP_KILL_CURSORS does not return a response
	finishedInfo := finishedInformation***REMOVED***
		cmdName:   "killCursors",
		requestID: startedInfo.requestID,
		connID:    startedInfo.connID,
	***REMOVED***

	err = conn.WriteWireMessage(ctx, wm)
	if err != nil ***REMOVED***
		err = Error***REMOVED***Message: err.Error(), Labels: []string***REMOVED***TransientTransactionError, NetworkError***REMOVED******REMOVED***
		if ep, ok := srvr.(ErrorProcessor); ok ***REMOVED***
			_ = ep.ProcessError(err, conn)
		***REMOVED***

		finishedInfo.cmdErr = err
		op.publishFinishedEvent(ctx, finishedInfo)
		return err
	***REMOVED***

	ridx, response := bsoncore.AppendDocumentStart(nil)
	response = bsoncore.AppendInt32Element(response, "ok", 1)
	response = bsoncore.AppendArrayElement(response, "cursorsUnknown", startedInfo.cmd.Lookup("cursors").Array())
	response, _ = bsoncore.AppendDocumentEnd(response, ridx)

	finishedInfo.response = response
	op.publishFinishedEvent(ctx, finishedInfo)
	return nil
***REMOVED***

func (op Operation) createLegacyKillCursorsWiremessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) ***REMOVED***
	info := startedInformation***REMOVED***
		cmdName:   "killCursors",
		requestID: wiremessage.NextRequestID(),
	***REMOVED***

	var cmdDoc bsoncore.Document
	var cmdIdx int32
	var err error

	cmdIdx, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	cmdDoc, err = op.CommandFn(cmdDoc, desc)
	if err != nil ***REMOVED***
		return nil, info, "", err
	***REMOVED***
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIdx)
	info.cmd = cmdDoc

	cmdElems, err := cmdDoc.Elements()
	if err != nil ***REMOVED***
		return nil, info, "", err
	***REMOVED***

	var collName string
	var cursors bsoncore.Array
	for _, elem := range cmdElems ***REMOVED***
		switch elem.Key() ***REMOVED***
		case "killCursors":
			collName = elem.Value().StringValue()
		case "cursors":
			cursors = elem.Value().Array()
		***REMOVED***
	***REMOVED***

	var cursorIDs []int64
	if cursors != nil ***REMOVED***
		cursorValues, err := cursors.Values()
		if err != nil ***REMOVED***
			return nil, info, "", err
		***REMOVED***

		for _, cursorVal := range cursorValues ***REMOVED***
			cursorIDs = append(cursorIDs, cursorVal.Int64())
		***REMOVED***
	***REMOVED***

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpKillCursors)
	dst = wiremessage.AppendKillCursorsZero(dst)
	dst = wiremessage.AppendKillCursorsNumberIDs(dst, int32(len(cursorIDs)))
	dst = wiremessage.AppendKillCursorsCursorIDs(dst, cursorIDs)

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, collName, nil
***REMOVED***

func (op Operation) legacyListCollections(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error ***REMOVED***
	wm, startedInfo, collName, err := op.createLegacyListCollectionsWiremessage(dst, desc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	finishedInfo := finishedInformation***REMOVED***
		cmdName:   startedInfo.cmdName,
		requestID: startedInfo.requestID,
		startTime: time.Now(),
		connID:    startedInfo.connID,
	***REMOVED***

	finishedInfo.response, finishedInfo.cmdErr = op.roundTripLegacyCursor(ctx, wm, srvr, conn, collName, firstBatchIdentifier)
	op.publishFinishedEvent(ctx, finishedInfo)

	if finishedInfo.cmdErr != nil ***REMOVED***
		return finishedInfo.cmdErr
	***REMOVED***

	if op.ProcessResponseFn != nil ***REMOVED***
		// CurrentIndex is always 0 in this mode.
		info := ResponseInfo***REMOVED***
			ServerResponse:        finishedInfo.response,
			Server:                srvr,
			Connection:            conn,
			ConnectionDescription: desc.Server,
		***REMOVED***
		return op.ProcessResponseFn(info)
	***REMOVED***
	return nil
***REMOVED***

func (op Operation) createLegacyListCollectionsWiremessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) ***REMOVED***
	info := startedInformation***REMOVED***
		cmdName:   "find",
		requestID: wiremessage.NextRequestID(),
	***REMOVED***

	var cmdDoc bsoncore.Document
	var cmdIdx int32
	var err error

	cmdIdx, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	if cmdDoc, err = op.CommandFn(cmdDoc, desc); err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIdx)
	info.cmd, err = op.convertCommandToFind(cmdDoc, listCollectionsNamespace)
	if err != nil ***REMOVED***
		return nil, info, "", err
	***REMOVED***

	// lookup filter directly instead of calling cmdDoc.Elements() because the only listCollections option is nameOnly,
	// which doesn't apply to legacy servers
	var originalFilter bsoncore.Document
	if filterVal, err := cmdDoc.LookupErr("filter"); err == nil ***REMOVED***
		originalFilter = filterVal.Document()
	***REMOVED***

	var optsElems []byte
	filter, err := op.transformListCollectionsFilter(originalFilter)
	if err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***
	rp, err := op.createReadPref(desc, true)
	if err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***
	if len(rp) > 0 ***REMOVED***
		optsElems = bsoncore.AppendDocumentElement(optsElems, "$readPreference", rp)
	***REMOVED***

	var batchSize int32
	if val, ok := cmdDoc.Lookup("cursor", "batchSize").AsInt32OK(); ok ***REMOVED***
		batchSize = val
	***REMOVED***

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, op.secondaryOK(desc))
	dst = wiremessage.AppendQueryFullCollectionName(dst, op.getFullCollectionName(listCollectionsNamespace))
	dst = wiremessage.AppendQueryNumberToSkip(dst, 0)
	dst = wiremessage.AppendQueryNumberToReturn(dst, batchSize)
	dst = op.appendLegacyQueryDocument(dst, filter, optsElems)
	// leave out returnFieldsSelector because it is optional

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, listCollectionsNamespace, nil
***REMOVED***

func (op Operation) transformListCollectionsFilter(filter bsoncore.Document) (bsoncore.Document, error) ***REMOVED***
	// filter out results containing $ because those represent indexes
	var regexFilter bsoncore.Document
	var ridx int32
	ridx, regexFilter = bsoncore.AppendDocumentStart(regexFilter)
	regexFilter = bsoncore.AppendRegexElement(regexFilter, "name", "^[^$]*$", "")
	regexFilter, _ = bsoncore.AppendDocumentEnd(regexFilter, ridx)

	if len(filter) == 0 ***REMOVED***
		return regexFilter, nil
	***REMOVED***

	convertedIdx, convertedFilter := bsoncore.AppendDocumentStart(nil)
	elems, err := filter.Elements()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, elem := range elems ***REMOVED***
		if elem.Key() != "name" ***REMOVED***
			convertedFilter = append(convertedFilter, elem...)
			continue
		***REMOVED***

		// the name value in a filter for legacy list collections must be a string and has to be prepended
		// with the database name
		nameVal := elem.Value()
		if nameVal.Type != bsontype.String ***REMOVED***
			return nil, ErrFilterType
		***REMOVED***
		convertedFilter = bsoncore.AppendStringElement(convertedFilter, "name", op.getFullCollectionName(nameVal.StringValue()))
	***REMOVED***
	convertedFilter, _ = bsoncore.AppendDocumentEnd(convertedFilter, convertedIdx)

	// combine regexFilter and convertedFilter with $and
	var combinedFilter bsoncore.Document
	var cidx, aidx int32
	cidx, combinedFilter = bsoncore.AppendDocumentStart(combinedFilter)
	aidx, combinedFilter = bsoncore.AppendArrayElementStart(combinedFilter, "$and")
	combinedFilter = bsoncore.AppendDocumentElement(combinedFilter, "0", regexFilter)
	combinedFilter = bsoncore.AppendDocumentElement(combinedFilter, "1", convertedFilter)
	combinedFilter, _ = bsoncore.AppendArrayEnd(combinedFilter, aidx)
	combinedFilter, _ = bsoncore.AppendDocumentEnd(combinedFilter, cidx)

	return combinedFilter, nil
***REMOVED***

func (op Operation) legacyListIndexes(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error ***REMOVED***
	wm, startedInfo, collName, err := op.createLegacyListIndexesWiremessage(dst, desc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	finishedInfo := finishedInformation***REMOVED***
		cmdName:   startedInfo.cmdName,
		requestID: startedInfo.requestID,
		startTime: time.Now(),
		connID:    startedInfo.connID,
	***REMOVED***

	finishedInfo.response, finishedInfo.cmdErr = op.roundTripLegacyCursor(ctx, wm, srvr, conn, collName, firstBatchIdentifier)
	op.publishFinishedEvent(ctx, finishedInfo)

	if finishedInfo.cmdErr != nil ***REMOVED***
		return finishedInfo.cmdErr
	***REMOVED***

	if op.ProcessResponseFn != nil ***REMOVED***
		// CurrentIndex is always 0 in this mode.
		info := ResponseInfo***REMOVED***
			ServerResponse:        finishedInfo.response,
			Server:                srvr,
			Connection:            conn,
			ConnectionDescription: desc.Server,
		***REMOVED***
		return op.ProcessResponseFn(info)
	***REMOVED***
	return nil
***REMOVED***

func (op Operation) createLegacyListIndexesWiremessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) ***REMOVED***
	info := startedInformation***REMOVED***
		cmdName:   "find",
		requestID: wiremessage.NextRequestID(),
	***REMOVED***

	var cmdDoc bsoncore.Document
	var cmdIndex int32
	var err error

	cmdIndex, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	cmdDoc, err = op.CommandFn(cmdDoc, desc)
	if err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIndex)
	info.cmd, err = op.convertCommandToFind(cmdDoc, listIndexesNamespace)
	if err != nil ***REMOVED***
		return nil, info, "", err
	***REMOVED***

	cmdElems, err := cmdDoc.Elements()
	if err != nil ***REMOVED***
		return nil, info, "", err
	***REMOVED***

	var filterCollName string
	var batchSize int32
	var optsElems []byte // options elements
	for _, elem := range cmdElems ***REMOVED***
		switch elem.Key() ***REMOVED***
		case "listIndexes":
			filterCollName = elem.Value().StringValue()
		case "cursor":
			// the batchSize option is embedded in a cursor subdocument
			cursorDoc := elem.Value().Document()
			if val, err := cursorDoc.LookupErr("batchSize"); err == nil ***REMOVED***
				batchSize = val.Int32()
			***REMOVED***
		case "maxTimeMS":
			optsElems = bsoncore.AppendValueElement(optsElems, "$maxTimeMS", elem.Value())
		***REMOVED***
	***REMOVED***

	// always filter with ***REMOVED***ns: db.collection***REMOVED***
	fidx, filter := bsoncore.AppendDocumentStart(nil)
	filter = bsoncore.AppendStringElement(filter, "ns", op.getFullCollectionName(filterCollName))
	filter, _ = bsoncore.AppendDocumentEnd(filter, fidx)

	rp, err := op.createReadPref(desc, true)
	if err != nil ***REMOVED***
		return dst, info, "", err
	***REMOVED***
	if len(rp) > 0 ***REMOVED***
		optsElems = bsoncore.AppendDocumentElement(optsElems, "$readPreference", rp)
	***REMOVED***

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, op.secondaryOK(desc))
	dst = wiremessage.AppendQueryFullCollectionName(dst, op.getFullCollectionName(listIndexesNamespace))
	dst = wiremessage.AppendQueryNumberToSkip(dst, 0)
	dst = wiremessage.AppendQueryNumberToReturn(dst, batchSize)
	dst = op.appendLegacyQueryDocument(dst, filter, optsElems)
	// leave out returnFieldsSelector because it is optional

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, listIndexesNamespace, nil
***REMOVED***

// convertCommandToFind takes a non-legacy command document for a command that needs to be run as a find on legacy
// servers and converts it to a find command document for APM.
func (op Operation) convertCommandToFind(cmdDoc bsoncore.Document, collName string) (bsoncore.Document, error) ***REMOVED***
	cidx, converted := bsoncore.AppendDocumentStart(nil)
	elems, err := cmdDoc.Elements()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	converted = bsoncore.AppendStringElement(converted, "find", collName)
	// skip the first element because that will have the old command name
	for i := 1; i < len(elems); i++ ***REMOVED***
		converted = bsoncore.AppendValueElement(converted, elems[i].Key(), elems[i].Value())
	***REMOVED***

	converted, _ = bsoncore.AppendDocumentEnd(converted, cidx)
	return converted, nil
***REMOVED***

// appendLegacyQueryDocument takes a filter and a list of options elements for a legacy find operation, creates
// a query document, and appends it to dst.
func (op Operation) appendLegacyQueryDocument(dst []byte, filter bsoncore.Document, opts []byte) []byte ***REMOVED***
	if len(opts) == 0 ***REMOVED***
		dst = append(dst, filter...)
		return dst
	***REMOVED***

	// filter must be wrapped in $query if other $-modifiers are used
	var qidx int32
	qidx, dst = bsoncore.AppendDocumentStart(dst)
	dst = bsoncore.AppendDocumentElement(dst, "$query", filter)
	dst = append(dst, opts...)
	dst, _ = bsoncore.AppendDocumentEnd(dst, qidx)
	return dst
***REMOVED***

// roundTripLegacyCursor sends a wiremessage for an operation expecting a cursor result and converts the legacy
// document sequence into a cursor document.
func (op Operation) roundTripLegacyCursor(ctx context.Context, wm []byte, srvr Server, conn Connection, collName, identifier string) (bsoncore.Document, error) ***REMOVED***
	wm, err := op.roundTripLegacy(ctx, conn, wm)
	if ep, ok := srvr.(ErrorProcessor); ok ***REMOVED***
		_ = ep.ProcessError(err, conn)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return op.upconvertCursorResponse(wm, identifier, collName)
***REMOVED***

// roundTripLegacy handles writing a wire message and reading the response.
func (op Operation) roundTripLegacy(ctx context.Context, conn Connection, wm []byte) ([]byte, error) ***REMOVED***
	err := conn.WriteWireMessage(ctx, wm)
	if err != nil ***REMOVED***
		return nil, Error***REMOVED***Message: err.Error(), Labels: []string***REMOVED***TransientTransactionError, NetworkError***REMOVED***, Wrapped: err***REMOVED***
	***REMOVED***

	wm, err = conn.ReadWireMessage(ctx, wm[:0])
	if err != nil ***REMOVED***
		err = Error***REMOVED***Message: err.Error(), Labels: []string***REMOVED***TransientTransactionError, NetworkError***REMOVED***, Wrapped: err***REMOVED***
	***REMOVED***
	return wm, err
***REMOVED***

func (op Operation) upconvertCursorResponse(wm []byte, batchIdentifier string, collName string) (bsoncore.Document, error) ***REMOVED***
	reply := op.decodeOpReply(wm, true)
	if reply.err != nil ***REMOVED***
		return nil, reply.err
	***REMOVED***

	cursorIdx, cursorDoc := bsoncore.AppendDocumentStart(nil)
	// convert reply documents to BSON array
	var arrIdx int32
	arrIdx, cursorDoc = bsoncore.AppendArrayElementStart(cursorDoc, batchIdentifier)
	for i, doc := range reply.documents ***REMOVED***
		cursorDoc = bsoncore.AppendDocumentElement(cursorDoc, strconv.Itoa(i), doc)
	***REMOVED***
	cursorDoc, _ = bsoncore.AppendArrayEnd(cursorDoc, arrIdx)

	cursorDoc = bsoncore.AppendInt64Element(cursorDoc, "id", reply.cursorID)
	cursorDoc = bsoncore.AppendStringElement(cursorDoc, "ns", op.getFullCollectionName(collName))
	cursorDoc, _ = bsoncore.AppendDocumentEnd(cursorDoc, cursorIdx)

	resIdx, resDoc := bsoncore.AppendDocumentStart(nil)
	resDoc = bsoncore.AppendInt32Element(resDoc, "ok", 1)
	resDoc = bsoncore.AppendDocumentElement(resDoc, "cursor", cursorDoc)
	resDoc, _ = bsoncore.AppendDocumentEnd(resDoc, resIdx)

	return resDoc, nil
***REMOVED***
