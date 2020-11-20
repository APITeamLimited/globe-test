// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package windows

type WSAData struct ***REMOVED***
	Version      uint16
	HighVersion  uint16
	MaxSockets   uint16
	MaxUdpDg     uint16
	VendorInfo   *byte
	Description  [WSADESCRIPTION_LEN + 1]byte
	SystemStatus [WSASYS_STATUS_LEN + 1]byte
***REMOVED***

type Servent struct ***REMOVED***
	Name    *byte
	Aliases **byte
	Proto   *byte
	Port    uint16
***REMOVED***

type JOBOBJECT_BASIC_LIMIT_INFORMATION struct ***REMOVED***
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
***REMOVED***
