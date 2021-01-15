package parser

import (
	"bufio"
	"fmt"
	"goquark/goquark/core/ast"
	"goquark/goquark/core/lexer"
	"goquark/goquark/core/token"
	"os"
	"strconv"
)

const (
	NOOPTS = 0
)

type ErrorHandler func(tok token.Token, msg string)

type Parser struct {
	src []token.Token
	len int
	pos int

	opts int
	err  ErrorHandler

	uniqueIdentCounter int
}

type exprParseCtx struct {
	toBeAlphaConverted map[string]string
}

func makeEmptyExprParseCtx() exprParseCtx {
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

func (parser *Parser) Init(src []token.Token, err ErrorHandler, opts int) {
	// Explicitly initialize all fields since a parser may be reused.
	parser.len = len(src)
	parser.src = src
	parser.pos = 0

	parser.opts = opts
	parser.err = err

	parser.uniqueIdentCounter = 0
}

func (parser *Parser) getUniqueIdentVal() string {
	ret := "%" + strconv.Itoa(parser.uniqueIdentCounter)
	parser.uniqueIdentCounter += 1
	return ret
}

func (parser *Parser) currentToken() token.Token { return parser.src[parser.pos] }

// todo: update lexer methods for consistency
func (parser *Parser) getNextToken() token.Token {
	ret := parser.currentToken()
	parser.consumeToken()
	return ret
}

func (parser *Parser) reachedEndOfSource() bool { return parser.currentToken().Kind == token.EOS }

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

// todo: update lexer methods for consistency
func (parser *Parser) expectAnyTokenFrom(msg string, kinds ...token.TokenKind) {
	if parser.reachedEndOfSource() || !parser.matchAnyTokenFrom(kinds...) {
		parser.err(parser.currentToken(), msg)
	}
}

func (parser *Parser) parseIdentList() []token.Token {
	var identList []token.Token

	for gotComma := true; gotComma && !parser.reachedEndOfSource(); {
		parser.expectAnyTokenFrom("todo", token.IDENT)
		identList = append(identList, parser.getNextToken())
		if gotComma = parser.matchAnyTokenFrom(token.COMMA); gotComma {
			parser.consumeToken()
		}
	}

	return identList
}

func (parser *Parser) parseStmt() ast.Stmt {
	switch {
	case parser.matchAnyTokenFrom(token.IMPORT):
		return nil // todo
	case parser.matchAnyTokenFrom(token.EXPORT):
		return nil // todo
	case parser.matchAnyTokenFrom(token.DEF):
		return parser.parseDefStmt()
	case parser.matchAnyTokenFrom(token.DEFUN):
		return nil // todo
	default:
		parse := parser.parseExpr(makeEmptyExprParseCtx())
		return parse
	}
}

func (parser *Parser) parseStmtList() []ast.Stmt {
	var stmts []ast.Stmt

	for !parser.reachedEndOfSource() {
		stmts = append(stmts, parser.parseStmt())
		parser.expectAnyTokenFrom("todo", token.SEMICOLON)
		parser.consumeToken()
	}

	return stmts
}

func (parser *Parser) parseDefStmt() ast.Stmt {
	var names []token.Token
	var exprs []ast.Expr

	parser.consumeToken() // token.DEF
	names = parser.parseIdentList()
	parser.expectAnyTokenFrom("todo", token.EQUAL)
	parser.consumeToken()
	exprs = parser.parseExprList(makeEmptyExprParseCtx())

	return ast.DefStmt{Names: names, Exprs: exprs}
}

func (parser *Parser) parseExpr(ctx exprParseCtx) ast.Expr {
	switch {
	case parser.matchAnyTokenFrom(token.LET, token.LETREC):
		return parser.parseLetExpr(ctx)
	case parser.matchAnyTokenFrom(token.IF):
		return parser.parseConditionalExpr(ctx)
	case parser.matchAnyTokenFrom(token.FUN):
		return parser.parseFunctionExpr(ctx)
	default:
		return parser.parseOperatorExpr(ctx, 1)
	}
}

func (parser *Parser) parseExprList(ctx exprParseCtx) []ast.Expr {
	var exprList []ast.Expr

	for gotComma := true; gotComma && !parser.reachedEndOfSource(); {
		exprList = append(exprList, parser.parseExpr(ctx))
		if gotComma = parser.matchAnyTokenFrom(token.COMMA); gotComma {
			parser.consumeToken()
		}
	}

	return exprList
}

func (parser *Parser) parseLetExpr(ctx exprParseCtx) ast.Expr {
	var names []token.Token
	var initExprs []ast.Expr
	var bodyExpr ast.Expr

	parser.consumeToken() // token.LET
	names = parser.parseIdentList()
	parser.expectAnyTokenFrom("todo", token.EQUAL)
	parser.consumeToken()
	initExprs = parser.parseExprList(ctx)
	parser.expectAnyTokenFrom("todo", token.IN)
	parser.consumeToken()
	bodyExpr = parser.parseExpr(ctx)

	return ast.LetExpr{InitNames: names, InitExprs: initExprs, BodyExpr: bodyExpr}
}

func (parser *Parser) parseFunctionExpr(ctx exprParseCtx) ast.Expr {
	var argNames, alphaConvertedArgNames []token.Token
	var bodyExpr ast.Expr

	parser.consumeToken() // Token.FUN
	parser.expectAnyTokenFrom("todo", token.DOUBLE_COLON)
	parser.consumeToken()

	// Remap variables within the function to uniquely named internal variables
	for _, ident := range parser.parseIdentList() {
		newIdentVal := parser.getUniqueIdentVal()
		ctx.enqueueAlphaConversion(ident.Raw, newIdentVal)
		argNames = append(argNames, ident)
		ident.Raw = newIdentVal
		alphaConvertedArgNames = append(alphaConvertedArgNames, ident)
	}

	parser.expectAnyTokenFrom("todo", token.EQUAL_GREATER)
	parser.consumeToken()
	bodyExpr = parser.parseExpr(ctx)

	for _, name := range argNames {
		delete(ctx.toBeAlphaConverted, name.Raw)
	}

	return ast.FunctionExpr{ArgumentNames: alphaConvertedArgNames, BodyExpr: bodyExpr}
}

func (parser *Parser) parseConditionalExpr(ctx exprParseCtx) ast.Expr {
	var condition, consequent, alternative ast.Expr

	parser.consumeToken() // token.IF
	condition = parser.parseExpr(ctx)
	parser.expectAnyTokenFrom("todo", token.THEN)
	parser.consumeToken()
	consequent = parser.parseExpr(ctx)

	if parser.matchAnyTokenFrom(token.ELSE) {
		parser.consumeToken()
		alternative = parser.parseExpr(ctx)
	} else if parser.matchAnyTokenFrom(token.ELIF) {
		parser.consumeToken()
		alternative = parser.parseConditionalExpr(ctx)
	}

	return ast.ConditionalExpr{Condition: condition, Consequent: consequent, Alternative: alternative}
}

func (parser *Parser) parseListExpr(ctx exprParseCtx) ast.Expr {
	var exprList []ast.Expr

	exprList = parser.parseExprList(ctx)
	parser.expectAnyTokenFrom("todo", token.RIGHT_BRACKET)
	parser.consumeToken()

	return ast.ListExpr{Items: exprList}
}

func (parser *Parser) parseOperatorExpr(ctx exprParseCtx, prec int) ast.Expr {
	switch prec {
	case 1:
		return parser.parseBinaryExpr(ctx, prec, token.XOR, token.OR)
	case 2:
		return parser.parseBinaryExpr(ctx, prec, token.AND)
	case 3:
		return parser.parseUnaryExpr(ctx, prec, token.NOT)
	case 4:
		return parser.parseBinaryExpr(ctx, prec, token.DOUBLE_EQUAL, token.EXCLAMATION_EQUAL, token.GREATER,
			token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL)
	case 5:
		return parser.parseBinaryExpr(ctx, prec, token.PLUS, token.MINUS)
	case 6:
		return parser.parseBinaryExpr(ctx, prec, token.STAR, token.SLASH, token.DOUBLE_SLASH, token.PERCENT,
			token.SLASH_PERCENT)
	case 7:
		return parser.parseUnaryExpr(ctx, prec, token.PLUS, token.MINUS)
	case 8:
		return parser.parseBinaryExpr(ctx, prec, token.DOUBLE_STAR)
	case 9:
		return parser.parseUnaryExpr(ctx, prec, token.NIL)
	case 10:
		return parser.parseBinaryExpr(ctx, prec, token.AMPERSAND)
	case 11:
		return parser.parseUnaryExpr(ctx, prec, token.HEAD, token.TAIL)
	default:
		return parser.parseApplicationExpr(ctx)
	}
}

func (parser *Parser) parseBinaryExpr(ctx exprParseCtx, prec int, opTokTypes ...token.TokenKind) ast.Expr {
	var lhsExpr, rhsExpr, expr ast.Expr
	var operand token.Token

	lhsExpr = parser.parseOperatorExpr(ctx, prec+1)
	if parser.matchAnyTokenFrom(opTokTypes...) {
		operand = parser.getNextToken()

		if operand.Kind.IsLeftAssociative() {
			rhsExpr = parser.parseOperatorExpr(ctx, prec+1)
			expr = ast.BinaryExpr{LhsExpr: lhsExpr, Operand: operand, RhsExpr: rhsExpr}

			for parser.matchAnyTokenFrom(opTokTypes...) {
				operand = parser.getNextToken()
				rhsExpr = parser.parseOperatorExpr(ctx, prec+1)
				expr = ast.BinaryExpr{LhsExpr: expr, Operand: operand, RhsExpr: rhsExpr}
			}
		} else {
			rhsExpr = parser.parseBinaryExpr(ctx, prec, opTokTypes...)
			expr = ast.BinaryExpr{LhsExpr: lhsExpr, Operand: operand, RhsExpr: rhsExpr}
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
		operand = parser.getNextToken()
		if parser.matchAnyTokenFrom(opTokTypes...) {
			expr = parser.parseUnaryExpr(ctx, prec, opTokTypes...)
		} else {
			expr = parser.parseOperatorExpr(ctx, prec+1)
		}
		return ast.UnaryExpr{Operand: operand, Expr: expr}
	} else {
		return parser.parseOperatorExpr(ctx, prec+1)
	}
}

func (parser *Parser) parseApplicationExpr(ctx exprParseCtx) ast.Expr {
	var function, expr ast.Expr
	var args []ast.Expr

	function = parser.parseAtomExpr(ctx)
	// function composition
	if parser.matchAnyTokenFrom(token.PERIOD) {
		parser.consumeToken()
		return ast.ApplicationExpr{Function: function, Arguments: parser.parseExprList(ctx)}
	} else if parser.matchAnyTokenFrom(token.LEFT_PARENTHESIS) {
		parser.consumeToken()
		args = parser.parseExprList(ctx)
		parser.expectAnyTokenFrom("todo", token.RIGHT_PARENTHESIS)
		parser.consumeToken()
		expr = ast.ApplicationExpr{Function: function, Arguments: args}
		// function application is left associative
		for parser.matchAnyTokenFrom(token.LEFT_PARENTHESIS) {
			parser.consumeToken()
			args = parser.parseExprList(ctx)
			parser.expectAnyTokenFrom("todo", token.RIGHT_PARENTHESIS)
			parser.consumeToken()
			expr = ast.ApplicationExpr{Function: expr, Arguments: args}
		}
		return expr
	} else {
		return function
	}
}

func (parser *Parser) parseAtomExpr(ctx exprParseCtx) ast.Expr {
	var expr ast.Expr
	var tok token.Token

	if parser.matchAnyTokenFrom(token.LEFT_PARENTHESIS) {
		parser.consumeToken()
		expr = parser.parseExpr(ctx)
		parser.expectAnyTokenFrom("todo", token.RIGHT_PARENTHESIS)
		parser.consumeToken()
		return expr
	} else {
		tok = parser.getNextToken()
		if tok.Kind == token.IDENT && ctx.requiresAlphaConversion(tok.Raw) {
			newVal := ctx.toBeAlphaConverted[tok.Raw]
			tok.Raw = newVal
		}
		return ast.AtomicExpr{Token: tok}
	}
}

func (parser *Parser) GetAst() []ast.Stmt {
	return parser.parseStmtList()
}

func TestInputLoop() {
	reader := bufio.NewReader(os.Stdin)
	lexer_ := &lexer.Lexer{}
	parser := &Parser{}

	lexerErrHandler := func(tokPos token.TokenPos, msg string) {
		fmt.Printf("LexerError at (line: %d, col: %d): %s\n", tokPos.Line, tokPos.Column, msg)
		os.Exit(-1)
	}

	parserErrHandler := func(tok token.Token, msg string) {
		fmt.Printf("ParserError at (line: %d, col: %d): %s\n", tok.FPos.Line, tok.FPos.Column, msg)
		os.Exit(-1)
	}

	for {
		fmt.Printf(">>> ")
		text, _ := reader.ReadBytes('\n')
		lexer_.Init(text, lexerErrHandler, lexer.IGNORE_SKIPPABLES)
		parser.Init(lexer_.GetTokens(), parserErrHandler, 0)

		for _, stmt := range parser.GetAst() {
			fmt.Println(stmt.JsonRepr())
		}
	}
}
