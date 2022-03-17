/*
 *
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package insecure provides an implementation of the
// credentials.TransportCredentials interface which disables transport security.
package insecure

import (
	"context"
	"net"

	"google.golang.org/grpc/credentials"
)

// NewCredentials returns a credentials which disables transport security.
//
// Note that using this credentials with per-RPC credentials which require
// transport security is incompatible and will cause grpc.Dial() to fail.
func NewCredentials() credentials.TransportCredentials ***REMOVED***
	return insecureTC***REMOVED******REMOVED***
***REMOVED***

// insecureTC implements the insecure transport credentials. The handshake
// methods simply return the passed in net.Conn and set the security level to
// NoSecurity.
type insecureTC struct***REMOVED******REMOVED***

func (insecureTC) ClientHandshake(ctx context.Context, _ string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) ***REMOVED***
	return conn, info***REMOVED***credentials.CommonAuthInfo***REMOVED***SecurityLevel: credentials.NoSecurity***REMOVED******REMOVED***, nil
***REMOVED***

func (insecureTC) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) ***REMOVED***
	return conn, info***REMOVED***credentials.CommonAuthInfo***REMOVED***SecurityLevel: credentials.NoSecurity***REMOVED******REMOVED***, nil
***REMOVED***

func (insecureTC) Info() credentials.ProtocolInfo ***REMOVED***
	return credentials.ProtocolInfo***REMOVED***SecurityProtocol: "insecure"***REMOVED***
***REMOVED***

func (insecureTC) Clone() credentials.TransportCredentials ***REMOVED***
	return insecureTC***REMOVED******REMOVED***
***REMOVED***

func (insecureTC) OverrideServerName(string) error ***REMOVED***
	return nil
***REMOVED***

// info contains the auth information for an insecure connection.
// It implements the AuthInfo interface.
type info struct ***REMOVED***
	credentials.CommonAuthInfo
***REMOVED***

// AuthType returns the type of info as a string.
func (info) AuthType() string ***REMOVED***
	return "insecure"
***REMOVED***

// insecureBundle implements an insecure bundle.
// An insecure bundle provides a thin wrapper around insecureTC to support
// the credentials.Bundle interface.
type insecureBundle struct***REMOVED******REMOVED***

// NewBundle returns a bundle with disabled transport security and no per rpc credential.
func NewBundle() credentials.Bundle ***REMOVED***
	return insecureBundle***REMOVED******REMOVED***
***REMOVED***

// NewWithMode returns a new insecure Bundle. The mode is ignored.
func (insecureBundle) NewWithMode(string) (credentials.Bundle, error) ***REMOVED***
	return insecureBundle***REMOVED******REMOVED***, nil
***REMOVED***

// PerRPCCredentials returns an nil implementation as insecure
// bundle does not support a per rpc credential.
func (insecureBundle) PerRPCCredentials() credentials.PerRPCCredentials ***REMOVED***
	return nil
***REMOVED***

// TransportCredentials returns the underlying insecure transport credential.
func (insecureBundle) TransportCredentials() credentials.TransportCredentials ***REMOVED***
	return NewCredentials()
***REMOVED***
