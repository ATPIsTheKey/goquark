package lexer

import (
	"goquark/goquark/core/token"
)

type class int

const (
	keyword class = iota
	literal
	operator
	separator
	special
	none
)

func tokenClass(tok token.Token) class {
	switch {
	case tok.Kind.IsLiteral():
		return literal
	case tok.Kind.IsOperator():
		return operator
	case tok.Kind.IsKeyword():
		return keyword
	case tok.Kind.IsSeparator():
		return separator
	case tok.Kind.IsSpecial():
		return special
	default:
		return none
	}
}

type elt struct {
	tok   token.TokenKind
	lit   string
	class class
}

var testTokens = [...]elt{
	// Special testTokens
	{token.COMMENT, "/* a comment */", special},
	{token.COMMENT, "// a comment \n", special},
	{token.COMMENT, "/*\r*/", special},
	{token.COMMENT, "/**\r/*/", special}, // issue 11151
	{token.COMMENT, "/**\r\r/*/", special},
	{token.COMMENT, "//\r\n", special},

	// Identifiers and basic type literals
	{token.IDENT, "foobar", literal},
	{token.IDENT, "a۰۱۸", literal},
	{token.IDENT, "foo६४", literal},
	{token.IDENT, "bar９８７６", literal},
	{token.IDENT, "ŝ", literal},    // was bug (issue 4000)
	{token.IDENT, "ŝfoo", literal}, // was bug (issue 4000)
	{token.INTEGER, "0", literal},
	{token.INTEGER, "1", literal},
	{token.INTEGER, "123456789012345678890", literal},
	{token.INTEGER, "01234567", literal},
	{token.INTEGER, "0xcafebabe", literal},
	{token.REAL, "0.", literal},
	{token.REAL, ".0", literal},
	{token.REAL, "3.14159265", literal},
	{token.REAL, "1e0", literal},
	{token.REAL, "1e+100", literal},
	{token.REAL, "1e-100", literal},
	{token.REAL, "2.71828e-1000", literal},
	{token.COMPLEX, "0i", literal},
	{token.COMPLEX, "1i", literal},
	{token.COMPLEX, "012345678901234567889i", literal},
	{token.COMPLEX, "123456789012345678890i", literal},
	{token.COMPLEX, "0.i", literal},
	{token.COMPLEX, ".0i", literal},
	{token.COMPLEX, "3.14159265i", literal},
	{token.COMPLEX, "1e0i", literal},
	{token.COMPLEX, "1e+100i", literal},
	{token.COMPLEX, "1e-100i", literal},
	{token.COMPLEX, "2.71828e-1000i", literal},
	{token.STRING, "'a'", literal},
	{token.STRING, "`foobar`", literal},

	// Operators
	{token.PLUS, "+", operator},
	{token.MINUS, "-", operator},
	{token.STAR, "*", operator},
	{token.SLASH, "/", operator},
	{token.PERCENT, "%", operator},

	{token.TILDE, "~", operator},
	{token.AMPERSAND, "&", operator},
	{token.VERTICAL_BAR, "|", operator},
	{token.CIRCUMFLEX, "^", operator},

	{token.DOUBLE_EQUAL, "==", operator},
	{token.LESS, "<", operator},
	{token.GREATER, ">", operator},
	{token.EQUAL, "=", operator},

	{token.EXCLAMATION_EQUAL, "!=", operator},
	{token.LESS_EQUAL, "<=", operator},
	{token.GREATER_EQUAL, ">=", operator},
	{token.ELLIPSIS, "...", separator},
	{token.QUESTION_MARK, "?", operator},
	{token.QUESTION_MARK_LEFT_PARENTHESIS, "?(", separator},

	{token.LEFT_PARENTHESIS, "(", separator},
	{token.LEFT_BRACKET, "[", separator},
	{token.LEFT_BRACELET, "{", separator},
	{token.COMMA, ",", separator},
	{token.PERIOD, ".", separator},

	{token.RIGHT_PARENTHESIS, ")", separator},
	{token.RIGHT_BRACKET, "]", separator},
	{token.RIGHT_BRACELET, "}", separator},
	{token.SEMICOLON, ";", separator},
	{token.COLON, ":", separator},

	// Keywords
	{token.CASE, "const", keyword},
	{token.OF, "of", keyword},
	{token.COND, "cond", keyword},
	{token.OTHERWISE, "otherwise", keyword},
	{token.IF, "if", keyword},
	{token.THEN, "then", keyword},
	{token.ELIF, "elif", keyword},
	{token.ELSE, "else", keyword},
	{token.LET, "let", keyword},
	{token.LETREC, "letrec", keyword},
	{token.IN, "in", keyword},
	{token.DEF, "def", keyword},
	{token.DEFUN, "defun", keyword},
	{token.AS, "as", keyword},
	{token.FUN, "fun", keyword},
	{token.IMPORT, "import", keyword},
	{token.EXPORT, "export", keyword},
}
