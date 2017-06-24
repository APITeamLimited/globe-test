package html

import (
	"net"
	"net/url"
	"strings"

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
type AnchorElement struct***REMOVED*** HrefElement ***REMOVED***
type AreaElement struct***REMOVED*** HrefElement ***REMOVED***

type BaseElement struct***REMOVED*** Element ***REMOVED***
type ButtonElement struct***REMOVED*** Element ***REMOVED***

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

func (h HrefElement) Rel() string ***REMOVED***
	return h.attrAsString("rel")
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

func (h HrefElement) Target() string ***REMOVED***
	return h.attrAsString("target")
***REMOVED***

func (h HrefElement) Text() string ***REMOVED***
	return h.TextContent()
***REMOVED***

func (h HrefElement) Type() string ***REMOVED***
	return h.attrAsString("type")
***REMOVED***

func (h HrefElement) AccessKey() string ***REMOVED***
	return h.attrAsString("accesskey")
***REMOVED***

func (h HrefElement) HrefLang() string ***REMOVED***
	return h.attrAsString("hreflang")
***REMOVED***

func (h HrefElement) Media() string ***REMOVED***
	return h.attrAsString("media")
***REMOVED***

func (h HrefElement) ToString() string ***REMOVED***
	return h.attrAsString("href")
***REMOVED***

func (h HrefElement) Href() string ***REMOVED***
	return h.attrAsString("href")
***REMOVED***

func (h BaseElement) Href() string ***REMOVED***
	return h.attrAsString("href")
***REMOVED***

func (h BaseElement) Target() string ***REMOVED***
	return h.attrAsString("target")
***REMOVED***

func (b ButtonElement) AccessKey() string ***REMOVED***
	return b.attrAsString("accesskey")
***REMOVED***

func (b ButtonElement) Autofocus() bool ***REMOVED***
	return b.attrIsPresent("autofocus")
***REMOVED***

func (b ButtonElement) Disabled() bool ***REMOVED***
	return b.attrIsPresent("disabled")
***REMOVED***

func (b ButtonElement) Form() goja.Value ***REMOVED***
	formSel, exists := b.ownerFormSel()
	if !exists ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return selToElement(Selection***REMOVED***b.sel.rt, formSel***REMOVED***)
***REMOVED***

// Used by the formAction, formMethod, formTarget and formEnctype methods of Button and Input elements
// Attempts to read attribute "form" + attrName on the current element or attrName on the owning form element
func (e Element) formAttrOrElemOverride(attrName string) string ***REMOVED***
	if elemAttr, exists := e.sel.sel.Attr("form" + attrName); exists ***REMOVED***
		return elemAttr
	***REMOVED***

	formSel, exists := e.ownerFormSel()
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
func (e Element) formOrElemAttrString(attrName string) string ***REMOVED***
	if elemAttr, exists := e.sel.sel.Attr("form" + attrName); exists ***REMOVED***
		return elemAttr
	***REMOVED***

	formSel, exists := e.ownerFormSel()
	if !exists ***REMOVED***
		return ""
	***REMOVED***

	formAttr, exists := formSel.Attr(attrName)
	if !exists ***REMOVED***
		return ""
	***REMOVED***

	return formAttr
***REMOVED***

func (e Element) formOrElemAttrPresent(attrName string) bool ***REMOVED***
	if _, exists := e.sel.sel.Attr("form" + attrName); exists ***REMOVED***
		return true
	***REMOVED***

	formSel, exists := e.ownerFormSel()
	if !exists ***REMOVED***
		return false
	***REMOVED***

	_, exists = formSel.Attr(attrName)
	return exists
***REMOVED***

func (b ButtonElement) FormAction() string ***REMOVED***
	return b.formOrElemAttrString("action")
***REMOVED***

func (b ButtonElement) FormEnctype() string ***REMOVED***
	return b.formOrElemAttrString("enctype")
***REMOVED***

func (b ButtonElement) FormMethod() string ***REMOVED***
	return b.formOrElemAttrString("method")
***REMOVED***

func (b ButtonElement) FormNoValidate() bool ***REMOVED***
	return b.formOrElemAttrPresent("novalidate")
***REMOVED***

func (b ButtonElement) FormTarget() string ***REMOVED***
	return b.formOrElemAttrString("target")
***REMOVED***

func (e Element) elemLabels() (items []goja.Value) ***REMOVED***
	wrapperLbl := e.sel.sel.Closest("label")

	id := e.attrAsString("id")
	if id == "" ***REMOVED***
		return elemList(Selection***REMOVED***e.sel.rt, wrapperLbl***REMOVED***)
	***REMOVED***

	idLbl := e.sel.sel.Parents().Last().Find("label[for=\"" + id + "\"]")
	if idLbl.Size() == 0 ***REMOVED***
		return elemList(Selection***REMOVED***e.sel.rt, wrapperLbl***REMOVED***)
	***REMOVED***

	allLbls := wrapperLbl.AddSelection(idLbl)

	return elemList(Selection***REMOVED***e.sel.rt, allLbls***REMOVED***)
***REMOVED***

func (b ButtonElement) Labels() (items []goja.Value) ***REMOVED***
	return b.elemLabels()
***REMOVED***

func (b ButtonElement) Name() string ***REMOVED***
	return b.attrAsString("name")
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

func (b ButtonElement) Value() string ***REMOVED***
	return b.attrAsString("value")
***REMOVED***
