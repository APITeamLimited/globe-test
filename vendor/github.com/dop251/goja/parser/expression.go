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
		if len(literal) > 1 ***REMOVED***
			tkn, strict := token.IsKeyword(literal)
			if tkn == token.KEYWORD ***REMOVED***
				if !strict ***REMOVED***
					self.error(idx, "Unexpected reserved word")
				***REMOVED***
			***REMOVED***
		***REMOVED***
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
	case token.FUNCTION:
		return self.parseFunction(false)
	***REMOVED***

	self.errorUnexpectedToken(self.token)
	self.nextStatement()
	return &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
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

	if err == nil ***REMOVED***
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

func (self *_parser) parseBindingTarget() (target ast.BindingTarget) ***REMOVED***
	if self.token == token.LET ***REMOVED***
		self.token = token.IDENTIFIER
	***REMOVED***
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

func (self *_parser) parseObjectPropertyKey() (unistring.String, ast.Expression, token.Token) ***REMOVED***
	if self.token == token.LEFT_BRACKET ***REMOVED***
		self.next()
		expr := self.parseAssignmentExpression()
		self.expect(token.RIGHT_BRACKET)
		return "", expr, token.ILLEGAL
	***REMOVED***
	idx, tkn, literal, parsedLiteral := self.idx, self.token, self.literal, self.parsedLiteral
	var value ast.Expression
	self.next()
	switch tkn ***REMOVED***
	case token.IDENTIFIER:
		value = &ast.StringLiteral***REMOVED***
			Idx:     idx,
			Literal: literal,
			Value:   unistring.String(literal),
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
	case token.STRING:
		value = &ast.StringLiteral***REMOVED***
			Idx:     idx,
			Literal: literal,
			Value:   parsedLiteral,
		***REMOVED***
	default:
		// null, false, class, etc.
		if isId(tkn) ***REMOVED***
			value = &ast.StringLiteral***REMOVED***
				Idx:     idx,
				Literal: literal,
				Value:   unistring.String(literal),
			***REMOVED***
			tkn = token.KEYWORD
		***REMOVED***
	***REMOVED***
	return parsedLiteral, value, tkn
***REMOVED***

func (self *_parser) parseObjectProperty() ast.Property ***REMOVED***
	if self.token == token.ELLIPSIS ***REMOVED***
		self.next()
		return &ast.SpreadElement***REMOVED***
			Expression: self.parseAssignmentExpression(),
		***REMOVED***
	***REMOVED***
	literal, value, tkn := self.parseObjectPropertyKey()
	if tkn == token.IDENTIFIER || tkn == token.STRING || tkn == token.KEYWORD || tkn == token.ILLEGAL ***REMOVED***
		switch ***REMOVED***
		case self.token == token.LEFT_PARENTHESIS:
			idx := self.idx
			parameterList := self.parseFunctionParameterList()

			node := &ast.FunctionLiteral***REMOVED***
				Function:      idx,
				ParameterList: parameterList,
			***REMOVED***
			node.Body, node.DeclarationList = self.parseFunctionBlock()

			return &ast.PropertyKeyed***REMOVED***
				Key:   value,
				Kind:  ast.PropertyKindMethod,
				Value: node,
			***REMOVED***
		case self.token == token.COMMA || self.token == token.RIGHT_BRACE || self.token == token.ASSIGN: // shorthand property
			if tkn == token.IDENTIFIER || tkn == token.KEYWORD && literal == "let" ***REMOVED***
				var initializer ast.Expression
				if self.token == token.ASSIGN ***REMOVED***
					// allow the initializer syntax here in case the object literal
					// needs to be reinterpreted as an assignment pattern, enforce later if it doesn't.
					self.next()
					initializer = self.parseAssignmentExpression()
				***REMOVED***
				return &ast.PropertyShort***REMOVED***
					Name: ast.Identifier***REMOVED***
						Name: literal,
						Idx:  value.Idx0(),
					***REMOVED***,
					Initializer: initializer,
				***REMOVED***
			***REMOVED***
		case literal == "get" && self.token != token.COLON:
			idx := self.idx
			_, value, _ := self.parseObjectPropertyKey()
			idx1 := self.idx
			parameterList := self.parseFunctionParameterList()
			if len(parameterList.List) > 0 || parameterList.Rest != nil ***REMOVED***
				self.error(idx1, "Getter must not have any formal parameters.")
			***REMOVED***
			node := &ast.FunctionLiteral***REMOVED***
				Function:      idx,
				ParameterList: parameterList,
			***REMOVED***
			node.Body, node.DeclarationList = self.parseFunctionBlock()
			return &ast.PropertyKeyed***REMOVED***
				Key:   value,
				Kind:  ast.PropertyKindGet,
				Value: node,
			***REMOVED***
		case literal == "set" && self.token != token.COLON:
			idx := self.idx
			_, value, _ := self.parseObjectPropertyKey()
			parameterList := self.parseFunctionParameterList()

			node := &ast.FunctionLiteral***REMOVED***
				Function:      idx,
				ParameterList: parameterList,
			***REMOVED***

			node.Body, node.DeclarationList = self.parseFunctionBlock()

			return &ast.PropertyKeyed***REMOVED***
				Key:   value,
				Kind:  ast.PropertyKindSet,
				Value: node,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	self.expect(token.COLON)

	return &ast.PropertyKeyed***REMOVED***
		Key:   value,
		Kind:  ast.PropertyKindValue,
		Value: self.parseAssignmentExpression(),
	***REMOVED***
***REMOVED***

func (self *_parser) parseObjectLiteral() *ast.ObjectLiteral ***REMOVED***
	var value []ast.Property
	idx0 := self.expect(token.LEFT_BRACE)
	for self.token != token.RIGHT_BRACE && self.token != token.EOF ***REMOVED***
		property := self.parseObjectProperty()
		value = append(value, property)
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
	for self.chr != -1 ***REMOVED***
		start := self.idx + 1
		literal, parsed, finished, parseErr, err := self.parseTemplateCharacters()
		if err != nil ***REMOVED***
			self.error(self.idx, err.Error())
		***REMOVED***
		res.Elements = append(res.Elements, &ast.TemplateElement***REMOVED***
			Idx:     start,
			Literal: literal,
			Parsed:  parsed,
			Valid:   parseErr == nil,
		***REMOVED***)
		if !tagged && parseErr != nil ***REMOVED***
			self.error(self.idx, parseErr.Error())
		***REMOVED***
		end := self.idx + 1
		self.next()
		if finished ***REMOVED***
			res.CloseQuote = end
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
	if self.token != token.RIGHT_PARENTHESIS ***REMOVED***
		for ***REMOVED***
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
	period := self.expect(token.PERIOD)

	literal := self.parsedLiteral
	idx := self.idx

	if self.token != token.IDENTIFIER && !isId(self.token) ***REMOVED***
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
		prop := self.parseIdentifier()
		if prop.Name == "target" ***REMOVED***
			if !self.scope.inFunction ***REMOVED***
				self.error(idx, "new.target expression is not allowed here")
			***REMOVED***
			return &ast.MetaProperty***REMOVED***
				Meta: &ast.Identifier***REMOVED***
					Name: unistring.String(token.NEW.String()),
					Idx:  idx,
				***REMOVED***,
				Property: prop,
			***REMOVED***
		***REMOVED***
		self.errorUnexpectedToken(token.IDENTIFIER)
	***REMOVED***
	callee := self.parseLeftHandSideExpression()
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
		case token.LEFT_PARENTHESIS:
			left = self.parseCallExpression(left)
		case token.BACKTICK:
			left = self.parseTaggedTemplateLiteral(left)
		default:
			break L
		***REMOVED***
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
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
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
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
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

func (self *_parser) parseMultiplicativeExpression() ast.Expression ***REMOVED***
	next := self.parseUnaryExpression
	left := next()

	for self.token == token.MULTIPLY || self.token == token.SLASH ||
		self.token == token.REMAINDER ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    next(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseAdditiveExpression() ast.Expression ***REMOVED***
	next := self.parseMultiplicativeExpression
	left := next()

	for self.token == token.PLUS || self.token == token.MINUS ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    next(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseShiftExpression() ast.Expression ***REMOVED***
	next := self.parseAdditiveExpression
	left := next()

	for self.token == token.SHIFT_LEFT || self.token == token.SHIFT_RIGHT ||
		self.token == token.UNSIGNED_SHIFT_RIGHT ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    next(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseRelationalExpression() ast.Expression ***REMOVED***
	next := self.parseShiftExpression
	left := next()

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
	next := self.parseRelationalExpression
	left := next()

	for self.token == token.EQUAL || self.token == token.NOT_EQUAL ||
		self.token == token.STRICT_EQUAL || self.token == token.STRICT_NOT_EQUAL ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator:   tkn,
			Left:       left,
			Right:      next(),
			Comparison: true,
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseBitwiseAndExpression() ast.Expression ***REMOVED***
	next := self.parseEqualityExpression
	left := next()

	for self.token == token.AND ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    next(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseBitwiseExclusiveOrExpression() ast.Expression ***REMOVED***
	next := self.parseBitwiseAndExpression
	left := next()

	for self.token == token.EXCLUSIVE_OR ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    next(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseBitwiseOrExpression() ast.Expression ***REMOVED***
	next := self.parseBitwiseExclusiveOrExpression
	left := next()

	for self.token == token.OR ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    next(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseLogicalAndExpression() ast.Expression ***REMOVED***
	next := self.parseBitwiseOrExpression
	left := next()

	for self.token == token.LOGICAL_AND ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    next(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseLogicalOrExpression() ast.Expression ***REMOVED***
	next := self.parseLogicalAndExpression
	left := next()

	for self.token == token.LOGICAL_OR ***REMOVED***
		tkn := self.token
		self.next()
		left = &ast.BinaryExpression***REMOVED***
			Operator: tkn,
			Left:     left,
			Right:    next(),
		***REMOVED***
	***REMOVED***

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
	if self.token == token.LET ***REMOVED***
		self.token = token.IDENTIFIER
	***REMOVED*** else if self.token == token.LEFT_PARENTHESIS ***REMOVED***
		self.mark(&state)
		parenthesis = true
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
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
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
	if self.token == token.LET ***REMOVED***
		self.token = token.IDENTIFIER
	***REMOVED***
	next := self.parseAssignmentExpression
	left := next()

	if self.token == token.COMMA ***REMOVED***
		sequence := []ast.Expression***REMOVED***left***REMOVED***
		for ***REMOVED***
			if self.token != token.COMMA ***REMOVED***
				break
			***REMOVED***
			self.next()
			sequence = append(sequence, next())
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
	case ast.Pattern, *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
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
