// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"fmt"
)

// ServerAPIOptions represents options used to configure the API version sent to the server
// when running commands.
//
// Sending a specified server API version causes the server to behave in a manner compatible with that
// API version. It also causes the driver to behave in a manner compatible with the driver’s behavior as
// of the release when the driver first started to support the specified server API version.
//
// The user must specify a ServerAPIVersion if including ServerAPIOptions in their client. That version
// must also be currently supported by the driver. This version of the driver supports API version "1".
type ServerAPIOptions struct ***REMOVED***
	ServerAPIVersion  ServerAPIVersion
	Strict            *bool
	DeprecationErrors *bool
***REMOVED***

// ServerAPI creates a new ServerAPIOptions configured with the provided serverAPIversion.
func ServerAPI(serverAPIVersion ServerAPIVersion) *ServerAPIOptions ***REMOVED***
	return &ServerAPIOptions***REMOVED***ServerAPIVersion: serverAPIVersion***REMOVED***
***REMOVED***

// SetStrict specifies whether the server should return errors for features that are not part of the API version.
func (s *ServerAPIOptions) SetStrict(strict bool) *ServerAPIOptions ***REMOVED***
	s.Strict = &strict
	return s
***REMOVED***

// SetDeprecationErrors specifies whether the server should return errors for deprecated features.
func (s *ServerAPIOptions) SetDeprecationErrors(deprecationErrors bool) *ServerAPIOptions ***REMOVED***
	s.DeprecationErrors = &deprecationErrors
	return s
***REMOVED***

// ServerAPIVersion represents an API version that can be used in ServerAPIOptions.
type ServerAPIVersion string

const (
	// ServerAPIVersion1 is the first API version.
	ServerAPIVersion1 ServerAPIVersion = "1"
)

// Validate determines if the provided ServerAPIVersion is currently supported by the driver.
func (sav ServerAPIVersion) Validate() error ***REMOVED***
	switch sav ***REMOVED***
	case ServerAPIVersion1:
		return nil
	***REMOVED***
	return fmt.Errorf("api version %q not supported; this driver version only supports API version \"1\"", sav)
***REMOVED***
