package goquery

import (
	"bytes"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var rxClassTrim = regexp.MustCompile("[\t\r\n]")

// Attr gets the specified attribute's value for the first element in the
// Selection. To get the value for each element individually, use a looping
// construct such as Each or Map method.
func (s *Selection) Attr(attrName string) (val string, exists bool) ***REMOVED***
	if len(s.Nodes) == 0 ***REMOVED***
		return
	***REMOVED***
	return getAttributeValue(attrName, s.Nodes[0])
***REMOVED***

// AttrOr works like Attr but returns default value if attribute is not present.
func (s *Selection) AttrOr(attrName, defaultValue string) string ***REMOVED***
	if len(s.Nodes) == 0 ***REMOVED***
		return defaultValue
	***REMOVED***

	val, exists := getAttributeValue(attrName, s.Nodes[0])
	if !exists ***REMOVED***
		return defaultValue
	***REMOVED***

	return val
***REMOVED***

// RemoveAttr removes the named attribute from each element in the set of matched elements.
func (s *Selection) RemoveAttr(attrName string) *Selection ***REMOVED***
	for _, n := range s.Nodes ***REMOVED***
		removeAttr(n, attrName)
	***REMOVED***

	return s
***REMOVED***

// SetAttr sets the given attribute on each element in the set of matched elements.
func (s *Selection) SetAttr(attrName, val string) *Selection ***REMOVED***
	for _, n := range s.Nodes ***REMOVED***
		attr := getAttributePtr(attrName, n)
		if attr == nil ***REMOVED***
			n.Attr = append(n.Attr, html.Attribute***REMOVED***Key: attrName, Val: val***REMOVED***)
		***REMOVED*** else ***REMOVED***
			attr.Val = val
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***

// Text gets the combined text contents of each element in the set of matched
// elements, including their descendants.
func (s *Selection) Text() string ***REMOVED***
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) ***REMOVED***
		if n.Type == html.TextNode ***REMOVED***
			// Keep newlines and spaces, like jQuery
			buf.WriteString(n.Data)
		***REMOVED***
		if n.FirstChild != nil ***REMOVED***
			for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
				f(c)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, n := range s.Nodes ***REMOVED***
		f(n)
	***REMOVED***

	return buf.String()
***REMOVED***

// Size is an alias for Length.
func (s *Selection) Size() int ***REMOVED***
	return s.Length()
***REMOVED***

// Length returns the number of elements in the Selection object.
func (s *Selection) Length() int ***REMOVED***
	return len(s.Nodes)
***REMOVED***

// Html gets the HTML contents of the first element in the set of matched
// elements. It includes text and comment nodes.
func (s *Selection) Html() (ret string, e error) ***REMOVED***
	// Since there is no .innerHtml, the HTML content must be re-created from
	// the nodes using html.Render.
	var buf bytes.Buffer

	if len(s.Nodes) > 0 ***REMOVED***
		for c := s.Nodes[0].FirstChild; c != nil; c = c.NextSibling ***REMOVED***
			e = html.Render(&buf, c)
			if e != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		ret = buf.String()
	***REMOVED***

	return
***REMOVED***

// AddClass adds the given class(es) to each element in the set of matched elements.
// Multiple class names can be specified, separated by a space or via multiple arguments.
func (s *Selection) AddClass(class ...string) *Selection ***REMOVED***
	classStr := strings.TrimSpace(strings.Join(class, " "))

	if classStr == "" ***REMOVED***
		return s
	***REMOVED***

	tcls := getClassesSlice(classStr)
	for _, n := range s.Nodes ***REMOVED***
		curClasses, attr := getClassesAndAttr(n, true)
		for _, newClass := range tcls ***REMOVED***
			if !strings.Contains(curClasses, " "+newClass+" ") ***REMOVED***
				curClasses += newClass + " "
			***REMOVED***
		***REMOVED***

		setClasses(n, attr, curClasses)
	***REMOVED***

	return s
***REMOVED***

// HasClass determines whether any of the matched elements are assigned the
// given class.
func (s *Selection) HasClass(class string) bool ***REMOVED***
	class = " " + class + " "
	for _, n := range s.Nodes ***REMOVED***
		classes, _ := getClassesAndAttr(n, false)
		if strings.Contains(classes, class) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// RemoveClass removes the given class(es) from each element in the set of matched elements.
// Multiple class names can be specified, separated by a space or via multiple arguments.
// If no class name is provided, all classes are removed.
func (s *Selection) RemoveClass(class ...string) *Selection ***REMOVED***
	var rclasses []string

	classStr := strings.TrimSpace(strings.Join(class, " "))
	remove := classStr == ""

	if !remove ***REMOVED***
		rclasses = getClassesSlice(classStr)
	***REMOVED***

	for _, n := range s.Nodes ***REMOVED***
		if remove ***REMOVED***
			removeAttr(n, "class")
		***REMOVED*** else ***REMOVED***
			classes, attr := getClassesAndAttr(n, true)
			for _, rcl := range rclasses ***REMOVED***
				classes = strings.Replace(classes, " "+rcl+" ", " ", -1)
			***REMOVED***

			setClasses(n, attr, classes)
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***

// ToggleClass adds or removes the given class(es) for each element in the set of matched elements.
// Multiple class names can be specified, separated by a space or via multiple arguments.
func (s *Selection) ToggleClass(class ...string) *Selection ***REMOVED***
	classStr := strings.TrimSpace(strings.Join(class, " "))

	if classStr == "" ***REMOVED***
		return s
	***REMOVED***

	tcls := getClassesSlice(classStr)

	for _, n := range s.Nodes ***REMOVED***
		classes, attr := getClassesAndAttr(n, true)
		for _, tcl := range tcls ***REMOVED***
			if strings.Contains(classes, " "+tcl+" ") ***REMOVED***
				classes = strings.Replace(classes, " "+tcl+" ", " ", -1)
			***REMOVED*** else ***REMOVED***
				classes += tcl + " "
			***REMOVED***
		***REMOVED***

		setClasses(n, attr, classes)
	***REMOVED***

	return s
***REMOVED***

func getAttributePtr(attrName string, n *html.Node) *html.Attribute ***REMOVED***
	if n == nil ***REMOVED***
		return nil
	***REMOVED***

	for i, a := range n.Attr ***REMOVED***
		if a.Key == attrName ***REMOVED***
			return &n.Attr[i]
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Private function to get the specified attribute's value from a node.
func getAttributeValue(attrName string, n *html.Node) (val string, exists bool) ***REMOVED***
	if a := getAttributePtr(attrName, n); a != nil ***REMOVED***
		val = a.Val
		exists = true
	***REMOVED***
	return
***REMOVED***

// Get and normalize the "class" attribute from the node.
func getClassesAndAttr(n *html.Node, create bool) (classes string, attr *html.Attribute) ***REMOVED***
	// Applies only to element nodes
	if n.Type == html.ElementNode ***REMOVED***
		attr = getAttributePtr("class", n)
		if attr == nil && create ***REMOVED***
			n.Attr = append(n.Attr, html.Attribute***REMOVED***
				Key: "class",
				Val: "",
			***REMOVED***)
			attr = &n.Attr[len(n.Attr)-1]
		***REMOVED***
	***REMOVED***

	if attr == nil ***REMOVED***
		classes = " "
	***REMOVED*** else ***REMOVED***
		classes = rxClassTrim.ReplaceAllString(" "+attr.Val+" ", " ")
	***REMOVED***

	return
***REMOVED***

func getClassesSlice(classes string) []string ***REMOVED***
	return strings.Split(rxClassTrim.ReplaceAllString(" "+classes+" ", " "), " ")
***REMOVED***

func removeAttr(n *html.Node, attrName string) ***REMOVED***
	for i, a := range n.Attr ***REMOVED***
		if a.Key == attrName ***REMOVED***
			n.Attr[i], n.Attr[len(n.Attr)-1], n.Attr =
				n.Attr[len(n.Attr)-1], html.Attribute***REMOVED******REMOVED***, n.Attr[:len(n.Attr)-1]
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func setClasses(n *html.Node, attr *html.Attribute, classes string) ***REMOVED***
	classes = strings.TrimSpace(classes)
	if classes == "" ***REMOVED***
		removeAttr(n, "class")
		return
	***REMOVED***

	attr.Val = classes
***REMOVED***
