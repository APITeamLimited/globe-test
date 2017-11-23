// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package windows

import (
	"syscall"
	"unsafe"
)

const (
	STANDARD_RIGHTS_REQUIRED = 0xf0000
	STANDARD_RIGHTS_READ     = 0x20000
	STANDARD_RIGHTS_WRITE    = 0x20000
	STANDARD_RIGHTS_EXECUTE  = 0x20000
	STANDARD_RIGHTS_ALL      = 0x1F0000
)

const (
	NameUnknown          = 0
	NameFullyQualifiedDN = 1
	NameSamCompatible    = 2
	NameDisplay          = 3
	NameUniqueId         = 6
	NameCanonical        = 7
	NameUserPrincipal    = 8
	NameCanonicalEx      = 9
	NameServicePrincipal = 10
	NameDnsDomain        = 12
)

// This function returns 1 byte BOOLEAN rather than the 4 byte BOOL.
// http://blogs.msdn.com/b/drnick/archive/2007/12/19/windows-and-upn-format-credentials.aspx
//sys	TranslateName(accName *uint16, accNameFormat uint32, desiredNameFormat uint32, translatedName *uint16, nSize *uint32) (err error) [failretval&0xff==0] = secur32.TranslateNameW
//sys	GetUserNameEx(nameFormat uint32, nameBuffre *uint16, nSize *uint32) (err error) [failretval&0xff==0] = secur32.GetUserNameExW

// TranslateAccountName converts a directory service
// object name from one format to another.
func TranslateAccountName(username string, from, to uint32, initSize int) (string, error) ***REMOVED***
	u, e := UTF16PtrFromString(username)
	if e != nil ***REMOVED***
		return "", e
	***REMOVED***
	n := uint32(50)
	for ***REMOVED***
		b := make([]uint16, n)
		e = TranslateName(u, from, to, &b[0], &n)
		if e == nil ***REMOVED***
			return UTF16ToString(b[:n]), nil
		***REMOVED***
		if e != ERROR_INSUFFICIENT_BUFFER ***REMOVED***
			return "", e
		***REMOVED***
		if n <= uint32(len(b)) ***REMOVED***
			return "", e
		***REMOVED***
	***REMOVED***
***REMOVED***

const (
	// do not reorder
	NetSetupUnknownStatus = iota
	NetSetupUnjoined
	NetSetupWorkgroupName
	NetSetupDomainName
)

type UserInfo10 struct ***REMOVED***
	Name       *uint16
	Comment    *uint16
	UsrComment *uint16
	FullName   *uint16
***REMOVED***

//sys	NetUserGetInfo(serverName *uint16, userName *uint16, level uint32, buf **byte) (neterr error) = netapi32.NetUserGetInfo
//sys	NetGetJoinInformation(server *uint16, name **uint16, bufType *uint32) (neterr error) = netapi32.NetGetJoinInformation
//sys	NetApiBufferFree(buf *byte) (neterr error) = netapi32.NetApiBufferFree

const (
	// do not reorder
	SidTypeUser = 1 + iota
	SidTypeGroup
	SidTypeDomain
	SidTypeAlias
	SidTypeWellKnownGroup
	SidTypeDeletedAccount
	SidTypeInvalid
	SidTypeUnknown
	SidTypeComputer
	SidTypeLabel
)

type SidIdentifierAuthority struct ***REMOVED***
	Value [6]byte
***REMOVED***

var (
	SECURITY_NULL_SID_AUTHORITY        = SidIdentifierAuthority***REMOVED***[6]byte***REMOVED***0, 0, 0, 0, 0, 0***REMOVED******REMOVED***
	SECURITY_WORLD_SID_AUTHORITY       = SidIdentifierAuthority***REMOVED***[6]byte***REMOVED***0, 0, 0, 0, 0, 1***REMOVED******REMOVED***
	SECURITY_LOCAL_SID_AUTHORITY       = SidIdentifierAuthority***REMOVED***[6]byte***REMOVED***0, 0, 0, 0, 0, 2***REMOVED******REMOVED***
	SECURITY_CREATOR_SID_AUTHORITY     = SidIdentifierAuthority***REMOVED***[6]byte***REMOVED***0, 0, 0, 0, 0, 3***REMOVED******REMOVED***
	SECURITY_NON_UNIQUE_AUTHORITY      = SidIdentifierAuthority***REMOVED***[6]byte***REMOVED***0, 0, 0, 0, 0, 4***REMOVED******REMOVED***
	SECURITY_NT_AUTHORITY              = SidIdentifierAuthority***REMOVED***[6]byte***REMOVED***0, 0, 0, 0, 0, 5***REMOVED******REMOVED***
	SECURITY_MANDATORY_LABEL_AUTHORITY = SidIdentifierAuthority***REMOVED***[6]byte***REMOVED***0, 0, 0, 0, 0, 16***REMOVED******REMOVED***
)

const (
	SECURITY_NULL_RID                   = 0
	SECURITY_WORLD_RID                  = 0
	SECURITY_LOCAL_RID                  = 0
	SECURITY_CREATOR_OWNER_RID          = 0
	SECURITY_CREATOR_GROUP_RID          = 1
	SECURITY_DIALUP_RID                 = 1
	SECURITY_NETWORK_RID                = 2
	SECURITY_BATCH_RID                  = 3
	SECURITY_INTERACTIVE_RID            = 4
	SECURITY_LOGON_IDS_RID              = 5
	SECURITY_SERVICE_RID                = 6
	SECURITY_LOCAL_SYSTEM_RID           = 18
	SECURITY_BUILTIN_DOMAIN_RID         = 32
	SECURITY_PRINCIPAL_SELF_RID         = 10
	SECURITY_CREATOR_OWNER_SERVER_RID   = 0x2
	SECURITY_CREATOR_GROUP_SERVER_RID   = 0x3
	SECURITY_LOGON_IDS_RID_COUNT        = 0x3
	SECURITY_ANONYMOUS_LOGON_RID        = 0x7
	SECURITY_PROXY_RID                  = 0x8
	SECURITY_ENTERPRISE_CONTROLLERS_RID = 0x9
	SECURITY_SERVER_LOGON_RID           = SECURITY_ENTERPRISE_CONTROLLERS_RID
	SECURITY_AUTHENTICATED_USER_RID     = 0xb
	SECURITY_RESTRICTED_CODE_RID        = 0xc
	SECURITY_NT_NON_UNIQUE_RID          = 0x15
)

//sys	LookupAccountSid(systemName *uint16, sid *SID, name *uint16, nameLen *uint32, refdDomainName *uint16, refdDomainNameLen *uint32, use *uint32) (err error) = advapi32.LookupAccountSidW
//sys	LookupAccountName(systemName *uint16, accountName *uint16, sid *SID, sidLen *uint32, refdDomainName *uint16, refdDomainNameLen *uint32, use *uint32) (err error) = advapi32.LookupAccountNameW
//sys	ConvertSidToStringSid(sid *SID, stringSid **uint16) (err error) = advapi32.ConvertSidToStringSidW
//sys	ConvertStringSidToSid(stringSid *uint16, sid **SID) (err error) = advapi32.ConvertStringSidToSidW
//sys	GetLengthSid(sid *SID) (len uint32) = advapi32.GetLengthSid
//sys	CopySid(destSidLen uint32, destSid *SID, srcSid *SID) (err error) = advapi32.CopySid
//sys	AllocateAndInitializeSid(identAuth *SidIdentifierAuthority, subAuth byte, subAuth0 uint32, subAuth1 uint32, subAuth2 uint32, subAuth3 uint32, subAuth4 uint32, subAuth5 uint32, subAuth6 uint32, subAuth7 uint32, sid **SID) (err error) = advapi32.AllocateAndInitializeSid
//sys	FreeSid(sid *SID) (err error) [failretval!=0] = advapi32.FreeSid
//sys	EqualSid(sid1 *SID, sid2 *SID) (isEqual bool) = advapi32.EqualSid

// The security identifier (SID) structure is a variable-length
// structure used to uniquely identify users or groups.
type SID struct***REMOVED******REMOVED***

// StringToSid converts a string-format security identifier
// sid into a valid, functional sid.
func StringToSid(s string) (*SID, error) ***REMOVED***
	var sid *SID
	p, e := UTF16PtrFromString(s)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***
	e = ConvertStringSidToSid(p, &sid)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***
	defer LocalFree((Handle)(unsafe.Pointer(sid)))
	return sid.Copy()
***REMOVED***

// LookupSID retrieves a security identifier sid for the account
// and the name of the domain on which the account was found.
// System specify target computer to search.
func LookupSID(system, account string) (sid *SID, domain string, accType uint32, err error) ***REMOVED***
	if len(account) == 0 ***REMOVED***
		return nil, "", 0, syscall.EINVAL
	***REMOVED***
	acc, e := UTF16PtrFromString(account)
	if e != nil ***REMOVED***
		return nil, "", 0, e
	***REMOVED***
	var sys *uint16
	if len(system) > 0 ***REMOVED***
		sys, e = UTF16PtrFromString(system)
		if e != nil ***REMOVED***
			return nil, "", 0, e
		***REMOVED***
	***REMOVED***
	n := uint32(50)
	dn := uint32(50)
	for ***REMOVED***
		b := make([]byte, n)
		db := make([]uint16, dn)
		sid = (*SID)(unsafe.Pointer(&b[0]))
		e = LookupAccountName(sys, acc, sid, &n, &db[0], &dn, &accType)
		if e == nil ***REMOVED***
			return sid, UTF16ToString(db), accType, nil
		***REMOVED***
		if e != ERROR_INSUFFICIENT_BUFFER ***REMOVED***
			return nil, "", 0, e
		***REMOVED***
		if n <= uint32(len(b)) ***REMOVED***
			return nil, "", 0, e
		***REMOVED***
	***REMOVED***
***REMOVED***

// String converts sid to a string format
// suitable for display, storage, or transmission.
func (sid *SID) String() (string, error) ***REMOVED***
	var s *uint16
	e := ConvertSidToStringSid(sid, &s)
	if e != nil ***REMOVED***
		return "", e
	***REMOVED***
	defer LocalFree((Handle)(unsafe.Pointer(s)))
	return UTF16ToString((*[256]uint16)(unsafe.Pointer(s))[:]), nil
***REMOVED***

// Len returns the length, in bytes, of a valid security identifier sid.
func (sid *SID) Len() int ***REMOVED***
	return int(GetLengthSid(sid))
***REMOVED***

// Copy creates a duplicate of security identifier sid.
func (sid *SID) Copy() (*SID, error) ***REMOVED***
	b := make([]byte, sid.Len())
	sid2 := (*SID)(unsafe.Pointer(&b[0]))
	e := CopySid(uint32(len(b)), sid2, sid)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***
	return sid2, nil
***REMOVED***

// LookupAccount retrieves the name of the account for this sid
// and the name of the first domain on which this sid is found.
// System specify target computer to search for.
func (sid *SID) LookupAccount(system string) (account, domain string, accType uint32, err error) ***REMOVED***
	var sys *uint16
	if len(system) > 0 ***REMOVED***
		sys, err = UTF16PtrFromString(system)
		if err != nil ***REMOVED***
			return "", "", 0, err
		***REMOVED***
	***REMOVED***
	n := uint32(50)
	dn := uint32(50)
	for ***REMOVED***
		b := make([]uint16, n)
		db := make([]uint16, dn)
		e := LookupAccountSid(sys, sid, &b[0], &n, &db[0], &dn, &accType)
		if e == nil ***REMOVED***
			return UTF16ToString(b), UTF16ToString(db), accType, nil
		***REMOVED***
		if e != ERROR_INSUFFICIENT_BUFFER ***REMOVED***
			return "", "", 0, e
		***REMOVED***
		if n <= uint32(len(b)) ***REMOVED***
			return "", "", 0, e
		***REMOVED***
	***REMOVED***
***REMOVED***

const (
	// do not reorder
	TOKEN_ASSIGN_PRIMARY = 1 << iota
	TOKEN_DUPLICATE
	TOKEN_IMPERSONATE
	TOKEN_QUERY
	TOKEN_QUERY_SOURCE
	TOKEN_ADJUST_PRIVILEGES
	TOKEN_ADJUST_GROUPS
	TOKEN_ADJUST_DEFAULT

	TOKEN_ALL_ACCESS = STANDARD_RIGHTS_REQUIRED |
		TOKEN_ASSIGN_PRIMARY |
		TOKEN_DUPLICATE |
		TOKEN_IMPERSONATE |
		TOKEN_QUERY |
		TOKEN_QUERY_SOURCE |
		TOKEN_ADJUST_PRIVILEGES |
		TOKEN_ADJUST_GROUPS |
		TOKEN_ADJUST_DEFAULT
	TOKEN_READ  = STANDARD_RIGHTS_READ | TOKEN_QUERY
	TOKEN_WRITE = STANDARD_RIGHTS_WRITE |
		TOKEN_ADJUST_PRIVILEGES |
		TOKEN_ADJUST_GROUPS |
		TOKEN_ADJUST_DEFAULT
	TOKEN_EXECUTE = STANDARD_RIGHTS_EXECUTE
)

const (
	// do not reorder
	TokenUser = 1 + iota
	TokenGroups
	TokenPrivileges
	TokenOwner
	TokenPrimaryGroup
	TokenDefaultDacl
	TokenSource
	TokenType
	TokenImpersonationLevel
	TokenStatistics
	TokenRestrictedSids
	TokenSessionId
	TokenGroupsAndPrivileges
	TokenSessionReference
	TokenSandBoxInert
	TokenAuditPolicy
	TokenOrigin
	TokenElevationType
	TokenLinkedToken
	TokenElevation
	TokenHasRestrictions
	TokenAccessInformation
	TokenVirtualizationAllowed
	TokenVirtualizationEnabled
	TokenIntegrityLevel
	TokenUIAccess
	TokenMandatoryPolicy
	TokenLogonSid
	MaxTokenInfoClass
)

type SIDAndAttributes struct ***REMOVED***
	Sid        *SID
	Attributes uint32
***REMOVED***

type Tokenuser struct ***REMOVED***
	User SIDAndAttributes
***REMOVED***

type Tokenprimarygroup struct ***REMOVED***
	PrimaryGroup *SID
***REMOVED***

type Tokengroups struct ***REMOVED***
	GroupCount uint32
	Groups     [1]SIDAndAttributes
***REMOVED***

//sys	OpenProcessToken(h Handle, access uint32, token *Token) (err error) = advapi32.OpenProcessToken
//sys	GetTokenInformation(t Token, infoClass uint32, info *byte, infoLen uint32, returnedLen *uint32) (err error) = advapi32.GetTokenInformation
//sys	GetUserProfileDirectory(t Token, dir *uint16, dirLen *uint32) (err error) = userenv.GetUserProfileDirectoryW

// An access token contains the security information for a logon session.
// The system creates an access token when a user logs on, and every
// process executed on behalf of the user has a copy of the token.
// The token identifies the user, the user's groups, and the user's
// privileges. The system uses the token to control access to securable
// objects and to control the ability of the user to perform various
// system-related operations on the local computer.
type Token Handle

// OpenCurrentProcessToken opens the access token
// associated with current process.
func OpenCurrentProcessToken() (Token, error) ***REMOVED***
	p, e := GetCurrentProcess()
	if e != nil ***REMOVED***
		return 0, e
	***REMOVED***
	var t Token
	e = OpenProcessToken(p, TOKEN_QUERY, &t)
	if e != nil ***REMOVED***
		return 0, e
	***REMOVED***
	return t, nil
***REMOVED***

// Close releases access to access token.
func (t Token) Close() error ***REMOVED***
	return CloseHandle(Handle(t))
***REMOVED***

// getInfo retrieves a specified type of information about an access token.
func (t Token) getInfo(class uint32, initSize int) (unsafe.Pointer, error) ***REMOVED***
	n := uint32(initSize)
	for ***REMOVED***
		b := make([]byte, n)
		e := GetTokenInformation(t, class, &b[0], uint32(len(b)), &n)
		if e == nil ***REMOVED***
			return unsafe.Pointer(&b[0]), nil
		***REMOVED***
		if e != ERROR_INSUFFICIENT_BUFFER ***REMOVED***
			return nil, e
		***REMOVED***
		if n <= uint32(len(b)) ***REMOVED***
			return nil, e
		***REMOVED***
	***REMOVED***
***REMOVED***

// GetTokenUser retrieves access token t user account information.
func (t Token) GetTokenUser() (*Tokenuser, error) ***REMOVED***
	i, e := t.getInfo(TokenUser, 50)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***
	return (*Tokenuser)(i), nil
***REMOVED***

// GetTokenGroups retrieves group accounts associated with access token t.
func (t Token) GetTokenGroups() (*Tokengroups, error) ***REMOVED***
	i, e := t.getInfo(TokenGroups, 50)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***
	return (*Tokengroups)(i), nil
***REMOVED***

// GetTokenPrimaryGroup retrieves access token t primary group information.
// A pointer to a SID structure representing a group that will become
// the primary group of any objects created by a process using this access token.
func (t Token) GetTokenPrimaryGroup() (*Tokenprimarygroup, error) ***REMOVED***
	i, e := t.getInfo(TokenPrimaryGroup, 50)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***
	return (*Tokenprimarygroup)(i), nil
***REMOVED***

// GetUserProfileDirectory retrieves path to the
// root directory of the access token t user's profile.
func (t Token) GetUserProfileDirectory() (string, error) ***REMOVED***
	n := uint32(100)
	for ***REMOVED***
		b := make([]uint16, n)
		e := GetUserProfileDirectory(t, &b[0], &n)
		if e == nil ***REMOVED***
			return UTF16ToString(b), nil
		***REMOVED***
		if e != ERROR_INSUFFICIENT_BUFFER ***REMOVED***
			return "", e
		***REMOVED***
		if n <= uint32(len(b)) ***REMOVED***
			return "", e
		***REMOVED***
	***REMOVED***
***REMOVED***
