//Copyright 2013 Thomson Reuters Global Resources. BSD License please see License file for more information

package ntlm

// Define KXKEY(SessionBaseKey, LmChallengeResponse, ServerChallenge) as
func kxKey(flags uint32, sessionBaseKey []byte, lmChallengeResponse []byte, serverChallenge []byte, lmnowf []byte) (keyExchangeKey []byte, err error) ***REMOVED***
	if NTLMSSP_NEGOTIATE_LM_KEY.IsSet(flags) ***REMOVED***
		var part1, part2 []byte
		part1, err = des(lmnowf[0:7], lmChallengeResponse[0:8])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		key := append([]byte***REMOVED***lmnowf[7]***REMOVED***, []byte***REMOVED***0xBD, 0xBD, 0xBD, 0xBD, 0xBD, 0xBD***REMOVED***...)
		part2, err = des(key, lmChallengeResponse[0:8])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		keyExchangeKey = concat(part1, part2)
	***REMOVED*** else if NTLMSSP_REQUEST_NON_NT_SESSION_KEY.IsSet(flags) ***REMOVED***
		keyExchangeKey = concat(lmnowf[0:8], zeroBytes(8))
	***REMOVED*** else ***REMOVED***
		keyExchangeKey = sessionBaseKey
	***REMOVED***

	return
***REMOVED***

// Define SIGNKEY(NegFlg, RandomSessionKey, Mode) as
func signKey(flags uint32, randomSessionKey []byte, mode string) (signKey []byte) ***REMOVED***
	if NTLMSSP_NEGOTIATE_EXTENDED_SESSIONSECURITY.IsSet(flags) ***REMOVED***
		if mode == "Client" ***REMOVED***
			signKey = md5(concat(randomSessionKey, []byte("session key to client-to-server signing key magic constant\x00")))
		***REMOVED*** else ***REMOVED***
			signKey = md5(concat(randomSessionKey, []byte("session key to server-to-client signing key magic constant\x00")))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		signKey = nil
	***REMOVED***
	return
***REMOVED***

// 	Define SEALKEY(NegotiateFlags, RandomSessionKey, Mode) as
func sealKey(flags uint32, randomSessionKey []byte, mode string) (sealKey []byte) ***REMOVED***
	if NTLMSSP_NEGOTIATE_EXTENDED_SESSIONSECURITY.IsSet(flags) ***REMOVED***
		if NTLMSSP_NEGOTIATE_128.IsSet(flags) ***REMOVED***
			sealKey = randomSessionKey
		***REMOVED*** else if NTLMSSP_NEGOTIATE_56.IsSet(flags) ***REMOVED***
			sealKey = randomSessionKey[0:7]
		***REMOVED*** else ***REMOVED***
			sealKey = randomSessionKey[0:5]
		***REMOVED***
		if mode == "Client" ***REMOVED***
			sealKey = md5(concat(sealKey, []byte("session key to client-to-server sealing key magic constant\x00")))
		***REMOVED*** else ***REMOVED***
			sealKey = md5(concat(sealKey, []byte("session key to server-to-client sealing key magic constant\x00")))
		***REMOVED***
	***REMOVED*** else if NTLMSSP_NEGOTIATE_LM_KEY.IsSet(flags) ***REMOVED***
		if NTLMSSP_NEGOTIATE_56.IsSet(flags) ***REMOVED***
			sealKey = concat(randomSessionKey[0:7], []byte***REMOVED***0xA0***REMOVED***)
		***REMOVED*** else ***REMOVED***
			sealKey = concat(randomSessionKey[0:5], []byte***REMOVED***0xE5, 0x38, 0xB0***REMOVED***)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		sealKey = randomSessionKey
	***REMOVED***

	return
***REMOVED***
