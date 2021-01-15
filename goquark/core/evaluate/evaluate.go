package evaluate

import (
	"bufio"
	"fmt"
	"goquark/goquark/core/ast"
	"goquark/goquark/core/lexer"
	"goquark/goquark/core/parser"
	"goquark/goquark/core/runtime"
	"goquark/goquark/core/token"
	"os"
	"strconv"
)

func EvalAtomicExpr(expr ast.AtomicExpr, closure Closure, argStack ArgStack) runtime.Object {
	switch expr.Token.Kind {
	case token.BOOLEAN:
		if expr.Token.Raw == "True" {
			return runtime.NewBooleanObject(true)
		} else {
			return runtime.NewBooleanObject(false)
		}
	case token.INTEGER:
		val, _ := strconv.ParseInt(expr.Token.Raw, 10, 64)
		return runtime.NewIntObject(val)
	case token.REAL:
		val, _ := strconv.ParseFloat(expr.Token.Raw, 64)
		return runtime.NewRealObject(val)
	case token.COMPLEX:
		// todo: implement own function to parse complex literals
		val, _ := strconv.ParseComplex(expr.Token.Raw, 128)
		return runtime.NewComplexObject(val)
	case token.IDENT:
		entry, ok := closure.LookupNode(expr.Token.Raw)
		if ok {
			if entry.EvalClosure == nil {
				return Eval(entry.Node, closure, argStack)
			} else {
				return Eval(entry.Node, *entry.EvalClosure, argStack)
			}
		} else {
			// todo: error handling
			return nil
		}
	default:
		// todo: error handling
		return nil
	}
}

func EvalUnaryExpr(expr ast.UnaryExpr, closure Closure, argStack ArgStack) runtime.Object {
	Eval := Eval(expr.Expr, closure, argStack)

	switch expr.Operand.Kind {
	case token.NOT:
		return Eval.LNot()
	case token.NIL:
		return Eval.Nil()
	case token.TAIL:
		return Eval.Tail()
	case token.HEAD:
		return Eval.Head()
	default:
		// todo: error handling
		return nil
	}
}

func EvalBinaryExpr(expr ast.BinaryExpr, closure Closure, argStack ArgStack) runtime.Object {
	lhsEval, rhsEval := Eval(expr.LhsExpr, closure, argStack), Eval(expr.RhsExpr, closure, argStack)

	switch expr.Operand.Kind {
	case token.DOUBLE_EQUAL:
		return lhsEval.Equal(rhsEval)
	case token.EXCLAMATION_EQUAL:
		return lhsEval.NotEqual(rhsEval)
	case token.GREATER:
		return lhsEval.Greater(rhsEval)
	case token.GREATER_EQUAL:
		return lhsEval.GreaterEqual(rhsEval)
	case token.LESS:
		return lhsEval.Less(rhsEval)
	case token.LESS_EQUAL:
		return lhsEval.Less(rhsEval)
	case token.AND:
		return lhsEval.LAnd(rhsEval)
	case token.OR:
		return lhsEval.LOr(rhsEval)
	case token.XOR:
		return lhsEval.LXor(rhsEval)
	case token.PLUS:
		return lhsEval.Add(rhsEval)
	case token.MINUS:
		return lhsEval.Sub(rhsEval)
	case token.PERCENT:
		return lhsEval.Mod(rhsEval)
	case token.STAR:
		return lhsEval.Mul(rhsEval)
	case token.SLASH:
		return lhsEval.Div(rhsEval)
	case token.DOUBLE_SLASH:
		return lhsEval.FloorDiv(rhsEval)
	case token.DOUBLE_STAR:
		return lhsEval.Pow(rhsEval)
	}
	// todo: error handling
	return nil
}

func EvalConditionalExpr(expr ast.ConditionalExpr, closure Closure, argStack ArgStack) runtime.Object {
	conditionEval := Eval(expr.Condition, closure, argStack)

	if conditionEval.AsBool().GetVal().Boolean == true {
		return Eval(expr.Consequent, closure, argStack)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative, closure, argStack)
	} else {
		return runtime.NewNilObject()
	}
}

func EvalFunctionExpr(expr ast.FunctionExpr, closure Closure, argStack ArgStack) runtime.Object {
	if !argStack.IsEmpty() {
		for _, name := range expr.ArgumentNames {
			arg, _ := argStack.Pop()
			closure.PushNode(name.Raw, arg, closure.Parent)
		}
	}

	return Eval(expr.BodyExpr, closure, argStack)
}

func EvalApplicationExpr(expr ast.ApplicationExpr, closure Closure, argStack ArgStack) runtime.Object {
	argStack.PushFromSliceReverse(&expr.Arguments)
	return Eval(expr.Function, closure.MakeChild(), argStack)
}

func EvalLetExpr(expr ast.LetExpr, closure Closure, argStack ArgStack) runtime.Object {
	closureChild := closure.MakeChild()

	for i, l := 0, len(expr.InitNames); i < l; i++ {
		name := expr.InitNames[i].Raw
		_, namePresent := closure.EntryMap[name]

		if namePresent {
			closureChild.PushNode(expr.InitNames[i].Raw, expr.InitExprs[i], closureChild.Parent)
		} else {
			closureChild.PushNode(expr.InitNames[i].Raw, expr.InitExprs[i], nil)
		}
	}

	return Eval(expr.BodyExpr, closureChild, argStack)
}

func EvalDefStmt(stmt ast.DefStmt, closure Closure, _ ArgStack) runtime.Object {
	evalClosure := closure.Copy()

	for i, l := 0, len(stmt.Names); i < l; i++ {
		closure.PushNode(stmt.Names[i].Raw, stmt.Exprs[i], &evalClosure)
	}

	return runtime.NewNilObject()
}

func Eval(node ast.Node, closure Closure, argStack ArgStack) runtime.Object {
	switch node.(type) {
	case ast.ImportStmt:
		return nil // todo
	case ast.ExportStmt:
		return nil // todo
	case ast.DefStmt:
		return EvalDefStmt(node.(ast.DefStmt), closure, argStack)
	case ast.ConditionalExpr:
		return EvalConditionalExpr(node.(ast.ConditionalExpr), closure, argStack)
	case ast.FunctionExpr:
		return EvalFunctionExpr(node.(ast.FunctionExpr), closure, argStack)
	case ast.ApplicationExpr:
		return EvalApplicationExpr(node.(ast.ApplicationExpr), closure, argStack)
	case ast.BinaryExpr:
		return EvalBinaryExpr(node.(ast.BinaryExpr), closure, argStack)
	case ast.UnaryExpr:
		return EvalUnaryExpr(node.(ast.UnaryExpr), closure, argStack)
	case ast.LetExpr:
		return EvalLetExpr(node.(ast.LetExpr), closure, argStack)
	case ast.ListExpr:
		return nil // todo
	case ast.CellExpr:
		return nil // todo
	case ast.AtomicExpr:
		return EvalAtomicExpr(node.(ast.AtomicExpr), closure, argStack)
	default:
		return nil // todo: runtime error
	}
}

func TestInputLoop() {
	reader := bufio.NewReader(os.Stdin)
	lexer_ := &lexer.Lexer{}
	parser_ := &parser.Parser{}

	lexerErrHandler := func(tokPos token.TokenPos, msg string) {
		fmt.Printf("LexerError at (line: %d, col: %d): %s\n", tokPos.Line, tokPos.Column, msg)
		os.Exit(-1)
	}

	parserErrHandler := func(tok token.Token, msg string) {
		fmt.Printf("ParserError at (line: %d, col: %d): %s\n", tok.FPos.Line, tok.FPos.Column, msg)
		os.Exit(-1)
	}

	globalClosure := MakeVoidClosure()
	var argStack ArgStack

	for {
		fmt.Printf(">>> ")
		text, _ := reader.ReadBytes('\n')
		lexer_.Init(text, lexerErrHandler, lexer.IGNORE_SKIPPABLES)
		parser_.Init(lexer_.GetTokens(), parserErrHandler, 0)

		for _, stmt := range parser_.GetAst() {
			fmt.Printf("%+v\n", Eval(stmt, globalClosure, argStack).Repr())
			//fmt.Printf("%+v\n", globalClosure)
		}

	}
}
