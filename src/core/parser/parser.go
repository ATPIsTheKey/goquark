package parser

import (
	"bufio"
	"fmt"
	"goquark/src/core/ast"
	"goquark/src/core/lexer"
	"goquark/src/core/token"
	"os"
	"strconv"
	"strings"
)

var uniqueIdentCounter = 0

const (
	NoOpts = 0
)

type ErrHandler func(tok token.Token, msg string)

func BaseErrHandler(tok token.Token, msg string) {
	fmt.Printf("Illegal syntax at %s: %s\n", tok.Pos.Fmt(), msg)
	os.Exit(-1)
}

type Parser struct {
	src []token.Token
	len int
	pos int

	opts       int
	errHandler ErrHandler
}

type exprParseCtx struct {
	toBeAlphaConverted map[string]string
}

func newExprParseCtx() exprParseCtx {
	ctx := exprParseCtx{}
	ctx.toBeAlphaConverted = make(map[string]string)
	return ctx
}

func (ctx *exprParseCtx) enqueueAlphaConversion(identVal string, newIdentVal string) {
	ctx.toBeAlphaConverted[identVal] = newIdentVal
}

func (ctx *exprParseCtx) requiresAlphaConversion(identVal string) bool {
	_, ok := ctx.toBeAlphaConverted[identVal]
	return ok
}

func (parser *Parser) Init(src []token.Token, errHandler ErrHandler, opts int) {
	// Explicitly initialize all fields since a parser may be reused.
	parser.len = len(src)
	parser.src = src
	parser.pos = 0

	parser.opts = opts
	parser.errHandler = errHandler
}

func (parser *Parser) getUniqueIdentVal() string {
	ret := "%" + strconv.Itoa(uniqueIdentCounter)
	uniqueIdentCounter += 1
	return ret
}

func (parser *Parser) currentToken() token.Token { return parser.src[parser.pos] }

func (parser *Parser) getAndConsumeNextToken() token.Token {
	ret := parser.currentToken()
	parser.consumeToken()
	return ret
}

func (parser *Parser) reachedEndOfSource() bool { return parser.currentToken().Kind == token.Eos }

func (parser *Parser) consumeToken() {
	parser.pos += 1
}

func (parser *Parser) matchAnyTokenFrom(kinds ...token.TokenKind) bool {
	if parser.reachedEndOfSource() {
		return false
	} else {
		for _, kind := range kinds {
			if parser.currentToken().Kind == kind {
				return true
			}
		}
		return false
	}
}

func (parser *Parser) expectAnyTokenFrom(kinds ...token.TokenKind) {
	if parser.reachedEndOfSource() || !parser.matchAnyTokenFrom(kinds...) {
		var tokStrings []string

		for _, v := range kinds {
			tokStrings = append(tokStrings, v.String())
		}

		failMsg := fmt.Sprint("Expected any from { ", strings.Join(tokStrings, ", "), " } but got ",
			parser.currentToken().Kind.String())
		parser.errHandler(parser.currentToken(), failMsg)
	}
}

func (parser *Parser) parseIdentList() []string {
	var identList []string

	for gotComma := true; gotComma && !parser.reachedEndOfSource(); {
		parser.expectAnyTokenFrom(token.Ident)
		identList = append(identList, parser.getAndConsumeNextToken().Raw)
		if gotComma = parser.matchAnyTokenFrom(token.Comma); gotComma {
			parser.consumeToken()
		}
	}

	return identList
}

func (parser *Parser) parseStmt() ast.Stmt {
	switch {
	case parser.matchAnyTokenFrom(token.Def):
		return parser.parseDefStmt()
	default:
		parse := parser.parseExpr(newExprParseCtx())
		return parse
	}
}

func (parser *Parser) parseStmtList() []ast.Stmt {
	var stmts []ast.Stmt

	for !parser.reachedEndOfSource() {
		stmts = append(stmts, parser.parseStmt())
		parser.expectAnyTokenFrom(token.Semicolon)
		parser.consumeToken()
	}

	return stmts
}

func (parser *Parser) parseDefStmt() ast.Stmt {
	var names []string
	var exprs []ast.Expr
	var isRec bool

	pos := parser.getAndConsumeNextToken().Pos // def
	if isRec = parser.matchAnyTokenFrom(token.Rec); isRec {
		parser.consumeToken()
	}
	for gotComma := true; gotComma; {
		names = append(names, parser.getAndConsumeNextToken().Raw)
		parser.expectAnyTokenFrom(token.Equal)
		parser.consumeToken()
		exprs = append(exprs, parser.parseExpr(newExprParseCtx()))

		if gotComma = parser.matchAnyTokenFrom(token.Comma); gotComma {
			parser.consumeToken()
		}
	}

	return &ast.DefStmt{
		Pos:   pos,
		IsRec: isRec,
		Names: names,
		Exprs: exprs,
	}
}

func (parser *Parser) parseExpr(ctx exprParseCtx) ast.Expr {
	switch {
	case parser.matchAnyTokenFrom(token.Let):
		return parser.parseLetExpr(ctx)
	case parser.matchAnyTokenFrom(token.If):
		return parser.parseConditionalExpr(ctx)
	case parser.matchAnyTokenFrom(token.Fn):
		return parser.parseFunctionExpr(ctx)
	default:
		return parser.parsePrecedenceExpr(ctx, 1)
	}
}

func (parser *Parser) parseExprList(ctx exprParseCtx) []ast.Expr {
	var exprList []ast.Expr

	for gotComma := true; gotComma && !parser.reachedEndOfSource(); {
		exprList = append(exprList, parser.parseExpr(ctx))
		if gotComma = parser.matchAnyTokenFrom(token.Comma); gotComma {
			parser.consumeToken()
		}
	}

	return exprList
}

func (parser *Parser) parseLetExpr(ctx exprParseCtx) ast.Expr {
	var initNames []string
	var initExprs []ast.Expr
	var bodyExpr ast.Expr
	var isRec bool

	pos := parser.getAndConsumeNextToken().Pos // let
	if isRec = parser.matchAnyTokenFrom(token.Rec); isRec {
		parser.consumeToken()
	}
	for gotComma := true; gotComma; {
		initNames = append(initNames, parser.getAndConsumeNextToken().Raw)
		parser.expectAnyTokenFrom(token.Equal)
		parser.consumeToken()
		initExprs = append(initExprs, parser.parseExpr(ctx))

		if gotComma = parser.matchAnyTokenFrom(token.Comma); gotComma {
			parser.consumeToken()
		}
	}
	parser.expectAnyTokenFrom(token.In)
	parser.consumeToken()
	bodyExpr = parser.parseExpr(ctx)

	return &ast.LetExpr{
		Pos:       pos,
		InitNames: initNames,
		InitExprs: initExprs,
		BodyExpr:  bodyExpr,
		IsRec:     isRec,
	}
}

func (parser *Parser) parseFunctionExpr(ctx exprParseCtx) ast.Expr {
	var argNames, alphaConvertedArgNames []string
	var bodyExpr ast.Expr

	pos := parser.getAndConsumeNextToken().Pos // fun

	// Alpha convert expressions at parse time
	for _, ident := range parser.parseIdentList() {
		newIdentVal := parser.getUniqueIdentVal()
		ctx.enqueueAlphaConversion(ident, newIdentVal)
		argNames = append(argNames, ident)
		ident = newIdentVal
		alphaConvertedArgNames = append(alphaConvertedArgNames, ident)
	}

	parser.expectAnyTokenFrom(token.DashGreater)
	parser.consumeToken()
	bodyExpr = parser.parseExpr(ctx)

	for _, name := range argNames {
		delete(ctx.toBeAlphaConverted, name)
	}

	return &ast.FunExpr{
		Pos:      pos,
		ArgNames: alphaConvertedArgNames,
		BodyExpr: bodyExpr,
	}
}

func (parser *Parser) parseConditionalExpr(ctx exprParseCtx) ast.Expr {
	var condition, consequent, alternative ast.Expr

	pos := parser.getAndConsumeNextToken().Pos // if
	condition = parser.parseExpr(ctx)
	parser.expectAnyTokenFrom(token.Then)
	parser.consumeToken()
	consequent = parser.parseExpr(ctx)

	if parser.matchAnyTokenFrom(token.Else) {
		parser.consumeToken()
		alternative = parser.parseExpr(ctx)
	} else if parser.matchAnyTokenFrom(token.Elif) {
		parser.consumeToken()
		alternative = parser.parseConditionalExpr(ctx)
	}

	return &ast.ConditionalExpr{
		Pos:         pos,
		Condition:   condition,
		Consequent:  consequent,
		Alternative: alternative,
	}
}

func (parser *Parser) parseListExpr(ctx exprParseCtx) ast.Expr {
	var exprList []ast.Expr

	pos := parser.getAndConsumeNextToken().Pos // right bracket
	exprList = parser.parseExprList(ctx)
	parser.expectAnyTokenFrom(token.RightBracket)
	parser.consumeToken()

	return &ast.ListExpr{
		Pos:   pos,
		Items: exprList,
	}
}

func (parser *Parser) parsePrecedenceExpr(ctx exprParseCtx, prec int) ast.Expr {
	switch prec {
	case 1:
		return parser.parseBinaryExpr(ctx, prec, token.Xor, token.Or)
	case 2:
		return parser.parseBinaryExpr(ctx, prec, token.And)
	case 3:
		return parser.parseUnaryExpr(ctx, prec, token.Not)
	case 4:
		return parser.parseBinaryExpr(ctx, prec, token.DoubleEqual, token.ExclamationEqual, token.Greater, token.GreaterEqual, token.Less, token.LessEqual)
	case 5:
		return parser.parseBinaryExpr(ctx, prec, token.DoublePlus)
	case 6:
		return parser.parseBinaryExpr(ctx, prec, token.Plus, token.Minus)
	case 7:
		return parser.parseBinaryExpr(ctx, prec, token.Star, token.Slash, token.DoubleSlash, token.Percent, token.SlashPercent)
	case 8:
		return parser.parseUnaryExpr(ctx, prec, token.Plus, token.Minus)
	case 9:
		return parser.parseBinaryExpr(ctx, prec, token.DoubleStar)
	case 10:
		return parser.parseBinaryExpr(ctx, prec, token.DoubleExclamation)
	case 11:
		return parser.parseApplicationExpr(ctx, prec)
	default:
		return parser.parseAtomExpr(ctx)
	}
}

func (parser *Parser) parseBinaryExpr(ctx exprParseCtx, prec int, opTokKinds ...token.TokenKind) ast.Expr {
	var lhsExpr, rhsExpr, expr ast.Expr
	var operand token.Token

	lhsExpr = parser.parsePrecedenceExpr(ctx, prec+1)
	if parser.matchAnyTokenFrom(opTokKinds...) {
		operand = parser.getAndConsumeNextToken()

		if operand.Kind.IsLeftAssociative() {
			rhsExpr = parser.parsePrecedenceExpr(ctx, prec+1)
			expr = &ast.BinaryExpr{
				LhsExpr: lhsExpr,
				Operand: operand,
				RhsExpr: rhsExpr,
			}

			for parser.matchAnyTokenFrom(opTokKinds...) {
				operand = parser.getAndConsumeNextToken()
				rhsExpr = parser.parsePrecedenceExpr(ctx, prec+1)
				expr = &ast.BinaryExpr{
					LhsExpr: lhsExpr,
					Operand: operand,
					RhsExpr: rhsExpr,
				}
			}
		} else {
			rhsExpr = parser.parseBinaryExpr(ctx, prec, opTokKinds...)
			expr = &ast.BinaryExpr{
				LhsExpr: lhsExpr,
				Operand: operand,
				RhsExpr: rhsExpr,
			}
		}

		return expr
	} else {
		return lhsExpr
	}
}

func (parser *Parser) parseUnaryExpr(ctx exprParseCtx, prec int, opTokTypes ...token.TokenKind) ast.Expr {
	var expr ast.Expr
	var operand token.Token

	if parser.matchAnyTokenFrom(opTokTypes...) {
		operand = parser.getAndConsumeNextToken()
		if parser.matchAnyTokenFrom(opTokTypes...) {
			expr = parser.parseUnaryExpr(ctx, prec, opTokTypes...)
		} else {
			expr = parser.parsePrecedenceExpr(ctx, prec+1)
		}
		return &ast.UnaryExpr{
			Operand: operand,
			Expr:    expr,
		}
	} else {
		return parser.parsePrecedenceExpr(ctx, prec+1)
	}
}

func (parser *Parser) parseApplicationExpr(ctx exprParseCtx, prec int) ast.Expr {
	var function, expr ast.Expr
	var args []ast.Expr

	pos := parser.currentToken().Pos
	function = parser.parsePrecedenceExpr(ctx, prec+1)

	if parser.matchAnyTokenFrom(token.LeftParenthesis) {
		parser.consumeToken()
		args = parser.parseExprList(ctx)
		parser.expectAnyTokenFrom(token.RightParenthesis)
		parser.consumeToken()
		expr = &ast.ApplicationExpr{
			Pos:  pos,
			Fun:  function,
			Args: args,
		}
		// function application is left associative
		for parser.matchAnyTokenFrom(token.LeftParenthesis) {
			parser.consumeToken()
			args = parser.parseExprList(ctx)
			parser.expectAnyTokenFrom(token.RightParenthesis)
			parser.consumeToken()
			expr = &ast.ApplicationExpr{
				Pos:  pos,
				Fun:  expr,
				Args: args,
			}
		}
		return expr
	} else {
		return function
	}
}

func (parser *Parser) parseAtomExpr(ctx exprParseCtx) ast.Expr {
	var expr ast.Expr
	var tok token.Token

	if parser.matchAnyTokenFrom(token.LeftBracket) {
		return parser.parseListExpr(ctx)
	} else if parser.matchAnyTokenFrom(token.LeftParenthesis) {
		parser.consumeToken()
		expr = parser.parseExpr(ctx)
		parser.expectAnyTokenFrom(token.RightParenthesis)
		parser.consumeToken()
		return expr
	} else {
		tok = parser.getAndConsumeNextToken()
		if tok.Kind == token.Ident && ctx.requiresAlphaConversion(tok.Raw) {
			newVal := ctx.toBeAlphaConverted[tok.Raw]
			tok.Raw = newVal
		}
		return &ast.AtomicExpr{
			Pos:     tok.Pos,
			Val:     tok.Raw,
			ValKind: tok.Kind,
		}
	}
}

func (parser *Parser) GetProgramStmts() []ast.Stmt {
	return parser.parseStmtList()
}

func TestInputLoop() {
	newReader := bufio.NewReader(os.Stdin)
	newLexer := &lexer.Lexer{}
	newParser := &Parser{}

	for {
		fmt.Printf(">>> ")
		text, _ := newReader.ReadBytes('\n')
		newLexer.Init(text, lexer.BaseErrHandler, lexer.IgnoreSkippables)
		newParser.Init(newLexer.GetTokens(), BaseErrHandler, NoOpts)

		for _, stmt := range newParser.GetProgramStmts() {
			fmt.Println(stmt.GetJsonRepr())
		}
	}
}
