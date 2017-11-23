package humanize

import (
	"math/big"
)

// order of magnitude (to a max order)
func oomm(n, b *big.Int, maxmag int) (float64, int) ***REMOVED***
	mag := 0
	m := &big.Int***REMOVED******REMOVED***
	for n.Cmp(b) >= 0 ***REMOVED***
		n.DivMod(n, b, m)
		mag++
		if mag == maxmag && maxmag >= 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return float64(n.Int64()) + (float64(m.Int64()) / float64(b.Int64())), mag
***REMOVED***

// total order of magnitude
// (same as above, but with no upper limit)
func oom(n, b *big.Int) (float64, int) ***REMOVED***
	mag := 0
	m := &big.Int***REMOVED******REMOVED***
	for n.Cmp(b) >= 0 ***REMOVED***
		n.DivMod(n, b, m)
		mag++
	***REMOVED***
	return float64(n.Int64()) + (float64(m.Int64()) / float64(b.Int64())), mag
***REMOVED***
