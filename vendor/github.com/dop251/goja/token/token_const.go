package token

const (
	_ Token = iota

	ILLEGAL
	EOF
	COMMENT
	KEYWORD

	STRING
	BOOLEAN
	NULL
	NUMBER
	IDENTIFIER

	PLUS      // +
	MINUS     // -
	MULTIPLY  // *
	SLASH     // /
	REMAINDER // %

	AND                  // &
	OR                   // |
	EXCLUSIVE_OR         // ^
	SHIFT_LEFT           // <<
	SHIFT_RIGHT          // >>
	UNSIGNED_SHIFT_RIGHT // >>>
	AND_NOT              // &^

	ADD_ASSIGN       // +=
	SUBTRACT_ASSIGN  // -=
	MULTIPLY_ASSIGN  // *=
	QUOTIENT_ASSIGN  // /=
	REMAINDER_ASSIGN // %=

	AND_ASSIGN                  // &=
	OR_ASSIGN                   // |=
	EXCLUSIVE_OR_ASSIGN         // ^=
	SHIFT_LEFT_ASSIGN           // <<=
	SHIFT_RIGHT_ASSIGN          // >>=
	UNSIGNED_SHIFT_RIGHT_ASSIGN // >>>=
	AND_NOT_ASSIGN              // &^=

	LOGICAL_AND // &&
	LOGICAL_OR  // ||
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

	firstKeyword
	IF
	IN
	DO

	VAR
	FOR
	NEW
	TRY

	THIS
	ELSE
	CASE
	VOID
	WITH

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
	lastKeyword
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
	MULTIPLY:                    "*",
	SLASH:                       "/",
	REMAINDER:                   "%",
	AND:                         "&",
	OR:                          "|",
	EXCLUSIVE_OR:                "^",
	SHIFT_LEFT:                  "<<",
	SHIFT_RIGHT:                 ">>",
	UNSIGNED_SHIFT_RIGHT:        ">>>",
	AND_NOT:                     "&^",
	ADD_ASSIGN:                  "+=",
	SUBTRACT_ASSIGN:             "-=",
	MULTIPLY_ASSIGN:             "*=",
	QUOTIENT_ASSIGN:             "/=",
	REMAINDER_ASSIGN:            "%=",
	AND_ASSIGN:                  "&=",
	OR_ASSIGN:                   "|=",
	EXCLUSIVE_OR_ASSIGN:         "^=",
	SHIFT_LEFT_ASSIGN:           "<<=",
	SHIFT_RIGHT_ASSIGN:          ">>=",
	UNSIGNED_SHIFT_RIGHT_ASSIGN: ">>>=",
	AND_NOT_ASSIGN:              "&^=",
	LOGICAL_AND:                 "&&",
	LOGICAL_OR:                  "||",
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
	IF:                          "if",
	IN:                          "in",
	DO:                          "do",
	VAR:                         "var",
	FOR:                         "for",
	NEW:                         "new",
	TRY:                         "try",
	THIS:                        "this",
	ELSE:                        "else",
	CASE:                        "case",
	VOID:                        "void",
	WITH:                        "with",
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
	"if": _keyword***REMOVED***
		token: IF,
	***REMOVED***,
	"in": _keyword***REMOVED***
		token: IN,
	***REMOVED***,
	"do": _keyword***REMOVED***
		token: DO,
	***REMOVED***,
	"var": _keyword***REMOVED***
		token: VAR,
	***REMOVED***,
	"for": _keyword***REMOVED***
		token: FOR,
	***REMOVED***,
	"new": _keyword***REMOVED***
		token: NEW,
	***REMOVED***,
	"try": _keyword***REMOVED***
		token: TRY,
	***REMOVED***,
	"this": _keyword***REMOVED***
		token: THIS,
	***REMOVED***,
	"else": _keyword***REMOVED***
		token: ELSE,
	***REMOVED***,
	"case": _keyword***REMOVED***
		token: CASE,
	***REMOVED***,
	"void": _keyword***REMOVED***
		token: VOID,
	***REMOVED***,
	"with": _keyword***REMOVED***
		token: WITH,
	***REMOVED***,
	"while": _keyword***REMOVED***
		token: WHILE,
	***REMOVED***,
	"break": _keyword***REMOVED***
		token: BREAK,
	***REMOVED***,
	"catch": _keyword***REMOVED***
		token: CATCH,
	***REMOVED***,
	"throw": _keyword***REMOVED***
		token: THROW,
	***REMOVED***,
	"return": _keyword***REMOVED***
		token: RETURN,
	***REMOVED***,
	"typeof": _keyword***REMOVED***
		token: TYPEOF,
	***REMOVED***,
	"delete": _keyword***REMOVED***
		token: DELETE,
	***REMOVED***,
	"switch": _keyword***REMOVED***
		token: SWITCH,
	***REMOVED***,
	"default": _keyword***REMOVED***
		token: DEFAULT,
	***REMOVED***,
	"finally": _keyword***REMOVED***
		token: FINALLY,
	***REMOVED***,
	"function": _keyword***REMOVED***
		token: FUNCTION,
	***REMOVED***,
	"continue": _keyword***REMOVED***
		token: CONTINUE,
	***REMOVED***,
	"debugger": _keyword***REMOVED***
		token: DEBUGGER,
	***REMOVED***,
	"instanceof": _keyword***REMOVED***
		token: INSTANCEOF,
	***REMOVED***,
	"const": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"class": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"enum": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"export": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"extends": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"import": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"super": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
	***REMOVED***,
	"implements": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"interface": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"let": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"package": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"private": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"protected": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"public": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
	"static": _keyword***REMOVED***
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	***REMOVED***,
***REMOVED***
