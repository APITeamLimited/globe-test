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

package lib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"crypto/tls"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

const groupSeparator = "::"

var ErrNameContainsGroupSeparator = errors.Errorf("group and check names may not contain '%s'", groupSeparator)

type SourceData struct ***REMOVED***
	Data     []byte
	Filename string
***REMOVED***

type Stage struct ***REMOVED***
	Duration time.Duration `json:"duration"`
	Target   null.Int      `json:"target"`
***REMOVED***

func (s *Stage) UnmarshalJSON(data []byte) error ***REMOVED***
	var fields struct ***REMOVED***
		Duration string   `json:"duration"`
		Target   null.Int `json:"target"`
	***REMOVED***
	if err := json.Unmarshal(data, &fields); err != nil ***REMOVED***
		return err
	***REMOVED***

	s.Target = fields.Target

	if fields.Duration != "" ***REMOVED***
		d, err := time.ParseDuration(fields.Duration)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		s.Duration = d
	***REMOVED***

	return nil
***REMOVED***

type TLSVersion struct ***REMOVED***
	Min int
	Max int
***REMOVED***

func (v *TLSVersion) UnmarshalJSON(data []byte) error ***REMOVED***
	// From https://golang.org/pkg/crypto/tls/#pkg-constants
	versionMap := map[string]int***REMOVED***
		"ssl3.0": tls.VersionSSL30,
		"tls1.0": tls.VersionTLS10,
		"tls1.1": tls.VersionTLS11,
		"tls1.2": tls.VersionTLS12,
	***REMOVED***

	// Version might be a string or an object with separate min & max fields
	var fields struct ***REMOVED***
		Min string `json:"min"`
		Max string `json:"max"`
	***REMOVED***
	if err := json.Unmarshal(data, &fields); err != nil ***REMOVED***
		switch err.(type) ***REMOVED***
		case *json.UnmarshalTypeError:
			// Check if it's a type error and the user has passed a string
			var version string
			if otherErr := json.Unmarshal(data, &version); otherErr != nil ***REMOVED***
				// Some other error occurred or none of the types match
				return otherErr
			***REMOVED***
			// It was a string, assign it to both min & max
			fields.Min = version
			fields.Max = version
		default:
			return err
		***REMOVED***
	***REMOVED***

	var minVersion int
	var maxVersion int
	var ok bool
	if minVersion, ok = versionMap[fields.Min]; !ok ***REMOVED***
		return errors.New("Unknown TLS version : " + fields.Min)
	***REMOVED***

	if maxVersion, ok = versionMap[fields.Max]; !ok ***REMOVED***
		return errors.New("Unknown TLS version : " + fields.Max)
	***REMOVED***

	v.Min = minVersion
	v.Max = maxVersion

	return nil
***REMOVED***

type TLSCipherSuites struct ***REMOVED***
	Values []uint16
***REMOVED***

func (s *TLSCipherSuites) UnmarshalJSON(data []byte) error ***REMOVED***
	// From https://golang.org/pkg/crypto/tls#pkg-constants
	suiteMap := map[string]uint16***REMOVED***
		"TLS_RSA_WITH_RC4_128_SHA":                0x0005,
		"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           0x000a,
		"TLS_RSA_WITH_AES_128_CBC_SHA":            0x002f,
		"TLS_RSA_WITH_AES_256_CBC_SHA":            0x0035,
		"TLS_RSA_WITH_AES_128_CBC_SHA256":         0x003c,
		"TLS_RSA_WITH_AES_128_GCM_SHA256":         0x009c,
		"TLS_RSA_WITH_AES_256_GCM_SHA384":         0x009d,
		"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        0xc007,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    0xc009,
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    0xc00a,
		"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          0xc011,
		"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     0xc012,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      0xc013,
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      0xc014,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": 0xc023,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   0xc027,
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   0xc02f,
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": 0xc02b,
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   0xc030,
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": 0xc02c,
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    0xcca8,
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  0xcca9,
	***REMOVED***

	var suiteNames []string
	if err := json.Unmarshal(data, &suiteNames); err != nil ***REMOVED***
		return err
	***REMOVED***

	var suiteIDs []uint16
	for _, name := range suiteNames ***REMOVED***
		if suiteID, ok := suiteMap[name]; ok ***REMOVED***
			suiteIDs = append(suiteIDs, suiteID)
		***REMOVED*** else ***REMOVED***
			return errors.New("Unknown cipher suite: " + name)
		***REMOVED***
	***REMOVED***

	s.Values = suiteIDs

	return nil
***REMOVED***

type Group struct ***REMOVED***
	ID     string            `json:"id"`
	Path   string            `json:"path"`
	Name   string            `json:"name"`
	Parent *Group            `json:"parent"`
	Groups map[string]*Group `json:"groups"`
	Checks map[string]*Check `json:"checks"`

	groupMutex sync.Mutex
	checkMutex sync.Mutex
***REMOVED***

func NewGroup(name string, parent *Group) (*Group, error) ***REMOVED***
	if strings.Contains(name, groupSeparator) ***REMOVED***
		return nil, ErrNameContainsGroupSeparator
	***REMOVED***

	path := name
	if parent != nil ***REMOVED***
		path = parent.Path + groupSeparator + path
	***REMOVED***

	hash := md5.Sum([]byte(path))
	id := hex.EncodeToString(hash[:])

	return &Group***REMOVED***
		ID:     id,
		Path:   path,
		Name:   name,
		Parent: parent,
		Groups: make(map[string]*Group),
		Checks: make(map[string]*Check),
	***REMOVED***, nil
***REMOVED***

func (g *Group) Group(name string) (*Group, error) ***REMOVED***
	snapshot := g.Groups
	group, ok := snapshot[name]
	if !ok ***REMOVED***
		g.groupMutex.Lock()
		defer g.groupMutex.Unlock()

		group, ok := g.Groups[name]
		if !ok ***REMOVED***
			group, err := NewGroup(name, g)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			g.Groups[name] = group
			return group, nil
		***REMOVED***
		return group, nil
	***REMOVED***
	return group, nil
***REMOVED***

func (g *Group) Check(name string) (*Check, error) ***REMOVED***
	snapshot := g.Checks
	check, ok := snapshot[name]
	if !ok ***REMOVED***
		g.checkMutex.Lock()
		defer g.checkMutex.Unlock()
		check, ok := g.Checks[name]
		if !ok ***REMOVED***
			check, err := NewCheck(name, g)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			g.Checks[name] = check
			return check, nil
		***REMOVED***
		return check, nil
	***REMOVED***
	return check, nil
***REMOVED***

type Check struct ***REMOVED***
	ID    string `json:"id"`
	Path  string `json:"path"`
	Group *Group `json:"group"`
	Name  string `json:"name"`

	Passes int64 `json:"passes"`
	Fails  int64 `json:"fails"`
***REMOVED***

func NewCheck(name string, group *Group) (*Check, error) ***REMOVED***
	if strings.Contains(name, groupSeparator) ***REMOVED***
		return nil, ErrNameContainsGroupSeparator
	***REMOVED***

	path := group.Path + groupSeparator + name
	hash := md5.Sum([]byte(path))
	id := hex.EncodeToString(hash[:])

	return &Check***REMOVED***
		ID:    id,
		Path:  path,
		Group: group,
		Name:  name,
	***REMOVED***, nil
***REMOVED***
