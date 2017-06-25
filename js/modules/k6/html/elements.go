package html

import (
	"net"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
)

//go:generate go run gen/main.go

var defaultPorts = map[string]string***REMOVED***
	"http":  "80",
	"https": "443",
	"ftp":   "21",
***REMOVED***

const (
	AnchorTagName = "a"
	AreaTagName   = "area"
	BaseTagName   = "base"
	ButtonTagName = "button"
)

type HrefElement struct***REMOVED*** Element ***REMOVED***
type FormFieldElement struct***REMOVED*** Element ***REMOVED***

type AnchorElement struct***REMOVED*** HrefElement ***REMOVED***
type AreaElement struct***REMOVED*** HrefElement ***REMOVED***

type BaseElement struct***REMOVED*** Element ***REMOVED***
type ButtonElement struct***REMOVED*** FormFieldElement ***REMOVED***

func (h HrefElement) hrefURL() *url.URL ***REMOVED***
	url, err := url.Parse(h.attrAsString("href"))
	if err != nil ***REMOVED***
		url, _ = url.Parse("")
	***REMOVED***

	return url
***REMOVED***

func (h HrefElement) Hash() string ***REMOVED***
	return "#" + h.hrefURL().Fragment
***REMOVED***

func (h HrefElement) Host() string ***REMOVED***
	url := h.hrefURL()
	hostAndPort := url.Host

	host, port, err := net.SplitHostPort(hostAndPort)
	if err != nil ***REMOVED***
		return hostAndPort
	***REMOVED***

	defaultPort := defaultPorts[url.Scheme]
	if defaultPort != "" && port == defaultPort ***REMOVED***
		return strings.TrimSuffix(host, ":"+defaultPort)
	***REMOVED***

	return hostAndPort
***REMOVED***

func (h HrefElement) Hostname() string ***REMOVED***
	host, _, err := net.SplitHostPort(h.hrefURL().Host)
	if err != nil ***REMOVED***
		return h.hrefURL().Host
	***REMOVED***
	return host
***REMOVED***

func (h HrefElement) Port() string ***REMOVED***
	_, port, err := net.SplitHostPort(h.hrefURL().Host)
	if err != nil ***REMOVED***
		return ""
	***REMOVED***
	return port
***REMOVED***

func (h HrefElement) Username() string ***REMOVED***
	user := h.hrefURL().User
	if user == nil ***REMOVED***
		return ""
	***REMOVED***
	return user.Username()
***REMOVED***

func (h HrefElement) Password() goja.Value ***REMOVED***
	user := h.hrefURL().User
	if user == nil ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	pwd, defined := user.Password()
	if !defined ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	return h.sel.rt.ToValue(pwd)
***REMOVED***

func (h HrefElement) Origin() string ***REMOVED***
	href := h.hrefURL()

	if href.Scheme == "file" ***REMOVED***
		return h.Href()
	***REMOVED***

	return href.Scheme + "://" + href.Host
***REMOVED***

func (h HrefElement) Pathname() string ***REMOVED***
	return h.hrefURL().Path
***REMOVED***

func (h HrefElement) Protocol() string ***REMOVED***
	return h.hrefURL().Scheme
***REMOVED***

func (h HrefElement) RelList() []string ***REMOVED***
	rel := h.attrAsString("rel")

	if rel == "" ***REMOVED***
		return make([]string, 0)
	***REMOVED***

	return strings.Split(rel, " ")
***REMOVED***

func (h HrefElement) Search() string ***REMOVED***
	q := h.hrefURL().RawQuery
	if q == "" ***REMOVED***
		return q
	***REMOVED***
	return "?" + q
***REMOVED***

func (h HrefElement) Text() string ***REMOVED***
	return h.TextContent()
***REMOVED***

func (f FormFieldElement) Form() goja.Value ***REMOVED***
	formSel, exists := f.ownerFormSel()
	if !exists ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return selToElement(Selection***REMOVED***f.sel.rt, formSel***REMOVED***)
***REMOVED***

// Used by the formAction, formMethod, formTarget and formEnctype methods of Button and Input elements
// Attempts to read attribute "form" + attrName on the current element or attrName on the owning form element
func (f FormFieldElement) formAttrOrElemOverride(attrName string) string ***REMOVED***
	if elemAttr, exists := f.sel.sel.Attr("form" + attrName); exists ***REMOVED***
		return elemAttr
	***REMOVED***

	formSel, exists := f.ownerFormSel()
	if !exists ***REMOVED***
		return ""
	***REMOVED***

	formAttr, exists := formSel.Attr(attrName)
	if !exists ***REMOVED***
		return ""
	***REMOVED***

	return formAttr
***REMOVED***

// Used by the formAction, formMethod, formTarget and formEnctype methods of Button and Input elements
// Attempts to read attribute "form" + attrName on the current element or attrName on the owning form element
func (f FormFieldElement) formOrElemAttrString(attrName string) string ***REMOVED***
	if elemAttr, exists := f.sel.sel.Attr("form" + attrName); exists ***REMOVED***
		return elemAttr
	***REMOVED***

	formSel, exists := f.ownerFormSel()
	if !exists ***REMOVED***
		return ""
	***REMOVED***

	formAttr, exists := formSel.Attr(attrName)
	if !exists ***REMOVED***
		return ""
	***REMOVED***

	return formAttr
***REMOVED***

func (f FormFieldElement) formOrElemAttrPresent(attrName string) bool ***REMOVED***
	if _, exists := f.sel.sel.Attr("form" + attrName); exists ***REMOVED***
		return true
	***REMOVED***

	formSel, exists := f.ownerFormSel()
	if !exists ***REMOVED***
		return false
	***REMOVED***

	_, exists = formSel.Attr(attrName)
	return exists
***REMOVED***

func (f FormFieldElement) FormAction() string ***REMOVED***
	return f.formOrElemAttrString("action")
***REMOVED***

func (f FormFieldElement) FormEnctype() string ***REMOVED***
	return f.formOrElemAttrString("enctype")
***REMOVED***

func (f FormFieldElement) FormMethod() string ***REMOVED***
	return f.formOrElemAttrString("method")
***REMOVED***

func (f FormFieldElement) FormNoValidate() bool ***REMOVED***
	return f.formOrElemAttrPresent("novalidate")
***REMOVED***

func (f FormFieldElement) FormTarget() string ***REMOVED***
	return f.formOrElemAttrString("target")
***REMOVED***

func (f FormFieldElement) elemLabels() (items []goja.Value) ***REMOVED***
	wrapperLbl := f.sel.sel.Closest("label")

	id := f.attrAsString("id")
	if id == "" ***REMOVED***
		return elemList(Selection***REMOVED***f.sel.rt, wrapperLbl***REMOVED***)
	***REMOVED***

	idLbl := f.sel.sel.Parents().Last().Find("label[for=\"" + id + "\"]")
	if idLbl.Size() == 0 ***REMOVED***
		return elemList(Selection***REMOVED***f.sel.rt, wrapperLbl***REMOVED***)
	***REMOVED***

	allLbls := wrapperLbl.AddSelection(idLbl)

	return elemList(Selection***REMOVED***f.sel.rt, allLbls***REMOVED***)
***REMOVED***

func (f FormFieldElement) Labels() (items []goja.Value) ***REMOVED***
	return f.elemLabels()
***REMOVED***

func (b ButtonElement) Type() string ***REMOVED***
	switch b.attrAsString("type") ***REMOVED***
	case "button":
		return "button"
	case "menu":
		return "menu"
	case "reset":
		return "reset"
	default:
		return "submit"
	***REMOVED***
***REMOVED***
