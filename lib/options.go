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
	"encoding/json"
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

type Options struct ***REMOVED***
	Paused     null.Bool   `json:"paused"`
	VUs        null.Int    `json:"vus"`
	VUsMax     null.Int    `json:"vusMax"`
	Duration   null.String `json:"duration"`
	Iterations null.Int    `json:"iterations"`
	Stages     []Stage     `json:"stages"`

	Linger        null.Bool `json:"linger"`
	NoUsageReport null.Bool `json:"noUsageReport"`

	MaxRedirects          null.Int    `json:"maxRedirects"`
	InsecureSkipTLSVerify null.Bool   `json:"insecureSkipTLSVerify"`
	NoConnectionReuse     null.Bool   `json:"noConnectionReuse"`
	UserAgent             null.String `json:"userAgent"`
	Throw                 null.Bool   `json:"throw"`

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
	if opts.NoConnectionReuse.Valid ***REMOVED***
		o.NoConnectionReuse = opts.NoConnectionReuse
	***REMOVED***
	if opts.Thresholds != nil ***REMOVED***
		o.Thresholds = opts.Thresholds
	***REMOVED***
	if opts.External != nil ***REMOVED***
		o.External = opts.External
	***REMOVED***
	return o
***REMOVED***
