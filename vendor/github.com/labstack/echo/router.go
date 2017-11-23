package echo

import "strings"

type (
	// Router is the registry of all registered routes for an `Echo` instance for
	// request matching and URL path parameter parsing.
	Router struct ***REMOVED***
		tree   *node
		routes map[string]*Route
		echo   *Echo
	***REMOVED***
	node struct ***REMOVED***
		kind          kind
		label         byte
		prefix        string
		parent        *node
		children      children
		ppath         string
		pnames        []string
		methodHandler *methodHandler
	***REMOVED***
	kind          uint8
	children      []*node
	methodHandler struct ***REMOVED***
		connect HandlerFunc
		delete  HandlerFunc
		get     HandlerFunc
		head    HandlerFunc
		options HandlerFunc
		patch   HandlerFunc
		post    HandlerFunc
		put     HandlerFunc
		trace   HandlerFunc
	***REMOVED***
)

const (
	skind kind = iota
	pkind
	akind
)

// NewRouter returns a new Router instance.
func NewRouter(e *Echo) *Router ***REMOVED***
	return &Router***REMOVED***
		tree: &node***REMOVED***
			methodHandler: new(methodHandler),
		***REMOVED***,
		routes: map[string]*Route***REMOVED******REMOVED***,
		echo:   e,
	***REMOVED***
***REMOVED***

// Add registers a new route for method and path with matching handler.
func (r *Router) Add(method, path string, h HandlerFunc) ***REMOVED***
	// Validate path
	if path == "" ***REMOVED***
		panic("echo: path cannot be empty")
	***REMOVED***
	if path[0] != '/' ***REMOVED***
		path = "/" + path
	***REMOVED***
	ppath := path        // Pristine path
	pnames := []string***REMOVED******REMOVED*** // Param names

	for i, l := 0, len(path); i < l; i++ ***REMOVED***
		if path[i] == ':' ***REMOVED***
			j := i + 1

			r.insert(method, path[:i], nil, skind, "", nil)
			for ; i < l && path[i] != '/'; i++ ***REMOVED***
			***REMOVED***

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l ***REMOVED***
				r.insert(method, path[:i], h, pkind, ppath, pnames)
				return
			***REMOVED***
			r.insert(method, path[:i], nil, pkind, ppath, pnames)
		***REMOVED*** else if path[i] == '*' ***REMOVED***
			r.insert(method, path[:i], nil, skind, "", nil)
			pnames = append(pnames, "*")
			r.insert(method, path[:i+1], h, akind, ppath, pnames)
			return
		***REMOVED***
	***REMOVED***

	r.insert(method, path, h, skind, ppath, pnames)
***REMOVED***

func (r *Router) insert(method, path string, h HandlerFunc, t kind, ppath string, pnames []string) ***REMOVED***
	// Adjust max param
	l := len(pnames)
	if *r.echo.maxParam < l ***REMOVED***
		*r.echo.maxParam = l
	***REMOVED***

	cn := r.tree // Current node as root
	if cn == nil ***REMOVED***
		panic("echo: invalid method")
	***REMOVED***
	search := path

	for ***REMOVED***
		sl := len(search)
		pl := len(cn.prefix)
		l := 0

		// LCP
		max := pl
		if sl < max ***REMOVED***
			max = sl
		***REMOVED***
		for ; l < max && search[l] == cn.prefix[l]; l++ ***REMOVED***
		***REMOVED***

		if l == 0 ***REMOVED***
			// At root node
			cn.label = search[0]
			cn.prefix = search
			if h != nil ***REMOVED***
				cn.kind = t
				cn.addHandler(method, h)
				cn.ppath = ppath
				cn.pnames = pnames
			***REMOVED***
		***REMOVED*** else if l < pl ***REMOVED***
			// Split node
			n := newNode(cn.kind, cn.prefix[l:], cn, cn.children, cn.methodHandler, cn.ppath, cn.pnames)

			// Reset parent node
			cn.kind = skind
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil
			cn.methodHandler = new(methodHandler)
			cn.ppath = ""
			cn.pnames = nil

			cn.addChild(n)

			if l == sl ***REMOVED***
				// At parent node
				cn.kind = t
				cn.addHandler(method, h)
				cn.ppath = ppath
				cn.pnames = pnames
			***REMOVED*** else ***REMOVED***
				// Create child node
				n = newNode(t, search[l:], cn, nil, new(methodHandler), ppath, pnames)
				n.addHandler(method, h)
				cn.addChild(n)
			***REMOVED***
		***REMOVED*** else if l < sl ***REMOVED***
			search = search[l:]
			c := cn.findChildWithLabel(search[0])
			if c != nil ***REMOVED***
				// Go deeper
				cn = c
				continue
			***REMOVED***
			// Create child node
			n := newNode(t, search, cn, nil, new(methodHandler), ppath, pnames)
			n.addHandler(method, h)
			cn.addChild(n)
		***REMOVED*** else ***REMOVED***
			// Node already exists
			if h != nil ***REMOVED***
				cn.addHandler(method, h)
				cn.ppath = ppath
				if len(cn.pnames) == 0 ***REMOVED*** // Issue #729
					cn.pnames = pnames
				***REMOVED***
				for i, n := range pnames ***REMOVED***
					// Param name aliases
					if i < len(cn.pnames) && !strings.Contains(cn.pnames[i], n) ***REMOVED***
						cn.pnames[i] += "," + n
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***
***REMOVED***

func newNode(t kind, pre string, p *node, c children, mh *methodHandler, ppath string, pnames []string) *node ***REMOVED***
	return &node***REMOVED***
		kind:          t,
		label:         pre[0],
		prefix:        pre,
		parent:        p,
		children:      c,
		ppath:         ppath,
		pnames:        pnames,
		methodHandler: mh,
	***REMOVED***
***REMOVED***

func (n *node) addChild(c *node) ***REMOVED***
	n.children = append(n.children, c)
***REMOVED***

func (n *node) findChild(l byte, t kind) *node ***REMOVED***
	for _, c := range n.children ***REMOVED***
		if c.label == l && c.kind == t ***REMOVED***
			return c
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (n *node) findChildWithLabel(l byte) *node ***REMOVED***
	for _, c := range n.children ***REMOVED***
		if c.label == l ***REMOVED***
			return c
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (n *node) findChildByKind(t kind) *node ***REMOVED***
	for _, c := range n.children ***REMOVED***
		if c.kind == t ***REMOVED***
			return c
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (n *node) addHandler(method string, h HandlerFunc) ***REMOVED***
	switch method ***REMOVED***
	case GET:
		n.methodHandler.get = h
	case POST:
		n.methodHandler.post = h
	case PUT:
		n.methodHandler.put = h
	case DELETE:
		n.methodHandler.delete = h
	case PATCH:
		n.methodHandler.patch = h
	case OPTIONS:
		n.methodHandler.options = h
	case HEAD:
		n.methodHandler.head = h
	case CONNECT:
		n.methodHandler.connect = h
	case TRACE:
		n.methodHandler.trace = h
	***REMOVED***
***REMOVED***

func (n *node) findHandler(method string) HandlerFunc ***REMOVED***
	switch method ***REMOVED***
	case GET:
		return n.methodHandler.get
	case POST:
		return n.methodHandler.post
	case PUT:
		return n.methodHandler.put
	case DELETE:
		return n.methodHandler.delete
	case PATCH:
		return n.methodHandler.patch
	case OPTIONS:
		return n.methodHandler.options
	case HEAD:
		return n.methodHandler.head
	case CONNECT:
		return n.methodHandler.connect
	case TRACE:
		return n.methodHandler.trace
	default:
		return nil
	***REMOVED***
***REMOVED***

func (n *node) checkMethodNotAllowed() HandlerFunc ***REMOVED***
	for _, m := range methods ***REMOVED***
		if h := n.findHandler(m); h != nil ***REMOVED***
			return MethodNotAllowedHandler
		***REMOVED***
	***REMOVED***
	return NotFoundHandler
***REMOVED***

// Find lookup a handler registered for method and path. It also parses URL for path
// parameters and load them into context.
//
// For performance:
//
// - Get context from `Echo#AcquireContext()`
// - Reset it `Context#Reset()`
// - Return it `Echo#ReleaseContext()`.
func (r *Router) Find(method, path string, c Context) ***REMOVED***
	ctx := c.(*context)
	ctx.path = path
	cn := r.tree // Current node as root

	var (
		search  = path
		child   *node         // Child node
		n       int           // Param counter
		nk      kind          // Next kind
		nn      *node         // Next node
		ns      string        // Next search
		pvalues = ctx.pvalues // Use the internal slice so the interface can keep the illusion of a dynamic slice
	)

	// Search order static > param > any
	for ***REMOVED***
		if search == "" ***REMOVED***
			goto End
		***REMOVED***

		pl := 0 // Prefix length
		l := 0  // LCP length

		if cn.label != ':' ***REMOVED***
			sl := len(search)
			pl = len(cn.prefix)

			// LCP
			max := pl
			if sl < max ***REMOVED***
				max = sl
			***REMOVED***
			for ; l < max && search[l] == cn.prefix[l]; l++ ***REMOVED***
			***REMOVED***
		***REMOVED***

		if l == pl ***REMOVED***
			// Continue search
			search = search[l:]
		***REMOVED*** else ***REMOVED***
			cn = nn
			search = ns
			if nk == pkind ***REMOVED***
				goto Param
			***REMOVED*** else if nk == akind ***REMOVED***
				goto Any
			***REMOVED***
			// Not found
			return
		***REMOVED***

		if search == "" ***REMOVED***
			goto End
		***REMOVED***

		// Static node
		if child = cn.findChild(search[0], skind); child != nil ***REMOVED***
			// Save next
			if cn.prefix[len(cn.prefix)-1] == '/' ***REMOVED*** // Issue #623
				nk = pkind
				nn = cn
				ns = search
			***REMOVED***
			cn = child
			continue
		***REMOVED***

		// Param node
	Param:
		if child = cn.findChildByKind(pkind); child != nil ***REMOVED***
			// Issue #378
			if len(pvalues) == n ***REMOVED***
				continue
			***REMOVED***

			// Save next
			if cn.prefix[len(cn.prefix)-1] == '/' ***REMOVED*** // Issue #623
				nk = akind
				nn = cn
				ns = search
			***REMOVED***

			cn = child
			i, l := 0, len(search)
			for ; i < l && search[i] != '/'; i++ ***REMOVED***
			***REMOVED***
			pvalues[n] = search[:i]
			n++
			search = search[i:]
			continue
		***REMOVED***

		// Any node
	Any:
		if cn = cn.findChildByKind(akind); cn == nil ***REMOVED***
			if nn != nil ***REMOVED***
				cn = nn
				nn = cn.parent // Next (Issue #954)
				search = ns
				if nk == pkind ***REMOVED***
					goto Param
				***REMOVED*** else if nk == akind ***REMOVED***
					goto Any
				***REMOVED***
			***REMOVED***
			// Not found
			return
		***REMOVED***
		pvalues[len(cn.pnames)-1] = search
		goto End
	***REMOVED***

End:
	ctx.handler = cn.findHandler(method)
	ctx.path = cn.ppath
	ctx.pnames = cn.pnames

	// NOTE: Slow zone...
	if ctx.handler == nil ***REMOVED***
		ctx.handler = cn.checkMethodNotAllowed()

		// Dig further for any, might have an empty value for *, e.g.
		// serving a directory. Issue #207.
		if cn = cn.findChildByKind(akind); cn == nil ***REMOVED***
			return
		***REMOVED***
		if h := cn.findHandler(method); h != nil ***REMOVED***
			ctx.handler = h
		***REMOVED*** else ***REMOVED***
			ctx.handler = cn.checkMethodNotAllowed()
		***REMOVED***
		ctx.path = cn.ppath
		ctx.pnames = cn.pnames
		pvalues[len(cn.pnames)-1] = ""
	***REMOVED***

	return
***REMOVED***
