%***REMOVED***
package protoparse

//lint:file-ignore SA4006 generated parser has unused values

import (
	"math"

	"github.com/jhump/protoreflect/desc/protoparse/ast"
)

%***REMOVED***

// fields inside this union end up as the fields in a structure known
// as $***REMOVED***PREFIX***REMOVED***SymType, of which a reference is passed to the lexer.
%union***REMOVED***
	file      *ast.FileNode
	syn       *ast.SyntaxNode
	fileDecl  ast.FileElement
	fileDecls []ast.FileElement
	pkg       *ast.PackageNode
	imprt     *ast.ImportNode
	msg       *ast.MessageNode
	msgDecl   ast.MessageElement
	msgDecls  []ast.MessageElement
	fld       *ast.FieldNode
	mapFld    *ast.MapFieldNode
	mapType   *ast.MapTypeNode
	grp       *ast.GroupNode
	oo        *ast.OneOfNode
	ooDecl    ast.OneOfElement
	ooDecls   []ast.OneOfElement
	ext       *ast.ExtensionRangeNode
	resvd     *ast.ReservedNode
	en        *ast.EnumNode
	enDecl    ast.EnumElement
	enDecls   []ast.EnumElement
	env       *ast.EnumValueNode
	extend    *ast.ExtendNode
	extDecl   ast.ExtendElement
	extDecls  []ast.ExtendElement
	svc       *ast.ServiceNode
	svcDecl   ast.ServiceElement
	svcDecls  []ast.ServiceElement
	mtd       *ast.RPCNode
	rpcType   *ast.RPCTypeNode
	rpcDecl   ast.RPCElement
	rpcDecls  []ast.RPCElement
	opt       *ast.OptionNode
	opts      *compactOptionList
	ref       *ast.FieldReferenceNode
	optNms    *fieldRefList
	cmpctOpts *ast.CompactOptionsNode
	rng       *ast.RangeNode
	rngs      *rangeList
	names     *nameList
	cid       *identList
	tid       ast.IdentValueNode
	sl        *valueList
	msgField  *ast.MessageFieldNode
	msgEntry  *messageFieldEntry
	msgLit    *messageFieldList
	v         ast.ValueNode
	il        ast.IntValueNode
	str       *stringList
	s         *ast.StringLiteralNode
	i         *ast.UintLiteralNode
	f         *ast.FloatLiteralNode
	id        *ast.IdentNode
	b         *ast.RuneNode
	err       error
***REMOVED***

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <file>      file
%type <syn>       syntax
%type <fileDecl>  fileDecl
%type <fileDecls> fileDecls
%type <imprt>     import
%type <pkg>       package
%type <opt>       option compactOption
%type <opts>      compactOptionDecls
%type <rpcDecl>   rpcDecl
%type <rpcDecls>  rpcDecls
%type <ref>       optionNameComponent aggName
%type <optNms>    optionName
%type <cmpctOpts> compactOptions
%type <v>         constant scalarConstant aggregate numLit
%type <il>        intLit
%type <id>        name keyType msgElementName extElementName oneofElementName enumElementName
%type <cid>       ident msgElementIdent extElementIdent oneofElementIdent
%type <tid>       typeIdent msgElementTypeIdent extElementTypeIdent oneofElementTypeIdent
%type <sl>        constantList
%type <msgField>  aggFieldEntry
%type <msgEntry>  aggField
%type <msgLit>    aggFields
%type <fld>       oneofField msgField extField
%type <oo>        oneof
%type <grp>       group oneofGroup
%type <mapFld>    mapField
%type <mapType>   mapType
%type <msg>       message
%type <msgDecl>   messageDecl
%type <msgDecls>  messageDecls
%type <ooDecl>    ooDecl
%type <ooDecls>   ooDecls
%type <names>     fieldNames
%type <resvd>     msgReserved enumReserved reservedNames
%type <rng>       tagRange enumRange
%type <rngs>      tagRanges enumRanges
%type <ext>       extensions
%type <en>        enum
%type <enDecl>    enumDecl
%type <enDecls>   enumDecls
%type <env>       enumValue
%type <extend>    extend
%type <extDecl>   extendDecl
%type <extDecls>  extendDecls
%type <str>       stringLit
%type <svc>       service
%type <svcDecl>   serviceDecl
%type <svcDecls>  serviceDecls
%type <mtd>       rpc
%type <rpcType>   rpcType

// same for terminals
%token <s>   _STRING_LIT
%token <i>   _INT_LIT
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
		$$ = ast.NewFileNode($1, nil)
		protolex.(*protoLex).res = $$
	***REMOVED***
	| fileDecls  ***REMOVED***
		$$ = ast.NewFileNode(nil, $1)
		protolex.(*protoLex).res = $$
	***REMOVED***
	| syntax fileDecls ***REMOVED***
		$$ = ast.NewFileNode($1, $2)
		protolex.(*protoLex).res = $$
	***REMOVED***
	| ***REMOVED***
	***REMOVED***

fileDecls : fileDecls fileDecl ***REMOVED***
		if $2 != nil ***REMOVED***
			$$ = append($1, $2)
		***REMOVED*** else ***REMOVED***
			$$ = $1
		***REMOVED***
	***REMOVED***
	| fileDecl ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = []ast.FileElement***REMOVED***$1***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***

fileDecl : import ***REMOVED***
		$$ = $1
	***REMOVED***
	| package ***REMOVED***
		$$ = $1
	***REMOVED***
	| option ***REMOVED***
		$$ = $1
	***REMOVED***
	| message ***REMOVED***
		$$ = $1
	***REMOVED***
	| enum ***REMOVED***
		$$ = $1
	***REMOVED***
	| extend ***REMOVED***
		$$ = $1
	***REMOVED***
	| service ***REMOVED***
		$$ = $1
	***REMOVED***
	| ';' ***REMOVED***
		$$ = ast.NewEmptyDeclNode($1)
	***REMOVED***
	| error ';' ***REMOVED***
		$$ = nil
	***REMOVED***
	| error ***REMOVED***
		$$ = nil
	***REMOVED***

syntax : _SYNTAX '=' stringLit ';' ***REMOVED***
		$$ = ast.NewSyntaxNode($1.ToKeyword(), $2, $3.toStringValueNode(), $4)
	***REMOVED***

import : _IMPORT stringLit ';' ***REMOVED***
		$$ = ast.NewImportNode($1.ToKeyword(), nil, nil, $2.toStringValueNode(), $3)
	***REMOVED***
	| _IMPORT _WEAK stringLit ';' ***REMOVED***
		$$ = ast.NewImportNode($1.ToKeyword(), nil, $2.ToKeyword(), $3.toStringValueNode(), $4)
	***REMOVED***
	| _IMPORT _PUBLIC stringLit ';' ***REMOVED***
		$$ = ast.NewImportNode($1.ToKeyword(), $2.ToKeyword(), nil, $3.toStringValueNode(), $4)
	***REMOVED***

package : _PACKAGE ident ';' ***REMOVED***
		$$ = ast.NewPackageNode($1.ToKeyword(), $2.toIdentValueNode(nil), $3)
	***REMOVED***

ident : name ***REMOVED***
		$$ = &identList***REMOVED***$1, nil, nil***REMOVED***
	***REMOVED***
	| name '.' ident ***REMOVED***
		$$ = &identList***REMOVED***$1, $2, $3***REMOVED***
	***REMOVED***

// to mimic limitations of protoc recursive-descent parser,
// we don't allowed message statement keywords as identifiers
// (or oneof statement keywords [e.g. "option"] below)

msgElementIdent : msgElementName ***REMOVED***
		$$ = &identList***REMOVED***$1, nil, nil***REMOVED***
	***REMOVED***
	| msgElementName '.' ident ***REMOVED***
		$$ = &identList***REMOVED***$1, $2, $3***REMOVED***
	***REMOVED***

extElementIdent : extElementName ***REMOVED***
		$$ = &identList***REMOVED***$1, nil, nil***REMOVED***
	***REMOVED***
	| extElementName '.' ident ***REMOVED***
		$$ = &identList***REMOVED***$1, $2, $3***REMOVED***
	***REMOVED***

oneofElementIdent : oneofElementName ***REMOVED***
		$$ = &identList***REMOVED***$1, nil, nil***REMOVED***
	***REMOVED***
	| oneofElementName '.' ident ***REMOVED***
		$$ = &identList***REMOVED***$1, $2, $3***REMOVED***
	***REMOVED***

option : _OPTION optionName '=' constant ';' ***REMOVED***
		refs, dots := $2.toNodes()
		optName := ast.NewOptionNameNode(refs, dots)
		$$ = ast.NewOptionNode($1.ToKeyword(), optName, $3, $4, $5)
	***REMOVED***

optionName : optionNameComponent ***REMOVED***
		$$ = &fieldRefList***REMOVED***$1, nil, nil***REMOVED***
	***REMOVED***
	| optionNameComponent '.' optionName ***REMOVED***
		$$ = &fieldRefList***REMOVED***$1, $2, $3***REMOVED***
	***REMOVED***

optionNameComponent : name ***REMOVED***
		$$ = ast.NewFieldReferenceNode($1)
	***REMOVED***
	| '(' typeIdent ')' ***REMOVED***
		$$ = ast.NewExtensionFieldReferenceNode($1, $2, $3)
	***REMOVED***

constant : scalarConstant
	| aggregate

scalarConstant : stringLit ***REMOVED***
		$$ = $1.toStringValueNode()
	***REMOVED***
	| numLit
	| name ***REMOVED***
        $$ = $1
	***REMOVED***

numLit : _FLOAT_LIT ***REMOVED***
		$$ = $1
	***REMOVED***
	| '-' _FLOAT_LIT ***REMOVED***
		$$ = ast.NewSignedFloatLiteralNode($1, $2)
	***REMOVED***
	| '+' _FLOAT_LIT ***REMOVED***
		$$ = ast.NewSignedFloatLiteralNode($1, $2)
	***REMOVED***
	| '+' _INF ***REMOVED***
		f := ast.NewSpecialFloatLiteralNode($2.ToKeyword())
		$$ = ast.NewSignedFloatLiteralNode($1, f)
	***REMOVED***
	| '-' _INF ***REMOVED***
		f := ast.NewSpecialFloatLiteralNode($2.ToKeyword())
		$$ = ast.NewSignedFloatLiteralNode($1, f)
	***REMOVED***
	| _INT_LIT ***REMOVED***
		$$ = $1
	***REMOVED***
	| '+' _INT_LIT ***REMOVED***
		$$ = ast.NewPositiveUintLiteralNode($1, $2)
	***REMOVED***
	| '-' _INT_LIT ***REMOVED***
		if $2.Val > math.MaxInt64 + 1 ***REMOVED***
			// can't represent as int so treat as float literal
			$$ = ast.NewSignedFloatLiteralNode($1, $2)
		***REMOVED*** else ***REMOVED***
			$$ = ast.NewNegativeIntLiteralNode($1, $2)
		***REMOVED***
	***REMOVED***

stringLit : _STRING_LIT ***REMOVED***
		$$ = &stringList***REMOVED***$1, nil***REMOVED***
	***REMOVED***
	| _STRING_LIT stringLit  ***REMOVED***
		$$ = &stringList***REMOVED***$1, $2***REMOVED***
	***REMOVED***

aggregate : '***REMOVED***' aggFields '***REMOVED***' ***REMOVED***
		fields, delims := $2.toNodes()
		$$ = ast.NewMessageLiteralNode($1, fields, delims, $3)
	***REMOVED***

aggFields : aggField ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = &messageFieldList***REMOVED***$1, nil***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggField aggFields ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = &messageFieldList***REMOVED***$1, $2***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = $2
		***REMOVED***
	***REMOVED***
	| ***REMOVED***
		$$ = nil
	***REMOVED***

aggField : aggFieldEntry ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = &messageFieldEntry***REMOVED***$1, nil***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggFieldEntry ',' ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = &messageFieldEntry***REMOVED***$1, $2***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggFieldEntry ';' ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = &messageFieldEntry***REMOVED***$1, $2***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| error ',' ***REMOVED***
		$$ = nil
	***REMOVED***
	| error ';' ***REMOVED***
		$$ = nil
	***REMOVED***
	| error ***REMOVED***
		$$ = nil
	***REMOVED***

aggFieldEntry : aggName ':' scalarConstant ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = ast.NewMessageFieldNode($1, $2, $3)
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggName '[' ']' ***REMOVED***
		if $1 != nil ***REMOVED***
			val := ast.NewArrayLiteralNode($2, nil, nil, $3)
			$$ = ast.NewMessageFieldNode($1, nil, val)
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggName ':' '[' ']' ***REMOVED***
		if $1 != nil ***REMOVED***
			val := ast.NewArrayLiteralNode($3, nil, nil, $4)
			$$ = ast.NewMessageFieldNode($1, $2, val)
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggName '[' constantList ']' ***REMOVED***
		if $1 != nil ***REMOVED***
			vals, commas := $3.toNodes()
			val := ast.NewArrayLiteralNode($2, vals, commas, $4)
			$$ = ast.NewMessageFieldNode($1, nil, val)
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggName ':' '[' constantList ']' ***REMOVED***
		if $1 != nil ***REMOVED***
			vals, commas := $4.toNodes()
			val := ast.NewArrayLiteralNode($3, vals, commas, $5)
			$$ = ast.NewMessageFieldNode($1, $2, val)
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggName ':' '[' error ']' ***REMOVED***
		$$ = nil
	***REMOVED***
	| aggName ':' aggregate ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = ast.NewMessageFieldNode($1, $2, $3)
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggName aggregate ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = ast.NewMessageFieldNode($1, nil, $2)
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggName ':' '<' aggFields '>' ***REMOVED***
		if $1 != nil ***REMOVED***
			fields, delims := $4.toNodes()
			msg := ast.NewMessageLiteralNode($3, fields, delims, $5)
			$$ = ast.NewMessageFieldNode($1, $2, msg)
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggName '<' aggFields '>' ***REMOVED***
		if $1 != nil ***REMOVED***
			fields, delims := $3.toNodes()
			msg := ast.NewMessageLiteralNode($2, fields, delims, $4)
			$$ = ast.NewMessageFieldNode($1, nil, msg)
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| aggName ':' '<' error '>' ***REMOVED***
		$$ = nil
	***REMOVED***
	| aggName '<' error '>' ***REMOVED***
		$$ = nil
	***REMOVED***

aggName : name ***REMOVED***
		$$ = ast.NewFieldReferenceNode($1)
	***REMOVED***
	| '[' ident ']' ***REMOVED***
		$$ = ast.NewExtensionFieldReferenceNode($1, $2.toIdentValueNode(nil), $3)
	***REMOVED***
	| '[' ident '/' ident ']' ***REMOVED***
		$$ = ast.NewAnyTypeReferenceNode($1, $2.toIdentValueNode(nil), $3, $4.toIdentValueNode(nil), $5)
	***REMOVED***
	| '[' error ']' ***REMOVED***
		$$ = nil
	***REMOVED***

constantList : constant ***REMOVED***
		$$ = &valueList***REMOVED***$1, nil, nil***REMOVED***
	***REMOVED***
	| constant ',' constantList ***REMOVED***
		$$ = &valueList***REMOVED***$1, $2, $3***REMOVED***
	***REMOVED***
	| '<' aggFields '>' ***REMOVED***
		fields, delims := $2.toNodes()
		msg := ast.NewMessageLiteralNode($1, fields, delims, $3)
		$$ = &valueList***REMOVED***msg, nil, nil***REMOVED***
	***REMOVED***
	| '<' aggFields '>' ',' constantList ***REMOVED***
		fields, delims := $2.toNodes()
		msg := ast.NewMessageLiteralNode($1, fields, delims, $3)
		$$ = &valueList***REMOVED***msg, $4, $5***REMOVED***
	***REMOVED***
	| '<' error '>' ***REMOVED***
		$$ = nil
	***REMOVED***
	| '<' error '>' ',' constantList ***REMOVED***
		$$ = $5
	***REMOVED***

typeIdent : ident ***REMOVED***
		$$ = $1.toIdentValueNode(nil)
	***REMOVED***
	| '.' ident ***REMOVED***
		$$ = $2.toIdentValueNode($1)
	***REMOVED***

msgElementTypeIdent : msgElementIdent ***REMOVED***
		$$ = $1.toIdentValueNode(nil)
	***REMOVED***
	| '.' ident ***REMOVED***
		$$ = $2.toIdentValueNode($1)
	***REMOVED***

extElementTypeIdent : extElementIdent ***REMOVED***
		$$ = $1.toIdentValueNode(nil)
	***REMOVED***
	| '.' ident ***REMOVED***
		$$ = $2.toIdentValueNode($1)
	***REMOVED***

oneofElementTypeIdent : oneofElementIdent ***REMOVED***
		$$ = $1.toIdentValueNode(nil)
	***REMOVED***
	| '.' ident ***REMOVED***
		$$ = $2.toIdentValueNode($1)
	***REMOVED***

msgField : _REQUIRED typeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, nil, $6)
	***REMOVED***
	| _OPTIONAL typeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, nil, $6)
	***REMOVED***
	| _REPEATED typeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, nil, $6)
	***REMOVED***
	| _REQUIRED typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, $6, $7)
	***REMOVED***
	| _OPTIONAL typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, $6, $7)
	***REMOVED***
	| _REPEATED typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, $6, $7)
	***REMOVED***
	| msgElementTypeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewFieldNode(nil, $1, $2, $3, $4, nil, $5)
	***REMOVED***
	| msgElementTypeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewFieldNode(nil, $1, $2, $3, $4, $5, $6)
	***REMOVED***

extField : _REQUIRED typeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, nil, $6)
	***REMOVED***
	| _OPTIONAL typeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, nil, $6)
	***REMOVED***
	| _REPEATED typeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, nil, $6)
	***REMOVED***
	| _REQUIRED typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, $6, $7)
	***REMOVED***
	| _OPTIONAL typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, $6, $7)
	***REMOVED***
	| _REPEATED typeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewFieldNode($1.ToKeyword(), $2, $3, $4, $5, $6, $7)
	***REMOVED***
	| extElementTypeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewFieldNode(nil, $1, $2, $3, $4, nil, $5)
	***REMOVED***
	| extElementTypeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewFieldNode(nil, $1, $2, $3, $4, $5, $6)
	***REMOVED***

compactOptions: '[' compactOptionDecls ']' ***REMOVED***
		opts, commas := $2.toNodes()
		$$ = ast.NewCompactOptionsNode($1, opts, commas, $3)
	***REMOVED***

compactOptionDecls : compactOption ***REMOVED***
		$$ = &compactOptionList***REMOVED***$1, nil, nil***REMOVED***
	***REMOVED***
	| compactOption ',' compactOptionDecls ***REMOVED***
		$$ = &compactOptionList***REMOVED***$1, $2, $3***REMOVED***
	***REMOVED***

compactOption: optionName '=' constant ***REMOVED***
		refs, dots := $1.toNodes()
		optName := ast.NewOptionNameNode(refs, dots)
		$$ = ast.NewCompactOptionNode(optName, $2, $3)
	***REMOVED***

group : _REQUIRED _GROUP name '=' _INT_LIT '***REMOVED***' messageDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewGroupNode($1.ToKeyword(), $2.ToKeyword(), $3, $4, $5, nil, $6, $7, $8)
	***REMOVED***
	| _OPTIONAL _GROUP name '=' _INT_LIT '***REMOVED***' messageDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewGroupNode($1.ToKeyword(), $2.ToKeyword(), $3, $4, $5, nil, $6, $7, $8)
	***REMOVED***
	| _REPEATED _GROUP name '=' _INT_LIT '***REMOVED***' messageDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewGroupNode($1.ToKeyword(), $2.ToKeyword(), $3, $4, $5, nil, $6, $7, $8)
	***REMOVED***
	| _REQUIRED _GROUP name '=' _INT_LIT compactOptions '***REMOVED***' messageDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewGroupNode($1.ToKeyword(), $2.ToKeyword(), $3, $4, $5, $6, $7, $8, $9)
	***REMOVED***
	| _OPTIONAL _GROUP name '=' _INT_LIT compactOptions '***REMOVED***' messageDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewGroupNode($1.ToKeyword(), $2.ToKeyword(), $3, $4, $5, $6, $7, $8, $9)
	***REMOVED***
	| _REPEATED _GROUP name '=' _INT_LIT compactOptions '***REMOVED***' messageDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewGroupNode($1.ToKeyword(), $2.ToKeyword(), $3, $4, $5, $6, $7, $8, $9)
	***REMOVED***

oneof : _ONEOF name '***REMOVED***' ooDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewOneOfNode($1.ToKeyword(), $2, $3, $4, $5)
	***REMOVED***

ooDecls : ooDecls ooDecl ***REMOVED***
		if $2 != nil ***REMOVED***
			$$ = append($1, $2)
		***REMOVED*** else ***REMOVED***
			$$ = $1
		***REMOVED***
	***REMOVED***
	| ooDecl ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = []ast.OneOfElement***REMOVED***$1***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| ***REMOVED***
		$$ = nil
	***REMOVED***

ooDecl : option ***REMOVED***
		$$ = $1
	***REMOVED***
	| oneofField ***REMOVED***
		$$ = $1
	***REMOVED***
	| oneofGroup ***REMOVED***
		$$ = $1
	***REMOVED***
	| ';' ***REMOVED***
		$$ = ast.NewEmptyDeclNode($1)
	***REMOVED***
	| error ';' ***REMOVED***
		$$ = nil
	***REMOVED***
	| error ***REMOVED***
		$$ = nil
	***REMOVED***

oneofField : oneofElementTypeIdent name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewFieldNode(nil, $1, $2, $3, $4, nil, $5)
	***REMOVED***
	| oneofElementTypeIdent name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewFieldNode(nil, $1, $2, $3, $4, $5, $6)
	***REMOVED***

oneofGroup : _GROUP name '=' _INT_LIT '***REMOVED***' messageDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewGroupNode(nil, $1.ToKeyword(), $2, $3, $4, nil, $5, $6, $7)
	***REMOVED***
	| _GROUP name '=' _INT_LIT compactOptions '***REMOVED***' messageDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewGroupNode(nil, $1.ToKeyword(), $2, $3, $4, $5, $6, $7, $8)
	***REMOVED***

mapField : mapType name '=' _INT_LIT ';' ***REMOVED***
		$$ = ast.NewMapFieldNode($1, $2, $3, $4, nil, $5)
	***REMOVED***
	| mapType name '=' _INT_LIT compactOptions ';' ***REMOVED***
		$$ = ast.NewMapFieldNode($1, $2, $3, $4, $5, $6)
	***REMOVED***

mapType : _MAP '<' keyType ',' typeIdent '>' ***REMOVED***
		$$ = ast.NewMapTypeNode($1.ToKeyword(), $2, $3, $4, $5, $6)
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
		ranges, commas := $2.toNodes()
		$$ = ast.NewExtensionRangeNode($1.ToKeyword(), ranges, commas, nil, $3)
	***REMOVED***
	| _EXTENSIONS tagRanges compactOptions ';' ***REMOVED***
		ranges, commas := $2.toNodes()
		$$ = ast.NewExtensionRangeNode($1.ToKeyword(), ranges, commas, $3, $4)
	***REMOVED***

tagRanges : tagRange ***REMOVED***
		$$ = &rangeList***REMOVED***$1, nil, nil***REMOVED***
	***REMOVED***
	| tagRange ',' tagRanges ***REMOVED***
		$$ = &rangeList***REMOVED***$1, $2, $3***REMOVED***
	***REMOVED***

tagRange : _INT_LIT ***REMOVED***
		$$ = ast.NewRangeNode($1, nil, nil, nil)
	***REMOVED***
	| _INT_LIT _TO _INT_LIT ***REMOVED***
		$$ = ast.NewRangeNode($1, $2.ToKeyword(), $3, nil)
	***REMOVED***
	| _INT_LIT _TO _MAX ***REMOVED***
		$$ = ast.NewRangeNode($1, $2.ToKeyword(), nil, $3.ToKeyword())
	***REMOVED***

enumRanges : enumRange ***REMOVED***
		$$ = &rangeList***REMOVED***$1, nil, nil***REMOVED***
	***REMOVED***
	| enumRange ',' enumRanges ***REMOVED***
		$$ = &rangeList***REMOVED***$1, $2, $3***REMOVED***
	***REMOVED***

enumRange : intLit ***REMOVED***
		$$ = ast.NewRangeNode($1, nil, nil, nil)
	***REMOVED***
	| intLit _TO intLit ***REMOVED***
		$$ = ast.NewRangeNode($1, $2.ToKeyword(), $3, nil)
	***REMOVED***
	| intLit _TO _MAX ***REMOVED***
		$$ = ast.NewRangeNode($1, $2.ToKeyword(), nil, $3.ToKeyword())
	***REMOVED***

intLit : _INT_LIT ***REMOVED***
		$$ = $1
	***REMOVED***
	| '-' _INT_LIT ***REMOVED***
		$$ = ast.NewNegativeIntLiteralNode($1, $2)
	***REMOVED***

msgReserved : _RESERVED tagRanges ';' ***REMOVED***
		ranges, commas := $2.toNodes()
		$$ = ast.NewReservedRangesNode($1.ToKeyword(), ranges, commas, $3)
	***REMOVED***
	| reservedNames

enumReserved : _RESERVED enumRanges ';' ***REMOVED***
		ranges, commas := $2.toNodes()
		$$ = ast.NewReservedRangesNode($1.ToKeyword(), ranges, commas, $3)
	***REMOVED***
	| reservedNames

reservedNames : _RESERVED fieldNames ';' ***REMOVED***
		names, commas := $2.toNodes()
		$$ = ast.NewReservedNamesNode($1.ToKeyword(), names, commas, $3)
	***REMOVED***

fieldNames : stringLit ***REMOVED***
		$$ = &nameList***REMOVED***$1.toStringValueNode(), nil, nil***REMOVED***
	***REMOVED***
	| stringLit ',' fieldNames ***REMOVED***
		$$ = &nameList***REMOVED***$1.toStringValueNode(), $2, $3***REMOVED***
	***REMOVED***

enum : _ENUM name '***REMOVED***' enumDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewEnumNode($1.ToKeyword(), $2, $3, $4, $5)
	***REMOVED***

enumDecls : enumDecls enumDecl ***REMOVED***
		if $2 != nil ***REMOVED***
			$$ = append($1, $2)
		***REMOVED*** else ***REMOVED***
			$$ = $1
		***REMOVED***
	***REMOVED***
	| enumDecl ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = []ast.EnumElement***REMOVED***$1***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| ***REMOVED***
		$$ = nil
	***REMOVED***

enumDecl : option ***REMOVED***
		$$ = $1
	***REMOVED***
	| enumValue ***REMOVED***
		$$ = $1
	***REMOVED***
	| enumReserved ***REMOVED***
		$$ = $1
	***REMOVED***
	| ';' ***REMOVED***
		$$ = ast.NewEmptyDeclNode($1)
	***REMOVED***
	| error ';' ***REMOVED***
		$$ = nil
	***REMOVED***
	| error ***REMOVED***
		$$ = nil
	***REMOVED***

enumValue : enumElementName '=' intLit ';' ***REMOVED***
		$$ = ast.NewEnumValueNode($1, $2, $3, nil, $4)
	***REMOVED***
	|  enumElementName '=' intLit compactOptions ';' ***REMOVED***
		$$ = ast.NewEnumValueNode($1, $2, $3, $4, $5)
	***REMOVED***

message : _MESSAGE name '***REMOVED***' messageDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewMessageNode($1.ToKeyword(), $2, $3, $4, $5)
	***REMOVED***

messageDecls : messageDecls messageDecl ***REMOVED***
		if $2 != nil ***REMOVED***
			$$ = append($1, $2)
		***REMOVED*** else ***REMOVED***
			$$ = $1
		***REMOVED***
	***REMOVED***
	| messageDecl ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = []ast.MessageElement***REMOVED***$1***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| ***REMOVED***
		$$ = nil
	***REMOVED***

messageDecl : msgField ***REMOVED***
		$$ = $1
	***REMOVED***
	| enum ***REMOVED***
		$$ = $1
	***REMOVED***
	| message ***REMOVED***
		$$ = $1
	***REMOVED***
	| extend ***REMOVED***
		$$ = $1
	***REMOVED***
	| extensions ***REMOVED***
		$$ = $1
	***REMOVED***
	| group ***REMOVED***
		$$ = $1
	***REMOVED***
	| option ***REMOVED***
		$$ = $1
	***REMOVED***
	| oneof ***REMOVED***
		$$ = $1
	***REMOVED***
	| mapField ***REMOVED***
		$$ = $1
	***REMOVED***
	| msgReserved ***REMOVED***
		$$ = $1
	***REMOVED***
	| ';' ***REMOVED***
		$$ = ast.NewEmptyDeclNode($1)
	***REMOVED***
	| error ';' ***REMOVED***
		$$ = nil
	***REMOVED***
	| error ***REMOVED***
		$$ = nil
	***REMOVED***

extend : _EXTEND typeIdent '***REMOVED***' extendDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewExtendNode($1.ToKeyword(), $2, $3, $4, $5)
	***REMOVED***

extendDecls : extendDecls extendDecl ***REMOVED***
		if $2 != nil ***REMOVED***
			$$ = append($1, $2)
		***REMOVED*** else ***REMOVED***
			$$ = $1
		***REMOVED***
	***REMOVED***
	| extendDecl ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = []ast.ExtendElement***REMOVED***$1***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| ***REMOVED***
		$$ = nil
	***REMOVED***

extendDecl : extField ***REMOVED***
		$$ = $1
	***REMOVED***
	| group ***REMOVED***
		$$ = $1
	***REMOVED***
	| ';' ***REMOVED***
		$$ = ast.NewEmptyDeclNode($1)
	***REMOVED***
	| error ';' ***REMOVED***
		$$ = nil
	***REMOVED***
	| error ***REMOVED***
		$$ = nil
	***REMOVED***

service : _SERVICE name '***REMOVED***' serviceDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewServiceNode($1.ToKeyword(), $2, $3, $4, $5)
	***REMOVED***

serviceDecls : serviceDecls serviceDecl ***REMOVED***
		if $2 != nil ***REMOVED***
			$$ = append($1, $2)
		***REMOVED*** else ***REMOVED***
			$$ = $1
		***REMOVED***
	***REMOVED***
	| serviceDecl ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = []ast.ServiceElement***REMOVED***$1***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| ***REMOVED***
		$$ = nil
	***REMOVED***

// NB: doc suggests support for "stream" declaration, separate from "rpc", but
// it does not appear to be supported in protoc (doc is likely from grammar for
// Google-internal version of protoc, with support for streaming stubby)
serviceDecl : option ***REMOVED***
		$$ = $1
	***REMOVED***
	| rpc ***REMOVED***
		$$ = $1
	***REMOVED***
	| ';' ***REMOVED***
		$$ = ast.NewEmptyDeclNode($1)
	***REMOVED***
	| error ';' ***REMOVED***
		$$ = nil
	***REMOVED***
	| error ***REMOVED***
		$$ = nil
	***REMOVED***

rpc : _RPC name rpcType _RETURNS rpcType ';' ***REMOVED***
		$$ = ast.NewRPCNode($1.ToKeyword(), $2, $3, $4.ToKeyword(), $5, $6)
	***REMOVED***
	| _RPC name rpcType _RETURNS rpcType '***REMOVED***' rpcDecls '***REMOVED***' ***REMOVED***
		$$ = ast.NewRPCNodeWithBody($1.ToKeyword(), $2, $3, $4.ToKeyword(), $5, $6, $7, $8)
	***REMOVED***

rpcType : '(' _STREAM typeIdent ')' ***REMOVED***
		$$ = ast.NewRPCTypeNode($1, $2.ToKeyword(), $3, $4)
	***REMOVED***
	| '(' typeIdent ')' ***REMOVED***
		$$ = ast.NewRPCTypeNode($1, nil, $2, $3)
	***REMOVED***

rpcDecls : rpcDecls rpcDecl ***REMOVED***
		if $2 != nil ***REMOVED***
			$$ = append($1, $2)
		***REMOVED*** else ***REMOVED***
			$$ = $1
		***REMOVED***
	***REMOVED***
	| rpcDecl ***REMOVED***
		if $1 != nil ***REMOVED***
			$$ = []ast.RPCElement***REMOVED***$1***REMOVED***
		***REMOVED*** else ***REMOVED***
			$$ = nil
		***REMOVED***
	***REMOVED***
	| ***REMOVED***
		$$ = nil
	***REMOVED***

rpcDecl : option ***REMOVED***
		$$ = $1
	***REMOVED***
	| ';' ***REMOVED***
		$$ = ast.NewEmptyDeclNode($1)
	***REMOVED***
	| error ';' ***REMOVED***
		$$ = nil
	***REMOVED***
	| error ***REMOVED***
		$$ = nil
	***REMOVED***

// excludes message, enum, oneof, extensions, reserved, extend,
//   option, optional, required, and repeated
msgElementName : _NAME
	| _SYNTAX
	| _IMPORT
	| _WEAK
	| _PUBLIC
	| _PACKAGE
	| _TRUE
	| _FALSE
	| _INF
	| _NAN
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
	| _MAP
	| _TO
	| _MAX
	| _SERVICE
	| _RPC
	| _STREAM
	| _RETURNS

// excludes optional, required, and repeated
extElementName : _NAME
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

// excludes reserved, option
enumElementName : _NAME
	| _SYNTAX
	| _IMPORT
	| _WEAK
	| _PUBLIC
	| _PACKAGE
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
	| _ENUM
	| _MESSAGE
	| _EXTEND
	| _SERVICE
	| _RPC
	| _STREAM
	| _RETURNS

// excludes option, optional, required, and repeated
oneofElementName : _NAME
	| _SYNTAX
	| _IMPORT
	| _WEAK
	| _PUBLIC
	| _PACKAGE
	| _TRUE
	| _FALSE
	| _INF
	| _NAN
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
