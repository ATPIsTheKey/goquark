package ast

import (
	"fmt"
	"github.com/scylladb/go-set/strset"
	"goquark/goquark/core/token"
	"strings"
)

type Node interface {
	StringRepr() string
	JsonRepr() string
}

type NodeList interface {
	Node
	Append() // todo: types
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	Variables() strset.Set
	FreeVariables() strset.Set
	BoundVariables() strset.Set
}

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
}

type (
	AtomExpr struct {
		Token token.Token
	}

	ListExpr struct {
		Items []Expr
	}

	UnaryExpr struct {
		Operand token.Token
		Expr    Expr
	}

	BinaryExpr struct {
		LhsExpr Expr
		Operand token.Token
		RhsExpr Expr
	}

	ApplicationExpr struct {
		Function  Expr
		Arguments []Expr
	}

	ConditionalExpr struct {
		Condition   Expr
		Consequent  Expr
		Alternative Expr
	}

	FunctionExpr struct {
		ArgumentNames []token.Token
		BodyExpr      Expr
	}

	LetExpr struct {
		Names     []token.Token
		InitExprs []Expr
		BodyExpr  Expr
	}

	AssignmentStmt struct {
		Names []token.Token
		Exprs []Expr
	}

	ImportStmt struct { /* todo */
	}

	ExportStmt struct { /* todo */
	}
)

func commaJoinStringReprsIdent(idents ...token.Token) string {
	var tmp []string

	for _, e := range idents {
		tmp = append(tmp, e.Val)
	}

	return strings.Join(tmp, ", ")
}

func commaJoinJsonReprsIdent(idents ...token.Token) string {
	var tmp []string

	for _, e := range idents {
		tmp = append(tmp, fmt.Sprintf(`"%s"`, e.Val))
	}

	return strings.Join(tmp, ", ")
}

func commaJoinStringReprsExpr(items ...Expr) string {
	var tmp []string

	for _, e := range items {
		tmp = append(tmp, e.StringRepr())
	}

	return strings.Join(tmp, ", ")
}

func commaJoinJsonReprsExpr(items ...Expr) string {
	var tmp []string

	for _, e := range items {
		tmp = append(tmp, e.JsonRepr())
	}

	return strings.Join(tmp, ", ")
}

func (expr AtomExpr) StringRepr() string { return expr.Token.Val }

func (expr AtomExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "AtomExpr", "value": %s}`, expr.Token.Val)
}

func (expr AtomExpr) Variables() strset.Set {
	return strset.Set{} // todo
}

func (expr AtomExpr) FreeVariables() strset.Set {
	return strset.Set{} // todo
}

func (expr AtomExpr) BoundVariables() strset.Set {
	return strset.Set{} // todo
}

////////////////////////////////////

func (expr ListExpr) StringRepr() string {
	return fmt.Sprintf("[%s]", commaJoinStringReprsExpr(expr.Items...))
}

func (expr ListExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "ListExpr", "items": [%s]}`, commaJoinJsonReprsExpr(expr.Items...))
}

func (expr ListExpr) Variables() strset.Set {
	return strset.Set{} // todo
}

func (expr ListExpr) FreeVariables() strset.Set {
	return strset.Set{} // todo
}

func (expr ListExpr) BoundVariables() strset.Set {
	return strset.Set{} // todo
}

////////////////////////////////////

func (expr UnaryExpr) StringRepr() string {
	return fmt.Sprintf("%s %s", expr.Operand.Kind.String(), expr.Expr.StringRepr())
}

func (expr UnaryExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "UnaryExpr", "expr": %s}`, expr.Expr.JsonRepr())
}

func (expr UnaryExpr) Variables() strset.Set {
	return strset.Set{} // todo
}

func (expr UnaryExpr) FreeVariables() strset.Set {
	return strset.Set{} // todo
}

func (expr UnaryExpr) BoundVariables() strset.Set {
	return strset.Set{} // todo
}

////////////////////////////////////

func (expr BinaryExpr) StringRepr() string {
	return fmt.Sprintf(
		"(%s) %s (%s)", expr.LhsExpr.StringRepr(), expr.Operand.Kind.String(), expr.RhsExpr.StringRepr())
}

func (expr BinaryExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "BinaryExpr", "operand": "%s", "lhsExpr": %s, "rhsExpr": %s}`,
		expr.Operand.Val, expr.LhsExpr.JsonRepr(), expr.RhsExpr.JsonRepr())
}

func (expr BinaryExpr) Variables() strset.Set {
	return strset.Set{} // todo
}

func (expr BinaryExpr) FreeVariables() strset.Set {
	return strset.Set{} // todo
}

func (expr BinaryExpr) BoundVariables() strset.Set {
	return strset.Set{} // todo
}

////////////////////////////////////

func (expr ApplicationExpr) StringRepr() string {
	return fmt.Sprintf("(%s)(%s)", expr.Function.StringRepr(), commaJoinStringReprsExpr(expr.Arguments...))
}

func (expr ApplicationExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "ApplicationExpr", "functions": "%s", "args": [%s]`,
		expr.Function.JsonRepr(), commaJoinJsonReprsExpr(expr.Arguments...))
}

func (expr ApplicationExpr) Variables() strset.Set {
	return strset.Set{} // todo
}

func (expr ApplicationExpr) FreeVariables() strset.Set {
	return strset.Set{} // todo
}

func (expr ApplicationExpr) BoundVariables() strset.Set {
	return strset.Set{} // todo
}

////////////////////////////////////

func (expr ConditionalExpr) StringRepr() string {
	alternative := "nil"

	if expr.Alternative != nil {
		alternative = expr.Alternative.StringRepr()
	}

	return fmt.Sprintf("if (%s) then (%s) else (%s)",
		expr.Condition.StringRepr(), expr.Consequent.StringRepr(), alternative)
}

func (expr ConditionalExpr) JsonRepr() string {
	var alternative string

	if expr.Alternative != nil {
		alternative = expr.Alternative.JsonRepr()
	}

	return fmt.Sprintf(`{"node_name": "ConditionalExpr", "condition": "%s", "consequent": %s, "alternative": "%s"`,
		expr.Condition.JsonRepr(), expr.Consequent.JsonRepr(), alternative)
}

func (expr ConditionalExpr) Variables() strset.Set {
	return strset.Set{} // todo
}

func (expr ConditionalExpr) FreeVariables() strset.Set {
	return strset.Set{} // todo
}

func (expr ConditionalExpr) BoundVariables() strset.Set {
	return strset.Set{} // todo
}

////////////////////////////////////

func (expr FunctionExpr) StringRepr() string {
	return fmt.Sprintf("fun :: %s => (%s)", commaJoinStringReprsIdent(expr.ArgumentNames...), expr.BodyExpr.StringRepr())
}

func (expr FunctionExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "FunctionExpr", "argument_names": [%s]", "body_expr": %s`,
		commaJoinJsonReprsIdent(expr.ArgumentNames...), expr.BodyExpr.JsonRepr())
}

func (expr FunctionExpr) Variables() strset.Set {
	return strset.Set{} // todo
}

func (expr FunctionExpr) FreeVariables() strset.Set {
	return strset.Set{} // todo
}

func (expr FunctionExpr) BoundVariables() strset.Set {
	return strset.Set{} // todo
}

////////////////////////////////////

func (expr LetExpr) StringRepr() string {
	return fmt.Sprintf("let %s = %s in %s", commaJoinStringReprsIdent(expr.Names...),
		commaJoinStringReprsExpr(expr.InitExprs...), expr.BodyExpr.StringRepr())
}

func (expr LetExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "LetExpr", "init_names": [%s]", "init_exprs": [%s], "body_expr": %s`,
		commaJoinJsonReprsIdent(expr.Names...), commaJoinJsonReprsExpr(expr.InitExprs...), expr.BodyExpr.JsonRepr())
}

func (expr LetExpr) Variables() strset.Set {
	return strset.Set{} // todo
}

func (expr LetExpr) FreeVariables() strset.Set {
	return strset.Set{} // todo
}

func (expr LetExpr) BoundVariables() strset.Set {
	return strset.Set{} // todo
}

////////////////////////////////////

func (stmt AssignmentStmt) StringRepr() string {
	return fmt.Sprintf("def %s = %s",
		commaJoinStringReprsIdent(stmt.Names...), commaJoinStringReprsExpr(stmt.Exprs...))
}

func (stmt AssignmentStmt) JsonRepr() string {
	return fmt.Sprintf(`"node_name": "AssignmentStmt", "names": [%s], "exprs": [%s]`,
		commaJoinJsonReprsIdent(stmt.Names...), commaJoinJsonReprsExpr(stmt.Exprs...))
}
