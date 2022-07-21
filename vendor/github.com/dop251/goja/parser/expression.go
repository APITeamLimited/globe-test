package parser

import (
	"strings"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/dop251/goja/unistring"
)

func (self *_parser) parseIdentifier() *ast.Identifier ***REMOVED***
	literal := self.parsedLiteral
	idx := self.idx
	self.next()
	return &ast.Identifier***REMOVED***
		Name: literal,
		Idx:  idx,
	***REMOVED***
***REMOVED***

func (self *_parser) parsePrimaryExpression() ast.Expression ***REMOVED***
	literal, parsedLiteral := self.literal, self.parsedLiteral
	idx := self.idx
	switch self.token ***REMOVED***
	case token.IDENTIFIER:
		self.next()
		return &ast.Identifier***REMOVED***
			Name: parsedLiteral,
			Idx:  idx,
		***REMOVED***
	case token.NULL:
		self.next()
		return &ast.NullLiteral***REMOVED***
			Idx:     idx,
			Literal: literal,
		***REMOVED***
	case token.BOOLEAN:
		self.next()
		value := false
		switch parsedLiteral ***REMOVED***
		case "true":
			value = true
		case "false":
			value = false
		default:
			self.error(idx, "Illegal boolean literal")
		***REMOVED***
		return &ast.BooleanLiteral***REMOVED***
			Idx:     idx,
			Literal: literal,
			Value:   value,
		***REMOVED***
	case token.STRING:
		self.next()
		return &ast.StringLiteral***REMOVED***
			Idx:     idx,
			Literal: literal,
			Value:   parsedLiteral,
		***REMOVED***
	case token.NUMBER:
		self.next()
		value, err := parseNumberLiteral(literal)
		if err != nil ***REMOVED***
			self.error(idx, err.Error())
			value = 0
		***REMOVED***
		return &ast.NumberLiteral***REMOVED***
			Idx:     idx,
			Literal: literal,
			Value:   value,
		***REMOVED***
	case token.SLASH, token.QUOTIENT_ASSIGN:
		return self.parseRegExpLiteral()
	case token.LEFT_BRACE:
		return self.parseObjectLiteral()
	case token.LEFT_BRACKET:
		return self.parseArrayLiteral()
	case token.LEFT_PARENTHESIS:
		return self.parseParenthesisedExpression()
	case token.BACKTICK:
		return self.parseTemplateLiteral(false)
	case token.THIS:
		self.next()
		return &ast.ThisExpression***REMOVED***
			Idx: idx,
		***REMOVED***
	case token.SUPER:
		return self.parseSuperProperty()
	case token.FUNCTION:
		return self.parseFunction(false)
	case token.CLASS:
		return self.parseClass(false)
	***REMOVED***

	if isBindingId(self.token, parsedLiteral) ***REMOVED***
		self.next()
		return &ast.Identifier***REMOVED***
			Name: parsedLiteral,
			Idx:  idx,
		***REMOVED***
	***REMOVED***

	self.errorUnexpectedToken(self.token)
	self.nextStatement()
	return &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
***REMOVED***

func (self *_parser) parseSuperProperty() ast.Expression ***REMOVED***
	idx := self.idx
	self.next()
	switch self.token ***REMOVED***
	case token.PERIOD:
		self.next()
		if !token.IsId(self.token) ***REMOVED***
			self.expect(token.IDENTIFIER)
			self.nextStatement()
			return &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
		***REMOVED***
		idIdx := self.idx
		parsedLiteral := self.parsedLiteral
		self.next()
		return &ast.DotExpression***REMOVED***
			Left: &ast.SuperExpression***REMOVED***
				Idx: idx,
			***REMOVED***,
			Identifier: ast.Identifier***REMOVED***
				Name: parsedLiteral,
				Idx:  idIdx,
			***REMOVED***,
		***REMOVED***
	case token.LEFT_BRACKET:
		return self.parseBracketMember(&ast.SuperExpression***REMOVED***
			Idx: idx,
		***REMOVED***)
	case token.LEFT_PARENTHESIS:
		return self.parseCallExpression(&ast.SuperExpression***REMOVED***
			Idx: idx,
		***REMOVED***)
	default:
		self.error(idx, "'super' keyword unexpected here")
		self.nextStatement()
		return &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
	***REMOVED***
***REMOVED***

func (self *_parser) reinterpretSequenceAsArrowFuncParams(seq *ast.SequenceExpression) *ast.ParameterList ***REMOVED***
	firstRestIdx := -1
	params := make([]*ast.Binding, 0, len(seq.Sequence))
	for i, item := range seq.Sequence ***REMOVED***
		if _, ok := item.(*ast.SpreadElement); ok ***REMOVED***
			if firstRestIdx == -1 ***REMOVED***
				firstRestIdx = i
				continue
			***REMOVED***
		***REMOVED***
		if firstRestIdx != -1 ***REMOVED***
			self.error(seq.Sequence[firstRestIdx].Idx0(), "Rest parameter must be last formal parameter")
			return &ast.ParameterList***REMOVED******REMOVED***
		***REMOVED***
		params = append(params, self.reinterpretAsBinding(item))
	***REMOVED***
	var rest ast.Expression
	if firstRestIdx != -1 ***REMOVED***
		rest = self.reinterpretAsBindingRestElement(seq.Sequence[firstRestIdx])
	***REMOVED***
	return &ast.ParameterList***REMOVED***
		List: params,
		Rest: rest,
	***REMOVED***
***REMOVED***

func (self *_parser) parseParenthesisedExpression() ast.Expression ***REMOVED***
	opening := self.idx
	self.expect(token.LEFT_PARENTHESIS)
	var list []ast.Expression
	if self.token != token.RIGHT_PARENTHESIS ***REMOVED***
		for ***REMOVED***
			if self.token == token.ELLIPSIS ***REMOVED***
				start := self.idx
				self.errorUnexpectedToken(token.ELLIPSIS)
				self.next()
				expr := self.parseAssignmentExpression()
				list = append(list, &ast.BadExpression***REMOVED***
					From: start,
					To:   expr.Idx1(),
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				list = append(list, self.parseAssignmentExpression())
			***REMOVED***
			if self.token != token.COMMA ***REMOVED***
				break
			***REMOVED***
			self.next()
			if self.token == token.RIGHT_PARENTHESIS ***REMOVED***
				self.errorUnexpectedToken(token.RIGHT_PARENTHESIS)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	self.expect(token.RIGHT_PARENTHESIS)
	if len(list) == 1 && len(self.errors) == 0 ***REMOVED***
		return list[0]
	***REMOVED***
	if len(list) == 0 ***REMOVED***
		self.errorUnexpectedToken(token.RIGHT_PARENTHESIS)
		return &ast.BadExpression***REMOVED***
			From: opening,
			To:   self.idx,
		***REMOVED***
	***REMOVED***
	return &ast.SequenceExpression***REMOVED***
		Sequence: list,
	***REMOVED***
***REMOVED***

func (self *_parser) parseRegExpLiteral() *ast.RegExpLiteral ***REMOVED***

	offset := self.chrOffset - 1 // Opening slash already gotten
	if self.token == token.QUOTIENT_ASSIGN ***REMOVED***
		offset -= 1 // =
	***REMOVED***
	idx := self.idxOf(offset)

	pattern, _, err := self.scanString(offset, false)
	endOffset := self.chrOffset

	if err == "" ***REMOVED***
		pattern = pattern[1 : len(pattern)-1]
	***REMOVED***

	flags := ""
	if !isLineTerminator(self.chr) && !isLineWhiteSpace(self.chr) ***REMOVED***
		self.next()

		if self.token == token.IDENTIFIER ***REMOVED*** // gim

			flags = self.literal
			self.next()
			endOffset = self.chrOffset - 1
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		self.next()
	***REMOVED***

	literal := self.str[offset:endOffset]

	return &ast.RegExpLiteral***REMOVED***
		Idx:     idx,
		Literal: literal,
		Pattern: pattern,
		Flags:   flags,
	***REMOVED***
***REMOVED***

func isBindingId(tok token.Token, parsedLiteral unistring.String) bool ***REMOVED***
	if tok == token.IDENTIFIER ***REMOVED***
		return true
	***REMOVED***
	if token.IsId(tok) ***REMOVED***
		switch parsedLiteral ***REMOVED***
		case "yield", "await":
			return true
		***REMOVED***
		if token.IsUnreservedWord(tok) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (self *_parser) tokenToBindingId() ***REMOVED***
	if isBindingId(self.token, self.parsedLiteral) ***REMOVED***
		self.token = token.IDENTIFIER
	***REMOVED***
***REMOVED***

func (self *_parser) parseBindingTarget() (target ast.BindingTarget) ***REMOVED***
	self.tokenToBindingId()
	switch self.token ***REMOVED***
	case token.IDENTIFIER:
		target = &ast.Identifier***REMOVED***
			Name: self.parsedLiteral,
			Idx:  self.idx,
		***REMOVED***
		self.next()
	case token.LEFT_BRACKET:
		target = self.parseArrayBindingPattern()
	case token.LEFT_BRACE:
		target = self.parseObjectBindingPattern()
	default:
		idx := self.expect(token.IDENTIFIER)
		self.nextStatement()
		target = &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
	***REMOVED***

	return
***REMOVED***

func (self *_parser) parseVariableDeclaration(declarationList *[]*ast.Binding) ast.Expression ***REMOVED***
	node := &ast.Binding***REMOVED***
		Target: self.parseBindingTarget(),
	***REMOVED***

	if declarationList != nil ***REMOVED***
		*declarationList = append(*declarationList, node)
	***REMOVED***

	if self.token == token.ASSIGN ***REMOVED***
		self.next()
		node.Initializer = self.parseAssignmentExpression()
	***REMOVED***

	return node
***REMOVED***

func (self *_parser) parseVariableDeclarationList() (declarationList []*ast.Binding) ***REMOVED***
	for ***REMOVED***
		self.parseVariableDeclaration(&declarationList)
		if self.token != token.COMMA ***REMOVED***
			break
		***REMOVED***
		self.next()
	***REMOVED***
	return
***REMOVED***

func (self *_parser) parseVarDeclarationList(var_ file.Idx) []*ast.Binding ***REMOVED***
	declarationList := self.parseVariableDeclarationList()

	self.scope.declare(&ast.VariableDeclaration***REMOVED***
		Var:  var_,
		List: declarationList,
	***REMOVED***)

	return declarationList
***REMOVED***

func (self *_parser) parseObjectPropertyKey() (string, unistring.String, ast.Expression, token.Token) ***REMOVED***
	if self.token == token.LEFT_BRACKET ***REMOVED***
		self.next()
		expr := self.parseAssignmentExpression()
		self.expect(token.RIGHT_BRACKET)
		return "", "", expr, token.ILLEGAL
	***REMOVED***
	idx, tkn, literal, parsedLiteral := self.idx, self.token, self.literal, self.parsedLiteral
	var value ast.Expression
	self.next()
	switch tkn ***REMOVED***
	case token.IDENTIFIER, token.STRING, token.KEYWORD, token.ESCAPED_RESERVED_WORD:
		value = &ast.StringLiteral***REMOVED***
			Idx:     idx,
			Literal: literal,
			Value:   parsedLiteral,
		***REMOVED***
	case token.NUMBER:
		num, err := parseNumberLiteral(literal)
		if err != nil ***REMOVED***
			self.error(idx, err.Error())
		***REMOVED*** else ***REMOVED***
			value = &ast.NumberLiteral***REMOVED***
				Idx:     idx,
				Literal: literal,
				Value:   num,
			***REMOVED***
		***REMOVED***
	case token.PRIVATE_IDENTIFIER:
		value = &ast.PrivateIdentifier***REMOVED***
			Identifier: ast.Identifier***REMOVED***
				Idx:  idx,
				Name: parsedLiteral,
			***REMOVED***,
		***REMOVED***
	default:
		// null, false, class, etc.
		if token.IsId(tkn) ***REMOVED***
			value = &ast.StringLiteral***REMOVED***
				Idx:     idx,
				Literal: literal,
				Value:   unistring.String(literal),
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			self.errorUnexpectedToken(tkn)
		***REMOVED***
	***REMOVED***
	return literal, parsedLiteral, value, tkn
***REMOVED***

func (self *_parser) parseObjectProperty() ast.Property ***REMOVED***
	if self.token == token.ELLIPSIS ***REMOVED***
		self.next()
		return &ast.SpreadElement***REMOVED***
			Expression: self.parseAssignmentExpression(),
		***REMOVED***
	***REMOVED***
	keyStartIdx := self.idx
	literal, parsedLiteral, value, tkn := self.parseObjectPropertyKey()
	if value == nil ***REMOVED***
		return nil
	***REMOVED***
	if token.IsId(tkn) || tkn == token.STRING || tkn == token.ILLEGAL ***REMOVED***
		switch ***REMOVED***
		case self.token == token.LEFT_PARENTHESIS:
			parameterList := self.parseFunctionParameterList()

			node := &ast.FunctionLiteral***REMOVED***
				Function:      keyStartIdx,
				ParameterList: parameterList,
			***REMOVED***
			node.Body, node.DeclarationList = self.parseFunctionBlock()
			node.Source = self.slice(keyStartIdx, node.Body.Idx1())

			return &ast.PropertyKeyed***REMOVED***
				Key:   value,
				Kind:  ast.PropertyKindMethod,
				Value: node,
			***REMOVED***
		case self.token == token.COMMA || self.token == token.RIGHT_BRACE || self.token == token.ASSIGN: // shorthand property
			if isBindingId(tkn, parsedLiteral) ***REMOVED***
				var initializer ast.Expression
				if self.token == token.ASSIGN ***REMOVED***
					// allow the initializer syntax here in case the object literal
					// needs to be reinterpreted as an assignment pattern, enforce later if it doesn't.
					self.next()
					initializer = self.parseAssignmentExpression()
				***REMOVED***
				return &ast.PropertyShort***REMOVED***
					Name: ast.Identifier***REMOVED***
						Name: parsedLiteral,
						Idx:  value.Idx0(),
					***REMOVED***,
					Initializer: initializer,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				self.errorUnexpectedToken(self.token)
			***REMOVED***
		case (literal == "get" || literal == "set") && self.token != token.COLON:
			_, _, keyValue, _ := self.parseObjectPropertyKey()
			if keyValue == nil ***REMOVED***
				return nil
			***REMOVED***

			var kind ast.PropertyKind
			if literal == "get" ***REMOVED***
				kind = ast.PropertyKindGet
			***REMOVED*** else ***REMOVED***
				kind = ast.PropertyKindSet
			***REMOVED***

			return &ast.PropertyKeyed***REMOVED***
				Key:   keyValue,
				Kind:  kind,
				Value: self.parseMethodDefinition(keyStartIdx, kind),
			***REMOVED***
		***REMOVED***
	***REMOVED***

	self.expect(token.COLON)
	return &ast.PropertyKeyed***REMOVED***
		Key:      value,
		Kind:     ast.PropertyKindValue,
		Value:    self.parseAssignmentExpression(),
		Computed: tkn == token.ILLEGAL,
	***REMOVED***
***REMOVED***

func (self *_parser) parseMethodDefinition(keyStartIdx file.Idx, kind ast.PropertyKind) *ast.FunctionLiteral ***REMOVED***
	idx1 := self.idx
	parameterList := self.parseFunctionParameterList()
	switch kind ***REMOVED***
	case ast.PropertyKindGet:
		if len(parameterList.List) > 0 || parameterList.Rest != nil ***REMOVED***
			self.error(idx1, "Getter must not have any formal parameters.")
		***REMOVED***
	case ast.PropertyKindSet:
		if len(parameterList.List) != 1 || parameterList.Rest != nil ***REMOVED***
			self.error(idx1, "Setter must have exactly one formal parameter.")
		***REMOVED***
	***REMOVED***
	node := &ast.FunctionLiteral***REMOVED***
		Function:      keyStartIdx,
		ParameterList: parameterList,
	***REMOVED***
	node.Body, node.DeclarationList = self.parseFunctionBlock()
	node.Source = self.slice(keyStartIdx, node.Body.Idx1())
	return node
***REMOVED***

func (self *_parser) parseObjectLiteral() *ast.ObjectLiteral ***REMOVED***
	var value []ast.Property
	idx0 := self.expect(token.LEFT_BRACE)
	for self.token != token.RIGHT_BRACE && self.token != token.EOF ***REMOVED***
		property := self.parseObjectProperty()
		if property != nil ***REMOVED***
			value = append(value, property)
		***REMOVED***
		if self.token != token.RIGHT_BRACE ***REMOVED***
			self.expect(token.COMMA)
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	idx1 := self.expect(token.RIGHT_BRACE)

	return &ast.ObjectLiteral***REMOVED***
		LeftBrace:  idx0,
		RightBrace: idx1,
		Value:      value,
	***REMOVED***
***REMOVED***

func (self *_parser) parseArrayLiteral() *ast.ArrayLiteral ***REMOVED***

	idx0 := self.expect(token.LEFT_BRACKET)
	var value []ast.Expression
	for self.token != token.RIGHT_BRACKET && self.token != token.EOF ***REMOVED***
		if self.token == token.COMMA ***REMOVED***
			self.next()
			value = append(value, nil)
			continue
		***REMOVED***
		if self.token == token.ELLIPSIS ***REMOVED***
			self.next()
			value = append(value, &ast.SpreadElement***REMOVED***
				Expression: self.parseAssignmentExpression(),
			***REMOVED***)
		***REMOVED*** else ***REMOVED***
			value = append(value, self.parseAssignmentExpression())
		***REMOVED***
		if self.token != token.RIGHT_BRACKET ***REMOVED***
			self.expect(token.COMMA)
		***REMOVED***
	***REMOVED***
	idx1 := self.expect(token.RIGHT_BRACKET)

	return &ast.ArrayLiteral***REMOVED***
		LeftBracket:  idx0,
		RightBracket: idx1,
		Value:        value,
	***REMOVED***
***REMOVED***

func (self *_parser) parseTemplateLiteral(tagged bool) *ast.TemplateLiteral ***REMOVED***
	res := &ast.TemplateLiteral***REMOVED***
		OpenQuote: self.idx,
	***REMOVED***
	for ***REMOVED***
		start := self.offset
		literal, parsed, finished, parseErr, err := self.parseTemplateCharacters()
		if err != "" ***REMOVED***
			self.error(self.offset, err)
		***REMOVED***
		res.Elements = append(res.Elements, &ast.TemplateElement***REMOVED***
			Idx:     self.idxOf(start),
			Literal: literal,
			Parsed:  parsed,
			Valid:   parseErr == "",
		***REMOVED***)
		if !tagged && parseErr != "" ***REMOVED***
			self.error(self.offset, parseErr)
		***REMOVED***
		end := self.chrOffset - 1
		self.next()
		if finished ***REMOVED***
			res.CloseQuote = self.idxOf(end)
			break
		***REMOVED***
		expr := self.parseExpression()
		res.Expressions = append(res.Expressions, expr)
		if self.token != token.RIGHT_BRACE ***REMOVED***
			self.errorUnexpectedToken(self.token)
		***REMOVED***
	***REMOVED***
	return res
***REMOVED***

func (self *_parser) parseTaggedTemplateLiteral(tag ast.Expression) *ast.TemplateLiteral ***REMOVED***
	l := self.parseTemplateLiteral(true)
	l.Tag = tag
	return l
***REMOVED***

func (self *_parser) parseArgumentList() (argumentList []ast.Expression, idx0, idx1 file.Idx) ***REMOVED***
	idx0 = self.expect(token.LEFT_PARENTHESIS)
	for self.token != token.RIGHT_PARENTHESIS ***REMOVED***
		var item ast.Expression
		if self.token == token.ELLIPSIS ***REMOVED***
			self.next()
			item = &ast.SpreadElement***REMOVED***
				Expression: self.parseAssignmentExpression(),
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			item = self.parseAssignmentExpression()
		***REMOVED***
		argumentList = append(argumentList, item)
		if self.token != token.COMMA ***REMOVED***
			break
		***REMOVED***
		self.next()
	***REMOVED***
	idx1 = self.expect(token.RIGHT_PARENTHESIS)
	return
***REMOVED***

func (self *_parser) parseCallExpression(left ast.Expression) ast.Expression ***REMOVED***
	argumentList, idx0, idx1 := self.parseArgumentList()
	return &ast.CallExpression***REMOVED***
		Callee:           left,
		LeftParenthesis:  idx0,
		ArgumentList:     argumentList,
		RightParenthesis: idx1,
	***REMOVED***
***REMOVED***

func (self *_parser) parseDotMember(left ast.Expression) ast.Expression ***REMOVED***
	period := self.idx
	self.next()

	literal := self.parsedLiteral
	idx := self.idx

	if self.token == token.PRIVATE_IDENTIFIER ***REMOVED***
		self.next()
		return &ast.PrivateDotExpression***REMOVED***
			Left: left,
			Identifier: ast.PrivateIdentifier***REMOVED***
				Identifier: ast.Identifier***REMOVED***
					Idx:  idx,
					Name: literal,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	if !token.IsId(self.token) ***REMOVED***
		self.expect(token.IDENTIFIER)
		self.nextStatement()
		return &ast.BadExpression***REMOVED***From: period, To: self.idx***REMOVED***
	***REMOVED***

	self.next()

	return &ast.DotExpression***REMOVED***
		Left: left,
		Identifier: ast.Identifier***REMOVED***
			Idx:  idx,
			Name: literal,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (self *_parser) parseBracketMember(left ast.Expression) ast.Expression ***REMOVED***
	idx0 := self.expect(token.LEFT_BRACKET)
	member := self.parseExpression()
	idx1 := self.expect(token.RIGHT_BRACKET)
	return &ast.BracketExpression***REMOVED***
		LeftBracket:  idx0,
		Left:         left,
		Member:       member,
		RightBracket: idx1,
	***REMOVED***
***REMOVED***

func (self *_parser) parseNewExpression() ast.Expression ***REMOVED***
	idx := self.expect(token.NEW)
	if self.token == token.PERIOD ***REMOVED***
		self.next()
		if self.literal == "target" ***REMOVED***
			return &ast.MetaProperty***REMOVED***
				Meta: &ast.Identifier***REMOVED***
					Name: unistring.String(token.NEW.String()),
					Idx:  idx,
				***REMOVED***,
				Property: self.parseIdentifier(),
			***REMOVED***
		***REMOVED***
		self.errorUnexpectedToken(token.IDENTIFIER)
	***REMOVED***
	callee := self.parseLeftHandSideExpression()
	if bad, ok := callee.(*ast.BadExpression); ok ***REMOVED***
		bad.From = idx
		return bad
	***REMOVED***
	node := &ast.NewExpression***REMOVED***
		New:    idx,
		Callee: callee,
	***REMOVED***
	if self.token == token.LEFT_PARENTHESIS ***REMOVED***
		argumentList, idx0, idx1 := self.parseArgumentList()
		node.ArgumentList = argumentList
		node.LeftParenthesis = idx0
		node.RightParenthesis = idx1
	***REMOVED***
	return node
***REMOVED***

func (self *_parser) parseLeftHandSideExpression() ast.Expression ***REMOVED***

	var left ast.Expression
	if self.token == token.NEW ***REMOVED***
		left = self.parseNewExpression()
	***REMOVED*** else ***REMOVED***
		left = self.parsePrimaryExpression()
	***REMOVED***
L:
	for ***REMOVED***
		switch self.token ***REMOVED***
		case token.PERIOD:
			left = self.parseDotMember(left)
		case token.LEFT_BRACKET:
			left = self.parseBracketMember(left)
		case token.BACKTICK:
			left = self.parseTaggedTemplateLiteral(left)
		default:
			break L
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseLeftHandSideExpressionAllowCall() ast.Expression ***REMOVED***

	allowIn := self.scope.allowIn
	self.scope.allowIn = true
	defer func() ***REMOVED***
		self.scope.allowIn = allowIn
	***REMOVED***()

	var left ast.Expression
	start := self.idx
	if self.token == token.NEW ***REMOVED***
		left = self.parseNewExpression()
	***REMOVED*** else ***REMOVED***
		left = self.parsePrimaryExpression()
	***REMOVED***

	optionalChain := false
L:
	for ***REMOVED***
		switch self.token ***REMOVED***
		case token.PERIOD:
			left = self.parseDotMember(left)
		case token.LEFT_BRACKET:
			left = self.parseBracketMember(left)
		case token.LEFT_PARENTHESIS:
			left = self.parseCallExpression(left)
		case token.BACKTICK:
			if optionalChain ***REMOVED***
				self.error(self.idx, "Invalid template literal on optional chain")
				self.nextStatement()
				return &ast.BadExpression***REMOVED***From: start, To: self.idx***REMOVED***
			***REMOVED***
			left = self.parseTaggedTemplateLiteral(left)
		case token.QUESTION_DOT:
			optionalChain = true
			left = &ast.Optional***REMOVED***Expression: left***REMOVED***

			switch self.peek() ***REMOVED***
			case token.LEFT_BRACKET, token.LEFT_PARENTHESIS, token.BACKTICK:
				self.next()
			default:
				left = self.parseDotMember(left)
			***REMOVED***
		default:
			break L
		***REMOVED***
	***REMOVED***

	if optionalChain ***REMOVED***
		left = &ast.OptionalChain***REMOVED***Expression: left***REMOVED***
	***REMOVED***
	return left
***REMOVED***

func (self *_parser) parsePostfixExpression() ast.Expression ***REMOVED***
	operand := self.parseLeftHandSideExpressionAllowCall()

	switch self.token ***REMOVED***
	case token.INCREMENT, token.DECREMENT:
		// Make sure there is no line terminator here
		if self.implicitSemicolon ***REMOVED***
			break
		***REMOVED***
		tkn := self.token
		idx := self.idx
		self.next()
		switch operand.(type) ***REMOVED***
		case *ast.Identifier, *ast.DotExpression, *ast.PrivateDotExpression, *ast.BracketExpression:
		default:
			self.error(idx, "Invalid left-hand side in assignment")
			self.nextStatement()
			return &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
		***REMOVED***
		return &ast.UnaryExpression***REMOVED***
			Operator: tkn,
			Idx:      idx,
			Operand:  operand,
			Postfix:  true,
		***REMOVED***
	***REMOVED***

	return operand
***REMOVED***

func (self *_parser) parseUnaryExpression() ast.Expression ***REMOVED***

	switch self.token ***REMOVED***
	case token.PLUS, token.MINUS, token.NOT, token.BITWISE_NOT:
		fallthrough
	case token.DELETE, token.VOID, token.TYPEOF:
		tkn := self.token
		idx := self.idx
		self.next()
		return &ast.UnaryExpression***REMOVED***
			Operator: tkn,
			Idx:      idx,
			Operand:  self.parseUnaryExpression(),
		***REMOVED***
	case token.INCREMENT, token.DECREMENT:
		tkn := self.token
		idx := self.idx
		self.next()
		operand := self.parseUnaryExpression()
		switch operand.(type) ***REMOVED***
		case *ast.Identifier, *ast.DotExpression, *ast.PrivateDotExpression, *ast.BracketExpression:
		default:
			self.error(idx, "Invalid left-hand side in assignment")
			self.nextStatement()
			return &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
		***REMOVED***
		return &ast.UnaryExpression***REMOVED***
			Operator: tkn,
			Idx:      idx,
			Operand:  operand,
		***REMOVED***
	***REMOVED***

	return self.parsePostfixExpression()
***REMOVED***

func isUpdateExpression(expr ast.Expression) bool ***REMOVED***
	if ux, ok := expr.(*ast.UnaryExpression); ok ***REMOVED***
		return ux.Operator == token.INCREMENT || ux.Operator == token.DECREMENT
	***REMOVED***
	return true
***REMOVED***

func (self *_parser) parseExponentiationExpression() ast.Expression ***REMOVED***
	left := self.parseUnaryExpression()

	for self.token == token.EXPONENT && isUpdateExpression(left) ***REMOVED***
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: token.EXPONENT,
			Left:     left,
			Right:    self.parseExponentiationExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseMultiplicativeExpression() ast.Expression ***REMOVED***
	left := self.parseExponentiationExpression()

	for self.token == token.MULTIPLY || self.token == token.SLASH ||
		self.token == token.REMAINDER ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    self.parseExponentiationExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseAdditiveExpression() ast.Expression ***REMOVED***
	left := self.parseMultiplicativeExpression()

	for self.token == token.PLUS || self.token == token.MINUS ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    self.parseMultiplicativeExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseShiftExpression() ast.Expression ***REMOVED***
	left := self.parseAdditiveExpression()

	for self.token == token.SHIFT_LEFT || self.token == token.SHIFT_RIGHT ||
		self.token == token.UNSIGNED_SHIFT_RIGHT ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    self.parseAdditiveExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseRelationalExpression() ast.Expression ***REMOVED***
	if self.scope.allowIn && self.token == token.PRIVATE_IDENTIFIER ***REMOVED***
		left := &ast.PrivateIdentifier***REMOVED***
			Identifier: ast.Identifier***REMOVED***
				Idx:  self.idx,
				Name: self.parsedLiteral,
			***REMOVED***,
		***REMOVED***
		self.next()
		if self.token == token.IN ***REMOVED***
			self.next()
			return &ast.BinaryExpression***REMOVED***
				Operator: self.token,
				Left:     left,
				Right:    self.parseShiftExpression(),
			***REMOVED***
		***REMOVED***
		return left
	***REMOVED***
	left := self.parseShiftExpression()

	allowIn := self.scope.allowIn
	self.scope.allowIn = true
	defer func() ***REMOVED***
		self.scope.allowIn = allowIn
	***REMOVED***()

	switch self.token ***REMOVED***
	case token.LESS, token.LESS_OR_EQUAL, token.GREATER, token.GREATER_OR_EQUAL:
		tkn := self.token
		self.next()
		return &ast.BinaryExpression***REMOVED***
			Operator:   tkn,
			Left:       left,
			Right:      self.parseRelationalExpression(),
			Comparison: true,
		***REMOVED***
	case token.INSTANCEOF:
		tkn := self.token
		self.next()
		return &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    self.parseRelationalExpression(),
		***REMOVED***
	case token.IN:
		if !allowIn ***REMOVED***
			return left
		***REMOVED***
		tkn := self.token
		self.next()
		return &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    self.parseRelationalExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseEqualityExpression() ast.Expression ***REMOVED***
	left := self.parseRelationalExpression()

	for self.token == token.EQUAL || self.token == token.NOT_EQUAL ||
		self.token == token.STRICT_EQUAL || self.token == token.STRICT_NOT_EQUAL ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator:   tkn,
			Left:       left,
			Right:      self.parseRelationalExpression(),
			Comparison: true,
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseBitwiseAndExpression() ast.Expression ***REMOVED***
	left := self.parseEqualityExpression()

	for self.token == token.AND ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    self.parseEqualityExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseBitwiseExclusiveOrExpression() ast.Expression ***REMOVED***
	left := self.parseBitwiseAndExpression()

	for self.token == token.EXCLUSIVE_OR ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    self.parseBitwiseAndExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseBitwiseOrExpression() ast.Expression ***REMOVED***
	left := self.parseBitwiseExclusiveOrExpression()

	for self.token == token.OR ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    self.parseBitwiseExclusiveOrExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseLogicalAndExpression() ast.Expression ***REMOVED***
	left := self.parseBitwiseOrExpression()

	for self.token == token.LOGICAL_AND ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    self.parseBitwiseOrExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func isLogicalAndExpr(expr ast.Expression) bool ***REMOVED***
	if bexp, ok := expr.(*ast.BinaryExpression); ok && bexp.Operator == token.LOGICAL_AND ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (self *_parser) parseLogicalOrExpression() ast.Expression ***REMOVED***
	var idx file.Idx
	parenthesis := self.token == token.LEFT_PARENTHESIS
	left := self.parseLogicalAndExpression()

	if self.token == token.LOGICAL_OR || !parenthesis && isLogicalAndExpr(left) ***REMOVED***
		for ***REMOVED***
			switch self.token ***REMOVED***
			case token.LOGICAL_OR:
				self.next()
				left = &ast.BinaryExpression***REMOVED***
					Operator: token.LOGICAL_OR,
					Left:     left,
					Right:    self.parseLogicalAndExpression(),
				***REMOVED***
			case token.COALESCE:
				idx = self.idx
				goto mixed
			default:
				return left
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for ***REMOVED***
			switch self.token ***REMOVED***
			case token.COALESCE:
				idx = self.idx
				self.next()

				parenthesis := self.token == token.LEFT_PARENTHESIS
				right := self.parseLogicalAndExpression()
				if !parenthesis && isLogicalAndExpr(right) ***REMOVED***
					goto mixed
				***REMOVED***

				left = &ast.BinaryExpression***REMOVED***
					Operator: token.COALESCE,
					Left:     left,
					Right:    right,
				***REMOVED***
			case token.LOGICAL_OR:
				idx = self.idx
				goto mixed
			default:
				return left
			***REMOVED***
		***REMOVED***
	***REMOVED***

mixed:
	self.error(idx, "Logical expressions and coalesce expressions cannot be mixed. Wrap either by parentheses")
	return left
***REMOVED***

func (self *_parser) parseConditionalExpression() ast.Expression ***REMOVED***
	left := self.parseLogicalOrExpression()

	if self.token == token.QUESTION_MARK ***REMOVED***
		self.next()
		consequent := self.parseAssignmentExpression()
		self.expect(token.COLON)
		return &ast.ConditionalExpression***REMOVED***
			Test:       left,
			Consequent: consequent,
			Alternate:  self.parseAssignmentExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseAssignmentExpression() ast.Expression ***REMOVED***
	start := self.idx
	parenthesis := false
	var state parserState
	if self.token == token.LEFT_PARENTHESIS ***REMOVED***
		self.mark(&state)
		parenthesis = true
	***REMOVED*** else ***REMOVED***
		self.tokenToBindingId()
	***REMOVED***
	left := self.parseConditionalExpression()
	var operator token.Token
	switch self.token ***REMOVED***
	case token.ASSIGN:
		operator = self.token
	case token.ADD_ASSIGN:
		operator = token.PLUS
	case token.SUBTRACT_ASSIGN:
		operator = token.MINUS
	case token.MULTIPLY_ASSIGN:
		operator = token.MULTIPLY
	case token.EXPONENT_ASSIGN:
		operator = token.EXPONENT
	case token.QUOTIENT_ASSIGN:
		operator = token.SLASH
	case token.REMAINDER_ASSIGN:
		operator = token.REMAINDER
	case token.AND_ASSIGN:
		operator = token.AND
	case token.OR_ASSIGN:
		operator = token.OR
	case token.EXCLUSIVE_OR_ASSIGN:
		operator = token.EXCLUSIVE_OR
	case token.SHIFT_LEFT_ASSIGN:
		operator = token.SHIFT_LEFT
	case token.SHIFT_RIGHT_ASSIGN:
		operator = token.SHIFT_RIGHT
	case token.UNSIGNED_SHIFT_RIGHT_ASSIGN:
		operator = token.UNSIGNED_SHIFT_RIGHT
	case token.ARROW:
		var paramList *ast.ParameterList
		if id, ok := left.(*ast.Identifier); ok ***REMOVED***
			paramList = &ast.ParameterList***REMOVED***
				Opening: id.Idx,
				Closing: id.Idx1(),
				List: []*ast.Binding***REMOVED******REMOVED***
					Target: id,
				***REMOVED******REMOVED***,
			***REMOVED***
		***REMOVED*** else if parenthesis ***REMOVED***
			if seq, ok := left.(*ast.SequenceExpression); ok && len(self.errors) == 0 ***REMOVED***
				paramList = self.reinterpretSequenceAsArrowFuncParams(seq)
			***REMOVED*** else ***REMOVED***
				self.restore(&state)
				paramList = self.parseFunctionParameterList()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			self.error(left.Idx0(), "Malformed arrow function parameter list")
			return &ast.BadExpression***REMOVED***From: left.Idx0(), To: left.Idx1()***REMOVED***
		***REMOVED***
		self.expect(token.ARROW)
		node := &ast.ArrowFunctionLiteral***REMOVED***
			Start:         start,
			ParameterList: paramList,
		***REMOVED***
		node.Body, node.DeclarationList = self.parseArrowFunctionBody()
		node.Source = self.slice(node.Start, node.Body.Idx1())
		return node
	***REMOVED***

	if operator != 0 ***REMOVED***
		idx := self.idx
		self.next()
		ok := false
		switch l := left.(type) ***REMOVED***
		case *ast.Identifier, *ast.DotExpression, *ast.PrivateDotExpression, *ast.BracketExpression:
			ok = true
		case *ast.ArrayLiteral:
			if !parenthesis && operator == token.ASSIGN ***REMOVED***
				left = self.reinterpretAsArrayAssignmentPattern(l)
				ok = true
			***REMOVED***
		case *ast.ObjectLiteral:
			if !parenthesis && operator == token.ASSIGN ***REMOVED***
				left = self.reinterpretAsObjectAssignmentPattern(l)
				ok = true
			***REMOVED***
		***REMOVED***
		if ok ***REMOVED***
			return &ast.AssignExpression***REMOVED***
				Left:     left,
				Operator: operator,
				Right:    self.parseAssignmentExpression(),
			***REMOVED***
		***REMOVED***
		self.error(left.Idx0(), "Invalid left-hand side in assignment")
		self.nextStatement()
		return &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseExpression() ast.Expression ***REMOVED***
	left := self.parseAssignmentExpression()

	if self.token == token.COMMA ***REMOVED***
		sequence := []ast.Expression***REMOVED***left***REMOVED***
		for ***REMOVED***
			if self.token != token.COMMA ***REMOVED***
				break
			***REMOVED***
			self.next()
			sequence = append(sequence, self.parseAssignmentExpression())
		***REMOVED***
		return &ast.SequenceExpression***REMOVED***
			Sequence: sequence,
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) checkComma(from, to file.Idx) ***REMOVED***
	if pos := strings.IndexByte(self.str[int(from)-self.base:int(to)-self.base], ','); pos >= 0 ***REMOVED***
		self.error(from+file.Idx(pos), "Comma is not allowed here")
	***REMOVED***
***REMOVED***

func (self *_parser) reinterpretAsArrayAssignmentPattern(left *ast.ArrayLiteral) ast.Expression ***REMOVED***
	value := left.Value
	var rest ast.Expression
	for i, item := range value ***REMOVED***
		if spread, ok := item.(*ast.SpreadElement); ok ***REMOVED***
			if i != len(value)-1 ***REMOVED***
				self.error(item.Idx0(), "Rest element must be last element")
				return &ast.BadExpression***REMOVED***From: left.Idx0(), To: left.Idx1()***REMOVED***
			***REMOVED***
			self.checkComma(spread.Expression.Idx1(), left.RightBracket)
			rest = self.reinterpretAsDestructAssignTarget(spread.Expression)
			value = value[:len(value)-1]
		***REMOVED*** else ***REMOVED***
			value[i] = self.reinterpretAsAssignmentElement(item)
		***REMOVED***
	***REMOVED***
	return &ast.ArrayPattern***REMOVED***
		LeftBracket:  left.LeftBracket,
		RightBracket: left.RightBracket,
		Elements:     value,
		Rest:         rest,
	***REMOVED***
***REMOVED***

func (self *_parser) reinterpretArrayAssignPatternAsBinding(pattern *ast.ArrayPattern) *ast.ArrayPattern ***REMOVED***
	for i, item := range pattern.Elements ***REMOVED***
		pattern.Elements[i] = self.reinterpretAsDestructBindingTarget(item)
	***REMOVED***
	if pattern.Rest != nil ***REMOVED***
		pattern.Rest = self.reinterpretAsDestructBindingTarget(pattern.Rest)
	***REMOVED***
	return pattern
***REMOVED***

func (self *_parser) reinterpretAsArrayBindingPattern(left *ast.ArrayLiteral) ast.BindingTarget ***REMOVED***
	value := left.Value
	var rest ast.Expression
	for i, item := range value ***REMOVED***
		if spread, ok := item.(*ast.SpreadElement); ok ***REMOVED***
			if i != len(value)-1 ***REMOVED***
				self.error(item.Idx0(), "Rest element must be last element")
				return &ast.BadExpression***REMOVED***From: left.Idx0(), To: left.Idx1()***REMOVED***
			***REMOVED***
			self.checkComma(spread.Expression.Idx1(), left.RightBracket)
			rest = self.reinterpretAsDestructBindingTarget(spread.Expression)
			value = value[:len(value)-1]
		***REMOVED*** else ***REMOVED***
			value[i] = self.reinterpretAsBindingElement(item)
		***REMOVED***
	***REMOVED***
	return &ast.ArrayPattern***REMOVED***
		LeftBracket:  left.LeftBracket,
		RightBracket: left.RightBracket,
		Elements:     value,
		Rest:         rest,
	***REMOVED***
***REMOVED***

func (self *_parser) parseArrayBindingPattern() ast.BindingTarget ***REMOVED***
	return self.reinterpretAsArrayBindingPattern(self.parseArrayLiteral())
***REMOVED***

func (self *_parser) parseObjectBindingPattern() ast.BindingTarget ***REMOVED***
	return self.reinterpretAsObjectBindingPattern(self.parseObjectLiteral())
***REMOVED***

func (self *_parser) reinterpretArrayObjectPatternAsBinding(pattern *ast.ObjectPattern) *ast.ObjectPattern ***REMOVED***
	for _, prop := range pattern.Properties ***REMOVED***
		if keyed, ok := prop.(*ast.PropertyKeyed); ok ***REMOVED***
			keyed.Value = self.reinterpretAsBindingElement(keyed.Value)
		***REMOVED***
	***REMOVED***
	if pattern.Rest != nil ***REMOVED***
		pattern.Rest = self.reinterpretAsBindingRestElement(pattern.Rest)
	***REMOVED***
	return pattern
***REMOVED***

func (self *_parser) reinterpretAsObjectBindingPattern(expr *ast.ObjectLiteral) ast.BindingTarget ***REMOVED***
	var rest ast.Expression
	value := expr.Value
	for i, prop := range value ***REMOVED***
		ok := false
		switch prop := prop.(type) ***REMOVED***
		case *ast.PropertyKeyed:
			if prop.Kind == ast.PropertyKindValue ***REMOVED***
				prop.Value = self.reinterpretAsBindingElement(prop.Value)
				ok = true
			***REMOVED***
		case *ast.PropertyShort:
			ok = true
		case *ast.SpreadElement:
			if i != len(expr.Value)-1 ***REMOVED***
				self.error(prop.Idx0(), "Rest element must be last element")
				return &ast.BadExpression***REMOVED***From: expr.Idx0(), To: expr.Idx1()***REMOVED***
			***REMOVED***
			// TODO make sure there is no trailing comma
			rest = self.reinterpretAsBindingRestElement(prop.Expression)
			value = value[:i]
			ok = true
		***REMOVED***
		if !ok ***REMOVED***
			self.error(prop.Idx0(), "Invalid destructuring binding target")
			return &ast.BadExpression***REMOVED***From: expr.Idx0(), To: expr.Idx1()***REMOVED***
		***REMOVED***
	***REMOVED***
	return &ast.ObjectPattern***REMOVED***
		LeftBrace:  expr.LeftBrace,
		RightBrace: expr.RightBrace,
		Properties: value,
		Rest:       rest,
	***REMOVED***
***REMOVED***

func (self *_parser) reinterpretAsObjectAssignmentPattern(l *ast.ObjectLiteral) ast.Expression ***REMOVED***
	var rest ast.Expression
	value := l.Value
	for i, prop := range value ***REMOVED***
		ok := false
		switch prop := prop.(type) ***REMOVED***
		case *ast.PropertyKeyed:
			if prop.Kind == ast.PropertyKindValue ***REMOVED***
				prop.Value = self.reinterpretAsAssignmentElement(prop.Value)
				ok = true
			***REMOVED***
		case *ast.PropertyShort:
			ok = true
		case *ast.SpreadElement:
			if i != len(l.Value)-1 ***REMOVED***
				self.error(prop.Idx0(), "Rest element must be last element")
				return &ast.BadExpression***REMOVED***From: l.Idx0(), To: l.Idx1()***REMOVED***
			***REMOVED***
			// TODO make sure there is no trailing comma
			rest = prop.Expression
			value = value[:i]
			ok = true
		***REMOVED***
		if !ok ***REMOVED***
			self.error(prop.Idx0(), "Invalid destructuring assignment target")
			return &ast.BadExpression***REMOVED***From: l.Idx0(), To: l.Idx1()***REMOVED***
		***REMOVED***
	***REMOVED***
	return &ast.ObjectPattern***REMOVED***
		LeftBrace:  l.LeftBrace,
		RightBrace: l.RightBrace,
		Properties: value,
		Rest:       rest,
	***REMOVED***
***REMOVED***

func (self *_parser) reinterpretAsAssignmentElement(expr ast.Expression) ast.Expression ***REMOVED***
	switch expr := expr.(type) ***REMOVED***
	case *ast.AssignExpression:
		if expr.Operator == token.ASSIGN ***REMOVED***
			expr.Left = self.reinterpretAsDestructAssignTarget(expr.Left)
			return expr
		***REMOVED*** else ***REMOVED***
			self.error(expr.Idx0(), "Invalid destructuring assignment target")
			return &ast.BadExpression***REMOVED***From: expr.Idx0(), To: expr.Idx1()***REMOVED***
		***REMOVED***
	default:
		return self.reinterpretAsDestructAssignTarget(expr)
	***REMOVED***
***REMOVED***

func (self *_parser) reinterpretAsBindingElement(expr ast.Expression) ast.Expression ***REMOVED***
	switch expr := expr.(type) ***REMOVED***
	case *ast.AssignExpression:
		if expr.Operator == token.ASSIGN ***REMOVED***
			expr.Left = self.reinterpretAsDestructBindingTarget(expr.Left)
			return expr
		***REMOVED*** else ***REMOVED***
			self.error(expr.Idx0(), "Invalid destructuring assignment target")
			return &ast.BadExpression***REMOVED***From: expr.Idx0(), To: expr.Idx1()***REMOVED***
		***REMOVED***
	default:
		return self.reinterpretAsDestructBindingTarget(expr)
	***REMOVED***
***REMOVED***

func (self *_parser) reinterpretAsBinding(expr ast.Expression) *ast.Binding ***REMOVED***
	switch expr := expr.(type) ***REMOVED***
	case *ast.AssignExpression:
		if expr.Operator == token.ASSIGN ***REMOVED***
			return &ast.Binding***REMOVED***
				Target:      self.reinterpretAsDestructBindingTarget(expr.Left),
				Initializer: expr.Right,
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			self.error(expr.Idx0(), "Invalid destructuring assignment target")
			return &ast.Binding***REMOVED***
				Target: &ast.BadExpression***REMOVED***From: expr.Idx0(), To: expr.Idx1()***REMOVED***,
			***REMOVED***
		***REMOVED***
	default:
		return &ast.Binding***REMOVED***
			Target: self.reinterpretAsDestructBindingTarget(expr),
		***REMOVED***
	***REMOVED***
***REMOVED***

func (self *_parser) reinterpretAsDestructAssignTarget(item ast.Expression) ast.Expression ***REMOVED***
	switch item := item.(type) ***REMOVED***
	case nil:
		return nil
	case *ast.ArrayLiteral:
		return self.reinterpretAsArrayAssignmentPattern(item)
	case *ast.ObjectLiteral:
		return self.reinterpretAsObjectAssignmentPattern(item)
	case ast.Pattern, *ast.Identifier, *ast.DotExpression, *ast.PrivateDotExpression, *ast.BracketExpression:
		return item
	***REMOVED***
	self.error(item.Idx0(), "Invalid destructuring assignment target")
	return &ast.BadExpression***REMOVED***From: item.Idx0(), To: item.Idx1()***REMOVED***
***REMOVED***

func (self *_parser) reinterpretAsDestructBindingTarget(item ast.Expression) ast.BindingTarget ***REMOVED***
	switch item := item.(type) ***REMOVED***
	case nil:
		return nil
	case *ast.ArrayPattern:
		return self.reinterpretArrayAssignPatternAsBinding(item)
	case *ast.ObjectPattern:
		return self.reinterpretArrayObjectPatternAsBinding(item)
	case *ast.ArrayLiteral:
		return self.reinterpretAsArrayBindingPattern(item)
	case *ast.ObjectLiteral:
		return self.reinterpretAsObjectBindingPattern(item)
	case *ast.Identifier:
		return item
	***REMOVED***
	self.error(item.Idx0(), "Invalid destructuring binding target")
	return &ast.BadExpression***REMOVED***From: item.Idx0(), To: item.Idx1()***REMOVED***
***REMOVED***

func (self *_parser) reinterpretAsBindingRestElement(expr ast.Expression) ast.Expression ***REMOVED***
	if _, ok := expr.(*ast.Identifier); ok ***REMOVED***
		return expr
	***REMOVED***
	self.error(expr.Idx0(), "Invalid binding rest")
	return &ast.BadExpression***REMOVED***From: expr.Idx0(), To: expr.Idx1()***REMOVED***
***REMOVED***
