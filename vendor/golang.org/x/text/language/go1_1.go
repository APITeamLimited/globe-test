// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.2
// +build !go1.2

package language

import "sort"

func sortStable(s sort.Interface) ***REMOVED***
	ss := stableSort***REMOVED***
		s:   s,
		pos: make([]int, s.Len()),
	***REMOVED***
	for i := range ss.pos ***REMOVED***
		ss.pos[i] = i
	***REMOVED***
	sort.Sort(&ss)
***REMOVED***

type stableSort struct ***REMOVED***
	s   sort.Interface
	pos []int
***REMOVED***

func (s *stableSort) Len() int ***REMOVED***
	return len(s.pos)
***REMOVED***

func (s *stableSort) Less(i, j int) bool ***REMOVED***
	return s.s.Less(i, j) || !s.s.Less(j, i) && s.pos[i] < s.pos[j]
***REMOVED***

func (s *stableSort) Swap(i, j int) ***REMOVED***
	s.s.Swap(i, j)
	s.pos[i], s.pos[j] = s.pos[j], s.pos[i]
***REMOVED***
