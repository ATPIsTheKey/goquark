package ast

type Visitor interface {
	Visit(node Node) (w Visitor)
}

func walkStmtList(v Visitor, list []Stmt) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkIDList(v Visitor, list []AtomExpr) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkExprList(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}

func Walk(v Visitor, node Node) {}
