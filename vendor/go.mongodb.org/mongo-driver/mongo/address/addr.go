// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package address provides structured representations of network addresses.
package address // import "go.mongodb.org/mongo-driver/mongo/address"

import (
	"net"
	"strings"
)

const defaultPort = "27017"

// Address is a network address. It can either be an IP address or a DNS name.
type Address string

// Network is the network protocol for this address. In most cases this will be
// "tcp" or "unix".
func (a Address) Network() string ***REMOVED***
	if strings.HasSuffix(string(a), "sock") ***REMOVED***
		return "unix"
	***REMOVED***
	return "tcp"
***REMOVED***

// String is the canonical version of this address, e.g. localhost:27017,
// 1.2.3.4:27017, example.com:27017.
func (a Address) String() string ***REMOVED***
	// TODO: unicode case folding?
	s := strings.ToLower(string(a))
	if len(s) == 0 ***REMOVED***
		return ""
	***REMOVED***
	if a.Network() != "unix" ***REMOVED***
		_, _, err := net.SplitHostPort(s)
		if err != nil && strings.Contains(err.Error(), "missing port in address") ***REMOVED***
			s += ":" + defaultPort
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***

// Canonicalize creates a canonicalized address.
func (a Address) Canonicalize() Address ***REMOVED***
	return Address(a.String())
***REMOVED***
