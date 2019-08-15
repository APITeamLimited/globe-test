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

	"github.com/sirupsen/logrus"
)

// Generate elements_gen.go. There are two sections of code which need to be generated. The selToElement function and the attribute accessor methods.

// Thhe first step to generate the selToElement function is parse the TagName constants and Element structs in elements.go using ast.Inspect
// One of NodeHandlerFunc methods is called for each ast.Node parsed by ast.Inspect
// The NodeHandlerFunc methods build ElemInfo structs and populate elemInfos in AstInspectState
// The template later iterates over elemInfos to build the selToElement function

type NodeHandlerFunc func(node ast.Node) NodeHandlerFunc

type AstInspectState struct ***REMOVED***
	handler   NodeHandlerFunc
	elemName  string
	elemInfos map[string]*ElemInfo
***REMOVED***

type ElemInfo struct ***REMOVED***
	StructName    string
	PrtStructName string
***REMOVED***

// The attribute accessors are build using function definitions. Each funcion definition has a TemplateType.
// The number of TemplateArgs varies based on the TenplateType and is documented below.
type TemplateType string
type TemplateArg string

const (
	stringTemplate       TemplateType = "typeString"
	urlTemplate          TemplateType = "typeUrl"
	boolTemplate         TemplateType = "typeBool"
	intTemplate          TemplateType = "typeInt"
	constTemplate        TemplateType = "typeConst"
	enumTemplate         TemplateType = "typeEnum"
	nullableEnumTemplate TemplateType = "typeEnumNullable"
)

// Some common TemplateArgs
var (
	// Default return values for urlTemplate functions. Either an empty string or the current URL.
	defaultURLEmpty   = []TemplateArg***REMOVED***"\"\""***REMOVED***
	defaultURLCurrent = []TemplateArg***REMOVED***"e.sel.URL"***REMOVED***

	// Common default return values for intTemplates
	defaultInt0      = []TemplateArg***REMOVED***"0"***REMOVED***
	defaultIntMinus1 = []TemplateArg***REMOVED***"-1"***REMOVED***
	defaultIntPlus1  = []TemplateArg***REMOVED***"1"***REMOVED***

	// The following are the for various attributes using enumTemplate.
	// The first item in the list is the default value.
	autocompleteOpts = []TemplateArg***REMOVED***"on", "off"***REMOVED***
	referrerOpts     = []TemplateArg***REMOVED***"", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "unsafe-url"***REMOVED***
	preloadOpts      = []TemplateArg***REMOVED***"auto", "metadata", "none"***REMOVED***
	btnTypeOpts      = []TemplateArg***REMOVED***"submit", "button", "menu", "reset"***REMOVED***
	encTypeOpts      = []TemplateArg***REMOVED***"application/x-www-form-urlencoded", "multipart/form-data", "text/plain"***REMOVED***
	inputTypeOpts    = []TemplateArg***REMOVED***"text", "button", "checkbox", "color", "date", "datetime-local", "email", "file", "hidden", "image", "month", "number", "password", "radio", "range", "reset", "search", "submit", "tel", "time", "url", "week"***REMOVED***
	keyTypeOpts      = []TemplateArg***REMOVED***"RSA", "DSA", "EC"***REMOVED***
	keygenTypeOpts   = []TemplateArg***REMOVED***"keygen"***REMOVED***
	liTypeOpts       = []TemplateArg***REMOVED***"", "1", "a", "A", "i", "I", "disc", "square", "circle"***REMOVED***
	httpEquivOpts    = []TemplateArg***REMOVED***"content-type", "default-style", "refresh"***REMOVED***
	olistTypeOpts    = []TemplateArg***REMOVED***"1", "a", "A", "i", "I"***REMOVED***
	scopeOpts        = []TemplateArg***REMOVED***"", "row", "col", "colgroup", "rowgroup"***REMOVED***
	autocapOpts      = []TemplateArg***REMOVED***"sentences", "none", "off", "characters", "words"***REMOVED***
	wrapOpts         = []TemplateArg***REMOVED***"soft", "hard", "off"***REMOVED***
	kindOpts         = []TemplateArg***REMOVED***"subtitle", "captions", "descriptions", "chapters", "metadata"***REMOVED***

	// These are the values allowed for the crossorigin attribute, used by the nullableEnumTemplates is always goja.Undefined
	crossOriginOpts = []TemplateArg***REMOVED***"anonymous", "use-credentials"***REMOVED***
)

// Elem is one of the Element struct names from elements.go
// Method is the go method name to be generated.
// Attr is the name of the DOM attribute the method will access, usually the Method name but lowercased.
// TemplateType determines which type of function is generation by the template
// TemplateArgs is a list of values to be interpolated in the template.

// The number of TemplateArgs depends on the template type.
//   stringTemplate: doesn't use any TemplateArgs
//   boolTemplate: doesn't use any TemplateArgs
//   constTemplate: uses 1 Template Arg, the generated function always returns that value
//   intTemplate: needs 1 TemplateArg, used as the default return value (when the attribute was empty).
//   urlTemplate: needs 1 TemplateArg, used as the default, either "defaultURLEmpty" or "defaultURLCurrent"
//   enumTemplate: uses any number or more TemplateArg, the gen'd func always returns one of the values in the TemplateArgs.
//                 The first item in the list is used as the default when the attribute was invalid or unset.
//   nullableEnumTemplate: similar to the enumTemplate except the default is goja.Undefined and the return type is goja.Value
var funcDefs = []struct ***REMOVED***
	Elem, Method, Attr string
	TemplateType       TemplateType
	TemplateArgs       []TemplateArg
***REMOVED******REMOVED***
	***REMOVED***"HrefElement", "Download", "download", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "ReferrerPolicy", "referrerpolicy", enumTemplate, referrerOpts***REMOVED***,
	***REMOVED***"HrefElement", "Rel", "rel", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "Href", "href", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"HrefElement", "Target", "target", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "AccessKey", "accesskey", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "HrefLang", "hreflang", stringTemplate, nil***REMOVED***,
	***REMOVED***"HrefElement", "ToString", "href", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"MediaElement", "Autoplay", "autoplay", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "Controls", "controls", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "Loop", "loop", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "Muted", "muted", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "Preload", "preload", enumTemplate, preloadOpts***REMOVED***,
	***REMOVED***"MediaElement", "Src", "src", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"MediaElement", "CrossOrigin", "crossorigin", nullableEnumTemplate, crossOriginOpts***REMOVED***,
	***REMOVED***"MediaElement", "CurrentSrc", "src", stringTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "DefaultMuted", "muted", boolTemplate, nil***REMOVED***,
	***REMOVED***"MediaElement", "MediaGroup", "mediagroup", stringTemplate, nil***REMOVED***,
	***REMOVED***"BaseElement", "Href", "href", urlTemplate, defaultURLCurrent***REMOVED***,
	***REMOVED***"BaseElement", "Target", "target", stringTemplate, nil***REMOVED***,
	***REMOVED***"ButtonElement", "AccessKey", "accesskey", stringTemplate, nil***REMOVED***,
	***REMOVED***"ButtonElement", "Autofocus", "autofocus", boolTemplate, nil***REMOVED***,
	***REMOVED***"ButtonElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"ButtonElement", "TabIndex", "tabindex", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"ButtonElement", "Type", "type", enumTemplate, btnTypeOpts***REMOVED***,
	***REMOVED***"DataElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"EmbedElement", "Height", "height", stringTemplate, nil***REMOVED***,
	***REMOVED***"EmbedElement", "Width", "width", stringTemplate, nil***REMOVED***,
	***REMOVED***"EmbedElement", "Src", "src", stringTemplate, nil***REMOVED***,
	***REMOVED***"EmbedElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"FieldSetElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"FieldSetElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"FormElement", "Action", "action", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"FormElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"FormElement", "Target", "target", stringTemplate, nil***REMOVED***,
	***REMOVED***"FormElement", "Enctype", "enctype", enumTemplate, encTypeOpts***REMOVED***,
	***REMOVED***"FormElement", "Encoding", "enctype", enumTemplate, encTypeOpts***REMOVED***,
	***REMOVED***"FormElement", "AcceptCharset", "accept-charset", stringTemplate, nil***REMOVED***,
	***REMOVED***"FormElement", "Autocomplete", "autocomplete", enumTemplate, autocompleteOpts***REMOVED***,
	***REMOVED***"FormElement", "NoValidate", "novalidate", boolTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "Allowfullscreen", "allowfullscreen", boolTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "ReferrerPolicy", "referrerpolicy", enumTemplate, referrerOpts***REMOVED***,
	***REMOVED***"IFrameElement", "Height", "height", stringTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "Width", "width", stringTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"IFrameElement", "Src", "src", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"ImageElement", "CurrentSrc", "src", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"ImageElement", "Sizes", "sizes", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "Srcset", "srcset", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "Alt", "alt", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "CrossOrigin", "crossorigin", nullableEnumTemplate, crossOriginOpts***REMOVED***,
	***REMOVED***"ImageElement", "Height", "height", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"ImageElement", "Width", "width", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"ImageElement", "IsMap", "ismap", boolTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "Src", "src", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"ImageElement", "UseMap", "usemap", stringTemplate, nil***REMOVED***,
	***REMOVED***"ImageElement", "ReferrerPolicy", "referrerpolicy", enumTemplate, referrerOpts***REMOVED***,
	***REMOVED***"InputElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "TabIndex", "tabindex", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"InputElement", "Type", "type", enumTemplate, inputTypeOpts***REMOVED***,
	***REMOVED***"InputElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Autofocus", "autofocus", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Required", "required", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Checked", "checked", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "DefaultChecked", "checked", boolTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Alt", "alt", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Src", "src", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"InputElement", "Height", "height", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Width", "width", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Accept", "accept", stringTemplate, nil***REMOVED***,
	***REMOVED***"InputElement", "Autocomplete", "autocomplete", enumTemplate, autocompleteOpts***REMOVED***,
	***REMOVED***"InputElement", "MaxLength", "maxlength", intTemplate, defaultIntMinus1***REMOVED***,
	***REMOVED***"InputElement", "Size", "size", intTemplate, defaultInt0***REMOVED***,
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
	***REMOVED***"KeygenElement", "Keytype", "keytype", enumTemplate, keyTypeOpts***REMOVED***,
	***REMOVED***"KeygenElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"KeygenElement", "Type", "type", constTemplate, keygenTypeOpts***REMOVED***,
	***REMOVED***"LabelElement", "HtmlFor", "for", stringTemplate, nil***REMOVED***,
	***REMOVED***"LegendElement", "AccessKey", "accesskey", stringTemplate, nil***REMOVED***,
	***REMOVED***"LiElement", "Value", "value", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"LiElement", "Type", "type", enumTemplate, liTypeOpts***REMOVED***,
	***REMOVED***"LinkElement", "CrossOrigin", "crossorigin", nullableEnumTemplate, crossOriginOpts***REMOVED***,
	***REMOVED***"LinkElement", "ReferrerPolicy", "referrerpolicy", enumTemplate, referrerOpts***REMOVED***,
	***REMOVED***"LinkElement", "Href", "href", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"LinkElement", "Hreflang", "hreflang", stringTemplate, nil***REMOVED***,
	***REMOVED***"LinkElement", "Media", "media", stringTemplate, nil***REMOVED***,
	***REMOVED***"LinkElement", "Rel", "rel", stringTemplate, nil***REMOVED***,
	***REMOVED***"LinkElement", "Target", "target", stringTemplate, nil***REMOVED***,
	***REMOVED***"LinkElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"MapElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"MetaElement", "Content", "content", stringTemplate, nil***REMOVED***,
	***REMOVED***"MetaElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"MetaElement", "HttpEquiv", "http-equiv", enumTemplate, httpEquivOpts***REMOVED***,
	***REMOVED***"MeterElement", "Min", "min", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"MeterElement", "Max", "max", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"MeterElement", "High", "high", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"MeterElement", "Low", "low", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"MeterElement", "Optimum", "optimum", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"ModElement", "Cite", "cite", stringTemplate, nil***REMOVED***,
	***REMOVED***"ModElement", "Datetime", "datetime", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "Data", "data", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"ObjectElement", "Height", "height", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "TabIndex", "tabindex", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"ObjectElement", "TypeMustMatch", "typemustmatch", boolTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "UseMap", "usemap", stringTemplate, nil***REMOVED***,
	***REMOVED***"ObjectElement", "Width", "width", stringTemplate, nil***REMOVED***,
	***REMOVED***"OListElement", "Reversed", "reversed", boolTemplate, nil***REMOVED***,
	***REMOVED***"OListElement", "Start", "start", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"OListElement", "Type", "type", enumTemplate, olistTypeOpts***REMOVED***,
	***REMOVED***"OptGroupElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"OptGroupElement", "Label", "label", stringTemplate, nil***REMOVED***,
	***REMOVED***"OptionElement", "DefaultSelected", "selected", boolTemplate, nil***REMOVED***,
	***REMOVED***"OptionElement", "Selected", "selected", boolTemplate, nil***REMOVED***,
	***REMOVED***"OutputElement", "HtmlFor", "for", stringTemplate, nil***REMOVED***,
	***REMOVED***"OutputElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"OutputElement", "Type", "type", constTemplate, []TemplateArg***REMOVED***"output"***REMOVED******REMOVED***,
	***REMOVED***"ParamElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"ParamElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"PreElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"PreElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"QuoteElement", "Cite", "cite", stringTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "CrossOrigin", "crossorigin", stringTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "Src", "src", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"ScriptElement", "Charset", "charset", stringTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "Async", "async", boolTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "Defer", "defer", boolTemplate, nil***REMOVED***,
	***REMOVED***"ScriptElement", "NoModule", "nomodule", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Autofocus", "autofocus", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Disabled", "disabled", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Multiple", "multiple", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Name", "name", stringTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "Required", "required", boolTemplate, nil***REMOVED***,
	***REMOVED***"SelectElement", "TabIndex", "tabindex", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"SourceElement", "KeySystem", "keysystem", stringTemplate, nil***REMOVED***,
	***REMOVED***"SourceElement", "Media", "media", stringTemplate, nil***REMOVED***,
	***REMOVED***"SourceElement", "Sizes", "sizes", stringTemplate, nil***REMOVED***,
	***REMOVED***"SourceElement", "Src", "src", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"SourceElement", "Srcset", "srcset", stringTemplate, nil***REMOVED***,
	***REMOVED***"SourceElement", "Type", "type", stringTemplate, nil***REMOVED***,
	***REMOVED***"StyleElement", "Media", "media", stringTemplate, nil***REMOVED***,
	***REMOVED***"TableElement", "Sortable", "sortable", boolTemplate, nil***REMOVED***,
	***REMOVED***"TableCellElement", "ColSpan", "colspan", intTemplate, defaultIntPlus1***REMOVED***,
	***REMOVED***"TableCellElement", "RowSpan", "rowspan", intTemplate, defaultIntPlus1***REMOVED***,
	***REMOVED***"TableCellElement", "Headers", "headers", stringTemplate, nil***REMOVED***,
	***REMOVED***"TableHeaderCellElement", "Abbr", "abbr", stringTemplate, nil***REMOVED***,
	***REMOVED***"TableHeaderCellElement", "Scope", "scope", enumTemplate, scopeOpts***REMOVED***,
	***REMOVED***"TableHeaderCellElement", "Sorted", "sorted", boolTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Type", "type", constTemplate, []TemplateArg***REMOVED***"textarea"***REMOVED******REMOVED***,
	***REMOVED***"TextAreaElement", "Value", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "DefaultValue", "value", stringTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Placeholder", "placeholder", stringTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Rows", "rows", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"TextAreaElement", "Cols", "cols", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"TextAreaElement", "MaxLength", "maxlength", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"TextAreaElement", "TabIndex", "tabindex", intTemplate, defaultInt0***REMOVED***,
	***REMOVED***"TextAreaElement", "AccessKey", "accesskey", stringTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "ReadOnly", "readonly", boolTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Required", "required", boolTemplate, nil***REMOVED***,
	***REMOVED***"TextAreaElement", "Autocomplete", "autocomplete", enumTemplate, autocompleteOpts***REMOVED***,
	***REMOVED***"TextAreaElement", "Autocapitalize", "autocapitalize", enumTemplate, autocapOpts***REMOVED***,
	***REMOVED***"TextAreaElement", "Wrap", "wrap", enumTemplate, wrapOpts***REMOVED***,
	***REMOVED***"TimeElement", "Datetime", "datetime", stringTemplate, nil***REMOVED***,
	***REMOVED***"TrackElement", "Kind", "kind", enumTemplate, kindOpts***REMOVED***,
	***REMOVED***"TrackElement", "Src", "src", urlTemplate, defaultURLEmpty***REMOVED***,
	***REMOVED***"TrackElement", "Srclang", "srclang", stringTemplate, nil***REMOVED***,
	***REMOVED***"TrackElement", "Label", "label", stringTemplate, nil***REMOVED***,
	***REMOVED***"TrackElement", "Default", "default", boolTemplate, nil***REMOVED***,
	***REMOVED***"UListElement", "Type", "type", stringTemplate, nil***REMOVED***,
***REMOVED***

func main() ***REMOVED***
	fs := token.NewFileSet()
	parsedFile, parseErr := parser.ParseFile(fs, "elements.go", nil, 0)
	if parseErr != nil ***REMOVED***
		logrus.WithError(parseErr).Fatal("Could not parse elements.go")
	***REMOVED***

	// Initialise the AstInspectState
	var collector = &AstInspectState***REMOVED******REMOVED***

	collector.handler = collector.defaultHandler
	collector.elemInfos = make(map[string]*ElemInfo)

	// Populate collector.elemInfos
	ast.Inspect(parsedFile, func(n ast.Node) bool ***REMOVED***
		if n != nil ***REMOVED***
			collector.handler = collector.handler(n)
		***REMOVED***
		return true
	***REMOVED***)

	// elemInfos and funcDefs are now complete and the template can be executed.
	var buf bytes.Buffer
	err := elemFuncsTemplate.Execute(&buf, struct ***REMOVED***
		ElemInfos map[string]*ElemInfo
		FuncDefs  []struct ***REMOVED***
			Elem, Method, Attr string
			TemplateType       TemplateType
			TemplateArgs       []TemplateArg
		***REMOVED***
		TemplateTypes struct***REMOVED*** String, Url, Enum, Bool, GojaEnum, Int, Const TemplateType ***REMOVED***
	***REMOVED******REMOVED***
		collector.elemInfos,
		funcDefs,
		struct***REMOVED*** String, Url, Enum, Bool, GojaEnum, Int, Const TemplateType ***REMOVED******REMOVED***stringTemplate, urlTemplate, enumTemplate, boolTemplate, nullableEnumTemplate, intTemplate, constTemplate***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		logrus.WithError(err).Fatal("Unable to execute template")
	***REMOVED***

	src, err := format.Source(buf.Bytes())
	if err != nil ***REMOVED***
		logrus.WithError(err).Fatal("format.Source on generated code failed")
	***REMOVED***

	f, err := os.Create("elements_gen.go")
	if err != nil ***REMOVED***
		logrus.WithError(err).Fatal("Unable to create the file 'elements_gen.go'")
	***REMOVED***

	if _, err = f.Write(src); err != nil ***REMOVED***
		logrus.WithError(err).Fatal("Unable to write to 'elements_gen.go'")
	***REMOVED***

	err = f.Close()
	if err != nil ***REMOVED***
		logrus.WithError(err).Fatal("Unable to close 'elements_gen.go'")
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

***REMOVED******REMOVED*** $templateTypes := .TemplateTypes ***REMOVED******REMOVED***
***REMOVED******REMOVED*** range $funcDef := .FuncDefs -***REMOVED******REMOVED*** 

func (e ***REMOVED******REMOVED***$funcDef.Elem***REMOVED******REMOVED***) ***REMOVED******REMOVED***$funcDef.Method***REMOVED******REMOVED***() ***REMOVED******REMOVED*** returnType $funcDef.TemplateType ***REMOVED******REMOVED*** ***REMOVED***
***REMOVED******REMOVED***- if eq $funcDef.TemplateType $templateTypes.Int ***REMOVED******REMOVED***
	return e.attrAsInt("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***", ***REMOVED******REMOVED*** index $funcDef.TemplateArgs 0 ***REMOVED******REMOVED***)
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateTypes.Enum ***REMOVED******REMOVED***
	attrVal := e.attrAsString("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***")
	switch attrVal ***REMOVED*** 
	***REMOVED******REMOVED***- range $optIdx, $optVal := $funcDef.TemplateArgs ***REMOVED******REMOVED***
	***REMOVED******REMOVED***- if ne $optIdx 0 ***REMOVED******REMOVED***
	case "***REMOVED******REMOVED***$optVal***REMOVED******REMOVED***":
		return attrVal
	***REMOVED******REMOVED***- end ***REMOVED******REMOVED***
	***REMOVED******REMOVED***- end***REMOVED******REMOVED***
	default: 
		return "***REMOVED******REMOVED*** index $funcDef.TemplateArgs 0 ***REMOVED******REMOVED***" 
	***REMOVED***
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateTypes.GojaEnum ***REMOVED******REMOVED***
	attrVal, exists := e.sel.sel.Attr("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***")
	if !exists ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	switch attrVal ***REMOVED*** 
	***REMOVED******REMOVED***- range $optVal := $funcDef.TemplateArgs ***REMOVED******REMOVED***
	case "***REMOVED******REMOVED***$optVal***REMOVED******REMOVED***":
		return e.sel.rt.ToValue(attrVal)
	***REMOVED******REMOVED***- end***REMOVED******REMOVED***
	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateTypes.Const ***REMOVED******REMOVED***
	return "***REMOVED******REMOVED*** index $funcDef.TemplateArgs 0 ***REMOVED******REMOVED***"
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateTypes.Url ***REMOVED******REMOVED***
	return e.attrAsURLString("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***", ***REMOVED******REMOVED*** index $funcDef.TemplateArgs 0 ***REMOVED******REMOVED***)
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateTypes.String ***REMOVED******REMOVED***
	return e.attrAsString("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***")
***REMOVED******REMOVED***- else if eq $funcDef.TemplateType $templateTypes.Bool ***REMOVED******REMOVED***
	return e.attrIsPresent("***REMOVED******REMOVED*** $funcDef.Attr ***REMOVED******REMOVED***")
***REMOVED******REMOVED***- end***REMOVED******REMOVED***
***REMOVED***
***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
`))

// generate the nested struct, either one or two levels of nesting, ie "BaseElement***REMOVED***elem***REMOVED***" or "ButtonElement***REMOVED***FormFieldElement***REMOVED***elem***REMOVED******REMOVED***)"
func buildStruct(elemInfo ElemInfo) string ***REMOVED***
	if elemInfo.PrtStructName == "Element" ***REMOVED***
		return elemInfo.StructName + "***REMOVED***elem***REMOVED***"
	***REMOVED*** else ***REMOVED***
		return elemInfo.StructName + "***REMOVED***" + elemInfo.PrtStructName + "***REMOVED***elem***REMOVED******REMOVED***"
	***REMOVED***
***REMOVED***

// Select the correct return type for one of the attribute accessor methods
func returnType(templateType TemplateType) string ***REMOVED***
	switch templateType ***REMOVED***
	case boolTemplate:
		return "bool"
	case intTemplate:
		return "int"
	case nullableEnumTemplate:
		return "goja.Value"
	default:
		return "string"
	***REMOVED***
***REMOVED***

// Default node handler functions for ast.Inspect. Return itself unless it's found a "const" or "struct" keyword
func (ce *AstInspectState) defaultHandler(node ast.Node) NodeHandlerFunc ***REMOVED***
	ce.elemName = ""
	switch node.(type) ***REMOVED***
	case *ast.TypeSpec: // struct keyword
		return ce.elementStructHandler

	case *ast.ValueSpec: // const keyword
		return ce.elementTagNameHandler

	default:
		return ce.defaultHandler
	***REMOVED***
***REMOVED***

// Found a tagname constant. The code 'const AnchorTagName = "a"' will add an ElemInfo called "Anchor", like elemInfos["Anchor"] = ElemInfo***REMOVED***"", ""***REMOVED***
func (ce *AstInspectState) elementTagNameHandler(node ast.Node) NodeHandlerFunc ***REMOVED***
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

// A struct definition was found, keep the elem handler if it's for an Element struct
// Element structs nest the "Element" struct or an intermediate struct like "HrefElement", the name of the 'parent' struct is contained in the
// *ast.Ident node located a few nodes after the TypeSpec node containing struct keyword
// The nodes in between the ast.TypeSpec and ast.Ident are ignored
func (ce *AstInspectState) elementStructHandler(node ast.Node) NodeHandlerFunc ***REMOVED***
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
			return ce.elementStructHandler
		***REMOVED*** else ***REMOVED***
			ce.elemInfos[ce.elemName].PrtStructName = x.Name
			return ce.defaultHandler
		***REMOVED***

	case *ast.StructType:
		return ce.elementStructHandler

	case *ast.FieldList:
		return ce.elementStructHandler

	case *ast.Field:
		return ce.elementStructHandler

	default:
		return ce.defaultHandler
	***REMOVED***
***REMOVED***
