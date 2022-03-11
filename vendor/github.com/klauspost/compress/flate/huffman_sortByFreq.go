// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flate

// Sort sorts data.
// It makes one call to data.Len to determine n, and O(n*log(n)) calls to
// data.Less and data.Swap. The sort is not guaranteed to be stable.
func sortByFreq(data []literalNode) ***REMOVED***
	n := len(data)
	quickSortByFreq(data, 0, n, maxDepth(n))
***REMOVED***

func quickSortByFreq(data []literalNode, a, b, maxDepth int) ***REMOVED***
	for b-a > 12 ***REMOVED*** // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 ***REMOVED***
			heapSort(data, a, b)
			return
		***REMOVED***
		maxDepth--
		mlo, mhi := doPivotByFreq(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi ***REMOVED***
			quickSortByFreq(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSortByFreq(data, mhi, b)
		***REMOVED*** else ***REMOVED***
			quickSortByFreq(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSortByFreq(data, a, mlo)
		***REMOVED***
	***REMOVED***
	if b-a > 1 ***REMOVED***
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ ***REMOVED***
			if data[i].freq == data[i-6].freq && data[i].literal < data[i-6].literal || data[i].freq < data[i-6].freq ***REMOVED***
				data[i], data[i-6] = data[i-6], data[i]
			***REMOVED***
		***REMOVED***
		insertionSortByFreq(data, a, b)
	***REMOVED***
***REMOVED***

// siftDownByFreq implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDownByFreq(data []literalNode, lo, hi, first int) ***REMOVED***
	root := lo
	for ***REMOVED***
		child := 2*root + 1
		if child >= hi ***REMOVED***
			break
		***REMOVED***
		if child+1 < hi && (data[first+child].freq == data[first+child+1].freq && data[first+child].literal < data[first+child+1].literal || data[first+child].freq < data[first+child+1].freq) ***REMOVED***
			child++
		***REMOVED***
		if data[first+root].freq == data[first+child].freq && data[first+root].literal > data[first+child].literal || data[first+root].freq > data[first+child].freq ***REMOVED***
			return
		***REMOVED***
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	***REMOVED***
***REMOVED***
func doPivotByFreq(data []literalNode, lo, hi int) (midlo, midhi int) ***REMOVED***
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	if hi-lo > 40 ***REMOVED***
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThreeSortByFreq(data, lo, lo+s, lo+2*s)
		medianOfThreeSortByFreq(data, m, m-s, m+s)
		medianOfThreeSortByFreq(data, hi-1, hi-1-s, hi-1-2*s)
	***REMOVED***
	medianOfThreeSortByFreq(data, lo, m, hi-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && (data[a].freq == data[pivot].freq && data[a].literal < data[pivot].literal || data[a].freq < data[pivot].freq); a++ ***REMOVED***
	***REMOVED***
	b := a
	for ***REMOVED***
		for ; b < c && (data[pivot].freq == data[b].freq && data[pivot].literal > data[b].literal || data[pivot].freq > data[b].freq); b++ ***REMOVED*** // data[b] <= pivot
		***REMOVED***
		for ; b < c && (data[pivot].freq == data[c-1].freq && data[pivot].literal < data[c-1].literal || data[pivot].freq < data[c-1].freq); c-- ***REMOVED*** // data[c-1] > pivot
		***REMOVED***
		if b >= c ***REMOVED***
			break
		***REMOVED***
		// data[b] > pivot; data[c-1] <= pivot
		data[b], data[c-1] = data[c-1], data[b]
		b++
		c--
	***REMOVED***
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let's be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 ***REMOVED***
		// Lets test some points for equality to pivot
		dups := 0
		if data[pivot].freq == data[hi-1].freq && data[pivot].literal > data[hi-1].literal || data[pivot].freq > data[hi-1].freq ***REMOVED*** // data[hi-1] = pivot
			data[c], data[hi-1] = data[hi-1], data[c]
			c++
			dups++
		***REMOVED***
		if data[b-1].freq == data[pivot].freq && data[b-1].literal > data[pivot].literal || data[b-1].freq > data[pivot].freq ***REMOVED*** // data[b-1] = pivot
			b--
			dups++
		***REMOVED***
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if data[m].freq == data[pivot].freq && data[m].literal > data[pivot].literal || data[m].freq > data[pivot].freq ***REMOVED*** // data[m] = pivot
			data[m], data[b-1] = data[b-1], data[m]
			b--
			dups++
		***REMOVED***
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	***REMOVED***
	if protect ***REMOVED***
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for ***REMOVED***
			for ; a < b && (data[b-1].freq == data[pivot].freq && data[b-1].literal > data[pivot].literal || data[b-1].freq > data[pivot].freq); b-- ***REMOVED*** // data[b] == pivot
			***REMOVED***
			for ; a < b && (data[a].freq == data[pivot].freq && data[a].literal < data[pivot].literal || data[a].freq < data[pivot].freq); a++ ***REMOVED*** // data[a] < pivot
			***REMOVED***
			if a >= b ***REMOVED***
				break
			***REMOVED***
			// data[a] == pivot; data[b-1] < pivot
			data[a], data[b-1] = data[b-1], data[a]
			a++
			b--
		***REMOVED***
	***REMOVED***
	// Swap pivot into middle
	data[pivot], data[b-1] = data[b-1], data[pivot]
	return b - 1, c
***REMOVED***

// Insertion sort
func insertionSortByFreq(data []literalNode, a, b int) ***REMOVED***
	for i := a + 1; i < b; i++ ***REMOVED***
		for j := i; j > a && (data[j].freq == data[j-1].freq && data[j].literal < data[j-1].literal || data[j].freq < data[j-1].freq); j-- ***REMOVED***
			data[j], data[j-1] = data[j-1], data[j]
		***REMOVED***
	***REMOVED***
***REMOVED***

// quickSortByFreq, loosely following Bentley and McIlroy,
// ``Engineering a Sort Function,'' SP&E November 1993.

// medianOfThreeSortByFreq moves the median of the three values data[m0], data[m1], data[m2] into data[m1].
func medianOfThreeSortByFreq(data []literalNode, m1, m0, m2 int) ***REMOVED***
	// sort 3 elements
	if data[m1].freq == data[m0].freq && data[m1].literal < data[m0].literal || data[m1].freq < data[m0].freq ***REMOVED***
		data[m1], data[m0] = data[m0], data[m1]
	***REMOVED***
	// data[m0] <= data[m1]
	if data[m2].freq == data[m1].freq && data[m2].literal < data[m1].literal || data[m2].freq < data[m1].freq ***REMOVED***
		data[m2], data[m1] = data[m1], data[m2]
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if data[m1].freq == data[m0].freq && data[m1].literal < data[m0].literal || data[m1].freq < data[m0].freq ***REMOVED***
			data[m1], data[m0] = data[m0], data[m1]
		***REMOVED***
	***REMOVED***
	// now data[m0] <= data[m1] <= data[m2]
***REMOVED***
