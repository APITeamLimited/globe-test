// +build ignore

package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

// The ast parser populates this struct using the <ElemName>TagName const and <ElemName>Element struct
type ElemInfo struct ***REMOVED***
	ConstName     string
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

// Built from an entry in funcDefs by buildFuncDef(funcDef)
type FuncDef struct ***REMOVED***
	ElemName   string
	ElemMethod string
	AttrMethod string
	AttrName   string
	ReturnType string
	ReturnBody string
	ReturnOpts []string
***REMOVED***

var funcDefs = []string***REMOVED***
	"Href Download string",
	"Href ReferrerPolicy enum=,no-referrer,no-referrer-when-downgrade,origin,origin-when-cross-origin,unsafe-url",
	"Href Rel string",
	"Href Href url-or-blank",
	"Href Target string",
	"Href Type string",
	"Href AccessKey string",
	"Href HrefLang string",
	"Href ToString=href url-or-blank",
	"Media Autoplay bool",
	"Media Controls bool",
	"Media Loop bool",
	"Media Muted bool",
	"Media Preload enum=auto,metadata,none",
	"Media Src url-or-blank",
	"Media CrossOrigin enum-nullable=anonymous,use-credentials",
	"Media CurrentSrc=src string",
	"Media DefaultMuted=muted bool",
	"Media MediaGroup string",
	"Base Href url-or-currpage",
	"Base Target string",
	"Button AccessKey string",
	"Button Autofocus bool",
	"Button Disabled bool",
	"Button TabIndex int",
	"Button Type enum=submit,button,menu,reset",
	"Data Value string",
	"Embed Height string",
	"Embed Width string",
	"Embed Src string",
	"Embed Type string",
	"FieldSet Disabled bool",
	"FieldSet Name string",
	"Form Action url-or-blank",
	"Form Name string",
	"Form Target string",
	"Form Enctype enum=application/x-www-form-urlencoded,multipart/form-data,text/plain",
	"Form Encoding=enctype enum=application/x-www-form-urlencoded,multipart/form-data,text/plain",
	"Form AcceptCharset=accept-charset string",
	"Form Autocomplete enum=on,off",
	"Form NoValidate bool",
	"IFrame Allowfullscreen bool",
	"IFrame ReferrerPolicy enum=,no-referrer,no-referrer-when-downgrade,origin,origin-when-cross-origin,unsafe-url",
	"IFrame Height string",
	"IFrame Width string",
	"IFrame Name string",
	"IFrame Src url-or-blank",
	"Image CurrentSrc=src url-or-blank",
	"Image Sizes string",
	"Image Srcset string",
	"Image Alt string",
	"Image CrossOrigin enum-nullable=anonymous,use-credentials",
	"Image Height int",
	"Image Width int",
	"Image IsMap bool",
	"Image Name string",
	"Image Src url-or-blank",
	"Image UseMap string",
	"Image ReferrerPolicy enum=,no-referrer,no-referrer-when-downgrade,origin,origin-when-cross-origin,unsafe-url",
	"Input Name string",
	"Input TabIndex int",
	"Input Type enum=text,button,checkbox,color,date,datetime-local,email,file,hidden,image,month,number,password,radio,range,reset,search,submit,tel,time,url,week",
	"Input Disabled bool",
	"Input Autofocus bool",
	"Input Required bool",
	"Input Value string",
	"Input Checked bool",
	"Input DefaultChecked=checked bool",
	"Input Alt string",
	"Input Src url-or-blank",
	"Input Height string",
	"Input Width string",
	"Input Accept string",
	"Input Autocomplete enum=on,off",
	"Input MaxLength int=-1",
	"Input Size int",
	"Input Pattern string",
	"Input Placeholder string",
	"Input Readonly bool",
	"Input Min string",
	"Input Max string",
	"Input DefaultValue=value string",
	"Input DirName string",
	"Input AccessKey string",
	"Input Multiple bool",
	"Input Step string",
	"Keygen Autofocus bool",
	"Keygen Challenge string",
	"Keygen Disabled bool",
	"Keygen Keytype enum=RSA,DSA,EC",
	"Keygen Name string",
	"Keygen Type const=keygen",
	"Label HtmlFor=for string",
	"Legend AccessKey string",
	"Li Value int=0",
	"Li Type enum=,1,a,A,i,I,disc,square,circle",
	"Link CrossOrigin enum-nullable=anonymous,use-credentials",
	"Link ReferrerPolicy enum=,no-referrer,no-referrer-when-downgrade,origin,origin-when-cross-origin,unsafe-url",
	"Link Href url-or-blank",
	"Link Hreflang string",
	"Link Media string",
	"Link Rel string",
	"Link Target string",
	"Link Type string",
	"Map Name string",
	"Meta Content string",
	"Meta Name string",
	"Meta HttpEquiv=http-equiv enum=content-type,default-style,refresh",
	"Meter Min int",
	"Meter Max int",
	"Meter High int",
	"Meter Low int",
	"Meter Optimum int",
	"Mod Cite string",
	"Mod Datetime string",
	"Object Data url-or-blank",
	"Object Height string",
	"Object Name string",
	"Object Type string",
	"Object TabIndex int=0",
	"Object TypeMustMatch bool",
	"Object UseMap string",
	"Object Width string",
	"OList Reversed bool",
	"OList Start int",
	"OList Type enum=1,a,A,i,I",
	"OptGroup Disabled bool",
	"OptGroup Label string",
	"Option DefaultSelected=selected bool",
	"Option Selected bool",
	"Output HtmlFor=for string",
	"Output Name string",
	"Output Type const=output",
	"Param Name string",
	"Param Value string",
	"Pre Name string",
	"Pre Value string",
	"Quote Cite string",
	"Script CrossOrigin string",
	"Script Type string",
	"Script Src url-or-blank",
	"Script Charset string",
	"Script Async bool",
	"Script Defer bool",
	"Script NoModule bool",
	"Select Autofocus bool",
	"Select Disabled bool",
	"Select Multiple bool",
	"Select Name string",
	"Select Required bool",
	"Select TabIndex int",
	"Source KeySystem string",
	"Source Media string",
	"Source Sizes string",
	"Source Src url-or-blank",
	"Source Srcset string",
	"Source Type string",
	"Style Media string",
	"Table Sortable bool",
	"TableCell ColSpan int=1",
	"TableCell RowSpan int=1",
	"TableCell Headers string",
	"TableHeaderCell Abbr string",
	"TableHeaderCell Scope enum=,row,col,colgroup,rowgroup",
	"TableHeaderCell Sorted bool",
	"TextArea Type const=textarea",
	"TextArea Value string",
	"TextArea DefaultValue=value string",
	"TextArea Placeholder string",
	"TextArea Rows int",
	"TextArea Cols int",
	"TextArea MaxLength int",
	"TextArea TabIndex int",
	"TextArea AccessKey string",
	"TextArea ReadOnly bool",
	"TextArea Required bool",
	"TextArea Autocomplete enum=on,off",
	"TextArea Autocapitalize enum=sentences,none,off,characters,words",
	"TextArea Wrap enum=soft,hard,off",
	"Time Datetime string",
	"Track Kind enum=subtitle,captions,descriptions,chapters,metadata",
	"Track Src url-or-blank",
	"Track Srclang string",
	"Track Label string",
	"Track Default bool",
	"UList Type string",
***REMOVED***

var collector = &ElemInfoCollectorState***REMOVED******REMOVED***

func main() ***REMOVED***
	fs := token.NewFileSet()
	parsedFile, parseErr := parser.ParseFile(fs, "elements.go", nil, 0)
	if parseErr != nil ***REMOVED***
		log.Fatalf("warning: internal error: could not parse elements.go: %s", parseErr)
		return
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

	f, err := os.Create("elements_gen.go")
	if err != nil ***REMOVED***
		log.Println("warning: Unable to create the file 'elements_gen.go'", err)
	***REMOVED***

	// Generate consts used in the enum items in funcDefs. Consts are needed to avoid warnings for the repeated use of strings
	var enumConsts = map[string]string***REMOVED******REMOVED***
	for _, def := range funcDefs ***REMOVED***
		funcDef := buildFuncDef(def)

		if !strings.HasPrefix(funcDef.ReturnBody, "enum") ***REMOVED***
			continue
		***REMOVED***
		for _, opt := range funcDef.ReturnOpts ***REMOVED***
			enumConsts[opt] = toConst(opt)
		***REMOVED***
	***REMOVED***

	//
	err = elemFuncsTemplate.Execute(f, struct ***REMOVED***
		ElemInfos  map[string]*ElemInfo
		FuncDefs   []string
		EnumConsts map[string]string
	***REMOVED******REMOVED***
		collector.elemInfos,
		funcDefs,
		enumConsts,
	***REMOVED***)
	if err != nil ***REMOVED***
		log.Println("error, unable to execute template:", err)
	***REMOVED***

	err = f.Close()
	if err != nil ***REMOVED***
		log.Println("Unable to close the file 'elements_gen.go': ", err)
	***REMOVED***
***REMOVED***

var elemFuncsTemplate = template.Must(template.New("").Funcs(template.FuncMap***REMOVED***
	"buildStruct":  buildStruct,
	"buildFuncDef": buildFuncDef,
	"toConst":      toConst,
***REMOVED***).Parse(`// go generate
// generated by js/modules/k6/html/gen/main.go directed by js/modules/k6/html/elements.go;  DO NOT EDIT
package html

import "github.com/dop251/goja"

const (
	***REMOVED******REMOVED*** range $constVal, $constName := .EnumConsts -***REMOVED******REMOVED***
		***REMOVED******REMOVED***$constName***REMOVED******REMOVED*** = "***REMOVED******REMOVED***$constVal***REMOVED******REMOVED***"
	***REMOVED******REMOVED*** end -***REMOVED******REMOVED***
)

func selToElement(sel Selection) goja.Value ***REMOVED***
	if sel.sel.Length() == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	elem := Element***REMOVED***sel.sel.Nodes[0], &sel***REMOVED***

	switch elem.node.Data ***REMOVED*** ***REMOVED******REMOVED*** range $elemInfo := .ElemInfos ***REMOVED******REMOVED***
	case ***REMOVED******REMOVED*** $elemInfo.ConstName ***REMOVED******REMOVED***:
		return sel.rt.ToValue(***REMOVED******REMOVED*** buildStruct $elemInfo ***REMOVED******REMOVED***)
***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
	default:
		return sel.rt.ToValue(elem)
	***REMOVED***
 ***REMOVED***

***REMOVED******REMOVED*** range $funcDefStr := .FuncDefs -***REMOVED******REMOVED*** 
***REMOVED******REMOVED*** $funcDef := buildFuncDef $funcDefStr -***REMOVED******REMOVED***
func (e ***REMOVED******REMOVED***$funcDef.ElemName***REMOVED******REMOVED***) ***REMOVED******REMOVED***$funcDef.ElemMethod***REMOVED******REMOVED***() ***REMOVED******REMOVED***$funcDef.ReturnType***REMOVED******REMOVED*** ***REMOVED***
***REMOVED******REMOVED***- if eq $funcDef.ReturnBody "int" ***REMOVED******REMOVED***
	return e.***REMOVED******REMOVED*** $funcDef.AttrMethod ***REMOVED******REMOVED***("***REMOVED******REMOVED*** $funcDef.AttrName ***REMOVED******REMOVED***", ***REMOVED******REMOVED*** index $funcDef.ReturnOpts 0 ***REMOVED******REMOVED***)
***REMOVED******REMOVED***- else if eq $funcDef.ReturnBody "enum" ***REMOVED******REMOVED***
	attrVal := e.attrAsString("***REMOVED******REMOVED*** $funcDef.AttrName ***REMOVED******REMOVED***")
	switch attrVal ***REMOVED*** 
	***REMOVED******REMOVED***- range $optIdx, $optVal := $funcDef.ReturnOpts ***REMOVED******REMOVED***
	***REMOVED******REMOVED***- $optConst := toConst $optVal ***REMOVED******REMOVED***
		***REMOVED******REMOVED***- if ne $optIdx 0 ***REMOVED******REMOVED***
	case ***REMOVED******REMOVED***$optConst***REMOVED******REMOVED***:
		return attrVal
		***REMOVED******REMOVED***- end ***REMOVED******REMOVED***
	***REMOVED******REMOVED***- end***REMOVED******REMOVED***
	default: 
		return ***REMOVED******REMOVED*** index $funcDef.ReturnOpts 0 | toConst ***REMOVED******REMOVED*** 
	***REMOVED***
***REMOVED******REMOVED***- else if eq $funcDef.ReturnBody "enum-nullable" ***REMOVED******REMOVED***
	attrVal, exists := e.sel.sel.Attr("***REMOVED******REMOVED*** $funcDef.AttrName ***REMOVED******REMOVED***")
	if !exists ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	switch attrVal ***REMOVED*** 
	***REMOVED******REMOVED***- range $optVal := $funcDef.ReturnOpts ***REMOVED******REMOVED***
	***REMOVED******REMOVED***- $optConst := toConst $optVal ***REMOVED******REMOVED***
	case ***REMOVED******REMOVED***$optConst***REMOVED******REMOVED***:
		return e.sel.rt.ToValue(attrVal)
	***REMOVED******REMOVED***- end***REMOVED******REMOVED***
	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED******REMOVED***- else if eq $funcDef.ReturnBody "const" ***REMOVED******REMOVED***
	return "***REMOVED******REMOVED*** index $funcDef.ReturnOpts 0 ***REMOVED******REMOVED***"
***REMOVED******REMOVED***- else if eq $funcDef.ReturnBody "url" ***REMOVED******REMOVED***
	return e.***REMOVED******REMOVED*** $funcDef.AttrMethod ***REMOVED******REMOVED***("***REMOVED******REMOVED*** $funcDef.AttrName ***REMOVED******REMOVED***", ***REMOVED******REMOVED*** index $funcDef.ReturnOpts 0 ***REMOVED******REMOVED***)
***REMOVED******REMOVED***- else ***REMOVED******REMOVED***
	return e.***REMOVED******REMOVED*** $funcDef.AttrMethod ***REMOVED******REMOVED***("***REMOVED******REMOVED*** $funcDef.AttrName ***REMOVED******REMOVED***")
***REMOVED******REMOVED***- end***REMOVED******REMOVED***
***REMOVED***
***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
`))

func constNameMapper(r rune) rune ***REMOVED***
	if r == '-' || r == '/' ***REMOVED***
		return '_'
	***REMOVED***
	return r
***REMOVED***

func toConst(optName string) string ***REMOVED***
	if optName == "" ***REMOVED***
		return "const_Blank"
	***REMOVED***
	return "const_" + strings.Map(constNameMapper, optName)
***REMOVED***

func buildStruct(elemInfo ElemInfo) string ***REMOVED***
	if elemInfo.PrtStructName == "Element" ***REMOVED***
		return elemInfo.StructName + "***REMOVED***elem***REMOVED***"
	***REMOVED*** else ***REMOVED***
		return elemInfo.StructName + "***REMOVED***" + elemInfo.PrtStructName + "***REMOVED***elem***REMOVED******REMOVED***"
	***REMOVED***
***REMOVED***

func buildFuncDef(funcDef string) FuncDef ***REMOVED***
	parts := strings.Split(funcDef, " ")
	// parts[0] is the element struct name (without the Element suffix)
	// parts[1] is either:
	//   MethodName               The name of method added onto that struct and converted to lowercase thenn used as the argument to elem.attrAsString(...) or elem.AttrIsPresent(...)
	//   MethodName=attrname      The MethodName is added to the struct. The attrname is the argument for attrAsString or AttrIsPresent
	// parts[2] is the return type, either string, const, bool, int, enum or enum-nullable.
	elemName := parts[0] + "Element"
	elemMethod := parts[1]
	attrName := strings.ToLower(parts[1])
	returnType := parts[2]
	returnOpts := ""

	if eqPos := strings.Index(parts[1], "="); eqPos != -1 ***REMOVED***
		attrName = elemMethod[eqPos+1:]
		elemMethod = elemMethod[0:eqPos]
	***REMOVED***

	if eqPos := strings.Index(returnType, "="); eqPos != -1 ***REMOVED***
		returnOpts = returnType[eqPos+1:]
		returnType = returnType[0:eqPos]
	***REMOVED***

	switch returnType ***REMOVED***
	case "int":
		// The number following 'int=' is a default value used when the attribute is not defined.
		// "TableCell ColSpan int=1"
		// => ***REMOVED***"TableCellElement" "ColSpan", "attrAsInt", "colspan", "int", "int", []string***REMOVED***"1"***REMOVED******REMOVED***
		// => `func (e TableCellElement) ColSpan() int***REMOVED*** return e.attrAsInt("colspan", 1) ***REMOVED***``
		if returnOpts == "" ***REMOVED***
			return FuncDef***REMOVED***elemName, elemMethod, "attrAsInt", attrName, "int", returnType, []string***REMOVED***"0"***REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			return FuncDef***REMOVED***elemName, elemMethod, "attrAsInt", attrName, "int", returnType, []string***REMOVED***returnOpts***REMOVED******REMOVED***
		***REMOVED***

	case "enum":
		// "Button Type enum=submit,button,menu,reset"
		// The items in the comma separated list are used in a switch statement. The first value in the list is used as the default.
		return FuncDef***REMOVED***elemName, elemMethod, "", attrName, "string", returnType, strings.Split(returnOpts, ",")***REMOVED***

	case "enum-nullable":
		// Similar to the enum except the default is goja.Undefined()
		return FuncDef***REMOVED***elemName, elemMethod, "", attrName, "goja.Value", returnType, strings.Split(returnOpts, ",")***REMOVED***

	case "string":
		// "Button AccessKey string"
		// => ***REMOVED***"ButtonElement" "AccessKey", "attrIsString", "accesskey", "string", "string", nil***REMOVED***
		// => `func (e ButtonElement) AccessKey() string***REMOVED*** return e.attrAsString("accessKey") ***REMOVED***``
		return FuncDef***REMOVED***elemName, elemMethod, "attrAsString", attrName, returnType, returnType, nil***REMOVED***

	// This url function uses the current page URL as default when attribute is empty and empty string as default when the attribute is undefined
	case "url-or-blank":
		// "Href Href url-or-blank"
		// => ***REMOVED***"HrefElement" "Href", "attrIsString", "accesskey", "string", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***
		// => `func (e HrefElement) Href() string***REMOVED*** return e.attrAsUrlString("href", """) ***REMOVED***``
		return FuncDef***REMOVED***elemName, elemMethod, "attrAsURLString", attrName, "string", "url", []string***REMOVED***"\"\""***REMOVED******REMOVED***

	// This url function uses current page URL for empty and undefined attributes
	case "url-or-currpage":
		// "Base Href url-or-currpage"
		// => ***REMOVED***"BaseElement" "Href", "attrIsString", "accesskey", "string", "url", []string***REMOVED***"e.sel.URL"***REMOVED******REMOVED***
		// => `func (e HrefElement) Href() string***REMOVED*** return e.attrAsUrlString("href", e.sel.URL) ***REMOVED***``
		return FuncDef***REMOVED***elemName, elemMethod, "attrAsURLString", attrName, "string", "url", []string***REMOVED***"e.sel.URL"***REMOVED******REMOVED***

	case "const":
		// "Output Type const=output"
		// => ***REMOVED***"OutputElement" "Type", "", "type", "string", "const", []***REMOVED***"output"***REMOVED******REMOVED***
		// => `func (e OutputElement) Type() string***REMOVED*** return "output" ***REMOVED***``
		return FuncDef***REMOVED***elemName, elemMethod, "", attrName, "string", returnType, []string***REMOVED***returnOpts***REMOVED******REMOVED***

	case "bool":
		// "Button Autofocus bool"
		// => ***REMOVED***"Button" "Autofocus", "attrIsPresent", "autofocus", "bool", "bool", nil***REMOVED***
		// => `func (e ButtonElement) ToString() bool ***REMOVED*** return e.attrIsPresent("autofocus") ***REMOVED***``
		return FuncDef***REMOVED***elemName, elemMethod, "attrIsPresent", attrName, returnType, returnType, nil***REMOVED***

	default:
		panic("Unknown return type in a funcDef: " + returnType)
	***REMOVED***
***REMOVED***

// Node handler functions used in ast.Inspect to scrape TagName consts and the names of Element structs and their parent/nested struct

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

// Found a const, a Tagname suffix means it's for an Element.
func (ce *ElemInfoCollectorState) tagNameValueSpecHandler(node ast.Node) ElemInfoCollector ***REMOVED***
	switch x := node.(type) ***REMOVED***
	case *ast.Ident:
		if strings.HasSuffix(x.Name, "TagName") ***REMOVED***
			ce.elemName = strings.TrimSuffix(x.Name, "TagName")
			ce.elemInfos[ce.elemName] = &ElemInfo***REMOVED***x.Name, "", ""***REMOVED***
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
			// Ignore elements which haven't had an elemInfo created by the TagName handler.
			// (Specifically intermediate structs like HrefElement, MediaElement & FormFieldElement)
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
