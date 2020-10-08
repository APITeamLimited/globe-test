package parser

import (
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
		self.expect(token.LEFT_PARENTHESIS)
		expression := self.parseExpression()
		self.expect(token.RIGHT_PARENTHESIS)
		return expression
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

func (self *_parser) parseVariableDeclaration(declarationList *[]*ast.VariableExpression) ast.Expression ***REMOVED***

	if self.token != token.IDENTIFIER ***REMOVED***
		idx := self.expect(token.IDENTIFIER)
		self.nextStatement()
		return &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
	***REMOVED***

	name := self.parsedLiteral
	idx := self.idx
	self.next()
	node := &ast.VariableExpression***REMOVED***
		Name: name,
		Idx:  idx,
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

func (self *_parser) parseVariableDeclarationList(var_ file.Idx) []ast.Expression ***REMOVED***

	var declarationList []*ast.VariableExpression // Avoid bad expressions
	var list []ast.Expression

	for ***REMOVED***
		list = append(list, self.parseVariableDeclaration(&declarationList))
		if self.token != token.COMMA ***REMOVED***
			break
		***REMOVED***
		self.next()
	***REMOVED***

	self.scope.declare(&ast.VariableDeclaration***REMOVED***
		Var:  var_,
		List: declarationList,
	***REMOVED***)

	return list
***REMOVED***

func (self *_parser) parseObjectPropertyKey() (string, unistring.String) ***REMOVED***
	idx, tkn, literal, parsedLiteral := self.idx, self.token, self.literal, self.parsedLiteral
	var value unistring.String
	self.next()
	switch tkn ***REMOVED***
	case token.IDENTIFIER:
		value = parsedLiteral
	case token.NUMBER:
		var err error
		_, err = parseNumberLiteral(literal)
		if err != nil ***REMOVED***
			self.error(idx, err.Error())
		***REMOVED*** else ***REMOVED***
			value = unistring.String(literal)
		***REMOVED***
	case token.STRING:
		value = parsedLiteral
	default:
		// null, false, class, etc.
		if isId(tkn) ***REMOVED***
			value = unistring.String(literal)
		***REMOVED***
	***REMOVED***
	return literal, value
***REMOVED***

func (self *_parser) parseObjectProperty() ast.Property ***REMOVED***

	literal, value := self.parseObjectPropertyKey()
	if literal == "get" && self.token != token.COLON ***REMOVED***
		idx := self.idx
		_, value := self.parseObjectPropertyKey()
		parameterList := self.parseFunctionParameterList()

		node := &ast.FunctionLiteral***REMOVED***
			Function:      idx,
			ParameterList: parameterList,
		***REMOVED***
		self.parseFunctionBlock(node)
		return ast.Property***REMOVED***
			Key:   value,
			Kind:  "get",
			Value: node,
		***REMOVED***
	***REMOVED*** else if literal == "set" && self.token != token.COLON ***REMOVED***
		idx := self.idx
		_, value := self.parseObjectPropertyKey()
		parameterList := self.parseFunctionParameterList()

		node := &ast.FunctionLiteral***REMOVED***
			Function:      idx,
			ParameterList: parameterList,
		***REMOVED***
		self.parseFunctionBlock(node)
		return ast.Property***REMOVED***
			Key:   value,
			Kind:  "set",
			Value: node,
		***REMOVED***
	***REMOVED***

	self.expect(token.COLON)

	return ast.Property***REMOVED***
		Key:   value,
		Kind:  "value",
		Value: self.parseAssignmentExpression(),
	***REMOVED***
***REMOVED***

func (self *_parser) parseObjectLiteral() ast.Expression ***REMOVED***
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

func (self *_parser) parseArrayLiteral() ast.Expression ***REMOVED***

	idx0 := self.expect(token.LEFT_BRACKET)
	var value []ast.Expression
	for self.token != token.RIGHT_BRACKET && self.token != token.EOF ***REMOVED***
		if self.token == token.COMMA ***REMOVED***
			self.next()
			value = append(value, nil)
			continue
		***REMOVED***
		value = append(value, self.parseAssignmentExpression())
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

func (self *_parser) parseArgumentList() (argumentList []ast.Expression, idx0, idx1 file.Idx) ***REMOVED***
	idx0 = self.expect(token.LEFT_PARENTHESIS)
	if self.token != token.RIGHT_PARENTHESIS ***REMOVED***
		for ***REMOVED***
			argumentList = append(argumentList, self.parseAssignmentExpression())
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

	for ***REMOVED***
		if self.token == token.PERIOD ***REMOVED***
			left = self.parseDotMember(left)
		***REMOVED*** else if self.token == token.LEFT_BRACKET ***REMOVED***
			left = self.parseBracketMember(left)
		***REMOVED*** else ***REMOVED***
			break
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

	for ***REMOVED***
		if self.token == token.PERIOD ***REMOVED***
			left = self.parseDotMember(left)
		***REMOVED*** else if self.token == token.LEFT_BRACKET ***REMOVED***
			left = self.parseBracketMember(left)
		***REMOVED*** else if self.token == token.LEFT_PARENTHESIS ***REMOVED***
			left = self.parseCallExpression(left)
		***REMOVED*** else ***REMOVED***
			break
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

func (self *_parser) parseConditionlExpression() ast.Expression ***REMOVED***
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
	left := self.parseConditionlExpression()
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
	***REMOVED***

	if operator != 0 ***REMOVED***
		idx := self.idx
		self.next()
		switch left.(type) ***REMOVED***
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
		default:
			self.error(left.Idx0(), "Invalid left-hand side in assignment")
			self.nextStatement()
			return &ast.BadExpression***REMOVED***From: idx, To: self.idx***REMOVED***
		***REMOVED***
		return &ast.AssignExpression***REMOVED***
			Left:     left,
			Operator: operator,
			Right:    self.parseAssignmentExpression(),
		***REMOVED***
	***REMOVED***

	return left
***REMOVED***

func (self *_parser) parseExpression() ast.Expression ***REMOVED***
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
