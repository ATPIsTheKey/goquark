package token

import "goquark/src/core/utils"

type TokenKind int

//go:generate stringer -type=TokenKind
const (
	specialBeg TokenKind = iota
	Comment
	Skip
	Mismatch
	Eof
	specialEnd

	literalBeg
	Ident
	Integer
	Real
	Complex
	String
	Boolean
	literalEnd

	operatorBeg
	strictlyUnaryOperatorBeg
	Not
	BNot
	strictlyUnaryOperatorEnd

	unaryBinaryOperatorBeg
	Plus
	Minus
	unaryBinaryOperatorEnd

	Percent
	Star
	DoubleStar
	Slash
	DoubleSlash
	SlashPercent
	Colon
	DoubleColon
	At
	On
	Less
	LessEqual
	Greater
	GreaterEqual
	Equal
	DoubleEqual
	ExclamationEqual
	DoubleExclamation
	DoublePlus
	And
	BAnd
	Or
	BOr
	Xor
	BXor
	DashGreater
	VerticalBarGreater
	operatorEnd

	keywordBeg
	Module
	Using
	As
	Case
	Of
	If
	Then
	Elif
	Else
	Let
	Rec
	In
	Def
	Fn
	Otherwise
	keywordEnd

	separatorBeg
	Comma
	Period
	LeftParenthesis
	RightParenthesis
	LeftBracket
	LeftBracelet
	RightBracelet
	RightBracket
	Quote
	Ellipsis
	Semicolon
	QuestionMark
	Tab
	Space
	Newline
	Eos
	separatorEnd
)

type Token struct {
	Pos  utils.SourceIndex
	Kind TokenKind
	Raw  string
}

var SingleCharTokens = map[string]TokenKind{
	"+":  Plus,
	"-":  Minus,
	"%":  Percent,
	"@":  At,
	",":  Comma,
	".":  Period,
	"{":  LeftBracelet,
	"}":  RightBracelet,
	"(":  LeftParenthesis,
	")":  RightParenthesis,
	"[":  LeftBracket,
	"]":  RightBracket,
	"\"": Quote,
	";":  Semicolon,
	"*":  Star,
	"/":  Slash,
	"<":  Less,
	">":  Greater,
	"=":  Equal,
	":":  Colon,
	"?":  QuestionMark,
	"\n": Newline,
}

var DoubleCharTokens = map[string]TokenKind{
	"!=": ExclamationEqual,
	"!!": DoubleExclamation,
	"++": DoublePlus,
	"|>": VerticalBarGreater,
	"**": DoubleStar,
	"//": DoubleSlash,
	"/%": SlashPercent,
	"<=": LessEqual,
	">=": GreaterEqual,
	"==": DoubleEqual,
	"->": DashGreater,
	"::": DoubleColon,
}

var TripleCharTokens = map[string]TokenKind{
	"...": Ellipsis,
}

var KeywordTokens = map[string]TokenKind{
	"module":    Module,
	"using":     Using,
	"as":        As,
	"case":      Case,
	"of":        Of,
	"otherwise": Otherwise,
	"not":       Not,
	"bnot":      BNot,
	"and":       And,
	"band":      BAnd,
	"or":        Or,
	"bor":       BOr,
	"xor":       Xor,
	"bxor":      BXor,
	"if":        If,
	"then":      Then,
	"elif":      Elif,
	"else":      Else,
	"let":       Let,
	"rec":       Rec,
	"in":        In,
	"def":       Def,
	"fn":        Fn,
}

func (kind TokenKind) IsLiteral() bool { return literalBeg < kind && kind < literalEnd }

func (kind TokenKind) IsOperator() bool { return operatorBeg < kind && kind < operatorEnd }

func (kind TokenKind) IsUnaryOperator() bool {
	return strictlyUnaryOperatorBeg < kind && kind < strictlyUnaryOperatorEnd
}

func (kind TokenKind) IsLeftAssociative() bool { return !kind.IsRightAssociative() }

func (kind TokenKind) IsRightAssociative() bool { return kind == DoubleStar || kind == At }

func (kind TokenKind) IsKeyword() bool { return keywordBeg < kind && kind < keywordEnd }

func (kind TokenKind) IsSeparator() bool { return separatorBeg < kind && kind < separatorEnd }

func (kind TokenKind) IsSpecial() bool { return specialBeg < kind && kind < specialEnd }
