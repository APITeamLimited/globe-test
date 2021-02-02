// +build go1.12

/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
)

func TestTLS13Support(t *testing.T) ***REMOVED***
	tb, state, _, rt, _ := newRuntime(t)
	defer tb.Cleanup()

	tb.Mux.HandleFunc("/tls-version", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) ***REMOVED***
		ver := req.TLS.Version
		fmt.Fprint(resp, lib.SupportedTLSVersionsToString[lib.TLSVersion(ver)])
	***REMOVED***))

	// We don't expect any failed requests
	state.Options.Throw = null.BoolFrom(true)
	state.Options.Apply(lib.Options***REMOVED***TLSVersion: &lib.TLSVersions***REMOVED***Max: lib.TLSVersion13***REMOVED******REMOVED***)

	_, err := rt.RunString(tb.Replacer.Replace(`
		var resp = http.get("HTTPSBIN_URL/tls-version");
		if (resp.body != "tls1.3") ***REMOVED***
			throw new Error("unexpected tls version: " + resp.body);
		***REMOVED***
	`))
	assert.NoError(t, err)
***REMOVED***
