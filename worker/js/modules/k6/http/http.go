package http

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker/netext"
	"github.com/APITeamLimited/globe-test/worker/libWorker/netext/httpext"
	"github.com/dop251/goja"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create HTTP module instances for each VU.
//
// TODO: add sync.Once for all of the deprecation warnings we might want to do
// for the old k6/http APIs here, so they are shown only once in a test run.
type RootModule struct***REMOVED******REMOVED***

// ModuleInstance represents an instance of the HTTP module for every VU.
type ModuleInstance struct ***REMOVED***
	vu            modules.VU
	rootModule    *RootModule
	defaultClient *Client
	exports       *goja.Object
***REMOVED***

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &ModuleInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new HTTP RootModule.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance returns an HTTP module instance for each VU.
func (r *RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	rt := vu.Runtime()
	mi := &ModuleInstance***REMOVED***
		vu:         vu,
		rootModule: r,
		exports:    rt.NewObject(),
	***REMOVED***
	mi.defineConstants()

	mi.defaultClient = &Client***REMOVED***
		// TODO: configure this from libWorker.Options and get rid of some of the
		// things in the VU State struct that should be here. See
		// https://github.com/grafana/k6/issues/2293
		moduleInstance:   mi,
		responseCallback: defaultExpectedStatuses.match,
	***REMOVED***

	mustExport := func(name string, value interface***REMOVED******REMOVED***) ***REMOVED***
		if err := mi.exports.Set(name, value); err != nil ***REMOVED***
			common.Throw(rt, err)
		***REMOVED***
	***REMOVED***

	mustExport("url", mi.URL)
	mustExport("CookieJar", mi.newCookieJar)
	mustExport("cookieJar", mi.getVUCookieJar)
	mustExport("file", mi.file) // TODO: deprecate or refactor?

	// TODO: refactor so the Client actually has better APIs and these are
	// wrappers (facades) that convert the old k6 idiosyncratic APIs to the new
	// proper Client ones that accept Request objects and don't suck
	mustExport("get", func(url goja.Value, args ...goja.Value) (*Response, error) ***REMOVED***
		// http.get(url, params) doesn't have a body argument, so we add undefined
		// as the third argument to http.request(method, url, body, params)
		args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
		return mi.defaultClient.Request(http.MethodGet, url, args...)
	***REMOVED***)
	mustExport("head", func(url goja.Value, args ...goja.Value) (*Response, error) ***REMOVED***
		// http.head(url, params) doesn't have a body argument, so we add undefined
		// as the third argument to http.request(method, url, body, params)
		args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
		return mi.defaultClient.Request(http.MethodHead, url, args...)
	***REMOVED***)
	mustExport("post", mi.defaultClient.getMethodClosure(http.MethodPost))
	mustExport("put", mi.defaultClient.getMethodClosure(http.MethodPut))
	mustExport("patch", mi.defaultClient.getMethodClosure(http.MethodPatch))
	mustExport("del", mi.defaultClient.getMethodClosure(http.MethodDelete))
	mustExport("options", mi.defaultClient.getMethodClosure(http.MethodOptions))
	mustExport("request", mi.defaultClient.Request)
	mustExport("batch", mi.defaultClient.Batch)
	mustExport("setResponseCallback", mi.defaultClient.SetResponseCallback)

	mustExport("expectedStatuses", mi.expectedStatuses) // TODO: refactor?

	// TODO: actually expose the default client as k6/http.defaultClient when we
	// have a better HTTP API (e.g. proper Client constructor, an actual Request
	// object, custom Transport implementations you can pass the Client, etc.).
	// This will allow us to find solutions to many of the issues with the
	// current HTTP API that plague us:
	// https://github.com/grafana/k6/issues?q=is%3Aopen+is%3Aissue+label%3Anew-http

	return mi
***REMOVED***

// Exports returns the JS values this module exports.
func (mi *ModuleInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Default: mi.exports,
		// TODO: add new HTTP APIs like Client, Request (see above comment in
		// NewModuleInstance()), etc. as named exports?
	***REMOVED***
***REMOVED***

func (mi *ModuleInstance) defineConstants() ***REMOVED***
	rt := mi.vu.Runtime()
	mustAddProp := func(name, val string) ***REMOVED***
		err := mi.exports.DefineDataProperty(
			name, rt.ToValue(val), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
		)
		if err != nil ***REMOVED***
			common.Throw(rt, err)
		***REMOVED***
	***REMOVED***
	mustAddProp("TLS_1_0", netext.TLS_1_0)
	mustAddProp("TLS_1_1", netext.TLS_1_1)
	mustAddProp("TLS_1_2", netext.TLS_1_2)
	mustAddProp("TLS_1_3", netext.TLS_1_3)
	mustAddProp("OCSP_STATUS_GOOD", netext.OCSP_STATUS_GOOD)
	mustAddProp("OCSP_STATUS_REVOKED", netext.OCSP_STATUS_REVOKED)
	mustAddProp("OCSP_STATUS_SERVER_FAILED", netext.OCSP_STATUS_SERVER_FAILED)
	mustAddProp("OCSP_STATUS_UNKNOWN", netext.OCSP_STATUS_UNKNOWN)
	mustAddProp("OCSP_REASON_UNSPECIFIED", netext.OCSP_REASON_UNSPECIFIED)
	mustAddProp("OCSP_REASON_KEY_COMPROMISE", netext.OCSP_REASON_KEY_COMPROMISE)
	mustAddProp("OCSP_REASON_CA_COMPROMISE", netext.OCSP_REASON_CA_COMPROMISE)
	mustAddProp("OCSP_REASON_AFFILIATION_CHANGED", netext.OCSP_REASON_AFFILIATION_CHANGED)
	mustAddProp("OCSP_REASON_SUPERSEDED", netext.OCSP_REASON_SUPERSEDED)
	mustAddProp("OCSP_REASON_CESSATION_OF_OPERATION", netext.OCSP_REASON_CESSATION_OF_OPERATION)
	mustAddProp("OCSP_REASON_CERTIFICATE_HOLD", netext.OCSP_REASON_CERTIFICATE_HOLD)
	mustAddProp("OCSP_REASON_REMOVE_FROM_CRL", netext.OCSP_REASON_REMOVE_FROM_CRL)
	mustAddProp("OCSP_REASON_PRIVILEGE_WITHDRAWN", netext.OCSP_REASON_PRIVILEGE_WITHDRAWN)
	mustAddProp("OCSP_REASON_AA_COMPROMISE", netext.OCSP_REASON_AA_COMPROMISE)
***REMOVED***

func (mi *ModuleInstance) newCookieJar(call goja.ConstructorCall) *goja.Object ***REMOVED***
	rt := mi.vu.Runtime()
	jar, err := cookiejar.New(nil)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	return rt.ToValue(&CookieJar***REMOVED***mi, jar***REMOVED***).ToObject(rt)
***REMOVED***

// getVUCookieJar returns the active cookie jar for the current VU.
func (mi *ModuleInstance) getVUCookieJar(call goja.FunctionCall) goja.Value ***REMOVED***
	rt := mi.vu.Runtime()
	if state := mi.vu.State(); state != nil ***REMOVED***
		return rt.ToValue(&CookieJar***REMOVED***mi, state.CookieJar***REMOVED***)
	***REMOVED***
	common.Throw(rt, ErrJarForbiddenInInitContext)
	return nil
***REMOVED***

// URL creates a new URL wrapper from the provided parts.
func (mi *ModuleInstance) URL(parts []string, pieces ...string) (httpext.URL, error) ***REMOVED***
	var name, urlstr string
	for i, part := range parts ***REMOVED***
		name += part
		urlstr += part
		if i < len(pieces) ***REMOVED***
			name += "$***REMOVED******REMOVED***"
			urlstr += pieces[i]
		***REMOVED***
	***REMOVED***
	return httpext.NewURL(urlstr, name)
***REMOVED***

// Client represents a stand-alone HTTP client.
//
// TODO: move to its own file
type Client struct ***REMOVED***
	moduleInstance   *ModuleInstance
	responseCallback func(int) bool
***REMOVED***
