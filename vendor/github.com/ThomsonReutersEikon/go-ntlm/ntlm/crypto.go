//Copyright 2013 Thomson Reuters Global Resources. BSD License please see License file for more information

package ntlm

import (
	desP "crypto/des"
	hmacP "crypto/hmac"
	md5P "crypto/md5"
	"crypto/rand"
	rc4P "crypto/rc4"
	crc32P "hash/crc32"

	md4P "github.com/ThomsonReutersEikon/go-ntlm/ntlm/md4"
)

func md4(data []byte) []byte ***REMOVED***
	md4 := md4P.New()
	md4.Write(data)
	return md4.Sum(nil)
***REMOVED***

func md5(data []byte) []byte ***REMOVED***
	md5 := md5P.New()
	md5.Write(data)
	return md5.Sum(nil)
***REMOVED***

// Indicates the computation of a 16-byte HMAC-keyed MD5 message digest of the byte string M using the key K.
func hmacMd5(key []byte, data []byte) []byte ***REMOVED***
	mac := hmacP.New(md5P.New, key)
	mac.Write(data)
	return mac.Sum(nil)
***REMOVED***

// Indicates the computation of an N-byte cryptographic- strength random number.
func nonce(length int) []byte ***REMOVED***
	result := make([]byte, length)
	rand.Read(result)
	return result
***REMOVED***

func crc32(bytes []byte) uint32 ***REMOVED***
	crc := crc32P.New(crc32P.IEEETable)
	crc.Write(bytes)
	return crc.Sum32()
***REMOVED***

// Indicates the encryption of data item D with the key K using the RC4 algorithm.
func rc4K(key []byte, ciphertext []byte) ([]byte, error) ***REMOVED***
	cipher, err := rc4P.NewCipher(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	result := make([]byte, len(ciphertext))
	cipher.XORKeyStream(result, ciphertext)
	return result, nil
***REMOVED***

func rc4Init(key []byte) (cipher *rc4P.Cipher, err error) ***REMOVED***
	cipher, err = rc4P.NewCipher(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return cipher, nil
***REMOVED***

func rc4(cipher *rc4P.Cipher, ciphertext []byte) []byte ***REMOVED***
	result := make([]byte, len(ciphertext))
	cipher.XORKeyStream(result, ciphertext)
	return result
***REMOVED***

// Indicates the encryption of an 8-byte data item D with the 7-byte key K using the Data Encryption Standard (DES)
// algorithm in Electronic Codebook (ECB) mode. The result is 8 bytes in length ([FIPS46-2]).
func des(key []byte, ciphertext []byte) ([]byte, error) ***REMOVED***
	calcKey := createDesKey(key)
	cipher, err := desP.NewCipher(calcKey)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	result := make([]byte, len(ciphertext))
	cipher.Encrypt(result, ciphertext)

	return result, nil
***REMOVED***

// Indicates the encryption of an 8-byte data item D with the 16-byte key K using the Data Encryption Standard Long (DESL) algorithm.
// The result is 24 bytes in length. DESL(K, D) is computed as follows.
// Note K[] implies a key represented as a character array.
func desL(key []byte, cipherText []byte) ([]byte, error) ***REMOVED***
	out1, err := des(zeroPaddedBytes(key, 0, 7), cipherText)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	out2, err := des(zeroPaddedBytes(key, 7, 7), cipherText)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	out3, err := des(zeroPaddedBytes(key, 14, 7), cipherText)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return concat(out1, out2, out3), nil
***REMOVED***

// Creates a DES encryption key from the given 7 byte key material.
func createDesKey(keyBytes []byte) []byte ***REMOVED***
	material := zeroBytes(8)
	material[0] = keyBytes[0]
	material[1] = (byte)(keyBytes[0]<<7 | (keyBytes[1]&0xff)>>1)
	material[2] = (byte)(keyBytes[1]<<6 | (keyBytes[2]&0xff)>>2)
	material[3] = (byte)(keyBytes[2]<<5 | (keyBytes[3]&0xff)>>3)
	material[4] = (byte)(keyBytes[3]<<4 | (keyBytes[4]&0xff)>>4)
	material[5] = (byte)(keyBytes[4]<<3 | (keyBytes[5]&0xff)>>5)
	material[6] = (byte)(keyBytes[5]<<2 | (keyBytes[6]&0xff)>>6)
	material[7] = (byte)(keyBytes[6] << 1)
	oddParity(material)
	return material
***REMOVED***

// Applies odd parity to the given byte array.
func oddParity(bytes []byte) ***REMOVED***
	for i := 0; i < len(bytes); i++ ***REMOVED***
		b := bytes[i]
		needsParity := (((b >> 7) ^ (b >> 6) ^ (b >> 5) ^ (b >> 4) ^ (b >> 3) ^ (b >> 2) ^ (b >> 1)) & 0x01) == 0
		if needsParity ***REMOVED***
			bytes[i] = bytes[i] | byte(0x01)
		***REMOVED*** else ***REMOVED***
			bytes[i] = bytes[i] & byte(0xfe)
		***REMOVED***
	***REMOVED***
***REMOVED***
