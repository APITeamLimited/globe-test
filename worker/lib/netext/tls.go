package netext

import (
	"crypto/tls"

	"golang.org/x/crypto/ocsp"

	"github.com/APITeamLimited/k6-worker/lib"
)

//nolint:golint
const (
	OCSP_STATUS_GOOD                   = "good"
	OCSP_STATUS_REVOKED                = "revoked"
	OCSP_STATUS_SERVER_FAILED          = "server_failed"
	OCSP_STATUS_UNKNOWN                = "unknown"
	OCSP_REASON_UNSPECIFIED            = "unspecified"
	OCSP_REASON_KEY_COMPROMISE         = "key_compromise"
	OCSP_REASON_CA_COMPROMISE          = "ca_compromise"
	OCSP_REASON_AFFILIATION_CHANGED    = "affiliation_changed"
	OCSP_REASON_SUPERSEDED             = "superseded"
	OCSP_REASON_CESSATION_OF_OPERATION = "cessation_of_operation"
	OCSP_REASON_CERTIFICATE_HOLD       = "certificate_hold"
	OCSP_REASON_REMOVE_FROM_CRL        = "remove_from_crl"
	OCSP_REASON_PRIVILEGE_WITHDRAWN    = "privilege_withdrawn"
	OCSP_REASON_AA_COMPROMISE          = "aa_compromise"
	TLS_1_0                            = "tls1.0"
	TLS_1_1                            = "tls1.1"
	TLS_1_2                            = "tls1.2"
	TLS_1_3                            = "tls1.3"
)

type TLSInfo struct ***REMOVED***
	Version     string
	CipherSuite string
***REMOVED***
type OCSP struct ***REMOVED***
	ProducedAt       int64  `json:"produced_at"`
	ThisUpdate       int64  `json:"this_update"`
	NextUpdate       int64  `json:"next_update"`
	RevokedAt        int64  `json:"revoked_at"`
	RevocationReason string `json:"revocation_reason"`
	Status           string `json:"status"`
***REMOVED***

func ParseTLSConnState(tlsState *tls.ConnectionState) (TLSInfo, OCSP) ***REMOVED***
	tlsInfo := TLSInfo***REMOVED******REMOVED***
	switch tlsState.Version ***REMOVED***
	case tls.VersionTLS10:
		tlsInfo.Version = TLS_1_0
	case tls.VersionTLS11:
		tlsInfo.Version = TLS_1_1
	case tls.VersionTLS12:
		tlsInfo.Version = TLS_1_2
	case tls.VersionTLS13:
		tlsInfo.Version = TLS_1_3
	***REMOVED***

	tlsInfo.CipherSuite = lib.SupportedTLSCipherSuitesToString[tlsState.CipherSuite]
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

	return tlsInfo, ocspStapledRes
***REMOVED***
