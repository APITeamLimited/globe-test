package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
)

const (
	stringTemplate   = "string"
	urlTemplate      = "url"
	enumTemplate     = "enum"
	boolTemplate     = "bool"
	gojaEnumTemplate = "enum-goja"
	intTemplate      = "int"
	constTemplate    = "const"
)

// The ast parser populates this struct using the <ElemName>TagName const and <ElemName>Element struct
type ElemInfo struct ***REMOVED***
	StructName    string
	PrtStructName string
***REMOVED***

// The signature for functions which inspect the ast nodes
type ElemInfoCollector func(node ast.Node) ElemInfoCollector

type ElemInfoCollectorState struct ***REMOVED***
	handler   ElemInfoCollector // The current function to check ast nodes
	elemName  string            // Only valid when a TagName const or Element struct encountered and used as an index into elemInfos
	elemInfos map[string]*ElemInfo
***REMOVED***

// "Elem" is the struct name for an Element. "Method" is the go method name. "Attr" is the name of the DOM attribute the method will access.
// "TemplateType" is used the "elemFuncsTemplate" and the returnType function - either string
// The Opts property:
//     for "string" and "bool" template types it is nil.
//     for "const", "int" and "url" templates it is the default return value when the attribute is unset. The "url" default should be an empty string or e.sel.URL
//     for "enum" template it is the list of valid options used in a switch statement - the default case is the first in the Opts list
//     for "enum-goja" template - same as the "enum" template except the function returns a goja.Value and the default value is always null
var funcDefs = []struct ***REMOVED***
	Elem, Method, Attr, TemplateType string
	Opts                             []string
***REMOVED******REMOVED***
	***REMOVED***"HrefElement", "Download", "download", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "ReferrerPolicy", "referrerpolicy", enumTemplate, []string***REMOVED***"", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "unsafe-url"***REMOVED******REMOVED***,
	***REMOVED***"HrefElement", "Rel", "rel", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "Href", "href", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"HrefElement", "Target", "target", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "AccessKey", "accesskey", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "HrefLang", "hreflang", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "ToString", "href", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"MediaElement", "Autoplay", "autoplay", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "Controls", "controls", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "Loop", "loop", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "Muted", "muted", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "Preload", "preload", enumTemplate, []string***REMOVED***"auto", "metadata", "none"***REMOVED******REMOVED***,
	***REMOVED***"MediaElement", "Src", "src", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"MediaElement", "CrossOrigin", "crossorigin", gojaEnumTemplate, []string***REMOVED***"anonymous", "use-credentials"***REMOVED******REMOVED***,
	***REMOVED***"MediaElement", "CurrentSrc", "src", stringTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "DefaultMuted", "muted", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "MediaGroup", "mediagroup", stringTemplate, nil***REMOVED***,
	***REMOVED***"BaseElement", "Href", "href", urlTemplate, []string***REMOVED***"e.sel.URL"***REMOVED******REMOVED***,
	***REMOVED***"BaseElement", "Target", "target", stringTemplate, nil***REMOVED***,
	***REMOVED***"ButtonElement", "AccessKey", "accesskey", stringTemplate, nil***REMOVED***,
	***REMOVED***"ButtonElement", "Autofocus", "autofocus", boolTemplate, nil***REMOVED***,
	***REMOVED***"ButtonElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"ButtonElement", "TabIndex", "tabindex", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ButtonElement", "Type", "type", enumTemplate, []string***REMOVED***"submit", "button", "menu", "reset"***REMOVED******REMOVED***,
	***REMOVED***"DataElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"EmbedElement", "Height", "height", stringTemplate, nil***REMOVED***,
	***REMOVED***"EmbedElement", "Width", "width", stringTemplate, nil***REMOVED***,
	***REMOVED***"EmbedElement", "Src", "src", stringTemplate, nil***REMOVED***,
	***REMOVED***"EmbedElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"FieldSetElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"FieldSetElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"FormElement", "Action", "action", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"FormElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"FormElement", "Target", "target", stringTemplate, nil***REMOVED***,
	***REMOVED***"FormElement", "Enctype", "enctype", enumTemplate, []string***REMOVED***"application/x-www-form-urlencoded", "multipart/form-data", "text/plain"***REMOVED******REMOVED***,
	***REMOVED***"FormElement", "Encoding", "enctype", enumTemplate, []string***REMOVED***"application/x-www-form-urlencoded", "multipart/form-data", "text/plain"***REMOVED******REMOVED***,
	***REMOVED***"FormElement", "AcceptCharset", "accept-charset", stringTemplate, nil***REMOVED***,
	***REMOVED***"FormElement", "Autocomplete", "autocomplete", enumTemplate, []string***REMOVED***"on", "off"***REMOVED******REMOVED***,
	***REMOVED***"FormElement", "NoValidate", "novalidate", boolTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "Allowfullscreen", "allowfullscreen", boolTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "ReferrerPolicy", "referrerpolicy", enumTemplate, []string***REMOVED***"", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "unsafe-url"***REMOVED******REMOVED***,
	***REMOVED***"IFrameElement", "Height", "height", stringTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "Width", "width", stringTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "Src", "src", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "CurrentSrc", "src", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "Sizes", "sizes", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "Srcset", "srcset", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "Alt", "alt", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "CrossOrigin", "crossorigin", gojaEnumTemplate, []string***REMOVED***"anonymous", "use-credentials"***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "Height", "height", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "Width", "width", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "IsMap", "ismap", boolTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "Src", "src", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "UseMap", "usemap", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "ReferrerPolicy", "referrerpolicy", enumTemplate, []string***REMOVED***"", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "unsafe-url"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "TabIndex", "tabindex", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Type", "type", enumTemplate, []string***REMOVED***"text", "button", "checkbox", "color", "date", "datetime-local", "email", "file", "hidden", "image", "month", "number", "password", "radio", "range", "reset", "search", "submit", "tel", "time", "url", "week"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Autofocus", "autofocus", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Required", "required", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Checked", "checked", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "DefaultChecked", "checked", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Alt", "alt", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Src", "src", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Height", "height", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Width", "width", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Accept", "accept", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Autocomplete", "autocomplete", enumTemplate, []string***REMOVED***"on", "off"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "MaxLength", "maxlength", intTemplate, []string***REMOVED***"-1"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Size", "size", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Pattern", "pattern", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Placeholder", "placeholder", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Readonly", "readonly", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Min", "min", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Max", "max", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "DefaultValue", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "DirName", "dirname", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "AccessKey", "accesskey", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Multiple", "multiple", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Step", "step", stringTemplate, nil***REMOVED***,
	***REMOVED***"KeygenElement", "Autofocus", "autofocus", boolTemplate, nil***REMOVED***,
	***REMOVED***"KeygenElement", "Challenge", "challenge", stringTemplate, nil***REMOVED***,
	***REMOVED***"KeygenElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"KeygenElement", "Keytype", "keytype", enumTemplate, []string***REMOVED***"RSA", "DSA", "EC"***REMOVED******REMOVED***,
	***REMOVED***"KeygenElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"KeygenElement", "Type", "type", constTemplate, []string***REMOVED***"keygen"***REMOVED******REMOVED***,
	***REMOVED***"LabelElement", "HtmlFor", "for", stringTemplate, nil***REMOVED***,
	***REMOVED***"LegendElement", "AccessKey", "accesskey", stringTemplate, nil***REMOVED***,
	***REMOVED***"LiElement", "Value", "value", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"LiElement", "Type", "type", enumTemplate, []string***REMOVED***"", "1", "a", "A", "i", "I", "disc", "square", "circle"***REMOVED******REMOVED***,
	***REMOVED***"LinkElement", "CrossOrigin", "crossorigin", gojaEnumTemplate, []string***REMOVED***"anonymous", "use-credentials"***REMOVED******REMOVED***,
	***REMOVED***"LinkElement", "ReferrerPolicy", "referrerpolicy", enumTemplate, []string***REMOVED***"", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "unsafe-url"***REMOVED******REMOVED***,
	***REMOVED***"LinkElement", "Href", "href", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"LinkElement", "Hreflang", "hreflang", stringTemplate, nil***REMOVED***,
	***REMOVED***"LinkElement", "Media", "media", stringTemplate, nil***REMOVED***,
	***REMOVED***"LinkElement", "Rel", "rel", stringTemplate, nil***REMOVED***,
	***REMOVED***"LinkElement", "Target", "target", stringTemplate, nil***REMOVED***,
	***REMOVED***"LinkElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"MapElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"MetaElement", "Content", "content", stringTemplate, nil***REMOVED***,
	***REMOVED***"MetaElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"MetaElement", "HttpEquiv", "http-equiv", enumTemplate, []string***REMOVED***"content-type", "default-style", "refresh"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "Min", "min", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "Max", "max", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "High", "high", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "Low", "low", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "Optimum", "optimum", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ModElement", "Cite", "cite", stringTemplate, nil***REMOVED***,
	***REMOVED***"ModElement", "Datetime", "datetime", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "Data", "data", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ObjectElement", "Height", "height", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "TabIndex", "tabindex", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ObjectElement", "TypeMustMatch", "typemustmatch", boolTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "UseMap", "usemap", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "Width", "width", stringTemplate, nil***REMOVED***,
	***REMOVED***"OListElement", "Reversed", "reversed", boolTemplate, nil***REMOVED***,
	***REMOVED***"OListElement", "Start", "start", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"OListElement", "Type", "type", enumTemplate, []string***REMOVED***"1", "a", "A", "i", "I"***REMOVED******REMOVED***,
	***REMOVED***"OptGroupElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"OptGroupElement", "Label", "label", stringTemplate, nil***REMOVED***,
	***REMOVED***"OptionElement", "DefaultSelected", "selected", boolTemplate, nil***REMOVED***,
	***REMOVED***"OptionElement", "Selected", "selected", boolTemplate, nil***REMOVED***,
	***REMOVED***"OutputElement", "HtmlFor", "for", stringTemplate, nil***REMOVED***,
	***REMOVED***"OutputElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"OutputElement", "Type", "type", constTemplate, []string***REMOVED***"output"***REMOVED******REMOVED***,
	***REMOVED***"ParamElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"ParamElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"PreElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"PreElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"QuoteElement", "Cite", "cite", stringTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "CrossOrigin", "crossorigin", stringTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "Src", "src", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ScriptElement", "Charset", "charset", stringTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "Async", "async", boolTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "Defer", "defer", boolTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "NoModule", "nomodule", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Autofocus", "autofocus", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Multiple", "multiple", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Required", "required", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "TabIndex", "tabindex", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"SourceElement", "KeySystem", "keysystem", stringTemplate, nil***REMOVED***,
	***REMOVED***"SourceElement", "Media", "media", stringTemplate, nil***REMOVED***,
	***REMOVED***"SourceElement", "Sizes", "sizes", stringTemplate, nil***REMOVED***,
	***REMOVED***"SourceElement", "Src", "src", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"SourceElement", "Srcset", "srcset", stringTemplate, nil***REMOVED***,
	***REMOVED***"SourceElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"StyleElement", "Media", "media", stringTemplate, nil***REMOVED***,
	***REMOVED***"TableElement", "Sortable", "sortable", boolTemplate, nil***REMOVED***,
	***REMOVED***"TableCellElement", "ColSpan", "colspan", intTemplate, []string***REMOVED***"1"***REMOVED******REMOVED***,
	***REMOVED***"TableCellElement", "RowSpan", "rowspan", intTemplate, []string***REMOVED***"1"***REMOVED******REMOVED***,
	***REMOVED***"TableCellElement", "Headers", "headers", stringTemplate, nil***REMOVED***,
	***REMOVED***"TableHeaderCellElement", "Abbr", "abbr", stringTemplate, nil***REMOVED***,
	***REMOVED***"TableHeaderCellElement", "Scope", "scope", enumTemplate, []string***REMOVED***"", "row", "col", "colgroup", "rowgroup"***REMOVED******REMOVED***,
	***REMOVED***"TableHeaderCellElement", "Sorted", "sorted", boolTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Type", "type", constTemplate, []string***REMOVED***"textarea"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "DefaultValue", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Placeholder", "placeholder", stringTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Rows", "rows", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "Cols", "cols", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "MaxLength", "maxlength", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "TabIndex", "tabindex", intTemplate, []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "AccessKey", "accesskey", stringTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "ReadOnly", "readonly", boolTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Required", "required", boolTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Autocomplete", "autocomplete", enumTemplate, []string***REMOVED***"on", "off"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "Autocapitalize", "autocapitalize", enumTemplate, []string***REMOVED***"sentences", "none", "off", "characters", "words"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "Wrap", "wrap", enumTemplate, []string***REMOVED***"soft", "hard", "off"***REMOVED******REMOVED***,
	***REMOVED***"TimeElement", "Datetime", "datetime", stringTemplate, nil***REMOVED***,
	***REMOVED***"TrackElement", "Kind", "kind", enumTemplate, []string***REMOVED***"subtitle", "captions", "descriptions", "chapters", "metadata"***REMOVED******REMOVED***,
	***REMOVED***"TrackElement", "Src", "src", urlTemplate, []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"TrackElement", "Srclang", "srclang", stringTemplate, nil***REMOVED***,
	***REMOVED***"TrackElement", "Label", "label", stringTemplate, nil***REMOVED***,
	***REMOVED***"TrackElement", "Default", "default", boolTemplate, nil***REMOVED***,
	***REMOVED***"UListElement", "Type", "type", stringTemplate, nil***REMOVED***,
***REMOVED***

var collector = &ElemInfoCollectorState***REMOVED******REMOVED***

func main() ***REMOVED***
	fs := token.NewFileSet()
	parsedFile, parseErr := parser.ParseFile(fs, "elements.go", nil, 0)
	if parseErr != nil ***REMOVED***
		log.WithError(parseErr).Fatal("Could not parse elements.go")
	***REMOVED***

	collector.handler = collector.defaultHandler
	collector.elemInfos = make(map[string]*ElemInfo)

	// Populate the elemInfos data
	ast.Inspect(parsedFile, func(n ast.Node) bool ***REMOVED***
		if n != nil ***REMOVED***
			collector.handler = collector.handler(n)
		***REMOVED***
		return true
	***REMOVED***)

	var buf bytes.Buffer
	err := elemFuncsTemplate.Execute(&buf, struct ***REMOVED***
		ElemInfos map[string]*ElemInfo
		FuncDefs  []struct ***REMOVED***
			Elem, Method, Attr, TemplateType string
			Opts                             []string
		***REMOVED***
		TemplateType struct***REMOVED*** String, Url, Enum, Bool, GojaEnum, Int, Const string ***REMOVED***
	***REMOVED******REMOVED***
		collector.elemInfos,
		funcDefs,
		struct***REMOVED*** String, Url, Enum, Bool, GojaEnum, Int, Const string ***REMOVED******REMOVED***stringTemplate, urlTemplate, enumTemplate, boolTemplate, gojaEnumTemplate, intTemplate, constTemplate***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Unable to execute template")
	***REMOVED***

	src, err := format.Source(buf.Bytes())
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("format.Source on generated code failed")
	***REMOVED***

	f, err := os.Create("elements_gen.go")
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Unable to create the file 'elements_gen.go'")
	***REMOVED***

	if _, err = f.Write(src); err != nil ***REMOVED***
		log.WithError(err).Fatal("Unable to write to 'elements_gen.go'")
	***REMOVED***

	err = f.Close()
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Unable to close 'elements_gen.go'")
	***REMOVED***
***REMOVED***

var elemFuncsTemplate = template.Must(template.New("").Funcs(template.FuncMap***REMOVED***
	"buildStruct": buildStruct,
	"returnType":  returnType,
***REMOVED***).Parse(`// generated by js/modules/k6/html/gen/gen_elements.go directed by js/modules/k6/html/elements.go;  DO NOT EDIT
// nolint: goconst
package html

import "github.com/dop251/goja"

func selToElement(sel Selection) goja.Value ***REMOVED***
	if sel.sel.Length() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	elem := Element***REMOVED***sel.sel.Nodes[0], &sel***REMOVED***

	switch elem.node.Data ***REMOVED*** 
***REMOVED******REMOVED***- range $elemName, $elemInfo := .ElemInfos ***REMOVED******REMOVED***
	case ***REMOVED******REMOVED*** $elemName ***REMOVED******REMOVED***TagName:
		return sel.rt.ToValue(***REMOVED******REMOVED*** buildStruct $elemInfo ***REMOVED******REMOVED***)
***REMOVED******REMOVED***- end ***REMOVED******REMOVED***
	default:
		return sel.rt.ToValue(elem)
	***REMOVED***
 ***REMOVED***

***REMOVED******REMOVED*** $templateType := .TemplateType ***REMOVED******REMOVED***
***REMOVED******REMOVED*** range $funcDef := .FuncDefs -***REMOVED******REMOVED*** 

func (e ***REMOVED******REMOVED***$funcDef.Elem***REMOVED******REMOVED***) ***REMOVED******REMOVED***$funcDef.Method***REMOVED******REMOVED***() ***REMOVED******REMOVED*** returnType $funcDef.TemplateType ***REMOVED******REMOVED*** ***REMOVED***
***REMOVED******REMOVED***- if eq $funcDef.TemplateType $templateType.Int ***REMOVED******REMOVED***
	return e.attrAsInt("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***", ***REMOVED******REMOVED*** index $funcDef.Opts 0 ***REMOVED******REMOVED***)
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateType.Enum ***REMOVED******REMOVED***
	attrVal := e.attrAsString("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***")
	switch attrVal ***REMOVED*** 
	***REMOVED******REMOVED***- range $optIdx, $optVal := $funcDef.Opts ***REMOVED******REMOVED***
	***REMOVED******REMOVED***- if ne $optIdx 0 ***REMOVED******REMOVED***
	case "***REMOVED******REMOVED***$optVal***REMOVED******REMOVED***":
		return attrVal
	***REMOVED******REMOVED***- end ***REMOVED******REMOVED***
	***REMOVED******REMOVED***- end***REMOVED******REMOVED***
	default: 
		return "***REMOVED******REMOVED*** index $funcDef.Opts 0 ***REMOVED******REMOVED***" 
	***REMOVED***
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateType.GojaEnum ***REMOVED******REMOVED***
	attrVal, exists := e.sel.sel.Attr("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***")
	if !exists ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	switch attrVal ***REMOVED*** 
	***REMOVED******REMOVED***- range $optVal := $funcDef.Opts ***REMOVED******REMOVED***
	case "***REMOVED******REMOVED***$optVal***REMOVED******REMOVED***":
		return e.sel.rt.ToValue(attrVal)
	***REMOVED******REMOVED***- end***REMOVED******REMOVED***
	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateType.Const ***REMOVED******REMOVED***
	return "***REMOVED******REMOVED*** index $funcDef.Opts 0 ***REMOVED******REMOVED***"
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateType.Url ***REMOVED******REMOVED***
	return e.attrAsURLString("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***", ***REMOVED******REMOVED*** index $funcDef.Opts 0 ***REMOVED******REMOVED***)
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateType.String ***REMOVED******REMOVED***
	return e.attrAsString("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***")
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateType.Bool ***REMOVED******REMOVED***
	return e.attrIsPresent("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***")
***REMOVED******REMOVED***- end***REMOVED******REMOVED***
***REMOVED***
***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
`))

func buildStruct(elemInfo ElemInfo) string ***REMOVED***
	if elemInfo.PrtStructName == "Element" ***REMOVED***
		return elemInfo.StructName + "***REMOVED***elem***REMOVED***"
	***REMOVED*** else ***REMOVED***
		return elemInfo.StructName + "***REMOVED***" + elemInfo.PrtStructName + "***REMOVED***elem***REMOVED******REMOVED***"
	***REMOVED***
***REMOVED***

func returnType(templateType string) string ***REMOVED***
	switch templateType ***REMOVED***
	case boolTemplate:
		return "bool"
	case intTemplate:
		return "int"
	case gojaEnumTemplate:
		return "goja.Value"
	default:
		return "string"
	***REMOVED***
***REMOVED***

// Node handler functions for ast.Inspect.
func (ce *ElemInfoCollectorState) defaultHandler(node ast.Node) ElemInfoCollector ***REMOVED***
	ce.elemName = ""
	switch node.(type) ***REMOVED***
	case *ast.TypeSpec:
		return ce.elemTypeSpecHandler

	case *ast.ValueSpec:
		return ce.tagNameValueSpecHandler

	default:
		return ce.defaultHandler
	***REMOVED***
***REMOVED***

// Found a tagname constant. eg AnchorTagName = "a" adds the entry ce.elemInfos["Anchor"] = ***REMOVED***""", ""***REMOVED***
func (ce *ElemInfoCollectorState) tagNameValueSpecHandler(node ast.Node) ElemInfoCollector ***REMOVED***
	switch x := node.(type) ***REMOVED***
	case *ast.Ident:
		if strings.HasSuffix(x.Name, "TagName") ***REMOVED***
			ce.elemName = strings.TrimSuffix(x.Name, "TagName")
			ce.elemInfos[ce.elemName] = &ElemInfo***REMOVED***"", ""***REMOVED***
		***REMOVED***

		return ce.defaultHandler

	default:
		return ce.defaultHandler
	***REMOVED***
***REMOVED***

func (ce *ElemInfoCollectorState) elemTypeSpecHandler(node ast.Node) ElemInfoCollector ***REMOVED***
	switch x := node.(type) ***REMOVED***
	case *ast.Ident:
		if !strings.HasSuffix(x.Name, "Element") ***REMOVED***
			return ce.defaultHandler
		***REMOVED***

		if ce.elemName == "" ***REMOVED***
			ce.elemName = strings.TrimSuffix(x.Name, "Element")
			// Ignore elements which don't have a tag name constant meaning no elemInfo structure was created by the TagName handle.
			// It skips the Href, Media, FormField, Mod, TableSection or TableCell structs as these structs are inherited by other elements and not created indepedently.
			if _, ok := ce.elemInfos[ce.elemName]; !ok ***REMOVED***
				return ce.defaultHandler
			***REMOVED***

			ce.elemInfos[ce.elemName].StructName = x.Name
			return ce.elemTypeSpecHandler
		***REMOVED*** else ***REMOVED***
			ce.elemInfos[ce.elemName].PrtStructName = x.Name
			return ce.defaultHandler
		***REMOVED***

	case *ast.StructType:
		return ce.elemTypeSpecHandler

	case *ast.FieldList:
		return ce.elemTypeSpecHandler

	case *ast.Field:
		return ce.elemTypeSpecHandler

	default:
		return ce.defaultHandler
	***REMOVED***
***REMOVED***
