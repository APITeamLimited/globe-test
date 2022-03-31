/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package html

import (
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/dop251/goja"
)

//go:generate go run gen/gen_elements.go
//go:generate gofumpt -w .

//nolint:gochecknoglobals
var defaultPorts = map[string]string***REMOVED***
	"http":  "80",
	"https": "443",
	"ftp":   "21",
***REMOVED***

// The code generator depends on the TagName constants being defined before the Element structs
const (
	AnchorTagName          = "a"
	AreaTagName            = "area"
	AudioTagName           = "audio"
	BaseTagName            = "base"
	ButtonTagName          = "button"
	CanvasTagName          = "canvas"
	DataTagName            = "data"
	DataListTagName        = "datalist"
	DelTagName             = "del"
	EmbedTagName           = "embed"
	FieldSetTagName        = "fieldset"
	FormTagName            = "form"
	IFrameTagName          = "iframe"
	ImageTagName           = "img"
	InputTagName           = "input"
	InsTagName             = "ins"
	KeygenTagName          = "keygen"
	LabelTagName           = "label"
	LegendTagName          = "legend"
	LiTagName              = "li"
	LinkTagName            = "link"
	MapTagName             = "map"
	MetaTagName            = "meta"
	MeterTagName           = "meter"
	ObjectTagName          = "object"
	OListTagName           = "ol"
	OptGroupTagName        = "optgroup"
	OptionTagName          = "option"
	OutputTagName          = "output"
	ParamTagName           = "param"
	PreTagName             = "pre"
	ProgressTagName        = "progress"
	QuoteTagName           = "quote"
	ScriptTagName          = "script"
	SelectTagName          = "select"
	SourceTagName          = "source"
	StyleTagName           = "style"
	TableTagName           = "table"
	TableHeadTagName       = "thead"
	TableFootTagName       = "tfoot"
	TableBodyTagName       = "tbody"
	TableRowTagName        = "tr"
	TableColTagName        = "col"
	TableDataCellTagName   = "td"
	TableHeaderCellTagName = "th"
	TextAreaTagName        = "textarea"
	TimeTagName            = "time"
	TitleTagName           = "title"
	TrackTagName           = "track"
	UListTagName           = "ul"
	VideoTagName           = "video"

	methodPost = "post"
	methodGet  = "get"
)

//revive:disable:exported

type (
	HrefElement            struct***REMOVED*** Element ***REMOVED***
	MediaElement           struct***REMOVED*** Element ***REMOVED***
	FormFieldElement       struct***REMOVED*** Element ***REMOVED***
	ModElement             struct***REMOVED*** Element ***REMOVED***
	TableSectionElement    struct***REMOVED*** Element ***REMOVED***
	TableCellElement       struct***REMOVED*** Element ***REMOVED***
	AnchorElement          struct***REMOVED*** HrefElement ***REMOVED***
	AreaElement            struct***REMOVED*** HrefElement ***REMOVED***
	AudioElement           struct***REMOVED*** MediaElement ***REMOVED***
	BaseElement            struct***REMOVED*** Element ***REMOVED***
	ButtonElement          struct***REMOVED*** FormFieldElement ***REMOVED***
	CanvasElement          struct***REMOVED*** Element ***REMOVED***
	DataElement            struct***REMOVED*** Element ***REMOVED***
	DataListElement        struct***REMOVED*** Element ***REMOVED***
	DelElement             struct***REMOVED*** ModElement ***REMOVED***
	InsElement             struct***REMOVED*** ModElement ***REMOVED***
	EmbedElement           struct***REMOVED*** Element ***REMOVED***
	FieldSetElement        struct***REMOVED*** Element ***REMOVED***
	FormElement            struct***REMOVED*** Element ***REMOVED***
	IFrameElement          struct***REMOVED*** Element ***REMOVED***
	ImageElement           struct***REMOVED*** Element ***REMOVED***
	InputElement           struct***REMOVED*** FormFieldElement ***REMOVED***
	KeygenElement          struct***REMOVED*** Element ***REMOVED***
	LabelElement           struct***REMOVED*** Element ***REMOVED***
	LegendElement          struct***REMOVED*** Element ***REMOVED***
	LiElement              struct***REMOVED*** Element ***REMOVED***
	LinkElement            struct***REMOVED*** Element ***REMOVED***
	MapElement             struct***REMOVED*** Element ***REMOVED***
	MetaElement            struct***REMOVED*** Element ***REMOVED***
	MeterElement           struct***REMOVED*** Element ***REMOVED***
	ObjectElement          struct***REMOVED*** Element ***REMOVED***
	OListElement           struct***REMOVED*** Element ***REMOVED***
	OptGroupElement        struct***REMOVED*** Element ***REMOVED***
	OptionElement          struct***REMOVED*** Element ***REMOVED***
	OutputElement          struct***REMOVED*** Element ***REMOVED***
	ParamElement           struct***REMOVED*** Element ***REMOVED***
	PreElement             struct***REMOVED*** Element ***REMOVED***
	ProgressElement        struct***REMOVED*** Element ***REMOVED***
	QuoteElement           struct***REMOVED*** Element ***REMOVED***
	ScriptElement          struct***REMOVED*** Element ***REMOVED***
	SelectElement          struct***REMOVED*** Element ***REMOVED***
	SourceElement          struct***REMOVED*** Element ***REMOVED***
	StyleElement           struct***REMOVED*** Element ***REMOVED***
	TableElement           struct***REMOVED*** Element ***REMOVED***
	TableHeadElement       struct***REMOVED*** TableSectionElement ***REMOVED***
	TableFootElement       struct***REMOVED*** TableSectionElement ***REMOVED***
	TableBodyElement       struct***REMOVED*** TableSectionElement ***REMOVED***
	TableRowElement        struct***REMOVED*** Element ***REMOVED***
	TableColElement        struct***REMOVED*** Element ***REMOVED***
	TableDataCellElement   struct***REMOVED*** TableCellElement ***REMOVED***
	TableHeaderCellElement struct***REMOVED*** TableCellElement ***REMOVED***
	TextAreaElement        struct***REMOVED*** Element ***REMOVED***
	TimeElement            struct***REMOVED*** Element ***REMOVED***
	TitleElement           struct***REMOVED*** Element ***REMOVED***
	TrackElement           struct***REMOVED*** Element ***REMOVED***
	UListElement           struct***REMOVED*** Element ***REMOVED***
	VideoElement           struct***REMOVED*** MediaElement ***REMOVED***
)

//revive:enable:exported

func (h HrefElement) hrefURL() *url.URL ***REMOVED***
	href, exists := h.attrAsURL("href")
	if !exists ***REMOVED***
		return &url.URL***REMOVED******REMOVED***
	***REMOVED***
	return href
***REMOVED***

func (h HrefElement) Hash() string ***REMOVED***
	frag := h.hrefURL().Fragment
	if frag == "" ***REMOVED***
		return ""
	***REMOVED***
	return "#" + frag
***REMOVED***

func (h HrefElement) Host() string ***REMOVED***
	href := h.hrefURL()
	if href.Host == "" ***REMOVED***
		return ""
	***REMOVED***

	host, port, err := net.SplitHostPort(href.Host)
	if err != nil ***REMOVED***
		return href.Host
	***REMOVED***

	defaultPort := defaultPorts[href.Scheme]
	if defaultPort != "" && port == defaultPort ***REMOVED***
		return strings.TrimSuffix(host, ":"+defaultPort)
	***REMOVED***

	return href.Host
***REMOVED***

func (h HrefElement) Hostname() string ***REMOVED***
	hostAndPort := h.hrefURL().Host
	if hostAndPort == "" ***REMOVED***
		return ""
	***REMOVED***

	host, _, err := net.SplitHostPort(hostAndPort)
	if err != nil ***REMOVED***
		return hostAndPort
	***REMOVED***

	return host
***REMOVED***

func (h HrefElement) Port() string ***REMOVED***
	hostAndPort := h.hrefURL().Host
	if hostAndPort == "" ***REMOVED***
		return ""
	***REMOVED***

	_, port, err := net.SplitHostPort(hostAndPort)
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

func (h HrefElement) Password() string ***REMOVED***
	user := h.hrefURL().User
	if user == nil ***REMOVED***
		return ""
	***REMOVED***

	pwd, defined := user.Password()
	if !defined ***REMOVED***
		return ""
	***REMOVED***

	return pwd
***REMOVED***

func (h HrefElement) Origin() string ***REMOVED***
	href := h.hrefURL()

	if href.Scheme == "" ***REMOVED***
		return ""
	***REMOVED***

	if href.Scheme == "file" ***REMOVED***
		return h.Href()
	***REMOVED***

	return href.Scheme + "://" + href.Host
***REMOVED***

func (h HrefElement) Pathname() string ***REMOVED***
	return h.hrefURL().Path
***REMOVED***

func (h HrefElement) Protocol() string ***REMOVED***
	scheme := h.hrefURL().Scheme
	if scheme == "" ***REMOVED***
		return ":"
	***REMOVED***
	return scheme
***REMOVED***

func (h HrefElement) RelList() []string ***REMOVED***
	return h.splitAttr("rel")
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
	return f.ownerFormVal()
***REMOVED***

func (f FormFieldElement) formOrElemAttr(attrName string) (string, bool) ***REMOVED***
	if elemAttr, exists := f.sel.sel.Attr("form" + attrName); exists ***REMOVED***
		return elemAttr, true
	***REMOVED***

	formSel, exists := f.ownerFormSel()
	if !exists ***REMOVED***
		return "", false
	***REMOVED***

	formAttr, exists := formSel.Attr(attrName)
	if !exists ***REMOVED***
		return "", false
	***REMOVED***

	return formAttr, true
***REMOVED***

func (f FormFieldElement) FormAction() string ***REMOVED***
	action, exists := f.formOrElemAttr("action")
	if f.sel.URL == "" ***REMOVED***
		return action
	***REMOVED***

	if !exists || action == "" ***REMOVED***
		return f.sel.URL
	***REMOVED***

	actionURL, ok := f.resolveURL(action)
	if !ok ***REMOVED***
		return action
	***REMOVED***

	return actionURL.String()
***REMOVED***

func (f FormFieldElement) FormEnctype() string ***REMOVED***
	enctype, _ := f.formOrElemAttr("enctype")

	switch enctype ***REMOVED***
	case "multipart/form-data":
		return enctype
	case "text/plain":
		return enctype
	default:
		return "application/x-www-form-urlencoded"
	***REMOVED***
***REMOVED***

func (f FormFieldElement) FormMethod() string ***REMOVED***
	method, _ := f.formOrElemAttr("method")

	switch strings.ToLower(method) ***REMOVED***
	case methodPost:
		return methodPost
	default:
		return methodGet
	***REMOVED***
***REMOVED***

func (f FormFieldElement) FormNoValidate() bool ***REMOVED***
	_, exists := f.formOrElemAttr("novalidate")
	return exists
***REMOVED***

func (f FormFieldElement) FormTarget() string ***REMOVED***
	target, _ := f.formOrElemAttr("target")
	return target
***REMOVED***

func (f FormFieldElement) Labels() []goja.Value ***REMOVED***
	return f.elemLabels()
***REMOVED***

func (f FormFieldElement) Name() string ***REMOVED***
	return f.attrAsString("name")
***REMOVED***

func (b ButtonElement) Value() string ***REMOVED***
	return valueOrHTML(b.sel.sel)
***REMOVED***

func (c CanvasElement) Width() int ***REMOVED***
	return c.attrAsInt("width", 150)
***REMOVED***

func (c CanvasElement) Height() int ***REMOVED***
	return c.attrAsInt("height", 150)
***REMOVED***

func (d DataListElement) Options() []goja.Value ***REMOVED***
	return elemList(d.sel.Find("option"))
***REMOVED***

func (f FieldSetElement) Form() goja.Value ***REMOVED***
	formSel, exists := f.ownerFormSel()
	if !exists ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return selToElement(Selection***REMOVED***f.sel.rt, formSel, f.sel.URL***REMOVED***)
***REMOVED***

func (f FieldSetElement) Type() string ***REMOVED***
	return "fieldset"
***REMOVED***

func (f FieldSetElement) Elements() []goja.Value ***REMOVED***
	return elemList(f.sel.Find("input,select,button,textarea"))
***REMOVED***

func (f FieldSetElement) Validity() goja.Value ***REMOVED***
	return goja.Undefined()
***REMOVED***

func (f FormElement) Elements() []goja.Value ***REMOVED***
	return elemList(f.sel.Find("input,select,button,textarea,fieldset"))
***REMOVED***

func (f FormElement) Length() int ***REMOVED***
	return f.sel.sel.Find("input,select,button,textarea,fieldset").Size()
***REMOVED***

func (f FormElement) Method() string ***REMOVED***
	if method := f.attrAsString("method"); method == methodPost ***REMOVED***
		return methodPost
	***REMOVED***

	return methodGet
***REMOVED***

func (i InputElement) List() goja.Value ***REMOVED***
	listID := i.attrAsString("list")

	if listID == "" ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	switch i.attrAsString("type") ***REMOVED***
	case "hidden":
		return goja.Undefined()
	case "checkbox":
		return goja.Undefined()
	case "radio":
		return goja.Undefined()
	case "file":
		return goja.Undefined()
	case "button":
		return goja.Undefined()
	***REMOVED***

	datalist := i.sel.sel.Parents().Last().Find("datalist[id=\"" + listID + "\"]")
	if datalist.Length() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	return selToElement(Selection***REMOVED***i.sel.rt, datalist.Eq(0), i.sel.URL***REMOVED***)
***REMOVED***

func (k KeygenElement) Form() goja.Value ***REMOVED***
	return k.ownerFormVal()
***REMOVED***

func (k KeygenElement) Labels() []goja.Value ***REMOVED***
	return k.elemLabels()
***REMOVED***

func (l LabelElement) Control() goja.Value ***REMOVED***
	forAttr, exists := l.sel.sel.Attr("for")
	if !exists ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	findControl := l.sel.sel.Parents().Last().Find("#" + forAttr)
	if findControl.Length() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	return selToElement(Selection***REMOVED***l.sel.rt, findControl.Eq(0), l.sel.URL***REMOVED***)
***REMOVED***

func (l LabelElement) Form() goja.Value ***REMOVED***
	return l.ownerFormVal()
***REMOVED***

func (l LegendElement) Form() goja.Value ***REMOVED***
	return l.ownerFormVal()
***REMOVED***

func (l LinkElement) RelList() []string ***REMOVED***
	return l.splitAttr("rel")
***REMOVED***

func (m MapElement) Areas() []goja.Value ***REMOVED***
	return elemList(m.sel.Find("area"))
***REMOVED***

func (m MapElement) Images() []goja.Value ***REMOVED***
	name, exists := m.idOrNameAttr()

	if !exists ***REMOVED***
		return make([]goja.Value, 0)
	***REMOVED***

	imgs := m.sel.sel.Parents().Last().Find("img[usemap=\"#" + name + "\"],object[usemap=\"#" + name + "\"]")
	return elemList(Selection***REMOVED***m.sel.rt, imgs, m.sel.URL***REMOVED***)
***REMOVED***

func (m MeterElement) Labels() []goja.Value ***REMOVED***
	return m.elemLabels()
***REMOVED***

func (o ObjectElement) Form() goja.Value ***REMOVED***
	return o.ownerFormVal()
***REMOVED***

func (o OptionElement) Disabled() bool ***REMOVED***
	if o.attrIsPresent("disabled") ***REMOVED***
		return true
	***REMOVED***

	optGroup := o.sel.sel.ParentsFiltered("optgroup")
	if optGroup.Length() == 0 ***REMOVED***
		return false
	***REMOVED***

	_, exists := optGroup.Attr("disabled")
	return exists
***REMOVED***

func (o OptionElement) Form() goja.Value ***REMOVED***
	prtForm := o.sel.sel.ParentsFiltered("form")
	if prtForm.Length() != 0 ***REMOVED***
		return selToElement(Selection***REMOVED***o.sel.rt, prtForm.First(), o.sel.URL***REMOVED***)
	***REMOVED***

	prtSelect := o.sel.sel.ParentsFiltered("select")
	formID, exists := prtSelect.Attr("form")
	if !exists ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	ownerForm := prtSelect.Parents().Last().Find("form#" + formID)
	if ownerForm.Length() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	return selToElement(Selection***REMOVED***o.sel.rt, ownerForm.First(), o.sel.URL***REMOVED***)
***REMOVED***

func (o OptionElement) Index() int ***REMOVED***
	optsHolder := o.sel.sel.ParentsFiltered("select,datalist")
	if optsHolder.Length() == 0 ***REMOVED***
		return 0
	***REMOVED***

	return optsHolder.Find("option").IndexOfSelection(o.sel.sel)
***REMOVED***

func (o OptionElement) Label() string ***REMOVED***
	if lbl, exists := o.sel.sel.Attr("label"); exists ***REMOVED***
		return lbl
	***REMOVED***

	return o.TextContent()
***REMOVED***

func (o OptionElement) Text() string ***REMOVED***
	return o.TextContent()
***REMOVED***

func (o OptionElement) Value() string ***REMOVED***
	return valueOrHTML(o.sel.sel)
***REMOVED***

func (o OutputElement) Form() goja.Value ***REMOVED***
	return o.ownerFormVal()
***REMOVED***

func (o OutputElement) Labels() []goja.Value ***REMOVED***
	return o.elemLabels()
***REMOVED***

func (o OutputElement) Value() string ***REMOVED***
	return o.TextContent()
***REMOVED***

func (o OutputElement) DefaultValue() string ***REMOVED***
	return o.TextContent()
***REMOVED***

func (p ProgressElement) Max() float64 ***REMOVED***
	maxStr, exists := p.sel.sel.Attr("max")
	if !exists ***REMOVED***
		return 1.0
	***REMOVED***

	maxVal, err := strconv.ParseFloat(maxStr, 64)
	if err != nil || maxVal < 0 ***REMOVED***
		return 1.0
	***REMOVED***

	return maxVal
***REMOVED***

func (p ProgressElement) calcProgress(defaultVal float64) float64 ***REMOVED***
	valStr, exists := p.sel.sel.Attr("value")
	if !exists ***REMOVED***
		return defaultVal
	***REMOVED***

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil || val < 0 ***REMOVED***
		return defaultVal
	***REMOVED***

	return val / p.Max()
***REMOVED***

func (p ProgressElement) Value() float64 ***REMOVED***
	return p.calcProgress(0.0)
***REMOVED***

func (p ProgressElement) Position() float64 ***REMOVED***
	return p.calcProgress(-1.0)
***REMOVED***

func (p ProgressElement) Labels() []goja.Value ***REMOVED***
	return p.elemLabels()
***REMOVED***

func (s ScriptElement) Text() string ***REMOVED***
	return s.TextContent()
***REMOVED***

func (s SelectElement) Form() goja.Value ***REMOVED***
	return s.ownerFormVal()
***REMOVED***

func (s SelectElement) Labels() []goja.Value ***REMOVED***
	return s.elemLabels()
***REMOVED***

func (s SelectElement) Length() int ***REMOVED***
	return s.sel.Find("option").Size()
***REMOVED***

func (s SelectElement) Options() []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***s.sel.rt, s.sel.sel.Find("option"), s.sel.URL***REMOVED***)
***REMOVED***

func (s SelectElement) SelectedIndex() int ***REMOVED***
	option := s.sel.sel.Find("option[selected]")
	if option.Length() == 0 ***REMOVED***
		return -1
	***REMOVED***
	return s.sel.sel.Find("option").IndexOfSelection(option)
***REMOVED***

func (s SelectElement) SelectedOptions() []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***s.sel.rt, s.sel.sel.Find("option[selected]"), s.sel.URL***REMOVED***)
***REMOVED***

func (s SelectElement) Size() int ***REMOVED***
	if s.attrIsPresent("multiple") ***REMOVED***
		return 4
	***REMOVED***
	return 1
***REMOVED***

func (s SelectElement) Type() string ***REMOVED***
	if s.attrIsPresent("multiple") ***REMOVED***
		return "select-multiple"
	***REMOVED***
	return "select"
***REMOVED***

func (s SelectElement) Value() string ***REMOVED***
	option := s.sel.sel.Find("option[selected]")
	if option.Length() == 0 ***REMOVED***
		return ""
	***REMOVED***
	return valueOrHTML(option.First())
***REMOVED***

func (s StyleElement) Type() string ***REMOVED***
	typeVal := s.attrAsString("type")
	if typeVal == "" ***REMOVED***
		return "text/css"
	***REMOVED***
	return typeVal
***REMOVED***

func (t TableElement) firstChild(elemName string) goja.Value ***REMOVED***
	child := t.sel.sel.ChildrenFiltered(elemName)
	if child.Size() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return selToElement(Selection***REMOVED***t.sel.rt, child, t.sel.URL***REMOVED***)
***REMOVED***

func (t TableElement) Caption() goja.Value ***REMOVED***
	return t.firstChild("caption")
***REMOVED***

func (t TableElement) THead() goja.Value ***REMOVED***
	return t.firstChild("thead")
***REMOVED***

func (t TableElement) TFoot() goja.Value ***REMOVED***
	return t.firstChild("tfoot")
***REMOVED***

func (t TableElement) Rows() []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***t.sel.rt, t.sel.sel.Find("tr"), t.sel.URL***REMOVED***)
***REMOVED***

func (t TableElement) TBodies() []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***t.sel.rt, t.sel.sel.Find("tbody"), t.sel.URL***REMOVED***)
***REMOVED***

func (t TableSectionElement) Rows() []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***t.sel.rt, t.sel.sel.Find("tr"), t.sel.URL***REMOVED***)
***REMOVED***

func (t TableCellElement) CellIndex() int ***REMOVED***
	prtRow := t.sel.sel.ParentsFiltered("tr")
	if prtRow.Length() == 0 ***REMOVED***
		return -1
	***REMOVED***
	return prtRow.Find("th,td").IndexOfSelection(t.sel.sel)
***REMOVED***

func (t TableRowElement) Cells() []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***t.sel.rt, t.sel.sel.Find("th,td"), t.sel.URL***REMOVED***)
***REMOVED***

func (t TableRowElement) RowIndex() int ***REMOVED***
	table := t.sel.sel.ParentsFiltered("table")
	if table.Length() == 0 ***REMOVED***
		return -1
	***REMOVED***
	return table.Find("tr").IndexOfSelection(t.sel.sel)
***REMOVED***

func (t TableRowElement) SectionRowIndex() int ***REMOVED***
	section := t.sel.sel.ParentsFiltered("thead,tbody,tfoot")
	if section.Length() == 0 ***REMOVED***
		return -1
	***REMOVED***
	return section.Find("tr").IndexOfSelection(t.sel.sel)
***REMOVED***

func (t TextAreaElement) Form() goja.Value ***REMOVED***
	return t.ownerFormVal()
***REMOVED***

func (t TextAreaElement) Length() int ***REMOVED***
	return len(t.attrAsString("value"))
***REMOVED***

func (t TextAreaElement) Labels() []goja.Value ***REMOVED***
	return t.elemLabels()
***REMOVED***

func (t TableColElement) Span() int ***REMOVED***
	span := t.attrAsInt("span", 1)
	if span < 1 ***REMOVED***
		return 1
	***REMOVED***
	return span
***REMOVED***

func (m MediaElement) TextTracks() []goja.Value ***REMOVED***
	return elemList(Selection***REMOVED***m.sel.rt, m.sel.sel.Find("track"), m.sel.URL***REMOVED***)
***REMOVED***

func (t TitleElement) Text() string ***REMOVED***
	return t.TextContent()
***REMOVED***
