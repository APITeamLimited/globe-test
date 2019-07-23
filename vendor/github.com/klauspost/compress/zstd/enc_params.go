// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

type encParams struct ***REMOVED***
	// largest match distance : larger == more compression, more memory needed during decompression
	windowLog uint8

	// fully searched segment : larger == more compression, slower, more memory (useless for fast)
	chainLog uint8

	//  dispatch table : larger == faster, more memory
	hashLog uint8

	// < nb of searches : larger == more compression, slower
	searchLog uint8

	// < match length searched : larger == faster decompression, sometimes less compression
	minMatch uint8

	// acceptable match size for optimal parser (only) : larger == more compression, slower
	targetLength uint32

	// see ZSTD_strategy definition above
	strategy strategy
***REMOVED***

// strategy defines the algorithm to use when generating sequences.
type strategy uint8

const (
	// Compression strategies, listed from fastest to strongest
	strategyFast strategy = iota + 1
	strategyDfast
	strategyGreedy
	strategyLazy
	strategyLazy2
	strategyBtlazy2
	strategyBtopt
	strategyBtultra
	strategyBtultra2
	// note : new strategies _might_ be added in the future.
	//   Only the order (from fast to strong) is guaranteed

)

var defEncParams = [4][]encParams***REMOVED***
	***REMOVED*** // "default" - for any srcSize > 256 KB
		// W,  C,  H,  S,  L, TL, strat
		***REMOVED***19, 12, 13, 1, 6, 1, strategyFast***REMOVED***,       // base for negative levels
		***REMOVED***19, 13, 14, 1, 7, 0, strategyFast***REMOVED***,       // level  1
		***REMOVED***20, 15, 16, 1, 6, 0, strategyFast***REMOVED***,       // level  2
		***REMOVED***21, 16, 17, 1, 5, 1, strategyDfast***REMOVED***,      // level  3
		***REMOVED***21, 18, 18, 1, 5, 1, strategyDfast***REMOVED***,      // level  4
		***REMOVED***21, 18, 19, 2, 5, 2, strategyGreedy***REMOVED***,     // level  5
		***REMOVED***21, 19, 19, 3, 5, 4, strategyGreedy***REMOVED***,     // level  6
		***REMOVED***21, 19, 19, 3, 5, 8, strategyLazy***REMOVED***,       // level  7
		***REMOVED***21, 19, 19, 3, 5, 16, strategyLazy2***REMOVED***,     // level  8
		***REMOVED***21, 19, 20, 4, 5, 16, strategyLazy2***REMOVED***,     // level  9
		***REMOVED***22, 20, 21, 4, 5, 16, strategyLazy2***REMOVED***,     // level 10
		***REMOVED***22, 21, 22, 4, 5, 16, strategyLazy2***REMOVED***,     // level 11
		***REMOVED***22, 21, 22, 5, 5, 16, strategyLazy2***REMOVED***,     // level 12
		***REMOVED***22, 21, 22, 5, 5, 32, strategyBtlazy2***REMOVED***,   // level 13
		***REMOVED***22, 22, 23, 5, 5, 32, strategyBtlazy2***REMOVED***,   // level 14
		***REMOVED***22, 23, 23, 6, 5, 32, strategyBtlazy2***REMOVED***,   // level 15
		***REMOVED***22, 22, 22, 5, 5, 48, strategyBtopt***REMOVED***,     // level 16
		***REMOVED***23, 23, 22, 5, 4, 64, strategyBtopt***REMOVED***,     // level 17
		***REMOVED***23, 23, 22, 6, 3, 64, strategyBtultra***REMOVED***,   // level 18
		***REMOVED***23, 24, 22, 7, 3, 256, strategyBtultra2***REMOVED***, // level 19
		***REMOVED***25, 25, 23, 7, 3, 256, strategyBtultra2***REMOVED***, // level 20
		***REMOVED***26, 26, 24, 7, 3, 512, strategyBtultra2***REMOVED***, // level 21
		***REMOVED***27, 27, 25, 9, 3, 999, strategyBtultra2***REMOVED***, // level 22
	***REMOVED***,
	***REMOVED*** // for srcSize <= 256 KB
		// W,  C,  H,  S,  L,  T, strat
		***REMOVED***18, 12, 13, 1, 5, 1, strategyFast***REMOVED***,        // base for negative levels
		***REMOVED***18, 13, 14, 1, 6, 0, strategyFast***REMOVED***,        // level  1
		***REMOVED***18, 14, 14, 1, 5, 1, strategyDfast***REMOVED***,       // level  2
		***REMOVED***18, 16, 16, 1, 4, 1, strategyDfast***REMOVED***,       // level  3
		***REMOVED***18, 16, 17, 2, 5, 2, strategyGreedy***REMOVED***,      // level  4.
		***REMOVED***18, 18, 18, 3, 5, 2, strategyGreedy***REMOVED***,      // level  5.
		***REMOVED***18, 18, 19, 3, 5, 4, strategyLazy***REMOVED***,        // level  6.
		***REMOVED***18, 18, 19, 4, 4, 4, strategyLazy***REMOVED***,        // level  7
		***REMOVED***18, 18, 19, 4, 4, 8, strategyLazy2***REMOVED***,       // level  8
		***REMOVED***18, 18, 19, 5, 4, 8, strategyLazy2***REMOVED***,       // level  9
		***REMOVED***18, 18, 19, 6, 4, 8, strategyLazy2***REMOVED***,       // level 10
		***REMOVED***18, 18, 19, 5, 4, 12, strategyBtlazy2***REMOVED***,    // level 11.
		***REMOVED***18, 19, 19, 7, 4, 12, strategyBtlazy2***REMOVED***,    // level 12.
		***REMOVED***18, 18, 19, 4, 4, 16, strategyBtopt***REMOVED***,      // level 13
		***REMOVED***18, 18, 19, 4, 3, 32, strategyBtopt***REMOVED***,      // level 14.
		***REMOVED***18, 18, 19, 6, 3, 128, strategyBtopt***REMOVED***,     // level 15.
		***REMOVED***18, 19, 19, 6, 3, 128, strategyBtultra***REMOVED***,   // level 16.
		***REMOVED***18, 19, 19, 8, 3, 256, strategyBtultra***REMOVED***,   // level 17.
		***REMOVED***18, 19, 19, 6, 3, 128, strategyBtultra2***REMOVED***,  // level 18.
		***REMOVED***18, 19, 19, 8, 3, 256, strategyBtultra2***REMOVED***,  // level 19.
		***REMOVED***18, 19, 19, 10, 3, 512, strategyBtultra2***REMOVED***, // level 20.
		***REMOVED***18, 19, 19, 12, 3, 512, strategyBtultra2***REMOVED***, // level 21.
		***REMOVED***18, 19, 19, 13, 3, 999, strategyBtultra2***REMOVED***, // level 22.
	***REMOVED***,
	***REMOVED*** // for srcSize <= 128 KB
		// W,  C,  H,  S,  L,  T, strat
		***REMOVED***17, 12, 12, 1, 5, 1, strategyFast***REMOVED***,        // base for negative levels
		***REMOVED***17, 12, 13, 1, 6, 0, strategyFast***REMOVED***,        // level  1
		***REMOVED***17, 13, 15, 1, 5, 0, strategyFast***REMOVED***,        // level  2
		***REMOVED***17, 15, 16, 2, 5, 1, strategyDfast***REMOVED***,       // level  3
		***REMOVED***17, 17, 17, 2, 4, 1, strategyDfast***REMOVED***,       // level  4
		***REMOVED***17, 16, 17, 3, 4, 2, strategyGreedy***REMOVED***,      // level  5
		***REMOVED***17, 17, 17, 3, 4, 4, strategyLazy***REMOVED***,        // level  6
		***REMOVED***17, 17, 17, 3, 4, 8, strategyLazy2***REMOVED***,       // level  7
		***REMOVED***17, 17, 17, 4, 4, 8, strategyLazy2***REMOVED***,       // level  8
		***REMOVED***17, 17, 17, 5, 4, 8, strategyLazy2***REMOVED***,       // level  9
		***REMOVED***17, 17, 17, 6, 4, 8, strategyLazy2***REMOVED***,       // level 10
		***REMOVED***17, 17, 17, 5, 4, 8, strategyBtlazy2***REMOVED***,     // level 11
		***REMOVED***17, 18, 17, 7, 4, 12, strategyBtlazy2***REMOVED***,    // level 12
		***REMOVED***17, 18, 17, 3, 4, 12, strategyBtopt***REMOVED***,      // level 13.
		***REMOVED***17, 18, 17, 4, 3, 32, strategyBtopt***REMOVED***,      // level 14.
		***REMOVED***17, 18, 17, 6, 3, 256, strategyBtopt***REMOVED***,     // level 15.
		***REMOVED***17, 18, 17, 6, 3, 128, strategyBtultra***REMOVED***,   // level 16.
		***REMOVED***17, 18, 17, 8, 3, 256, strategyBtultra***REMOVED***,   // level 17.
		***REMOVED***17, 18, 17, 10, 3, 512, strategyBtultra***REMOVED***,  // level 18.
		***REMOVED***17, 18, 17, 5, 3, 256, strategyBtultra2***REMOVED***,  // level 19.
		***REMOVED***17, 18, 17, 7, 3, 512, strategyBtultra2***REMOVED***,  // level 20.
		***REMOVED***17, 18, 17, 9, 3, 512, strategyBtultra2***REMOVED***,  // level 21.
		***REMOVED***17, 18, 17, 11, 3, 999, strategyBtultra2***REMOVED***, // level 22.
	***REMOVED***,
	***REMOVED*** // for srcSize <= 16 KB
		// W,  C,  H,  S,  L,  T, strat
		***REMOVED***14, 12, 13, 1, 5, 1, strategyFast***REMOVED***,        // base for negative levels
		***REMOVED***14, 14, 15, 1, 5, 0, strategyFast***REMOVED***,        // level  1
		***REMOVED***14, 14, 15, 1, 4, 0, strategyFast***REMOVED***,        // level  2
		***REMOVED***14, 14, 15, 2, 4, 1, strategyDfast***REMOVED***,       // level  3
		***REMOVED***14, 14, 14, 4, 4, 2, strategyGreedy***REMOVED***,      // level  4
		***REMOVED***14, 14, 14, 3, 4, 4, strategyLazy***REMOVED***,        // level  5.
		***REMOVED***14, 14, 14, 4, 4, 8, strategyLazy2***REMOVED***,       // level  6
		***REMOVED***14, 14, 14, 6, 4, 8, strategyLazy2***REMOVED***,       // level  7
		***REMOVED***14, 14, 14, 8, 4, 8, strategyLazy2***REMOVED***,       // level  8.
		***REMOVED***14, 15, 14, 5, 4, 8, strategyBtlazy2***REMOVED***,     // level  9.
		***REMOVED***14, 15, 14, 9, 4, 8, strategyBtlazy2***REMOVED***,     // level 10.
		***REMOVED***14, 15, 14, 3, 4, 12, strategyBtopt***REMOVED***,      // level 11.
		***REMOVED***14, 15, 14, 4, 3, 24, strategyBtopt***REMOVED***,      // level 12.
		***REMOVED***14, 15, 14, 5, 3, 32, strategyBtultra***REMOVED***,    // level 13.
		***REMOVED***14, 15, 15, 6, 3, 64, strategyBtultra***REMOVED***,    // level 14.
		***REMOVED***14, 15, 15, 7, 3, 256, strategyBtultra***REMOVED***,   // level 15.
		***REMOVED***14, 15, 15, 5, 3, 48, strategyBtultra2***REMOVED***,   // level 16.
		***REMOVED***14, 15, 15, 6, 3, 128, strategyBtultra2***REMOVED***,  // level 17.
		***REMOVED***14, 15, 15, 7, 3, 256, strategyBtultra2***REMOVED***,  // level 18.
		***REMOVED***14, 15, 15, 8, 3, 256, strategyBtultra2***REMOVED***,  // level 19.
		***REMOVED***14, 15, 15, 8, 3, 512, strategyBtultra2***REMOVED***,  // level 20.
		***REMOVED***14, 15, 15, 9, 3, 512, strategyBtultra2***REMOVED***,  // level 21.
		***REMOVED***14, 15, 15, 10, 3, 999, strategyBtultra2***REMOVED***, // level 22.
	***REMOVED***,
***REMOVED***
