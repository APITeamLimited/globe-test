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
	"time"

	"github.com/loadimpact/k6/stats"
	"gopkg.in/guregu/null.v3"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error ***REMOVED***
	var str string
	if err := json.Unmarshal(data, &str); err != nil ***REMOVED***
		return err
	***REMOVED***

	v, err := time.ParseDuration(str)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	*d = Duration(v)

	return nil
***REMOVED***

type TLSVersion struct ***REMOVED***
	Min int
	Max int
***REMOVED***

func (v *TLSVersion) UnmarshalJSON(data []byte) error ***REMOVED***
	// From https://golang.org/pkg/crypto/tls/#pkg-constants
	versionMap := map[string]int***REMOVED***
		"ssl3.0": tls.VersionSSL30,
		"tls1.0": tls.VersionTLS10,
		"tls1.1": tls.VersionTLS11,
		"tls1.2": tls.VersionTLS12,
	***REMOVED***

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

	var minVersion int
	var maxVersion int
	var ok bool
	if minVersion, ok = versionMap[fields.Min]; !ok ***REMOVED***
		return errors.New("Unknown TLS version : " + fields.Min)
	***REMOVED***

	if maxVersion, ok = versionMap[fields.Max]; !ok ***REMOVED***
		return errors.New("Unknown TLS version : " + fields.Max)
	***REMOVED***

	v.Min = minVersion
	v.Max = maxVersion

	return nil
***REMOVED***

type TLSCipherSuites struct ***REMOVED***
	Values []uint16
***REMOVED***

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

	s.Values = suiteIDs

	return nil
***REMOVED***

type Options struct ***REMOVED***
	Paused     null.Bool   `json:"paused"`
	VUs        null.Int    `json:"vus"`
	VUsMax     null.Int    `json:"vusMax"`
	Duration   null.String `json:"duration"`
	Iterations null.Int    `json:"iterations"`
	Stages     []Stage     `json:"stages"`

	Linger        null.Bool `json:"linger"`
	NoUsageReport null.Bool `json:"noUsageReport"`

	MaxRedirects          null.Int         `json:"maxRedirects"`
	InsecureSkipTLSVerify null.Bool        `json:"insecureSkipTLSVerify"`
	TLSCipherSuites       *TLSCipherSuites `json:"tlsCipherSuites"`
	TLSVersion            *TLSVersion      `json:"tlsVersion"`
	NoConnectionReuse     null.Bool        `json:"noConnectionReuse"`
	UserAgent             null.String      `json:"userAgent"`
	Throw                 null.Bool        `json:"throw"`

	Thresholds map[string]stats.Thresholds `json:"thresholds"`

	// These values are for third party collectors' benefit.
	External map[string]interface***REMOVED******REMOVED*** `json:"ext"`
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
	if opts.Linger.Valid ***REMOVED***
		o.Linger = opts.Linger
	***REMOVED***
	if opts.NoUsageReport.Valid ***REMOVED***
		o.NoUsageReport = opts.NoUsageReport
	***REMOVED***
	if opts.MaxRedirects.Valid ***REMOVED***
		o.MaxRedirects = opts.MaxRedirects
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
	if opts.External != nil ***REMOVED***
		o.External = opts.External
	***REMOVED***
	return o
***REMOVED***
