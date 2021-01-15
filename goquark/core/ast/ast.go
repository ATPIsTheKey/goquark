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
	AtomicExpr struct {
		Token token.Token
	}

	CellExpr struct {
		First  Expr
		Second Expr
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
		InitNames []token.Token
		InitExprs []Expr
		BodyExpr  Expr
	}

	DefStmt struct {
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
		tmp = append(tmp, e.Raw)
	}

	return strings.Join(tmp, ", ")
}

func commaJoinJsonReprsIdent(idents ...token.Token) string {
	var tmp []string

	for _, e := range idents {
		tmp = append(tmp, fmt.Sprintf(`"%s"`, e.Raw))
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

func (expr AtomicExpr) StringRepr() string { return expr.Token.Raw }

func (expr AtomicExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "AtomicExpr", "value": %#v}`, expr.Token.Raw)
}

func (expr AtomicExpr) Variables() strset.Set {
	if expr.Token.Kind == token.IDENT {
		return *strset.New(expr.Token.Raw)
	} else {
		return *strset.New()
	}
}

func (expr AtomicExpr) FreeVariables() strset.Set {
	return expr.Variables()
}

func (expr AtomicExpr) BoundVariables() strset.Set {
	return strset.Set{} // empty set
}

////////////////////////////////////

func (expr ListExpr) StringRepr() string {
	return fmt.Sprintf("[%s]", commaJoinStringReprsExpr(expr.Items...))
}

func (expr ListExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "ListExpr", "items": [%s]}`, commaJoinJsonReprsExpr(expr.Items...))
}

func (expr ListExpr) Variables() strset.Set {
	set := strset.New()

	for _, item := range expr.Items {
		itemSet := item.Variables()
		set.Merge(&itemSet)
	}

	return *set
}

func (expr ListExpr) FreeVariables() strset.Set {
	set := strset.New()

	for _, item := range expr.Items {
		itemSet := item.FreeVariables()
		set.Merge(&itemSet)
	}

	return *set
}

func (expr ListExpr) BoundVariables() strset.Set {
	set := strset.New()

	for _, item := range expr.Items {
		itemSet := item.BoundVariables()
		set.Merge(&itemSet)
	}

	return *set
}

////////////////////////////////////

func (expr CellExpr) StringRepr() string {
	return fmt.Sprintf("<%s, %s>", expr.First.StringRepr(), expr.Second.StringRepr())
}

func (expr CellExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "CellExpr", "first": %s, "second": %s}`, expr.First.StringRepr(), expr.Second.StringRepr())
}

func (expr CellExpr) Variables() strset.Set { return strset.Set{} /* todo */ }

func (expr CellExpr) FreeVariables() strset.Set { return strset.Set{} /* todo */ }

func (expr CellExpr) BoundVariables() strset.Set { return strset.Set{} /* todo */ }

////////////////////////////////////

func (expr UnaryExpr) StringRepr() string {
	return fmt.Sprintf("%s %s", expr.Operand.Kind.String(), expr.Expr.StringRepr())
}

func (expr UnaryExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "UnaryExpr", "expr": %s}`, expr.Expr.JsonRepr())
}

func (expr UnaryExpr) Variables() strset.Set {
	return expr.Expr.Variables()
}

func (expr UnaryExpr) FreeVariables() strset.Set {
	return expr.Expr.FreeVariables()
}

func (expr UnaryExpr) BoundVariables() strset.Set {
	return expr.Expr.BoundVariables()
}

////////////////////////////////////

func (expr BinaryExpr) StringRepr() string {
	return fmt.Sprintf(
		"(%s) %s (%s)", expr.LhsExpr.StringRepr(), expr.Operand.Kind.String(), expr.RhsExpr.StringRepr())
}

func (expr BinaryExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "BinaryExpr", "operand": "%s", "lhsExpr": %s, "rhsExpr": %s}`,
		expr.Operand.Raw, expr.LhsExpr.JsonRepr(), expr.RhsExpr.JsonRepr())
}

func (expr BinaryExpr) Variables() strset.Set {
	lhsVars, rhsVars := expr.LhsExpr.Variables(), expr.RhsExpr.Variables()
	vars := strset.Union(&lhsVars, &rhsVars)
	return *vars
}

func (expr BinaryExpr) FreeVariables() strset.Set {
	lhsSet, rhsSet := expr.LhsExpr.FreeVariables(), expr.RhsExpr.FreeVariables()
	vars := strset.Union(&lhsSet, &rhsSet)
	return *vars
}

func (expr BinaryExpr) BoundVariables() strset.Set {
	lhsSet, rhsSet := expr.LhsExpr.BoundVariables(), expr.RhsExpr.BoundVariables()
	vars := strset.Union(&lhsSet, &rhsSet)
	return *vars
}

////////////////////////////////////

func (expr ApplicationExpr) StringRepr() string {
	return fmt.Sprintf("(%s)(%s)", expr.Function.StringRepr(), commaJoinStringReprsExpr(expr.Arguments...))
}

func (expr ApplicationExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "ApplicationExpr", "functions": %s, "args": [%s]}`,
		expr.Function.JsonRepr(), commaJoinJsonReprsExpr(expr.Arguments...))
}

func (expr ApplicationExpr) Variables() strset.Set {
	set := strset.New()

	for _, arg := range expr.Arguments {
		itemSet := arg.Variables()
		set.Merge(&itemSet)
	}

	return *set
}

func (expr ApplicationExpr) FreeVariables() strset.Set {
	set := strset.New()

	for _, arg := range expr.Arguments {
		itemSet := arg.FreeVariables()
		set.Merge(&itemSet)
	}

	return *set
}

func (expr ApplicationExpr) BoundVariables() strset.Set {
	set := strset.New()

	for _, arg := range expr.Arguments {
		itemSet := arg.BoundVariables()
		set.Merge(&itemSet)
	}

	return *set
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
	} else {
		alternative = `"nil"`
	}

	return fmt.Sprintf(`{"node_name": "ConditionalExpr", "condition": %s, "consequent": %s, "alternative": %s}`,
		expr.Condition.JsonRepr(), expr.Consequent.JsonRepr(), alternative)
}

func (expr ConditionalExpr) Variables() strset.Set {
	conditionSet, consequentSet := expr.Condition.Variables(), expr.Consequent.Variables()

	if expr.Alternative != nil {
		alternativeSet := expr.Alternative.Variables()
		set := strset.Union(&conditionSet, &consequentSet, &alternativeSet)
		return *set
	} else {
		set := strset.Union(&conditionSet, &consequentSet)
		return *set
	}
}

func (expr ConditionalExpr) FreeVariables() strset.Set {
	conditionSet, consequentSet := expr.Condition.FreeVariables(), expr.Consequent.FreeVariables()

	if expr.Alternative != nil {
		alternativeSet := expr.Alternative.FreeVariables()
		set := strset.Union(&conditionSet, &consequentSet, &alternativeSet)
		return *set
	} else {
		set := strset.Union(&conditionSet, &consequentSet)
		return *set
	}
}

func (expr ConditionalExpr) BoundVariables() strset.Set {
	conditionSet, consequentSet := expr.Condition.BoundVariables(), expr.Consequent.BoundVariables()

	if expr.Alternative != nil {
		alternativeSet := expr.Alternative.BoundVariables()
		set := strset.Union(&conditionSet, &consequentSet, &alternativeSet)
		return *set
	} else {
		set := strset.Union(&conditionSet, &consequentSet)
		return *set
	}
}

////////////////////////////////////

func (expr FunctionExpr) StringRepr() string {
	return fmt.Sprintf("fun :: %s => (%s)", commaJoinStringReprsIdent(expr.ArgumentNames...), expr.BodyExpr.StringRepr())
}

func (expr FunctionExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "FunctionExpr", "argument_names": [%s], "body_expr": %s}`,
		commaJoinJsonReprsIdent(expr.ArgumentNames...), expr.BodyExpr.JsonRepr())
}

func (expr FunctionExpr) Variables() strset.Set {
	return expr.BodyExpr.Variables()
}

func (expr FunctionExpr) FreeVariables() strset.Set {
	bodySet := expr.BodyExpr.FreeVariables()

	argNamesSet := strset.New()
	for _, tok := range expr.ArgumentNames {
		argNamesSet.Add(tok.Raw)
	}

	set := strset.Difference(&bodySet, strset.Union(argNamesSet, &bodySet))
	return *set
}

func (expr FunctionExpr) BoundVariables() strset.Set {
	bodySet := expr.BodyExpr.FreeVariables()

	argNamesSet := strset.New()
	for _, tok := range expr.ArgumentNames {
		argNamesSet.Add(tok.Raw)
	}

	set := strset.Intersection(argNamesSet, &bodySet)
	return *set
}

////////////////////////////////////

func (expr LetExpr) StringRepr() string {
	return fmt.Sprintf("let %s = %s in %s", commaJoinStringReprsIdent(expr.InitNames...),
		commaJoinStringReprsExpr(expr.InitExprs...), expr.BodyExpr.StringRepr())
}

func (expr LetExpr) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "LetExpr", "init_names": [%s], "init_exprs": [%s], "body_expr": %s}`,
		commaJoinJsonReprsIdent(expr.InitNames...), commaJoinJsonReprsExpr(expr.InitExprs...), expr.BodyExpr.JsonRepr())
}

func (expr LetExpr) Variables() strset.Set {
	bodyExprSet := expr.BodyExpr.Variables()

	initExprsSet := strset.New()
	for _, initExpr := range expr.InitExprs {
		itemSet := initExpr.Variables()
		initExprsSet.Merge(&itemSet)
	}
	set := strset.Union(&bodyExprSet, initExprsSet)

	return *set
}

func (expr LetExpr) FreeVariables() strset.Set {
	bodyExprSet := expr.BodyExpr.FreeVariables()

	initNamesSet := strset.Set{}
	for _, name := range expr.InitNames {
		initNamesSet.Add(name.Raw)
	}

	initExprsSet := strset.New()
	for _, initExpr := range expr.InitExprs {
		itemSet := initExpr.Variables()
		initExprsSet.Merge(&itemSet)
	}

	set := strset.Difference(strset.Union(&bodyExprSet, initExprsSet), &initNamesSet)

	return *set
}

func (expr LetExpr) BoundVariables() strset.Set {
	bodyExprSet := expr.BodyExpr.BoundVariables()

	initNamesSet := strset.New()
	for _, name := range expr.InitNames {
		initNamesSet.Add(name.Raw)
	}

	initExprsSet := strset.New()
	for _, initExpr := range expr.InitExprs {
		itemSet := initExpr.Variables()
		initExprsSet.Merge(&itemSet)
	}

	set := strset.Intersection(strset.Union(&bodyExprSet, initExprsSet), initNamesSet)

	return *set
}

////////////////////////////////////

func (stmt DefStmt) StringRepr() string {
	return fmt.Sprintf("def %s = %s",
		commaJoinStringReprsIdent(stmt.Names...), commaJoinStringReprsExpr(stmt.Exprs...))
}

func (stmt DefStmt) JsonRepr() string {
	return fmt.Sprintf(`{"node_name": "DefStmt", "names": [%s], "exprs": [%s]}`,
		commaJoinJsonReprsIdent(stmt.Names...), commaJoinJsonReprsExpr(stmt.Exprs...))
}

func (stmt ImportStmt) StringRepr() string { return "" /* todo */ }

func (stmt ImportStmt) JsonRepr() string { return "" /* todo */ }

func (stmt ExportStmt) StringRepr() string { return "" /* todo */ }

func (stmt ExportStmt) JsonRepr() string { return "" /* todo */ }
