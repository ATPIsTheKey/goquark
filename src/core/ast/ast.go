package ast

import (
	"fmt"
	"github.com/scylladb/go-set/strset"
	"goquark/src/core/token"
	"goquark/src/core/utils"
	"strings"
)

type Node interface {
	GetPos() utils.SourceIndex
	GetStringRepr() string
	GetJsonRepr() string
}

type NodeList interface {
	Node
	Append() // todo: objects
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	GetVariables() *strset.Set
	GetFreeVariables() *strset.Set
	GetBoundVariables() *strset.Set
}

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
}

type (
	AtomicExpr struct {
		Pos     utils.SourceIndex
		Val     string
		ValKind token.TokenKind
	}

	TupleExpr struct {
		Pos   utils.SourceIndex
		Items []Expr
	}

	ListExpr struct {
		Pos   utils.SourceIndex
		Items []Expr
	}

	UnaryExpr struct {
		Pos     utils.SourceIndex
		Operand token.Token
		Expr    Expr
	}

	BinaryExpr struct {
		Pos     utils.SourceIndex
		LhsExpr Expr
		Operand token.Token
		RhsExpr Expr
	}

	ApplicationExpr struct {
		Pos  utils.SourceIndex
		Fun  Expr
		Args []Expr
	}

	ConditionalExpr struct {
		Pos         utils.SourceIndex
		Condition   Expr
		Consequent  Expr
		Alternative Expr
	}

	FunExpr struct {
		Pos      utils.SourceIndex
		ArgNames []string
		BodyExpr Expr
	}

	LetExpr struct {
		Pos       utils.SourceIndex
		InitNames []string
		InitExprs []Expr
		BodyExpr  Expr
		IsRec     bool
	}

	ModuleDeclStmt struct {
		Pos             utils.SourceIndex
		Name            string
		UsedModuleNames []string
		Aliases         map[string]string
		Stmts           []Stmt
	}

	LoadStmt struct {
		Pos     utils.SourceIndex
		Files   []string
		Aliases map[string]string
	}

	DefStmt struct {
		Pos   utils.SourceIndex
		Names []string
		Exprs []Expr
		IsRec bool
	}
)

func commaJoinStringReprsExpr(items ...Expr) string {
	var tmp []string

	for _, e := range items {
		tmp = append(tmp, e.GetStringRepr())
	}

	return strings.Join(tmp, ", ")
}

func commaJoinJsonReprsExpr(items ...Expr) string {
	var tmp []string

	for _, e := range items {
		tmp = append(tmp, e.GetJsonRepr())
	}

	return strings.Join(tmp, ", ")
}

func commaJoinQuotedIdents(idents ...string) string {
	var tmp []string

	for _, n := range idents {
		tmp = append(tmp, "\""+n+"\"")
	}

	return strings.Join(tmp, ", ")
}

func (expr AtomicExpr) GetPos() utils.SourceIndex { return expr.Pos }

func (expr AtomicExpr) GetStringRepr() string { return expr.Val }

func (expr AtomicExpr) GetJsonRepr() string {
	return fmt.Sprintf(`{"NodeName": "AtomicExpr", "Value": %#v}`, expr.Val)
}

func (expr AtomicExpr) GetVariables() *strset.Set {
	if expr.ValKind == token.Ident {
		return strset.New(expr.Val)
	} else {
		return strset.New()
	}
}

func (expr AtomicExpr) GetFreeVariables() *strset.Set {
	return expr.GetVariables()
}

func (expr AtomicExpr) GetBoundVariables() *strset.Set {
	return strset.New() // empty set
}

func (expr ListExpr) GetPos() utils.SourceIndex { return expr.Pos }

func (expr ListExpr) GetStringRepr() string {
	return fmt.Sprintf("[%s]", commaJoinStringReprsExpr(expr.Items...))
}

func (expr ListExpr) GetJsonRepr() string {
	return fmt.Sprintf(`{"NodeName": "ListExpr", "Items": [%s]}`, commaJoinJsonReprsExpr(expr.Items...))
}

func (expr ListExpr) GetVariables() *strset.Set {
	set := strset.New()

	for _, item := range expr.Items {
		set.Merge(item.GetVariables())
	}

	return set
}

func (expr ListExpr) GetFreeVariables() *strset.Set {
	set := strset.New()

	for _, item := range expr.Items {
		set.Merge(item.GetFreeVariables())
	}

	return set
}

func (expr ListExpr) GetBoundVariables() *strset.Set {
	set := strset.New()

	for _, item := range expr.Items {
		set.Merge(item.GetBoundVariables())
	}

	return set
}

func (expr UnaryExpr) GetPos() utils.SourceIndex { return expr.Pos }

func (expr UnaryExpr) GetStringRepr() string {
	return fmt.Sprintf("%s %s", expr.Operand.Kind.String(), expr.Expr.GetStringRepr())
}

func (expr UnaryExpr) GetJsonRepr() string {
	return fmt.Sprintf(`{"NodeName": "UnaryExpr", "Expr": %s}`, expr.Expr.GetJsonRepr())
}

func (expr UnaryExpr) GetVariables() *strset.Set {
	return expr.Expr.GetVariables()
}

func (expr UnaryExpr) GetFreeVariables() *strset.Set {
	return expr.Expr.GetFreeVariables()
}

func (expr UnaryExpr) GetBoundVariables() *strset.Set {
	return expr.Expr.GetBoundVariables()
}

func (expr BinaryExpr) GetPos() utils.SourceIndex { return expr.Pos }

func (expr BinaryExpr) GetStringRepr() string {
	return fmt.Sprintf(
		"(%s) %s (%s)", expr.LhsExpr.GetStringRepr(), expr.Operand.Kind.String(), expr.RhsExpr.GetStringRepr())
}

func (expr BinaryExpr) GetJsonRepr() string {
	return fmt.Sprintf(`{"NodeName": "BinaryExpr", "Operand": "%s", "LhsExpr": %s, "RhsExpr": %s}`,
		expr.Operand.Raw, expr.LhsExpr.GetJsonRepr(), expr.RhsExpr.GetJsonRepr())
}

func (expr BinaryExpr) GetVariables() *strset.Set {
	lhsVars, rhsVars := expr.LhsExpr.GetVariables(), expr.RhsExpr.GetVariables()
	return strset.Union(lhsVars, rhsVars)
}

func (expr BinaryExpr) GetFreeVariables() *strset.Set {
	lhsSet, rhsSet := expr.LhsExpr.GetFreeVariables(), expr.RhsExpr.GetFreeVariables()
	return strset.Union(lhsSet, rhsSet)
}

func (expr BinaryExpr) GetBoundVariables() *strset.Set {
	lhsSet, rhsSet := expr.LhsExpr.GetBoundVariables(), expr.RhsExpr.GetBoundVariables()
	return strset.Union(lhsSet, rhsSet)
}

func (expr ApplicationExpr) GetPos() utils.SourceIndex { return expr.Pos }

func (expr ApplicationExpr) GetStringRepr() string {
	return fmt.Sprintf("(%s)(%s)", expr.Fun.GetStringRepr(), commaJoinStringReprsExpr(expr.Args...))
}

func (expr ApplicationExpr) GetJsonRepr() string {
	return fmt.Sprintf(`{"NodeName": "ApplicationExpr", "Functions": %s, "Args": [%s]}`,
		expr.Fun.GetJsonRepr(), commaJoinJsonReprsExpr(expr.Args...))
}

func (expr ApplicationExpr) GetVariables() *strset.Set {
	set := strset.New()

	for _, arg := range expr.Args {
		set.Merge(arg.GetVariables())
	}

	return set
}

func (expr ApplicationExpr) GetFreeVariables() *strset.Set {
	set := strset.New()

	for _, arg := range expr.Args {
		set.Merge(arg.GetFreeVariables())
	}

	return set
}

func (expr ApplicationExpr) GetBoundVariables() *strset.Set {
	set := strset.New()

	for _, arg := range expr.Args {
		set.Merge(arg.GetBoundVariables())
	}

	return set
}

func (expr ConditionalExpr) GetPos() utils.SourceIndex { return expr.Pos }

func (expr ConditionalExpr) GetStringRepr() string {
	alternative := "nil"

	if expr.Alternative != nil {
		alternative = expr.Alternative.GetStringRepr()
	}

	return fmt.Sprintf("if (%s) then (%s) else (%s)",
		expr.Condition.GetStringRepr(), expr.Consequent.GetStringRepr(), alternative)
}

func (expr ConditionalExpr) GetJsonRepr() string {
	var alternative string

	if expr.Alternative != nil {
		alternative = expr.Alternative.GetJsonRepr()
	} else {
		alternative = `"nil"`
	}

	return fmt.Sprintf(`{"NodeName": "ConditionalExpr", "Condition": %s, "Consequent": %s, "Alternative": %s}`,
		expr.Condition.GetJsonRepr(), expr.Consequent.GetJsonRepr(), alternative)
}

func (expr ConditionalExpr) GetVariables() *strset.Set {
	if expr.Alternative != nil {
		return strset.Union(expr.Condition.GetVariables(), expr.Alternative.GetVariables(), expr.Consequent.GetVariables())
	} else {
		return strset.Union(expr.Condition.GetVariables(), expr.Consequent.GetVariables())
	}
}

func (expr ConditionalExpr) GetFreeVariables() *strset.Set {
	if expr.Alternative != nil {
		return strset.Union(expr.Condition.GetFreeVariables(), expr.Alternative.GetFreeVariables(), expr.Consequent.GetFreeVariables())
	} else {
		return strset.Union(expr.Condition.GetFreeVariables(), expr.Consequent.GetFreeVariables())
	}
}

func (expr ConditionalExpr) GetBoundVariables() *strset.Set {
	if expr.Alternative != nil {
		return strset.Union(expr.Condition.GetBoundVariables(), expr.Alternative.GetBoundVariables(), expr.Consequent.GetBoundVariables())
	} else {
		return strset.Union(expr.Condition.GetBoundVariables(), expr.Consequent.GetBoundVariables())
	}
}

func (expr FunExpr) GetPos() utils.SourceIndex { return expr.Pos }

func (expr FunExpr) GetStringRepr() string {
	return fmt.Sprintf("fn %s -> (%s)", strings.Join(expr.ArgNames, ", "), expr.BodyExpr.GetStringRepr())
}

func (expr FunExpr) GetJsonRepr() string {
	return fmt.Sprintf(`{"NodeName": "FunExpr", "ArgNames": [%s], "BodyExpr": %s}`,
		strings.Join(expr.ArgNames, ", "), expr.BodyExpr.GetJsonRepr())
}

func (expr FunExpr) GetVariables() *strset.Set {
	return expr.BodyExpr.GetVariables()
}

func (expr FunExpr) GetFreeVariables() *strset.Set {
	bodyExprSet, argNamesSet := expr.BodyExpr.GetFreeVariables(), strset.New()

	for _, name := range expr.ArgNames {
		argNamesSet.Add(name)
	}

	return strset.Difference(bodyExprSet, strset.Union(argNamesSet, bodyExprSet))
}

func (expr FunExpr) GetBoundVariables() *strset.Set {
	argNamesSet := strset.New()

	for _, name := range expr.ArgNames {
		argNamesSet.Add(name)
	}

	return strset.Intersection(argNamesSet, expr.BodyExpr.GetFreeVariables())
}

func (expr FunExpr) Arity() int {
	return len(expr.ArgNames)
}

func (expr FunExpr) Partialise(n int) FunExpr {
	return FunExpr{ArgNames: expr.ArgNames[n:], BodyExpr: expr.BodyExpr}
}

func (expr LetExpr) GetPos() utils.SourceIndex { return expr.Pos }

func (expr LetExpr) GetStringRepr() string {
	return fmt.Sprintf("let %s = %s in %s", strings.Join(expr.InitNames, ", "),
		commaJoinStringReprsExpr(expr.InitExprs...), expr.BodyExpr.GetStringRepr())
}

func (expr LetExpr) GetJsonRepr() string {
	return fmt.Sprintf(`{"NodeName": "LetExpr", "InitNames": [%s], "InitExprs": [%s], "BodyExpr": %s}`,
		commaJoinQuotedIdents(expr.InitNames...), commaJoinJsonReprsExpr(expr.InitExprs...), expr.BodyExpr.GetJsonRepr())
}

func (expr LetExpr) GetVariables() *strset.Set {
	initExprsSet := strset.New()

	for _, initExpr := range expr.InitExprs {
		initExprsSet.Merge(initExpr.GetVariables())
	}

	return strset.Union(expr.BodyExpr.GetVariables(), initExprsSet)
}

func (expr LetExpr) GetFreeVariables() *strset.Set {
	initNamesSet, initExprsSet := strset.New(), strset.New()

	for i := range expr.InitNames {
		initNamesSet.Add(expr.InitNames[i])
		initExprsSet.Merge(expr.InitExprs[i].GetVariables())
	}

	return strset.Difference(strset.Union(expr.BodyExpr.GetFreeVariables(), initExprsSet), initNamesSet)
}

func (expr LetExpr) GetBoundVariables() *strset.Set {
	initNamesSet, initExprsSet := strset.New(), strset.New()

	for i := range expr.InitNames {
		initNamesSet.Add(expr.InitNames[i])
		initExprsSet.Merge(expr.InitExprs[i].GetVariables())
	}

	return strset.Intersection(strset.Union(expr.BodyExpr.GetBoundVariables(), initExprsSet), initNamesSet)
}

func (stmt DefStmt) GetPos() utils.SourceIndex { return stmt.Pos }

func (stmt DefStmt) GetStringRepr() string {
	return fmt.Sprintf("def %s = %s",
		strings.Join(stmt.Names, ", "), commaJoinStringReprsExpr(stmt.Exprs...))
}

func (stmt DefStmt) GetJsonRepr() string {
	return fmt.Sprintf(`{"NodeName": "DefStmt", "Names": [%s], "Exprs": [%s]}`,
		commaJoinQuotedIdents(stmt.Names...), commaJoinJsonReprsExpr(stmt.Exprs...))
}

func (stmt ModuleDeclStmt) GetPos() utils.SourceIndex { return stmt.Pos }

func (stmt ModuleDeclStmt) GetStringRepr() string {
	var usedModulesWithAliases []string
	var declStmts []string

	for _, name := range stmt.UsedModuleNames {
		if alias, ok := stmt.Aliases[name]; ok {
			usedModulesWithAliases = append(usedModulesWithAliases, fmt.Sprint(name, "as", alias))
		} else {
			usedModulesWithAliases = append(usedModulesWithAliases, name)
		}
	}

	for _, modStmt := range stmt.Stmts {
		declStmts = append(declStmts, modStmt.GetStringRepr())
	}

	return fmt.Sprint("module", stmt.Name, "using", strings.Join(usedModulesWithAliases, ", "),
		"{\n", strings.Join(declStmts, "\n"), "\n}")
}

func (stmt ModuleDeclStmt) GetJsonRepr() string {
	var usedModuleNames []string
	var aliases []string
	var stmts []string

	for _, name := range stmt.UsedModuleNames {
		usedModuleNames = append(usedModuleNames, "\""+name+"\"")
		if alias, ok := stmt.Aliases[name]; ok {
			aliases = append(aliases, fmt.Sprint("\""+name+"\"", ":", "\""+alias+"\""))
		}
	}

	for _, modStmt := range stmt.Stmts {
		stmts = append(stmts, modStmt.GetJsonRepr())
	}

	return fmt.Sprintf(`{"NodeName": "ModuleDeclStmt", "UsedModuleNames": [%s], "Aliases": {%s}, "Stmts"": [%s]}`,
		strings.Join(usedModuleNames, ", "), strings.Join(aliases, ", "), strings.Join(stmts, ", "))
}
