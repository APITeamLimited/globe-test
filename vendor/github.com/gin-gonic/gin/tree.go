// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

package gin

import (
	"net/url"
	"strings"
	"unicode"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct ***REMOVED***
	Key   string
	Value string
***REMOVED***

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// Get returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) Get(name string) (string, bool) ***REMOVED***
	for _, entry := range ps ***REMOVED***
		if entry.Key == name ***REMOVED***
			return entry.Value, true
		***REMOVED***
	***REMOVED***
	return "", false
***REMOVED***

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) (va string) ***REMOVED***
	va, _ = ps.Get(name)
	return
***REMOVED***

type methodTree struct ***REMOVED***
	method string
	root   *node
***REMOVED***

type methodTrees []methodTree

func (trees methodTrees) get(method string) *node ***REMOVED***
	for _, tree := range trees ***REMOVED***
		if tree.method == method ***REMOVED***
			return tree.root
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

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
	handlers  HandlersChain
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
		tmpN := n.children[newPos-1]
		n.children[newPos-1] = n.children[newPos]
		n.children[newPos] = tmpN

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
func (n *node) addRoute(path string, handlers HandlersChain) ***REMOVED***
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
					indices:   n.indices,
					children:  n.children,
					handlers:  n.handlers,
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
				n.handlers = nil
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
					if len(path) >= len(n.path) && n.path == path[:len(n.path)] ***REMOVED***
						// check for longer wildcard, e.g. :name and :names
						if len(n.path) >= len(path) || path[len(n.path)] == '/' ***REMOVED***
							continue walk
						***REMOVED***
					***REMOVED***

					panic("path segment '" + path +
						"' conflicts with existing wildcard '" + n.path +
						"' in path '" + fullPath + "'")
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
				n.insertChild(numParams, path, fullPath, handlers)
				return

			***REMOVED*** else if i == len(path) ***REMOVED*** // Make node a (in-path) leaf
				if n.handlers != nil ***REMOVED***
					panic("handlers are already registered for path ''" + fullPath + "'")
				***REMOVED***
				n.handlers = handlers
			***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED*** // Empty tree
		n.insertChild(numParams, path, fullPath, handlers)
		n.nType = root
	***REMOVED***
***REMOVED***

func (n *node) insertChild(numParams uint8, path string, fullPath string, handlers HandlersChain) ***REMOVED***
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
				handlers:  handlers,
				priority:  1,
			***REMOVED***
			n.children = []*node***REMOVED***child***REMOVED***

			return
		***REMOVED***
	***REMOVED***

	// insert remaining path part and handle to the leaf
	n.path = path[offset:]
	n.handlers = handlers
***REMOVED***

// Returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *node) getValue(path string, po Params, unescape bool) (handlers HandlersChain, p Params, tsr bool) ***REMOVED***
	p = po
walk: // Outer loop for walking the tree
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
					tsr = (path == "/" && n.handlers != nil)
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
					if cap(p) < int(n.maxParams) ***REMOVED***
						p = make(Params, 0, n.maxParams)
					***REMOVED***
					i := len(p)
					p = p[:i+1] // expand slice within preallocated capacity
					p[i].Key = n.path[1:]
					val := path[:end]
					if unescape ***REMOVED***
						var err error
						if p[i].Value, err = url.QueryUnescape(val); err != nil ***REMOVED***
							p[i].Value = val // fallback, in case of error
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						p[i].Value = val
					***REMOVED***

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

					if handlers = n.handlers; handlers != nil ***REMOVED***
						return
					***REMOVED***
					if len(n.children) == 1 ***REMOVED***
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[0]
						tsr = (n.path == "/" && n.handlers != nil)
					***REMOVED***

					return

				case catchAll:
					// save param value
					if cap(p) < int(n.maxParams) ***REMOVED***
						p = make(Params, 0, n.maxParams)
					***REMOVED***
					i := len(p)
					p = p[:i+1] // expand slice within preallocated capacity
					p[i].Key = n.path[2:]
					if unescape ***REMOVED***
						var err error
						if p[i].Value, err = url.QueryUnescape(path); err != nil ***REMOVED***
							p[i].Value = path // fallback, in case of error
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						p[i].Value = path
					***REMOVED***

					handlers = n.handlers
					return

				default:
					panic("invalid node type")
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if path == n.path ***REMOVED***
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if handlers = n.handlers; handlers != nil ***REMOVED***
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
					tsr = (len(n.path) == 1 && n.handlers != nil) ||
						(n.nType == catchAll && n.children[0].handlers != nil)
					return
				***REMOVED***
			***REMOVED***

			return
		***REMOVED***

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		tsr = (path == "/") ||
			(len(n.path) == len(path)+1 && n.path[len(path)] == '/' &&
				path == n.path[:len(n.path)-1] && n.handlers != nil)
		return
	***REMOVED***
***REMOVED***

// Makes a case-insensitive lookup of the given path and tries to find a handler.
// It can optionally also fix trailing slashes.
// It returns the case-corrected path and a bool indicating whether the lookup
// was successful.
func (n *node) findCaseInsensitivePath(path string, fixTrailingSlash bool) (ciPath []byte, found bool) ***REMOVED***
	ciPath = make([]byte, 0, len(path)+1) // preallocate enough memory

	// Outer loop for walking the tree
	for len(path) >= len(n.path) && strings.ToLower(path[:len(n.path)]) == strings.ToLower(n.path) ***REMOVED***
		path = path[len(n.path):]
		ciPath = append(ciPath, n.path...)

		if len(path) > 0 ***REMOVED***
			// If this node does not have a wildcard (param or catchAll) child,
			// we can just look up the next child node and continue to walk down
			// the tree
			if !n.wildChild ***REMOVED***
				r := unicode.ToLower(rune(path[0]))
				for i, index := range n.indices ***REMOVED***
					// must use recursive approach since both index and
					// ToLower(index) could exist. We must check both.
					if r == unicode.ToLower(index) ***REMOVED***
						out, found := n.children[i].findCaseInsensitivePath(path, fixTrailingSlash)
						if found ***REMOVED***
							return append(ciPath, out...), true
						***REMOVED***
					***REMOVED***
				***REMOVED***

				// Nothing found. We can recommend to redirect to the same URL
				// without a trailing slash if a leaf exists for that path
				found = (fixTrailingSlash && path == "/" && n.handlers != nil)
				return
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
						path = path[k:]
						n = n.children[0]
						continue
					***REMOVED***

					// ... but we can't
					if fixTrailingSlash && len(path) == k+1 ***REMOVED***
						return ciPath, true
					***REMOVED***
					return
				***REMOVED***

				if n.handlers != nil ***REMOVED***
					return ciPath, true
				***REMOVED*** else if fixTrailingSlash && len(n.children) == 1 ***REMOVED***
					// No handle found. Check if a handle for this path + a
					// trailing slash exists
					n = n.children[0]
					if n.path == "/" && n.handlers != nil ***REMOVED***
						return append(ciPath, '/'), true
					***REMOVED***
				***REMOVED***
				return

			case catchAll:
				return append(ciPath, path...), true

			default:
				panic("invalid node type")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if n.handlers != nil ***REMOVED***
				return ciPath, true
			***REMOVED***

			// No handle found.
			// Try to fix the path by adding a trailing slash
			if fixTrailingSlash ***REMOVED***
				for i := 0; i < len(n.indices); i++ ***REMOVED***
					if n.indices[i] == '/' ***REMOVED***
						n = n.children[i]
						if (len(n.path) == 1 && n.handlers != nil) ||
							(n.nType == catchAll && n.children[0].handlers != nil) ***REMOVED***
							return append(ciPath, '/'), true
						***REMOVED***
						return
					***REMOVED***
				***REMOVED***
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	// Nothing found.
	// Try to fix the path by adding / removing a trailing slash
	if fixTrailingSlash ***REMOVED***
		if path == "/" ***REMOVED***
			return ciPath, true
		***REMOVED***
		if len(path)+1 == len(n.path) && n.path[len(path)] == '/' &&
			strings.ToLower(path) == strings.ToLower(n.path[:len(path)]) &&
			n.handlers != nil ***REMOVED***
			return append(ciPath, n.path...), true
		***REMOVED***
	***REMOVED***
	return
***REMOVED***
