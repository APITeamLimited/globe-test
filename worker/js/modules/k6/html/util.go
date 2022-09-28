package html

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/serenize/snaker"
	gohtml "golang.org/x/net/html"
)

func attrToProperty(s string) string ***REMOVED***
	if idx := strings.Index(s, "-"); idx != -1 ***REMOVED***
		return s[0:idx] + snaker.SnakeToCamel(strings.ReplaceAll(s[idx+1:], "-", "_"))
	***REMOVED***
	return s
***REMOVED***

func propertyToAttr(attrName string) string ***REMOVED***
	return strings.ReplaceAll(snaker.CamelToSnake(attrName), "_", "-")
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

func valueOrHTML(s *goquery.Selection) string ***REMOVED***
	if val, exists := s.Attr("value"); exists ***REMOVED***
		return val
	***REMOVED***

	if val, err := s.Html(); err == nil ***REMOVED***
		return val
	***REMOVED***

	return ""
***REMOVED***

func getHTMLAttr(node *gohtml.Node, name string) *gohtml.Attribute ***REMOVED***
	for i := 0; i < len(node.Attr); i++ ***REMOVED***
		if node.Attr[i].Key == name ***REMOVED***
			return &node.Attr[i]
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func elemList(s Selection) (items []goja.Value) ***REMOVED***
	items = make([]goja.Value, s.Size())
	for i := 0; i < s.Size(); i++ ***REMOVED***
		items[i] = selToElement(s.Eq(i))
	***REMOVED***
	return items
***REMOVED***

func nodeToElement(e Element, node *gohtml.Node) goja.Value ***REMOVED***
	// Goquery does not expose a way to build a goquery.Selection with an arbitrary html.Node.
	// Workaround by adding a node to an empty Selection
	emptySel := e.sel.emptySelection()
	emptySel.sel.Nodes = append(emptySel.sel.Nodes, node)

	return selToElement(emptySel)
***REMOVED***

// Try to read numeric values in data- attributes.
// Return numeric value when the representation is unchanged by conversion to float and back.
// Other potentially numeric values (ie "101.00" "1E02") remain as strings.
func toNumeric(val string) (float64, bool) ***REMOVED***
	if fltVal, err := strconv.ParseFloat(val, 64); err != nil ***REMOVED***
		return 0, false
	***REMOVED*** else if repr := strconv.FormatFloat(fltVal, 'f', -1, 64); repr == val ***REMOVED***
		return fltVal, true
	***REMOVED*** else ***REMOVED***
		return 0, false
	***REMOVED***
***REMOVED***

func convertDataAttrVal(val string) interface***REMOVED******REMOVED*** ***REMOVED***
	if len(val) == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	if val[0] == '***REMOVED***' || val[0] == '[' ***REMOVED***
		var subdata interface***REMOVED******REMOVED***

		err := json.Unmarshal([]byte(val), &subdata)
		if err == nil ***REMOVED***
			return subdata
		***REMOVED***
		return val
	***REMOVED***

	switch val ***REMOVED***
	case "true":
		return true

	case "false":
		return false

	case "null":
		return goja.Undefined()

	case "undefined":
		return goja.Undefined()

	default:
		if fltVal, isOk := toNumeric(val); isOk ***REMOVED***
			return fltVal
		***REMOVED***
		return val
	***REMOVED***
***REMOVED***
