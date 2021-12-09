package ntlmssp

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

type authenicateMessage struct ***REMOVED***
	LmChallengeResponse []byte
	NtChallengeResponse []byte

	TargetName string
	UserName   string

	// only set if negotiateFlag_NTLMSSP_NEGOTIATE_KEY_EXCH
	EncryptedRandomSessionKey []byte

	NegotiateFlags negotiateFlags

	MIC []byte
***REMOVED***

type authenticateMessageFields struct ***REMOVED***
	messageHeader
	LmChallengeResponse varField
	NtChallengeResponse varField
	TargetName          varField
	UserName            varField
	Workstation         varField
	_                   [8]byte
	NegotiateFlags      negotiateFlags
***REMOVED***

func (m authenicateMessage) MarshalBinary() ([]byte, error) ***REMOVED***
	if !m.NegotiateFlags.Has(negotiateFlagNTLMSSPNEGOTIATEUNICODE) ***REMOVED***
		return nil, errors.New("Only unicode is supported")
	***REMOVED***

	target, user := toUnicode(m.TargetName), toUnicode(m.UserName)
	workstation := toUnicode("")

	ptr := binary.Size(&authenticateMessageFields***REMOVED******REMOVED***)
	f := authenticateMessageFields***REMOVED***
		messageHeader:       newMessageHeader(3),
		NegotiateFlags:      m.NegotiateFlags,
		LmChallengeResponse: newVarField(&ptr, len(m.LmChallengeResponse)),
		NtChallengeResponse: newVarField(&ptr, len(m.NtChallengeResponse)),
		TargetName:          newVarField(&ptr, len(target)),
		UserName:            newVarField(&ptr, len(user)),
		Workstation:         newVarField(&ptr, len(workstation)),
	***REMOVED***

	f.NegotiateFlags.Unset(negotiateFlagNTLMSSPNEGOTIATEVERSION)

	b := bytes.Buffer***REMOVED******REMOVED***
	if err := binary.Write(&b, binary.LittleEndian, &f); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&b, binary.LittleEndian, &m.LmChallengeResponse); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&b, binary.LittleEndian, &m.NtChallengeResponse); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&b, binary.LittleEndian, &target); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&b, binary.LittleEndian, &user); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&b, binary.LittleEndian, &workstation); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return b.Bytes(), nil
***REMOVED***

//ProcessChallenge crafts an AUTHENTICATE message in response to the CHALLENGE message
//that was received from the server
func ProcessChallenge(challengeMessageData []byte, user, password string) ([]byte, error) ***REMOVED***
	if user == "" && password == "" ***REMOVED***
		return nil, errors.New("Anonymous authentication not supported")
	***REMOVED***

	var cm challengeMessage
	if err := cm.UnmarshalBinary(challengeMessageData); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if cm.NegotiateFlags.Has(negotiateFlagNTLMSSPNEGOTIATELMKEY) ***REMOVED***
		return nil, errors.New("Only NTLM v2 is supported, but server requested v1 (NTLMSSP_NEGOTIATE_LM_KEY)")
	***REMOVED***
	if cm.NegotiateFlags.Has(negotiateFlagNTLMSSPNEGOTIATEKEYEXCH) ***REMOVED***
		return nil, errors.New("Key exchange requested but not supported (NTLMSSP_NEGOTIATE_KEY_EXCH)")
	***REMOVED***

	am := authenicateMessage***REMOVED***
		UserName:       user,
		TargetName:     cm.TargetName,
		NegotiateFlags: cm.NegotiateFlags,
	***REMOVED***

	timestamp := cm.TargetInfo[avIDMsvAvTimestamp]
	if timestamp == nil ***REMOVED*** // no time sent, take current time
		ft := uint64(time.Now().UnixNano()) / 100
		ft += 116444736000000000 // add time between unix & windows offset
		timestamp = make([]byte, 8)
		binary.LittleEndian.PutUint64(timestamp, ft)
	***REMOVED***

	clientChallenge := make([]byte, 8)
	rand.Reader.Read(clientChallenge)

	ntlmV2Hash := getNtlmV2Hash(password, user, cm.TargetName)

	am.NtChallengeResponse = computeNtlmV2Response(ntlmV2Hash,
		cm.ServerChallenge[:], clientChallenge, timestamp, cm.TargetInfoRaw)

	if cm.TargetInfoRaw == nil ***REMOVED***
		am.LmChallengeResponse = computeLmV2Response(ntlmV2Hash,
			cm.ServerChallenge[:], clientChallenge)
	***REMOVED***
	return am.MarshalBinary()
***REMOVED***

func ProcessChallengeWithHash(challengeMessageData []byte, user, hash string) ([]byte, error) ***REMOVED***
	if user == "" && hash == "" ***REMOVED***
		return nil, errors.New("Anonymous authentication not supported")
	***REMOVED***

	var cm challengeMessage
	if err := cm.UnmarshalBinary(challengeMessageData); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if cm.NegotiateFlags.Has(negotiateFlagNTLMSSPNEGOTIATELMKEY) ***REMOVED***
		return nil, errors.New("Only NTLM v2 is supported, but server requested v1 (NTLMSSP_NEGOTIATE_LM_KEY)")
	***REMOVED***
	if cm.NegotiateFlags.Has(negotiateFlagNTLMSSPNEGOTIATEKEYEXCH) ***REMOVED***
		return nil, errors.New("Key exchange requested but not supported (NTLMSSP_NEGOTIATE_KEY_EXCH)")
	***REMOVED***

	am := authenicateMessage***REMOVED***
		UserName:       user,
		TargetName:     cm.TargetName,
		NegotiateFlags: cm.NegotiateFlags,
	***REMOVED***

	timestamp := cm.TargetInfo[avIDMsvAvTimestamp]
	if timestamp == nil ***REMOVED*** // no time sent, take current time
		ft := uint64(time.Now().UnixNano()) / 100
		ft += 116444736000000000 // add time between unix & windows offset
		timestamp = make([]byte, 8)
		binary.LittleEndian.PutUint64(timestamp, ft)
	***REMOVED***

	clientChallenge := make([]byte, 8)
	rand.Reader.Read(clientChallenge)

	hashParts := strings.Split(hash, ":")
	if len(hashParts) > 1 ***REMOVED***
		hash = hashParts[1]
	***REMOVED***
	hashBytes, err := hex.DecodeString(hash)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ntlmV2Hash := hmacMd5(hashBytes, toUnicode(strings.ToUpper(user)+cm.TargetName))

	am.NtChallengeResponse = computeNtlmV2Response(ntlmV2Hash,
		cm.ServerChallenge[:], clientChallenge, timestamp, cm.TargetInfoRaw)

	if cm.TargetInfoRaw == nil ***REMOVED***
		am.LmChallengeResponse = computeLmV2Response(ntlmV2Hash,
			cm.ServerChallenge[:], clientChallenge)
	***REMOVED***
	return am.MarshalBinary()
***REMOVED***
