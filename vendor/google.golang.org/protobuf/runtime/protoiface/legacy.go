// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protoiface

type MessageV1 interface ***REMOVED***
	Reset()
	String() string
	ProtoMessage()
***REMOVED***

type ExtensionRangeV1 struct ***REMOVED***
	Start, End int32 // both inclusive
***REMOVED***
