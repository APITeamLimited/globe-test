package parser

import (
	"encoding/base64"
	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/go-sourcemap/sourcemap"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

func (self *_parser) parseBlockStatement() *ast.BlockStatement ***REMOVED***
	node := &ast.BlockStatement***REMOVED******REMOVED***
	node.LeftBrace = self.expect(token.LEFT_BRACE)
	node.List = self.parseStatementList()
	node.RightBrace = self.expect(token.RIGHT_BRACE)

	return node
***REMOVED***

func (self *_parser) parseEmptyStatement() ast.Statement ***REMOVED***
	idx := self.expect(token.SEMICOLON)
	return &ast.EmptyStatement***REMOVED***Semicolon: idx***REMOVED***
***REMOVED***

func (self *_parser) parseStatementList() (list []ast.Statement) ***REMOVED***
	for self.token != token.RIGHT_BRACE && self.token != token.EOF ***REMOVED***
		list = append(list, self.parseStatement())
	***REMOVED***

	return
***REMOVED***

func (self *_parser) parseStatement() ast.Statement ***REMOVED***

	if self.token == token.EOF ***REMOVED***
		self.errorUnexpectedToken(self.token)
		return &ast.BadStatement***REMOVED***From: self.idx, To: self.idx + 1***REMOVED***
	***REMOVED***

	switch self.token ***REMOVED***
	case token.SEMICOLON:
		return self.parseEmptyStatement()
	case token.LEFT_BRACE:
		return self.parseBlockStatement()
	case token.IF:
		return self.parseIfStatement()
	case token.DO:
		return self.parseDoWhileStatement()
	case token.WHILE:
		return self.parseWhileStatement()
	case token.FOR:
		return self.parseForOrForInStatement()
	case token.BREAK:
		return self.parseBreakStatement()
	case token.CONTINUE:
		return self.parseContinueStatement()
	case token.DEBUGGER:
		return self.parseDebuggerStatement()
	case token.WITH:
		return self.parseWithStatement()
	case token.VAR:
		return self.parseVariableStatement()
	case token.FUNCTION:
		self.parseFunction(true)
		// FIXME
		return &ast.EmptyStatement***REMOVED******REMOVED***
	case token.SWITCH:
		return self.parseSwitchStatement()
	case token.RETURN:
		return self.parseReturnStatement()
	case token.THROW:
		return self.parseThrowStatement()
	case token.TRY:
		return self.parseTryStatement()
	***REMOVED***

	expression := self.parseExpression()

	if identifier, isIdentifier := expression.(*ast.Identifier); isIdentifier && self.token == token.COLON ***REMOVED***
		// LabelledStatement
		colon := self.idx
		self.next() // :
		label := identifier.Name
		for _, value := range self.scope.labels ***REMOVED***
			if label == value ***REMOVED***
				self.error(identifier.Idx0(), "Label '%s' already exists", label)
			***REMOVED***
		***REMOVED***
		self.scope.labels = append(self.scope.labels, label) // Push the label
		statement := self.parseStatement()
		self.scope.labels = self.scope.labels[:len(self.scope.labels)-1] // Pop the label
		return &ast.LabelledStatement***REMOVED***
			Label:     identifier,
			Colon:     colon,
			Statement: statement,
		***REMOVED***
	***REMOVED***

	self.optionalSemicolon()

	return &ast.ExpressionStatement***REMOVED***
		Expression: expression,
	***REMOVED***
***REMOVED***

func (self *_parser) parseTryStatement() ast.Statement ***REMOVED***

	node := &ast.TryStatement***REMOVED***
		Try:  self.expect(token.TRY),
		Body: self.parseBlockStatement(),
	***REMOVED***

	if self.token == token.CATCH ***REMOVED***
		catch := self.idx
		self.next()
		self.expect(token.LEFT_PARENTHESIS)
		if self.token != token.IDENTIFIER ***REMOVED***
			self.expect(token.IDENTIFIER)
			self.nextStatement()
			return &ast.BadStatement***REMOVED***From: catch, To: self.idx***REMOVED***
		***REMOVED*** else ***REMOVED***
			identifier := self.parseIdentifier()
			self.expect(token.RIGHT_PARENTHESIS)
			node.Catch = &ast.CatchStatement***REMOVED***
				Catch:     catch,
				Parameter: identifier,
				Body:      self.parseBlockStatement(),
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if self.token == token.FINALLY ***REMOVED***
		self.next()
		node.Finally = self.parseBlockStatement()
	***REMOVED***

	if node.Catch == nil && node.Finally == nil ***REMOVED***
		self.error(node.Try, "Missing catch or finally after try")
		return &ast.BadStatement***REMOVED***From: node.Try, To: node.Body.Idx1()***REMOVED***
	***REMOVED***

	return node
***REMOVED***

func (self *_parser) parseFunctionParameterList() *ast.ParameterList ***REMOVED***
	opening := self.expect(token.LEFT_PARENTHESIS)
	var list []*ast.Identifier
	for self.token != token.RIGHT_PARENTHESIS && self.token != token.EOF ***REMOVED***
		if self.token != token.IDENTIFIER ***REMOVED***
			self.expect(token.IDENTIFIER)
		***REMOVED*** else ***REMOVED***
			list = append(list, self.parseIdentifier())
		***REMOVED***
		if self.token != token.RIGHT_PARENTHESIS ***REMOVED***
			self.expect(token.COMMA)
		***REMOVED***
	***REMOVED***
	closing := self.expect(token.RIGHT_PARENTHESIS)

	return &ast.ParameterList***REMOVED***
		Opening: opening,
		List:    list,
		Closing: closing,
	***REMOVED***
***REMOVED***

func (self *_parser) parseParameterList() (list []string) ***REMOVED***
	for self.token != token.EOF ***REMOVED***
		if self.token != token.IDENTIFIER ***REMOVED***
			self.expect(token.IDENTIFIER)
		***REMOVED***
		list = append(list, self.literal)
		self.next()
		if self.token != token.EOF ***REMOVED***
			self.expect(token.COMMA)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (self *_parser) parseFunction(declaration bool) *ast.FunctionLiteral ***REMOVED***

	node := &ast.FunctionLiteral***REMOVED***
		Function: self.expect(token.FUNCTION),
	***REMOVED***

	var name *ast.Identifier
	if self.token == token.IDENTIFIER ***REMOVED***
		name = self.parseIdentifier()
		if declaration ***REMOVED***
			self.scope.declare(&ast.FunctionDeclaration***REMOVED***
				Function: node,
			***REMOVED***)
		***REMOVED***
	***REMOVED*** else if declaration ***REMOVED***
		// Use expect error handling
		self.expect(token.IDENTIFIER)
	***REMOVED***
	node.Name = name
	node.ParameterList = self.parseFunctionParameterList()
	self.parseFunctionBlock(node)
	node.Source = self.slice(node.Idx0(), node.Idx1())

	return node
***REMOVED***

func (self *_parser) parseFunctionBlock(node *ast.FunctionLiteral) ***REMOVED***
	***REMOVED***
		self.openScope()
		inFunction := self.scope.inFunction
		self.scope.inFunction = true
		defer func() ***REMOVED***
			self.scope.inFunction = inFunction
			self.closeScope()
		***REMOVED***()
		node.Body = self.parseBlockStatement()
		node.DeclarationList = self.scope.declarationList
	***REMOVED***
***REMOVED***

func (self *_parser) parseDebuggerStatement() ast.Statement ***REMOVED***
	idx := self.expect(token.DEBUGGER)

	node := &ast.DebuggerStatement***REMOVED***
		Debugger: idx,
	***REMOVED***

	self.semicolon()

	return node
***REMOVED***

func (self *_parser) parseReturnStatement() ast.Statement ***REMOVED***
	idx := self.expect(token.RETURN)

	if !self.scope.inFunction ***REMOVED***
		self.error(idx, "Illegal return statement")
		self.nextStatement()
		return &ast.BadStatement***REMOVED***From: idx, To: self.idx***REMOVED***
	***REMOVED***

	node := &ast.ReturnStatement***REMOVED***
		Return: idx,
	***REMOVED***

	if !self.implicitSemicolon && self.token != token.SEMICOLON && self.token != token.RIGHT_BRACE && self.token != token.EOF ***REMOVED***
		node.Argument = self.parseExpression()
	***REMOVED***

	self.semicolon()

	return node
***REMOVED***

func (self *_parser) parseThrowStatement() ast.Statement ***REMOVED***
	idx := self.expect(token.THROW)

	if self.implicitSemicolon ***REMOVED***
		if self.chr == -1 ***REMOVED*** // Hackish
			self.error(idx, "Unexpected end of input")
		***REMOVED*** else ***REMOVED***
			self.error(idx, "Illegal newline after throw")
		***REMOVED***
		self.nextStatement()
		return &ast.BadStatement***REMOVED***From: idx, To: self.idx***REMOVED***
	***REMOVED***

	node := &ast.ThrowStatement***REMOVED***
		Argument: self.parseExpression(),
	***REMOVED***

	self.semicolon()

	return node
***REMOVED***

func (self *_parser) parseSwitchStatement() ast.Statement ***REMOVED***
	self.expect(token.SWITCH)
	self.expect(token.LEFT_PARENTHESIS)
	node := &ast.SwitchStatement***REMOVED***
		Discriminant: self.parseExpression(),
		Default:      -1,
	***REMOVED***
	self.expect(token.RIGHT_PARENTHESIS)

	self.expect(token.LEFT_BRACE)

	inSwitch := self.scope.inSwitch
	self.scope.inSwitch = true
	defer func() ***REMOVED***
		self.scope.inSwitch = inSwitch
	***REMOVED***()

	for index := 0; self.token != token.EOF; index++ ***REMOVED***
		if self.token == token.RIGHT_BRACE ***REMOVED***
			self.next()
			break
		***REMOVED***

		clause := self.parseCaseStatement()
		if clause.Test == nil ***REMOVED***
			if node.Default != -1 ***REMOVED***
				self.error(clause.Case, "Already saw a default in switch")
			***REMOVED***
			node.Default = index
		***REMOVED***
		node.Body = append(node.Body, clause)
	***REMOVED***

	return node
***REMOVED***

func (self *_parser) parseWithStatement() ast.Statement ***REMOVED***
	self.expect(token.WITH)
	self.expect(token.LEFT_PARENTHESIS)
	node := &ast.WithStatement***REMOVED***
		Object: self.parseExpression(),
	***REMOVED***
	self.expect(token.RIGHT_PARENTHESIS)

	node.Body = self.parseStatement()

	return node
***REMOVED***

func (self *_parser) parseCaseStatement() *ast.CaseStatement ***REMOVED***

	node := &ast.CaseStatement***REMOVED***
		Case: self.idx,
	***REMOVED***
	if self.token == token.DEFAULT ***REMOVED***
		self.next()
	***REMOVED*** else ***REMOVED***
		self.expect(token.CASE)
		node.Test = self.parseExpression()
	***REMOVED***
	self.expect(token.COLON)

	for ***REMOVED***
		if self.token == token.EOF ||
			self.token == token.RIGHT_BRACE ||
			self.token == token.CASE ||
			self.token == token.DEFAULT ***REMOVED***
			break
		***REMOVED***
		node.Consequent = append(node.Consequent, self.parseStatement())

	***REMOVED***

	return node
***REMOVED***

func (self *_parser) parseIterationStatement() ast.Statement ***REMOVED***
	inIteration := self.scope.inIteration
	self.scope.inIteration = true
	defer func() ***REMOVED***
		self.scope.inIteration = inIteration
	***REMOVED***()
	return self.parseStatement()
***REMOVED***

func (self *_parser) parseForIn(idx file.Idx, into ast.Expression) *ast.ForInStatement ***REMOVED***

	// Already have consumed "<into> in"

	source := self.parseExpression()
	self.expect(token.RIGHT_PARENTHESIS)

	return &ast.ForInStatement***REMOVED***
		For:    idx,
		Into:   into,
		Source: source,
		Body:   self.parseIterationStatement(),
	***REMOVED***
***REMOVED***

func (self *_parser) parseForOf(idx file.Idx, into ast.Expression) *ast.ForOfStatement ***REMOVED***

	// Already have consumed "<into> of"

	source := self.parseExpression()
	self.expect(token.RIGHT_PARENTHESIS)

	return &ast.ForOfStatement***REMOVED***
		For:    idx,
		Into:   into,
		Source: source,
		Body:   self.parseIterationStatement(),
	***REMOVED***
***REMOVED***

func (self *_parser) parseFor(idx file.Idx, initializer ast.Expression) *ast.ForStatement ***REMOVED***

	// Already have consumed "<initializer> ;"

	var test, update ast.Expression

	if self.token != token.SEMICOLON ***REMOVED***
		test = self.parseExpression()
	***REMOVED***
	self.expect(token.SEMICOLON)

	if self.token != token.RIGHT_PARENTHESIS ***REMOVED***
		update = self.parseExpression()
	***REMOVED***
	self.expect(token.RIGHT_PARENTHESIS)

	return &ast.ForStatement***REMOVED***
		For:         idx,
		Initializer: initializer,
		Test:        test,
		Update:      update,
		Body:        self.parseIterationStatement(),
	***REMOVED***
***REMOVED***

func (self *_parser) parseForOrForInStatement() ast.Statement ***REMOVED***
	idx := self.expect(token.FOR)
	self.expect(token.LEFT_PARENTHESIS)

	var left []ast.Expression

	forIn := false
	forOf := false
	if self.token != token.SEMICOLON ***REMOVED***

		allowIn := self.scope.allowIn
		self.scope.allowIn = false
		if self.token == token.VAR ***REMOVED***
			var_ := self.idx
			self.next()
			list := self.parseVariableDeclarationList(var_)
			if len(list) == 1 ***REMOVED***
				if self.token == token.IN ***REMOVED***
					self.next() // in
					forIn = true
				***REMOVED*** else if self.token == token.IDENTIFIER ***REMOVED***
					if self.literal == "of" ***REMOVED***
						self.next()
						forOf = true
					***REMOVED***
				***REMOVED***
			***REMOVED***
			left = list
		***REMOVED*** else ***REMOVED***
			left = append(left, self.parseExpression())
			if self.token == token.IN ***REMOVED***
				self.next()
				forIn = true
			***REMOVED*** else if self.token == token.IDENTIFIER ***REMOVED***
				if self.literal == "of" ***REMOVED***
					self.next()
					forOf = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
		self.scope.allowIn = allowIn
	***REMOVED***

	if forIn || forOf ***REMOVED***
		switch left[0].(type) ***REMOVED***
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression, *ast.VariableExpression:
			// These are all acceptable
		default:
			self.error(idx, "Invalid left-hand side in for-in or for-of")
			self.nextStatement()
			return &ast.BadStatement***REMOVED***From: idx, To: self.idx***REMOVED***
		***REMOVED***
		if forIn ***REMOVED***
			return self.parseForIn(idx, left[0])
		***REMOVED***
		return self.parseForOf(idx, left[0])
	***REMOVED***

	self.expect(token.SEMICOLON)
	return self.parseFor(idx, &ast.SequenceExpression***REMOVED***Sequence: left***REMOVED***)
***REMOVED***

func (self *_parser) parseVariableStatement() *ast.VariableStatement ***REMOVED***

	idx := self.expect(token.VAR)

	list := self.parseVariableDeclarationList(idx)
	self.semicolon()

	return &ast.VariableStatement***REMOVED***
		Var:  idx,
		List: list,
	***REMOVED***
***REMOVED***

func (self *_parser) parseDoWhileStatement() ast.Statement ***REMOVED***
	inIteration := self.scope.inIteration
	self.scope.inIteration = true
	defer func() ***REMOVED***
		self.scope.inIteration = inIteration
	***REMOVED***()

	self.expect(token.DO)
	node := &ast.DoWhileStatement***REMOVED******REMOVED***
	if self.token == token.LEFT_BRACE ***REMOVED***
		node.Body = self.parseBlockStatement()
	***REMOVED*** else ***REMOVED***
		node.Body = self.parseStatement()
	***REMOVED***

	self.expect(token.WHILE)
	self.expect(token.LEFT_PARENTHESIS)
	node.Test = self.parseExpression()
	self.expect(token.RIGHT_PARENTHESIS)
	if self.token == token.SEMICOLON ***REMOVED***
		self.next()
	***REMOVED***

	return node
***REMOVED***

func (self *_parser) parseWhileStatement() ast.Statement ***REMOVED***
	self.expect(token.WHILE)
	self.expect(token.LEFT_PARENTHESIS)
	node := &ast.WhileStatement***REMOVED***
		Test: self.parseExpression(),
	***REMOVED***
	self.expect(token.RIGHT_PARENTHESIS)
	node.Body = self.parseIterationStatement()

	return node
***REMOVED***

func (self *_parser) parseIfStatement() ast.Statement ***REMOVED***
	self.expect(token.IF)
	self.expect(token.LEFT_PARENTHESIS)
	node := &ast.IfStatement***REMOVED***
		Test: self.parseExpression(),
	***REMOVED***
	self.expect(token.RIGHT_PARENTHESIS)

	if self.token == token.LEFT_BRACE ***REMOVED***
		node.Consequent = self.parseBlockStatement()
	***REMOVED*** else ***REMOVED***
		node.Consequent = self.parseStatement()
	***REMOVED***

	if self.token == token.ELSE ***REMOVED***
		self.next()
		node.Alternate = self.parseStatement()
	***REMOVED***

	return node
***REMOVED***

func (self *_parser) parseSourceElement() ast.Statement ***REMOVED***
	return self.parseStatement()
***REMOVED***

func (self *_parser) parseSourceElements() []ast.Statement ***REMOVED***
	body := []ast.Statement(nil)

	for ***REMOVED***
		if self.token != token.STRING ***REMOVED***
			break
		***REMOVED***

		body = append(body, self.parseSourceElement())
	***REMOVED***

	for self.token != token.EOF ***REMOVED***
		body = append(body, self.parseSourceElement())
	***REMOVED***

	return body
***REMOVED***

func (self *_parser) parseProgram() *ast.Program ***REMOVED***
	self.openScope()
	defer self.closeScope()
	return &ast.Program***REMOVED***
		Body:            self.parseSourceElements(),
		DeclarationList: self.scope.declarationList,
		File:            self.file,
		SourceMap:       self.parseSourceMap(),
	***REMOVED***
***REMOVED***

func (self *_parser) parseSourceMap() *sourcemap.Consumer ***REMOVED***
	lastLine := self.str[strings.LastIndexByte(self.str, '\n')+1:]
	if strings.HasPrefix(lastLine, "//# sourceMappingURL") ***REMOVED***
		urlIndex := strings.Index(lastLine, "=")
		urlStr := lastLine[urlIndex+1:]

		var data []byte
		if strings.HasPrefix(urlStr, "data:application/json") ***REMOVED***
			b64Index := strings.Index(urlStr, ",")
			b64 := urlStr[b64Index+1:]
			if d, err := base64.StdEncoding.DecodeString(b64); err == nil ***REMOVED***
				data = d
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if smUrl, err := url.Parse(urlStr); err == nil ***REMOVED***
				if smUrl.Scheme == "" || smUrl.Scheme == "file" ***REMOVED***
					if f, err := os.Open(smUrl.Path); err == nil ***REMOVED***
						if d, err := ioutil.ReadAll(f); err == nil ***REMOVED***
							data = d
						***REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					// Not implemented - compile error?
					return nil
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if data == nil ***REMOVED***
			return nil
		***REMOVED***

		if sm, err := sourcemap.Parse(self.file.Name(), data); err == nil ***REMOVED***
			return sm
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (self *_parser) parseBreakStatement() ast.Statement ***REMOVED***
	idx := self.expect(token.BREAK)
	semicolon := self.implicitSemicolon
	if self.token == token.SEMICOLON ***REMOVED***
		semicolon = true
		self.next()
	***REMOVED***

	if semicolon || self.token == token.RIGHT_BRACE ***REMOVED***
		self.implicitSemicolon = false
		if !self.scope.inIteration && !self.scope.inSwitch ***REMOVED***
			goto illegal
		***REMOVED***
		return &ast.BranchStatement***REMOVED***
			Idx:   idx,
			Token: token.BREAK,
		***REMOVED***
	***REMOVED***

	if self.token == token.IDENTIFIER ***REMOVED***
		identifier := self.parseIdentifier()
		if !self.scope.hasLabel(identifier.Name) ***REMOVED***
			self.error(idx, "Undefined label '%s'", identifier.Name)
			return &ast.BadStatement***REMOVED***From: idx, To: identifier.Idx1()***REMOVED***
		***REMOVED***
		self.semicolon()
		return &ast.BranchStatement***REMOVED***
			Idx:   idx,
			Token: token.BREAK,
			Label: identifier,
		***REMOVED***
	***REMOVED***

	self.expect(token.IDENTIFIER)

illegal:
	self.error(idx, "Illegal break statement")
	self.nextStatement()
	return &ast.BadStatement***REMOVED***From: idx, To: self.idx***REMOVED***
***REMOVED***

func (self *_parser) parseContinueStatement() ast.Statement ***REMOVED***
	idx := self.expect(token.CONTINUE)
	semicolon := self.implicitSemicolon
	if self.token == token.SEMICOLON ***REMOVED***
		semicolon = true
		self.next()
	***REMOVED***

	if semicolon || self.token == token.RIGHT_BRACE ***REMOVED***
		self.implicitSemicolon = false
		if !self.scope.inIteration ***REMOVED***
			goto illegal
		***REMOVED***
		return &ast.BranchStatement***REMOVED***
			Idx:   idx,
			Token: token.CONTINUE,
		***REMOVED***
	***REMOVED***

	if self.token == token.IDENTIFIER ***REMOVED***
		identifier := self.parseIdentifier()
		if !self.scope.hasLabel(identifier.Name) ***REMOVED***
			self.error(idx, "Undefined label '%s'", identifier.Name)
			return &ast.BadStatement***REMOVED***From: idx, To: identifier.Idx1()***REMOVED***
		***REMOVED***
		if !self.scope.inIteration ***REMOVED***
			goto illegal
		***REMOVED***
		self.semicolon()
		return &ast.BranchStatement***REMOVED***
			Idx:   idx,
			Token: token.CONTINUE,
			Label: identifier,
		***REMOVED***
	***REMOVED***

	self.expect(token.IDENTIFIER)

illegal:
	self.error(idx, "Illegal continue statement")
	self.nextStatement()
	return &ast.BadStatement***REMOVED***From: idx, To: self.idx***REMOVED***
***REMOVED***

// Find the next statement after an error (recover)
func (self *_parser) nextStatement() ***REMOVED***
	for ***REMOVED***
		switch self.token ***REMOVED***
		case token.BREAK, token.CONTINUE,
			token.FOR, token.IF, token.RETURN, token.SWITCH,
			token.VAR, token.DO, token.TRY, token.WITH,
			token.WHILE, token.THROW, token.CATCH, token.FINALLY:
			// Return only if parser made some progress since last
			// sync or if it has not reached 10 next calls without
			// progress. Otherwise consume at least one token to
			// avoid an endless parser loop
			if self.idx == self.recover.idx && self.recover.count < 10 ***REMOVED***
				self.recover.count++
				return
			***REMOVED***
			if self.idx > self.recover.idx ***REMOVED***
				self.recover.idx = self.idx
				self.recover.count = 0
				return
			***REMOVED***
			// Reaching here indicates a parser bug, likely an
			// incorrect token list in this function, but it only
			// leads to skipping of possibly correct code if a
			// previous error is present, and thus is preferred
			// over a non-terminating parse.
		case token.EOF:
			return
		***REMOVED***
		self.next()
	***REMOVED***
***REMOVED***
