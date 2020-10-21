%***REMOVED***
package protoparse

//lint:file-ignore SA4006 generated parser has unused values

import (
	"fmt"
	"math"
)

%***REMOVED***

// fields inside this union end up as the fields in a structure known
// as $***REMOVED***PREFIX***REMOVED***SymType, of which a reference is passed to the lexer.
%union***REMOVED***
	file      *fileNode
	fileDecls []*fileElement
	syn       *syntaxNode
	pkg       *packageNode
	imprt     *importNode
	msg       *messageNode
	msgDecls  []*messageElement
	fld       *fieldNode
	mapFld    *mapFieldNode
	mapType   *mapTypeNode
	grp       *groupNode
	oo        *oneOfNode
	ooDecls   []*oneOfElement
	ext       *extensionRangeNode
	resvd     *reservedNode
	en        *enumNode
	enDecls   []*enumElement
	env       *enumValueNode
	extend    *extendNode
	extDecls  []*extendElement
	svc       *serviceNode
	svcDecls  []*serviceElement
	mtd       *methodNode
	rpcType   *rpcTypeNode
	opts      []*optionNode
	optNm     []*optionNamePartNode
	cmpctOpts *compactOptionsNode
	rngs      []*rangeNode
	names     []*compoundStringNode
	cid       *compoundIdentNode
	sl        []valueNode
	agg       []*aggregateEntryNode
	aggName   *aggregateNameNode
	v         valueNode
	il        *compoundIntNode
	str       *compoundStringNode
	s         *stringLiteralNode
	i         *intLiteralNode
	f         *floatLiteralNode
	id        *identNode
	b         *basicNode
	err       error
***REMOVED***

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <file>      file
%type <syn>       syntax
%type <fileDecls> fileDecl fileDecls
%type <imprt>     import
%type <pkg>       package
%type <opts>      option compactOption compactOptionDecls rpcOption rpcOptions
%type <optNm>     optionName optionNameComponent
%type <cmpctOpts> compactOptions
%type <v>         constant scalarConstant aggregate numLit
%type <il>        intLit
%type <id>        name keyType
%type <cid>       ident typeIdent
%type <aggName>   aggName
%type <sl>        constantList
%type <agg>       aggFields aggField aggFieldEntry
%type <fld>       field oneofField
%type <oo>        oneof
%type <grp>       group oneofGroup
%type <mapFld>    mapField
%type <mapType>   mapType
%type <msg>       message
%type <msgDecls>  messageItem messageBody
%type <ooDecls>   oneofItem oneofBody
%type <names>     fieldNames
%type <resvd>     msgReserved enumReserved reservedNames
%type <rngs>      tagRange tagRanges enumRange enumRanges
%type <ext>       extensions
%type <en>        enum
%type <enDecls>   enumItem enumBody
%type <env>       enumField
%type <extend>    extend
%type <extDecls>  extendItem extendBody
%type <str>       stringLit
%type <svc>       service
%type <svcDecls>  serviceItem serviceBody
%type <mtd>       rpc
%type <rpcType>   rpcType

// same for terminals
%token <s> _STRING_LIT
%token <i>  _INT_LIT
%token <f>   _FLOAT_LIT
%token <id>  _NAME
%token <id>  _SYNTAX _IMPORT _WEAK _PUBLIC _PACKAGE _OPTION _TRUE _FALSE _INF _NAN _REPEATED _OPTIONAL _REQUIRED
%token <id>  _DOUBLE _FLOAT _INT32 _INT64 _UINT32 _UINT64 _SINT32 _SINT64 _FIXED32 _FIXED64 _SFIXED32 _SFIXED64
%token <id>  _BOOL _STRING _BYTES _GROUP _ONEOF _MAP _EXTENSIONS _TO _MAX _RESERVED _ENUM _MESSAGE _EXTEND
%token <id>  _SERVICE _RPC _STREAM _RETURNS
%token <err> _ERROR
// we define all of these, even ones that aren't used, to improve error messages
// so it shows the unexpected symbol instead of showing "$unk"
%token <b>   '=' ';' ':' '***REMOVED***' '***REMOVED***' '\\' '/' '?' '.' ',' '>' '<' '+' '-' '(' ')' '[' ']' '*' '&' '^' '%' '$' '#' '@' '!' '~' '`'

%%

file : syntax ***REMOVED***
		$$ = &fileNode***REMOVED***syntax: $1***REMOVED***
		$$.setRange($1, $1)
		protolex.(*protoLex).res = $$
	***REMOVED***
	| fileDecls  ***REMOVED***
		$$ = &fileNode***REMOVED***decls: $1***REMOVED***
		if len($1) > 0 ***REMOVED***
			$$.setRange($1[0], $1[len($1)-1])
		***REMOVED***
		protolex.(*protoLex).res = $$
	***REMOVED***
	| syntax fileDecls ***REMOVED***
		$$ = &fileNode***REMOVED***syntax: $1, decls: $2***REMOVED***
		var end node
		if len($2) > 0 ***REMOVED***
			end = $2[len($2)-1]
		***REMOVED*** else ***REMOVED***
			end = $1
		***REMOVED***
		$$.setRange($1, end)
		protolex.(*protoLex).res = $$
	***REMOVED***
	| ***REMOVED***
	***REMOVED***

fileDecls : fileDecls fileDecl ***REMOVED***
		$$ = append($1, $2...)
	***REMOVED***
	| fileDecl

fileDecl : import ***REMOVED***
		$$ = []*fileElement***REMOVED******REMOVED***imp: $1***REMOVED******REMOVED***
	***REMOVED***
	| package ***REMOVED***
		$$ = []*fileElement***REMOVED******REMOVED***pkg: $1***REMOVED******REMOVED***
	***REMOVED***
	| option ***REMOVED***
		$$ = []*fileElement***REMOVED******REMOVED***option: $1[0]***REMOVED******REMOVED***
	***REMOVED***
	| message ***REMOVED***
		$$ = []*fileElement***REMOVED******REMOVED***message: $1***REMOVED******REMOVED***
	***REMOVED***
	| enum ***REMOVED***
		$$ = []*fileElement***REMOVED******REMOVED***enum: $1***REMOVED******REMOVED***
	***REMOVED***
	| extend ***REMOVED***
		$$ = []*fileElement***REMOVED******REMOVED***extend: $1***REMOVED******REMOVED***
	***REMOVED***
	| service ***REMOVED***
		$$ = []*fileElement***REMOVED******REMOVED***service: $1***REMOVED******REMOVED***
	***REMOVED***
	| ';' ***REMOVED***
		$$ = []*fileElement***REMOVED******REMOVED***empty: $1***REMOVED******REMOVED***
	***REMOVED***
	| error ';' ***REMOVED***
	***REMOVED***
	| error ***REMOVED***
	***REMOVED***

syntax : _SYNTAX '=' stringLit ';' ***REMOVED***
		$$ = &syntaxNode***REMOVED***syntax: $3***REMOVED***
		$$.setRange($1, $4)
	***REMOVED***

import : _IMPORT stringLit ';' ***REMOVED***
		$$ = &importNode***REMOVED*** name: $2 ***REMOVED***
		$$.setRange($1, $3)
	***REMOVED***
	| _IMPORT _WEAK stringLit ';' ***REMOVED***
		$$ = &importNode***REMOVED*** name: $3, weak: true ***REMOVED***
		$$.setRange($1, $4)
	***REMOVED***
	| _IMPORT _PUBLIC stringLit ';' ***REMOVED***
		$$ = &importNode***REMOVED*** name: $3, public: true ***REMOVED***
		$$.setRange($1, $4)
	***REMOVED***

package : _PACKAGE ident ';' ***REMOVED***
		$$ = &packageNode***REMOVED***name: $2***REMOVED***
		$$.setRange($1, $3)
	***REMOVED***

ident : name ***REMOVED***
        $$ = &compoundIdentNode***REMOVED***val: $1.val***REMOVED***
        $$.setRange($1, $1)
    ***REMOVED***
	| ident '.' name ***REMOVED***
        $$ = &compoundIdentNode***REMOVED***val: $1.val + "." + $3.val***REMOVED***
        $$.setRange($1, $3)
	***REMOVED***

option : _OPTION optionName '=' constant ';' ***REMOVED***
		n := &optionNameNode***REMOVED***parts: $2***REMOVED***
		n.setRange($2[0], $2[len($2)-1])
		o := &optionNode***REMOVED***name: n, val: $4***REMOVED***
		o.setRange($1, $5)
		$$ = []*optionNode***REMOVED***o***REMOVED***
	***REMOVED***

optionName : optionNameComponent
    |
    optionName '.' optionNameComponent ***REMOVED***
		$$ = append($1, $3...)
	***REMOVED***


optionNameComponent : name ***REMOVED***
        nm := &compoundIdentNode***REMOVED***val: $1.val***REMOVED***
        nm.setRange($1, $1)
		$$ = toNameParts(nm)
	***REMOVED***
	| '(' typeIdent ')' ***REMOVED***
		p := &optionNamePartNode***REMOVED***text: $2, isExtension: true***REMOVED***
		p.setRange($1, $3)
		$$ = []*optionNamePartNode***REMOVED***p***REMOVED***
	***REMOVED***

constant : scalarConstant
	| aggregate

scalarConstant : stringLit ***REMOVED***
		$$ = $1
	***REMOVED***
	| numLit
	| name ***REMOVED***
		if $1.val == "true" ***REMOVED***
			$$ = &boolLiteralNode***REMOVED***identNode: $1, val: true***REMOVED***
		***REMOVED*** else if $1.val == "false" ***REMOVED***
			$$ = &boolLiteralNode***REMOVED***identNode: $1, val: false***REMOVED***
		***REMOVED*** else if $1.val == "inf" ***REMOVED***
			f := &compoundFloatNode***REMOVED***val: math.Inf(1)***REMOVED***
			f.setRange($1, $1)
			$$ = f
		***REMOVED*** else if $1.val == "nan" ***REMOVED***
			f := &compoundFloatNode***REMOVED***val: math.NaN()***REMOVED***
			f.setRange($1, $1)
			$$ = f
		***REMOVED*** else ***REMOVED***
			$$ = $1
		***REMOVED***
	***REMOVED***

numLit : _FLOAT_LIT ***REMOVED***
        $$ = $1
    ***REMOVED***
	| '-' _FLOAT_LIT ***REMOVED***
		f := &compoundFloatNode***REMOVED***val: -$2.val***REMOVED***
		f.setRange($1, $2)
		$$ = f
	***REMOVED***
	| '+' _FLOAT_LIT ***REMOVED***
		f := &compoundFloatNode***REMOVED***val: $2.val***REMOVED***
		f.setRange($1, $2)
		$$ = f
	***REMOVED***
	| '+' _INF ***REMOVED***
		f := &compoundFloatNode***REMOVED***val: math.Inf(1)***REMOVED***
		f.setRange($1, $2)
		$$ = f
	***REMOVED***
	| '-' _INF ***REMOVED***
		f := &compoundFloatNode***REMOVED***val: math.Inf(-1)***REMOVED***
		f.setRange($1, $2)
		$$ = f
	***REMOVED***
	| _INT_LIT ***REMOVED***
        $$ = $1
    ***REMOVED***
    | '+' _INT_LIT ***REMOVED***
          i := &compoundUintNode***REMOVED***val: $2.val***REMOVED***
          i.setRange($1, $2)
          $$ = i
    ***REMOVED***
    | '-' _INT_LIT ***REMOVED***
        if $2.val > math.MaxInt64 + 1 ***REMOVED***
            // can't represent as int so treat as float literal
            f := &compoundFloatNode***REMOVED***val: -float64($2.val)***REMOVED***
            f.setRange($1, $2)
            $$ = f
        ***REMOVED*** else ***REMOVED***
            i := &compoundIntNode***REMOVED***val: -int64($2.val)***REMOVED***
            i.setRange($1, $2)
            $$ = i
        ***REMOVED***
    ***REMOVED***

stringLit : _STRING_LIT ***REMOVED***
        $$ = &compoundStringNode***REMOVED***val: $1.val***REMOVED***
        $$.setRange($1, $1)
    ***REMOVED***
    | stringLit _STRING_LIT ***REMOVED***
        $$ = &compoundStringNode***REMOVED***val: $1.val + $2.val***REMOVED***
        $$.setRange($1, $2)
    ***REMOVED***

aggregate : '***REMOVED***' aggFields '***REMOVED***' ***REMOVED***
		a := &aggregateLiteralNode***REMOVED***elements: $2***REMOVED***
		a.setRange($1, $3)
		$$ = a
	***REMOVED***

aggFields : aggField
	| aggFields aggField ***REMOVED***
		$$ = append($1, $2...)
	***REMOVED***
	| ***REMOVED***
		$$ = nil
	***REMOVED***

aggField : aggFieldEntry
	| aggFieldEntry ',' ***REMOVED***
		$$ = $1
	***REMOVED***
	| aggFieldEntry ';' ***REMOVED***
		$$ = $1
	***REMOVED***
	| error ',' ***REMOVED***
	***REMOVED***
	| error ';' ***REMOVED***
	***REMOVED***
	| error ***REMOVED***
	***REMOVED***

aggFieldEntry : aggName ':' scalarConstant ***REMOVED***
		a := &aggregateEntryNode***REMOVED***name: $1, val: $3***REMOVED***
		a.setRange($1, $3)
		$$ = []*aggregateEntryNode***REMOVED***a***REMOVED***
	***REMOVED***
	| aggName ':' '[' ']' ***REMOVED***
		s := &sliceLiteralNode***REMOVED******REMOVED***
		s.setRange($3, $4)
		a := &aggregateEntryNode***REMOVED***name: $1, val: s***REMOVED***
		a.setRange($1, $4)
		$$ = []*aggregateEntryNode***REMOVED***a***REMOVED***
	***REMOVED***
	| aggName ':' '[' constantList ']' ***REMOVED***
		s := &sliceLiteralNode***REMOVED***elements: $4***REMOVED***
		s.setRange($3, $5)
		a := &aggregateEntryNode***REMOVED***name: $1, val: s***REMOVED***
		a.setRange($1, $5)
		$$ = []*aggregateEntryNode***REMOVED***a***REMOVED***
	***REMOVED***
	| aggName ':' '[' error ']' ***REMOVED***
	***REMOVED***
	| aggName ':' aggregate ***REMOVED***
		a := &aggregateEntryNode***REMOVED***name: $1, val: $3***REMOVED***
		a.setRange($1, $3)
		$$ = []*aggregateEntryNode***REMOVED***a***REMOVED***
	***REMOVED***
	| aggName aggregate ***REMOVED***
		a := &aggregateEntryNode***REMOVED***name: $1, val: $2***REMOVED***
		a.setRange($1, $2)
		$$ = []*aggregateEntryNode***REMOVED***a***REMOVED***
	***REMOVED***
	| aggName ':' '<' aggFields '>' ***REMOVED***
		s := &aggregateLiteralNode***REMOVED***elements: $4***REMOVED***
		s.setRange($3, $5)
		a := &aggregateEntryNode***REMOVED***name: $1, val: s***REMOVED***
		a.setRange($1, $5)
		$$ = []*aggregateEntryNode***REMOVED***a***REMOVED***
	***REMOVED***
	| aggName '<' aggFields '>' ***REMOVED***
		s := &aggregateLiteralNode***REMOVED***elements: $3***REMOVED***
		s.setRange($2, $4)
		a := &aggregateEntryNode***REMOVED***name: $1, val: s***REMOVED***
		a.setRange($1, $4)
		$$ = []*aggregateEntryNode***REMOVED***a***REMOVED***
	***REMOVED***
	| aggName ':' '<' error '>' ***REMOVED***
	***REMOVED***
	| aggName '<' error '>' ***REMOVED***
	***REMOVED***

aggName : name ***REMOVED***
        n := &compoundIdentNode***REMOVED***val: $1.val***REMOVED***
        n.setRange($1, $1)
		$$ = &aggregateNameNode***REMOVED***name: n***REMOVED***
		$$.setRange($1, $1)
	***REMOVED***
	| '[' typeIdent ']' ***REMOVED***
		$$ = &aggregateNameNode***REMOVED***name: $2, isExtension: true***REMOVED***
		$$.setRange($1, $3)
	***REMOVED***
	| '[' error ']' ***REMOVED***
	***REMOVED***

constantList : constant ***REMOVED***
		$$ = []valueNode***REMOVED***$1***REMOVED***
	***REMOVED***
	| constantList ',' constant ***REMOVED***
		$$ = append($1, $3)
	***REMOVED***
	| constantList ';' constant ***REMOVED***
		$$ = append($1, $3)
	***REMOVED***
	| '<' aggFields '>' ***REMOVED***
		s := &aggregateLiteralNode***REMOVED***elements: $2***REMOVED***
		s.setRange($1, $3)
		$$ = []valueNode***REMOVED***s***REMOVED***
	***REMOVED***
	| constantList ','  '<' aggFields '>' ***REMOVED***
		s := &aggregateLiteralNode***REMOVED***elements: $4***REMOVED***
		s.setRange($3, $5)
		$$ = append($1, s)
	***REMOVED***
	| constantList ';'  '<' aggFields '>' ***REMOVED***
		s := &aggregateLiteralNode***REMOVED***elements: $4***REMOVED***
		s.setRange($3, $5)
		$$ = append($1, s)
	***REMOVED***
	| '<' error '>' ***REMOVED***
	***REMOVED***
	| constantList ','  '<' error '>' ***REMOVED***
	***REMOVED***
	| constantList ';'  '<' error '>' ***REMOVED***
	***REMOVED***

typeIdent : ident
    | '.' ident ***REMOVED***
          $$ = &compoundIdentNode***REMOVED***val: "." + $2.val***REMOVED***
          $$.setRange($1, $2)
    ***REMOVED***

field : _REQUIRED typeIdent name '=' _INT_LIT ';' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1, required: true***REMOVED***
		$$ = &fieldNode***REMOVED***label: lbl, fldType: $2, name: $3, tag: $5***REMOVED***
		$$.setRange($1, $6)
	***REMOVED***
	| _OPTIONAL typeIdent name '=' _INT_LIT ';' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1***REMOVED***
		$$ = &fieldNode***REMOVED***label: lbl, fldType: $2, name: $3, tag: $5***REMOVED***
		$$.setRange($1, $6)
	***REMOVED***
	| _REPEATED typeIdent name '=' _INT_LIT ';' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1, repeated: true***REMOVED***
		$$ = &fieldNode***REMOVED***label: lbl, fldType: $2, name: $3, tag: $5***REMOVED***
		$$.setRange($1, $6)
	***REMOVED***
	| typeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = &fieldNode***REMOVED***fldType: $1, name: $2, tag: $4***REMOVED***
		$$.setRange($1, $5)
	***REMOVED***
	| _REQUIRED typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1, required: true***REMOVED***
		$$ = &fieldNode***REMOVED***label: lbl, fldType: $2, name: $3, tag: $5, options: $6***REMOVED***
		$$.setRange($1, $7)
	***REMOVED***
	| _OPTIONAL typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1***REMOVED***
		$$ = &fieldNode***REMOVED***label: lbl, fldType: $2, name: $3, tag: $5, options: $6***REMOVED***
		$$.setRange($1, $7)
	***REMOVED***
	| _REPEATED typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1, repeated: true***REMOVED***
		$$ = &fieldNode***REMOVED***label: lbl, fldType: $2, name: $3, tag: $5, options: $6***REMOVED***
		$$.setRange($1, $7)
	***REMOVED***
	| typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = &fieldNode***REMOVED***fldType: $1, name: $2, tag: $4, options: $5***REMOVED***
		$$.setRange($1, $6)
	***REMOVED***

compactOptions: '[' compactOptionDecls ']' ***REMOVED***
        $$ = &compactOptionsNode***REMOVED***decls: $2***REMOVED***
        $$.setRange($1, $3)
    ***REMOVED***

compactOptionDecls : compactOptionDecls ',' compactOption ***REMOVED***
		$$ = append($1, $3...)
	***REMOVED***
	| compactOption

compactOption: optionName '=' constant ***REMOVED***
		n := &optionNameNode***REMOVED***parts: $1***REMOVED***
		n.setRange($1[0], $1[len($1)-1])
		o := &optionNode***REMOVED***name: n, val: $3***REMOVED***
		o.setRange($1[0], $3)
		$$ = []*optionNode***REMOVED***o***REMOVED***
	***REMOVED***

group : _REQUIRED _GROUP name '=' _INT_LIT '***REMOVED***' messageBody '***REMOVED***' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1, required: true***REMOVED***
		$$ = &groupNode***REMOVED***groupKeyword: $2, label: lbl, name: $3, tag: $5, decls: $7***REMOVED***
		$$.setRange($1, $8)
	***REMOVED***
	| _OPTIONAL _GROUP name '=' _INT_LIT '***REMOVED***' messageBody '***REMOVED***' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1***REMOVED***
		$$ = &groupNode***REMOVED***groupKeyword: $2, label: lbl, name: $3, tag: $5, decls: $7***REMOVED***
		$$.setRange($1, $8)
	***REMOVED***
	| _REPEATED _GROUP name '=' _INT_LIT '***REMOVED***' messageBody '***REMOVED***' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1, repeated: true***REMOVED***
		$$ = &groupNode***REMOVED***groupKeyword: $2, label: lbl, name: $3, tag: $5, decls: $7***REMOVED***
		$$.setRange($1, $8)
	***REMOVED***
	| _REQUIRED _GROUP name '=' _INT_LIT compactOptions '***REMOVED***' messageBody '***REMOVED***' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1, required: true***REMOVED***
		$$ = &groupNode***REMOVED***groupKeyword: $2, label: lbl, name: $3, tag: $5, options: $6, decls: $8***REMOVED***
		$$.setRange($1, $9)
	***REMOVED***
	| _OPTIONAL _GROUP name '=' _INT_LIT compactOptions '***REMOVED***' messageBody '***REMOVED***' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1***REMOVED***
		$$ = &groupNode***REMOVED***groupKeyword: $2, label: lbl, name: $3, tag: $5, options: $6, decls: $8***REMOVED***
		$$.setRange($1, $9)
	***REMOVED***
	| _REPEATED _GROUP name '=' _INT_LIT compactOptions '***REMOVED***' messageBody '***REMOVED***' ***REMOVED***
		lbl := fieldLabel***REMOVED***identNode: $1, repeated: true***REMOVED***
		$$ = &groupNode***REMOVED***groupKeyword: $2, label: lbl, name: $3, tag: $5, options: $6, decls: $8***REMOVED***
		$$.setRange($1, $9)
	***REMOVED***

oneof : _ONEOF name '***REMOVED***' oneofBody '***REMOVED***' ***REMOVED***
		$$ = &oneOfNode***REMOVED***name: $2, decls: $4***REMOVED***
		$$.setRange($1, $5)
	***REMOVED***

oneofBody : oneofBody oneofItem ***REMOVED***
		$$ = append($1, $2...)
	***REMOVED***
	| oneofItem
	| ***REMOVED***
		$$ = nil
	***REMOVED***

oneofItem : option ***REMOVED***
		$$ = []*oneOfElement***REMOVED******REMOVED***option: $1[0]***REMOVED******REMOVED***
	***REMOVED***
	| oneofField ***REMOVED***
		$$ = []*oneOfElement***REMOVED******REMOVED***field: $1***REMOVED******REMOVED***
	***REMOVED***
	| oneofGroup ***REMOVED***
		$$ = []*oneOfElement***REMOVED******REMOVED***group: $1***REMOVED******REMOVED***
	***REMOVED***
	| ';' ***REMOVED***
		$$ = []*oneOfElement***REMOVED******REMOVED***empty: $1***REMOVED******REMOVED***
	***REMOVED***
	| error ';' ***REMOVED***
	***REMOVED***
	| error ***REMOVED***
	***REMOVED***

oneofField : typeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = &fieldNode***REMOVED***fldType: $1, name: $2, tag: $4***REMOVED***
		$$.setRange($1, $5)
	***REMOVED***
	| typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = &fieldNode***REMOVED***fldType: $1, name: $2, tag: $4, options: $5***REMOVED***
		$$.setRange($1, $6)
	***REMOVED***

oneofGroup : _GROUP name '=' _INT_LIT '***REMOVED***' messageBody '***REMOVED***' ***REMOVED***
		$$ = &groupNode***REMOVED***groupKeyword: $1, name: $2, tag: $4, decls: $6***REMOVED***
		$$.setRange($1, $7)
	***REMOVED***
	| _GROUP name '=' _INT_LIT compactOptions '***REMOVED***' messageBody '***REMOVED***' ***REMOVED***
		$$ = &groupNode***REMOVED***groupKeyword: $1, name: $2, tag: $4, options: $5, decls: $7***REMOVED***
		$$.setRange($1, $8)
	***REMOVED***

mapField : mapType name '=' _INT_LIT ';' ***REMOVED***
		$$ = &mapFieldNode***REMOVED***mapType: $1, name: $2, tag: $4***REMOVED***
		$$.setRange($1, $5)
	***REMOVED***
	| mapType name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = &mapFieldNode***REMOVED***mapType: $1, name: $2, tag: $4, options: $5***REMOVED***
		$$.setRange($1, $6)
	***REMOVED***

mapType : _MAP '<' keyType ',' typeIdent '>' ***REMOVED***
        $$ = &mapTypeNode***REMOVED***mapKeyword: $1, keyType: $3, valueType: $5***REMOVED***
        $$.setRange($1, $6)
***REMOVED***

keyType : _INT32
	| _INT64
	| _UINT32
	| _UINT64
	| _SINT32
	| _SINT64
	| _FIXED32
	| _FIXED64
	| _SFIXED32
	| _SFIXED64
	| _BOOL
	| _STRING

extensions : _EXTENSIONS tagRanges ';' ***REMOVED***
		$$ = &extensionRangeNode***REMOVED***ranges: $2***REMOVED***
		$$.setRange($1, $3)
	***REMOVED***
	| _EXTENSIONS tagRanges compactOptions ';' ***REMOVED***
		$$ = &extensionRangeNode***REMOVED***ranges: $2, options: $3***REMOVED***
		$$.setRange($1, $4)
	***REMOVED***

tagRanges : tagRanges ',' tagRange ***REMOVED***
		$$ = append($1, $3...)
	***REMOVED***
	| tagRange

tagRange : _INT_LIT ***REMOVED***
		r := &rangeNode***REMOVED***startNode: $1***REMOVED***
		r.setRange($1, $1)
		$$ = []*rangeNode***REMOVED***r***REMOVED***
	***REMOVED***
	| _INT_LIT _TO _INT_LIT ***REMOVED***
		r := &rangeNode***REMOVED***startNode: $1, endNode: $3***REMOVED***
		r.setRange($1, $3)
		$$ = []*rangeNode***REMOVED***r***REMOVED***
	***REMOVED***
	| _INT_LIT _TO _MAX ***REMOVED***
		r := &rangeNode***REMOVED***startNode: $1, endNode: $3, endMax: true***REMOVED***
		r.setRange($1, $3)
		$$ = []*rangeNode***REMOVED***r***REMOVED***
	***REMOVED***

enumRanges : enumRanges ',' enumRange ***REMOVED***
		$$ = append($1, $3...)
	***REMOVED***
	| enumRange

enumRange : intLit ***REMOVED***
		r := &rangeNode***REMOVED***startNode: $1***REMOVED***
		r.setRange($1, $1)
		$$ = []*rangeNode***REMOVED***r***REMOVED***
	***REMOVED***
	| intLit _TO intLit ***REMOVED***
		r := &rangeNode***REMOVED***startNode: $1, endNode: $3***REMOVED***
		r.setRange($1, $3)
		$$ = []*rangeNode***REMOVED***r***REMOVED***
	***REMOVED***
	| intLit _TO _MAX ***REMOVED***
		r := &rangeNode***REMOVED***startNode: $1, endNode: $3, endMax: true***REMOVED***
		r.setRange($1, $3)
		$$ = []*rangeNode***REMOVED***r***REMOVED***
	***REMOVED***

intLit : _INT_LIT ***REMOVED***
		i := &compoundIntNode***REMOVED***val: int64($1.val)***REMOVED***
		i.setRange($1, $1)
		$$ = i
	***REMOVED***
	| '-' _INT_LIT ***REMOVED***
		if $2.val > math.MaxInt64 + 1 ***REMOVED***
			lexError(protolex, $2.start(), fmt.Sprintf("numeric constant %d would underflow 64-bit signed int (allowed range is %d to %d)", $2.val, int64(math.MinInt64), int64(math.MaxInt64)))
		***REMOVED***
		i := &compoundIntNode***REMOVED***val: -int64($2.val)***REMOVED***
		i.setRange($1, $2)
		$$ = i
	***REMOVED***

msgReserved : _RESERVED tagRanges ';' ***REMOVED***
		$$ = &reservedNode***REMOVED***ranges: $2***REMOVED***
		$$.setRange($1, $3)
	***REMOVED***
	| reservedNames

enumReserved : _RESERVED enumRanges ';' ***REMOVED***
		$$ = &reservedNode***REMOVED***ranges: $2***REMOVED***
		$$.setRange($1, $3)
	***REMOVED***
	| reservedNames

reservedNames : _RESERVED fieldNames ';' ***REMOVED***
		$$ = &reservedNode***REMOVED***names: $2***REMOVED***
		$$.setRange($1, $3)
	***REMOVED***

fieldNames : fieldNames ',' stringLit ***REMOVED***
		$$ = append($1, $3)
	***REMOVED***
	| stringLit ***REMOVED***
		$$ = []*compoundStringNode***REMOVED***$1***REMOVED***
	***REMOVED***

enum : _ENUM name '***REMOVED***' enumBody '***REMOVED***' ***REMOVED***
		$$ = &enumNode***REMOVED***name: $2, decls: $4***REMOVED***
		$$.setRange($1, $5)
	***REMOVED***

enumBody : enumBody enumItem ***REMOVED***
		$$ = append($1, $2...)
	***REMOVED***
	| enumItem
	| ***REMOVED***
		$$ = nil
	***REMOVED***

enumItem : option ***REMOVED***
		$$ = []*enumElement***REMOVED******REMOVED***option: $1[0]***REMOVED******REMOVED***
	***REMOVED***
	| enumField ***REMOVED***
		$$ = []*enumElement***REMOVED******REMOVED***value: $1***REMOVED******REMOVED***
	***REMOVED***
	| enumReserved ***REMOVED***
		$$ = []*enumElement***REMOVED******REMOVED***reserved: $1***REMOVED******REMOVED***
	***REMOVED***
	| ';' ***REMOVED***
		$$ = []*enumElement***REMOVED******REMOVED***empty: $1***REMOVED******REMOVED***
	***REMOVED***
	| error ';' ***REMOVED***
	***REMOVED***
	| error ***REMOVED***
	***REMOVED***

enumField : name '=' intLit ';' ***REMOVED***
		$$ = &enumValueNode***REMOVED***name: $1, number: $3***REMOVED***
		$$.setRange($1, $4)
	***REMOVED***
	|  name '=' intLit compactOptions ';' ***REMOVED***
		$$ = &enumValueNode***REMOVED***name: $1, number: $3, options: $4***REMOVED***
		$$.setRange($1, $5)
	***REMOVED***

message : _MESSAGE name '***REMOVED***' messageBody '***REMOVED***' ***REMOVED***
		$$ = &messageNode***REMOVED***name: $2, decls: $4***REMOVED***
		$$.setRange($1, $5)
	***REMOVED***

messageBody : messageBody messageItem ***REMOVED***
		$$ = append($1, $2...)
	***REMOVED***
	| messageItem
	| ***REMOVED***
		$$ = nil
	***REMOVED***

messageItem : field ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***field: $1***REMOVED******REMOVED***
	***REMOVED***
	| enum ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***enum: $1***REMOVED******REMOVED***
	***REMOVED***
	| message ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***nested: $1***REMOVED******REMOVED***
	***REMOVED***
	| extend ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***extend: $1***REMOVED******REMOVED***
	***REMOVED***
	| extensions ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***extensionRange: $1***REMOVED******REMOVED***
	***REMOVED***
	| group ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***group: $1***REMOVED******REMOVED***
	***REMOVED***
	| option ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***option: $1[0]***REMOVED******REMOVED***
	***REMOVED***
	| oneof ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***oneOf: $1***REMOVED******REMOVED***
	***REMOVED***
	| mapField ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***mapField: $1***REMOVED******REMOVED***
	***REMOVED***
	| msgReserved ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***reserved: $1***REMOVED******REMOVED***
	***REMOVED***
	| ';' ***REMOVED***
		$$ = []*messageElement***REMOVED******REMOVED***empty: $1***REMOVED******REMOVED***
	***REMOVED***
	| error ';' ***REMOVED***
	***REMOVED***
	| error ***REMOVED***
	***REMOVED***

extend : _EXTEND typeIdent '***REMOVED***' extendBody '***REMOVED***' ***REMOVED***
		$$ = &extendNode***REMOVED***extendee: $2, decls: $4***REMOVED***
		$$.setRange($1, $5)
	***REMOVED***

extendBody : extendBody extendItem ***REMOVED***
		$$ = append($1, $2...)
	***REMOVED***
	| extendItem
	| ***REMOVED***
		$$ = nil
	***REMOVED***

extendItem : field ***REMOVED***
		$$ = []*extendElement***REMOVED******REMOVED***field: $1***REMOVED******REMOVED***
	***REMOVED***
	| group ***REMOVED***
		$$ = []*extendElement***REMOVED******REMOVED***group: $1***REMOVED******REMOVED***
	***REMOVED***
	| ';' ***REMOVED***
		$$ = []*extendElement***REMOVED******REMOVED***empty: $1***REMOVED******REMOVED***
	***REMOVED***
	| error ';' ***REMOVED***
	***REMOVED***
	| error ***REMOVED***
	***REMOVED***

service : _SERVICE name '***REMOVED***' serviceBody '***REMOVED***' ***REMOVED***
		$$ = &serviceNode***REMOVED***name: $2, decls: $4***REMOVED***
		$$.setRange($1, $5)
	***REMOVED***

serviceBody : serviceBody serviceItem ***REMOVED***
		$$ = append($1, $2...)
	***REMOVED***
	| serviceItem
	| ***REMOVED***
		$$ = nil
	***REMOVED***

// NB: doc suggests support for "stream" declaration, separate from "rpc", but
// it does not appear to be supported in protoc (doc is likely from grammar for
// Google-internal version of protoc, with support for streaming stubby)
serviceItem : option ***REMOVED***
		$$ = []*serviceElement***REMOVED******REMOVED***option: $1[0]***REMOVED******REMOVED***
	***REMOVED***
	| rpc ***REMOVED***
		$$ = []*serviceElement***REMOVED******REMOVED***rpc: $1***REMOVED******REMOVED***
	***REMOVED***
	| ';' ***REMOVED***
		$$ = []*serviceElement***REMOVED******REMOVED***empty: $1***REMOVED******REMOVED***
	***REMOVED***
	| error ';' ***REMOVED***
	***REMOVED***
	| error ***REMOVED***
	***REMOVED***

rpc : _RPC name '(' rpcType ')' _RETURNS '(' rpcType ')' ';' ***REMOVED***
		$$ = &methodNode***REMOVED***name: $2, input: $4, output: $8***REMOVED***
		$$.setRange($1, $10)
	***REMOVED***
	| _RPC name '(' rpcType ')' _RETURNS '(' rpcType ')' '***REMOVED***' rpcOptions '***REMOVED***' ***REMOVED***
		$$ = &methodNode***REMOVED***name: $2, input: $4, output: $8, options: $11***REMOVED***
		$$.setRange($1, $12)
	***REMOVED***

rpcType : _STREAM typeIdent ***REMOVED***
		$$ = &rpcTypeNode***REMOVED***msgType: $2, streamKeyword: $1***REMOVED***
		$$.setRange($1, $2)
	***REMOVED***
	| typeIdent ***REMOVED***
		$$ = &rpcTypeNode***REMOVED***msgType: $1***REMOVED***
		$$.setRange($1, $1)
	***REMOVED***

rpcOptions : rpcOptions rpcOption ***REMOVED***
		$$ = append($1, $2...)
	***REMOVED***
	| rpcOption
	| ***REMOVED***
		$$ = []*optionNode***REMOVED******REMOVED***
	***REMOVED***

rpcOption : option ***REMOVED***
		$$ = $1
	***REMOVED***
	| ';' ***REMOVED***
		$$ = []*optionNode***REMOVED******REMOVED***
	***REMOVED***
	| error ';' ***REMOVED***
	***REMOVED***
	| error ***REMOVED***
	***REMOVED***

name : _NAME
	| _SYNTAX
	| _IMPORT
	| _WEAK
	| _PUBLIC
	| _PACKAGE
	| _OPTION
	| _TRUE
	| _FALSE
	| _INF
	| _NAN
	| _REPEATED
	| _OPTIONAL
	| _REQUIRED
	| _DOUBLE
	| _FLOAT
	| _INT32
	| _INT64
	| _UINT32
	| _UINT64
	| _SINT32
	| _SINT64
	| _FIXED32
	| _FIXED64
	| _SFIXED32
	| _SFIXED64
	| _BOOL
	| _STRING
	| _BYTES
	| _GROUP
	| _ONEOF
	| _MAP
	| _EXTENSIONS
	| _TO
	| _MAX
	| _RESERVED
	| _ENUM
	| _MESSAGE
	| _EXTEND
	| _SERVICE
	| _RPC
	| _STREAM
	| _RETURNS

%%
