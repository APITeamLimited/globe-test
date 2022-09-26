package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/k6-worker/lib"
)

func TestTLS13Support(t *testing.T) ***REMOVED***
	tb, state, _, rt, _ := newRuntime(t)

	tb.Mux.HandleFunc("/tls-version", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) ***REMOVED***
		ver := req.TLS.Version
		fmt.Fprint(resp, lib.SupportedTLSVersionsToString[lib.TLSVersion(ver)])
	***REMOVED***))

	// We don't expect any failed requests
	state.Options.Throw = null.BoolFrom(true)
	state.Options.Apply(lib.Options***REMOVED***TLSVersion: &lib.TLSVersions***REMOVED***Max: tls.VersionTLS13***REMOVED******REMOVED***)

	_, err := rt.RunString(tb.Replacer.Replace(`
		var resp = http.get("HTTPSBIN_URL/tls-version");
		if (resp.body != "tls1.3") ***REMOVED***
			throw new Error("unexpected tls version: " + resp.body);
		***REMOVED***
	`))
	assert.NoError(t, err)
***REMOVED***
