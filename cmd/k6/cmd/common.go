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
	"bytes"
	"io"
	"sync"

	"github.com/loadimpact/k6/lib"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	null "gopkg.in/guregu/null.v3"
)

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

// A writer that syncs writes with a mutex and, if the output is a TTY, clears before newlines.
type consoleWriter struct ***REMOVED***
	Writer io.Writer
	IsTTY  bool
	Mutex  *sync.Mutex
***REMOVED***

func (w consoleWriter) Write(p []byte) (n int, err error) ***REMOVED***
	if w.IsTTY ***REMOVED***
		p = bytes.Replace(p, []byte***REMOVED***'\n'***REMOVED***, []byte***REMOVED***'\x1b', '[', '0', 'K', '\n'***REMOVED***, -1)
	***REMOVED***
	w.Mutex.Lock()
	n, err = w.Writer.Write(p)
	w.Mutex.Unlock()
	return
***REMOVED***

func registerOptions(flags *pflag.FlagSet) ***REMOVED***
	flags.Int64P("vus", "u", 1, "number of virtual users")
	flags.Int64P("max", "m", 0, "max available virtual users")
	flags.DurationP("duration", "d", 0, "test duration limit")
	flags.Int64P("iterations", "i", 0, "script iteration limit")
	flags.StringSliceP("stage", "s", nil, "add a `stage`, as `[duration]:[target]`")
	flags.BoolP("paused", "p", false, "start the test in a paused state")
	flags.Int64("max-redirects", 10, "follow at most n redirects")
	flags.String("user-agent", "", "user agent for http requests")
	flags.Bool("insecure-skip-tls-verify", false, "skip verification of TLS certificates")
	flags.Bool("no-connection-reuse", false, "don't reuse connections between iterations")
	flags.BoolP("throw", "w", false, "throw warnings (like failed http requests) as errors")
***REMOVED***

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

func getNullDuration(flags *pflag.FlagSet, key string) lib.NullDuration ***REMOVED***
	v, err := flags.GetDuration(key)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return lib.NullDuration***REMOVED***Duration: lib.Duration(v), Valid: flags.Changed(key)***REMOVED***
***REMOVED***

func getNullString(flags *pflag.FlagSet, key string) null.String ***REMOVED***
	v, err := flags.GetString(key)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return null.NewString(v, flags.Changed(key))
***REMOVED***

func getOptions(flags *pflag.FlagSet) (lib.Options, error) ***REMOVED***
	opts := lib.Options***REMOVED***
		VUs:                   getNullInt64(flags, "vus"),
		VUsMax:                getNullInt64(flags, "max"),
		Duration:              getNullDuration(flags, "duration"),
		Iterations:            getNullInt64(flags, "iterations"),
		Paused:                getNullBool(flags, "paused"),
		MaxRedirects:          getNullInt64(flags, "max-redirects"),
		UserAgent:             getNullString(flags, "user-agent"),
		InsecureSkipTLSVerify: getNullBool(flags, "insecure-skip-tls-verify"),
		NoConnectionReuse:     getNullBool(flags, "no-connection-reuse"),
		Throw:                 getNullBool(flags, "throw"),
	***REMOVED***

	stageStrings, err := flags.GetStringSlice("stage")
	if err != nil ***REMOVED***
		return opts, err
	***REMOVED***
	if len(stageStrings) > 0 ***REMOVED***
		opts.Stages = make([]lib.Stage, len(stageStrings))
		for i, s := range stageStrings ***REMOVED***
			var stage lib.Stage
			if err := stage.UnmarshalText([]byte(s)); err != nil ***REMOVED***
				return opts, errors.Wrapf(err, "stage %d", i)
			***REMOVED***
			opts.Stages[i] = stage
		***REMOVED***
	***REMOVED***
	return opts, nil
***REMOVED***
