package x509

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha1" // #nosec G505
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/APITeamLimited/k6-worker/js/common"
	"github.com/APITeamLimited/k6-worker/js/modules"
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct***REMOVED******REMOVED***

	// X509 represents an instance of the X509 certificate module.
	X509 struct ***REMOVED***
		vu modules.VU
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &X509***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &X509***REMOVED***vu: vu***REMOVED***
***REMOVED***

// Exports returns the exports of the execution module.
func (mi *X509) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"parse":       mi.parse,
			"getAltNames": mi.altNames,
			"getIssuer":   mi.issuer,
			"getSubject":  mi.subject,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Certificate is an X.509 certificate
type Certificate struct ***REMOVED***
	Subject            Subject
	Issuer             Issuer
	NotBefore          string    `js:"notBefore"`
	NotAfter           string    `js:"notAfter"`
	AltNames           []string  `js:"altNames"`
	SignatureAlgorithm string    `js:"signatureAlgorithm"`
	FingerPrint        []byte    `js:"fingerPrint"`
	PublicKey          PublicKey `js:"publicKey"`
***REMOVED***

// RDN is a component of an X.509 distinguished name
type RDN struct ***REMOVED***
	Type  string
	Value string
***REMOVED***

// Subject is a certificate subject
type Subject struct ***REMOVED***
	CommonName             string `js:"commonName"`
	Country                string
	PostalCode             string   `js:"postalCode"`
	StateOrProvinceName    string   `js:"stateOrProvinceName"`
	LocalityName           string   `js:"localityName"`
	StreetAddress          string   `js:"streetAddress"`
	OrganizationName       string   `js:"organizationName"`
	OrganizationalUnitName []string `js:"organizationalUnitName"`
	Names                  []RDN
***REMOVED***

// Issuer is a certificate issuer
type Issuer struct ***REMOVED***
	CommonName          string `js:"commonName"`
	Country             string
	StateOrProvinceName string `js:"stateOrProvinceName"`
	LocalityName        string `js:"localityName"`
	OrganizationName    string `js:"organizationName"`
	Names               []RDN
***REMOVED***

// PublicKey is used for decryption and signature verification
type PublicKey struct ***REMOVED***
	Algorithm string
	Key       interface***REMOVED******REMOVED***
***REMOVED***

// parse produces an entire X.509 certificate
func (mi X509) parse(encoded []byte) (Certificate, error) ***REMOVED***
	parsed, err := parseCertificate(encoded)
	if err != nil ***REMOVED***
		return Certificate***REMOVED******REMOVED***, err
	***REMOVED***
	certificate, err := makeCertificate(parsed)
	if err != nil ***REMOVED***
		return Certificate***REMOVED******REMOVED***, err
	***REMOVED***
	return certificate, nil
***REMOVED***

// altNames extracts alt names
func (mi X509) altNames(encoded []byte) ([]string, error) ***REMOVED***
	parsed, err := parseCertificate(encoded)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return altNames(parsed), nil
***REMOVED***

// issuer extracts certificate issuer
func (mi X509) issuer(encoded []byte) (Issuer, error) ***REMOVED***
	parsed, err := parseCertificate(encoded)
	if err != nil ***REMOVED***
		return Issuer***REMOVED******REMOVED***, err
	***REMOVED***
	return makeIssuer(parsed.Issuer), nil
***REMOVED***

// subject extracts certificate subject
func (mi X509) subject(encoded []byte) Subject ***REMOVED***
	parsed, err := parseCertificate(encoded)
	if err != nil ***REMOVED***
		common.Throw(mi.vu.Runtime(), err)
	***REMOVED***
	return makeSubject(parsed.Subject)
***REMOVED***

func parseCertificate(encoded []byte) (*x509.Certificate, error) ***REMOVED***
	decoded, _ := pem.Decode(encoded)
	if decoded == nil ***REMOVED***
		return nil, fmt.Errorf("failed to decode certificate PEM file")
	***REMOVED***
	parsed, err := x509.ParseCertificate(decoded.Bytes)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	***REMOVED***
	return parsed, nil
***REMOVED***

func makeCertificate(parsed *x509.Certificate) (Certificate, error) ***REMOVED***
	publicKey, err := makePublicKey(parsed.PublicKey)
	if err != nil ***REMOVED***
		return Certificate***REMOVED******REMOVED***, err
	***REMOVED***
	return Certificate***REMOVED***
		Subject:            makeSubject(parsed.Subject),
		Issuer:             makeIssuer(parsed.Issuer),
		NotBefore:          iso8601(parsed.NotBefore),
		NotAfter:           iso8601(parsed.NotAfter),
		AltNames:           altNames(parsed),
		SignatureAlgorithm: signatureAlgorithm(parsed.SignatureAlgorithm),
		FingerPrint:        fingerPrint(parsed),
		PublicKey:          publicKey,
	***REMOVED***, nil
***REMOVED***

func makeSubject(subject pkix.Name) Subject ***REMOVED***
	return Subject***REMOVED***
		CommonName:             subject.CommonName,
		Country:                first(subject.Country),
		PostalCode:             first(subject.PostalCode),
		StateOrProvinceName:    first(subject.Province),
		LocalityName:           first(subject.Locality),
		StreetAddress:          first(subject.StreetAddress),
		OrganizationName:       first(subject.Organization),
		OrganizationalUnitName: subject.OrganizationalUnit,
		Names:                  makeRdns(subject.Names),
	***REMOVED***
***REMOVED***

func makeIssuer(issuer pkix.Name) Issuer ***REMOVED***
	return Issuer***REMOVED***
		CommonName:          issuer.CommonName,
		Country:             first(issuer.Country),
		StateOrProvinceName: first(issuer.Province),
		LocalityName:        first(issuer.Locality),
		OrganizationName:    first(issuer.Organization),
		Names:               makeRdns(issuer.Names),
	***REMOVED***
***REMOVED***

func makePublicKey(parsed interface***REMOVED******REMOVED***) (PublicKey, error) ***REMOVED***
	var algorithm string
	switch parsed.(type) ***REMOVED***
	case *dsa.PublicKey:
		algorithm = "DSA"
	case *ecdsa.PublicKey:
		algorithm = "ECDSA"
	case *rsa.PublicKey:
		algorithm = "RSA"
	default:
		err := errors.New("unsupported public key algorithm")
		return PublicKey***REMOVED******REMOVED***, err
	***REMOVED***
	return PublicKey***REMOVED***
		Algorithm: algorithm,
		Key:       parsed,
	***REMOVED***, nil
***REMOVED***

func first(values []string) string ***REMOVED***
	if len(values) > 0 ***REMOVED***
		return values[0]
	***REMOVED***
	return ""
***REMOVED***

func iso8601(value time.Time) string ***REMOVED***
	return value.Format(time.RFC3339)
***REMOVED***

func makeRdns(names []pkix.AttributeTypeAndValue) []RDN ***REMOVED***
	result := make([]RDN, len(names))
	for i, name := range names ***REMOVED***
		result[i] = makeRdn(name)
	***REMOVED***
	return result
***REMOVED***

func makeRdn(name pkix.AttributeTypeAndValue) RDN ***REMOVED***
	return RDN***REMOVED***
		Type:  name.Type.String(),
		Value: fmt.Sprintf("%v", name.Value),
	***REMOVED***
***REMOVED***

func altNames(parsed *x509.Certificate) []string ***REMOVED***
	var names []string
	names = append(names, parsed.DNSNames...)
	names = append(names, parsed.EmailAddresses...)
	names = append(names, ipAddresses(parsed)...)
	names = append(names, uris(parsed)...)
	return names
***REMOVED***

func ipAddresses(parsed *x509.Certificate) []string ***REMOVED***
	strings := make([]string, len(parsed.IPAddresses))
	for i, item := range parsed.IPAddresses ***REMOVED***
		strings[i] = item.String()
	***REMOVED***
	return strings
***REMOVED***

func uris(parsed *x509.Certificate) []string ***REMOVED***
	strings := make([]string, len(parsed.URIs))
	for i, item := range parsed.URIs ***REMOVED***
		strings[i] = item.String()
	***REMOVED***
	return strings
***REMOVED***

func signatureAlgorithm(value x509.SignatureAlgorithm) string ***REMOVED***
	if value == x509.UnknownSignatureAlgorithm ***REMOVED***
		return "UnknownSignatureAlgorithm"
	***REMOVED***
	return value.String()
***REMOVED***

func fingerPrint(parsed *x509.Certificate) []byte ***REMOVED***
	bytes := sha1.Sum(parsed.Raw) // #nosec G401
	return bytes[:]
***REMOVED***
