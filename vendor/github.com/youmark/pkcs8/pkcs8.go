// Package pkcs8 implements functions to parse and convert private keys in PKCS#8 format, as defined in RFC5208 and RFC5958
package pkcs8

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

// Copy from crypto/x509
var (
	oidPublicKeyRSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 1***REMOVED***
	oidPublicKeyDSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10040, 4, 1***REMOVED***
	oidPublicKeyECDSA = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 2, 1***REMOVED***
)

// Copy from crypto/x509
var (
	oidNamedCurveP224 = asn1.ObjectIdentifier***REMOVED***1, 3, 132, 0, 33***REMOVED***
	oidNamedCurveP256 = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 3, 1, 7***REMOVED***
	oidNamedCurveP384 = asn1.ObjectIdentifier***REMOVED***1, 3, 132, 0, 34***REMOVED***
	oidNamedCurveP521 = asn1.ObjectIdentifier***REMOVED***1, 3, 132, 0, 35***REMOVED***
)

// Copy from crypto/x509
func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) ***REMOVED***
	switch curve ***REMOVED***
	case elliptic.P224():
		return oidNamedCurveP224, true
	case elliptic.P256():
		return oidNamedCurveP256, true
	case elliptic.P384():
		return oidNamedCurveP384, true
	case elliptic.P521():
		return oidNamedCurveP521, true
	***REMOVED***

	return nil, false
***REMOVED***

// Unecrypted PKCS8
var (
	oidPKCS5PBKDF2    = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 5, 12***REMOVED***
	oidPBES2          = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 5, 13***REMOVED***
	oidAES256CBC      = asn1.ObjectIdentifier***REMOVED***2, 16, 840, 1, 101, 3, 4, 1, 42***REMOVED***
	oidAES128CBC      = asn1.ObjectIdentifier***REMOVED***2, 16, 840, 1, 101, 3, 4, 1, 2***REMOVED***
	oidHMACWithSHA256 = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 2, 9***REMOVED***
	oidDESEDE3CBC     = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 3, 7***REMOVED***
)

type ecPrivateKey struct ***REMOVED***
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
***REMOVED***

type privateKeyInfo struct ***REMOVED***
	Version             int
	PrivateKeyAlgorithm []asn1.ObjectIdentifier
	PrivateKey          []byte
***REMOVED***

// Encrypted PKCS8
type prfParam struct ***REMOVED***
	IdPRF     asn1.ObjectIdentifier
	NullParam asn1.RawValue
***REMOVED***

type pbkdf2Params struct ***REMOVED***
	Salt           []byte
	IterationCount int
	PrfParam       prfParam `asn1:"optional"`
***REMOVED***

type pbkdf2Algorithms struct ***REMOVED***
	IdPBKDF2     asn1.ObjectIdentifier
	PBKDF2Params pbkdf2Params
***REMOVED***

type pbkdf2Encs struct ***REMOVED***
	EncryAlgo asn1.ObjectIdentifier
	IV        []byte
***REMOVED***

type pbes2Params struct ***REMOVED***
	KeyDerivationFunc pbkdf2Algorithms
	EncryptionScheme  pbkdf2Encs
***REMOVED***

type pbes2Algorithms struct ***REMOVED***
	IdPBES2     asn1.ObjectIdentifier
	PBES2Params pbes2Params
***REMOVED***

type encryptedPrivateKeyInfo struct ***REMOVED***
	EncryptionAlgorithm pbes2Algorithms
	EncryptedData       []byte
***REMOVED***

// ParsePKCS8PrivateKeyRSA parses encrypted/unencrypted private keys in PKCS#8 format. To parse encrypted private keys, a password of []byte type should be provided to the function as the second parameter.
//
// The function can decrypt the private key encrypted with AES-256-CBC mode, and stored in PKCS #5 v2.0 format.
func ParsePKCS8PrivateKeyRSA(der []byte, v ...[]byte) (*rsa.PrivateKey, error) ***REMOVED***
	key, err := ParsePKCS8PrivateKey(der, v...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	typedKey, ok := key.(*rsa.PrivateKey)
	if !ok ***REMOVED***
		return nil, errors.New("key block is not of type RSA")
	***REMOVED***
	return typedKey, nil
***REMOVED***

// ParsePKCS8PrivateKeyECDSA parses encrypted/unencrypted private keys in PKCS#8 format. To parse encrypted private keys, a password of []byte type should be provided to the function as the second parameter.
//
// The function can decrypt the private key encrypted with AES-256-CBC mode, and stored in PKCS #5 v2.0 format.
func ParsePKCS8PrivateKeyECDSA(der []byte, v ...[]byte) (*ecdsa.PrivateKey, error) ***REMOVED***
	key, err := ParsePKCS8PrivateKey(der, v...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	typedKey, ok := key.(*ecdsa.PrivateKey)
	if !ok ***REMOVED***
		return nil, errors.New("key block is not of type ECDSA")
	***REMOVED***
	return typedKey, nil
***REMOVED***

// ParsePKCS8PrivateKey parses encrypted/unencrypted private keys in PKCS#8 format. To parse encrypted private keys, a password of []byte type should be provided to the function as the second parameter.
//
// The function can decrypt the private key encrypted with AES-256-CBC mode, and stored in PKCS #5 v2.0 format.
func ParsePKCS8PrivateKey(der []byte, v ...[]byte) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// No password provided, assume the private key is unencrypted
	if v == nil ***REMOVED***
		return x509.ParsePKCS8PrivateKey(der)
	***REMOVED***

	// Use the password provided to decrypt the private key
	password := v[0]
	var privKey encryptedPrivateKeyInfo
	if _, err := asn1.Unmarshal(der, &privKey); err != nil ***REMOVED***
		return nil, errors.New("pkcs8: only PKCS #5 v2.0 supported")
	***REMOVED***

	if !privKey.EncryptionAlgorithm.IdPBES2.Equal(oidPBES2) ***REMOVED***
		return nil, errors.New("pkcs8: only PBES2 supported")
	***REMOVED***

	if !privKey.EncryptionAlgorithm.PBES2Params.KeyDerivationFunc.IdPBKDF2.Equal(oidPKCS5PBKDF2) ***REMOVED***
		return nil, errors.New("pkcs8: only PBKDF2 supported")
	***REMOVED***

	encParam := privKey.EncryptionAlgorithm.PBES2Params.EncryptionScheme
	kdfParam := privKey.EncryptionAlgorithm.PBES2Params.KeyDerivationFunc.PBKDF2Params

	iv := encParam.IV
	salt := kdfParam.Salt
	iter := kdfParam.IterationCount
	keyHash := sha1.New
	if kdfParam.PrfParam.IdPRF.Equal(oidHMACWithSHA256) ***REMOVED***
		keyHash = sha256.New
	***REMOVED***

	encryptedKey := privKey.EncryptedData
	var symkey []byte
	var block cipher.Block
	var err error
	switch ***REMOVED***
	case encParam.EncryAlgo.Equal(oidAES128CBC):
		symkey = pbkdf2.Key(password, salt, iter, 16, keyHash)
		block, err = aes.NewCipher(symkey)
	case encParam.EncryAlgo.Equal(oidAES256CBC):
		symkey = pbkdf2.Key(password, salt, iter, 32, keyHash)
		block, err = aes.NewCipher(symkey)
	case encParam.EncryAlgo.Equal(oidDESEDE3CBC):
		symkey = pbkdf2.Key(password, salt, iter, 24, keyHash)
		block, err = des.NewTripleDESCipher(symkey)
	default:
		return nil, errors.New("pkcs8: only AES-256-CBC, AES-128-CBC and DES-EDE3-CBC are supported")
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptedKey, encryptedKey)

	key, err := x509.ParsePKCS8PrivateKey(encryptedKey)
	if err != nil ***REMOVED***
		return nil, errors.New("pkcs8: incorrect password")
	***REMOVED***
	return key, nil
***REMOVED***

func convertPrivateKeyToPKCS8(priv interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	var pkey privateKeyInfo

	switch priv := priv.(type) ***REMOVED***
	case *ecdsa.PrivateKey:
		eckey, err := x509.MarshalECPrivateKey(priv)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		oidNamedCurve, ok := oidFromNamedCurve(priv.Curve)
		if !ok ***REMOVED***
			return nil, errors.New("pkcs8: unknown elliptic curve")
		***REMOVED***

		// Per RFC5958, if publicKey is present, then version is set to v2(1) else version is set to v1(0).
		// But openssl set to v1 even publicKey is present
		pkey.Version = 1
		pkey.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 2)
		pkey.PrivateKeyAlgorithm[0] = oidPublicKeyECDSA
		pkey.PrivateKeyAlgorithm[1] = oidNamedCurve
		pkey.PrivateKey = eckey
	case *rsa.PrivateKey:

		// Per RFC5958, if publicKey is present, then version is set to v2(1) else version is set to v1(0).
		// But openssl set to v1 even publicKey is present
		pkey.Version = 0
		pkey.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
		pkey.PrivateKeyAlgorithm[0] = oidPublicKeyRSA
		pkey.PrivateKey = x509.MarshalPKCS1PrivateKey(priv)
	***REMOVED***

	return asn1.Marshal(pkey)
***REMOVED***

func convertPrivateKeyToPKCS8Encrypted(priv interface***REMOVED******REMOVED***, password []byte) ([]byte, error) ***REMOVED***
	// Convert private key into PKCS8 format
	pkey, err := convertPrivateKeyToPKCS8(priv)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Calculate key from password based on PKCS5 algorithm
	// Use 8 byte salt, 16 byte IV, and 2048 iteration
	iter := 2048
	salt := make([]byte, 8)
	iv := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	_, err = rand.Read(iv)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	key := pbkdf2.Key(password, salt, iter, 32, sha256.New)

	// Use AES256-CBC mode, pad plaintext with PKCS5 padding scheme
	padding := aes.BlockSize - len(pkey)%aes.BlockSize
	if padding > 0 ***REMOVED***
		n := len(pkey)
		pkey = append(pkey, make([]byte, padding)...)
		for i := 0; i < padding; i++ ***REMOVED***
			pkey[n+i] = byte(padding)
		***REMOVED***
	***REMOVED***

	encryptedKey := make([]byte, len(pkey))
	block, err := aes.NewCipher(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptedKey, pkey)

	//	pbkdf2algo := pbkdf2Algorithms***REMOVED***oidPKCS5PBKDF2, pbkdf2Params***REMOVED***salt, iter, prfParam***REMOVED***oidHMACWithSHA256***REMOVED******REMOVED******REMOVED***

	pbkdf2algo := pbkdf2Algorithms***REMOVED***oidPKCS5PBKDF2, pbkdf2Params***REMOVED***salt, iter, prfParam***REMOVED***oidHMACWithSHA256, asn1.RawValue***REMOVED***Tag: asn1.TagNull***REMOVED******REMOVED******REMOVED******REMOVED***
	pbkdf2encs := pbkdf2Encs***REMOVED***oidAES256CBC, iv***REMOVED***
	pbes2algo := pbes2Algorithms***REMOVED***oidPBES2, pbes2Params***REMOVED***pbkdf2algo, pbkdf2encs***REMOVED******REMOVED***

	encryptedPkey := encryptedPrivateKeyInfo***REMOVED***pbes2algo, encryptedKey***REMOVED***

	return asn1.Marshal(encryptedPkey)
***REMOVED***

// ConvertPrivateKeyToPKCS8 converts the private key into PKCS#8 format.
// To encrypt the private key, the password of []byte type should be provided as the second parameter.
//
// The only supported key types are RSA and ECDSA (*rsa.PublicKey or *ecdsa.PublicKey for priv)
func ConvertPrivateKeyToPKCS8(priv interface***REMOVED******REMOVED***, v ...[]byte) ([]byte, error) ***REMOVED***
	if v == nil ***REMOVED***
		return convertPrivateKeyToPKCS8(priv)
	***REMOVED***

	password := string(v[0])
	return convertPrivateKeyToPKCS8Encrypted(priv, []byte(password))
***REMOVED***
