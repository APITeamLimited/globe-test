// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build go1.5,!go1.6

package codec

import "os"

func init() ***REMOVED***
	genCheckVendor = os.Getenv("GO15VENDOREXPERIMENT") == "1"
***REMOVED***