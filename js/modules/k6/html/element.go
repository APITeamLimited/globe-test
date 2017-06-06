package html

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	gohtml "golang.org/x/net/html"
)

type Element struct ***REMOVED***
	sel  *Selection
	rt   *goja.Runtime
	qsel *goquery.Selection
	node *gohtml.Node
***REMOVED***

type Attribute struct ***REMOVED***
	Name         string
	nsPrefix     string
	OwnerElement *Element
	Value        string
***REMOVED***

func namespaceURI(prefix string) string ***REMOVED***
	switch prefix ***REMOVED***
	case "svg":
		return "http://www.w3.org/2000/svg"
	case "math":
		return "http://www.w3.org/1998/Math/MathML"
	default:
		return "http://www.w3.org/1999/xhtml"
	***REMOVED***
***REMOVED***

func (a Attribute) Prefix() string ***REMOVED***
	return a.nsPrefix
***REMOVED***

func (a Attribute) NamespaceURI() string ***REMOVED***
	return namespaceURI(a.nsPrefix)
***REMOVED***

func (a Attribute) LocalName() string ***REMOVED***
	return a.Name
***REMOVED***

func (e Element) GetAttribute(name string) goja.Value ***REMOVED***
	return e.sel.Attr(name)
***REMOVED***

func (e Element) GetAttributeNode(name string) goja.Value ***REMOVED***
	if attr := getHtmlAttr(e.node, name); attr != nil ***REMOVED***
		return e.rt.ToValue(Attribute***REMOVED***attr.Key, attr.Namespace, &e, attr.Val***REMOVED***)
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

func (e Element) Attributes() map[string]Attribute ***REMOVED***
	attrs := make(map[string]Attribute)
	for i := 0; i < len(e.node.Attr); i++ ***REMOVED***
		attr := e.node.Attr[i]
		attrs[attr.Key] = Attribute***REMOVED***attr.Key, attr.Namespace, &e, attr.Val***REMOVED***
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
	if other, ok := v.Export().(Element); ok ***REMOVED***
		htmlA, errA := e.qsel.Html()
		htmlB, errB := other.qsel.Html()

		return errA == nil && errB == nil && htmlA == htmlB
	***REMOVED*** else ***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func (e Element) IsSameNode(v goja.Value) bool ***REMOVED***
	if other, ok := v.Export().(Element); ok ***REMOVED***
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

func (e Element) ChildElementCount() int ***REMOVED***
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
	return namespaceURI(e.node.Namespace)
***REMOVED***

func (e Element) IsDefaultNamespace() bool ***REMOVED***
	return e.node.Namespace == ""
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
	if other, ok := v.Export().(Element); ok ***REMOVED***
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

func valToElementList(val goja.Value) (elems []Element) ***REMOVED***
	vals := val.Export().([]goja.Value)
	for i := 0; i < len(vals); i++ ***REMOVED***
		elems = append(elems, vals[i].Export().(Element))
	***REMOVED***
	return
***REMOVED***

func selToElement(sel Selection) goja.Value ***REMOVED***
	if sel.sel.Length() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED*** else if sel.sel.Length() > 1 ***REMOVED***
		sel = sel.First()
	***REMOVED***

	elem := Element***REMOVED***&sel, sel.rt, sel.sel, sel.sel.Nodes[0]***REMOVED***

	return sel.rt.ToValue(elem)
***REMOVED***
