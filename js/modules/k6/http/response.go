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
	"errors"

	"fmt"
	"net/url"
	"strings"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/modules/k6/html"
	"github.com/loadimpact/k6/lib"
	"golang.org/x/crypto/ocsp"
)

type OCSP struct ***REMOVED***
	ProducedAt, ThisUpdate, NextUpdate, RevokedAt int64
	RevocationReason                              string
	Status                                        string
***REMOVED***

type HTTPResponseTimings struct ***REMOVED***
	Duration, Blocked, LookingUp, Connecting, TLSHandshaking, Sending, Waiting, Receiving float64
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
	Body           interface***REMOVED******REMOVED***
	Timings        HTTPResponseTimings
	TLSVersion     string
	TLSCipherSuite string
	OCSP           OCSP `js:"ocsp"`
	Error          string
	Request        HTTPRequest

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
		var body []byte
		switch b := res.Body.(type) ***REMOVED***
		case []byte:
			body = b
		case string:
			body = []byte(b)
		default:
			common.Throw(common.GetRuntime(res.ctx), errors.New("Invalid response type"))
		***REMOVED***
		if err := json.Unmarshal(body, &v); err != nil ***REMOVED***
			common.Throw(common.GetRuntime(res.ctx), err)
		***REMOVED***
		res.cachedJSON = common.GetRuntime(res.ctx).ToValue(v)
	***REMOVED***
	return res.cachedJSON
***REMOVED***

func (res *HTTPResponse) Html(selector ...string) html.Selection ***REMOVED***
	var body string
	switch b := res.Body.(type) ***REMOVED***
	case []byte:
		body = string(b)
	case string:
		body = b
	default:
		common.Throw(common.GetRuntime(res.ctx), errors.New("Invalid response type"))
	***REMOVED***

	sel, err := html.HTML***REMOVED******REMOVED***.ParseHTML(res.ctx, body)
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
	requestParams := goja.Null()
	if len(args) > 0 ***REMOVED***
		params := args[0].ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch k ***REMOVED***
			case "formSelector":
				formSelector = params.Get(k).String()
			case "submitSelector":
				submitSelector = params.Get(k).String()
			case "fields":
				if rt.ExportTo(params.Get(k), &fields) != nil ***REMOVED***
					fields = nil
				***REMOVED***
			case "params":
				requestParams = params.Get(k)
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
		requestMethod = HTTP_METHOD_GET
	***REMOVED*** else ***REMOVED***
		requestMethod = strings.ToUpper(methodAttr.String())
	***REMOVED***

	responseUrl, err := url.Parse(res.URL)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	actionAttr := form.Attr("action")
	var requestUrl *url.URL
	if actionAttr == goja.Undefined() ***REMOVED***
		// Use the url of the response if no action is set
		requestUrl = responseUrl
	***REMOVED*** else ***REMOVED***
		actionUrl, err := url.Parse(actionAttr.String())
		if err != nil ***REMOVED***
			common.Throw(rt, err)
		***REMOVED***
		requestUrl = responseUrl.ResolveReference(actionUrl)
	***REMOVED***

	// Set the body based on the form values
	values := form.SerializeObject()

	// Set the name + value of the submit button
	submit := form.Find(submitSelector)
	submitName := submit.Attr("name")
	submitValue := submit.Val()
	if submitName != goja.Undefined() && submitValue != goja.Undefined() ***REMOVED***
		values[submitName.String()] = submitValue
	***REMOVED***

	// Set the values supplied in the arguments, overriding automatically set values
	for k, v := range fields ***REMOVED***
		values[k] = v
	***REMOVED***

	if requestMethod == HTTP_METHOD_GET ***REMOVED***
		q := url.Values***REMOVED******REMOVED***
		for k, v := range values ***REMOVED***
			q.Add(k, v.String())
		***REMOVED***
		requestUrl.RawQuery = q.Encode()
		return New().Request(res.ctx, requestMethod, rt.ToValue(requestUrl.String()), goja.Null(), requestParams)
	***REMOVED***
	return New().Request(res.ctx, requestMethod, rt.ToValue(requestUrl.String()), rt.ToValue(values), requestParams)
***REMOVED***

func (res *HTTPResponse) ClickLink(args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	rt := common.GetRuntime(res.ctx)

	selector := "a[href]"
	requestParams := goja.Null()
	if len(args) > 0 ***REMOVED***
		params := args[0].ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch k ***REMOVED***
			case "selector":
				selector = params.Get(k).String()
			case "params":
				requestParams = params.Get(k)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	responseUrl, err := url.Parse(res.URL)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	link := res.Html(selector)
	if link.Size() == 0 ***REMOVED***
		common.Throw(rt, fmt.Errorf("no element found for selector '%s' in response '%s'", selector, res.URL))
	***REMOVED***
	hrefAttr := link.Attr("href")
	if hrefAttr == goja.Undefined() ***REMOVED***
		common.Throw(rt, fmt.Errorf("no valid href attribute value found on element '%s' in response '%s'", selector, res.URL))
	***REMOVED***
	hrefUrl, err := url.Parse(hrefAttr.String())
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	requestUrl := responseUrl.ResolveReference(hrefUrl)

	return New().Get(res.ctx, rt.ToValue(requestUrl.String()), requestParams)
***REMOVED***
