/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package lib

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"

	"github.com/loadimpact/k6/stats"
	"gopkg.in/guregu/null.v3"
)

type TLSVersion struct ***REMOVED***
	Min int
	Max int
***REMOVED***

func (v *TLSVersion) UnmarshalJSON(data []byte) error ***REMOVED***
	version := TLSVersion***REMOVED******REMOVED***

	// Version might be a string or an object with separate min & max fields
	var fields struct ***REMOVED***
		Min string `json:"min"`
		Max string `json:"max"`
	***REMOVED***
	if err := json.Unmarshal(data, &fields); err != nil ***REMOVED***
		switch err.(type) ***REMOVED***
		case *json.UnmarshalTypeError:
			// Check if it's a type error and the user has passed a string
			var version string
			if otherErr := json.Unmarshal(data, &version); otherErr != nil ***REMOVED***
				switch otherErr.(type) ***REMOVED***
				case *json.UnmarshalTypeError:
					return errors.New("Type error: the value of tlsVersion " +
						"should be an object with min/max fields or a string")
				***REMOVED***

				// Some other error occurred
				return otherErr
			***REMOVED***
			// It was a string, assign it to both min & max
			fields.Min = version
			fields.Max = version
		default:
			return err
		***REMOVED***
	***REMOVED***

	var ok bool
	if version.Min, ok = SupportedTLSVersions[fields.Min]; !ok ***REMOVED***
		return errors.New("Unknown TLS version : " + fields.Min)
	***REMOVED***

	if version.Max, ok = SupportedTLSVersions[fields.Max]; !ok ***REMOVED***
		return errors.New("Unknown TLS version : " + fields.Max)
	***REMOVED***

	*v = version

	return nil
***REMOVED***

type TLSCipherSuites []uint16

func (s *TLSCipherSuites) UnmarshalJSON(data []byte) error ***REMOVED***
	var suiteNames []string
	if err := json.Unmarshal(data, &suiteNames); err != nil ***REMOVED***
		return err
	***REMOVED***

	var suiteIDs []uint16
	for _, name := range suiteNames ***REMOVED***
		if suiteID, ok := SupportedTLSCipherSuites[name]; ok ***REMOVED***
			suiteIDs = append(suiteIDs, suiteID)
		***REMOVED*** else ***REMOVED***
			return errors.New("Unknown cipher suite: " + name)
		***REMOVED***
	***REMOVED***

	*s = suiteIDs

	return nil
***REMOVED***

type TLSAuthFields struct ***REMOVED***
	Cert    string   `json:"cert"`
	Key     string   `json:"key"`
	Domains []string `json:"domains"`
***REMOVED***

type TLSAuth struct ***REMOVED***
	TLSAuthFields
	certificate *tls.Certificate
***REMOVED***

func (c *TLSAuth) UnmarshalJSON(data []byte) error ***REMOVED***
	if err := json.Unmarshal(data, &c.TLSAuthFields); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := c.Certificate(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (c *TLSAuth) Certificate() (*tls.Certificate, error) ***REMOVED***
	if c.certificate == nil ***REMOVED***
		cert, err := tls.X509KeyPair([]byte(c.Cert), []byte(c.Key))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		c.certificate = &cert
	***REMOVED***
	return c.certificate, nil
***REMOVED***

type Options struct ***REMOVED***
	Paused     null.Bool    `json:"paused" envconfig:"paused"`
	VUs        null.Int     `json:"vus" envconfig:"vus"`
	VUsMax     null.Int     `json:"vusMax" envconfig:"vus_max"`
	Duration   NullDuration `json:"duration" envconfig:"duration"`
	Iterations null.Int     `json:"iterations" envconfig:"iterations"`
	Stages     []Stage      `json:"stages" envconfig:"stages"`

	MaxRedirects          null.Int         `json:"maxRedirects" envconfig:"max_redirects"`
	Batch                 null.Int         `json:"batch" envconfig:"batch"`
	InsecureSkipTLSVerify null.Bool        `json:"insecureSkipTLSVerify" envconfig:"insecure_skip_tls_verify"`
	TLSCipherSuites       *TLSCipherSuites `json:"tlsCipherSuites" envconfig:"tls_cipher_suites"`
	TLSVersion            *TLSVersion      `json:"tlsVersion" envconfig:"tls_version"`
	TLSAuth               []*TLSAuth       `json:"tlsAuth" envconfig:"tlsauth"`
	NoConnectionReuse     null.Bool        `json:"noConnectionReuse" envconfig:"no_connection_reuse"`
	UserAgent             null.String      `json:"userAgent" envconfig:"user_agent"`
	Throw                 null.Bool        `json:"throw" envconfig:"throw"`

	Thresholds   map[string]stats.Thresholds `json:"thresholds" envconfig:"thresholds"`
	BlacklistIPs []*net.IPNet                `json:"blacklistIPs" envconfig:"blacklist_ips"`

	// These values are for third party collectors' benefit.
	// Can't be set through env vars.
	External map[string]interface***REMOVED******REMOVED*** `json:"ext" ignored:"true"`
***REMOVED***

func (o Options) Apply(opts Options) Options ***REMOVED***
	if opts.Paused.Valid ***REMOVED***
		o.Paused = opts.Paused
	***REMOVED***
	if opts.VUs.Valid ***REMOVED***
		o.VUs = opts.VUs
	***REMOVED***
	if opts.VUsMax.Valid ***REMOVED***
		o.VUsMax = opts.VUsMax
	***REMOVED***
	if opts.Duration.Valid ***REMOVED***
		o.Duration = opts.Duration
	***REMOVED***
	if opts.Iterations.Valid ***REMOVED***
		o.Iterations = opts.Iterations
	***REMOVED***
	if opts.Stages != nil ***REMOVED***
		o.Stages = opts.Stages
	***REMOVED***
	if opts.MaxRedirects.Valid ***REMOVED***
		o.MaxRedirects = opts.MaxRedirects
	***REMOVED***
	if opts.Batch.Valid ***REMOVED***
		o.Batch = opts.Batch
	***REMOVED***
	if opts.InsecureSkipTLSVerify.Valid ***REMOVED***
		o.InsecureSkipTLSVerify = opts.InsecureSkipTLSVerify
	***REMOVED***
	if opts.TLSCipherSuites != nil ***REMOVED***
		o.TLSCipherSuites = opts.TLSCipherSuites
	***REMOVED***
	if opts.TLSVersion != nil ***REMOVED***
		o.TLSVersion = opts.TLSVersion
	***REMOVED***
	if opts.TLSAuth != nil ***REMOVED***
		o.TLSAuth = opts.TLSAuth
	***REMOVED***
	if opts.NoConnectionReuse.Valid ***REMOVED***
		o.NoConnectionReuse = opts.NoConnectionReuse
	***REMOVED***
	if opts.UserAgent.Valid ***REMOVED***
		o.UserAgent = opts.UserAgent
	***REMOVED***
	if opts.Throw.Valid ***REMOVED***
		o.Throw = opts.Throw
	***REMOVED***
	if opts.Thresholds != nil ***REMOVED***
		o.Thresholds = opts.Thresholds
	***REMOVED***
	if opts.BlacklistIPs != nil ***REMOVED***
		o.BlacklistIPs = opts.BlacklistIPs
	***REMOVED***
	if opts.External != nil ***REMOVED***
		o.External = opts.External
	***REMOVED***
	return o
***REMOVED***
