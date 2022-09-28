// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "go.mongodb.org/mongo-driver/mongo/readpref"

// RunCmdOptions represents options that can be used to configure a RunCommand operation.
type RunCmdOptions struct ***REMOVED***
	// The read preference to use for the operation. The default value is nil, which means that the primary read
	// preference will be used.
	ReadPreference *readpref.ReadPref
***REMOVED***

// RunCmd creates a new RunCmdOptions instance.
func RunCmd() *RunCmdOptions ***REMOVED***
	return &RunCmdOptions***REMOVED******REMOVED***
***REMOVED***

// SetReadPreference sets value for the ReadPreference field.
func (rc *RunCmdOptions) SetReadPreference(rp *readpref.ReadPref) *RunCmdOptions ***REMOVED***
	rc.ReadPreference = rp
	return rc
***REMOVED***

// MergeRunCmdOptions combines the given RunCmdOptions instances into one *RunCmdOptions in a last-one-wins fashion.
func MergeRunCmdOptions(opts ...*RunCmdOptions) *RunCmdOptions ***REMOVED***
	rc := RunCmd()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.ReadPreference != nil ***REMOVED***
			rc.ReadPreference = opt.ReadPreference
		***REMOVED***
	***REMOVED***

	return rc
***REMOVED***
