// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package httprouter

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func min(a, b int) int ***REMOVED***
	if a <= b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

func countParams(path string) uint8 ***REMOVED***
	var n uint
	for i := 0; i < len(path); i++ ***REMOVED***
		if path[i] != ':' && path[i] != '*' ***REMOVED***
			continue
		***REMOVED***
		n++
	***REMOVED***
	if n >= 255 ***REMOVED***
		return 255
	***REMOVED***
	return uint8(n)
***REMOVED***

type nodeType uint8

const (
	static nodeType = iota // default
	root
	param
	catchAll
)

type node struct ***REMOVED***
	path      string
	wildChild bool
	nType     nodeType
	maxParams uint8
	indices   string
	children  []*node
	handle    Handle
	priority  uint32
***REMOVED***

// increments priority of the given child and reorders if necessary
func (n *node) incrementChildPrio(pos int) int ***REMOVED***
	n.children[pos].priority++
	prio := n.children[pos].priority

	// adjust position (move to front)
	newPos := pos
	for newPos > 0 && n.children[newPos-1].priority < prio ***REMOVED***
		// swap node positions
		n.children[newPos-1], n.children[newPos] = n.children[newPos], n.children[newPos-1]

		newPos--
	***REMOVED***

	// build new index char string
	if newPos != pos ***REMOVED***
		n.indices = n.indices[:newPos] + // unchanged prefix, might be empty
			n.indices[pos:pos+1] + // the index char we move
			n.indices[newPos:pos] + n.indices[pos+1:] // rest without char at 'pos'
	***REMOVED***

	return newPos
***REMOVED***

// addRoute adds a node with the given handle to the path.
// Not concurrency-safe!
func (n *node) addRoute(path string, handle Handle) ***REMOVED***
	fullPath := path
	n.priority++
	numParams := countParams(path)

	// non-empty tree
	if len(n.path) > 0 || len(n.children) > 0 ***REMOVED***
	walk:
		for ***REMOVED***
			// Update maxParams of the current node
			if numParams > n.maxParams ***REMOVED***
				n.maxParams = numParams
			***REMOVED***

			// Find the longest common prefix.
			// This also implies that the common prefix contains no ':' or '*'
			// since the existing key can't contain those chars.
			i := 0
			max := min(len(path), len(n.path))
			for i < max && path[i] == n.path[i] ***REMOVED***
				i++
			***REMOVED***

			// Split edge
			if i < len(n.path) ***REMOVED***
				child := node***REMOVED***
					path:      n.path[i:],
					wildChild: n.wildChild,
					nType:     static,
					indices:   n.indices,
					children:  n.children,
					handle:    n.handle,
					priority:  n.priority - 1,
				***REMOVED***

				// Update maxParams (max of all children)
				for i := range child.children ***REMOVED***
					if child.children[i].maxParams > child.maxParams ***REMOVED***
						child.maxParams = child.children[i].maxParams
					***REMOVED***
				***REMOVED***

				n.children = []*node***REMOVED***&child***REMOVED***
				// []byte for proper unicode char conversion, see #65
				n.indices = string([]byte***REMOVED***n.path[i]***REMOVED***)
				n.path = path[:i]
				n.handle = nil
				n.wildChild = false
			***REMOVED***

			// Make new node a child of this node
			if i < len(path) ***REMOVED***
				path = path[i:]

				if n.wildChild ***REMOVED***
					n = n.children[0]
					n.priority++

					// Update maxParams of the child node
					if numParams > n.maxParams ***REMOVED***
						n.maxParams = numParams
					***REMOVED***
					numParams--

					// Check if the wildcard matches
					if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
						// Check for longer wildcard, e.g. :name and :names
						(len(n.path) >= len(path) || path[len(n.path)] == '/') ***REMOVED***
						continue walk
					***REMOVED*** else ***REMOVED***
						// Wildcard conflict
						var pathSeg string
						if n.nType == catchAll ***REMOVED***
							pathSeg = path
						***REMOVED*** else ***REMOVED***
							pathSeg = strings.SplitN(path, "/", 2)[0]
						***REMOVED***
						prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
						panic("'" + pathSeg +
							"' in new path '" + fullPath +
							"' conflicts with existing wildcard '" + n.path +
							"' in existing prefix '" + prefix +
							"'")
					***REMOVED***
				***REMOVED***

				c := path[0]

				// slash after param
				if n.nType == param && c == '/' && len(n.children) == 1 ***REMOVED***
					n = n.children[0]
					n.priority++
					continue walk
				***REMOVED***

				// Check if a child with the next path byte exists
				for i := 0; i < len(n.indices); i++ ***REMOVED***
					if c == n.indices[i] ***REMOVED***
						i = n.incrementChildPrio(i)
						n = n.children[i]
						continue walk
					***REMOVED***
				***REMOVED***

				// Otherwise insert it
				if c != ':' && c != '*' ***REMOVED***
					// []byte for proper unicode char conversion, see #65
					n.indices += string([]byte***REMOVED***c***REMOVED***)
					child := &node***REMOVED***
						maxParams: numParams,
					***REMOVED***
					n.children = append(n.children, child)
					n.incrementChildPrio(len(n.indices) - 1)
					n = child
				***REMOVED***
				n.insertChild(numParams, path, fullPath, handle)
				return

			***REMOVED*** else if i == len(path) ***REMOVED*** // Make node a (in-path) leaf
				if n.handle != nil ***REMOVED***
					panic("a handle is already registered for path '" + fullPath + "'")
				***REMOVED***
				n.handle = handle
			***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED*** // Empty tree
		n.insertChild(numParams, path, fullPath, handle)
		n.nType = root
	***REMOVED***
***REMOVED***

func (n *node) insertChild(numParams uint8, path, fullPath string, handle Handle) ***REMOVED***
	var offset int // already handled bytes of the path

	// find prefix until first wildcard (beginning with ':'' or '*'')
	for i, max := 0, len(path); numParams > 0; i++ ***REMOVED***
		c := path[i]
		if c != ':' && c != '*' ***REMOVED***
			continue
		***REMOVED***

		// find wildcard end (either '/' or path end)
		end := i + 1
		for end < max && path[end] != '/' ***REMOVED***
			switch path[end] ***REMOVED***
			// the wildcard name must not contain ':' and '*'
			case ':', '*':
				panic("only one wildcard per path segment is allowed, has: '" +
					path[i:] + "' in path '" + fullPath + "'")
			default:
				end++
			***REMOVED***
		***REMOVED***

		// check if this Node existing children which would be
		// unreachable if we insert the wildcard here
		if len(n.children) > 0 ***REMOVED***
			panic("wildcard route '" + path[i:end] +
				"' conflicts with existing children in path '" + fullPath + "'")
		***REMOVED***

		// check if the wildcard has a name
		if end-i < 2 ***REMOVED***
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		***REMOVED***

		if c == ':' ***REMOVED*** // param
			// split path at the beginning of the wildcard
			if i > 0 ***REMOVED***
				n.path = path[offset:i]
				offset = i
			***REMOVED***

			child := &node***REMOVED***
				nType:     param,
				maxParams: numParams,
			***REMOVED***
			n.children = []*node***REMOVED***child***REMOVED***
			n.wildChild = true
			n = child
			n.priority++
			numParams--

			// if the path doesn't end with the wildcard, then there
			// will be another non-wildcard subpath starting with '/'
			if end < max ***REMOVED***
				n.path = path[offset:end]
				offset = end

				child := &node***REMOVED***
					maxParams: numParams,
					priority:  1,
				***REMOVED***
				n.children = []*node***REMOVED***child***REMOVED***
				n = child
			***REMOVED***

		***REMOVED*** else ***REMOVED*** // catchAll
			if end != max || numParams > 1 ***REMOVED***
				panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
			***REMOVED***

			if len(n.path) > 0 && n.path[len(n.path)-1] == '/' ***REMOVED***
				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
			***REMOVED***

			// currently fixed width 1 for '/'
			i--
			if path[i] != '/' ***REMOVED***
				panic("no / before catch-all in path '" + fullPath + "'")
			***REMOVED***

			n.path = path[offset:i]

			// first node: catchAll node with empty path
			child := &node***REMOVED***
				wildChild: true,
				nType:     catchAll,
				maxParams: 1,
			***REMOVED***
			n.children = []*node***REMOVED***child***REMOVED***
			n.indices = string(path[i])
			n = child
			n.priority++

			// second node: node holding the variable
			child = &node***REMOVED***
				path:      path[i:],
				nType:     catchAll,
				maxParams: 1,
				handle:    handle,
				priority:  1,
			***REMOVED***
			n.children = []*node***REMOVED***child***REMOVED***

			return
		***REMOVED***
	***REMOVED***

	// insert remaining path part and handle to the leaf
	n.path = path[offset:]
	n.handle = handle
***REMOVED***

// Returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *node) getValue(path string) (handle Handle, p Params, tsr bool) ***REMOVED***
walk: // outer loop for walking the tree
	for ***REMOVED***
		if len(path) > len(n.path) ***REMOVED***
			if path[:len(n.path)] == n.path ***REMOVED***
				path = path[len(n.path):]
				// If this node does not have a wildcard (param or catchAll)
				// child,  we can just look up the next child node and continue
				// to walk down the tree
				if !n.wildChild ***REMOVED***
					c := path[0]
					for i := 0; i < len(n.indices); i++ ***REMOVED***
						if c == n.indices[i] ***REMOVED***
							n = n.children[i]
							continue walk
						***REMOVED***
					***REMOVED***

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					tsr = (path == "/" && n.handle != nil)
					return

				***REMOVED***

				// handle wildcard child
				n = n.children[0]
				switch n.nType ***REMOVED***
				case param:
					// find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' ***REMOVED***
						end++
					***REMOVED***

					// save param value
					if p == nil ***REMOVED***
						// lazy allocation
						p = make(Params, 0, n.maxParams)
					***REMOVED***
					i := len(p)
					p = p[:i+1] // expand slice within preallocated capacity
					p[i].Key = n.path[1:]
					p[i].Value = path[:end]

					// we need to go deeper!
					if end < len(path) ***REMOVED***
						if len(n.children) > 0 ***REMOVED***
							path = path[end:]
							n = n.children[0]
							continue walk
						***REMOVED***

						// ... but we can't
						tsr = (len(path) == end+1)
						return
					***REMOVED***

					if handle = n.handle; handle != nil ***REMOVED***
						return
					***REMOVED*** else if len(n.children) == 1 ***REMOVED***
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[0]
						tsr = (n.path == "/" && n.handle != nil)
					***REMOVED***

					return

				case catchAll:
					// save param value
					if p == nil ***REMOVED***
						// lazy allocation
						p = make(Params, 0, n.maxParams)
					***REMOVED***
					i := len(p)
					p = p[:i+1] // expand slice within preallocated capacity
					p[i].Key = n.path[2:]
					p[i].Value = path

					handle = n.handle
					return

				default:
					panic("invalid node type")
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if path == n.path ***REMOVED***
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if handle = n.handle; handle != nil ***REMOVED***
				return
			***REMOVED***

			if path == "/" && n.wildChild && n.nType != root ***REMOVED***
				tsr = true
				return
			***REMOVED***

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i := 0; i < len(n.indices); i++ ***REMOVED***
				if n.indices[i] == '/' ***REMOVED***
					n = n.children[i]
					tsr = (len(n.path) == 1 && n.handle != nil) ||
						(n.nType == catchAll && n.children[0].handle != nil)
					return
				***REMOVED***
			***REMOVED***

			return
		***REMOVED***

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		tsr = (path == "/") ||
			(len(n.path) == len(path)+1 && n.path[len(path)] == '/' &&
				path == n.path[:len(n.path)-1] && n.handle != nil)
		return
	***REMOVED***
***REMOVED***

// Makes a case-insensitive lookup of the given path and tries to find a handler.
// It can optionally also fix trailing slashes.
// It returns the case-corrected path and a bool indicating whether the lookup
// was successful.
func (n *node) findCaseInsensitivePath(path string, fixTrailingSlash bool) (ciPath []byte, found bool) ***REMOVED***
	return n.findCaseInsensitivePathRec(
		path,
		strings.ToLower(path),
		make([]byte, 0, len(path)+1), // preallocate enough memory for new path
		[4]byte***REMOVED******REMOVED***,                    // empty rune buffer
		fixTrailingSlash,
	)
***REMOVED***

// shift bytes in array by n bytes left
func shiftNRuneBytes(rb [4]byte, n int) [4]byte ***REMOVED***
	switch n ***REMOVED***
	case 0:
		return rb
	case 1:
		return [4]byte***REMOVED***rb[1], rb[2], rb[3], 0***REMOVED***
	case 2:
		return [4]byte***REMOVED***rb[2], rb[3]***REMOVED***
	case 3:
		return [4]byte***REMOVED***rb[3]***REMOVED***
	default:
		return [4]byte***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// recursive case-insensitive lookup function used by n.findCaseInsensitivePath
func (n *node) findCaseInsensitivePathRec(path, loPath string, ciPath []byte, rb [4]byte, fixTrailingSlash bool) ([]byte, bool) ***REMOVED***
	loNPath := strings.ToLower(n.path)

walk: // outer loop for walking the tree
	for len(loPath) >= len(loNPath) && (len(loNPath) == 0 || loPath[1:len(loNPath)] == loNPath[1:]) ***REMOVED***
		// add common path to result
		ciPath = append(ciPath, n.path...)

		if path = path[len(n.path):]; len(path) > 0 ***REMOVED***
			loOld := loPath
			loPath = loPath[len(loNPath):]

			// If this node does not have a wildcard (param or catchAll) child,
			// we can just look up the next child node and continue to walk down
			// the tree
			if !n.wildChild ***REMOVED***
				// skip rune bytes already processed
				rb = shiftNRuneBytes(rb, len(loNPath))

				if rb[0] != 0 ***REMOVED***
					// old rune not finished
					for i := 0; i < len(n.indices); i++ ***REMOVED***
						if n.indices[i] == rb[0] ***REMOVED***
							// continue with child node
							n = n.children[i]
							loNPath = strings.ToLower(n.path)
							continue walk
						***REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					// process a new rune
					var rv rune

					// find rune start
					// runes are up to 4 byte long,
					// -4 would definitely be another rune
					var off int
					for max := min(len(loNPath), 3); off < max; off++ ***REMOVED***
						if i := len(loNPath) - off; utf8.RuneStart(loOld[i]) ***REMOVED***
							// read rune from cached lowercase path
							rv, _ = utf8.DecodeRuneInString(loOld[i:])
							break
						***REMOVED***
					***REMOVED***

					// calculate lowercase bytes of current rune
					utf8.EncodeRune(rb[:], rv)
					// skipp already processed bytes
					rb = shiftNRuneBytes(rb, off)

					for i := 0; i < len(n.indices); i++ ***REMOVED***
						// lowercase matches
						if n.indices[i] == rb[0] ***REMOVED***
							// must use a recursive approach since both the
							// uppercase byte and the lowercase byte might exist
							// as an index
							if out, found := n.children[i].findCaseInsensitivePathRec(
								path, loPath, ciPath, rb, fixTrailingSlash,
							); found ***REMOVED***
								return out, true
							***REMOVED***
							break
						***REMOVED***
					***REMOVED***

					// same for uppercase rune, if it differs
					if up := unicode.ToUpper(rv); up != rv ***REMOVED***
						utf8.EncodeRune(rb[:], up)
						rb = shiftNRuneBytes(rb, off)

						for i := 0; i < len(n.indices); i++ ***REMOVED***
							// uppercase matches
							if n.indices[i] == rb[0] ***REMOVED***
								// continue with child node
								n = n.children[i]
								loNPath = strings.ToLower(n.path)
								continue walk
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***

				// Nothing found. We can recommend to redirect to the same URL
				// without a trailing slash if a leaf exists for that path
				return ciPath, (fixTrailingSlash && path == "/" && n.handle != nil)
			***REMOVED***

			n = n.children[0]
			switch n.nType ***REMOVED***
			case param:
				// find param end (either '/' or path end)
				k := 0
				for k < len(path) && path[k] != '/' ***REMOVED***
					k++
				***REMOVED***

				// add param value to case insensitive path
				ciPath = append(ciPath, path[:k]...)

				// we need to go deeper!
				if k < len(path) ***REMOVED***
					if len(n.children) > 0 ***REMOVED***
						// continue with child node
						n = n.children[0]
						loNPath = strings.ToLower(n.path)
						loPath = loPath[k:]
						path = path[k:]
						continue
					***REMOVED***

					// ... but we can't
					if fixTrailingSlash && len(path) == k+1 ***REMOVED***
						return ciPath, true
					***REMOVED***
					return ciPath, false
				***REMOVED***

				if n.handle != nil ***REMOVED***
					return ciPath, true
				***REMOVED*** else if fixTrailingSlash && len(n.children) == 1 ***REMOVED***
					// No handle found. Check if a handle for this path + a
					// trailing slash exists
					n = n.children[0]
					if n.path == "/" && n.handle != nil ***REMOVED***
						return append(ciPath, '/'), true
					***REMOVED***
				***REMOVED***
				return ciPath, false

			case catchAll:
				return append(ciPath, path...), true

			default:
				panic("invalid node type")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if n.handle != nil ***REMOVED***
				return ciPath, true
			***REMOVED***

			// No handle found.
			// Try to fix the path by adding a trailing slash
			if fixTrailingSlash ***REMOVED***
				for i := 0; i < len(n.indices); i++ ***REMOVED***
					if n.indices[i] == '/' ***REMOVED***
						n = n.children[i]
						if (len(n.path) == 1 && n.handle != nil) ||
							(n.nType == catchAll && n.children[0].handle != nil) ***REMOVED***
							return append(ciPath, '/'), true
						***REMOVED***
						return ciPath, false
					***REMOVED***
				***REMOVED***
			***REMOVED***
			return ciPath, false
		***REMOVED***
	***REMOVED***

	// Nothing found.
	// Try to fix the path by adding / removing a trailing slash
	if fixTrailingSlash ***REMOVED***
		if path == "/" ***REMOVED***
			return ciPath, true
		***REMOVED***
		if len(loPath)+1 == len(loNPath) && loNPath[len(loPath)] == '/' &&
			loPath[1:] == loNPath[1:len(loPath)] && n.handle != nil ***REMOVED***
			return append(ciPath, n.path...), true
		***REMOVED***
	***REMOVED***
	return ciPath, false
***REMOVED***
