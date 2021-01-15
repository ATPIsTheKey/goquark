package token

type TokenKind int

//go:generate stringer -type=TokenKind
const (
	special_beg TokenKind = iota
	COMMENT
	SKIP
	MISMATCH
	EOF
	special_end

	literal_beg
	IDENT
	INTEGER
	REAL
	COMPLEX
	STRING
	BOOLEAN
	literal_end

	operator_beg

	strictly_unary_operator_beg
	TILDE
	XOR
	HEAD
	TAIL
	NIL
	NOT
	MAX
	MIN
	CONJ
	ABS
	strictly_unary_operator_end

	unary_binary_operator_beg
	PLUS
	MINUS
	unary_binary_operator_end

	PERCENT
	STAR
	DOUBLE_STAR
	SLASH
	DOUBLE_SLASH
	SLASH_PERCENT
	COLON
	DOUBLE_COLON
	AT
	ON
	AMPERSAND
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	EQUAL
	DOUBLE_EQUAL
	EXCLAMATION_EQUAL
	AND
	OR
	DOUBLE_COLON_EQUAL
	EQUAL_GREATER
	CIRCUMFLEX
	VERTICAL_BAR
	operator_end

	keyword_beg
	CASE
	OF
	IF
	THEN
	ELIF
	ELSE
	LET
	DEF
	DEFUN
	REDEF
	LETREC
	CONST
	WITH
	IN
	LAMBDA
	FUN
	IMPORT
	EXPORT
	AS
	COND
	OTHERWISE
	keyword_end

	separator_beg
	COMMA
	PERIOD
	QUESTION_MARK_LEFT_PARENTHESIS
	LEFT_PARENTHESIS
	RIGHT_PARENTHESIS
	LEFT_BRACKET
	LEFT_BRACELET
	RIGHT_BRACELET
	RIGHT_BRACKET
	QUOTE
	ELLIPSIS
	SEMICOLON
	BACKSLASH
	QUESTION_MARK
	TAB
	SPACE
	NEWLINE
	EOS
	separator_end
)

type TokenPos struct {
	Line   int
	Column int
}

type Token struct {
	FPos TokenPos
	Kind TokenKind
	Raw  string
}

var Single_char_tokens = map[string]TokenKind{
	"+":  PLUS,
	"-":  MINUS,
	"%":  PERCENT,
	"^":  CIRCUMFLEX,
	"@":  AT,
	"&":  AMPERSAND,
	"\\": BACKSLASH,
	",":  COMMA,
	".":  PERIOD,
	"{":  LEFT_BRACELET,
	"}":  RIGHT_BRACELET,
	"(":  LEFT_PARENTHESIS,
	")":  RIGHT_PARENTHESIS,
	"[":  LEFT_BRACKET,
	"]":  RIGHT_BRACKET,
	"|":  VERTICAL_BAR,
	"~":  TILDE,
	"\"": QUOTE,
	";":  SEMICOLON,
	"*":  STAR,
	"/":  SLASH,
	"<":  LESS,
	">":  GREATER,
	"=":  EQUAL,
	":":  COLON,
	"?":  QUESTION_MARK,
	"\n": NEWLINE,
}

var Double_char_tokens = map[string]TokenKind{
	"!=": EXCLAMATION_EQUAL,
	"**": DOUBLE_STAR,
	"//": DOUBLE_SLASH,
	"/%": SLASH_PERCENT,
	"<=": LESS_EQUAL,
	">=": GREATER_EQUAL,
	"==": DOUBLE_EQUAL,
	"=>": EQUAL_GREATER,
	"::": DOUBLE_COLON,
	"?(": QUESTION_MARK_LEFT_PARENTHESIS,
}

var Triple_char_tokens = map[string]TokenKind{
	"::=": DOUBLE_COLON_EQUAL,
	"...": ELLIPSIS,
}

var Keyword_tokens = map[string]TokenKind{
	"not":       NOT,
	"and":       AND,
	"or":        OR,
	"case":      CASE,
	"of":        OF,
	"xor":       XOR,
	"head":      HEAD,
	"tail":      TAIL,
	"nil":       NIL,
	"min":       MIN,
	"max":       MAX,
	"conj":      CONJ,
	"abs":       ABS,
	"cond":      COND,
	"if":        IF,
	"then":      THEN,
	"elif":      ELIF,
	"else":      ELSE,
	"let":       LET,
	"letrec":    LETREC,
	"const":     CONST,
	"with":      WITH,
	"in":        IN,
	"on":        ON,
	"as":        AS,
	"def":       DEF,
	"redef":     REDEF,
	"defun":     DEFUN,
	"lambda":    LAMBDA,
	"fun":       FUN,
	"import":    IMPORT,
	"export":    EXPORT,
	"otherwise": OTHERWISE,
}

const (
	LowestPrec  = 0
	HighestPrec = 11
)

func (tok TokenKind) Precedence() int {
	switch tok {
	case OR, XOR:
		return 1
	case AND:
		return 2
	case DOUBLE_EQUAL, EXCLAMATION_EQUAL, GREATER, GREATER_EQUAL, LESS, LESS_EQUAL:
		return 3
	case PLUS, MINUS:
		return 4
	case STAR, SLASH, DOUBLE_SLASH, PERCENT, SLASH_PERCENT:
		return 5
	case DOUBLE_STAR:
		return 6
	case NIL:
		return 7
	case HEAD, TAIL:
		return 8
	case VERTICAL_BAR:
		return 9
	case ON:
		return 11
	}
	return LowestPrec
}

func (tok TokenKind) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

func (tok TokenKind) IsOperator() bool { return operator_beg < tok && tok < operator_end }

func (tok TokenKind) IsUnaryOperator() bool {
	return strictly_unary_operator_beg < tok && tok < strictly_unary_operator_end
}

func (tok TokenKind) IsLeftAssociative() bool { return !tok.IsRightAssociative() }

func (tok TokenKind) IsRightAssociative() bool { return tok == VERTICAL_BAR || tok == DOUBLE_STAR }

func (tok TokenKind) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

func (tok TokenKind) IsSeparator() bool { return separator_beg < tok && tok < separator_end }

func (tok TokenKind) IsSpecial() bool { return special_beg < tok && tok < special_end }

func MakeToken(linePos int, columnPos int, kind TokenKind, val string) Token {
	return Token{TokenPos{linePos, columnPos}, kind, val}
}
