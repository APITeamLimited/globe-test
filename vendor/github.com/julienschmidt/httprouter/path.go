// Copyright 2013 Julien Schmidt. All rights reserved.
// Based on the path package, Copyright 2009 The Go Authors.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package httprouter

// CleanPath is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//	1. Replace multiple slashes with a single slash.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned
func CleanPath(p string) string ***REMOVED***
	// Turn empty string into "/"
	if p == "" ***REMOVED***
		return "/"
	***REMOVED***

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' ***REMOVED***
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	***REMOVED***

	trailing := n > 1 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n ***REMOVED***
		switch ***REMOVED***
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r += 2

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 3

			if w > 1 ***REMOVED***
				// can backtrack
				w--

				if buf == nil ***REMOVED***
					for w > 1 && p[w] != '/' ***REMOVED***
						w--
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					for w > 1 && buf[w] != '/' ***REMOVED***
						w--
					***REMOVED***
				***REMOVED***
			***REMOVED***

		default:
			// real path element.
			// add slash if needed
			if w > 1 ***REMOVED***
				bufApp(&buf, p, w, '/')
				w++
			***REMOVED***

			// copy element
			for r < n && p[r] != '/' ***REMOVED***
				bufApp(&buf, p, w, p[r])
				w++
				r++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// re-append trailing slash
	if trailing && w > 1 ***REMOVED***
		bufApp(&buf, p, w, '/')
		w++
	***REMOVED***

	if buf == nil ***REMOVED***
		return p[:w]
	***REMOVED***
	return string(buf[:w])
***REMOVED***

// internal helper to lazily create a buffer if necessary
func bufApp(buf *[]byte, s string, w int, c byte) ***REMOVED***
	if *buf == nil ***REMOVED***
		if s[w] == c ***REMOVED***
			return
		***REMOVED***

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	***REMOVED***
	(*buf)[w] = c
***REMOVED***
