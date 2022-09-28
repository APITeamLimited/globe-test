// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// CPU affinity functions

package unix

import (
	"math/bits"
	"unsafe"
)

const cpuSetSize = _CPU_SETSIZE / _NCPUBITS

// CPUSet represents a CPU affinity mask.
type CPUSet [cpuSetSize]cpuMask

func schedAffinity(trap uintptr, pid int, set *CPUSet) error ***REMOVED***
	_, _, e := RawSyscall(trap, uintptr(pid), uintptr(unsafe.Sizeof(*set)), uintptr(unsafe.Pointer(set)))
	if e != 0 ***REMOVED***
		return errnoErr(e)
	***REMOVED***
	return nil
***REMOVED***

// SchedGetaffinity gets the CPU affinity mask of the thread specified by pid.
// If pid is 0 the calling thread is used.
func SchedGetaffinity(pid int, set *CPUSet) error ***REMOVED***
	return schedAffinity(SYS_SCHED_GETAFFINITY, pid, set)
***REMOVED***

// SchedSetaffinity sets the CPU affinity mask of the thread specified by pid.
// If pid is 0 the calling thread is used.
func SchedSetaffinity(pid int, set *CPUSet) error ***REMOVED***
	return schedAffinity(SYS_SCHED_SETAFFINITY, pid, set)
***REMOVED***

// Zero clears the set s, so that it contains no CPUs.
func (s *CPUSet) Zero() ***REMOVED***
	for i := range s ***REMOVED***
		s[i] = 0
	***REMOVED***
***REMOVED***

func cpuBitsIndex(cpu int) int ***REMOVED***
	return cpu / _NCPUBITS
***REMOVED***

func cpuBitsMask(cpu int) cpuMask ***REMOVED***
	return cpuMask(1 << (uint(cpu) % _NCPUBITS))
***REMOVED***

// Set adds cpu to the set s.
func (s *CPUSet) Set(cpu int) ***REMOVED***
	i := cpuBitsIndex(cpu)
	if i < len(s) ***REMOVED***
		s[i] |= cpuBitsMask(cpu)
	***REMOVED***
***REMOVED***

// Clear removes cpu from the set s.
func (s *CPUSet) Clear(cpu int) ***REMOVED***
	i := cpuBitsIndex(cpu)
	if i < len(s) ***REMOVED***
		s[i] &^= cpuBitsMask(cpu)
	***REMOVED***
***REMOVED***

// IsSet reports whether cpu is in the set s.
func (s *CPUSet) IsSet(cpu int) bool ***REMOVED***
	i := cpuBitsIndex(cpu)
	if i < len(s) ***REMOVED***
		return s[i]&cpuBitsMask(cpu) != 0
	***REMOVED***
	return false
***REMOVED***

// Count returns the number of CPUs in the set s.
func (s *CPUSet) Count() int ***REMOVED***
	c := 0
	for _, b := range s ***REMOVED***
		c += bits.OnesCount64(uint64(b))
	***REMOVED***
	return c
***REMOVED***
