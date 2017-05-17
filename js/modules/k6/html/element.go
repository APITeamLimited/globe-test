package html

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	gohtml "golang.org/x/net/html"
)

var (
	protoPrg *goja.Program
)

type Element struct ***REMOVED***
	sel  *Selection
	rt   *goja.Runtime
	qsel *goquery.Selection
	node *gohtml.Node
***REMOVED***

type Attribute struct ***REMOVED***
	Name         string
	NamespaceURI string
	LocalName    string
	Prefix       string
	OwnerElement goja.Value
	Value        string
***REMOVED***

func (e Element) GetAttribute(name string) goja.Value ***REMOVED***
	return e.sel.Attr(name)
***REMOVED***

func (e Element) GetAttributeNode(self goja.Value, name string) goja.Value ***REMOVED***
	if attr := getHtmlAttr(e.node, name); attr != nil ***REMOVED***
		return e.rt.ToValue(Attribute***REMOVED***attr.Key, attr.Namespace, attr.Namespace, attr.Namespace, self, attr.Val***REMOVED***)
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) HasAttribute(name string) bool ***REMOVED***
	return e.sel.Attr(name) != goja.Undefined()
***REMOVED***

func (e Element) HasAttributes() bool ***REMOVED***
	return e.qsel.Length() > 0 && len(e.node.Attr) > 0
***REMOVED***

func (e Element) Attributes(self goja.Value) map[string]Attribute ***REMOVED***
	attrs := make(map[string]Attribute)
	for i := 0; i < len(e.node.Attr); i++ ***REMOVED***
		attr := e.node.Attr[i]
		attrs[attr.Key] = Attribute***REMOVED***attr.Key, attr.Namespace, attr.Namespace, attr.Namespace, self, attr.Val***REMOVED***
	***REMOVED***
	return attrs
***REMOVED***

func (e Element) ToString() goja.Value ***REMOVED***
	if e.qsel.Length() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED*** else if e.node.Type == gohtml.ElementNode ***REMOVED***
		return e.rt.ToValue("[object html.Node]")
	***REMOVED*** else ***REMOVED***
		return e.rt.ToValue(fmt.Sprintf("[object %s]", e.NodeName()))
	***REMOVED***
***REMOVED***

func (e Element) HasChildNodes() bool ***REMOVED***
	return e.qsel.Length() > 0 && e.node.FirstChild != nil
***REMOVED***

func (e Element) TextContent() string ***REMOVED***
	return e.qsel.Text()
***REMOVED***

func (e Element) Id() goja.Value ***REMOVED***
	return e.GetAttribute("id")
***REMOVED***

func (e Element) IsEqualNode(v goja.Value) bool ***REMOVED***
	if other, ok := valToElement(v); ok ***REMOVED***
		htmlA, errA := e.qsel.Html()
		htmlB, errB := other.qsel.Html()

		return errA == nil && errB == nil && htmlA == htmlB
	***REMOVED*** else ***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func (e Element) IsSameNode(v goja.Value) bool ***REMOVED***
	if other, ok := valToElement(v); ok ***REMOVED***
		return e.node == other.node
	***REMOVED*** else ***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func (e Element) GetElementsByClassName(name string) []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***e.rt, e.qsel.Find("." + name)***REMOVED***)
***REMOVED***

func (e Element) GetElementsByTagName(name string) []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***e.rt, e.qsel.Find(name)***REMOVED***)
***REMOVED***

func (e Element) QuerySelector(selector string) goja.Value ***REMOVED***
	return selToElement(Selection***REMOVED***e.rt, e.qsel.Find(selector)***REMOVED***)
***REMOVED***

func (e Element) QuerySelectorAll(selector string) []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***e.rt, e.qsel.Find(selector)***REMOVED***)
***REMOVED***

func (e Element) NodeName() string ***REMOVED***
	return goquery.NodeName(e.qsel)
***REMOVED***

func (e Element) FirstChild() goja.Value ***REMOVED***
	return nodeToElement(e, e.node.FirstChild)
***REMOVED***

func (e Element) LastChild() goja.Value ***REMOVED***
	return nodeToElement(e, e.node.LastChild)
***REMOVED***

func (e Element) FirstElementChild() goja.Value ***REMOVED***
	if child := e.qsel.Children().First(); child.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.rt, child.First()***REMOVED***)
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) LastElementChild() goja.Value ***REMOVED***
	if child := e.qsel.Children(); child.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.rt, child.Last()***REMOVED***)
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) PreviousSibling() goja.Value ***REMOVED***
	return nodeToElement(e, e.node.PrevSibling)
***REMOVED***

func (e Element) NextSibling() goja.Value ***REMOVED***
	return nodeToElement(e, e.node.NextSibling)
***REMOVED***

func (e Element) PreviousElementSibling() goja.Value ***REMOVED***
	if prev := e.qsel.Prev(); prev.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.rt, prev***REMOVED***)
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) NextElementSibling() goja.Value ***REMOVED***
	if next := e.qsel.Next(); next.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.rt, next***REMOVED***)
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) ParentNode() goja.Value ***REMOVED***
	if e.node.Parent != nil ***REMOVED***
		return nodeToElement(e, e.node.Parent)
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) ParentElement() goja.Value ***REMOVED***
	if prt := e.qsel.Parent(); prt.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.rt, prt***REMOVED***)
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) ChildNodes() []goja.Value ***REMOVED***
	return elemList(e.sel.Contents())
***REMOVED***

func (e Element) Children() []goja.Value ***REMOVED***
	return elemList(e.sel.Children())
***REMOVED***

func (e Element) childElementCount() int ***REMOVED***
	return e.sel.Children().Size()
***REMOVED***

func (e Element) ClassList() []string ***REMOVED***
	if clsName, exists := e.qsel.Attr("class"); exists ***REMOVED***
		return strings.Fields(clsName)
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

func (e Element) ClassName() goja.Value ***REMOVED***
	return e.sel.Attr("class")
***REMOVED***

func (e Element) Lang() goja.Value ***REMOVED***
	if attr := getHtmlAttr(e.node, "lang"); attr != nil && attr.Namespace == "" ***REMOVED***
		return e.rt.ToValue(attr.Val)
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) OwnerDocument() goja.Value ***REMOVED***
	if node := getOwnerDocNode(e.node); node != nil ***REMOVED***
		return nodeToElement(e, node)
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) NamespaceURI() string ***REMOVED***
	return e.node.Namespace
***REMOVED***

func (e Element) IsDefaultNamespace() bool ***REMOVED***
	// 	TODO namespace value of node always seems to be blank?
	return false
***REMOVED***

func getOwnerDocNode(node *gohtml.Node) *gohtml.Node ***REMOVED***
	for ; node != nil; node = node.Parent ***REMOVED***
		if node.Type == gohtml.DocumentNode ***REMOVED***
			return node
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (e Element) InnerHTML() goja.Value ***REMOVED***
	return e.sel.Html()
***REMOVED***

func (e Element) NodeType() goja.Value ***REMOVED***
	switch e.node.Type ***REMOVED***
	case gohtml.TextNode:
		return e.rt.ToValue(3)

	case gohtml.DocumentNode:
		return e.rt.ToValue(9)

	case gohtml.ElementNode:
		return e.rt.ToValue(1)

	case gohtml.CommentNode:
		return e.rt.ToValue(8)

	case gohtml.DoctypeNode:
		return e.rt.ToValue(10)

	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) NodeValue() goja.Value ***REMOVED***
	switch e.node.Type ***REMOVED***
	case gohtml.TextNode:
		return e.rt.ToValue(e.sel.Text())

	case gohtml.CommentNode:
		return e.rt.ToValue(e.sel.Text())

	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) Contains(v goja.Value) bool ***REMOVED***
	if other, ok := valToElement(v); ok ***REMOVED***
		// when testing if a node contains itself, jquery + goquery Contains() return true, JS return false
		return other.node == e.node || e.qsel.Contains(other.node)
	***REMOVED*** else ***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func (e Element) Matches(selector string) bool ***REMOVED***
	return e.qsel.Is(selector)
***REMOVED***

//helper methods
func getHtmlAttr(node *gohtml.Node, name string) *gohtml.Attribute ***REMOVED***
	for i := 0; i < len(node.Attr); i++ ***REMOVED***
		if node.Attr[i].Key == name ***REMOVED***
			return &node.Attr[i]
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func elemList(s Selection) (items []goja.Value) ***REMOVED***
	for i := 0; i < s.Size(); i++ ***REMOVED***
		items = append(items, selToElement(s.Eq(i)))
	***REMOVED***
	return items
***REMOVED***

func nodeToElement(e Element, node *gohtml.Node) goja.Value ***REMOVED***
	emptySel := e.qsel.Eq(e.qsel.Length())
	emptySel.Nodes = append(emptySel.Nodes, node)

	sel := Selection***REMOVED***e.rt, emptySel***REMOVED***

	return selToElement(sel)
***REMOVED***

func valToElementList(val goja.Value) (elems []*Element) ***REMOVED***
	vals := val.Export().([]goja.Value)
	for i := 0; i < len(vals); i++ ***REMOVED***
		if elem, ok := valToElement(vals[i]); ok ***REMOVED***
			elems = append(elems, elem)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func valToElement(v goja.Value) (*Element, bool) ***REMOVED***
	obj, ok := v.Export().(map[string]interface***REMOVED******REMOVED***)

	if !ok ***REMOVED***
		return nil, false
	***REMOVED***

	other, ok := obj["__elem__"]

	if !ok ***REMOVED***
		return nil, false
	***REMOVED***

	if elem, ok := other.(*Element); ok ***REMOVED***
		return elem, true
	***REMOVED*** else ***REMOVED***
		return nil, false
	***REMOVED***
***REMOVED***

func selToElement(sel Selection) goja.Value ***REMOVED***
	if sel.sel.Length() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED*** else if sel.sel.Length() > 1 ***REMOVED***
		sel = sel.First()
	***REMOVED***

	elem := sel.rt.NewObject()

	e := Element***REMOVED***&sel, sel.rt, sel.sel, sel.sel.Nodes[0]***REMOVED***

	proto, ok := initJsElem(sel.rt)
	if !ok ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	elem.Set("__proto__", proto)
	elem.Set("__elem__", sel.rt.ToValue(&e))

	return sel.rt.ToValue(elem)
***REMOVED***

func initJsElem(rt *goja.Runtime) (goja.Value, bool) ***REMOVED***
	if protoPrg == nil ***REMOVED***
		compileProtoElem()
	***REMOVED***

	obj, err := rt.RunProgram(protoPrg)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	return obj, true
***REMOVED***

func compileProtoElem() ***REMOVED***
	protoPrg = common.MustCompile("Element proto", `Object.freeze(***REMOVED***
	get id() ***REMOVED*** return this.__elem__.id(); ***REMOVED***,
	get nodeName() ***REMOVED*** return this.__elem__.nodeName(); ***REMOVED***,
	get nodeType() ***REMOVED*** return this.__elem__.nodeType(); ***REMOVED***,
	get nodeValue() ***REMOVED*** return this.__elem__.nodeValue(); ***REMOVED***,
	get innerHTML() ***REMOVED*** return this.__elem__.innerHTML(); ***REMOVED***,
	get textContent() ***REMOVED*** return this.__elem__.textContent(); ***REMOVED***,

	get attributes() ***REMOVED*** return this.__elem__.attributes(this); ***REMOVED***,

	get firstChild() ***REMOVED*** return this.__elem__.firstChild(); ***REMOVED***,
	get lastChild() ***REMOVED*** return this.__elem__.lastChild(); ***REMOVED***,
	get firstElementChild() ***REMOVED*** return this.__elem__.firstElementChild(); ***REMOVED***,
	get lastElementChild() ***REMOVED*** return this.__elem__.lastElementChild(); ***REMOVED***,

	get previousSibling() ***REMOVED*** return this.__elem__.previousSibling(); ***REMOVED***,
	get nextSibling() ***REMOVED*** return this.__elem__.nextSibling(); ***REMOVED***,

	get previousElementSibling() ***REMOVED*** return this.__elem__.previousElementSibling(); ***REMOVED***,
	get nextElementSibling() ***REMOVED*** return this.__elem__.nextElementSibling(); ***REMOVED***,

	get parentNode() ***REMOVED*** return this.__elem__.parentNode(); ***REMOVED***,
	get parentElement() ***REMOVED*** return this.__elem__.parentElement(); ***REMOVED***,

	get childNodes() ***REMOVED*** return this.__elem__.childNodes(); ***REMOVED***,
	get childElementCount() ***REMOVED*** return this.__elem__.childElementCount(); ***REMOVED***,
	get children() ***REMOVED*** return this.__elem__.children(); ***REMOVED***,

	get classList() ***REMOVED*** return this.__elem__.classList(); ***REMOVED***,
	get className() ***REMOVED*** return this.__elem__.className(); ***REMOVED***,

	get lang() ***REMOVED*** return this.__elem__.lang(); ***REMOVED***,
	get ownerDocument() ***REMOVED*** return this.__elem__.ownerDocument(); ***REMOVED***,
	get namespaceURI() ***REMOVED*** return this.__elem__.namespaceURI(); ***REMOVED***,


	toString: function() ***REMOVED*** return this.__elem__.toString(); ***REMOVED***,
	hasAttribute: function(name) ***REMOVED*** return this.__elem__.hasAttribute(name); ***REMOVED***,
	getAttribute: function(name) ***REMOVED*** return this.__elem__.getAttribute(name); ***REMOVED***,
	getAttributeNode: function(name) ***REMOVED*** return this.__elem__.getAttributeNode(this, name); ***REMOVED***,
	hasAttributes: function() ***REMOVED*** return this.__elem__.hasAttributes(); ***REMOVED***,
	hasChildNodes: function() ***REMOVED*** return this.__elem__.hasChildNodes(); ***REMOVED***,
	isSameNode: function(val) ***REMOVED*** return this.__elem__.isSameNode(val); ***REMOVED***,
	isEqualNode: function(val) ***REMOVED*** return this.__elem__.isEqualNode(val); ***REMOVED***,
	getElementsByClassName: function(val) ***REMOVED*** return this.__elem__.getElementsByClassName(val); ***REMOVED***,
	getElementsByTagName: function(val) ***REMOVED*** return this.__elem__.getElementsByTagName(val); ***REMOVED***,

	querySelector: function(val) ***REMOVED*** return this.__elem__.querySelector(val); ***REMOVED***,
	querySelectorAll: function(val) ***REMOVED*** return this.__elem__.querySelectorAll(val); ***REMOVED***,

	contains: function(node) ***REMOVED*** return this.__elem__.contains(node); ***REMOVED***
	matches: function(str) ***REMOVED*** return this.__elem__.matches(str); ***REMOVED***

***REMOVED***);
`, true)
***REMOVED***
