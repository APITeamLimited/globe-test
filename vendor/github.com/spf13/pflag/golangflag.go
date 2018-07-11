// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	goflag "flag"
	"reflect"
	"strings"
)

// flagValueWrapper implements pflag.Value around a flag.Value.  The main
// difference here is the addition of the Type method that returns a string
// name of the type.  As this is generally unknown, we approximate that with
// reflection.
type flagValueWrapper struct ***REMOVED***
	inner    goflag.Value
	flagType string
***REMOVED***

// We are just copying the boolFlag interface out of goflag as that is what
// they use to decide if a flag should get "true" when no arg is given.
type goBoolFlag interface ***REMOVED***
	goflag.Value
	IsBoolFlag() bool
***REMOVED***

func wrapFlagValue(v goflag.Value) Value ***REMOVED***
	// If the flag.Value happens to also be a pflag.Value, just use it directly.
	if pv, ok := v.(Value); ok ***REMOVED***
		return pv
	***REMOVED***

	pv := &flagValueWrapper***REMOVED***
		inner: v,
	***REMOVED***

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Interface || t.Kind() == reflect.Ptr ***REMOVED***
		t = t.Elem()
	***REMOVED***

	pv.flagType = strings.TrimSuffix(t.Name(), "Value")
	return pv
***REMOVED***

func (v *flagValueWrapper) String() string ***REMOVED***
	return v.inner.String()
***REMOVED***

func (v *flagValueWrapper) Set(s string) error ***REMOVED***
	return v.inner.Set(s)
***REMOVED***

func (v *flagValueWrapper) Type() string ***REMOVED***
	return v.flagType
***REMOVED***

// PFlagFromGoFlag will return a *pflag.Flag given a *flag.Flag
// If the *flag.Flag.Name was a single character (ex: `v`) it will be accessiblei
// with both `-v` and `--v` in flags. If the golang flag was more than a single
// character (ex: `verbose`) it will only be accessible via `--verbose`
func PFlagFromGoFlag(goflag *goflag.Flag) *Flag ***REMOVED***
	// Remember the default value as a string; it won't change.
	flag := &Flag***REMOVED***
		Name:  goflag.Name,
		Usage: goflag.Usage,
		Value: wrapFlagValue(goflag.Value),
		// Looks like golang flags don't set DefValue correctly  :-(
		//DefValue: goflag.DefValue,
		DefValue: goflag.Value.String(),
	***REMOVED***
	// Ex: if the golang flag was -v, allow both -v and --v to work
	if len(flag.Name) == 1 ***REMOVED***
		flag.Shorthand = flag.Name
	***REMOVED***
	if fv, ok := goflag.Value.(goBoolFlag); ok && fv.IsBoolFlag() ***REMOVED***
		flag.NoOptDefVal = "true"
	***REMOVED***
	return flag
***REMOVED***

// AddGoFlag will add the given *flag.Flag to the pflag.FlagSet
func (f *FlagSet) AddGoFlag(goflag *goflag.Flag) ***REMOVED***
	if f.Lookup(goflag.Name) != nil ***REMOVED***
		return
	***REMOVED***
	newflag := PFlagFromGoFlag(goflag)
	f.AddFlag(newflag)
***REMOVED***

// AddGoFlagSet will add the given *flag.FlagSet to the pflag.FlagSet
func (f *FlagSet) AddGoFlagSet(newSet *goflag.FlagSet) ***REMOVED***
	if newSet == nil ***REMOVED***
		return
	***REMOVED***
	newSet.VisitAll(func(goflag *goflag.Flag) ***REMOVED***
		f.AddGoFlag(goflag)
	***REMOVED***)
	if f.addedGoFlagSets == nil ***REMOVED***
		f.addedGoFlagSets = make([]*goflag.FlagSet, 0)
	***REMOVED***
	f.addedGoFlagSets = append(f.addedGoFlagSets, newSet)
***REMOVED***
