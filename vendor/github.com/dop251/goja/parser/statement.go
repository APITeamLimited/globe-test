package parser

import (
	"encoding/base64"
	"fmt"
	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/go-sourcemap/sourcemap"
	"io/ioutil"
	"net/url"
	"path"
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
		self.scope.allowLet = true
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
	case token.LET:
		tok := self.peek()
		if tok == token.LEFT_BRACKET || self.scope.allowLet && (tok == token.IDENTIFIER || tok == token.LET || tok == token.LEFT_BRACE) ***REMOVED***
			return self.parseLexicalDeclaration(self.token)
		***REMOVED***
		self.insertSemicolon = true
	case token.CONST:
		return self.parseLexicalDeclaration(self.token)
	case token.FUNCTION:
		return &ast.FunctionDeclaration***REMOVED***
			Function: self.parseFunction(true),
		***REMOVED***
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
		self.scope.allowLet = false
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
		var parameter ast.BindingTarget
		if self.token == token.LEFT_PARENTHESIS ***REMOVED***
			self.next()
			parameter = self.parseBindingTarget()
			self.expect(token.RIGHT_PARENTHESIS)
		***REMOVED***
		node.Catch = &ast.CatchStatement***REMOVED***
			Catch:     catch,
			Parameter: parameter,
			Body:      self.parseBlockStatement(),
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
	var list []*ast.Binding
	var rest ast.Expression
	for self.token != token.RIGHT_PARENTHESIS && self.token != token.EOF ***REMOVED***
		if self.token == token.ELLIPSIS ***REMOVED***
			self.next()
			rest = self.reinterpretAsDestructBindingTarget(self.parseAssignmentExpression())
			break
		***REMOVED***
		self.parseVariableDeclaration(&list)
		if self.token != token.RIGHT_PARENTHESIS ***REMOVED***
			self.expect(token.COMMA)
		***REMOVED***
	***REMOVED***
	closing := self.expect(token.RIGHT_PARENTHESIS)

	return &ast.ParameterList***REMOVED***
		Opening: opening,
		List:    list,
		Rest:    rest,
		Closing: closing,
	***REMOVED***
***REMOVED***

func (self *_parser) parseFunction(declaration bool) *ast.FunctionLiteral ***REMOVED***

	node := &ast.FunctionLiteral***REMOVED***
		Function: self.expect(token.FUNCTION),
	***REMOVED***

	var name *ast.Identifier
	if self.token == token.IDENTIFIER ***REMOVED***
		name = self.parseIdentifier()
	***REMOVED*** else if declaration ***REMOVED***
		// Use expect error handling
		self.expect(token.IDENTIFIER)
	***REMOVED***
	node.Name = name
	node.ParameterList = self.parseFunctionParameterList()
	node.Body, node.DeclarationList = self.parseFunctionBlock()
	node.Source = self.slice(node.Idx0(), node.Idx1())

	return node
***REMOVED***

func (self *_parser) parseFunctionBlock() (body *ast.BlockStatement, declarationList []*ast.VariableDeclaration) ***REMOVED***
	self.openScope()
	inFunction := self.scope.inFunction
	self.scope.inFunction = true
	defer func() ***REMOVED***
		self.scope.inFunction = inFunction
		self.closeScope()
	***REMOVED***()
	body = self.parseBlockStatement()
	declarationList = self.scope.declarationList
	return
***REMOVED***

func (self *_parser) parseArrowFunctionBody() (ast.ConciseBody, []*ast.VariableDeclaration) ***REMOVED***
	if self.token == token.LEFT_BRACE ***REMOVED***
		return self.parseFunctionBlock()
	***REMOVED***
	return &ast.ExpressionBody***REMOVED***
		Expression: self.parseAssignmentExpression(),
	***REMOVED***, nil
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
	self.scope.allowLet = false
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
	self.scope.allowLet = false
	return self.parseStatement()
***REMOVED***

func (self *_parser) parseForIn(idx file.Idx, into ast.ForInto) *ast.ForInStatement ***REMOVED***

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

func (self *_parser) parseForOf(idx file.Idx, into ast.ForInto) *ast.ForOfStatement ***REMOVED***

	// Already have consumed "<into> of"

	source := self.parseAssignmentExpression()
	self.expect(token.RIGHT_PARENTHESIS)

	return &ast.ForOfStatement***REMOVED***
		For:    idx,
		Into:   into,
		Source: source,
		Body:   self.parseIterationStatement(),
	***REMOVED***
***REMOVED***

func (self *_parser) parseFor(idx file.Idx, initializer ast.ForLoopInitializer) *ast.ForStatement ***REMOVED***

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

	var initializer ast.ForLoopInitializer

	forIn := false
	forOf := false
	var into ast.ForInto
	if self.token != token.SEMICOLON ***REMOVED***

		allowIn := self.scope.allowIn
		self.scope.allowIn = false
		tok := self.token
		if tok == token.LET ***REMOVED***
			switch self.peek() ***REMOVED***
			case token.IDENTIFIER, token.LEFT_BRACKET, token.LEFT_BRACE:
			default:
				tok = token.IDENTIFIER
			***REMOVED***
		***REMOVED***
		if tok == token.VAR || tok == token.LET || tok == token.CONST ***REMOVED***
			idx := self.idx
			self.next()
			var list []*ast.Binding
			if tok == token.VAR ***REMOVED***
				list = self.parseVarDeclarationList(idx)
			***REMOVED*** else ***REMOVED***
				list = self.parseVariableDeclarationList()
			***REMOVED***
			if len(list) == 1 ***REMOVED***
				if self.token == token.IN ***REMOVED***
					self.next() // in
					forIn = true
				***REMOVED*** else if self.token == token.IDENTIFIER && self.literal == "of" ***REMOVED***
					self.next()
					forOf = true
				***REMOVED***
			***REMOVED***
			if forIn || forOf ***REMOVED***
				if tok == token.VAR ***REMOVED***
					into = &ast.ForIntoVar***REMOVED***
						Binding: list[0],
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if list[0].Initializer != nil ***REMOVED***
						self.error(list[0].Initializer.Idx0(), "for-in loop variable declaration may not have an initializer")
					***REMOVED***
					into = &ast.ForDeclaration***REMOVED***
						Idx:     idx,
						IsConst: tok == token.CONST,
						Target:  list[0].Target,
					***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				self.ensurePatternInit(list)
				if tok == token.VAR ***REMOVED***
					initializer = &ast.ForLoopInitializerVarDeclList***REMOVED***
						List: list,
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					initializer = &ast.ForLoopInitializerLexicalDecl***REMOVED***
						LexicalDeclaration: ast.LexicalDeclaration***REMOVED***
							Idx:   idx,
							Token: tok,
							List:  list,
						***REMOVED***,
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			expr := self.parseExpression()
			if self.token == token.IN ***REMOVED***
				self.next()
				forIn = true
			***REMOVED*** else if self.token == token.IDENTIFIER && self.literal == "of" ***REMOVED***
				self.next()
				forOf = true
			***REMOVED***
			if forIn || forOf ***REMOVED***
				switch e := expr.(type) ***REMOVED***
				case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression, *ast.Binding:
					// These are all acceptable
				case *ast.ObjectLiteral:
					expr = self.reinterpretAsObjectAssignmentPattern(e)
				case *ast.ArrayLiteral:
					expr = self.reinterpretAsArrayAssignmentPattern(e)
				default:
					self.error(idx, "Invalid left-hand side in for-in or for-of")
					self.nextStatement()
					return &ast.BadStatement***REMOVED***From: idx, To: self.idx***REMOVED***
				***REMOVED***
				into = &ast.ForIntoExpression***REMOVED***
					Expression: expr,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				initializer = &ast.ForLoopInitializerExpression***REMOVED***
					Expression: expr,
				***REMOVED***
			***REMOVED***
		***REMOVED***
		self.scope.allowIn = allowIn
	***REMOVED***

	if forIn ***REMOVED***
		return self.parseForIn(idx, into)
	***REMOVED***
	if forOf ***REMOVED***
		return self.parseForOf(idx, into)
	***REMOVED***

	self.expect(token.SEMICOLON)
	return self.parseFor(idx, initializer)
***REMOVED***

func (self *_parser) ensurePatternInit(list []*ast.Binding) ***REMOVED***
	for _, item := range list ***REMOVED***
		if _, ok := item.Target.(ast.Pattern); ok ***REMOVED***
			if item.Initializer == nil ***REMOVED***
				self.error(item.Idx1(), "Missing initializer in destructuring declaration")
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (self *_parser) parseVariableStatement() *ast.VariableStatement ***REMOVED***

	idx := self.expect(token.VAR)

	list := self.parseVarDeclarationList(idx)
	self.ensurePatternInit(list)
	self.semicolon()

	return &ast.VariableStatement***REMOVED***
		Var:  idx,
		List: list,
	***REMOVED***
***REMOVED***

func (self *_parser) parseLexicalDeclaration(tok token.Token) *ast.LexicalDeclaration ***REMOVED***
	idx := self.expect(tok)
	if !self.scope.allowLet ***REMOVED***
		self.error(idx, "Lexical declaration cannot appear in a single-statement context")
	***REMOVED***

	list := self.parseVariableDeclarationList()
	self.ensurePatternInit(list)
	self.semicolon()

	return &ast.LexicalDeclaration***REMOVED***
		Idx:   idx,
		Token: tok,
		List:  list,
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
		self.scope.allowLet = false
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
		self.scope.allowLet = false
		node.Consequent = self.parseStatement()
	***REMOVED***

	if self.token == token.ELSE ***REMOVED***
		self.next()
		self.scope.allowLet = false
		node.Alternate = self.parseStatement()
	***REMOVED***

	return node
***REMOVED***

func (self *_parser) parseSourceElements() (body []ast.Statement) ***REMOVED***
	for self.token != token.EOF ***REMOVED***
		self.scope.allowLet = true
		body = append(body, self.parseStatement())
	***REMOVED***

	return body
***REMOVED***

func (self *_parser) parseProgram() *ast.Program ***REMOVED***
	self.openScope()
	defer self.closeScope()
	prg := &ast.Program***REMOVED***
		Body:            self.parseSourceElements(),
		DeclarationList: self.scope.declarationList,
		File:            self.file,
	***REMOVED***
	self.file.SetSourceMap(self.parseSourceMap())
	return prg
***REMOVED***

func extractSourceMapLine(str string) string ***REMOVED***
	for ***REMOVED***
		p := strings.LastIndexByte(str, '\n')
		line := str[p+1:]
		if line != "" && line != "***REMOVED***)" ***REMOVED***
			if strings.HasPrefix(line, "//# sourceMappingURL=") ***REMOVED***
				return line
			***REMOVED***
			break
		***REMOVED***
		if p >= 0 ***REMOVED***
			str = str[:p]
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func (self *_parser) parseSourceMap() *sourcemap.Consumer ***REMOVED***
	if self.opts.disableSourceMaps ***REMOVED***
		return nil
	***REMOVED***
	if smLine := extractSourceMapLine(self.str); smLine != "" ***REMOVED***
		urlIndex := strings.Index(smLine, "=")
		urlStr := smLine[urlIndex+1:]

		var data []byte
		var err error
		if strings.HasPrefix(urlStr, "data:application/json") ***REMOVED***
			b64Index := strings.Index(urlStr, ",")
			b64 := urlStr[b64Index+1:]
			data, err = base64.StdEncoding.DecodeString(b64)
		***REMOVED*** else ***REMOVED***
			var smUrl *url.URL
			if smUrl, err = url.Parse(urlStr); err == nil ***REMOVED***
				p := smUrl.Path
				if !path.IsAbs(p) ***REMOVED***
					baseName := self.file.Name()
					baseUrl, err1 := url.Parse(baseName)
					if err1 == nil && baseUrl.Scheme != "" ***REMOVED***
						baseUrl.Path = path.Join(path.Dir(baseUrl.Path), p)
						p = baseUrl.String()
					***REMOVED*** else ***REMOVED***
						p = path.Join(path.Dir(baseName), p)
					***REMOVED***
				***REMOVED***
				if self.opts.sourceMapLoader != nil ***REMOVED***
					data, err = self.opts.sourceMapLoader(p)
				***REMOVED*** else ***REMOVED***
					if smUrl.Scheme == "" || smUrl.Scheme == "file" ***REMOVED***
						data, err = ioutil.ReadFile(p)
					***REMOVED*** else ***REMOVED***
						err = fmt.Errorf("unsupported source map URL scheme: %s", smUrl.Scheme)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err != nil ***REMOVED***
			self.error(file.Idx(0), "Could not load source map: %v", err)
			return nil
		***REMOVED***
		if data == nil ***REMOVED***
			return nil
		***REMOVED***

		if sm, err := sourcemap.Parse(self.file.Name(), data); err == nil ***REMOVED***
			return sm
		***REMOVED*** else ***REMOVED***
			self.error(file.Idx(0), "Could not parse source map: %v", err)
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
