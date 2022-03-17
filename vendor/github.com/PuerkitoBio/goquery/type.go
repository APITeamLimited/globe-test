package goquery

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

// Document represents an HTML document to be manipulated. Unlike jQuery, which
// is loaded as part of a DOM document, and thus acts upon its containing
// document, GoQuery doesn't know which HTML document to act upon. So it needs
// to be told, and that's what the Document class is for. It holds the root
// document node to manipulate, and can make selections on this document.
type Document struct ***REMOVED***
	*Selection
	Url      *url.URL
	rootNode *html.Node
***REMOVED***

// NewDocumentFromNode is a Document constructor that takes a root html Node
// as argument.
func NewDocumentFromNode(root *html.Node) *Document ***REMOVED***
	return newDocument(root, nil)
***REMOVED***

// NewDocument is a Document constructor that takes a string URL as argument.
// It loads the specified document, parses it, and stores the root Document
// node, ready to be manipulated.
//
// Deprecated: Use the net/http standard library package to make the request
// and validate the response before calling goquery.NewDocumentFromReader
// with the response's body.
func NewDocument(url string) (*Document, error) ***REMOVED***
	// Load the URL
	res, e := http.Get(url)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***
	return NewDocumentFromResponse(res)
***REMOVED***

// NewDocumentFromReader returns a Document from an io.Reader.
// It returns an error as second value if the reader's data cannot be parsed
// as html. It does not check if the reader is also an io.Closer, the
// provided reader is never closed by this call. It is the responsibility
// of the caller to close it if required.
func NewDocumentFromReader(r io.Reader) (*Document, error) ***REMOVED***
	root, e := html.Parse(r)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***
	return newDocument(root, nil), nil
***REMOVED***

// NewDocumentFromResponse is another Document constructor that takes an http response as argument.
// It loads the specified response's document, parses it, and stores the root Document
// node, ready to be manipulated. The response's body is closed on return.
//
// Deprecated: Use goquery.NewDocumentFromReader with the response's body.
func NewDocumentFromResponse(res *http.Response) (*Document, error) ***REMOVED***
	if res == nil ***REMOVED***
		return nil, errors.New("Response is nil")
	***REMOVED***
	defer res.Body.Close()
	if res.Request == nil ***REMOVED***
		return nil, errors.New("Response.Request is nil")
	***REMOVED***

	// Parse the HTML into nodes
	root, e := html.Parse(res.Body)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***

	// Create and fill the document
	return newDocument(root, res.Request.URL), nil
***REMOVED***

// CloneDocument creates a deep-clone of a document.
func CloneDocument(doc *Document) *Document ***REMOVED***
	return newDocument(cloneNode(doc.rootNode), doc.Url)
***REMOVED***

// Private constructor, make sure all fields are correctly filled.
func newDocument(root *html.Node, url *url.URL) *Document ***REMOVED***
	// Create and fill the document
	d := &Document***REMOVED***nil, url, root***REMOVED***
	d.Selection = newSingleSelection(root, d)
	return d
***REMOVED***

// Selection represents a collection of nodes matching some criteria. The
// initial Selection can be created by using Document.Find, and then
// manipulated using the jQuery-like chainable syntax and methods.
type Selection struct ***REMOVED***
	Nodes    []*html.Node
	document *Document
	prevSel  *Selection
***REMOVED***

// Helper constructor to create an empty selection
func newEmptySelection(doc *Document) *Selection ***REMOVED***
	return &Selection***REMOVED***nil, doc, nil***REMOVED***
***REMOVED***

// Helper constructor to create a selection of only one node
func newSingleSelection(node *html.Node, doc *Document) *Selection ***REMOVED***
	return &Selection***REMOVED***[]*html.Node***REMOVED***node***REMOVED***, doc, nil***REMOVED***
***REMOVED***

// Matcher is an interface that defines the methods to match
// HTML nodes against a compiled selector string. Cascadia's
// Selector implements this interface.
type Matcher interface ***REMOVED***
	Match(*html.Node) bool
	MatchAll(*html.Node) []*html.Node
	Filter([]*html.Node) []*html.Node
***REMOVED***

// Single compiles a selector string to a Matcher that stops after the first
// match is found.
//
// By default, Selection.Find and other functions that accept a selector string
// to select nodes will use all matches corresponding to that selector. By
// using the Matcher returned by Single, at most the first match will be
// selected.
//
// For example, those two statements are semantically equivalent:
//
//     sel1 := doc.Find("a").First()
//     sel2 := doc.FindMatcher(goquery.Single("a"))
//
// The one using Single is optimized to be potentially much faster on large
// documents.
//
// Only the behaviour of the MatchAll method of the Matcher interface is
// altered compared to standard Matchers. This means that the single-selection
// property of the Matcher only applies for Selection methods where the Matcher
// is used to select nodes, not to filter or check if a node matches the
// Matcher - in those cases, the behaviour of the Matcher is unchanged (e.g.
// FilterMatcher(Single("div")) will still result in a Selection with multiple
// "div"s if there were many "div"s in the Selection to begin with).
func Single(selector string) Matcher ***REMOVED***
	return singleMatcher***REMOVED***compileMatcher(selector)***REMOVED***
***REMOVED***

// SingleMatcher returns a Matcher matches the same nodes as m, but that stops
// after the first match is found.
//
// See the documentation of function Single for more details.
func SingleMatcher(m Matcher) Matcher ***REMOVED***
	if _, ok := m.(singleMatcher); ok ***REMOVED***
		// m is already a singleMatcher
		return m
	***REMOVED***
	return singleMatcher***REMOVED***m***REMOVED***
***REMOVED***

// compileMatcher compiles the selector string s and returns
// the corresponding Matcher. If s is an invalid selector string,
// it returns a Matcher that fails all matches.
func compileMatcher(s string) Matcher ***REMOVED***
	cs, err := cascadia.Compile(s)
	if err != nil ***REMOVED***
		return invalidMatcher***REMOVED******REMOVED***
	***REMOVED***
	return cs
***REMOVED***

type singleMatcher struct ***REMOVED***
	Matcher
***REMOVED***

func (m singleMatcher) MatchAll(n *html.Node) []*html.Node ***REMOVED***
	// Optimized version - stops finding at the first match (cascadia-compiled
	// matchers all use this code path).
	if mm, ok := m.Matcher.(interface***REMOVED*** MatchFirst(*html.Node) *html.Node ***REMOVED***); ok ***REMOVED***
		node := mm.MatchFirst(n)
		if node == nil ***REMOVED***
			return nil
		***REMOVED***
		return []*html.Node***REMOVED***node***REMOVED***
	***REMOVED***

	// Fallback version, for e.g. test mocks that don't provide the MatchFirst
	// method.
	nodes := m.Matcher.MatchAll(n)
	if len(nodes) > 0 ***REMOVED***
		return nodes[:1:1]
	***REMOVED***
	return nil
***REMOVED***

// invalidMatcher is a Matcher that always fails to match.
type invalidMatcher struct***REMOVED******REMOVED***

func (invalidMatcher) Match(n *html.Node) bool             ***REMOVED*** return false ***REMOVED***
func (invalidMatcher) MatchAll(n *html.Node) []*html.Node  ***REMOVED*** return nil ***REMOVED***
func (invalidMatcher) Filter(ns []*html.Node) []*html.Node ***REMOVED*** return nil ***REMOVED***
