// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package wiremessage

import (
	"bytes"
	"strings"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// WireMessage represents a MongoDB wire message in binary form.
type WireMessage []byte

var globalRequestID int32

// CurrentRequestID returns the current request ID.
func CurrentRequestID() int32 ***REMOVED*** return atomic.LoadInt32(&globalRequestID) ***REMOVED***

// NextRequestID returns the next request ID.
func NextRequestID() int32 ***REMOVED*** return atomic.AddInt32(&globalRequestID, 1) ***REMOVED***

// OpCode represents a MongoDB wire protocol opcode.
type OpCode int32

// These constants are the valid opcodes for the version of the wireprotocol
// supported by this library. The skipped OpCodes are historical OpCodes that
// are no longer used.
const (
	OpReply        OpCode = 1
	_              OpCode = 1001
	OpUpdate       OpCode = 2001
	OpInsert       OpCode = 2002
	_              OpCode = 2003
	OpQuery        OpCode = 2004
	OpGetMore      OpCode = 2005
	OpDelete       OpCode = 2006
	OpKillCursors  OpCode = 2007
	OpCommand      OpCode = 2010
	OpCommandReply OpCode = 2011
	OpCompressed   OpCode = 2012
	OpMsg          OpCode = 2013
)

// String implements the fmt.Stringer interface.
func (oc OpCode) String() string ***REMOVED***
	switch oc ***REMOVED***
	case OpReply:
		return "OP_REPLY"
	case OpUpdate:
		return "OP_UPDATE"
	case OpInsert:
		return "OP_INSERT"
	case OpQuery:
		return "OP_QUERY"
	case OpGetMore:
		return "OP_GET_MORE"
	case OpDelete:
		return "OP_DELETE"
	case OpKillCursors:
		return "OP_KILL_CURSORS"
	case OpCommand:
		return "OP_COMMAND"
	case OpCommandReply:
		return "OP_COMMANDREPLY"
	case OpCompressed:
		return "OP_COMPRESSED"
	case OpMsg:
		return "OP_MSG"
	default:
		return "<invalid opcode>"
	***REMOVED***
***REMOVED***

// QueryFlag represents the flags on an OP_QUERY message.
type QueryFlag int32

// These constants represent the individual flags on an OP_QUERY message.
const (
	_ QueryFlag = 1 << iota
	TailableCursor
	SecondaryOK
	OplogReplay
	NoCursorTimeout
	AwaitData
	Exhaust
	Partial
)

// String implements the fmt.Stringer interface.
func (qf QueryFlag) String() string ***REMOVED***
	strs := make([]string, 0)
	if qf&TailableCursor == TailableCursor ***REMOVED***
		strs = append(strs, "TailableCursor")
	***REMOVED***
	if qf&SecondaryOK == SecondaryOK ***REMOVED***
		strs = append(strs, "SecondaryOK")
	***REMOVED***
	if qf&OplogReplay == OplogReplay ***REMOVED***
		strs = append(strs, "OplogReplay")
	***REMOVED***
	if qf&NoCursorTimeout == NoCursorTimeout ***REMOVED***
		strs = append(strs, "NoCursorTimeout")
	***REMOVED***
	if qf&AwaitData == AwaitData ***REMOVED***
		strs = append(strs, "AwaitData")
	***REMOVED***
	if qf&Exhaust == Exhaust ***REMOVED***
		strs = append(strs, "Exhaust")
	***REMOVED***
	if qf&Partial == Partial ***REMOVED***
		strs = append(strs, "Partial")
	***REMOVED***
	str := "["
	str += strings.Join(strs, ", ")
	str += "]"
	return str
***REMOVED***

// MsgFlag represents the flags on an OP_MSG message.
type MsgFlag uint32

// These constants represent the individual flags on an OP_MSG message.
const (
	ChecksumPresent MsgFlag = 1 << iota
	MoreToCome

	ExhaustAllowed MsgFlag = 1 << 16
)

// ReplyFlag represents the flags of an OP_REPLY message.
type ReplyFlag int32

// These constants represent the individual flags of an OP_REPLY message.
const (
	CursorNotFound ReplyFlag = 1 << iota
	QueryFailure
	ShardConfigStale
	AwaitCapable
)

// String implements the fmt.Stringer interface.
func (rf ReplyFlag) String() string ***REMOVED***
	strs := make([]string, 0)
	if rf&CursorNotFound == CursorNotFound ***REMOVED***
		strs = append(strs, "CursorNotFound")
	***REMOVED***
	if rf&QueryFailure == QueryFailure ***REMOVED***
		strs = append(strs, "QueryFailure")
	***REMOVED***
	if rf&ShardConfigStale == ShardConfigStale ***REMOVED***
		strs = append(strs, "ShardConfigStale")
	***REMOVED***
	if rf&AwaitCapable == AwaitCapable ***REMOVED***
		strs = append(strs, "AwaitCapable")
	***REMOVED***
	str := "["
	str += strings.Join(strs, ", ")
	str += "]"
	return str
***REMOVED***

// SectionType represents the type for 1 section in an OP_MSG
type SectionType uint8

// These constants represent the individual section types for a section in an OP_MSG
const (
	SingleDocument SectionType = iota
	DocumentSequence
)

// OpmsgWireVersion is the minimum wire version needed to use OP_MSG
const OpmsgWireVersion = 6

// CompressorID is the ID for each type of Compressor.
type CompressorID uint8

// These constants represent the individual compressor IDs for an OP_COMPRESSED.
const (
	CompressorNoOp CompressorID = iota
	CompressorSnappy
	CompressorZLib
	CompressorZstd
)

// String implements the fmt.Stringer interface.
func (id CompressorID) String() string ***REMOVED***
	switch id ***REMOVED***
	case CompressorNoOp:
		return "CompressorNoOp"
	case CompressorSnappy:
		return "CompressorSnappy"
	case CompressorZLib:
		return "CompressorZLib"
	case CompressorZstd:
		return "CompressorZstd"
	default:
		return "CompressorInvalid"
	***REMOVED***
***REMOVED***

const (
	// DefaultZlibLevel is the default level for zlib compression
	DefaultZlibLevel = 6
	// DefaultZstdLevel is the default level for zstd compression.
	// Matches https://github.com/wiredtiger/wiredtiger/blob/f08bc4b18612ef95a39b12166abcccf207f91596/ext/compressors/zstd/zstd_compress.c#L299
	DefaultZstdLevel = 6
)

// AppendHeaderStart appends a header to the dst slice and returns an index where the wire message
// starts in dst and the updated slice.
func AppendHeaderStart(dst []byte, reqid, respto int32, opcode OpCode) (index int32, b []byte) ***REMOVED***
	index, dst = bsoncore.ReserveLength(dst)
	dst = appendi32(dst, reqid)
	dst = appendi32(dst, respto)
	dst = appendi32(dst, int32(opcode))
	return index, dst
***REMOVED***

// AppendHeader appends a header to dst.
func AppendHeader(dst []byte, length, reqid, respto int32, opcode OpCode) []byte ***REMOVED***
	dst = appendi32(dst, length)
	dst = appendi32(dst, reqid)
	dst = appendi32(dst, respto)
	dst = appendi32(dst, int32(opcode))
	return dst
***REMOVED***

// ReadHeader reads a wire message header from src.
func ReadHeader(src []byte) (length, requestID, responseTo int32, opcode OpCode, rem []byte, ok bool) ***REMOVED***
	if len(src) < 16 ***REMOVED***
		return 0, 0, 0, 0, src, false
	***REMOVED***
	length = (int32(src[0]) | int32(src[1])<<8 | int32(src[2])<<16 | int32(src[3])<<24)
	requestID = (int32(src[4]) | int32(src[5])<<8 | int32(src[6])<<16 | int32(src[7])<<24)
	responseTo = (int32(src[8]) | int32(src[9])<<8 | int32(src[10])<<16 | int32(src[11])<<24)
	opcode = OpCode(int32(src[12]) | int32(src[13])<<8 | int32(src[14])<<16 | int32(src[15])<<24)
	return length, requestID, responseTo, opcode, src[16:], true
***REMOVED***

// AppendQueryFlags appends the flags for an OP_QUERY wire message.
func AppendQueryFlags(dst []byte, flags QueryFlag) []byte ***REMOVED***
	return appendi32(dst, int32(flags))
***REMOVED***

// AppendMsgFlags appends the flags for an OP_MSG wire message.
func AppendMsgFlags(dst []byte, flags MsgFlag) []byte ***REMOVED***
	return appendi32(dst, int32(flags))
***REMOVED***

// AppendReplyFlags appends the flags for an OP_REPLY wire message.
func AppendReplyFlags(dst []byte, flags ReplyFlag) []byte ***REMOVED***
	return appendi32(dst, int32(flags))
***REMOVED***

// AppendMsgSectionType appends the section type to dst.
func AppendMsgSectionType(dst []byte, stype SectionType) []byte ***REMOVED***
	return append(dst, byte(stype))
***REMOVED***

// AppendQueryFullCollectionName appends the full collection name to dst.
func AppendQueryFullCollectionName(dst []byte, ns string) []byte ***REMOVED***
	return appendCString(dst, ns)
***REMOVED***

// AppendQueryNumberToSkip appends the number to skip to dst.
func AppendQueryNumberToSkip(dst []byte, skip int32) []byte ***REMOVED***
	return appendi32(dst, skip)
***REMOVED***

// AppendQueryNumberToReturn appends the number to return to dst.
func AppendQueryNumberToReturn(dst []byte, nor int32) []byte ***REMOVED***
	return appendi32(dst, nor)
***REMOVED***

// AppendReplyCursorID appends the cursor ID to dst.
func AppendReplyCursorID(dst []byte, id int64) []byte ***REMOVED***
	return appendi64(dst, id)
***REMOVED***

// AppendReplyStartingFrom appends the starting from field to dst.
func AppendReplyStartingFrom(dst []byte, sf int32) []byte ***REMOVED***
	return appendi32(dst, sf)
***REMOVED***

// AppendReplyNumberReturned appends the number returned to dst.
func AppendReplyNumberReturned(dst []byte, nr int32) []byte ***REMOVED***
	return appendi32(dst, nr)
***REMOVED***

// AppendCompressedOriginalOpCode appends the original opcode to dst.
func AppendCompressedOriginalOpCode(dst []byte, opcode OpCode) []byte ***REMOVED***
	return appendi32(dst, int32(opcode))
***REMOVED***

// AppendCompressedUncompressedSize appends the uncompressed size of a
// compressed wiremessage to dst.
func AppendCompressedUncompressedSize(dst []byte, size int32) []byte ***REMOVED*** return appendi32(dst, size) ***REMOVED***

// AppendCompressedCompressorID appends the ID of the compressor to dst.
func AppendCompressedCompressorID(dst []byte, id CompressorID) []byte ***REMOVED***
	return append(dst, byte(id))
***REMOVED***

// AppendCompressedCompressedMessage appends the compressed wiremessage to dst.
func AppendCompressedCompressedMessage(dst []byte, msg []byte) []byte ***REMOVED*** return append(dst, msg...) ***REMOVED***

// AppendGetMoreZero appends the zero field to dst.
func AppendGetMoreZero(dst []byte) []byte ***REMOVED***
	return appendi32(dst, 0)
***REMOVED***

// AppendGetMoreFullCollectionName appends the fullCollectionName field to dst.
func AppendGetMoreFullCollectionName(dst []byte, ns string) []byte ***REMOVED***
	return appendCString(dst, ns)
***REMOVED***

// AppendGetMoreNumberToReturn appends the numberToReturn field to dst.
func AppendGetMoreNumberToReturn(dst []byte, numToReturn int32) []byte ***REMOVED***
	return appendi32(dst, numToReturn)
***REMOVED***

// AppendGetMoreCursorID appends the cursorID field to dst.
func AppendGetMoreCursorID(dst []byte, cursorID int64) []byte ***REMOVED***
	return appendi64(dst, cursorID)
***REMOVED***

// AppendKillCursorsZero appends the zero field to dst.
func AppendKillCursorsZero(dst []byte) []byte ***REMOVED***
	return appendi32(dst, 0)
***REMOVED***

// AppendKillCursorsNumberIDs appends the numberOfCursorIDs field to dst.
func AppendKillCursorsNumberIDs(dst []byte, numIDs int32) []byte ***REMOVED***
	return appendi32(dst, numIDs)
***REMOVED***

// AppendKillCursorsCursorIDs appends each the cursorIDs field to dst.
func AppendKillCursorsCursorIDs(dst []byte, cursors []int64) []byte ***REMOVED***
	for _, cursor := range cursors ***REMOVED***
		dst = appendi64(dst, cursor)
	***REMOVED***
	return dst
***REMOVED***

// ReadMsgFlags reads the OP_MSG flags from src.
func ReadMsgFlags(src []byte) (flags MsgFlag, rem []byte, ok bool) ***REMOVED***
	i32, rem, ok := readi32(src)
	return MsgFlag(i32), rem, ok
***REMOVED***

// IsMsgMoreToCome returns if the provided wire message is an OP_MSG with the more to come flag set.
func IsMsgMoreToCome(wm []byte) bool ***REMOVED***
	return len(wm) >= 20 &&
		OpCode(readi32unsafe(wm[12:16])) == OpMsg &&
		MsgFlag(readi32unsafe(wm[16:20]))&MoreToCome == MoreToCome
***REMOVED***

// ReadMsgSectionType reads the section type from src.
func ReadMsgSectionType(src []byte) (stype SectionType, rem []byte, ok bool) ***REMOVED***
	if len(src) < 1 ***REMOVED***
		return 0, src, false
	***REMOVED***
	return SectionType(src[0]), src[1:], true
***REMOVED***

// ReadMsgSectionSingleDocument reads a single document from src.
func ReadMsgSectionSingleDocument(src []byte) (doc bsoncore.Document, rem []byte, ok bool) ***REMOVED***
	return bsoncore.ReadDocument(src)
***REMOVED***

// ReadMsgSectionDocumentSequence reads an identifier and document sequence from src and returns the document sequence
// data parsed into a slice of BSON documents.
func ReadMsgSectionDocumentSequence(src []byte) (identifier string, docs []bsoncore.Document, rem []byte, ok bool) ***REMOVED***
	length, rem, ok := readi32(src)
	if !ok || int(length) > len(src) ***REMOVED***
		return "", nil, rem, false
	***REMOVED***

	rem, ret := rem[:length-4], rem[length-4:] // reslice so we can just iterate a loop later

	identifier, rem, ok = readcstring(rem)
	if !ok ***REMOVED***
		return "", nil, rem, false
	***REMOVED***

	docs = make([]bsoncore.Document, 0)
	var doc bsoncore.Document
	for ***REMOVED***
		doc, rem, ok = bsoncore.ReadDocument(rem)
		if !ok ***REMOVED***
			break
		***REMOVED***
		docs = append(docs, doc)
	***REMOVED***
	if len(rem) > 0 ***REMOVED***
		return "", nil, append(rem, ret...), false
	***REMOVED***

	return identifier, docs, ret, true
***REMOVED***

// ReadMsgSectionRawDocumentSequence reads an identifier and document sequence from src and returns the raw document
// sequence data.
func ReadMsgSectionRawDocumentSequence(src []byte) (identifier string, data []byte, rem []byte, ok bool) ***REMOVED***
	length, rem, ok := readi32(src)
	if !ok || int(length) > len(src) ***REMOVED***
		return "", nil, rem, false
	***REMOVED***

	// After these assignments, rem will be the data containing the identifier string + the document sequence bytes and
	// rest will be the rest of the wire message after this document sequence.
	rem, rest := rem[:length-4], rem[length-4:]

	identifier, rem, ok = readcstring(rem)
	if !ok ***REMOVED***
		return "", nil, rem, false
	***REMOVED***

	return identifier, rem, rest, true
***REMOVED***

// ReadMsgChecksum reads a checksum from src.
func ReadMsgChecksum(src []byte) (checksum uint32, rem []byte, ok bool) ***REMOVED***
	i32, rem, ok := readi32(src)
	return uint32(i32), rem, ok
***REMOVED***

// ReadQueryFlags reads OP_QUERY flags from src.
func ReadQueryFlags(src []byte) (flags QueryFlag, rem []byte, ok bool) ***REMOVED***
	i32, rem, ok := readi32(src)
	return QueryFlag(i32), rem, ok
***REMOVED***

// ReadQueryFullCollectionName reads the full collection name from src.
func ReadQueryFullCollectionName(src []byte) (collname string, rem []byte, ok bool) ***REMOVED***
	return readcstring(src)
***REMOVED***

// ReadQueryNumberToSkip reads the number to skip from src.
func ReadQueryNumberToSkip(src []byte) (nts int32, rem []byte, ok bool) ***REMOVED***
	return readi32(src)
***REMOVED***

// ReadQueryNumberToReturn reads the number to return from src.
func ReadQueryNumberToReturn(src []byte) (ntr int32, rem []byte, ok bool) ***REMOVED***
	return readi32(src)
***REMOVED***

// ReadQueryQuery reads the query from src.
func ReadQueryQuery(src []byte) (query bsoncore.Document, rem []byte, ok bool) ***REMOVED***
	return bsoncore.ReadDocument(src)
***REMOVED***

// ReadQueryReturnFieldsSelector reads a return fields selector document from src.
func ReadQueryReturnFieldsSelector(src []byte) (rfs bsoncore.Document, rem []byte, ok bool) ***REMOVED***
	return bsoncore.ReadDocument(src)
***REMOVED***

// ReadReplyFlags reads OP_REPLY flags from src.
func ReadReplyFlags(src []byte) (flags ReplyFlag, rem []byte, ok bool) ***REMOVED***
	i32, rem, ok := readi32(src)
	return ReplyFlag(i32), rem, ok
***REMOVED***

// ReadReplyCursorID reads a cursor ID from src.
func ReadReplyCursorID(src []byte) (cursorID int64, rem []byte, ok bool) ***REMOVED***
	return readi64(src)
***REMOVED***

// ReadReplyStartingFrom reads the starting from from src.
func ReadReplyStartingFrom(src []byte) (startingFrom int32, rem []byte, ok bool) ***REMOVED***
	return readi32(src)
***REMOVED***

// ReadReplyNumberReturned reads the numbered returned from src.
func ReadReplyNumberReturned(src []byte) (numberReturned int32, rem []byte, ok bool) ***REMOVED***
	return readi32(src)
***REMOVED***

// ReadReplyDocuments reads as many documents as possible from src
func ReadReplyDocuments(src []byte) (docs []bsoncore.Document, rem []byte, ok bool) ***REMOVED***
	rem = src
	for ***REMOVED***
		var doc bsoncore.Document
		doc, rem, ok = bsoncore.ReadDocument(rem)
		if !ok ***REMOVED***
			break
		***REMOVED***

		docs = append(docs, doc)
	***REMOVED***

	return docs, rem, true
***REMOVED***

// ReadReplyDocument reads a reply document from src.
func ReadReplyDocument(src []byte) (doc bsoncore.Document, rem []byte, ok bool) ***REMOVED***
	return bsoncore.ReadDocument(src)
***REMOVED***

// ReadCompressedOriginalOpCode reads the original opcode from src.
func ReadCompressedOriginalOpCode(src []byte) (opcode OpCode, rem []byte, ok bool) ***REMOVED***
	i32, rem, ok := readi32(src)
	return OpCode(i32), rem, ok
***REMOVED***

// ReadCompressedUncompressedSize reads the uncompressed size of a
// compressed wiremessage to dst.
func ReadCompressedUncompressedSize(src []byte) (size int32, rem []byte, ok bool) ***REMOVED***
	return readi32(src)
***REMOVED***

// ReadCompressedCompressorID reads the ID of the compressor to dst.
func ReadCompressedCompressorID(src []byte) (id CompressorID, rem []byte, ok bool) ***REMOVED***
	if len(src) < 1 ***REMOVED***
		return 0, src, false
	***REMOVED***
	return CompressorID(src[0]), src[1:], true
***REMOVED***

// ReadCompressedCompressedMessage reads the compressed wiremessage to dst.
func ReadCompressedCompressedMessage(src []byte, length int32) (msg []byte, rem []byte, ok bool) ***REMOVED***
	if len(src) < int(length) ***REMOVED***
		return nil, src, false
	***REMOVED***
	return src[:length], src[length:], true
***REMOVED***

// ReadKillCursorsZero reads the zero field from src.
func ReadKillCursorsZero(src []byte) (zero int32, rem []byte, ok bool) ***REMOVED***
	return readi32(src)
***REMOVED***

// ReadKillCursorsNumberIDs reads the numberOfCursorIDs field from src.
func ReadKillCursorsNumberIDs(src []byte) (numIDs int32, rem []byte, ok bool) ***REMOVED***
	return readi32(src)
***REMOVED***

// ReadKillCursorsCursorIDs reads numIDs cursor IDs from src.
func ReadKillCursorsCursorIDs(src []byte, numIDs int32) (cursorIDs []int64, rem []byte, ok bool) ***REMOVED***
	var i int32
	var id int64
	for i = 0; i < numIDs; i++ ***REMOVED***
		id, src, ok = readi64(src)
		if !ok ***REMOVED***
			return cursorIDs, src, false
		***REMOVED***

		cursorIDs = append(cursorIDs, id)
	***REMOVED***
	return cursorIDs, src, true
***REMOVED***

func appendi32(dst []byte, i32 int32) []byte ***REMOVED***
	return append(dst, byte(i32), byte(i32>>8), byte(i32>>16), byte(i32>>24))
***REMOVED***

func appendi64(b []byte, i int64) []byte ***REMOVED***
	return append(b, byte(i), byte(i>>8), byte(i>>16), byte(i>>24), byte(i>>32), byte(i>>40), byte(i>>48), byte(i>>56))
***REMOVED***

func appendCString(b []byte, str string) []byte ***REMOVED***
	b = append(b, str...)
	return append(b, 0x00)
***REMOVED***

func readi32(src []byte) (int32, []byte, bool) ***REMOVED***
	if len(src) < 4 ***REMOVED***
		return 0, src, false
	***REMOVED***

	return (int32(src[0]) | int32(src[1])<<8 | int32(src[2])<<16 | int32(src[3])<<24), src[4:], true
***REMOVED***

func readi32unsafe(src []byte) int32 ***REMOVED***
	return (int32(src[0]) | int32(src[1])<<8 | int32(src[2])<<16 | int32(src[3])<<24)
***REMOVED***

func readi64(src []byte) (int64, []byte, bool) ***REMOVED***
	if len(src) < 8 ***REMOVED***
		return 0, src, false
	***REMOVED***
	i64 := (int64(src[0]) | int64(src[1])<<8 | int64(src[2])<<16 | int64(src[3])<<24 |
		int64(src[4])<<32 | int64(src[5])<<40 | int64(src[6])<<48 | int64(src[7])<<56)
	return i64, src[8:], true
***REMOVED***

func readcstring(src []byte) (string, []byte, bool) ***REMOVED***
	idx := bytes.IndexByte(src, 0x00)
	if idx < 0 ***REMOVED***
		return "", src, false
	***REMOVED***
	return string(src[:idx]), src[idx+1:], true
***REMOVED***
