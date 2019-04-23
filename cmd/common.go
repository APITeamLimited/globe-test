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

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/loadimpact/k6/lib/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	null "gopkg.in/guregu/null.v3"
)

// Use these when interacting with fs and writing to terminal, makes a command testable
var defaultFs = afero.NewOsFs()
var defaultWriter io.Writer = os.Stdout

// Panic if the given error is not nil.
func must(err error) ***REMOVED***
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

// Silently set an exit code.
type ExitCode struct ***REMOVED***
	error
	Code int
***REMOVED***

//TODO: refactor the CLI config so these functions aren't needed - they
// can mask errors by failing only at runtime, not at compile time
func getNullBool(flags *pflag.FlagSet, key string) null.Bool ***REMOVED***
	v, err := flags.GetBool(key)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return null.NewBool(v, flags.Changed(key))
***REMOVED***

func getNullInt64(flags *pflag.FlagSet, key string) null.Int ***REMOVED***
	v, err := flags.GetInt64(key)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return null.NewInt(v, flags.Changed(key))
***REMOVED***

func getNullDuration(flags *pflag.FlagSet, key string) types.NullDuration ***REMOVED***
	v, err := flags.GetDuration(key)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return types.NullDuration***REMOVED***Duration: types.Duration(v), Valid: flags.Changed(key)***REMOVED***
***REMOVED***

func getNullString(flags *pflag.FlagSet, key string) null.String ***REMOVED***
	v, err := flags.GetString(key)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return null.NewString(v, flags.Changed(key))
***REMOVED***

func exactArgsWithMsg(n int, msg string) cobra.PositionalArgs ***REMOVED***
	return func(cmd *cobra.Command, args []string) error ***REMOVED***
		if len(args) != n ***REMOVED***
			return fmt.Errorf("accepts %d arg(s), received %d: %s", n, len(args), msg)
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***
