// +build !appengine

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

// Package credentials defines APIs for parsing SPIFFE ID.
//
// All APIs in this package are experimental.
package credentials

import (
	"crypto/tls"
	"crypto/x509"
	"net/url"

	"google.golang.org/grpc/grpclog"
)

var logger = grpclog.Component("credentials")

// SPIFFEIDFromState parses the SPIFFE ID from State. If the SPIFFE ID format
// is invalid, return nil with warning.
func SPIFFEIDFromState(state tls.ConnectionState) *url.URL ***REMOVED***
	if len(state.PeerCertificates) == 0 || len(state.PeerCertificates[0].URIs) == 0 ***REMOVED***
		return nil
	***REMOVED***
	return SPIFFEIDFromCert(state.PeerCertificates[0])
***REMOVED***

// SPIFFEIDFromCert parses the SPIFFE ID from x509.Certificate. If the SPIFFE
// ID format is invalid, return nil with warning.
func SPIFFEIDFromCert(cert *x509.Certificate) *url.URL ***REMOVED***
	if cert == nil || cert.URIs == nil ***REMOVED***
		return nil
	***REMOVED***
	var spiffeID *url.URL
	for _, uri := range cert.URIs ***REMOVED***
		if uri == nil || uri.Scheme != "spiffe" || uri.Opaque != "" || (uri.User != nil && uri.User.Username() != "") ***REMOVED***
			continue
		***REMOVED***
		// From this point, we assume the uri is intended for a SPIFFE ID.
		if len(uri.String()) > 2048 ***REMOVED***
			logger.Warning("invalid SPIFFE ID: total ID length larger than 2048 bytes")
			return nil
		***REMOVED***
		if len(uri.Host) == 0 || len(uri.Path) == 0 ***REMOVED***
			logger.Warning("invalid SPIFFE ID: domain or workload ID is empty")
			return nil
		***REMOVED***
		if len(uri.Host) > 255 ***REMOVED***
			logger.Warning("invalid SPIFFE ID: domain length larger than 255 characters")
			return nil
		***REMOVED***
		// A valid SPIFFE certificate can only have exactly one URI SAN field.
		if len(cert.URIs) > 1 ***REMOVED***
			logger.Warning("invalid SPIFFE ID: multiple URI SANs")
			return nil
		***REMOVED***
		spiffeID = uri
	***REMOVED***
	return spiffeID
***REMOVED***
