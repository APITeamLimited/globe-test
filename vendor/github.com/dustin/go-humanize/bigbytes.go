package humanize

import (
	"fmt"
	"math/big"
	"strings"
	"unicode"
)

var (
	bigIECExp = big.NewInt(1024)

	// BigByte is one byte in bit.Ints
	BigByte = big.NewInt(1)
	// BigKiByte is 1,024 bytes in bit.Ints
	BigKiByte = (&big.Int***REMOVED******REMOVED***).Mul(BigByte, bigIECExp)
	// BigMiByte is 1,024 k bytes in bit.Ints
	BigMiByte = (&big.Int***REMOVED******REMOVED***).Mul(BigKiByte, bigIECExp)
	// BigGiByte is 1,024 m bytes in bit.Ints
	BigGiByte = (&big.Int***REMOVED******REMOVED***).Mul(BigMiByte, bigIECExp)
	// BigTiByte is 1,024 g bytes in bit.Ints
	BigTiByte = (&big.Int***REMOVED******REMOVED***).Mul(BigGiByte, bigIECExp)
	// BigPiByte is 1,024 t bytes in bit.Ints
	BigPiByte = (&big.Int***REMOVED******REMOVED***).Mul(BigTiByte, bigIECExp)
	// BigEiByte is 1,024 p bytes in bit.Ints
	BigEiByte = (&big.Int***REMOVED******REMOVED***).Mul(BigPiByte, bigIECExp)
	// BigZiByte is 1,024 e bytes in bit.Ints
	BigZiByte = (&big.Int***REMOVED******REMOVED***).Mul(BigEiByte, bigIECExp)
	// BigYiByte is 1,024 z bytes in bit.Ints
	BigYiByte = (&big.Int***REMOVED******REMOVED***).Mul(BigZiByte, bigIECExp)
)

var (
	bigSIExp = big.NewInt(1000)

	// BigSIByte is one SI byte in big.Ints
	BigSIByte = big.NewInt(1)
	// BigKByte is 1,000 SI bytes in big.Ints
	BigKByte = (&big.Int***REMOVED******REMOVED***).Mul(BigSIByte, bigSIExp)
	// BigMByte is 1,000 SI k bytes in big.Ints
	BigMByte = (&big.Int***REMOVED******REMOVED***).Mul(BigKByte, bigSIExp)
	// BigGByte is 1,000 SI m bytes in big.Ints
	BigGByte = (&big.Int***REMOVED******REMOVED***).Mul(BigMByte, bigSIExp)
	// BigTByte is 1,000 SI g bytes in big.Ints
	BigTByte = (&big.Int***REMOVED******REMOVED***).Mul(BigGByte, bigSIExp)
	// BigPByte is 1,000 SI t bytes in big.Ints
	BigPByte = (&big.Int***REMOVED******REMOVED***).Mul(BigTByte, bigSIExp)
	// BigEByte is 1,000 SI p bytes in big.Ints
	BigEByte = (&big.Int***REMOVED******REMOVED***).Mul(BigPByte, bigSIExp)
	// BigZByte is 1,000 SI e bytes in big.Ints
	BigZByte = (&big.Int***REMOVED******REMOVED***).Mul(BigEByte, bigSIExp)
	// BigYByte is 1,000 SI z bytes in big.Ints
	BigYByte = (&big.Int***REMOVED******REMOVED***).Mul(BigZByte, bigSIExp)
)

var bigBytesSizeTable = map[string]*big.Int***REMOVED***
	"b":   BigByte,
	"kib": BigKiByte,
	"kb":  BigKByte,
	"mib": BigMiByte,
	"mb":  BigMByte,
	"gib": BigGiByte,
	"gb":  BigGByte,
	"tib": BigTiByte,
	"tb":  BigTByte,
	"pib": BigPiByte,
	"pb":  BigPByte,
	"eib": BigEiByte,
	"eb":  BigEByte,
	"zib": BigZiByte,
	"zb":  BigZByte,
	"yib": BigYiByte,
	"yb":  BigYByte,
	// Without suffix
	"":   BigByte,
	"ki": BigKiByte,
	"k":  BigKByte,
	"mi": BigMiByte,
	"m":  BigMByte,
	"gi": BigGiByte,
	"g":  BigGByte,
	"ti": BigTiByte,
	"t":  BigTByte,
	"pi": BigPiByte,
	"p":  BigPByte,
	"ei": BigEiByte,
	"e":  BigEByte,
	"z":  BigZByte,
	"zi": BigZiByte,
	"y":  BigYByte,
	"yi": BigYiByte,
***REMOVED***

var ten = big.NewInt(10)

func humanateBigBytes(s, base *big.Int, sizes []string) string ***REMOVED***
	if s.Cmp(ten) < 0 ***REMOVED***
		return fmt.Sprintf("%d B", s)
	***REMOVED***
	c := (&big.Int***REMOVED******REMOVED***).Set(s)
	val, mag := oomm(c, base, len(sizes)-1)
	suffix := sizes[mag]
	f := "%.0f %s"
	if val < 10 ***REMOVED***
		f = "%.1f %s"
	***REMOVED***

	return fmt.Sprintf(f, val, suffix)

***REMOVED***

// BigBytes produces a human readable representation of an SI size.
//
// See also: ParseBigBytes.
//
// BigBytes(82854982) -> 83 MB
func BigBytes(s *big.Int) string ***REMOVED***
	sizes := []string***REMOVED***"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"***REMOVED***
	return humanateBigBytes(s, bigSIExp, sizes)
***REMOVED***

// BigIBytes produces a human readable representation of an IEC size.
//
// See also: ParseBigBytes.
//
// BigIBytes(82854982) -> 79 MiB
func BigIBytes(s *big.Int) string ***REMOVED***
	sizes := []string***REMOVED***"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"***REMOVED***
	return humanateBigBytes(s, bigIECExp, sizes)
***REMOVED***

// ParseBigBytes parses a string representation of bytes into the number
// of bytes it represents.
//
// See also: BigBytes, BigIBytes.
//
// ParseBigBytes("42 MB") -> 42000000, nil
// ParseBigBytes("42 mib") -> 44040192, nil
func ParseBigBytes(s string) (*big.Int, error) ***REMOVED***
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

	val := &big.Rat***REMOVED******REMOVED***
	_, err := fmt.Sscanf(num, "%f", val)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	extra := strings.ToLower(strings.TrimSpace(s[lastDigit:]))
	if m, ok := bigBytesSizeTable[extra]; ok ***REMOVED***
		mv := (&big.Rat***REMOVED******REMOVED***).SetInt(m)
		val.Mul(val, mv)
		rv := &big.Int***REMOVED******REMOVED***
		rv.Div(val.Num(), val.Denom())
		return rv, nil
	***REMOVED***

	return nil, fmt.Errorf("unhandled size name: %v", extra)
***REMOVED***
