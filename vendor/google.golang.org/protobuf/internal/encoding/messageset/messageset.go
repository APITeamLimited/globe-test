// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package messageset encodes and decodes the obsolete MessageSet wire format.
package messageset

import (
	"math"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/errors"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// The MessageSet wire format is equivalent to a message defined as follows,
// where each Item defines an extension field with a field number of 'type_id'
// and content of 'message'. MessageSet extensions must be non-repeated message
// fields.
//
//	message MessageSet ***REMOVED***
//		repeated group Item = 1 ***REMOVED***
//			required int32 type_id = 2;
//			required string message = 3;
//		***REMOVED***
//	***REMOVED***
const (
	FieldItem    = protowire.Number(1)
	FieldTypeID  = protowire.Number(2)
	FieldMessage = protowire.Number(3)
)

// ExtensionName is the field name for extensions of MessageSet.
//
// A valid MessageSet extension must be of the form:
//	message MyMessage ***REMOVED***
//		extend proto2.bridge.MessageSet ***REMOVED***
//			optional MyMessage message_set_extension = 1234;
//		***REMOVED***
//		...
//	***REMOVED***
const ExtensionName = "message_set_extension"

// IsMessageSet returns whether the message uses the MessageSet wire format.
func IsMessageSet(md pref.MessageDescriptor) bool ***REMOVED***
	xmd, ok := md.(interface***REMOVED*** IsMessageSet() bool ***REMOVED***)
	return ok && xmd.IsMessageSet()
***REMOVED***

// IsMessageSetExtension reports this field properly extends a MessageSet.
func IsMessageSetExtension(fd pref.FieldDescriptor) bool ***REMOVED***
	switch ***REMOVED***
	case fd.Name() != ExtensionName:
		return false
	case !IsMessageSet(fd.ContainingMessage()):
		return false
	case fd.FullName().Parent() != fd.Message().FullName():
		return false
	***REMOVED***
	return true
***REMOVED***

// SizeField returns the size of a MessageSet item field containing an extension
// with the given field number, not counting the contents of the message subfield.
func SizeField(num protowire.Number) int ***REMOVED***
	return 2*protowire.SizeTag(FieldItem) + protowire.SizeTag(FieldTypeID) + protowire.SizeVarint(uint64(num))
***REMOVED***

// Unmarshal parses a MessageSet.
//
// It calls fn with the type ID and value of each item in the MessageSet.
// Unknown fields are discarded.
//
// If wantLen is true, the item values include the varint length prefix.
// This is ugly, but simplifies the fast-path decoder in internal/impl.
func Unmarshal(b []byte, wantLen bool, fn func(typeID protowire.Number, value []byte) error) error ***REMOVED***
	for len(b) > 0 ***REMOVED***
		num, wtyp, n := protowire.ConsumeTag(b)
		if n < 0 ***REMOVED***
			return protowire.ParseError(n)
		***REMOVED***
		b = b[n:]
		if num != FieldItem || wtyp != protowire.StartGroupType ***REMOVED***
			n := protowire.ConsumeFieldValue(num, wtyp, b)
			if n < 0 ***REMOVED***
				return protowire.ParseError(n)
			***REMOVED***
			b = b[n:]
			continue
		***REMOVED***
		typeID, value, n, err := ConsumeFieldValue(b, wantLen)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		b = b[n:]
		if typeID == 0 ***REMOVED***
			continue
		***REMOVED***
		if err := fn(typeID, value); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// ConsumeFieldValue parses b as a MessageSet item field value until and including
// the trailing end group marker. It assumes the start group tag has already been parsed.
// It returns the contents of the type_id and message subfields and the total
// item length.
//
// If wantLen is true, the returned message value includes the length prefix.
func ConsumeFieldValue(b []byte, wantLen bool) (typeid protowire.Number, message []byte, n int, err error) ***REMOVED***
	ilen := len(b)
	for ***REMOVED***
		num, wtyp, n := protowire.ConsumeTag(b)
		if n < 0 ***REMOVED***
			return 0, nil, 0, protowire.ParseError(n)
		***REMOVED***
		b = b[n:]
		switch ***REMOVED***
		case num == FieldItem && wtyp == protowire.EndGroupType:
			if wantLen && len(message) == 0 ***REMOVED***
				// The message field was missing, which should never happen.
				// Be prepared for this case anyway.
				message = protowire.AppendVarint(message, 0)
			***REMOVED***
			return typeid, message, ilen - len(b), nil
		case num == FieldTypeID && wtyp == protowire.VarintType:
			v, n := protowire.ConsumeVarint(b)
			if n < 0 ***REMOVED***
				return 0, nil, 0, protowire.ParseError(n)
			***REMOVED***
			b = b[n:]
			if v < 1 || v > math.MaxInt32 ***REMOVED***
				return 0, nil, 0, errors.New("invalid type_id in message set")
			***REMOVED***
			typeid = protowire.Number(v)
		case num == FieldMessage && wtyp == protowire.BytesType:
			m, n := protowire.ConsumeBytes(b)
			if n < 0 ***REMOVED***
				return 0, nil, 0, protowire.ParseError(n)
			***REMOVED***
			if message == nil ***REMOVED***
				if wantLen ***REMOVED***
					message = b[:n:n]
				***REMOVED*** else ***REMOVED***
					message = m[:len(m):len(m)]
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// This case should never happen in practice, but handle it for
				// correctness: The MessageSet item contains multiple message
				// fields, which need to be merged.
				//
				// In the case where we're returning the length, this becomes
				// quite inefficient since we need to strip the length off
				// the existing data and reconstruct it with the combined length.
				if wantLen ***REMOVED***
					_, nn := protowire.ConsumeVarint(message)
					m0 := message[nn:]
					message = nil
					message = protowire.AppendVarint(message, uint64(len(m0)+len(m)))
					message = append(message, m0...)
					message = append(message, m...)
				***REMOVED*** else ***REMOVED***
					message = append(message, m...)
				***REMOVED***
			***REMOVED***
			b = b[n:]
		default:
			// We have no place to put it, so we just ignore unknown fields.
			n := protowire.ConsumeFieldValue(num, wtyp, b)
			if n < 0 ***REMOVED***
				return 0, nil, 0, protowire.ParseError(n)
			***REMOVED***
			b = b[n:]
		***REMOVED***
	***REMOVED***
***REMOVED***

// AppendFieldStart appends the start of a MessageSet item field containing
// an extension with the given number. The caller must add the message
// subfield (including the tag).
func AppendFieldStart(b []byte, num protowire.Number) []byte ***REMOVED***
	b = protowire.AppendTag(b, FieldItem, protowire.StartGroupType)
	b = protowire.AppendTag(b, FieldTypeID, protowire.VarintType)
	b = protowire.AppendVarint(b, uint64(num))
	return b
***REMOVED***

// AppendFieldEnd appends the trailing end group marker for a MessageSet item field.
func AppendFieldEnd(b []byte) []byte ***REMOVED***
	return protowire.AppendTag(b, FieldItem, protowire.EndGroupType)
***REMOVED***

// SizeUnknown returns the size of an unknown fields section in MessageSet format.
//
// See AppendUnknown.
func SizeUnknown(unknown []byte) (size int) ***REMOVED***
	for len(unknown) > 0 ***REMOVED***
		num, typ, n := protowire.ConsumeTag(unknown)
		if n < 0 || typ != protowire.BytesType ***REMOVED***
			return 0
		***REMOVED***
		unknown = unknown[n:]
		_, n = protowire.ConsumeBytes(unknown)
		if n < 0 ***REMOVED***
			return 0
		***REMOVED***
		unknown = unknown[n:]
		size += SizeField(num) + protowire.SizeTag(FieldMessage) + n
	***REMOVED***
	return size
***REMOVED***

// AppendUnknown appends unknown fields to b in MessageSet format.
//
// For historic reasons, unresolved items in a MessageSet are stored in a
// message's unknown fields section in non-MessageSet format. That is, an
// unknown item with typeID T and value V appears in the unknown fields as
// a field with number T and value V.
//
// This function converts the unknown fields back into MessageSet form.
func AppendUnknown(b, unknown []byte) ([]byte, error) ***REMOVED***
	for len(unknown) > 0 ***REMOVED***
		num, typ, n := protowire.ConsumeTag(unknown)
		if n < 0 || typ != protowire.BytesType ***REMOVED***
			return nil, errors.New("invalid data in message set unknown fields")
		***REMOVED***
		unknown = unknown[n:]
		_, n = protowire.ConsumeBytes(unknown)
		if n < 0 ***REMOVED***
			return nil, errors.New("invalid data in message set unknown fields")
		***REMOVED***
		b = AppendFieldStart(b, num)
		b = protowire.AppendTag(b, FieldMessage, protowire.BytesType)
		b = append(b, unknown[:n]...)
		b = AppendFieldEnd(b)
		unknown = unknown[n:]
	***REMOVED***
	return b, nil
***REMOVED***
