// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonrw

import (
	"fmt"
)

type mode int

const (
	_ mode = iota
	mTopLevel
	mDocument
	mArray
	mValue
	mElement
	mCodeWithScope
	mSpacer
)

func (m mode) String() string ***REMOVED***
	var str string

	switch m ***REMOVED***
	case mTopLevel:
		str = "TopLevel"
	case mDocument:
		str = "DocumentMode"
	case mArray:
		str = "ArrayMode"
	case mValue:
		str = "ValueMode"
	case mElement:
		str = "ElementMode"
	case mCodeWithScope:
		str = "CodeWithScopeMode"
	case mSpacer:
		str = "CodeWithScopeSpacerFrame"
	default:
		str = "UnknownMode"
	***REMOVED***

	return str
***REMOVED***

func (m mode) TypeString() string ***REMOVED***
	var str string

	switch m ***REMOVED***
	case mTopLevel:
		str = "TopLevel"
	case mDocument:
		str = "Document"
	case mArray:
		str = "Array"
	case mValue:
		str = "Value"
	case mElement:
		str = "Element"
	case mCodeWithScope:
		str = "CodeWithScope"
	case mSpacer:
		str = "CodeWithScopeSpacer"
	default:
		str = "Unknown"
	***REMOVED***

	return str
***REMOVED***

// TransitionError is an error returned when an invalid progressing a
// ValueReader or ValueWriter state machine occurs.
// If read is false, the error is for writing
type TransitionError struct ***REMOVED***
	name        string
	parent      mode
	current     mode
	destination mode
	modes       []mode
	action      string
***REMOVED***

func (te TransitionError) Error() string ***REMOVED***
	errString := fmt.Sprintf("%s can only %s", te.name, te.action)
	if te.destination != mode(0) ***REMOVED***
		errString = fmt.Sprintf("%s a %s", errString, te.destination.TypeString())
	***REMOVED***
	errString = fmt.Sprintf("%s while positioned on a", errString)
	for ind, m := range te.modes ***REMOVED***
		if ind != 0 && len(te.modes) > 2 ***REMOVED***
			errString = fmt.Sprintf("%s,", errString)
		***REMOVED***
		if ind == len(te.modes)-1 && len(te.modes) > 1 ***REMOVED***
			errString = fmt.Sprintf("%s or", errString)
		***REMOVED***
		errString = fmt.Sprintf("%s %s", errString, m.TypeString())
	***REMOVED***
	errString = fmt.Sprintf("%s but is positioned on a %s", errString, te.current.TypeString())
	if te.parent != mode(0) ***REMOVED***
		errString = fmt.Sprintf("%s with parent %s", errString, te.parent.TypeString())
	***REMOVED***
	return errString
***REMOVED***
