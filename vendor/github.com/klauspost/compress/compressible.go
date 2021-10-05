package compress

import "math"

// Estimate returns a normalized compressibility estimate of block b.
// Values close to zero are likely uncompressible.
// Values above 0.1 are likely to be compressible.
// Values above 0.5 are very compressible.
// Very small lengths will return 0.
func Estimate(b []byte) float64 ***REMOVED***
	if len(b) < 16 ***REMOVED***
		return 0
	***REMOVED***

	// Correctly predicted order 1
	hits := 0
	lastMatch := false
	var o1 [256]byte
	var hist [256]int
	c1 := byte(0)
	for _, c := range b ***REMOVED***
		if c == o1[c1] ***REMOVED***
			// We only count a hit if there was two correct predictions in a row.
			if lastMatch ***REMOVED***
				hits++
			***REMOVED***
			lastMatch = true
		***REMOVED*** else ***REMOVED***
			lastMatch = false
		***REMOVED***
		o1[c1] = c
		c1 = c
		hist[c]++
	***REMOVED***

	// Use x^0.6 to give better spread
	prediction := math.Pow(float64(hits)/float64(len(b)), 0.6)

	// Calculate histogram distribution
	variance := float64(0)
	avg := float64(len(b)) / 256

	for _, v := range hist ***REMOVED***
		Δ := float64(v) - avg
		variance += Δ * Δ
	***REMOVED***

	stddev := math.Sqrt(float64(variance)) / float64(len(b))
	exp := math.Sqrt(1 / float64(len(b)))

	// Subtract expected stddev
	stddev -= exp
	if stddev < 0 ***REMOVED***
		stddev = 0
	***REMOVED***
	stddev *= 1 + exp

	// Use x^0.4 to give better spread
	entropy := math.Pow(stddev, 0.4)

	// 50/50 weight between prediction and histogram distribution
	return math.Pow((prediction+entropy)/2, 0.9)
***REMOVED***

// ShannonEntropyBits returns the number of bits minimum required to represent
// an entropy encoding of the input bytes.
// https://en.wiktionary.org/wiki/Shannon_entropy
func ShannonEntropyBits(b []byte) int ***REMOVED***
	if len(b) == 0 ***REMOVED***
		return 0
	***REMOVED***
	var hist [256]int
	for _, c := range b ***REMOVED***
		hist[c]++
	***REMOVED***
	shannon := float64(0)
	invTotal := 1.0 / float64(len(b))
	for _, v := range hist[:] ***REMOVED***
		if v > 0 ***REMOVED***
			n := float64(v)
			shannon += math.Ceil(-math.Log2(n*invTotal) * n)
		***REMOVED***
	***REMOVED***
	return int(math.Ceil(shannon))
***REMOVED***
