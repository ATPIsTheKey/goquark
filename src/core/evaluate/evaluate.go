package evaluate

import (
	"bufio"
	"fmt"
	"goquark/src/core/ast"
	"goquark/src/core/lexer"
	"goquark/src/core/parser"
	"goquark/src/core/runtime"
	"goquark/src/core/token"
	"os"
	"strconv"
)

func EvalAtomicExpr(expr *ast.AtomicExpr, frame *runtime.Frame) runtime.Object {
	return runtime.NewThunkObject(
		func() runtime.Object {
			switch expr.ValKind {
			case token.Boolean:
				if expr.Val == "True" {
					return runtime.NewBoolObject(true)
				} else {
					return runtime.NewBoolObject(false)
				}
			case token.Integer:
				val, _ := strconv.ParseInt(expr.Val, 10, 64)
				return runtime.NewIntObject(val)
			case token.Real:
				val, _ := strconv.ParseFloat(expr.Val, 64)
				return runtime.NewRealObject(val)
			case token.Complex:
				val, _ := strconv.ParseComplex(expr.Val, 128)
				return runtime.NewComplexObject(val)
			case token.Ident:
				if obj, ok := frame.GetFromEnv(expr.Val); ok {
					return obj
				} else {
					return runtime.NewPoisonObject(expr.Val+" not defined in scope\n", frame)
				}
			default:
				panic("Invalid atomic expr kind " + expr.ValKind.String()) // todo
				return nil
			}
		},
	)
}

func EvalListExpr(expr *ast.ListExpr, frame *runtime.Frame) runtime.Object {
	return runtime.NewThunkObject(
		func() runtime.Object {
			var objs []runtime.Object
			for _, expr := range expr.Items {
				objs = append(objs, EvalExpr(expr, frame))
			}

			return runtime.NewListObject(objs...)
		},
	)
}

func EvalUnaryExpr(expr *ast.UnaryExpr, frame *runtime.Frame) runtime.Object {
	return runtime.NewThunkObject(
		func() runtime.Object {
			res := EvalExpr(expr.Expr, frame)

			switch expr.Operand.Kind {
			case token.Not:
				return res.LNot(frame.New(expr.Pos.Fmt() + ": __LogicalNot(_)"))
			default:
				panic("Unsupported unary operand " + expr.Operand.Kind.String())
				return nil
			}
		},
	)
}

func EvalBinaryExpr(expr *ast.BinaryExpr, frame *runtime.Frame) runtime.Object {
	return runtime.NewThunkObject(
		func() runtime.Object {
			lhsRes, rhsRes := EvalExpr(expr.LhsExpr, frame), EvalExpr(expr.RhsExpr, frame)

			switch expr.Operand.Kind {
			case token.DoubleEqual:
				return lhsRes.Equal(rhsRes, frame.New(expr.Pos.Fmt()+": __Equal(_, _)"))
			case token.ExclamationEqual:
				return lhsRes.NotEqual(rhsRes, frame.New(expr.Pos.Fmt()+": __NotEqual(_, _)"))
			case token.Greater:
				return lhsRes.Greater(rhsRes, frame.New(expr.Pos.Fmt()+": __Greater(_, _)"))
			case token.GreaterEqual:
				return lhsRes.GreaterEqual(rhsRes, frame.New(expr.Pos.Fmt()+": __GreaterEqual(_, _)"))
			case token.Less:
				return lhsRes.Less(rhsRes, frame.New(expr.Pos.Fmt()+": __Less(_, _)"))
			case token.LessEqual:
				return lhsRes.LessEqual(rhsRes, frame.New(expr.Pos.Fmt()+": __LessEqual(_, _)"))
			case token.And:
				return lhsRes.LAnd(rhsRes, frame.New(expr.Pos.Fmt()+": __LogicalAnd(_, _)"))
			case token.Or:
				return lhsRes.LOr(rhsRes, frame.New(expr.Pos.Fmt()+": __LogicalOr(_, _)"))
			case token.Xor:
				return lhsRes.LXor(rhsRes, frame.New(expr.Pos.Fmt()+": __LogicalXor(_, _)"))
			case token.Plus:
				return lhsRes.Add(rhsRes, frame.New(expr.Pos.Fmt()+": __Add(_, _)"))
			case token.Minus:
				return lhsRes.Sub(rhsRes, frame.New(expr.Pos.Fmt()+": __Sub(_, _)"))
			case token.Percent:
				return lhsRes.Mod(rhsRes, frame.New(expr.Pos.Fmt()+": __Mod(_, _)"))
			case token.Star:
				return lhsRes.Mul(rhsRes, frame.New(expr.Pos.Fmt()+": __Mul(_, _)"))
			case token.Slash:
				return lhsRes.Div(rhsRes, frame.New(expr.Pos.Fmt()+": __Div(_, _)"))
			case token.DoubleSlash:
				return lhsRes.FloorDiv(rhsRes, frame.New(expr.Pos.Fmt()+": __FloorDiv(_, _)"))
			case token.DoubleStar:
				return lhsRes.Pow(rhsRes, frame.New(expr.Pos.Fmt()+": __Pow(_, _)"))
			case token.DoubleExclamation:
				return lhsRes.GetItem(rhsRes, frame.New(expr.Pos.Fmt()+": __GetItem(_, _)"))
			case token.DoublePlus:
				return lhsRes.Concatenate(rhsRes, frame.New(expr.Pos.Fmt()+": __Concatenate(_, _)"))
			default:
				panic("Unsupported binary operand " + expr.Operand.Kind.String())
				return nil
			}
		},
	)
}

func EvalConditionalExpr(expr *ast.ConditionalExpr, frame *runtime.Frame) runtime.Object {
	return runtime.NewThunkObject(
		func() runtime.Object {
			conditionRes := EvalExpr(expr.Condition, frame)

			if conditionRes.AsBool(frame.New(expr.Pos.Fmt()+": __AsBool(_)")).GetVal().Bool == true {
				return EvalExpr(expr.Consequent, frame)
			} else if expr.Alternative != nil {
				return EvalExpr(expr.Alternative, frame)
			} else {
				return runtime.NewNilObject()
			}
		},
	)
}

func EvalFunExpr(expr *ast.FunExpr, frame *runtime.Frame) runtime.Object {
	return runtime.NewFunObject(
		func(callFrame *runtime.Frame) runtime.Object {
			evalFrame := frame.Copy()

			if callFrame.ArgStack.Size() < expr.Arity() {
				return runtime.NewFunObject(
					func(callFrame2 *runtime.Frame) runtime.Object {
						callFrame := callFrame.Copy()
						callFrame.ArgStack = callFrame2.ArgStack.GetMerged(callFrame.ArgStack)
						return EvalFunExpr(expr, evalFrame).Apply(callFrame)
					},
				)
			} else {
				for _, name := range expr.ArgNames {
					argObj, _ := callFrame.ArgStack.Pop()
					evalFrame.PushToEnv(name, argObj)
				}
				return EvalExpr(expr.BodyExpr, evalFrame).Apply(callFrame) // apply in case body expression is also a fun expression
			}
		},
	)
}

func EvalApplicationExpr(expr *ast.ApplicationExpr, frame *runtime.Frame) runtime.Object {
	return runtime.NewThunkObject(
		func() runtime.Object {
			funRes := EvalExpr(expr.Fun, frame)
			last := len(expr.Args) - 1

			for i := range expr.Args {
				frame.ArgStack.Push(EvalExpr(expr.Args[last-i], frame))
			}

			return funRes.Apply(frame)
		},
	)
}

func EvalLetExpr(expr *ast.LetExpr, frame *runtime.Frame) runtime.Object {
	evalFrame := frame.New(expr.Pos.Fmt() + ": __NewClosure()")

	if expr.IsRec {
		for i := range expr.InitNames {
			initName, initExpr := expr.InitNames[i], expr.InitExprs[i]
			evalFrame.PushToEnv(initName, EvalExpr(initExpr, evalFrame))
		}
	} else {
		for i := range expr.InitNames {
			initName, initExpr := expr.InitNames[i], expr.InitExprs[i]
			evalFrame.PushToEnv(initName, EvalExpr(initExpr, evalFrame.Copy()))
		}
	}

	return EvalExpr(expr.BodyExpr, evalFrame)
}

func EvalDefStmt(stmt *ast.DefStmt, frame *runtime.Frame) runtime.Object {
	if stmt.IsRec {
		for i := range stmt.Names {
			name, expr := stmt.Names[i], stmt.Exprs[i]
			frame.PushToEnv(name, EvalExpr(expr, frame))
		}
	} else {
		for i := range stmt.Names {
			name, expr := stmt.Names[i], stmt.Exprs[i]
			frame.PushToEnv(name, EvalExpr(expr, frame.Copy()))
		}
	}

	return runtime.NewNilObject()
}

func EvalExpr(expr ast.Expr, frame *runtime.Frame) runtime.Object {
	switch (expr).(type) {
	case *ast.ConditionalExpr:
		return EvalConditionalExpr(expr.(*ast.ConditionalExpr), frame)
	case *ast.FunExpr:
		return EvalFunExpr(expr.(*ast.FunExpr), frame)
	case *ast.ApplicationExpr:
		{
			pos := expr.(*ast.ApplicationExpr).Pos
			// todo: think of way to repr lambda functions
			frameDescription := pos.Fmt() + ": __FunctionApplication" // todo: think of way to repr lambda functions
			return EvalApplicationExpr(expr.(*ast.ApplicationExpr), frame.New(frameDescription))
		}
	case *ast.BinaryExpr:
		return EvalBinaryExpr(expr.(*ast.BinaryExpr), frame)
	case *ast.UnaryExpr:
		return EvalUnaryExpr(expr.(*ast.UnaryExpr), frame)
	case *ast.LetExpr:
		return EvalLetExpr(expr.(*ast.LetExpr), frame)
	case *ast.ListExpr:
		return EvalListExpr(expr.(*ast.ListExpr), frame)
	case *ast.AtomicExpr:
		return EvalAtomicExpr(expr.(*ast.AtomicExpr), frame)
	default:
		panic("Unsupported expr type")
		return nil
	}
}

func EvalStmt(stmt ast.Stmt, frame *runtime.Frame) (ret runtime.Object) {
	defer func() {
		if traceback := recover(); traceback != nil {
			fmt.Println(traceback)
			ret = runtime.NewNilObject()
		}
	}()

	switch stmt.(type) {
	case *ast.DefStmt:
		return EvalDefStmt(stmt.(*ast.DefStmt), frame)
	default:
		if expr, ok := interface{}(stmt).(ast.Expr); ok { // check if node implements expr interface
			return EvalExpr(expr, frame)
		} else {
			panic("Unsupported stmt type")
			return nil
		}
	}
}

func EvalMain(stmts ...ast.Stmt) runtime.Object {
	rootFrame := runtime.NewRootFrame()

	for _, stmt := range stmts {
		_ = EvalStmt(stmt, rootFrame)
	}

	if res, ok := rootFrame.GetFromEnv("Main"); ok {
		return res
	} else {
		return runtime.NewNilObject()
	}
}

func TestInputLoop() {
	newReader := bufio.NewReader(os.Stdin)
	newLexer := &lexer.Lexer{}
	newParser := &parser.Parser{}
	newRootFrame := runtime.NewRootFrame()

	for {
		fmt.Printf(">>> ")
		input, _ := newReader.ReadBytes('\n')
		newLexer.Init(input, lexer.BaseErrHandler, lexer.IgnoreSkippables)
		newParser.Init(newLexer.GetTokens(), parser.BaseErrHandler, parser.NoOpts)

		for _, stmt := range newParser.GetProgramStmts() {
			if res := EvalStmt(stmt, newRootFrame); res.Inspect().Type == runtime.PoisonType {
				fmt.Printf("RuntimeError: %s \nTraceback (most recent call last):\n", res.(runtime.ThunkObject).GetActualObject().(runtime.PoisonObject).ErrorMsg)
				fmt.Println(res.(runtime.ThunkObject).GetActualObject().(runtime.PoisonObject).ReleaseFrame.BuildTraceback())
			} else {
				fmt.Println(res.Repr())
			}
		}
	}
}
