package token

const (
	_ Token = iota

	ILLEGAL
	EOF
	COMMENT

	STRING
	NUMBER

	PLUS      // +
	MINUS     // -
	MULTIPLY  // *
	EXPONENT  // **
	SLASH     // /
	REMAINDER // %

	AND                  // &
	OR                   // |
	EXCLUSIVE_OR         // ^
	SHIFT_LEFT           // <<
	SHIFT_RIGHT          // >>
	UNSIGNED_SHIFT_RIGHT // >>>

	ADD_ASSIGN       // +=
	SUBTRACT_ASSIGN  // -=
	MULTIPLY_ASSIGN  // *=
	EXPONENT_ASSIGN  // **=
	QUOTIENT_ASSIGN  // /=
	REMAINDER_ASSIGN // %=

	AND_ASSIGN                  // &=
	OR_ASSIGN                   // |=
	EXCLUSIVE_OR_ASSIGN         // ^=
	SHIFT_LEFT_ASSIGN           // <<=
	SHIFT_RIGHT_ASSIGN          // >>=
	UNSIGNED_SHIFT_RIGHT_ASSIGN // >>>=

	LOGICAL_AND // &&
	LOGICAL_OR  // ||
	COALESCE    // ??
	INCREMENT   // ++
	DECREMENT   // --

	EQUAL        // ==
	STRICT_EQUAL // ===
	LESS         // <
	GREATER      // >
	ASSIGN       // =
	NOT          // !

	BITWISE_NOT // ~

	NOT_EQUAL        // !=
	STRICT_NOT_EQUAL // !==
	LESS_OR_EQUAL    // <=
	GREATER_OR_EQUAL // >=

	LEFT_PARENTHESIS // (
	LEFT_BRACKET     // [
	LEFT_BRACE       // ***REMOVED***
	COMMA            // ,
	PERIOD           // .

	RIGHT_PARENTHESIS // )
	RIGHT_BRACKET     // ]
	RIGHT_BRACE       // ***REMOVED***
	SEMICOLON         // ;
	COLON             // :
	QUESTION_MARK     // ?
	QUESTION_DOT      // ?.
	ARROW             // =>
	ELLIPSIS          // ...
	BACKTICK          // `

	// tokens below (and only them) are syntactically valid identifiers

	IDENTIFIER
	KEYWORD
	BOOLEAN
	NULL

	IF
	IN
	OF
	DO

	VAR
	LET
	FOR
	NEW
	TRY

	THIS
	ELSE
	CASE
	VOID
	WITH

	CONST
	WHILE
	BREAK
	CATCH
	THROW

	RETURN
	TYPEOF
	DELETE
	SWITCH

	DEFAULT
	FINALLY

	FUNCTION
	CONTINUE
	DEBUGGER

	INSTANCEOF
)

var token2string = [...]string***REMOVED***
	ILLEGAL:                     "ILLEGAL",
	EOF:                         "EOF",
	COMMENT:                     "COMMENT",
	KEYWORD:                     "KEYWORD",
	STRING:                      "STRING",
	BOOLEAN:                     "BOOLEAN",
	NULL:                        "NULL",
	NUMBER:                      "NUMBER",
	IDENTIFIER:                  "IDENTIFIER",
	PLUS:                        "+",
	MINUS:                       "-",
	EXPONENT:                    "**",
	MULTIPLY:                    "*",
	SLASH:                       "/",
	REMAINDER:                   "%",
	AND:                         "&",
	OR:                          "|",
	EXCLUSIVE_OR:                "^",
	SHIFT_LEFT:                  "<<",
	SHIFT_RIGHT:                 ">>",
	UNSIGNED_SHIFT_RIGHT:        ">>>",
	ADD_ASSIGN:                  "+=",
	SUBTRACT_ASSIGN:             "-=",
	MULTIPLY_ASSIGN:             "*=",
	EXPONENT_ASSIGN:             "**=",
	QUOTIENT_ASSIGN:             "/=",
	REMAINDER_ASSIGN:            "%=",
	AND_ASSIGN:                  "&=",
	OR_ASSIGN:                   "|=",
	EXCLUSIVE_OR_ASSIGN:         "^=",
	SHIFT_LEFT_ASSIGN:           "<<=",
	SHIFT_RIGHT_ASSIGN:          ">>=",
	UNSIGNED_SHIFT_RIGHT_ASSIGN: ">>>=",
	LOGICAL_AND:                 "&&",
	LOGICAL_OR:                  "||",
	COALESCE:                    "??",
	INCREMENT:                   "++",
	DECREMENT:                   "--",
	EQUAL:                       "==",
	STRICT_EQUAL:                "===",
	LESS:                        "<",
	GREATER:                     ">",
	ASSIGN:                      "=",
	NOT:                         "!",
	BITWISE_NOT:                 "~",
	NOT_EQUAL:                   "!=",
	STRICT_NOT_EQUAL:            "!==",
	LESS_OR_EQUAL:               "<=",
	GREATER_OR_EQUAL:            ">=",
	LEFT_PARENTHESIS:            "(",
	LEFT_BRACKET:                "[",
	LEFT_BRACE:                  "***REMOVED***",
	COMMA:                       ",",
	PERIOD:                      ".",
	RIGHT_PARENTHESIS:           ")",
	RIGHT_BRACKET:               "]",
	RIGHT_BRACE:                 "***REMOVED***",
	SEMICOLON:                   ";",
	COLON:                       ":",
	QUESTION_MARK:               "?",
	QUESTION_DOT:                "?.",
	ARROW:                       "=>",
	ELLIPSIS:                    "...",
	BACKTICK:                    "`",
	IF:                          "if",
	IN:                          "in",
	OF:                          "of",
	DO:                          "do",
	VAR:                         "var",
	LET:                         "let",
	FOR:                         "for",
	NEW:                         "new",
	TRY:                         "try",
	THIS:                        "this",
	ELSE:                        "else",
	CASE:                        "case",
	VOID:                        "void",
	WITH:                        "with",
	CONST:                       "const",
	WHILE:                       "while",
	BREAK:                       "break",
	CATCH:                       "catch",
	THROW:                       "throw",
	RETURN:                      "return",
	TYPEOF:                      "typeof",
	DELETE:                      "delete",
	SWITCH:                      "switch",
	DEFAULT:                     "default",
	FINALLY:                     "finally",
	FUNCTION:                    "function",
	CONTINUE:                    "continue",
	DEBUGGER:                    "debugger",
	INSTANCEOF:                  "instanceof",
***REMOVED***

var keywordTable = map[string]_keyword***REMOVED***
	"if": ***REMOVED***
		token: IF,
	***REMOVED***,
	"in": ***REMOVED***
		token: IN,
	***REMOVED***,
	"do": ***REMOVED***
		token: DO,
	***REMOVED***,
	"var": ***REMOVED***
		token: VAR,
	***REMOVED***,
	"for": ***REMOVED***
		token: FOR,
	***REMOVED***,
	"new": ***REMOVED***
		token: NEW,
	***REMOVED***,
	"try": ***REMOVED***
		token: TRY,
	***REMOVED***,
	"this": ***REMOVED***
		token: THIS,
	***REMOVED***,
	"else": ***REMOVED***
		token: ELSE,
	***REMOVED***,
	"case": ***REMOVED***
		token: CASE,
	***REMOVED***,
	"void": ***REMOVED***
		token: VOID,
	***REMOVED***,
	"with": ***REMOVED***
		token: WITH,
	***REMOVED***,
	"while": ***REMOVED***
		token: WHILE,
	***REMOVED***,
	"break": ***REMOVED***
		token: BREAK,
	***REMOVED***,
	"catch": ***REMOVED***
		token: CATCH,
	***REMOVED***,
	"throw": ***REMOVED***
		token: THROW,
	***REMOVED***,
	"return": ***REMOVED***
		token: RETURN,
	***REMOVED***,
	"typeof": ***REMOVED***
		token: TYPEOF,
	***REMOVED***,
	"delete": ***REMOVED***
		token: DELETE,
	***REMOVED***,
	"switch": ***REMOVED***
		token: SWITCH,
	***REMOVED***,
	"default": ***REMOVED***
		token: DEFAULT,
	***REMOVED***,
	"finally": ***REMOVED***
		token: FINALLY,
	***REMOVED***,
	"function": ***REMOVED***
		token: FUNCTION,
	***REMOVED***,
	"continue": ***REMOVED***
		token: CONTINUE,
	***REMOVED***,
	"debugger": ***REMOVED***
		token: DEBUGGER,
	***REMOVED***,
	"instanceof": ***REMOVED***
		token: INSTANCEOF,
	***REMOVED***,
	"const": ***REMOVED***
		token: CONST,
	***REMOVED***,
	"class": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"enum": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"export": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"extends": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"import": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"super": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"implements": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"interface": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"let": ***REMOVED***
		token:  LET,
		strict: true,
	***REMOVED***,
	"package": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"private": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"protected": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"public": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"static": ***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
***REMOVED***
