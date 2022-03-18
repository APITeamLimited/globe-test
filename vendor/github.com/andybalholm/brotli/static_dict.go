package brotli

import "encoding/binary"

/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Class to model the static dictionary. */

const maxStaticDictionaryMatchLen = 37

const kInvalidMatch uint32 = 0xFFFFFFF

/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/
func hash(data []byte) uint32 ***REMOVED***
	var h uint32 = binary.LittleEndian.Uint32(data) * kDictHashMul32

	/* The higher bits contain more mixture from the multiplication,
	   so we take our results from there. */
	return h >> uint(32-kDictNumBits)
***REMOVED***

func addMatch(distance uint, len uint, len_code uint, matches []uint32) ***REMOVED***
	var match uint32 = uint32((distance << 5) + len_code)
	matches[len] = brotli_min_uint32_t(matches[len], match)
***REMOVED***

func dictMatchLength(dict *dictionary, data []byte, id uint, len uint, maxlen uint) uint ***REMOVED***
	var offset uint = uint(dict.offsets_by_length[len]) + len*id
	return findMatchLengthWithLimit(dict.data[offset:], data, brotli_min_size_t(uint(len), maxlen))
***REMOVED***

func isMatch(d *dictionary, w dictWord, data []byte, max_length uint) bool ***REMOVED***
	if uint(w.len) > max_length ***REMOVED***
		return false
	***REMOVED*** else ***REMOVED***
		var offset uint = uint(d.offsets_by_length[w.len]) + uint(w.len)*uint(w.idx)
		var dict []byte = d.data[offset:]
		if w.transform == 0 ***REMOVED***
			/* Match against base dictionary word. */
			return findMatchLengthWithLimit(dict, data, uint(w.len)) == uint(w.len)
		***REMOVED*** else if w.transform == 10 ***REMOVED***
			/* Match against uppercase first transform.
			   Note that there are only ASCII uppercase words in the lookup table. */
			return dict[0] >= 'a' && dict[0] <= 'z' && (dict[0]^32) == data[0] && findMatchLengthWithLimit(dict[1:], data[1:], uint(w.len)-1) == uint(w.len-1)
		***REMOVED*** else ***REMOVED***
			/* Match against uppercase all transform.
			   Note that there are only ASCII uppercase words in the lookup table. */
			var i uint
			for i = 0; i < uint(w.len); i++ ***REMOVED***
				if dict[i] >= 'a' && dict[i] <= 'z' ***REMOVED***
					if (dict[i] ^ 32) != data[i] ***REMOVED***
						return false
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if dict[i] != data[i] ***REMOVED***
						return false
					***REMOVED***
				***REMOVED***
			***REMOVED***

			return true
		***REMOVED***
	***REMOVED***
***REMOVED***

func findAllStaticDictionaryMatches(dict *encoderDictionary, data []byte, min_length uint, max_length uint, matches []uint32) bool ***REMOVED***
	var has_found_match bool = false
	***REMOVED***
		var offset uint = uint(dict.buckets[hash(data)])
		var end bool = offset == 0
		for !end ***REMOVED***
			w := dict.dict_words[offset]
			offset++
			var l uint = uint(w.len) & 0x1F
			var n uint = uint(1) << dict.words.size_bits_by_length[l]
			var id uint = uint(w.idx)
			end = !(w.len&0x80 == 0)
			w.len = byte(l)
			if w.transform == 0 ***REMOVED***
				var matchlen uint = dictMatchLength(dict.words, data, id, l, max_length)
				var s []byte
				var minlen uint
				var maxlen uint
				var len uint

				/* Transform "" + BROTLI_TRANSFORM_IDENTITY + "" */
				if matchlen == l ***REMOVED***
					addMatch(id, l, l, matches)
					has_found_match = true
				***REMOVED***

				/* Transforms "" + BROTLI_TRANSFORM_OMIT_LAST_1 + "" and
				   "" + BROTLI_TRANSFORM_OMIT_LAST_1 + "ing " */
				if matchlen >= l-1 ***REMOVED***
					addMatch(id+12*n, l-1, l, matches)
					if l+2 < max_length && data[l-1] == 'i' && data[l] == 'n' && data[l+1] == 'g' && data[l+2] == ' ' ***REMOVED***
						addMatch(id+49*n, l+3, l, matches)
					***REMOVED***

					has_found_match = true
				***REMOVED***

				/* Transform "" + BROTLI_TRANSFORM_OMIT_LAST_# + "" (# = 2 .. 9) */
				minlen = min_length

				if l > 9 ***REMOVED***
					minlen = brotli_max_size_t(minlen, l-9)
				***REMOVED***
				maxlen = brotli_min_size_t(matchlen, l-2)
				for len = minlen; len <= maxlen; len++ ***REMOVED***
					var cut uint = l - len
					var transform_id uint = (cut << 2) + uint((dict.cutoffTransforms>>(cut*6))&0x3F)
					addMatch(id+transform_id*n, uint(len), l, matches)
					has_found_match = true
				***REMOVED***

				if matchlen < l || l+6 >= max_length ***REMOVED***
					continue
				***REMOVED***

				s = data[l:]

				/* Transforms "" + BROTLI_TRANSFORM_IDENTITY + <suffix> */
				if s[0] == ' ' ***REMOVED***
					addMatch(id+n, l+1, l, matches)
					if s[1] == 'a' ***REMOVED***
						if s[2] == ' ' ***REMOVED***
							addMatch(id+28*n, l+3, l, matches)
						***REMOVED*** else if s[2] == 's' ***REMOVED***
							if s[3] == ' ' ***REMOVED***
								addMatch(id+46*n, l+4, l, matches)
							***REMOVED***
						***REMOVED*** else if s[2] == 't' ***REMOVED***
							if s[3] == ' ' ***REMOVED***
								addMatch(id+60*n, l+4, l, matches)
							***REMOVED***
						***REMOVED*** else if s[2] == 'n' ***REMOVED***
							if s[3] == 'd' && s[4] == ' ' ***REMOVED***
								addMatch(id+10*n, l+5, l, matches)
							***REMOVED***
						***REMOVED***
					***REMOVED*** else if s[1] == 'b' ***REMOVED***
						if s[2] == 'y' && s[3] == ' ' ***REMOVED***
							addMatch(id+38*n, l+4, l, matches)
						***REMOVED***
					***REMOVED*** else if s[1] == 'i' ***REMOVED***
						if s[2] == 'n' ***REMOVED***
							if s[3] == ' ' ***REMOVED***
								addMatch(id+16*n, l+4, l, matches)
							***REMOVED***
						***REMOVED*** else if s[2] == 's' ***REMOVED***
							if s[3] == ' ' ***REMOVED***
								addMatch(id+47*n, l+4, l, matches)
							***REMOVED***
						***REMOVED***
					***REMOVED*** else if s[1] == 'f' ***REMOVED***
						if s[2] == 'o' ***REMOVED***
							if s[3] == 'r' && s[4] == ' ' ***REMOVED***
								addMatch(id+25*n, l+5, l, matches)
							***REMOVED***
						***REMOVED*** else if s[2] == 'r' ***REMOVED***
							if s[3] == 'o' && s[4] == 'm' && s[5] == ' ' ***REMOVED***
								addMatch(id+37*n, l+6, l, matches)
							***REMOVED***
						***REMOVED***
					***REMOVED*** else if s[1] == 'o' ***REMOVED***
						if s[2] == 'f' ***REMOVED***
							if s[3] == ' ' ***REMOVED***
								addMatch(id+8*n, l+4, l, matches)
							***REMOVED***
						***REMOVED*** else if s[2] == 'n' ***REMOVED***
							if s[3] == ' ' ***REMOVED***
								addMatch(id+45*n, l+4, l, matches)
							***REMOVED***
						***REMOVED***
					***REMOVED*** else if s[1] == 'n' ***REMOVED***
						if s[2] == 'o' && s[3] == 't' && s[4] == ' ' ***REMOVED***
							addMatch(id+80*n, l+5, l, matches)
						***REMOVED***
					***REMOVED*** else if s[1] == 't' ***REMOVED***
						if s[2] == 'h' ***REMOVED***
							if s[3] == 'e' ***REMOVED***
								if s[4] == ' ' ***REMOVED***
									addMatch(id+5*n, l+5, l, matches)
								***REMOVED***
							***REMOVED*** else if s[3] == 'a' ***REMOVED***
								if s[4] == 't' && s[5] == ' ' ***REMOVED***
									addMatch(id+29*n, l+6, l, matches)
								***REMOVED***
							***REMOVED***
						***REMOVED*** else if s[2] == 'o' ***REMOVED***
							if s[3] == ' ' ***REMOVED***
								addMatch(id+17*n, l+4, l, matches)
							***REMOVED***
						***REMOVED***
					***REMOVED*** else if s[1] == 'w' ***REMOVED***
						if s[2] == 'i' && s[3] == 't' && s[4] == 'h' && s[5] == ' ' ***REMOVED***
							addMatch(id+35*n, l+6, l, matches)
						***REMOVED***
					***REMOVED***
				***REMOVED*** else if s[0] == '"' ***REMOVED***
					addMatch(id+19*n, l+1, l, matches)
					if s[1] == '>' ***REMOVED***
						addMatch(id+21*n, l+2, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == '.' ***REMOVED***
					addMatch(id+20*n, l+1, l, matches)
					if s[1] == ' ' ***REMOVED***
						addMatch(id+31*n, l+2, l, matches)
						if s[2] == 'T' && s[3] == 'h' ***REMOVED***
							if s[4] == 'e' ***REMOVED***
								if s[5] == ' ' ***REMOVED***
									addMatch(id+43*n, l+6, l, matches)
								***REMOVED***
							***REMOVED*** else if s[4] == 'i' ***REMOVED***
								if s[5] == 's' && s[6] == ' ' ***REMOVED***
									addMatch(id+75*n, l+7, l, matches)
								***REMOVED***
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED*** else if s[0] == ',' ***REMOVED***
					addMatch(id+76*n, l+1, l, matches)
					if s[1] == ' ' ***REMOVED***
						addMatch(id+14*n, l+2, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == '\n' ***REMOVED***
					addMatch(id+22*n, l+1, l, matches)
					if s[1] == '\t' ***REMOVED***
						addMatch(id+50*n, l+2, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == ']' ***REMOVED***
					addMatch(id+24*n, l+1, l, matches)
				***REMOVED*** else if s[0] == '\'' ***REMOVED***
					addMatch(id+36*n, l+1, l, matches)
				***REMOVED*** else if s[0] == ':' ***REMOVED***
					addMatch(id+51*n, l+1, l, matches)
				***REMOVED*** else if s[0] == '(' ***REMOVED***
					addMatch(id+57*n, l+1, l, matches)
				***REMOVED*** else if s[0] == '=' ***REMOVED***
					if s[1] == '"' ***REMOVED***
						addMatch(id+70*n, l+2, l, matches)
					***REMOVED*** else if s[1] == '\'' ***REMOVED***
						addMatch(id+86*n, l+2, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == 'a' ***REMOVED***
					if s[1] == 'l' && s[2] == ' ' ***REMOVED***
						addMatch(id+84*n, l+3, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == 'e' ***REMOVED***
					if s[1] == 'd' ***REMOVED***
						if s[2] == ' ' ***REMOVED***
							addMatch(id+53*n, l+3, l, matches)
						***REMOVED***
					***REMOVED*** else if s[1] == 'r' ***REMOVED***
						if s[2] == ' ' ***REMOVED***
							addMatch(id+82*n, l+3, l, matches)
						***REMOVED***
					***REMOVED*** else if s[1] == 's' ***REMOVED***
						if s[2] == 't' && s[3] == ' ' ***REMOVED***
							addMatch(id+95*n, l+4, l, matches)
						***REMOVED***
					***REMOVED***
				***REMOVED*** else if s[0] == 'f' ***REMOVED***
					if s[1] == 'u' && s[2] == 'l' && s[3] == ' ' ***REMOVED***
						addMatch(id+90*n, l+4, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == 'i' ***REMOVED***
					if s[1] == 'v' ***REMOVED***
						if s[2] == 'e' && s[3] == ' ' ***REMOVED***
							addMatch(id+92*n, l+4, l, matches)
						***REMOVED***
					***REMOVED*** else if s[1] == 'z' ***REMOVED***
						if s[2] == 'e' && s[3] == ' ' ***REMOVED***
							addMatch(id+100*n, l+4, l, matches)
						***REMOVED***
					***REMOVED***
				***REMOVED*** else if s[0] == 'l' ***REMOVED***
					if s[1] == 'e' ***REMOVED***
						if s[2] == 's' && s[3] == 's' && s[4] == ' ' ***REMOVED***
							addMatch(id+93*n, l+5, l, matches)
						***REMOVED***
					***REMOVED*** else if s[1] == 'y' ***REMOVED***
						if s[2] == ' ' ***REMOVED***
							addMatch(id+61*n, l+3, l, matches)
						***REMOVED***
					***REMOVED***
				***REMOVED*** else if s[0] == 'o' ***REMOVED***
					if s[1] == 'u' && s[2] == 's' && s[3] == ' ' ***REMOVED***
						addMatch(id+106*n, l+4, l, matches)
					***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				var is_all_caps bool = (w.transform != transformUppercaseFirst)
				/* Set is_all_caps=0 for BROTLI_TRANSFORM_UPPERCASE_FIRST and
				    is_all_caps=1 otherwise (BROTLI_TRANSFORM_UPPERCASE_ALL)
				transform. */

				var s []byte
				if !isMatch(dict.words, w, data, max_length) ***REMOVED***
					continue
				***REMOVED***

				/* Transform "" + kUppercase***REMOVED***First,All***REMOVED*** + "" */
				var tmp int
				if is_all_caps ***REMOVED***
					tmp = 44
				***REMOVED*** else ***REMOVED***
					tmp = 9
				***REMOVED***
				addMatch(id+uint(tmp)*n, l, l, matches)

				has_found_match = true
				if l+1 >= max_length ***REMOVED***
					continue
				***REMOVED***

				/* Transforms "" + kUppercase***REMOVED***First,All***REMOVED*** + <suffix> */
				s = data[l:]

				if s[0] == ' ' ***REMOVED***
					var tmp int
					if is_all_caps ***REMOVED***
						tmp = 68
					***REMOVED*** else ***REMOVED***
						tmp = 4
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+1, l, matches)
				***REMOVED*** else if s[0] == '"' ***REMOVED***
					var tmp int
					if is_all_caps ***REMOVED***
						tmp = 87
					***REMOVED*** else ***REMOVED***
						tmp = 66
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+1, l, matches)
					if s[1] == '>' ***REMOVED***
						var tmp int
						if is_all_caps ***REMOVED***
							tmp = 97
						***REMOVED*** else ***REMOVED***
							tmp = 69
						***REMOVED***
						addMatch(id+uint(tmp)*n, l+2, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == '.' ***REMOVED***
					var tmp int
					if is_all_caps ***REMOVED***
						tmp = 101
					***REMOVED*** else ***REMOVED***
						tmp = 79
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+1, l, matches)
					if s[1] == ' ' ***REMOVED***
						var tmp int
						if is_all_caps ***REMOVED***
							tmp = 114
						***REMOVED*** else ***REMOVED***
							tmp = 88
						***REMOVED***
						addMatch(id+uint(tmp)*n, l+2, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == ',' ***REMOVED***
					var tmp int
					if is_all_caps ***REMOVED***
						tmp = 112
					***REMOVED*** else ***REMOVED***
						tmp = 99
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+1, l, matches)
					if s[1] == ' ' ***REMOVED***
						var tmp int
						if is_all_caps ***REMOVED***
							tmp = 107
						***REMOVED*** else ***REMOVED***
							tmp = 58
						***REMOVED***
						addMatch(id+uint(tmp)*n, l+2, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == '\'' ***REMOVED***
					var tmp int
					if is_all_caps ***REMOVED***
						tmp = 94
					***REMOVED*** else ***REMOVED***
						tmp = 74
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+1, l, matches)
				***REMOVED*** else if s[0] == '(' ***REMOVED***
					var tmp int
					if is_all_caps ***REMOVED***
						tmp = 113
					***REMOVED*** else ***REMOVED***
						tmp = 78
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+1, l, matches)
				***REMOVED*** else if s[0] == '=' ***REMOVED***
					if s[1] == '"' ***REMOVED***
						var tmp int
						if is_all_caps ***REMOVED***
							tmp = 105
						***REMOVED*** else ***REMOVED***
							tmp = 104
						***REMOVED***
						addMatch(id+uint(tmp)*n, l+2, l, matches)
					***REMOVED*** else if s[1] == '\'' ***REMOVED***
						var tmp int
						if is_all_caps ***REMOVED***
							tmp = 116
						***REMOVED*** else ***REMOVED***
							tmp = 108
						***REMOVED***
						addMatch(id+uint(tmp)*n, l+2, l, matches)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	/* Transforms with prefixes " " and "." */
	if max_length >= 5 && (data[0] == ' ' || data[0] == '.') ***REMOVED***
		var is_space bool = (data[0] == ' ')
		var offset uint = uint(dict.buckets[hash(data[1:])])
		var end bool = offset == 0
		for !end ***REMOVED***
			w := dict.dict_words[offset]
			offset++
			var l uint = uint(w.len) & 0x1F
			var n uint = uint(1) << dict.words.size_bits_by_length[l]
			var id uint = uint(w.idx)
			end = !(w.len&0x80 == 0)
			w.len = byte(l)
			if w.transform == 0 ***REMOVED***
				var s []byte
				if !isMatch(dict.words, w, data[1:], max_length-1) ***REMOVED***
					continue
				***REMOVED***

				/* Transforms " " + BROTLI_TRANSFORM_IDENTITY + "" and
				   "." + BROTLI_TRANSFORM_IDENTITY + "" */
				var tmp int
				if is_space ***REMOVED***
					tmp = 6
				***REMOVED*** else ***REMOVED***
					tmp = 32
				***REMOVED***
				addMatch(id+uint(tmp)*n, l+1, l, matches)

				has_found_match = true
				if l+2 >= max_length ***REMOVED***
					continue
				***REMOVED***

				/* Transforms " " + BROTLI_TRANSFORM_IDENTITY + <suffix> and
				   "." + BROTLI_TRANSFORM_IDENTITY + <suffix>
				*/
				s = data[l+1:]

				if s[0] == ' ' ***REMOVED***
					var tmp int
					if is_space ***REMOVED***
						tmp = 2
					***REMOVED*** else ***REMOVED***
						tmp = 77
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+2, l, matches)
				***REMOVED*** else if s[0] == '(' ***REMOVED***
					var tmp int
					if is_space ***REMOVED***
						tmp = 89
					***REMOVED*** else ***REMOVED***
						tmp = 67
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+2, l, matches)
				***REMOVED*** else if is_space ***REMOVED***
					if s[0] == ',' ***REMOVED***
						addMatch(id+103*n, l+2, l, matches)
						if s[1] == ' ' ***REMOVED***
							addMatch(id+33*n, l+3, l, matches)
						***REMOVED***
					***REMOVED*** else if s[0] == '.' ***REMOVED***
						addMatch(id+71*n, l+2, l, matches)
						if s[1] == ' ' ***REMOVED***
							addMatch(id+52*n, l+3, l, matches)
						***REMOVED***
					***REMOVED*** else if s[0] == '=' ***REMOVED***
						if s[1] == '"' ***REMOVED***
							addMatch(id+81*n, l+3, l, matches)
						***REMOVED*** else if s[1] == '\'' ***REMOVED***
							addMatch(id+98*n, l+3, l, matches)
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED*** else if is_space ***REMOVED***
				var is_all_caps bool = (w.transform != transformUppercaseFirst)
				/* Set is_all_caps=0 for BROTLI_TRANSFORM_UPPERCASE_FIRST and
				    is_all_caps=1 otherwise (BROTLI_TRANSFORM_UPPERCASE_ALL)
				transform. */

				var s []byte
				if !isMatch(dict.words, w, data[1:], max_length-1) ***REMOVED***
					continue
				***REMOVED***

				/* Transforms " " + kUppercase***REMOVED***First,All***REMOVED*** + "" */
				var tmp int
				if is_all_caps ***REMOVED***
					tmp = 85
				***REMOVED*** else ***REMOVED***
					tmp = 30
				***REMOVED***
				addMatch(id+uint(tmp)*n, l+1, l, matches)

				has_found_match = true
				if l+2 >= max_length ***REMOVED***
					continue
				***REMOVED***

				/* Transforms " " + kUppercase***REMOVED***First,All***REMOVED*** + <suffix> */
				s = data[l+1:]

				if s[0] == ' ' ***REMOVED***
					var tmp int
					if is_all_caps ***REMOVED***
						tmp = 83
					***REMOVED*** else ***REMOVED***
						tmp = 15
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+2, l, matches)
				***REMOVED*** else if s[0] == ',' ***REMOVED***
					if !is_all_caps ***REMOVED***
						addMatch(id+109*n, l+2, l, matches)
					***REMOVED***

					if s[1] == ' ' ***REMOVED***
						var tmp int
						if is_all_caps ***REMOVED***
							tmp = 111
						***REMOVED*** else ***REMOVED***
							tmp = 65
						***REMOVED***
						addMatch(id+uint(tmp)*n, l+3, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == '.' ***REMOVED***
					var tmp int
					if is_all_caps ***REMOVED***
						tmp = 115
					***REMOVED*** else ***REMOVED***
						tmp = 96
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+2, l, matches)
					if s[1] == ' ' ***REMOVED***
						var tmp int
						if is_all_caps ***REMOVED***
							tmp = 117
						***REMOVED*** else ***REMOVED***
							tmp = 91
						***REMOVED***
						addMatch(id+uint(tmp)*n, l+3, l, matches)
					***REMOVED***
				***REMOVED*** else if s[0] == '=' ***REMOVED***
					if s[1] == '"' ***REMOVED***
						var tmp int
						if is_all_caps ***REMOVED***
							tmp = 110
						***REMOVED*** else ***REMOVED***
							tmp = 118
						***REMOVED***
						addMatch(id+uint(tmp)*n, l+3, l, matches)
					***REMOVED*** else if s[1] == '\'' ***REMOVED***
						var tmp int
						if is_all_caps ***REMOVED***
							tmp = 119
						***REMOVED*** else ***REMOVED***
							tmp = 120
						***REMOVED***
						addMatch(id+uint(tmp)*n, l+3, l, matches)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if max_length >= 6 ***REMOVED***
		/* Transforms with prefixes "e ", "s ", ", " and "\xC2\xA0" */
		if (data[1] == ' ' && (data[0] == 'e' || data[0] == 's' || data[0] == ',')) || (data[0] == 0xC2 && data[1] == 0xA0) ***REMOVED***
			var offset uint = uint(dict.buckets[hash(data[2:])])
			var end bool = offset == 0
			for !end ***REMOVED***
				w := dict.dict_words[offset]
				offset++
				var l uint = uint(w.len) & 0x1F
				var n uint = uint(1) << dict.words.size_bits_by_length[l]
				var id uint = uint(w.idx)
				end = !(w.len&0x80 == 0)
				w.len = byte(l)
				if w.transform == 0 && isMatch(dict.words, w, data[2:], max_length-2) ***REMOVED***
					if data[0] == 0xC2 ***REMOVED***
						addMatch(id+102*n, l+2, l, matches)
						has_found_match = true
					***REMOVED*** else if l+2 < max_length && data[l+2] == ' ' ***REMOVED***
						var t uint = 13
						if data[0] == 'e' ***REMOVED***
							t = 18
						***REMOVED*** else if data[0] == 's' ***REMOVED***
							t = 7
						***REMOVED***
						addMatch(id+t*n, l+3, l, matches)
						has_found_match = true
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if max_length >= 9 ***REMOVED***
		/* Transforms with prefixes " the " and ".com/" */
		if (data[0] == ' ' && data[1] == 't' && data[2] == 'h' && data[3] == 'e' && data[4] == ' ') || (data[0] == '.' && data[1] == 'c' && data[2] == 'o' && data[3] == 'm' && data[4] == '/') ***REMOVED***
			var offset uint = uint(dict.buckets[hash(data[5:])])
			var end bool = offset == 0
			for !end ***REMOVED***
				w := dict.dict_words[offset]
				offset++
				var l uint = uint(w.len) & 0x1F
				var n uint = uint(1) << dict.words.size_bits_by_length[l]
				var id uint = uint(w.idx)
				end = !(w.len&0x80 == 0)
				w.len = byte(l)
				if w.transform == 0 && isMatch(dict.words, w, data[5:], max_length-5) ***REMOVED***
					var tmp int
					if data[0] == ' ' ***REMOVED***
						tmp = 41
					***REMOVED*** else ***REMOVED***
						tmp = 72
					***REMOVED***
					addMatch(id+uint(tmp)*n, l+5, l, matches)
					has_found_match = true
					if l+5 < max_length ***REMOVED***
						var s []byte = data[l+5:]
						if data[0] == ' ' ***REMOVED***
							if l+8 < max_length && s[0] == ' ' && s[1] == 'o' && s[2] == 'f' && s[3] == ' ' ***REMOVED***
								addMatch(id+62*n, l+9, l, matches)
								if l+12 < max_length && s[4] == 't' && s[5] == 'h' && s[6] == 'e' && s[7] == ' ' ***REMOVED***
									addMatch(id+73*n, l+13, l, matches)
								***REMOVED***
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return has_found_match
***REMOVED***
