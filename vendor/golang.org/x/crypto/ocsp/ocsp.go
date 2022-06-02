// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ocsp parses OCSP responses as specified in RFC 2560. OCSP responses
// are signed messages attesting to the validity of a certificate for a small
// period of time. This is used to manage revocation for X.509 certificates.
package ocsp // import "golang.org/x/crypto/ocsp"

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

var idPKIXOCSPBasic = asn1.ObjectIdentifier([]int***REMOVED***1, 3, 6, 1, 5, 5, 7, 48, 1, 1***REMOVED***)

// ResponseStatus contains the result of an OCSP request. See
// https://tools.ietf.org/html/rfc6960#section-2.3
type ResponseStatus int

const (
	Success       ResponseStatus = 0
	Malformed     ResponseStatus = 1
	InternalError ResponseStatus = 2
	TryLater      ResponseStatus = 3
	// Status code four is unused in OCSP. See
	// https://tools.ietf.org/html/rfc6960#section-4.2.1
	SignatureRequired ResponseStatus = 5
	Unauthorized      ResponseStatus = 6
)

func (r ResponseStatus) String() string ***REMOVED***
	switch r ***REMOVED***
	case Success:
		return "success"
	case Malformed:
		return "malformed"
	case InternalError:
		return "internal error"
	case TryLater:
		return "try later"
	case SignatureRequired:
		return "signature required"
	case Unauthorized:
		return "unauthorized"
	default:
		return "unknown OCSP status: " + strconv.Itoa(int(r))
	***REMOVED***
***REMOVED***

// ResponseError is an error that may be returned by ParseResponse to indicate
// that the response itself is an error, not just that it's indicating that a
// certificate is revoked, unknown, etc.
type ResponseError struct ***REMOVED***
	Status ResponseStatus
***REMOVED***

func (r ResponseError) Error() string ***REMOVED***
	return "ocsp: error from server: " + r.Status.String()
***REMOVED***

// These are internal structures that reflect the ASN.1 structure of an OCSP
// response. See RFC 2560, section 4.2.

type certID struct ***REMOVED***
	HashAlgorithm pkix.AlgorithmIdentifier
	NameHash      []byte
	IssuerKeyHash []byte
	SerialNumber  *big.Int
***REMOVED***

// https://tools.ietf.org/html/rfc2560#section-4.1.1
type ocspRequest struct ***REMOVED***
	TBSRequest tbsRequest
***REMOVED***

type tbsRequest struct ***REMOVED***
	Version       int              `asn1:"explicit,tag:0,default:0,optional"`
	RequestorName pkix.RDNSequence `asn1:"explicit,tag:1,optional"`
	RequestList   []request
***REMOVED***

type request struct ***REMOVED***
	Cert certID
***REMOVED***

type responseASN1 struct ***REMOVED***
	Status   asn1.Enumerated
	Response responseBytes `asn1:"explicit,tag:0,optional"`
***REMOVED***

type responseBytes struct ***REMOVED***
	ResponseType asn1.ObjectIdentifier
	Response     []byte
***REMOVED***

type basicResponse struct ***REMOVED***
	TBSResponseData    responseData
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Signature          asn1.BitString
	Certificates       []asn1.RawValue `asn1:"explicit,tag:0,optional"`
***REMOVED***

type responseData struct ***REMOVED***
	Raw            asn1.RawContent
	Version        int `asn1:"optional,default:0,explicit,tag:0"`
	RawResponderID asn1.RawValue
	ProducedAt     time.Time `asn1:"generalized"`
	Responses      []singleResponse
***REMOVED***

type singleResponse struct ***REMOVED***
	CertID           certID
	Good             asn1.Flag        `asn1:"tag:0,optional"`
	Revoked          revokedInfo      `asn1:"tag:1,optional"`
	Unknown          asn1.Flag        `asn1:"tag:2,optional"`
	ThisUpdate       time.Time        `asn1:"generalized"`
	NextUpdate       time.Time        `asn1:"generalized,explicit,tag:0,optional"`
	SingleExtensions []pkix.Extension `asn1:"explicit,tag:1,optional"`
***REMOVED***

type revokedInfo struct ***REMOVED***
	RevocationTime time.Time       `asn1:"generalized"`
	Reason         asn1.Enumerated `asn1:"explicit,tag:0,optional"`
***REMOVED***

var (
	oidSignatureMD2WithRSA      = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 2***REMOVED***
	oidSignatureMD5WithRSA      = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 4***REMOVED***
	oidSignatureSHA1WithRSA     = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 5***REMOVED***
	oidSignatureSHA256WithRSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 11***REMOVED***
	oidSignatureSHA384WithRSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 12***REMOVED***
	oidSignatureSHA512WithRSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 13***REMOVED***
	oidSignatureDSAWithSHA1     = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10040, 4, 3***REMOVED***
	oidSignatureDSAWithSHA256   = asn1.ObjectIdentifier***REMOVED***2, 16, 840, 1, 101, 3, 4, 3, 2***REMOVED***
	oidSignatureECDSAWithSHA1   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 4, 1***REMOVED***
	oidSignatureECDSAWithSHA256 = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 4, 3, 2***REMOVED***
	oidSignatureECDSAWithSHA384 = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 4, 3, 3***REMOVED***
	oidSignatureECDSAWithSHA512 = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 4, 3, 4***REMOVED***
)

var hashOIDs = map[crypto.Hash]asn1.ObjectIdentifier***REMOVED***
	crypto.SHA1:   asn1.ObjectIdentifier([]int***REMOVED***1, 3, 14, 3, 2, 26***REMOVED***),
	crypto.SHA256: asn1.ObjectIdentifier([]int***REMOVED***2, 16, 840, 1, 101, 3, 4, 2, 1***REMOVED***),
	crypto.SHA384: asn1.ObjectIdentifier([]int***REMOVED***2, 16, 840, 1, 101, 3, 4, 2, 2***REMOVED***),
	crypto.SHA512: asn1.ObjectIdentifier([]int***REMOVED***2, 16, 840, 1, 101, 3, 4, 2, 3***REMOVED***),
***REMOVED***

// TODO(rlb): This is also from crypto/x509, so same comment as AGL's below
var signatureAlgorithmDetails = []struct ***REMOVED***
	algo       x509.SignatureAlgorithm
	oid        asn1.ObjectIdentifier
	pubKeyAlgo x509.PublicKeyAlgorithm
	hash       crypto.Hash
***REMOVED******REMOVED***
	***REMOVED***x509.MD2WithRSA, oidSignatureMD2WithRSA, x509.RSA, crypto.Hash(0) /* no value for MD2 */***REMOVED***,
	***REMOVED***x509.MD5WithRSA, oidSignatureMD5WithRSA, x509.RSA, crypto.MD5***REMOVED***,
	***REMOVED***x509.SHA1WithRSA, oidSignatureSHA1WithRSA, x509.RSA, crypto.SHA1***REMOVED***,
	***REMOVED***x509.SHA256WithRSA, oidSignatureSHA256WithRSA, x509.RSA, crypto.SHA256***REMOVED***,
	***REMOVED***x509.SHA384WithRSA, oidSignatureSHA384WithRSA, x509.RSA, crypto.SHA384***REMOVED***,
	***REMOVED***x509.SHA512WithRSA, oidSignatureSHA512WithRSA, x509.RSA, crypto.SHA512***REMOVED***,
	***REMOVED***x509.DSAWithSHA1, oidSignatureDSAWithSHA1, x509.DSA, crypto.SHA1***REMOVED***,
	***REMOVED***x509.DSAWithSHA256, oidSignatureDSAWithSHA256, x509.DSA, crypto.SHA256***REMOVED***,
	***REMOVED***x509.ECDSAWithSHA1, oidSignatureECDSAWithSHA1, x509.ECDSA, crypto.SHA1***REMOVED***,
	***REMOVED***x509.ECDSAWithSHA256, oidSignatureECDSAWithSHA256, x509.ECDSA, crypto.SHA256***REMOVED***,
	***REMOVED***x509.ECDSAWithSHA384, oidSignatureECDSAWithSHA384, x509.ECDSA, crypto.SHA384***REMOVED***,
	***REMOVED***x509.ECDSAWithSHA512, oidSignatureECDSAWithSHA512, x509.ECDSA, crypto.SHA512***REMOVED***,
***REMOVED***

// TODO(rlb): This is also from crypto/x509, so same comment as AGL's below
func signingParamsForPublicKey(pub interface***REMOVED******REMOVED***, requestedSigAlgo x509.SignatureAlgorithm) (hashFunc crypto.Hash, sigAlgo pkix.AlgorithmIdentifier, err error) ***REMOVED***
	var pubType x509.PublicKeyAlgorithm

	switch pub := pub.(type) ***REMOVED***
	case *rsa.PublicKey:
		pubType = x509.RSA
		hashFunc = crypto.SHA256
		sigAlgo.Algorithm = oidSignatureSHA256WithRSA
		sigAlgo.Parameters = asn1.RawValue***REMOVED***
			Tag: 5,
		***REMOVED***

	case *ecdsa.PublicKey:
		pubType = x509.ECDSA

		switch pub.Curve ***REMOVED***
		case elliptic.P224(), elliptic.P256():
			hashFunc = crypto.SHA256
			sigAlgo.Algorithm = oidSignatureECDSAWithSHA256
		case elliptic.P384():
			hashFunc = crypto.SHA384
			sigAlgo.Algorithm = oidSignatureECDSAWithSHA384
		case elliptic.P521():
			hashFunc = crypto.SHA512
			sigAlgo.Algorithm = oidSignatureECDSAWithSHA512
		default:
			err = errors.New("x509: unknown elliptic curve")
		***REMOVED***

	default:
		err = errors.New("x509: only RSA and ECDSA keys supported")
	***REMOVED***

	if err != nil ***REMOVED***
		return
	***REMOVED***

	if requestedSigAlgo == 0 ***REMOVED***
		return
	***REMOVED***

	found := false
	for _, details := range signatureAlgorithmDetails ***REMOVED***
		if details.algo == requestedSigAlgo ***REMOVED***
			if details.pubKeyAlgo != pubType ***REMOVED***
				err = errors.New("x509: requested SignatureAlgorithm does not match private key type")
				return
			***REMOVED***
			sigAlgo.Algorithm, hashFunc = details.oid, details.hash
			if hashFunc == 0 ***REMOVED***
				err = errors.New("x509: cannot sign with hash function requested")
				return
			***REMOVED***
			found = true
			break
		***REMOVED***
	***REMOVED***

	if !found ***REMOVED***
		err = errors.New("x509: unknown SignatureAlgorithm")
	***REMOVED***

	return
***REMOVED***

// TODO(agl): this is taken from crypto/x509 and so should probably be exported
// from crypto/x509 or crypto/x509/pkix.
func getSignatureAlgorithmFromOID(oid asn1.ObjectIdentifier) x509.SignatureAlgorithm ***REMOVED***
	for _, details := range signatureAlgorithmDetails ***REMOVED***
		if oid.Equal(details.oid) ***REMOVED***
			return details.algo
		***REMOVED***
	***REMOVED***
	return x509.UnknownSignatureAlgorithm
***REMOVED***

// TODO(rlb): This is not taken from crypto/x509, but it's of the same general form.
func getHashAlgorithmFromOID(target asn1.ObjectIdentifier) crypto.Hash ***REMOVED***
	for hash, oid := range hashOIDs ***REMOVED***
		if oid.Equal(target) ***REMOVED***
			return hash
		***REMOVED***
	***REMOVED***
	return crypto.Hash(0)
***REMOVED***

func getOIDFromHashAlgorithm(target crypto.Hash) asn1.ObjectIdentifier ***REMOVED***
	for hash, oid := range hashOIDs ***REMOVED***
		if hash == target ***REMOVED***
			return oid
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// This is the exposed reflection of the internal OCSP structures.

// The status values that can be expressed in OCSP.  See RFC 6960.
const (
	// Good means that the certificate is valid.
	Good = iota
	// Revoked means that the certificate has been deliberately revoked.
	Revoked
	// Unknown means that the OCSP responder doesn't know about the certificate.
	Unknown
	// ServerFailed is unused and was never used (see
	// https://go-review.googlesource.com/#/c/18944). ParseResponse will
	// return a ResponseError when an error response is parsed.
	ServerFailed
)

// The enumerated reasons for revoking a certificate.  See RFC 5280.
const (
	Unspecified          = 0
	KeyCompromise        = 1
	CACompromise         = 2
	AffiliationChanged   = 3
	Superseded           = 4
	CessationOfOperation = 5
	CertificateHold      = 6

	RemoveFromCRL      = 8
	PrivilegeWithdrawn = 9
	AACompromise       = 10
)

// Request represents an OCSP request. See RFC 6960.
type Request struct ***REMOVED***
	HashAlgorithm  crypto.Hash
	IssuerNameHash []byte
	IssuerKeyHash  []byte
	SerialNumber   *big.Int
***REMOVED***

// Marshal marshals the OCSP request to ASN.1 DER encoded form.
func (req *Request) Marshal() ([]byte, error) ***REMOVED***
	hashAlg := getOIDFromHashAlgorithm(req.HashAlgorithm)
	if hashAlg == nil ***REMOVED***
		return nil, errors.New("Unknown hash algorithm")
	***REMOVED***
	return asn1.Marshal(ocspRequest***REMOVED***
		tbsRequest***REMOVED***
			Version: 0,
			RequestList: []request***REMOVED***
				***REMOVED***
					Cert: certID***REMOVED***
						pkix.AlgorithmIdentifier***REMOVED***
							Algorithm:  hashAlg,
							Parameters: asn1.RawValue***REMOVED***Tag: 5 /* ASN.1 NULL */***REMOVED***,
						***REMOVED***,
						req.IssuerNameHash,
						req.IssuerKeyHash,
						req.SerialNumber,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// Response represents an OCSP response containing a single SingleResponse. See
// RFC 6960.
type Response struct ***REMOVED***
	Raw []byte

	// Status is one of ***REMOVED***Good, Revoked, Unknown***REMOVED***
	Status                                        int
	SerialNumber                                  *big.Int
	ProducedAt, ThisUpdate, NextUpdate, RevokedAt time.Time
	RevocationReason                              int
	Certificate                                   *x509.Certificate
	// TBSResponseData contains the raw bytes of the signed response. If
	// Certificate is nil then this can be used to verify Signature.
	TBSResponseData    []byte
	Signature          []byte
	SignatureAlgorithm x509.SignatureAlgorithm

	// IssuerHash is the hash used to compute the IssuerNameHash and IssuerKeyHash.
	// Valid values are crypto.SHA1, crypto.SHA256, crypto.SHA384, and crypto.SHA512.
	// If zero, the default is crypto.SHA1.
	IssuerHash crypto.Hash

	// RawResponderName optionally contains the DER-encoded subject of the
	// responder certificate. Exactly one of RawResponderName and
	// ResponderKeyHash is set.
	RawResponderName []byte
	// ResponderKeyHash optionally contains the SHA-1 hash of the
	// responder's public key. Exactly one of RawResponderName and
	// ResponderKeyHash is set.
	ResponderKeyHash []byte

	// Extensions contains raw X.509 extensions from the singleExtensions field
	// of the OCSP response. When parsing certificates, this can be used to
	// extract non-critical extensions that are not parsed by this package. When
	// marshaling OCSP responses, the Extensions field is ignored, see
	// ExtraExtensions.
	Extensions []pkix.Extension

	// ExtraExtensions contains extensions to be copied, raw, into any marshaled
	// OCSP response (in the singleExtensions field). Values override any
	// extensions that would otherwise be produced based on the other fields. The
	// ExtraExtensions field is not populated when parsing certificates, see
	// Extensions.
	ExtraExtensions []pkix.Extension
***REMOVED***

// These are pre-serialized error responses for the various non-success codes
// defined by OCSP. The Unauthorized code in particular can be used by an OCSP
// responder that supports only pre-signed responses as a response to requests
// for certificates with unknown status. See RFC 5019.
var (
	MalformedRequestErrorResponse = []byte***REMOVED***0x30, 0x03, 0x0A, 0x01, 0x01***REMOVED***
	InternalErrorErrorResponse    = []byte***REMOVED***0x30, 0x03, 0x0A, 0x01, 0x02***REMOVED***
	TryLaterErrorResponse         = []byte***REMOVED***0x30, 0x03, 0x0A, 0x01, 0x03***REMOVED***
	SigRequredErrorResponse       = []byte***REMOVED***0x30, 0x03, 0x0A, 0x01, 0x05***REMOVED***
	UnauthorizedErrorResponse     = []byte***REMOVED***0x30, 0x03, 0x0A, 0x01, 0x06***REMOVED***
)

// CheckSignatureFrom checks that the signature in resp is a valid signature
// from issuer. This should only be used if resp.Certificate is nil. Otherwise,
// the OCSP response contained an intermediate certificate that created the
// signature. That signature is checked by ParseResponse and only
// resp.Certificate remains to be validated.
func (resp *Response) CheckSignatureFrom(issuer *x509.Certificate) error ***REMOVED***
	return issuer.CheckSignature(resp.SignatureAlgorithm, resp.TBSResponseData, resp.Signature)
***REMOVED***

// ParseError results from an invalid OCSP response.
type ParseError string

func (p ParseError) Error() string ***REMOVED***
	return string(p)
***REMOVED***

// ParseRequest parses an OCSP request in DER form. It only supports
// requests for a single certificate. Signed requests are not supported.
// If a request includes a signature, it will result in a ParseError.
func ParseRequest(bytes []byte) (*Request, error) ***REMOVED***
	var req ocspRequest
	rest, err := asn1.Unmarshal(bytes, &req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(rest) > 0 ***REMOVED***
		return nil, ParseError("trailing data in OCSP request")
	***REMOVED***

	if len(req.TBSRequest.RequestList) == 0 ***REMOVED***
		return nil, ParseError("OCSP request contains no request body")
	***REMOVED***
	innerRequest := req.TBSRequest.RequestList[0]

	hashFunc := getHashAlgorithmFromOID(innerRequest.Cert.HashAlgorithm.Algorithm)
	if hashFunc == crypto.Hash(0) ***REMOVED***
		return nil, ParseError("OCSP request uses unknown hash function")
	***REMOVED***

	return &Request***REMOVED***
		HashAlgorithm:  hashFunc,
		IssuerNameHash: innerRequest.Cert.NameHash,
		IssuerKeyHash:  innerRequest.Cert.IssuerKeyHash,
		SerialNumber:   innerRequest.Cert.SerialNumber,
	***REMOVED***, nil
***REMOVED***

// ParseResponse parses an OCSP response in DER form. The response must contain
// only one certificate status. To parse the status of a specific certificate
// from a response which may contain multiple statuses, use ParseResponseForCert
// instead.
//
// If the response contains an embedded certificate, then that certificate will
// be used to verify the response signature. If the response contains an
// embedded certificate and issuer is not nil, then issuer will be used to verify
// the signature on the embedded certificate.
//
// If the response does not contain an embedded certificate and issuer is not
// nil, then issuer will be used to verify the response signature.
//
// Invalid responses and parse failures will result in a ParseError.
// Error responses will result in a ResponseError.
func ParseResponse(bytes []byte, issuer *x509.Certificate) (*Response, error) ***REMOVED***
	return ParseResponseForCert(bytes, nil, issuer)
***REMOVED***

// ParseResponseForCert acts identically to ParseResponse, except it supports
// parsing responses that contain multiple statuses. If the response contains
// multiple statuses and cert is not nil, then ParseResponseForCert will return
// the first status which contains a matching serial, otherwise it will return an
// error. If cert is nil, then the first status in the response will be returned.
func ParseResponseForCert(bytes []byte, cert, issuer *x509.Certificate) (*Response, error) ***REMOVED***
	var resp responseASN1
	rest, err := asn1.Unmarshal(bytes, &resp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(rest) > 0 ***REMOVED***
		return nil, ParseError("trailing data in OCSP response")
	***REMOVED***

	if status := ResponseStatus(resp.Status); status != Success ***REMOVED***
		return nil, ResponseError***REMOVED***status***REMOVED***
	***REMOVED***

	if !resp.Response.ResponseType.Equal(idPKIXOCSPBasic) ***REMOVED***
		return nil, ParseError("bad OCSP response type")
	***REMOVED***

	var basicResp basicResponse
	rest, err = asn1.Unmarshal(resp.Response.Response, &basicResp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(rest) > 0 ***REMOVED***
		return nil, ParseError("trailing data in OCSP response")
	***REMOVED***

	if n := len(basicResp.TBSResponseData.Responses); n == 0 || cert == nil && n > 1 ***REMOVED***
		return nil, ParseError("OCSP response contains bad number of responses")
	***REMOVED***

	var singleResp singleResponse
	if cert == nil ***REMOVED***
		singleResp = basicResp.TBSResponseData.Responses[0]
	***REMOVED*** else ***REMOVED***
		match := false
		for _, resp := range basicResp.TBSResponseData.Responses ***REMOVED***
			if cert.SerialNumber.Cmp(resp.CertID.SerialNumber) == 0 ***REMOVED***
				singleResp = resp
				match = true
				break
			***REMOVED***
		***REMOVED***
		if !match ***REMOVED***
			return nil, ParseError("no response matching the supplied certificate")
		***REMOVED***
	***REMOVED***

	ret := &Response***REMOVED***
		Raw:                bytes,
		TBSResponseData:    basicResp.TBSResponseData.Raw,
		Signature:          basicResp.Signature.RightAlign(),
		SignatureAlgorithm: getSignatureAlgorithmFromOID(basicResp.SignatureAlgorithm.Algorithm),
		Extensions:         singleResp.SingleExtensions,
		SerialNumber:       singleResp.CertID.SerialNumber,
		ProducedAt:         basicResp.TBSResponseData.ProducedAt,
		ThisUpdate:         singleResp.ThisUpdate,
		NextUpdate:         singleResp.NextUpdate,
	***REMOVED***

	// Handle the ResponderID CHOICE tag. ResponderID can be flattened into
	// TBSResponseData once https://go-review.googlesource.com/34503 has been
	// released.
	rawResponderID := basicResp.TBSResponseData.RawResponderID
	switch rawResponderID.Tag ***REMOVED***
	case 1: // Name
		var rdn pkix.RDNSequence
		if rest, err := asn1.Unmarshal(rawResponderID.Bytes, &rdn); err != nil || len(rest) != 0 ***REMOVED***
			return nil, ParseError("invalid responder name")
		***REMOVED***
		ret.RawResponderName = rawResponderID.Bytes
	case 2: // KeyHash
		if rest, err := asn1.Unmarshal(rawResponderID.Bytes, &ret.ResponderKeyHash); err != nil || len(rest) != 0 ***REMOVED***
			return nil, ParseError("invalid responder key hash")
		***REMOVED***
	default:
		return nil, ParseError("invalid responder id tag")
	***REMOVED***

	if len(basicResp.Certificates) > 0 ***REMOVED***
		// Responders should only send a single certificate (if they
		// send any) that connects the responder's certificate to the
		// original issuer. We accept responses with multiple
		// certificates due to a number responders sending them[1], but
		// ignore all but the first.
		//
		// [1] https://github.com/golang/go/issues/21527
		ret.Certificate, err = x509.ParseCertificate(basicResp.Certificates[0].FullBytes)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if err := ret.CheckSignatureFrom(ret.Certificate); err != nil ***REMOVED***
			return nil, ParseError("bad signature on embedded certificate: " + err.Error())
		***REMOVED***

		if issuer != nil ***REMOVED***
			if err := issuer.CheckSignature(ret.Certificate.SignatureAlgorithm, ret.Certificate.RawTBSCertificate, ret.Certificate.Signature); err != nil ***REMOVED***
				return nil, ParseError("bad OCSP signature: " + err.Error())
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if issuer != nil ***REMOVED***
		if err := ret.CheckSignatureFrom(issuer); err != nil ***REMOVED***
			return nil, ParseError("bad OCSP signature: " + err.Error())
		***REMOVED***
	***REMOVED***

	for _, ext := range singleResp.SingleExtensions ***REMOVED***
		if ext.Critical ***REMOVED***
			return nil, ParseError("unsupported critical extension")
		***REMOVED***
	***REMOVED***

	for h, oid := range hashOIDs ***REMOVED***
		if singleResp.CertID.HashAlgorithm.Algorithm.Equal(oid) ***REMOVED***
			ret.IssuerHash = h
			break
		***REMOVED***
	***REMOVED***
	if ret.IssuerHash == 0 ***REMOVED***
		return nil, ParseError("unsupported issuer hash algorithm")
	***REMOVED***

	switch ***REMOVED***
	case bool(singleResp.Good):
		ret.Status = Good
	case bool(singleResp.Unknown):
		ret.Status = Unknown
	default:
		ret.Status = Revoked
		ret.RevokedAt = singleResp.Revoked.RevocationTime
		ret.RevocationReason = int(singleResp.Revoked.Reason)
	***REMOVED***

	return ret, nil
***REMOVED***

// RequestOptions contains options for constructing OCSP requests.
type RequestOptions struct ***REMOVED***
	// Hash contains the hash function that should be used when
	// constructing the OCSP request. If zero, SHA-1 will be used.
	Hash crypto.Hash
***REMOVED***

func (opts *RequestOptions) hash() crypto.Hash ***REMOVED***
	if opts == nil || opts.Hash == 0 ***REMOVED***
		// SHA-1 is nearly universally used in OCSP.
		return crypto.SHA1
	***REMOVED***
	return opts.Hash
***REMOVED***

// CreateRequest returns a DER-encoded, OCSP request for the status of cert. If
// opts is nil then sensible defaults are used.
func CreateRequest(cert, issuer *x509.Certificate, opts *RequestOptions) ([]byte, error) ***REMOVED***
	hashFunc := opts.hash()

	// OCSP seems to be the only place where these raw hash identifiers are
	// used. I took the following from
	// http://msdn.microsoft.com/en-us/library/ff635603.aspx
	_, ok := hashOIDs[hashFunc]
	if !ok ***REMOVED***
		return nil, x509.ErrUnsupportedAlgorithm
	***REMOVED***

	if !hashFunc.Available() ***REMOVED***
		return nil, x509.ErrUnsupportedAlgorithm
	***REMOVED***
	h := opts.hash().New()

	var publicKeyInfo struct ***REMOVED***
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	***REMOVED***
	if _, err := asn1.Unmarshal(issuer.RawSubjectPublicKeyInfo, &publicKeyInfo); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	h.Write(publicKeyInfo.PublicKey.RightAlign())
	issuerKeyHash := h.Sum(nil)

	h.Reset()
	h.Write(issuer.RawSubject)
	issuerNameHash := h.Sum(nil)

	req := &Request***REMOVED***
		HashAlgorithm:  hashFunc,
		IssuerNameHash: issuerNameHash,
		IssuerKeyHash:  issuerKeyHash,
		SerialNumber:   cert.SerialNumber,
	***REMOVED***
	return req.Marshal()
***REMOVED***

// CreateResponse returns a DER-encoded OCSP response with the specified contents.
// The fields in the response are populated as follows:
//
// The responder cert is used to populate the responder's name field, and the
// certificate itself is provided alongside the OCSP response signature.
//
// The issuer cert is used to populate the IssuerNameHash and IssuerKeyHash fields.
//
// The template is used to populate the SerialNumber, Status, RevokedAt,
// RevocationReason, ThisUpdate, and NextUpdate fields.
//
// If template.IssuerHash is not set, SHA1 will be used.
//
// The ProducedAt date is automatically set to the current date, to the nearest minute.
func CreateResponse(issuer, responderCert *x509.Certificate, template Response, priv crypto.Signer) ([]byte, error) ***REMOVED***
	var publicKeyInfo struct ***REMOVED***
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	***REMOVED***
	if _, err := asn1.Unmarshal(issuer.RawSubjectPublicKeyInfo, &publicKeyInfo); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if template.IssuerHash == 0 ***REMOVED***
		template.IssuerHash = crypto.SHA1
	***REMOVED***
	hashOID := getOIDFromHashAlgorithm(template.IssuerHash)
	if hashOID == nil ***REMOVED***
		return nil, errors.New("unsupported issuer hash algorithm")
	***REMOVED***

	if !template.IssuerHash.Available() ***REMOVED***
		return nil, fmt.Errorf("issuer hash algorithm %v not linked into binary", template.IssuerHash)
	***REMOVED***
	h := template.IssuerHash.New()
	h.Write(publicKeyInfo.PublicKey.RightAlign())
	issuerKeyHash := h.Sum(nil)

	h.Reset()
	h.Write(issuer.RawSubject)
	issuerNameHash := h.Sum(nil)

	innerResponse := singleResponse***REMOVED***
		CertID: certID***REMOVED***
			HashAlgorithm: pkix.AlgorithmIdentifier***REMOVED***
				Algorithm:  hashOID,
				Parameters: asn1.RawValue***REMOVED***Tag: 5 /* ASN.1 NULL */***REMOVED***,
			***REMOVED***,
			NameHash:      issuerNameHash,
			IssuerKeyHash: issuerKeyHash,
			SerialNumber:  template.SerialNumber,
		***REMOVED***,
		ThisUpdate:       template.ThisUpdate.UTC(),
		NextUpdate:       template.NextUpdate.UTC(),
		SingleExtensions: template.ExtraExtensions,
	***REMOVED***

	switch template.Status ***REMOVED***
	case Good:
		innerResponse.Good = true
	case Unknown:
		innerResponse.Unknown = true
	case Revoked:
		innerResponse.Revoked = revokedInfo***REMOVED***
			RevocationTime: template.RevokedAt.UTC(),
			Reason:         asn1.Enumerated(template.RevocationReason),
		***REMOVED***
	***REMOVED***

	rawResponderID := asn1.RawValue***REMOVED***
		Class:      2, // context-specific
		Tag:        1, // Name (explicit tag)
		IsCompound: true,
		Bytes:      responderCert.RawSubject,
	***REMOVED***
	tbsResponseData := responseData***REMOVED***
		Version:        0,
		RawResponderID: rawResponderID,
		ProducedAt:     time.Now().Truncate(time.Minute).UTC(),
		Responses:      []singleResponse***REMOVED***innerResponse***REMOVED***,
	***REMOVED***

	tbsResponseDataDER, err := asn1.Marshal(tbsResponseData)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	hashFunc, signatureAlgorithm, err := signingParamsForPublicKey(priv.Public(), template.SignatureAlgorithm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	responseHash := hashFunc.New()
	responseHash.Write(tbsResponseDataDER)
	signature, err := priv.Sign(rand.Reader, responseHash.Sum(nil), hashFunc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	response := basicResponse***REMOVED***
		TBSResponseData:    tbsResponseData,
		SignatureAlgorithm: signatureAlgorithm,
		Signature: asn1.BitString***REMOVED***
			Bytes:     signature,
			BitLength: 8 * len(signature),
		***REMOVED***,
	***REMOVED***
	if template.Certificate != nil ***REMOVED***
		response.Certificates = []asn1.RawValue***REMOVED***
			***REMOVED***FullBytes: template.Certificate.Raw***REMOVED***,
		***REMOVED***
	***REMOVED***
	responseDER, err := asn1.Marshal(response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return asn1.Marshal(responseASN1***REMOVED***
		Status: asn1.Enumerated(Success),
		Response: responseBytes***REMOVED***
			ResponseType: idPKIXOCSPBasic,
			Response:     responseDER,
		***REMOVED***,
	***REMOVED***)
***REMOVED***
