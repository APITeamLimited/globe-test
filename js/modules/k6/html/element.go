package html

import (
	"fmt"
	"net/url"
	"strconv"
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
	OwnerElement *Element `json:"owner_element"`
	Name         string   `json:"name"`
	nsPrefix     string
	Value        string `json:"value"`
***REMOVED***

func (e Element) attrAsString(name string) string ***REMOVED***
	val, exists := e.sel.sel.Attr(name)
	if !exists ***REMOVED***
		return ""
	***REMOVED***
	return val
***REMOVED***

func (e Element) resolveURL(val string) (*url.URL, bool) ***REMOVED***
	baseURL, err := url.Parse(e.sel.URL)
	if err != nil ***REMOVED***
		return nil, false
	***REMOVED***

	addURL, err := url.Parse(val)
	if err != nil ***REMOVED***
		return nil, false
	***REMOVED***

	return baseURL.ResolveReference(addURL), true
***REMOVED***

func (e Element) attrAsURL(name string) (*url.URL, bool) ***REMOVED***
	val, exists := e.sel.sel.Attr(name)
	if !exists ***REMOVED***
		return nil, false
	***REMOVED***

	return e.resolveURL(val)
***REMOVED***

func (e Element) attrAsURLString(name string, defaultWhenNoAttr string) string ***REMOVED***
	if e.sel.URL == "" ***REMOVED***
		return e.attrAsString(name)
	***REMOVED***

	url, ok := e.attrAsURL(name)
	if !ok ***REMOVED***
		return defaultWhenNoAttr
	***REMOVED***

	return url.String()
***REMOVED***

func (e Element) attrAsInt(name string, defaultVal int) int ***REMOVED***
	strVal, exists := e.sel.sel.Attr(name)
	if !exists ***REMOVED***
		return defaultVal
	***REMOVED***

	intVal, err := strconv.Atoi(strVal)
	if err != nil ***REMOVED***
		return defaultVal
	***REMOVED***

	return intVal
***REMOVED***

func (e Element) attrIsPresent(name string) bool ***REMOVED***
	_, exists := e.sel.sel.Attr(name)
	return exists
***REMOVED***

func (e Element) ownerFormSel() (*goquery.Selection, bool) ***REMOVED***
	prtForm := e.sel.sel.Closest("form")
	if prtForm.Length() > 0 ***REMOVED***
		return prtForm, true
	***REMOVED***

	formId := e.attrAsString("form")
	if formId == "" ***REMOVED***
		return nil, false
	***REMOVED***

	findForm := e.sel.sel.Parents().Last().Find("#" + formId)
	if findForm.Length() == 0 ***REMOVED***
		return nil, false
	***REMOVED***

	return findForm, true
***REMOVED***

func (e Element) ownerFormVal() goja.Value ***REMOVED***
	formSel, exists := e.ownerFormSel()
	if !exists ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return selToElement(Selection***REMOVED***e.sel.rt, formSel.Eq(0), e.sel.URL***REMOVED***)
***REMOVED***

func (e Element) elemLabels() []goja.Value ***REMOVED***
	wrapperLbl := e.sel.sel.Closest("label")

	id := e.attrAsString("id")
	if id == "" ***REMOVED***
		return elemList(Selection***REMOVED***e.sel.rt, wrapperLbl, e.sel.URL***REMOVED***)
	***REMOVED***

	idLbl := e.sel.sel.Parents().Last().Find("label[for=\"" + id + "\"]")
	if idLbl.Size() == 0 ***REMOVED***
		return elemList(Selection***REMOVED***e.sel.rt, wrapperLbl, e.sel.URL***REMOVED***)
	***REMOVED***

	allLbls := wrapperLbl.AddSelection(idLbl)

	return elemList(Selection***REMOVED***e.sel.rt, allLbls, e.sel.URL***REMOVED***)
***REMOVED***

func (e Element) splitAttr(attrName string) []string ***REMOVED***
	attr := e.attrAsString(attrName)

	if attr == "" ***REMOVED***
		return make([]string, 0)
	***REMOVED***

	return strings.Split(attr, " ")
***REMOVED***

func (e Element) idOrNameAttr() (string, bool) ***REMOVED***
	if id, exists := e.sel.sel.Attr("id"); exists ***REMOVED***
		return id, true
	***REMOVED***

	if name, exists := e.sel.sel.Attr("id"); exists ***REMOVED***
		return name, true
	***REMOVED***

	return "", false
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
		return e.sel.rt.ToValue(Attribute***REMOVED***&e, attr.Key, attr.Namespace, attr.Val***REMOVED***)
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
		attrs[attr.Key] = Attribute***REMOVED***&e, attr.Key, attr.Namespace, attr.Val***REMOVED***
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
	return elemList(Selection***REMOVED***e.sel.rt, e.sel.sel.Find("." + name), e.sel.URL***REMOVED***)
***REMOVED***

func (e Element) GetElementsByTagName(name string) []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***e.sel.rt, e.sel.sel.Find(name), e.sel.URL***REMOVED***)
***REMOVED***

func (e Element) QuerySelector(selector string) goja.Value ***REMOVED***
	return selToElement(Selection***REMOVED***e.sel.rt, e.sel.sel.Find(selector), e.sel.URL***REMOVED***)
***REMOVED***

func (e Element) QuerySelectorAll(selector string) []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***e.sel.rt, e.sel.sel.Find(selector), e.sel.URL***REMOVED***)
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
		return selToElement(Selection***REMOVED***e.sel.rt, child.First(), e.sel.URL***REMOVED***)
	***REMOVED***

	return goja.Undefined()
***REMOVED***

func (e Element) LastElementChild() goja.Value ***REMOVED***
	if child := e.sel.sel.Children(); child.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.sel.rt, child.Last(), e.sel.URL***REMOVED***)
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
		return selToElement(Selection***REMOVED***e.sel.rt, prev, e.sel.URL***REMOVED***)
	***REMOVED***

	return goja.Undefined()
***REMOVED***

func (e Element) NextElementSibling() goja.Value ***REMOVED***
	if next := e.sel.sel.Next(); next.Length() > 0 ***REMOVED***
		return selToElement(Selection***REMOVED***e.sel.rt, next, e.sel.URL***REMOVED***)
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
		return selToElement(Selection***REMOVED***e.sel.rt, prt, e.sel.URL***REMOVED***)
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
