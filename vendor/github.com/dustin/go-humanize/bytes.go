package humanize

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// IEC Sizes.
// kibis of bits
const (
	Byte = 1 << (iota * 10)
	KiByte
	MiByte
	GiByte
	TiByte
	PiByte
	EiByte
)

// SI Sizes.
const (
	IByte = 1
	KByte = IByte * 1000
	MByte = KByte * 1000
	GByte = MByte * 1000
	TByte = GByte * 1000
	PByte = TByte * 1000
	EByte = PByte * 1000
)

var bytesSizeTable = map[string]uint64***REMOVED***
	"b":   Byte,
	"kib": KiByte,
	"kb":  KByte,
	"mib": MiByte,
	"mb":  MByte,
	"gib": GiByte,
	"gb":  GByte,
	"tib": TiByte,
	"tb":  TByte,
	"pib": PiByte,
	"pb":  PByte,
	"eib": EiByte,
	"eb":  EByte,
	// Without suffix
	"":   Byte,
	"ki": KiByte,
	"k":  KByte,
	"mi": MiByte,
	"m":  MByte,
	"gi": GiByte,
	"g":  GByte,
	"ti": TiByte,
	"t":  TByte,
	"pi": PiByte,
	"p":  PByte,
	"ei": EiByte,
	"e":  EByte,
***REMOVED***

func logn(n, b float64) float64 ***REMOVED***
	return math.Log(n) / math.Log(b)
***REMOVED***

func humanateBytes(s uint64, base float64, sizes []string) string ***REMOVED***
	if s < 10 ***REMOVED***
		return fmt.Sprintf("%d B", s)
	***REMOVED***
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 ***REMOVED***
		f = "%.1f %s"
	***REMOVED***

	return fmt.Sprintf(f, val, suffix)
***REMOVED***

// Bytes produces a human readable representation of an SI size.
//
// See also: ParseBytes.
//
// Bytes(82854982) -> 83 MB
func Bytes(s uint64) string ***REMOVED***
	sizes := []string***REMOVED***"B", "kB", "MB", "GB", "TB", "PB", "EB"***REMOVED***
	return humanateBytes(s, 1000, sizes)
***REMOVED***

// IBytes produces a human readable representation of an IEC size.
//
// See also: ParseBytes.
//
// IBytes(82854982) -> 79 MiB
func IBytes(s uint64) string ***REMOVED***
	sizes := []string***REMOVED***"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"***REMOVED***
	return humanateBytes(s, 1024, sizes)
***REMOVED***

// ParseBytes parses a string representation of bytes into the number
// of bytes it represents.
//
// See Also: Bytes, IBytes.
//
// ParseBytes("42 MB") -> 42000000, nil
// ParseBytes("42 mib") -> 44040192, nil
func ParseBytes(s string) (uint64, error) ***REMOVED***
	lastDigit := 0
	hasComma := false
	for _, r := range s ***REMOVED***
		if !(unicode.IsDigit(r) || r == '.' || r == ',') ***REMOVED***
			break
		***REMOVED***
		if r == ',' ***REMOVED***
			hasComma = true
		***REMOVED***
		lastDigit++
	***REMOVED***

	num := s[:lastDigit]
	if hasComma ***REMOVED***
		num = strings.Replace(num, ",", "", -1)
	***REMOVED***

	f, err := strconv.ParseFloat(num, 64)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	extra := strings.ToLower(strings.TrimSpace(s[lastDigit:]))
	if m, ok := bytesSizeTable[extra]; ok ***REMOVED***
		f *= float64(m)
		if f >= math.MaxUint64 ***REMOVED***
			return 0, fmt.Errorf("too large: %v", s)
		***REMOVED***
		return uint64(f), nil
	***REMOVED***

	return 0, fmt.Errorf("unhandled size name: %v", extra)
***REMOVED***
