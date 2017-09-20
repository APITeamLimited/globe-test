package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
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
	***REMOVED***"HrefElement", "Download", "download", "string", nil***REMOVED***,
	***REMOVED***"HrefElement", "ReferrerPolicy", "referrerpolicy", "enum", []string***REMOVED***"", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "unsafe-url"***REMOVED******REMOVED***,
	***REMOVED***"HrefElement", "Rel", "rel", "string", nil***REMOVED***,
	***REMOVED***"HrefElement", "Href", "href", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"HrefElement", "Target", "target", "string", nil***REMOVED***,
	***REMOVED***"HrefElement", "Type", "type", "string", nil***REMOVED***,
	***REMOVED***"HrefElement", "AccessKey", "accesskey", "string", nil***REMOVED***,
	***REMOVED***"HrefElement", "HrefLang", "hreflang", "string", nil***REMOVED***,
	***REMOVED***"HrefElement", "ToString", "href", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"MediaElement", "Autoplay", "autoplay", "bool", nil***REMOVED***,
	***REMOVED***"MediaElement", "Controls", "controls", "bool", nil***REMOVED***,
	***REMOVED***"MediaElement", "Loop", "loop", "bool", nil***REMOVED***,
	***REMOVED***"MediaElement", "Muted", "muted", "bool", nil***REMOVED***,
	***REMOVED***"MediaElement", "Preload", "preload", "enum", []string***REMOVED***"auto", "metadata", "none"***REMOVED******REMOVED***,
	***REMOVED***"MediaElement", "Src", "src", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"MediaElement", "CrossOrigin", "crossorigin", "enum-goja", []string***REMOVED***"anonymous", "use-credentials"***REMOVED******REMOVED***,
	***REMOVED***"MediaElement", "CurrentSrc", "src", "string", nil***REMOVED***,
	***REMOVED***"MediaElement", "DefaultMuted", "muted", "bool", nil***REMOVED***,
	***REMOVED***"MediaElement", "MediaGroup", "mediagroup", "string", nil***REMOVED***,
	***REMOVED***"BaseElement", "Href", "href", "url", []string***REMOVED***"e.sel.URL"***REMOVED******REMOVED***,
	***REMOVED***"BaseElement", "Target", "target", "string", nil***REMOVED***,
	***REMOVED***"ButtonElement", "AccessKey", "accesskey", "string", nil***REMOVED***,
	***REMOVED***"ButtonElement", "Autofocus", "autofocus", "bool", nil***REMOVED***,
	***REMOVED***"ButtonElement", "Disabled", "disabled", "bool", nil***REMOVED***,
	***REMOVED***"ButtonElement", "TabIndex", "tabindex", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ButtonElement", "Type", "type", "enum", []string***REMOVED***"submit", "button", "menu", "reset"***REMOVED******REMOVED***,
	***REMOVED***"DataElement", "Value", "value", "string", nil***REMOVED***,
	***REMOVED***"EmbedElement", "Height", "height", "string", nil***REMOVED***,
	***REMOVED***"EmbedElement", "Width", "width", "string", nil***REMOVED***,
	***REMOVED***"EmbedElement", "Src", "src", "string", nil***REMOVED***,
	***REMOVED***"EmbedElement", "Type", "type", "string", nil***REMOVED***,
	***REMOVED***"FieldSetElement", "Disabled", "disabled", "bool", nil***REMOVED***,
	***REMOVED***"FieldSetElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"FormElement", "Action", "action", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"FormElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"FormElement", "Target", "target", "string", nil***REMOVED***,
	***REMOVED***"FormElement", "Enctype", "enctype", "enum", []string***REMOVED***"application/x-www-form-urlencoded", "multipart/form-data", "text/plain"***REMOVED******REMOVED***,
	***REMOVED***"FormElement", "Encoding", "enctype", "enum", []string***REMOVED***"application/x-www-form-urlencoded", "multipart/form-data", "text/plain"***REMOVED******REMOVED***,
	***REMOVED***"FormElement", "AcceptCharset", "accept-charset", "string", nil***REMOVED***,
	***REMOVED***"FormElement", "Autocomplete", "autocomplete", "enum", []string***REMOVED***"on", "off"***REMOVED******REMOVED***,
	***REMOVED***"FormElement", "NoValidate", "novalidate", "bool", nil***REMOVED***,
	***REMOVED***"IFrameElement", "Allowfullscreen", "allowfullscreen", "bool", nil***REMOVED***,
	***REMOVED***"IFrameElement", "ReferrerPolicy", "referrerpolicy", "enum", []string***REMOVED***"", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "unsafe-url"***REMOVED******REMOVED***,
	***REMOVED***"IFrameElement", "Height", "height", "string", nil***REMOVED***,
	***REMOVED***"IFrameElement", "Width", "width", "string", nil***REMOVED***,
	***REMOVED***"IFrameElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"IFrameElement", "Src", "src", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "CurrentSrc", "src", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "Sizes", "sizes", "string", nil***REMOVED***,
	***REMOVED***"ImageElement", "Srcset", "srcset", "string", nil***REMOVED***,
	***REMOVED***"ImageElement", "Alt", "alt", "string", nil***REMOVED***,
	***REMOVED***"ImageElement", "CrossOrigin", "crossorigin", "enum-goja", []string***REMOVED***"anonymous", "use-credentials"***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "Height", "height", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "Width", "width", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "IsMap", "ismap", "bool", nil***REMOVED***,
	***REMOVED***"ImageElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"ImageElement", "Src", "src", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ImageElement", "UseMap", "usemap", "string", nil***REMOVED***,
	***REMOVED***"ImageElement", "ReferrerPolicy", "referrerpolicy", "enum", []string***REMOVED***"", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "unsafe-url"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "TabIndex", "tabindex", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Type", "type", "enum", []string***REMOVED***"text", "button", "checkbox", "color", "date", "datetime-local", "email", "file", "hidden", "image", "month", "number", "password", "radio", "range", "reset", "search", "submit", "tel", "time", "url", "week"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Disabled", "disabled", "bool", nil***REMOVED***,
	***REMOVED***"InputElement", "Autofocus", "autofocus", "bool", nil***REMOVED***,
	***REMOVED***"InputElement", "Required", "required", "bool", nil***REMOVED***,
	***REMOVED***"InputElement", "Value", "value", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "Checked", "checked", "bool", nil***REMOVED***,
	***REMOVED***"InputElement", "DefaultChecked", "checked", "bool", nil***REMOVED***,
	***REMOVED***"InputElement", "Alt", "alt", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "Src", "src", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Height", "height", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "Width", "width", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "Accept", "accept", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "Autocomplete", "autocomplete", "enum", []string***REMOVED***"on", "off"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "MaxLength", "maxlength", "int", []string***REMOVED***"-1"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Size", "size", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"InputElement", "Pattern", "pattern", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "Placeholder", "placeholder", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "Readonly", "readonly", "bool", nil***REMOVED***,
	***REMOVED***"InputElement", "Min", "min", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "Max", "max", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "DefaultValue", "value", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "DirName", "dirname", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "AccessKey", "accesskey", "string", nil***REMOVED***,
	***REMOVED***"InputElement", "Multiple", "multiple", "bool", nil***REMOVED***,
	***REMOVED***"InputElement", "Step", "step", "string", nil***REMOVED***,
	***REMOVED***"KeygenElement", "Autofocus", "autofocus", "bool", nil***REMOVED***,
	***REMOVED***"KeygenElement", "Challenge", "challenge", "string", nil***REMOVED***,
	***REMOVED***"KeygenElement", "Disabled", "disabled", "bool", nil***REMOVED***,
	***REMOVED***"KeygenElement", "Keytype", "keytype", "enum", []string***REMOVED***"RSA", "DSA", "EC"***REMOVED******REMOVED***,
	***REMOVED***"KeygenElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"KeygenElement", "Type", "type", "const", []string***REMOVED***"keygen"***REMOVED******REMOVED***,
	***REMOVED***"LabelElement", "HtmlFor", "for", "string", nil***REMOVED***,
	***REMOVED***"LegendElement", "AccessKey", "accesskey", "string", nil***REMOVED***,
	***REMOVED***"LiElement", "Value", "value", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"LiElement", "Type", "type", "enum", []string***REMOVED***"", "1", "a", "A", "i", "I", "disc", "square", "circle"***REMOVED******REMOVED***,
	***REMOVED***"LinkElement", "CrossOrigin", "crossorigin", "enum-goja", []string***REMOVED***"anonymous", "use-credentials"***REMOVED******REMOVED***,
	***REMOVED***"LinkElement", "ReferrerPolicy", "referrerpolicy", "enum", []string***REMOVED***"", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "unsafe-url"***REMOVED******REMOVED***,
	***REMOVED***"LinkElement", "Href", "href", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"LinkElement", "Hreflang", "hreflang", "string", nil***REMOVED***,
	***REMOVED***"LinkElement", "Media", "media", "string", nil***REMOVED***,
	***REMOVED***"LinkElement", "Rel", "rel", "string", nil***REMOVED***,
	***REMOVED***"LinkElement", "Target", "target", "string", nil***REMOVED***,
	***REMOVED***"LinkElement", "Type", "type", "string", nil***REMOVED***,
	***REMOVED***"MapElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"MetaElement", "Content", "content", "string", nil***REMOVED***,
	***REMOVED***"MetaElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"MetaElement", "HttpEquiv", "http-equiv", "enum", []string***REMOVED***"content-type", "default-style", "refresh"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "Min", "min", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "Max", "max", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "High", "high", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "Low", "low", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"MeterElement", "Optimum", "optimum", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ModElement", "Cite", "cite", "string", nil***REMOVED***,
	***REMOVED***"ModElement", "Datetime", "datetime", "string", nil***REMOVED***,
	***REMOVED***"ObjectElement", "Data", "data", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ObjectElement", "Height", "height", "string", nil***REMOVED***,
	***REMOVED***"ObjectElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"ObjectElement", "Type", "type", "string", nil***REMOVED***,
	***REMOVED***"ObjectElement", "TabIndex", "tabindex", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"ObjectElement", "TypeMustMatch", "typemustmatch", "bool", nil***REMOVED***,
	***REMOVED***"ObjectElement", "UseMap", "usemap", "string", nil***REMOVED***,
	***REMOVED***"ObjectElement", "Width", "width", "string", nil***REMOVED***,
	***REMOVED***"OListElement", "Reversed", "reversed", "bool", nil***REMOVED***,
	***REMOVED***"OListElement", "Start", "start", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"OListElement", "Type", "type", "enum", []string***REMOVED***"1", "a", "A", "i", "I"***REMOVED******REMOVED***,
	***REMOVED***"OptGroupElement", "Disabled", "disabled", "bool", nil***REMOVED***,
	***REMOVED***"OptGroupElement", "Label", "label", "string", nil***REMOVED***,
	***REMOVED***"OptionElement", "DefaultSelected", "selected", "bool", nil***REMOVED***,
	***REMOVED***"OptionElement", "Selected", "selected", "bool", nil***REMOVED***,
	***REMOVED***"OutputElement", "HtmlFor", "for", "string", nil***REMOVED***,
	***REMOVED***"OutputElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"OutputElement", "Type", "type", "const", []string***REMOVED***"output"***REMOVED******REMOVED***,
	***REMOVED***"ParamElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"ParamElement", "Value", "value", "string", nil***REMOVED***,
	***REMOVED***"PreElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"PreElement", "Value", "value", "string", nil***REMOVED***,
	***REMOVED***"QuoteElement", "Cite", "cite", "string", nil***REMOVED***,
	***REMOVED***"ScriptElement", "CrossOrigin", "crossorigin", "string", nil***REMOVED***,
	***REMOVED***"ScriptElement", "Type", "type", "string", nil***REMOVED***,
	***REMOVED***"ScriptElement", "Src", "src", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"ScriptElement", "Charset", "charset", "string", nil***REMOVED***,
	***REMOVED***"ScriptElement", "Async", "async", "bool", nil***REMOVED***,
	***REMOVED***"ScriptElement", "Defer", "defer", "bool", nil***REMOVED***,
	***REMOVED***"ScriptElement", "NoModule", "nomodule", "bool", nil***REMOVED***,
	***REMOVED***"SelectElement", "Autofocus", "autofocus", "bool", nil***REMOVED***,
	***REMOVED***"SelectElement", "Disabled", "disabled", "bool", nil***REMOVED***,
	***REMOVED***"SelectElement", "Multiple", "multiple", "bool", nil***REMOVED***,
	***REMOVED***"SelectElement", "Name", "name", "string", nil***REMOVED***,
	***REMOVED***"SelectElement", "Required", "required", "bool", nil***REMOVED***,
	***REMOVED***"SelectElement", "TabIndex", "tabindex", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"SourceElement", "KeySystem", "keysystem", "string", nil***REMOVED***,
	***REMOVED***"SourceElement", "Media", "media", "string", nil***REMOVED***,
	***REMOVED***"SourceElement", "Sizes", "sizes", "string", nil***REMOVED***,
	***REMOVED***"SourceElement", "Src", "src", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"SourceElement", "Srcset", "srcset", "string", nil***REMOVED***,
	***REMOVED***"SourceElement", "Type", "type", "string", nil***REMOVED***,
	***REMOVED***"StyleElement", "Media", "media", "string", nil***REMOVED***,
	***REMOVED***"TableElement", "Sortable", "sortable", "bool", nil***REMOVED***,
	***REMOVED***"TableCellElement", "ColSpan", "colspan", "int", []string***REMOVED***"1"***REMOVED******REMOVED***,
	***REMOVED***"TableCellElement", "RowSpan", "rowspan", "int", []string***REMOVED***"1"***REMOVED******REMOVED***,
	***REMOVED***"TableCellElement", "Headers", "headers", "string", nil***REMOVED***,
	***REMOVED***"TableHeaderCellElement", "Abbr", "abbr", "string", nil***REMOVED***,
	***REMOVED***"TableHeaderCellElement", "Scope", "scope", "enum", []string***REMOVED***"", "row", "col", "colgroup", "rowgroup"***REMOVED******REMOVED***,
	***REMOVED***"TableHeaderCellElement", "Sorted", "sorted", "bool", nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Type", "type", "const", []string***REMOVED***"textarea"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "Value", "value", "string", nil***REMOVED***,
	***REMOVED***"TextAreaElement", "DefaultValue", "value", "string", nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Placeholder", "placeholder", "string", nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Rows", "rows", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "Cols", "cols", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "MaxLength", "maxlength", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "TabIndex", "tabindex", "int", []string***REMOVED***"0"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "AccessKey", "accesskey", "string", nil***REMOVED***,
	***REMOVED***"TextAreaElement", "ReadOnly", "readonly", "bool", nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Required", "required", "bool", nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Autocomplete", "autocomplete", "enum", []string***REMOVED***"on", "off"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "Autocapitalize", "autocapitalize", "enum", []string***REMOVED***"sentences", "none", "off", "characters", "words"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "Wrap", "wrap", "enum", []string***REMOVED***"soft", "hard", "off"***REMOVED******REMOVED***,
	***REMOVED***"TimeElement", "Datetime", "datetime", "string", nil***REMOVED***,
	***REMOVED***"TrackElement", "Kind", "kind", "enum", []string***REMOVED***"subtitle", "captions", "descriptions", "chapters", "metadata"***REMOVED******REMOVED***,
	***REMOVED***"TrackElement", "Src", "src", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***,
	***REMOVED***"TrackElement", "Srclang", "srclang", "string", nil***REMOVED***,
	***REMOVED***"TrackElement", "Label", "label", "string", nil***REMOVED***,
	***REMOVED***"TrackElement", "Default", "default", "bool", nil***REMOVED***,
	***REMOVED***"UListElement", "Type", "type", "string", nil***REMOVED***,
***REMOVED***

var collector = &ElemInfoCollectorState***REMOVED******REMOVED***

func main() ***REMOVED***
	fs := token.NewFileSet()
	parsedFile, parseErr := parser.ParseFile(fs, "elements.go", nil, 0)
	if parseErr != nil ***REMOVED***
		log.Fatalf("error: could not parse elements.go\n", parseErr)
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

	buf := new(bytes.Buffer)
	err := elemFuncsTemplate.Execute(buf, struct ***REMOVED***
		ElemInfos map[string]*ElemInfo
		FuncDefs  []struct ***REMOVED***
			Elem, Method, Attr, TemplateType string
			Opts                             []string
		***REMOVED***
	***REMOVED******REMOVED***
		collector.elemInfos,
		funcDefs,
	***REMOVED***)
	if err != nil ***REMOVED***
		log.Fatalf("error: unable to execute template\n", err)
	***REMOVED***

	src, err := format.Source(buf.Bytes())
	if err != nil ***REMOVED***
		log.Fatalf("error: format.Source on generated code failed\n", err)
	***REMOVED***

	f, err := os.Create("elements_gen.go")
	if err != nil ***REMOVED***
		log.Fatalf("error: Unable to create the file 'elements_gen.go'\n", err)
	***REMOVED***

	if _, err = f.Write(src); err != nil ***REMOVED***
		log.Fatalf("error: Unable to write to 'elements_gen.go'\n", err)
	***REMOVED***

	err = f.Close()
	if err != nil ***REMOVED***
		log.Fatalf("error: unable to close the file 'elements_gen.go'\n", err)
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

***REMOVED******REMOVED*** range $funcDef := .FuncDefs -***REMOVED******REMOVED*** 

func (e ***REMOVED******REMOVED***$funcDef.Elem***REMOVED******REMOVED***) ***REMOVED******REMOVED***$funcDef.Method***REMOVED******REMOVED***() ***REMOVED******REMOVED*** returnType $funcDef.TemplateType ***REMOVED******REMOVED*** ***REMOVED***
***REMOVED******REMOVED***- if eq $funcDef.TemplateType "int" ***REMOVED******REMOVED***
	return e.attrAsInt("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***", ***REMOVED******REMOVED*** index $funcDef.Opts 0 ***REMOVED******REMOVED***)
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType "enum" ***REMOVED******REMOVED***
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
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType "enum-goja" ***REMOVED******REMOVED***
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
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType "const" ***REMOVED******REMOVED***
	return "***REMOVED******REMOVED*** index $funcDef.Opts 0 ***REMOVED******REMOVED***"
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType "url" ***REMOVED******REMOVED***
	return e.attrAsURLString("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***", ***REMOVED******REMOVED*** index $funcDef.Opts 0 ***REMOVED******REMOVED***)
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType "string" ***REMOVED******REMOVED***
	return e.attrAsString("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***")
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType "bool" ***REMOVED******REMOVED***
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
	case "bool":
		return "bool"
	case "int":
		return "int"
	case "enum-goja":
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
