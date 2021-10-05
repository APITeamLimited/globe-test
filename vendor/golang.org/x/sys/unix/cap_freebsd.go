// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd
// +build freebsd

package unix

import (
	"errors"
	"fmt"
)

// Go implementation of C mostly found in /usr/src/sys/kern/subr_capability.c

const (
	// This is the version of CapRights this package understands. See C implementation for parallels.
	capRightsGoVersion = CAP_RIGHTS_VERSION_00
	capArSizeMin       = CAP_RIGHTS_VERSION_00 + 2
	capArSizeMax       = capRightsGoVersion + 2
)

var (
	bit2idx = []int***REMOVED***
		-1, 0, 1, -1, 2, -1, -1, -1, 3, -1, -1, -1, -1, -1, -1, -1,
		4, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	***REMOVED***
)

func capidxbit(right uint64) int ***REMOVED***
	return int((right >> 57) & 0x1f)
***REMOVED***

func rightToIndex(right uint64) (int, error) ***REMOVED***
	idx := capidxbit(right)
	if idx < 0 || idx >= len(bit2idx) ***REMOVED***
		return -2, fmt.Errorf("index for right 0x%x out of range", right)
	***REMOVED***
	return bit2idx[idx], nil
***REMOVED***

func caprver(right uint64) int ***REMOVED***
	return int(right >> 62)
***REMOVED***

func capver(rights *CapRights) int ***REMOVED***
	return caprver(rights.Rights[0])
***REMOVED***

func caparsize(rights *CapRights) int ***REMOVED***
	return capver(rights) + 2
***REMOVED***

// CapRightsSet sets the permissions in setrights in rights.
func CapRightsSet(rights *CapRights, setrights []uint64) error ***REMOVED***
	// This is essentially a copy of cap_rights_vset()
	if capver(rights) != CAP_RIGHTS_VERSION_00 ***REMOVED***
		return fmt.Errorf("bad rights version %d", capver(rights))
	***REMOVED***

	n := caparsize(rights)
	if n < capArSizeMin || n > capArSizeMax ***REMOVED***
		return errors.New("bad rights size")
	***REMOVED***

	for _, right := range setrights ***REMOVED***
		if caprver(right) != CAP_RIGHTS_VERSION_00 ***REMOVED***
			return errors.New("bad right version")
		***REMOVED***
		i, err := rightToIndex(right)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if i >= n ***REMOVED***
			return errors.New("index overflow")
		***REMOVED***
		if capidxbit(rights.Rights[i]) != capidxbit(right) ***REMOVED***
			return errors.New("index mismatch")
		***REMOVED***
		rights.Rights[i] |= right
		if capidxbit(rights.Rights[i]) != capidxbit(right) ***REMOVED***
			return errors.New("index mismatch (after assign)")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// CapRightsClear clears the permissions in clearrights from rights.
func CapRightsClear(rights *CapRights, clearrights []uint64) error ***REMOVED***
	// This is essentially a copy of cap_rights_vclear()
	if capver(rights) != CAP_RIGHTS_VERSION_00 ***REMOVED***
		return fmt.Errorf("bad rights version %d", capver(rights))
	***REMOVED***

	n := caparsize(rights)
	if n < capArSizeMin || n > capArSizeMax ***REMOVED***
		return errors.New("bad rights size")
	***REMOVED***

	for _, right := range clearrights ***REMOVED***
		if caprver(right) != CAP_RIGHTS_VERSION_00 ***REMOVED***
			return errors.New("bad right version")
		***REMOVED***
		i, err := rightToIndex(right)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if i >= n ***REMOVED***
			return errors.New("index overflow")
		***REMOVED***
		if capidxbit(rights.Rights[i]) != capidxbit(right) ***REMOVED***
			return errors.New("index mismatch")
		***REMOVED***
		rights.Rights[i] &= ^(right & 0x01FFFFFFFFFFFFFF)
		if capidxbit(rights.Rights[i]) != capidxbit(right) ***REMOVED***
			return errors.New("index mismatch (after assign)")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// CapRightsIsSet checks whether all the permissions in setrights are present in rights.
func CapRightsIsSet(rights *CapRights, setrights []uint64) (bool, error) ***REMOVED***
	// This is essentially a copy of cap_rights_is_vset()
	if capver(rights) != CAP_RIGHTS_VERSION_00 ***REMOVED***
		return false, fmt.Errorf("bad rights version %d", capver(rights))
	***REMOVED***

	n := caparsize(rights)
	if n < capArSizeMin || n > capArSizeMax ***REMOVED***
		return false, errors.New("bad rights size")
	***REMOVED***

	for _, right := range setrights ***REMOVED***
		if caprver(right) != CAP_RIGHTS_VERSION_00 ***REMOVED***
			return false, errors.New("bad right version")
		***REMOVED***
		i, err := rightToIndex(right)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		if i >= n ***REMOVED***
			return false, errors.New("index overflow")
		***REMOVED***
		if capidxbit(rights.Rights[i]) != capidxbit(right) ***REMOVED***
			return false, errors.New("index mismatch")
		***REMOVED***
		if (rights.Rights[i] & right) != right ***REMOVED***
			return false, nil
		***REMOVED***
	***REMOVED***

	return true, nil
***REMOVED***

func capright(idx uint64, bit uint64) uint64 ***REMOVED***
	return ((1 << (57 + idx)) | bit)
***REMOVED***

// CapRightsInit returns a pointer to an initialised CapRights structure filled with rights.
// See man cap_rights_init(3) and rights(4).
func CapRightsInit(rights []uint64) (*CapRights, error) ***REMOVED***
	var r CapRights
	r.Rights[0] = (capRightsGoVersion << 62) | capright(0, 0)
	r.Rights[1] = capright(1, 0)

	err := CapRightsSet(&r, rights)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &r, nil
***REMOVED***

// CapRightsLimit reduces the operations permitted on fd to at most those contained in rights.
// The capability rights on fd can never be increased by CapRightsLimit.
// See man cap_rights_limit(2) and rights(4).
func CapRightsLimit(fd uintptr, rights *CapRights) error ***REMOVED***
	return capRightsLimit(int(fd), rights)
***REMOVED***

// CapRightsGet returns a CapRights structure containing the operations permitted on fd.
// See man cap_rights_get(3) and rights(4).
func CapRightsGet(fd uintptr) (*CapRights, error) ***REMOVED***
	r, err := CapRightsInit(nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = capRightsGet(capRightsGoVersion, int(fd), r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r, nil
***REMOVED***
