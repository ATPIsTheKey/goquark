package lexer

import (
	"bufio"
	"fmt"
	"goquark/goquark/core/token"
	"os"
	"unicode"
)

const (
	NOOPTS            = 0
	IGNORE_SKIPPABLES = 1 << iota
)

type ErrorHandler func(tokPos token.TokenPos, msg string)

type Lexer struct {
	src  []byte
	len  int
	opts int

	start     int
	pos       int
	columnPos int
	linePos   int

	err ErrorHandler
}

func (lexer *Lexer) Init(src []byte, err ErrorHandler, opts int) {
	// Explicitly initialize all fields since a lexer may be reused.
	lexer.len = len(src)

	lexer.src = src
	lexer.opts = opts

	lexer.start, lexer.pos = 0, 0
	lexer.columnPos, lexer.linePos = 0, 0

	lexer.err = err
}

func (lexer *Lexer) currentFilePos() token.TokenPos {
	return token.TokenPos{Line: lexer.linePos, Column: lexer.columnPos}
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
		lexer.err(lexer.currentFilePos(), msg)
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
	var fPos token.TokenPos
	var kind token.TokenKind
	var lit string

	fPos = lexer.currentFilePos()

	if lexer.isSkippable() {
		lexer.consumeChar()
		for lexer.isSkippable() {
			lexer.consumeChar()
		}
		if lexer.opts&IGNORE_SKIPPABLES > 0 {
			lexer.discardConsumedChars()
			if lexer.reachedEndOfSource(0) {
				return token.Token{FPos: fPos, Kind: token.EOS, Raw: ""}
			} else {
				return lexer.Next()
			}
		} else {
			kind, lit = token.SKIP, string(lexer.consumedChars())
			lexer.discardConsumedChars()
		}

	} else if kind, ok = token.Triple_char_tokens[string(lexer.currentNChars(3))]; ok {
		/* Lex triple char Tokens */
		lexer.consumeNChars(3)
		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if lexer.isNumChar() {
		/* Lex INTEGER, REAL, COMPLEX */
		// Only REAL token are permitted to have leading zeros, as well as the INTEGER token zero itself.
		mustBeRealOrZero := lexer.matchChars('0')
		lexer.consumeChar()
		for lexer.isNumChar() {
			if mustBeRealOrZero {
				lexer.expectChars("leading zeros in decimal integer literals are not permitted", '0')
			}
			lexer.consumeChar()
		}
		// Check if token is of type REAL
		if lexer.matchChars('.') {
			lexer.consumeChar()
			// decimal part of REAL literal not necessary: f.e. "1." , "3432." are valid REAL literal
			for lexer.isNumChar() {
				lexer.consumeChar()
			}
			kind = token.REAL
		} else {
			kind = token.INTEGER
		}
		// If any number literal is prefixed with "im", it is a COMPLEX literal
		if lexer.matchChars('i', 'm') {
			kind = token.COMPLEX
			lexer.consumeNChars(2)
		}

		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if lexer.matchChars('.') {
		/* Lex PERIOD, REAL, COMPLEX */
		// REAL literals may omit the integer part.
		lexer.consumeChar()
		if lexer.isNumChar() {
			lexer.consumeChar()
			for lexer.isNumChar() {
				lexer.consumeChar()
			}
			// "im" can also be prefixed to REAL literals omitting the integer part
			if lexer.matchChars('i', 'm') {
				kind = token.COMPLEX
				lexer.consumeNChars(2)
			} else {
				kind = token.REAL
			}
		} else { // If no number chars are present, the token is a simple PERIOD
			kind = token.PERIOD
		}

		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if lexer.isStrStart() {
		/* ,Check if current byte is start of a string literal, then consume it */
		lexer.consumeChar() // byte : '"'
		for lexer.isStrChar() {
			lexer.consumeChar()
		}
		lexer.expectChars("EOL while scanning string literal", '"')
		lexer.consumeChar()
		kind, lit = token.STRING, string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if kind, ok = token.Double_char_tokens[string(lexer.currentNChars(2))]; ok {
		/* Lex double char Tokens */
		lexer.consumeNChars(2)
		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if kind, ok = token.Single_char_tokens[string(lexer.currentChar())]; ok {
		/* Lex single char Tokens */
		lexer.consumeChar()
		lit = string(lexer.consumedChars())
		lexer.discardConsumedChars()

	} else if lexer.isIDStart() {
		/* Lex IDENT, keyword Tokens */
		lexer.consumeChar()
		for lexer.isIDChar() {
			lexer.consumeChar()
		}

		lit = string(lexer.consumedChars())
		if kind, ok = token.Keyword_tokens[lit]; !ok {
			kind = token.IDENT
		}
		lexer.discardConsumedChars()

	} else {
		/* No token could be lexed */
		lexer.err(token.TokenPos{Line: lexer.linePos, Column: lexer.columnPos},
			fmt.Sprintf("invalid character %c in identifier", lexer.currentChar()))
		lit = string(lexer.consumedChars())
		kind = token.MISMATCH
		lexer.discardConsumedChars()
	}

	return token.Token{FPos: fPos, Kind: kind, Raw: lit}
}

func (lexer *Lexer) GetTokens() []token.Token {
	var toks []token.Token

	for !lexer.reachedEndOfSource(0) {
		toks = append(toks, lexer.Next())
	}

	return toks
}

func TestInputLoop() {
	reader := bufio.NewReader(os.Stdin)
	lexer := &Lexer{}

	errHandler := func(pos token.TokenPos, msg string) {
		fmt.Printf("LexerError at (line: %d, col: %d): %s\n", pos.Line, pos.Column, msg)
		os.Exit(-1)
	}

	for {
		fmt.Printf(">>> ")
		text, _ := reader.ReadBytes('\n')
		lexer.Init(text, errHandler, IGNORE_SKIPPABLES)
		for !lexer.reachedEndOfSource(0) {
			t := lexer.Next()
			fmt.Printf("<%s>: %#v\n", t.Kind.String(), t.Raw)
		}
	}
}
