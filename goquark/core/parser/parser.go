package parser

import (
	"goquark/goquark/core/ast"
	"goquark/goquark/core/lexer"
	"goquark/goquark/core/token"
)

type ErrorHandler func(msg string, tok token.Token)

type Parser struct {
	src []token.Token // source
	len int           // length of source

	pos int // position of parser in src

	opts int // options passed to lex
	err  ErrorHandler
}

type parseContext struct {
	remappedVars map[string]string
}

func (parser *Parser) Init(src []token.Token, err ErrorHandler, opts int) {
	// Explicitly initialize all fields since a parser may be reused.
	parser.len = len(src)
	parser.src = src

	parser.pos = 0

	parser.opts = opts
	parser.err = err
}

func (parser *Parser) reachedEndOfSource() bool { return parser.pos >= parser.len }

// todo: update lexer methods for consistency
func (parser *Parser) getNextToken() token.Token {
	ret := parser.src[parser.pos]
	parser.consumeToken()
	return ret
}

func (parser *Parser) consumeToken() {
	parser.pos += 1
}

func (parser *Parser) matchAnyTokenFrom(kinds ...token.TokenKind) bool {
	if parser.reachedEndOfSource() {
		return false
	} else {
		for _, kind := range kinds {
			if parser.getNextToken().Kind == kind {
				return true
			}
		}
		return false
	}
}

// todo: update lexer methods for consistency
func (parser *Parser) expectAnyTokenFrom(msg string, kinds ...token.TokenKind) {
	if !parser.matchAnyTokenFrom(kinds...) {
		parser.err(msg, parser.getNextToken())
	}
}

func (parser *Parser) parseStmt() ast.Stmt {
	switch {
	case parser.matchAnyTokenFrom(token.IMPORT):
		return nil // todo
	case parser.matchAnyTokenFrom(token.EXPORT):
		return nil // todo
	case parser.matchAnyTokenFrom(token.DEF):
		return nil // todo
	case parser.matchAnyTokenFrom(token.DEFUN):
		return nil // todo
	default:
		return parser.parseExpr(&parseContext{nil})
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

func (parser *Parser) parseExpr(ctx *parseContext) ast.Expr {
	switch {
	case parser.matchAnyTokenFrom(token.LET, token.LETREC):
		return parser.parseLetExpr(ctx)
	case parser.matchAnyTokenFrom(token.IF):
		return parser.parseConditionalExpr(ctx)
	case parser.matchAnyTokenFrom(token.FUN):
		return parser.parseFunctionExpr(ctx)
	case parser.matchAnyTokenFrom(token.QUESTION_MARK_LEFT_PARENTHESIS):
		return nil // todo
	default:
		return parser.parseOperatorExpr(ctx, 1)
	}
}

func (parser *Parser) parseExprList(ctx *parseContext) []ast.Expr {
	var exprList []ast.Expr

	for gotComma := true; gotComma && !parser.reachedEndOfSource(); {
		parser.expectAnyTokenFrom("todo", token.IDENT)
		exprList = append(exprList, parser.parseExpr(ctx))
		if gotComma = parser.matchAnyTokenFrom(token.COMMA); gotComma {
			parser.consumeToken()
		}
	}

	return exprList
}

func (parser *Parser) parseLetExpr(ctx *parseContext) ast.Expr {
	var names []token.Token
	var initExprs []ast.Expr
	var bodyExpr ast.Expr

	parser.consumeToken() // token.FUN
	names = parser.parseIdentList()
	parser.expectAnyTokenFrom("todo", token.EQUAL)
	parser.consumeToken()
	initExprs = parser.parseExprList(ctx)
	parser.expectAnyTokenFrom("todo", token.IN)
	parser.consumeToken()
	bodyExpr = parser.parseExpr(ctx)

	return ast.LetExpr{Names: names, InitExprs: initExprs, BodyExpr: bodyExpr}
}

func (parser *Parser) parseFunctionExpr(ctx *parseContext) ast.Expr {
	var argNames []token.Token
	var bodyExpr ast.Expr

	// todo: unction expressions are alpha converted at parse time
	parser.expectAnyTokenFrom("todo", token.DOUBLE_COLON)
	parser.consumeToken()
	argNames = parser.parseIdentList()
	parser.expectAnyTokenFrom("todo", token.EQUAL_GREATER)
	parser.consumeToken()
	bodyExpr = parser.parseExpr(ctx)

	return ast.FunctionExpr{ArgumentNames: argNames, BodyExpr: bodyExpr}
}

func (parser *Parser) parseConditionalExpr(ctx *parseContext) ast.Expr {
	var condition, consequent, alternative ast.Expr

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

func (parser *Parser) parseListExpr(ctx *parseContext) ast.Expr {
	var exprList []ast.Expr

	exprList = parser.parseExprList(ctx)
	parser.expectAnyTokenFrom("todo", token.RIGHT_BRACKET)
	parser.consumeToken()

	return ast.ListExpr{Items: exprList}
}

func (parser *Parser) parseOperatorExpr(ctx *parseContext, prec int) ast.Expr {
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

func (parser *Parser) parseBinaryExpr(ctx *parseContext, prec int, opTokTypes ...token.TokenKind) ast.Expr {
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
			expr = ast.BinaryExpr{LhsExpr: expr, Operand: operand, RhsExpr: rhsExpr}
		}

		return expr
	} else {
		return lhsExpr
	}
}

func (parser *Parser) parseUnaryExpr(ctx *parseContext, prec int, opTokTypes ...token.TokenKind) ast.Expr {
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

func (parser *Parser) parseApplicationExpr(ctx *parseContext) ast.Expr {
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

func (parser *Parser) parseAtomExpr(ctx *parseContext) ast.Expr {
	var expr ast.Expr

	if parser.matchAnyTokenFrom(token.LEFT_PARENTHESIS) {
		parser.consumeToken()
		expr = parser.parseExpr(ctx)
		parser.expectAnyTokenFrom("todo", token.RIGHT_PARENTHESIS)
		parser.consumeToken()
		return expr
	} else {
		return ast.AtomExpr{Token: parser.getNextToken()}
	}

	// todo: error handling, alpha conversion
}
