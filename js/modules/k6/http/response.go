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

package http

import (
	"context"
	"crypto/tls"
	"encoding/json"

	"fmt"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/modules/k6/html"
	"github.com/loadimpact/k6/lib"
	"golang.org/x/crypto/ocsp"
	"net/url"
	"strings"
)

type OCSP struct ***REMOVED***
	ProducedAt, ThisUpdate, NextUpdate, RevokedAt int64
	RevocationReason                              string
	Status                                        string
***REMOVED***

type HTTPResponseTimings struct ***REMOVED***
	Duration, Blocked, LookingUp, Connecting, Sending, Waiting, Receiving float64
***REMOVED***

type HTTPResponse struct ***REMOVED***
	ctx context.Context

	RemoteIP       string
	RemotePort     int
	URL            string
	Status         int
	Proto          string
	Headers        map[string]string
	Cookies        map[string][]*HTTPCookie
	Body           string
	Timings        HTTPResponseTimings
	TLSVersion     string
	TLSCipherSuite string
	OCSP           OCSP `js:"ocsp"`
	Error          string

	cachedJSON goja.Value
***REMOVED***

func (res *HTTPResponse) setTLSInfo(tlsState *tls.ConnectionState) ***REMOVED***
	switch tlsState.Version ***REMOVED***
	case tls.VersionSSL30:
		res.TLSVersion = SSL_3_0
	case tls.VersionTLS10:
		res.TLSVersion = TLS_1_0
	case tls.VersionTLS11:
		res.TLSVersion = TLS_1_1
	case tls.VersionTLS12:
		res.TLSVersion = TLS_1_2
	***REMOVED***

	res.TLSCipherSuite = lib.SupportedTLSCipherSuitesToString[tlsState.CipherSuite]
	ocspStapledRes := OCSP***REMOVED***Status: OCSP_STATUS_UNKNOWN***REMOVED***

	if ocspRes, err := ocsp.ParseResponse(tlsState.OCSPResponse, nil); err == nil ***REMOVED***
		switch ocspRes.Status ***REMOVED***
		case ocsp.Good:
			ocspStapledRes.Status = OCSP_STATUS_GOOD
		case ocsp.Revoked:
			ocspStapledRes.Status = OCSP_STATUS_REVOKED
		case ocsp.ServerFailed:
			ocspStapledRes.Status = OCSP_STATUS_SERVER_FAILED
		case ocsp.Unknown:
			ocspStapledRes.Status = OCSP_STATUS_UNKNOWN
		***REMOVED***
		switch ocspRes.RevocationReason ***REMOVED***
		case ocsp.Unspecified:
			ocspStapledRes.RevocationReason = OCSP_REASON_UNSPECIFIED
		case ocsp.KeyCompromise:
			ocspStapledRes.RevocationReason = OCSP_REASON_KEY_COMPROMISE
		case ocsp.CACompromise:
			ocspStapledRes.RevocationReason = OCSP_REASON_CA_COMPROMISE
		case ocsp.AffiliationChanged:
			ocspStapledRes.RevocationReason = OCSP_REASON_AFFILIATION_CHANGED
		case ocsp.Superseded:
			ocspStapledRes.RevocationReason = OCSP_REASON_SUPERSEDED
		case ocsp.CessationOfOperation:
			ocspStapledRes.RevocationReason = OCSP_REASON_CESSATION_OF_OPERATION
		case ocsp.CertificateHold:
			ocspStapledRes.RevocationReason = OCSP_REASON_CERTIFICATE_HOLD
		case ocsp.RemoveFromCRL:
			ocspStapledRes.RevocationReason = OCSP_REASON_REMOVE_FROM_CRL
		case ocsp.PrivilegeWithdrawn:
			ocspStapledRes.RevocationReason = OCSP_REASON_PRIVILEGE_WITHDRAWN
		case ocsp.AACompromise:
			ocspStapledRes.RevocationReason = OCSP_REASON_AA_COMPROMISE
		***REMOVED***
		ocspStapledRes.ProducedAt = ocspRes.ProducedAt.Unix()
		ocspStapledRes.ThisUpdate = ocspRes.ThisUpdate.Unix()
		ocspStapledRes.NextUpdate = ocspRes.NextUpdate.Unix()
		ocspStapledRes.RevokedAt = ocspRes.RevokedAt.Unix()
	***REMOVED***

	res.OCSP = ocspStapledRes
***REMOVED***

func (res *HTTPResponse) Json() goja.Value ***REMOVED***
	if res.cachedJSON == nil ***REMOVED***
		var v interface***REMOVED******REMOVED***
		if err := json.Unmarshal([]byte(res.Body), &v); err != nil ***REMOVED***
			common.Throw(common.GetRuntime(res.ctx), err)
		***REMOVED***
		res.cachedJSON = common.GetRuntime(res.ctx).ToValue(v)
	***REMOVED***
	return res.cachedJSON
***REMOVED***

func (res *HTTPResponse) Html(selector ...string) html.Selection ***REMOVED***
	sel, err := html.HTML***REMOVED******REMOVED***.ParseHTML(res.ctx, res.Body)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(res.ctx), err)
	***REMOVED***
	sel.URL = res.URL
	if len(selector) > 0 ***REMOVED***
		sel = sel.Find(selector[0])
	***REMOVED***
	return sel
***REMOVED***

func (res *HTTPResponse) SubmitForm(args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	rt := common.GetRuntime(res.ctx)

	formSelector := "form"
	submitSelector := "[type=\"submit\"]"
	var fields map[string]goja.Value
	var requestOptions goja.Value
	if len(args) > 0 ***REMOVED***
		params := args[0].ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch k ***REMOVED***
			case "formSelector":
				formSelector = params.Get(k).String()
			case "submitSelector":
				formSelector = params.Get(k).String()
			case "fields":
				rt.ExportTo(params.Get(k), &fields)
			case "options":
				requestOptions = params.Get(k)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	form := res.Html(formSelector)
	if form.Size() == 0 ***REMOVED***
		common.Throw(rt, fmt.Errorf("no form found for selector '%s' in response '%s'", formSelector, res.URL))
	***REMOVED***

	methodAttr := form.Attr("method")
	var requestMethod string
	if methodAttr == goja.Undefined() ***REMOVED***
		// Use GET by default
		requestMethod = "GET" // nolint: goconst
	***REMOVED*** else ***REMOVED***
		requestMethod = strings.ToUpper(methodAttr.String())
	***REMOVED***

	actionAttr := form.Attr("action")
	var requestUrl goja.Value
	if actionAttr == goja.Undefined() ***REMOVED***
		// Use the url of the response if no action is set
		requestUrl = rt.ToValue(res.URL)
	***REMOVED*** else ***REMOVED***
		// Resolve the action url from the response url
		responseUrl, _ := url.Parse(res.URL)
		actionUrl, _ := url.Parse(actionAttr.String())
		requestUrl = rt.ToValue(responseUrl.ResolveReference(actionUrl).String())
	***REMOVED***

	// Set the body based on the form values
	body := form.SerializeObject()

	// Set the name + value of the submit button
	submit := form.Find(submitSelector)
	submitName := submit.Attr("name")
	submitValue := submit.Val()
	if submitName != goja.Undefined() && submitValue != goja.Undefined() ***REMOVED***
		body[submitName.String()] = submitValue
	***REMOVED***

	// Set the values supplied in the arguments, overriding automatically set values
	if fields != nil ***REMOVED***
		for k, v := range fields ***REMOVED***
			body[k] = v
		***REMOVED***
	***REMOVED***

	if requestOptions == nil ***REMOVED***
		return New().Request(res.ctx, requestMethod, requestUrl, rt.ToValue(body))
	***REMOVED*** else ***REMOVED***
		return New().Request(res.ctx, requestMethod, requestUrl, rt.ToValue(body), requestOptions)
	***REMOVED***
***REMOVED***
