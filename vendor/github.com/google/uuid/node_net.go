// Copyright 2017 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !js

package uuid

import "net"

var interfaces []net.Interface // cached list of interfaces

// getHardwareInterface returns the name and hardware address of interface name.
// If name is "" then the name and hardware address of one of the system's
// interfaces is returned.  If no interfaces are found (name does not exist or
// there are no interfaces) then "", nil is returned.
//
// Only addresses of at least 6 bytes are returned.
func getHardwareInterface(name string) (string, []byte) ***REMOVED***
	if interfaces == nil ***REMOVED***
		var err error
		interfaces, err = net.Interfaces()
		if err != nil ***REMOVED***
			return "", nil
		***REMOVED***
	***REMOVED***
	for _, ifs := range interfaces ***REMOVED***
		if len(ifs.HardwareAddr) >= 6 && (name == "" || name == ifs.Name) ***REMOVED***
			return ifs.Name, ifs.HardwareAddr
		***REMOVED***
	***REMOVED***
	return "", nil
***REMOVED***
