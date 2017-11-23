// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.8

package websocket

import "crypto/tls"

func cloneTLSConfig(cfg *tls.Config) *tls.Config ***REMOVED***
	if cfg == nil ***REMOVED***
		return &tls.Config***REMOVED******REMOVED***
	***REMOVED***
	return cfg.Clone()
***REMOVED***
