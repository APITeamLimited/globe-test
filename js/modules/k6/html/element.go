package html

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	gohtml "golang.org/x/net/html"
)

const (
	ElementNode  = 1
	TextNode     = 3
	CommentNode  = 8
	DocumentNode = 9
	DoctypeNode  = 10
)

type Element struct ***REMOVED***
	node *gohtml.Node
	sel  *Selection
***REMOVED***

type Attribute struct ***REMOVED***
	Name         string
	nsPrefix     string
	OwnerElement *Element
	Value        string
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

func (e Element) attrAsString(name string) string ***REMOVED***
	val, exists := e.sel.sel.Attr(name)
	if !exists ***REMOVED***
		return ""
	***REMOVED***
	return val
***REMOVED***

func (e Element) attrIsPresent(name string) bool ***REMOVED***
	_, exists := e.sel.sel.Attr(name)
	return exists
***REMOVED***

func (e Element) GetAttribute(name string) goja.Value ***REMOVED***
	return e.sel.Attr(name)
***REMOVED***

func (e Element) GetAttributeNode(name string) goja.Value ***REMOVED***
	if attr := getHtmlAttr(e.node, name); attr != nil ***REMOVED***
		return e.sel.rt.ToValue(Attribute***REMOVED***attr.Key, attr.Namespace, &e, attr.Val***REMOVED***)
	***REMOVED***

	return goja.Undefined()
***REMOVED***

func (e Element) HasAttribute(name string) bool ***REMOVED***
	_, exists := e.sel.sel.Attr(name)
	return exists
***REMOVED***

func (e Element) HasAttributes() bool ***REMOVED***
	return e.sel.sel.Length() > 0 && len(e.node.Attr) > 0
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
	if e.sel.sel.Length() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	if e.node.Type == gohtml.ElementNode ***REMOVED***
		return e.sel.rt.ToValue("[object html.Node]")
	***REMOVED***

	return e.sel.rt.ToValue(fmt.Sprintf("[object %s]", e.NodeName()))
***REMOVED***

func (e Element) HasChildNodes() bool ***REMOVED***
	return e.sel.sel.Length() > 0 && e.node.FirstChild != nil
***REMOVED***

func (e Element) TextContent() string ***REMOVED***
	return e.sel.sel.Text()
***REMOVED***

func (e Element) Id() string ***REMOVED***
	return e.attrAsString("id")
***REMOVED***

func (e Element) IsEqualNode(v goja.Value) bool ***REMOVED***
	if other, ok := v.Export().(Element); ok ***REMOVED***
		htmlA, errA := e.sel.sel.Html()
		htmlB, errB := other.sel.sel.Html()

		return errA == nil && errB == nil && htmlA == htmlB
	***REMOVED***

	return false
***REMOVED***

func (e Element) IsSameNode(v goja.Value) bool ***REMOVED***
	if other, ok := v.Export().(Element); ok ***REMOVED***
		return e.node == other.node
	***REMOVED***

	return false
***REMOVED***

func (e Element) GetElementsByClassName(name string) []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***e.sel.rt, e.sel.sel.Find("." + name)***REMOVED***)
***REMOVED***

func (e Element) GetElementsByTagName(name string) []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***e.sel.rt, e.sel.sel.Find(name)***REMOVED***)
***REMOVED***

func (e Element) QuerySelector(selector string) goja.Value ***REMOVED***
	return selToElement(Selection***REMOVED***e.sel.rt, e.sel.sel.Find(selector)***REMOVED***)
***REMOVED***

func (e Element) QuerySelectorAll(selector string) []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***e.sel.rt, e.sel.sel.Find(selector)***REMOVED***)
***REMOVED***

func (e Element) NodeName() string ***REMOVED***
	return goquery.NodeName(e.sel.sel)
***REMOVED***

func (e Element) FirstChild() goja.Value ***REMOVED***
	return nodeToElement(e, e.node.FirstChild)
***REMOVED***

func (e Element) LastChild() goja.Value ***REMOVED***
	return nodeToElement(e, e.node.LastChild)
***REMOVED***

func (e Element) FirstElementChild() goja.Value ***REMOVED***
	if child := e.sel.sel.Children().First(); child.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.sel.rt, child.First()***REMOVED***)
	***REMOVED***

	return goja.Undefined()
***REMOVED***

func (e Element) LastElementChild() goja.Value ***REMOVED***
	if child := e.sel.sel.Children(); child.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.sel.rt, child.Last()***REMOVED***)
	***REMOVED***

	return goja.Undefined()
***REMOVED***

func (e Element) PreviousSibling() goja.Value ***REMOVED***
	return nodeToElement(e, e.node.PrevSibling)
***REMOVED***

func (e Element) NextSibling() goja.Value ***REMOVED***
	return nodeToElement(e, e.node.NextSibling)
***REMOVED***

func (e Element) PreviousElementSibling() goja.Value ***REMOVED***
	if prev := e.sel.sel.Prev(); prev.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.sel.rt, prev***REMOVED***)
	***REMOVED***

	return goja.Undefined()
***REMOVED***

func (e Element) NextElementSibling() goja.Value ***REMOVED***
	if next := e.sel.sel.Next(); next.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.sel.rt, next***REMOVED***)
	***REMOVED***

	return goja.Undefined()
***REMOVED***

func (e Element) ParentNode() goja.Value ***REMOVED***
	if e.node.Parent != nil ***REMOVED***
		return nodeToElement(e, e.node.Parent)
	***REMOVED***

	return goja.Undefined()
***REMOVED***

func (e Element) ParentElement() goja.Value ***REMOVED***
	if prt := e.sel.sel.Parent(); prt.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.sel.rt, prt***REMOVED***)
	***REMOVED***

	return goja.Undefined()
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
	if clsName, exists := e.sel.sel.Attr("class"); exists ***REMOVED***
		return strings.Fields(clsName)
	***REMOVED***

	return nil
***REMOVED***

func (e Element) ClassName() goja.Value ***REMOVED***
	return e.sel.Attr("class")
***REMOVED***

func (e Element) Lang() goja.Value ***REMOVED***
	if attr := getHtmlAttr(e.node, "lang"); attr != nil && attr.Namespace == "" ***REMOVED***
		return e.sel.rt.ToValue(attr.Val)
	***REMOVED***

	return goja.Undefined()
***REMOVED***

func (e Element) OwnerDocument() goja.Value ***REMOVED***
	if node := getOwnerDocNode(e.node); node != nil ***REMOVED***
		return nodeToElement(e, node)
	***REMOVED***

	return goja.Undefined()
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
		return e.sel.rt.ToValue(TextNode)

	case gohtml.DocumentNode:
		return e.sel.rt.ToValue(DocumentNode)

	case gohtml.ElementNode:
		return e.sel.rt.ToValue(ElementNode)

	case gohtml.CommentNode:
		return e.sel.rt.ToValue(CommentNode)

	case gohtml.DoctypeNode:
		return e.sel.rt.ToValue(DoctypeNode)

	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) NodeValue() goja.Value ***REMOVED***
	switch e.node.Type ***REMOVED***
	case gohtml.TextNode:
		return e.sel.rt.ToValue(e.sel.Text())

	case gohtml.CommentNode:
		return e.sel.rt.ToValue(e.sel.Text())

	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (e Element) Contains(v goja.Value) bool ***REMOVED***
	if other, ok := v.Export().(Element); ok ***REMOVED***
		// When testing if a node contains itself, jquery's + goquery's version of Contains() return true while the DOM API returns false.
		return other.node != e.node && e.sel.sel.Contains(other.node)
	***REMOVED***

	return false
***REMOVED***

func (e Element) Matches(selector string) bool ***REMOVED***
	return e.sel.sel.Is(selector)
***REMOVED***
