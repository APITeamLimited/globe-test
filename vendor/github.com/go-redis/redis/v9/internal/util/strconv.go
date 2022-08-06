package util

import "strconv"

func Atoi(b []byte) (int, error) ***REMOVED***
	return strconv.Atoi(BytesToString(b))
***REMOVED***

func ParseInt(b []byte, base int, bitSize int) (int64, error) ***REMOVED***
	return strconv.ParseInt(BytesToString(b), base, bitSize)
***REMOVED***

func ParseUint(b []byte, base int, bitSize int) (uint64, error) ***REMOVED***
	return strconv.ParseUint(BytesToString(b), base, bitSize)
***REMOVED***

func ParseFloat(b []byte, bitSize int) (float64, error) ***REMOVED***
	return strconv.ParseFloat(BytesToString(b), bitSize)
***REMOVED***
