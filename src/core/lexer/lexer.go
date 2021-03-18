package lexer

import (
	"bufio"
	"fmt"
	"goquark/src/core/token"
	"goquark/src/core/utils"
	"os"
	"unicode"
)

const (
	NOOPTS           = 0
	IgnoreSkippables = 1 << iota
)

type ErrHandler func(tokPos utils.SourceIndex, msg string)

func BaseErrHandler(pos utils.SourceIndex, msg string) {
	fmt.Printf("Illegal token at (line: %d, col: %d): %s\n", pos.Line, pos.Column, msg)
}

type Lexer struct {
	src  []byte
	len  int
	opts int

	start     int
	pos       int
	columnPos int
	linePos   int

	errHandler ErrHandler
}

func (lexer *Lexer) Init(src []byte, errHandler ErrHandler, opts int) {
	// Explicitly initialize all fields since a lexer may be reused.
	lexer.len = len(src)

	lexer.src = src
	lexer.opts = opts

	lexer.start, lexer.pos = 0, 0
	lexer.columnPos, lexer.linePos = 0, 0

	lexer.errHandler = errHandler
}

func (lexer *Lexer) currentFilePos() utils.SourceIndex {
	return utils.SourceIndex{Line: lexer.linePos, Column: lexer.columnPos}
}

func (lexer *Lexer) currentChar() byte {
	if lexer.reachedEndOfSource(0) {
		return 0x00
	} else {
		return lexer.src[lexer.pos]
	}
}

func (lexer *Lexer) currentNChars(n int) []byte {
	if lexer.reachedEndOfSource(n) {
		return nil
	} else {
		return lexer.src[lexer.pos : lexer.pos+n]
	}
}

func (lexer *Lexer) consumeChar() {
	if lexer.currentChar() == '\n' {
		lexer.columnPos = 1
		lexer.linePos += 1
	}
	lexer.columnPos += 1
	lexer.pos += 1
}

func (lexer *Lexer) consumeNChars(n int) {
	for i := 0; i < n; i++ {
		lexer.consumeChar()
	}
}

func (lexer *Lexer) consumedChars() []byte { return lexer.src[lexer.start:lexer.pos] }

func (lexer *Lexer) discardConsumedChars() {
	lexer.start = lexer.pos
}

func (lexer *Lexer) matchChars(chs ...byte) bool {
	if lexer.reachedEndOfSource(len(chs) - 1) {
		return false
	}
	for i, c := range chs {
		if lexer.src[lexer.pos+i] != c {
			return false
		}
	}
	return true
}

func (lexer *Lexer) expectChars(msg string, chs ...byte) {
	if !lexer.matchChars(chs...) {
		lexer.errHandler(lexer.currentFilePos(), msg)
	}
}

func (lexer *Lexer) reachedEndOfSource(off int) bool { return lexer.pos >= lexer.len-off }

func (lexer *Lexer) isIDStart() bool {
	ch := lexer.currentChar()
	return 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || ch == '_'
}

func (lexer *Lexer) isIDChar() bool {
	ch := lexer.currentChar()
	return 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || ch == '_' || '0' <= ch && ch <= '9'
}

func (lexer *Lexer) isNumChar() bool {
	ch := lexer.currentChar()
	return '0' <= ch && ch <= '9'
}

func (lexer *Lexer) isStrStart() bool {
	ch := lexer.currentChar()
	return ch == '"'
}

func (lexer *Lexer) isStrChar() bool {
	ch := lexer.currentChar()
	return ch != '"' && ch < unicode.MaxASCII
}

func (lexer *Lexer) isSkippable() bool {
	ch := lexer.currentChar()
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func (lexer *Lexer) Next() token.Token {
	var ok bool
	var sourcePos utils.SourceIndex
	var kind token.TokenKind
	var lit string

	sourcePos = lexer.currentFilePos()

	if lexer.isSkippable() {
		lexer.consumeChar()
		for lexer.isSkippable() {
			lexer.consumeChar()
		}
		if lexer.opts&IgnoreSkippables > 0 {
			lexer.discardConsumedChars()
			if lexer.reachedEndOfSource(0) {
				return token.Token{Pos: sourcePos, Kind: token.Eos, Raw: ""}
			} else {
				return lexer.Next()
			}
		} else {
			kind, lit = token.Skip, string(lexer.consumedChars())
			lexer.discardConsumedChars()
		}

	} else if kind, ok = token.TripleCharTokens[string(lexer.currentNChars(3))]; ok {
		/* Lex triple char Tokens */
		lexer.consumeNChars(3)
		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if lexer.isNumChar() {
		/* Lex Integer, Real, Complex */
		// Only Real token are permitted to have leading zeros, as well as the Integer token zero itself.
		mustBeRealOrZero := lexer.matchChars('0')
		lexer.consumeChar()
		for lexer.isNumChar() {
			if mustBeRealOrZero {
				lexer.expectChars("Leading zeros in decimal integer literals are not permitted!", '0')
			}
			lexer.consumeChar()
		}
		// Check if token is of type Real
		if lexer.matchChars('.') {
			lexer.consumeChar()
			// decimal part of Real literal not necessary: f.e. "1." , "3432." are valid Real literal
			for lexer.isNumChar() {
				lexer.consumeChar()
			}
			kind = token.Real
		} else {
			kind = token.Integer
		}
		// If any number literal is prefixed with "i", it is a Complex literal
		if lexer.matchChars('i') {
			kind = token.Complex
			lexer.consumeNChars(1)
		}

		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if lexer.matchChars('.') {
		/* Lex Period, Real, Complex */
		// Real literals may omit the integer part.
		lexer.consumeChar()
		if lexer.isNumChar() {
			lexer.consumeChar()
			for lexer.isNumChar() {
				lexer.consumeChar()
			}
			// "i" can also be prefixed to Real literals omitting the integer part
			if lexer.matchChars('i') {
				kind = token.Complex
				lexer.consumeNChars(1)
			} else {
				kind = token.Real
			}
		} else { // If no number chars are present, the token is a Period
			kind = token.Period
		}

		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if lexer.isStrStart() {
		/* Check if current byte is start of a string literal, then consume it */
		lexer.consumeChar() // byte : '"'
		for lexer.isStrChar() {
			lexer.consumeChar()
		}
		lexer.expectChars("EOL while scanning string literal!", '"')
		lexer.consumeChar()
		kind, lit = token.String, string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if kind, ok = token.DoubleCharTokens[string(lexer.currentNChars(2))]; ok {
		/* Lex double char Tokens */
		lexer.consumeNChars(2)
		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if kind, ok = token.SingleCharTokens[string(lexer.currentChar())]; ok {
		/* Lex single char Tokens */
		lexer.consumeChar()
		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if lexer.isIDStart() {
		/* Lex Ident, keyword Tokens */
		lexer.consumeChar()
		for lexer.isIDChar() {
			lexer.consumeChar()
		}

		lit = string(lexer.consumedChars())
		if kind, ok = token.KeywordTokens[lit]; !ok {
			kind = token.Ident
		}
		lexer.discardConsumedChars()

	} else {
		/* No token could be lexed */
		lexer.errHandler(utils.SourceIndex{Line: lexer.linePos, Column: lexer.columnPos},
			fmt.Sprintf("Invalid character %c in identifier!", lexer.currentChar()))
		lit = string(lexer.consumedChars())
		kind = token.Mismatch
		lexer.discardConsumedChars()
	}

	return token.Token{Pos: sourcePos, Kind: kind, Raw: lit}
}

func (lexer *Lexer) GetTokens() []token.Token {
	var toks []token.Token

	for !lexer.reachedEndOfSource(0) {
		toks = append(toks, lexer.Next())
	}

	return toks
}

func TestInputLoop() {
	newReader := bufio.NewReader(os.Stdin)
	newLexer := &Lexer{}

	for {
		fmt.Printf(">>> ")
		text, _ := newReader.ReadBytes('\n')
		newLexer.Init(text, BaseErrHandler, IgnoreSkippables)
		for !newLexer.reachedEndOfSource(0) {
			t := newLexer.Next()
			fmt.Printf("<%s>: %#v\n", t.Kind.String(), t.Raw)
		}
	}
}
