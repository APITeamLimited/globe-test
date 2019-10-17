/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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
	"fmt"
	"strings"

	"github.com/loadimpact/k6/js/compiler"
	null "gopkg.in/guregu/null.v3"
)

// RuntimeOptions are settings passed onto the goja JS runtime
type RuntimeOptions struct ***REMOVED***
	// Whether to pass the actual system environment variables to the JS runtime
	IncludeSystemEnvVars null.Bool `json:"includeSystemEnvVars" envconfig:"K6_INCLUDE_SYSTEM_ENV_VARS"`

	// JS compatibility mode: "extended" (Goja+Babel+core.js) or "base" (plain Goja)
	CompatibilityMode null.String `json:"compatibilityMode"`

	// Environment variables passed onto the runner
	Env map[string]string `json:"env" envconfig:"K6_ENV"`
***REMOVED***

// Apply overwrites the receiver RuntimeOptions' fields with any that are set
// on the argument struct and returns the receiver
func (o RuntimeOptions) Apply(opts RuntimeOptions) RuntimeOptions ***REMOVED***
	if opts.IncludeSystemEnvVars.Valid ***REMOVED***
		o.IncludeSystemEnvVars = opts.IncludeSystemEnvVars
	***REMOVED***
	if opts.CompatibilityMode.Valid ***REMOVED***
		o.CompatibilityMode = opts.CompatibilityMode
	***REMOVED***
	if opts.Env != nil ***REMOVED***
		o.Env = opts.Env
	***REMOVED***
	return o
***REMOVED***

// ValidateCompatibilityMode checks if the provided val is a valid compatibility mode
func ValidateCompatibilityMode(val string) (cm compiler.CompatibilityMode, err error) ***REMOVED***
	if val == "" ***REMOVED***
		return compiler.CompatibilityModeExtended, nil
	***REMOVED***
	if cm, err = compiler.CompatibilityModeString(val); err != nil ***REMOVED***
		var compatValues []string
		for _, v := range compiler.CompatibilityModeValues() ***REMOVED***
			compatValues = append(compatValues, v.String())
		***REMOVED***
		err = fmt.Errorf(`invalid compatibility mode "%s". Use: "%s"`,
			val, strings.Join(compatValues, `", "`))
	***REMOVED***
	return
***REMOVED***
