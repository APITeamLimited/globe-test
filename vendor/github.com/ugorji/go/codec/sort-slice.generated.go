// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from sort-slice.go.tmpl - DO NOT EDIT.

package codec

import "time"
import "reflect"
import "bytes"

type stringSlice []string

func (p stringSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p stringSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p stringSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)] < p[uint(j)]
***REMOVED***

type float64Slice []float64

func (p float64Slice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p float64Slice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p float64Slice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)] < p[uint(j)] || isNaN64(p[uint(i)]) && !isNaN64(p[uint(j)])
***REMOVED***

type uint64Slice []uint64

func (p uint64Slice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p uint64Slice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p uint64Slice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)] < p[uint(j)]
***REMOVED***

type uintptrSlice []uintptr

func (p uintptrSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p uintptrSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p uintptrSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)] < p[uint(j)]
***REMOVED***

type int64Slice []int64

func (p int64Slice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p int64Slice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p int64Slice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)] < p[uint(j)]
***REMOVED***

type boolSlice []bool

func (p boolSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p boolSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p boolSlice) Less(i, j int) bool ***REMOVED***
	return !p[uint(i)] && p[uint(j)]
***REMOVED***

type timeSlice []time.Time

func (p timeSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p timeSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p timeSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].Before(p[uint(j)])
***REMOVED***

type bytesSlice [][]byte

func (p bytesSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p bytesSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p bytesSlice) Less(i, j int) bool ***REMOVED***
	return bytes.Compare(p[uint(i)], p[uint(j)]) == -1
***REMOVED***

type stringRv struct ***REMOVED***
	v string
	r reflect.Value
***REMOVED***
type stringRvSlice []stringRv

func (p stringRvSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p stringRvSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p stringRvSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v
***REMOVED***

type stringIntf struct ***REMOVED***
	v string
	i interface***REMOVED******REMOVED***
***REMOVED***
type stringIntfSlice []stringIntf

func (p stringIntfSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p stringIntfSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p stringIntfSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v
***REMOVED***

type float64Rv struct ***REMOVED***
	v float64
	r reflect.Value
***REMOVED***
type float64RvSlice []float64Rv

func (p float64RvSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p float64RvSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p float64RvSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v || isNaN64(p[uint(i)].v) && !isNaN64(p[uint(j)].v)
***REMOVED***

type float64Intf struct ***REMOVED***
	v float64
	i interface***REMOVED******REMOVED***
***REMOVED***
type float64IntfSlice []float64Intf

func (p float64IntfSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p float64IntfSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p float64IntfSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v || isNaN64(p[uint(i)].v) && !isNaN64(p[uint(j)].v)
***REMOVED***

type uint64Rv struct ***REMOVED***
	v uint64
	r reflect.Value
***REMOVED***
type uint64RvSlice []uint64Rv

func (p uint64RvSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p uint64RvSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p uint64RvSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v
***REMOVED***

type uint64Intf struct ***REMOVED***
	v uint64
	i interface***REMOVED******REMOVED***
***REMOVED***
type uint64IntfSlice []uint64Intf

func (p uint64IntfSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p uint64IntfSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p uint64IntfSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v
***REMOVED***

type uintptrRv struct ***REMOVED***
	v uintptr
	r reflect.Value
***REMOVED***
type uintptrRvSlice []uintptrRv

func (p uintptrRvSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p uintptrRvSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p uintptrRvSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v
***REMOVED***

type uintptrIntf struct ***REMOVED***
	v uintptr
	i interface***REMOVED******REMOVED***
***REMOVED***
type uintptrIntfSlice []uintptrIntf

func (p uintptrIntfSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p uintptrIntfSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p uintptrIntfSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v
***REMOVED***

type int64Rv struct ***REMOVED***
	v int64
	r reflect.Value
***REMOVED***
type int64RvSlice []int64Rv

func (p int64RvSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p int64RvSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p int64RvSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v
***REMOVED***

type int64Intf struct ***REMOVED***
	v int64
	i interface***REMOVED******REMOVED***
***REMOVED***
type int64IntfSlice []int64Intf

func (p int64IntfSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p int64IntfSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p int64IntfSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v < p[uint(j)].v
***REMOVED***

type boolRv struct ***REMOVED***
	v bool
	r reflect.Value
***REMOVED***
type boolRvSlice []boolRv

func (p boolRvSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p boolRvSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p boolRvSlice) Less(i, j int) bool ***REMOVED***
	return !p[uint(i)].v && p[uint(j)].v
***REMOVED***

type boolIntf struct ***REMOVED***
	v bool
	i interface***REMOVED******REMOVED***
***REMOVED***
type boolIntfSlice []boolIntf

func (p boolIntfSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p boolIntfSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p boolIntfSlice) Less(i, j int) bool ***REMOVED***
	return !p[uint(i)].v && p[uint(j)].v
***REMOVED***

type timeRv struct ***REMOVED***
	v time.Time
	r reflect.Value
***REMOVED***
type timeRvSlice []timeRv

func (p timeRvSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p timeRvSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p timeRvSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v.Before(p[uint(j)].v)
***REMOVED***

type timeIntf struct ***REMOVED***
	v time.Time
	i interface***REMOVED******REMOVED***
***REMOVED***
type timeIntfSlice []timeIntf

func (p timeIntfSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p timeIntfSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p timeIntfSlice) Less(i, j int) bool ***REMOVED***
	return p[uint(i)].v.Before(p[uint(j)].v)
***REMOVED***

type bytesRv struct ***REMOVED***
	v []byte
	r reflect.Value
***REMOVED***
type bytesRvSlice []bytesRv

func (p bytesRvSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p bytesRvSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p bytesRvSlice) Less(i, j int) bool ***REMOVED***
	return bytes.Compare(p[uint(i)].v, p[uint(j)].v) == -1
***REMOVED***

type bytesIntf struct ***REMOVED***
	v []byte
	i interface***REMOVED******REMOVED***
***REMOVED***
type bytesIntfSlice []bytesIntf

func (p bytesIntfSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p bytesIntfSlice) Swap(i, j int) ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***
func (p bytesIntfSlice) Less(i, j int) bool ***REMOVED***
	return bytes.Compare(p[uint(i)].v, p[uint(j)].v) == -1
***REMOVED***
