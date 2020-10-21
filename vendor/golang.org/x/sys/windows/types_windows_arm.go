// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package windows

type WSAData struct ***REMOVED***
	Version      uint16
	HighVersion  uint16
	Description  [WSADESCRIPTION_LEN + 1]byte
	SystemStatus [WSASYS_STATUS_LEN + 1]byte
	MaxSockets   uint16
	MaxUdpDg     uint16
	VendorInfo   *byte
***REMOVED***

type Servent struct ***REMOVED***
	Name    *byte
	Aliases **byte
	Port    uint16
	Proto   *byte
***REMOVED***
