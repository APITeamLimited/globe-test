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
	"gopkg.in/guregu/null.v3"
)

type Options struct ***REMOVED***
	Paused     null.Bool   `json:"paused"`
	VUs        null.Int    `json:"vus"`
	VUsMax     null.Int    `json:"vusMax"`
	Duration   null.String `json:"duration"`
	Iterations null.Int    `json:"iterations"`
	Stages     []Stage     `json:"stage"`

	Linger        null.Bool `json:"linger"`
	NoUsageReport null.Bool `json:"noUsageReport"`

	MaxRedirects          null.Int  `json:"maxRedirects"`
	InsecureSkipTLSVerify null.Bool `json:"insecureSkipTLSVerify"`

	Thresholds map[string]Thresholds `json:"thresholds"`
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
	if opts.Thresholds != nil ***REMOVED***
		o.Thresholds = opts.Thresholds
	***REMOVED***
	return o
***REMOVED***

func (o Options) SetAllValid(valid bool) Options ***REMOVED***
	o.Paused.Valid = valid
	o.VUs.Valid = valid
	o.VUsMax.Valid = valid
	o.Duration.Valid = valid
	o.Linger.Valid = valid
	o.NoUsageReport.Valid = valid
	return o
***REMOVED***
