package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type ElemInfo struct ***REMOVED***
	ConstName     string
	TagName       string
	StructName    string
	PrtStructName string
***REMOVED***

type NodeHandler func(node ast.Node) NodeHandler

type CollectElements struct ***REMOVED***
	handler   NodeHandler
	elemName  string
	elemInfos map[string]*ElemInfo
***REMOVED***

type FuncDef struct ***REMOVED***
	ElemName   string
	ElemMethod string
	AttrMethod string
	AttrName   string
	ReturnType string
	ReturnOpts []string
***REMOVED***

var renameTestElems = map[string]string***REMOVED***
	"href":            "a",
	"mod":             "del",
	"tablecell":       "",
	"tableheadercell": "",
***REMOVED***

var funcDefs = []string***REMOVED***
	"Href Rel string",
	"Href Href string",
	"Href Target string",
	"Href Type string",
	"Href AccessKey string",
	"Href HrefLang string",
	"Href Media string",
	"Href ToString=href string",

	"Base Href bool",
	"Base Target bool",

	"Button AccessKey string",
	"Button Autofocus bool",
	"Button Disabled bool",
	"Button Type enum=submit,button,menu,reset",

	"Data Value string",

	"Embed Height string",
	"Embed Width string",
	"Embed Src string",
	"Embed Type string",

	"FieldSet Disabled bool",
	"FieldSet Name string",

	"Form Name string",
	"Form Target string",
	"Form Action string",
	"Form Enctype string",
	"Form Encoding=enctype string",
	"Form AcceptCharset=accept-charset string",
	"Form Autocomplete string",
	"Form NoValidate bool",

	"IFrame Height string",
	"IFrame Width string",
	"IFrame Name string",
	"IFrame Src string",

	"Image Alt string",
	"Image CrossOrigin enum=anonymous,use-credentials",
	"Image Height int",
	"Image Width int",
	"Image IsMap bool",
	"Image Name string",
	"Image Src string",
	"Image UseMap string",

	"Input Type enum=text,button,checkbox,color,date,datetime-local,email,file,hidden,image,month,number,password,radio,range,reset,search,submit,tel,time,url,week",
	"Input Disabled bool",
	"Input Autofocus bool",
	"Input Required bool",
	"Input Value string",

	"Input Checked bool",
	"Input DefaultChecked=checked bool",

	"Input Alt string",
	"Input Src string",
	"Input Height string",
	"Input Width string",

	"Input Accept string",
	"Input Autocomplete enum=off,on",
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
	"Keygen Type enum=keygen",

	"Label HtmlFor=for string",

	"Legend AccessKey string",
	"Legend Value string",

	"Li Value string",
	"Li Type enum=,1,a,A,i,I,disc,square,circle",

	"Link Href string",
	"Link Hreflang string",
	"Link Media string",
	// The first value in enum lists gets used as the default. Putting "," at the start makes "" the default value for Rel instead of "alternate"
	"Link Rel enum=,alternate,author,dns-prefetch,help,icon,license,next,pingback,preconnect,prefetch,preload,prerender,prev,search,stylesheet",
	"Link Target string",
	"Link Type string",

	"Map Name string",

	"Meta Content string",
	"Meta HttpEquiv=http-equiv enum=content-type,default-style,refresh",
	"Meta Name enum=application-name,author,description,generator,keywords,viewport",

	"Meter Min int",
	"Meter Max int",
	"Meter High int",
	"Meter Low int",
	"Meter Optimum int",

	"Mod Cite string",
	"Mod DateTime string",

	"Object Data string",
	"Object Height string",
	"Object Name string",
	"Object Type string",
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
	"Output Type enum=output",

	"Param Name string",
	"Param Value string",

	"Pre Name string",
	"Pre Value string",

	"Quote Cite string",

	"Script Type string",
	"Script Src string",
	"Script HtmlFor=for string",
	"Script Charset string",
	"Script Async bool",
	"Script Defer bool",
	"Script NoModule bool",

	"Select Autofocus bool",
	"Select Disabled bool",
	"Select Multiple bool",
	"Select Name string",
	"Select Required bool",

	"Source KeySystem string",
	"Source Media string",
	"Source Sizes string",
	"Source Src string",
	"Source SrcSet string",
	"Source Type string",

	"Style Media string",
	"Style Type string",
	"Style Disabled bool",

	"TableCell ColSpan int=1",
	"TableCell RowSpan int=1",
	"TableCell Headers string",

	"TableHeaderCell Abbr string",
	"TableHeaderCell Scope string",
	"TableHeaderCell Sorted bool",

	"TextArea Type enum=textarea",
	"TextArea Value string",
	"TextArea DefaultValue=value string",
	"TextArea Placeholder string",
	"TextArea Rows int",
	"TextArea Cols int",
	"TextArea MaxLength int",
	"TextArea AccessKey string",
	"TextArea ReadOnly bool",
	"TextArea Required bool",
	"TextArea Autocomplete bool",
	"TextArea Autocapitalize enum=none,off,characters,words,sentences",
	"TextArea Wrap string",

	"Time DateTime string",

	"Title Text string",

	"UList Type string",
***REMOVED***

type TestDef struct ***REMOVED***
	ElemHtmlName string
	ElemMethod   string
	AttrName     string
	AttrVal      string
	ReturnType   string
	ReturnOpts   []string
***REMOVED***

var collector = &CollectElements***REMOVED******REMOVED***

func main() ***REMOVED***
	fs := token.NewFileSet()
	parsedFile, parseErr := parser.ParseFile(fs, "elements.go", nil, 0)
	if parseErr != nil ***REMOVED***
		log.Fatalf("warning: internal error: could not parse elements.go: %s", parseErr)
		return
	***REMOVED***

	collector.handler = collector.defaultHandler
	collector.elemInfos = make(map[string]*ElemInfo)

	ast.Inspect(parsedFile, func(n ast.Node) bool ***REMOVED***
		if n != nil ***REMOVED***
			collector.handler = collector.handler(n)
		***REMOVED***
		return true
	***REMOVED***)

	f, err := os.Create("elements_gen.go")
	if err != nil ***REMOVED***
		log.Println("warning: internal error: invalid Go generated:", err)
	***REMOVED***

	elemFuncsTemplate.Execute(f, struct ***REMOVED***
		ElemInfos map[string]*ElemInfo
		FuncDefs  []string
	***REMOVED******REMOVED***
		collector.elemInfos,
		funcDefs,
	***REMOVED***)
	f.Close()

	f, err = os.Create("elements_gen_test.go")
	if err != nil ***REMOVED***
		log.Println("warning: internal error: invalid Go generated:", err)
	***REMOVED***

	testFuncTemplate.Execute(f, struct ***REMOVED***
		FuncDefs []string
	***REMOVED******REMOVED***
		funcDefs,
	***REMOVED***)
	f.Close()
***REMOVED***

var elemFuncsTemplate = template.Must(template.New("").Funcs(template.FuncMap***REMOVED***
	"buildStruct":  buildStruct,
	"buildFuncDef": collector.buildFuncDef,
***REMOVED***).Parse(`// go generate
// generated by js/modules/k6/html/gen/main.go directed by js/modules/k6/html/elements.go;  DO NOT EDIT
package html

import "github.com/dop251/goja"

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

***REMOVED******REMOVED*** range $funcDefStr := .FuncDefs ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** $funcDef := buildFuncDef $funcDefStr ***REMOVED******REMOVED***
func (e ***REMOVED******REMOVED***$funcDef.ElemName***REMOVED******REMOVED***) ***REMOVED******REMOVED***$funcDef.ElemMethod***REMOVED******REMOVED***() ***REMOVED******REMOVED*** if ne $funcDef.ReturnType "enum" ***REMOVED******REMOVED*** ***REMOVED******REMOVED***$funcDef.ReturnType***REMOVED******REMOVED******REMOVED******REMOVED***else***REMOVED******REMOVED*** string ***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED***
***REMOVED******REMOVED*** if ne $funcDef.ReturnType "enum" ***REMOVED******REMOVED*** return e.***REMOVED******REMOVED*** $funcDef.AttrMethod ***REMOVED******REMOVED***("***REMOVED******REMOVED*** $funcDef.AttrName ***REMOVED******REMOVED***"***REMOVED******REMOVED*** if $funcDef.ReturnOpts ***REMOVED******REMOVED***, ***REMOVED******REMOVED*** index $funcDef.ReturnOpts 0 ***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***)***REMOVED******REMOVED*** else ***REMOVED******REMOVED*** attrVal := e.attrAsString("***REMOVED******REMOVED*** $funcDef.AttrName ***REMOVED******REMOVED***")
	switch attrVal ***REMOVED*** ***REMOVED******REMOVED*** range $optVal := $funcDef.ReturnOpts ***REMOVED******REMOVED***
		case "***REMOVED******REMOVED***$optVal***REMOVED******REMOVED***": 
			return "***REMOVED******REMOVED***$optVal***REMOVED******REMOVED***"
		***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
	***REMOVED***
	return "***REMOVED******REMOVED*** index $funcDef.ReturnOpts 0***REMOVED******REMOVED***" ***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
***REMOVED***
***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
`))

var testFuncTemplate = template.Must(template.New("").Funcs(template.FuncMap***REMOVED***
	"buildTestDef": collector.buildTestDef,
***REMOVED***).Parse(`// go generate
// generated by js/modules/k6/html/gen/main.go directed by js/modules/k6/html/elements.go;  DO NOT EDIT
package html

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/stretchr/testify/assert"
)

const testGenElems = ` + "`" + `<html><body>
***REMOVED******REMOVED***- range $index, $testDefStr := .FuncDefs -***REMOVED******REMOVED***
	***REMOVED******REMOVED*** $def := buildTestDef $index $testDefStr -***REMOVED******REMOVED***
	***REMOVED******REMOVED*** if eq $def.ElemHtmlName "" -***REMOVED******REMOVED***
	***REMOVED******REMOVED*** else if eq $def.ReturnType "enum" -***REMOVED******REMOVED***
	***REMOVED******REMOVED*** range $optIdx, $optVal := $def.ReturnOpts -***REMOVED******REMOVED***
		<***REMOVED******REMOVED***$def.ElemHtmlName***REMOVED******REMOVED*** id="elem_***REMOVED******REMOVED***$index***REMOVED******REMOVED***_***REMOVED******REMOVED***$optIdx***REMOVED******REMOVED***" ***REMOVED******REMOVED***$def.AttrName***REMOVED******REMOVED***="***REMOVED******REMOVED***$optVal***REMOVED******REMOVED***"> ***REMOVED******REMOVED***end***REMOVED******REMOVED***
 	***REMOVED******REMOVED***- else if eq $def.ReturnType "bool" -***REMOVED******REMOVED***
	  <***REMOVED******REMOVED***$def.ElemHtmlName***REMOVED******REMOVED*** id="elem_***REMOVED******REMOVED***$index***REMOVED******REMOVED***" ***REMOVED******REMOVED***$def.AttrName***REMOVED******REMOVED***></***REMOVED******REMOVED*** $def.ElemHtmlName ***REMOVED******REMOVED***>
	***REMOVED******REMOVED***else -***REMOVED******REMOVED*** 
	  <***REMOVED******REMOVED***$def.ElemHtmlName***REMOVED******REMOVED*** id="elem_***REMOVED******REMOVED***$index***REMOVED******REMOVED***" ***REMOVED******REMOVED***$def.AttrName***REMOVED******REMOVED***="***REMOVED******REMOVED***$def.AttrVal***REMOVED******REMOVED***"></***REMOVED******REMOVED*** $def.ElemHtmlName ***REMOVED******REMOVED***>
	***REMOVED******REMOVED***end -***REMOVED******REMOVED***
***REMOVED******REMOVED***- end***REMOVED******REMOVED***
</body></html>` + "`" + `

func TestGenElements(t *testing.T) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	ctx := common.WithRuntime(context.Background(), rt)
	rt.Set("src", testGenElems)
	rt.Set("html", common.Bind(rt, &HTML***REMOVED******REMOVED***, &ctx))
	// compileProtoElem()

	_, err := common.RunString(rt, "let doc = html.parseHTML(src)")

	assert.NoError(t, err)
	assert.IsType(t, Selection***REMOVED******REMOVED***, rt.Get("doc").Export())
***REMOVED******REMOVED*** range $index, $testDefStr := .FuncDefs ***REMOVED******REMOVED*** 
***REMOVED******REMOVED*** $def := buildTestDef $index $testDefStr ***REMOVED******REMOVED*** 
	***REMOVED******REMOVED*** if ne $def.ElemHtmlName "" -***REMOVED******REMOVED***
		t.Run("***REMOVED******REMOVED***$def.ElemHtmlName***REMOVED******REMOVED***.***REMOVED******REMOVED***$def.ElemMethod***REMOVED******REMOVED***", func(t *testing.T) ***REMOVED*** 
	***REMOVED******REMOVED*** if eq $def.ReturnType "enum" -***REMOVED******REMOVED*** 
		***REMOVED******REMOVED*** range $optIdx, $optVal := $def.ReturnOpts -***REMOVED******REMOVED***
			if v, err := common.RunString(rt, "doc.find(\"#elem_***REMOVED******REMOVED***$index***REMOVED******REMOVED***_***REMOVED******REMOVED***$optIdx***REMOVED******REMOVED***\").get(0).***REMOVED******REMOVED***$def.ElemMethod***REMOVED******REMOVED***()"); assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, "***REMOVED******REMOVED***$optVal***REMOVED******REMOVED***", v.Export()) 
			***REMOVED*** 
		***REMOVED******REMOVED***end -***REMOVED******REMOVED***
	***REMOVED******REMOVED***else -***REMOVED******REMOVED*** 
	  	if v, err := common.RunString(rt, "doc.find(\"#elem_***REMOVED******REMOVED***$index***REMOVED******REMOVED***\").get(0).***REMOVED******REMOVED***$def.ElemMethod***REMOVED******REMOVED***()"); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, ***REMOVED******REMOVED*** if eq $def.ReturnType "bool" ***REMOVED******REMOVED******REMOVED******REMOVED***$def.AttrVal***REMOVED******REMOVED*** ***REMOVED******REMOVED***else if eq $def.ReturnType "string" ***REMOVED******REMOVED*** "***REMOVED******REMOVED***$def.AttrVal***REMOVED******REMOVED***" ***REMOVED******REMOVED***else if eq $def.ReturnType "int"***REMOVED******REMOVED*** int64(***REMOVED******REMOVED***$def.AttrVal***REMOVED******REMOVED***) ***REMOVED******REMOVED***end***REMOVED******REMOVED***, v.Export()) 
			***REMOVED*** 
	***REMOVED******REMOVED***end -***REMOVED******REMOVED***
***REMOVED***)
***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
***REMOVED******REMOVED***- end -***REMOVED******REMOVED***
***REMOVED***
`))

func buildStruct(elemInfo ElemInfo) string ***REMOVED***
	if elemInfo.PrtStructName == "Element" ***REMOVED***
		return elemInfo.StructName + "***REMOVED***elem***REMOVED***"
	***REMOVED*** else ***REMOVED***
		return elemInfo.StructName + "***REMOVED***" + elemInfo.PrtStructName + "***REMOVED***elem***REMOVED******REMOVED***"
	***REMOVED***
***REMOVED***

func (ce *CollectElements) buildFuncDef(funcDef string) FuncDef ***REMOVED***
	parts := strings.Split(funcDef, " ")
	// parts[0] is the element struct name (without the Element suffix for brevity)
	// parts[1] is either:
	//   MethodName               The name of method added onto that struct and converted to lowercase thenn used as the argument to elem.attrAsString(...) or elem.AttrIsPresent(...)
	//   MethodName=attrname      The MethodName is added to the struct. The attrname is the argument for attrAsString or AttrIsPresent
	// parts[2] is the return type, either string or bool
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
		// "Button AccessKey string" => ***REMOVED***"ButtonElement" "AccessKey", "attrIsString", "accesskey", "string"***REMOVED*** => `func (e ButtonElement) AccessKey() string***REMOVED*** return e.attrAsString("accessKey") ***REMOVED***``
		if returnOpts == "" ***REMOVED***
			return FuncDef***REMOVED***elemName, elemMethod, "attrAsInt", attrName, returnType, []string***REMOVED***"0"***REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			return FuncDef***REMOVED***elemName, elemMethod, "attrAsInt", attrName, returnType, []string***REMOVED***returnOpts***REMOVED******REMOVED***
		***REMOVED***

	case "enum":
		return FuncDef***REMOVED***elemName, elemMethod, "attrAsInt", attrName, returnType, strings.Split(returnOpts, ",")***REMOVED***

	case "string":
		// "Button AccessKey string" => ***REMOVED***"ButtonElement" "AccessKey", "attrIsString", "accesskey", "string"***REMOVED*** => `func (e ButtonElement) AccessKey() string***REMOVED*** return e.attrAsString("accessKey") ***REMOVED***``
		return FuncDef***REMOVED***elemName, elemMethod, "attrAsString", attrName, returnType, nil***REMOVED***

	case "bool":
		// "Button Autofocus bool" ***REMOVED***"Button" "Autofocus", "attrIsPresent", "autofocus", "bool"***REMOVED*** => `func (e ButtonElement) ToString() bool ***REMOVED*** return e.attrIsPresent("autofocus") ***REMOVED***``
		return FuncDef***REMOVED***elemName, elemMethod, "attrIsPresent", attrName, returnType, nil***REMOVED***
	default:
		panic("Unknown return type in a funcDef: " + returnType)
	***REMOVED***
***REMOVED***

func (ce *CollectElements) buildTestDef(index int, testDef string) TestDef ***REMOVED***
	parts := strings.Split(testDef, " ")

	elemHtmlName := strings.ToLower(parts[0])
	elemMethod := strings.ToLower(parts[1][0:1]) + parts[1][1:]
	attrName := strings.ToLower(parts[1])
	returnType := parts[2]
	returnOpts := ""

	if useElemName, ok := renameTestElems[elemHtmlName]; ok ***REMOVED***
		elemHtmlName = useElemName
	***REMOVED*** else if elemInfo, ok := ce.elemInfos[parts[0]]; ok ***REMOVED***
		elemHtmlName = strings.Trim(elemInfo.TagName, "\"")
	***REMOVED***

	if eqPos := strings.Index(elemMethod, "="); eqPos != -1 ***REMOVED***
		attrName = elemMethod[eqPos+1:]
		elemMethod = elemMethod[0:eqPos]
	***REMOVED***

	if eqPos := strings.Index(returnType, "="); eqPos != -1 ***REMOVED***
		returnOpts = returnType[eqPos+1:]
		returnType = returnType[0:eqPos]
	***REMOVED***

	switch returnType ***REMOVED***
	case "bool":
		return TestDef***REMOVED***elemHtmlName, elemMethod, attrName, "true", returnType, nil***REMOVED***
	case "string":
		return TestDef***REMOVED***elemHtmlName, elemMethod, attrName, "attrval_" + strconv.Itoa(index), returnType, nil***REMOVED***
	case "int":
		return TestDef***REMOVED***elemHtmlName, elemMethod, attrName, strconv.Itoa(index), returnType, nil***REMOVED***
	case "enum":
		return TestDef***REMOVED***elemHtmlName, elemMethod, attrName, "attrval_" + strconv.Itoa(index), returnType, strings.Split(returnOpts, ",")***REMOVED***
	default:
		panic("Unknown return type in a funcDef:" + returnType)

	***REMOVED***
***REMOVED***

// Node handler functions used in ast.Inspect to scrape TagName consts and the names of Element structs and their parent/nested struct

func (ce *CollectElements) defaultHandler(node ast.Node) NodeHandler ***REMOVED***
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

func (ce *CollectElements) tagNameValueSpecHandler(node ast.Node) NodeHandler ***REMOVED***
	switch x := node.(type) ***REMOVED***
	case *ast.Ident:
		if strings.HasSuffix(x.Name, "TagName") ***REMOVED***
			ce.elemName = strings.TrimSuffix(x.Name, "TagName")
			ce.elemInfos[ce.elemName] = &ElemInfo***REMOVED***x.Name, "", "", ""***REMOVED***
			return ce.tagNameValueSpecHandler
		***REMOVED***

		return ce.defaultHandler
	case *ast.BasicLit:
		if _, ok := ce.elemInfos[ce.elemName]; !ok ***REMOVED***
			return ce.defaultHandler
		***REMOVED***

		ce.elemInfos[ce.elemName].TagName = x.Value
		return ce.defaultHandler

	default:
		return ce.defaultHandler
	***REMOVED***
***REMOVED***

func (ce *CollectElements) elemTypeSpecHandler(node ast.Node) NodeHandler ***REMOVED***
	switch x := node.(type) ***REMOVED***
	case *ast.Ident:
		if !strings.HasSuffix(x.Name, "Element") ***REMOVED***
			return ce.defaultHandler
		***REMOVED***

		if ce.elemName == "" ***REMOVED***
			ce.elemName = strings.TrimSuffix(x.Name, "Element")
			// Ignore HrefElement and MediaElement structs. They are subclassed by AnchorElement/AreaElement/VideoElement and do not have their own entry in ElemInfos
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
