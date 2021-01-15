package evaluate

import "goquark/goquark/core/ast"

type ArgStack []ast.Expr

func (stack *ArgStack) IsEmpty() bool {
	return len(*stack) == 0
}

func (stack *ArgStack) Push(node ast.Expr) {
	*stack = append(*stack, node)
}

func (stack *ArgStack) PushFromSliceReverse(nodes *[]ast.Expr) {
	last := len(*nodes)
	for i := range *nodes {
		stack.Push((*nodes)[last-1-i])
	}
}

func (stack *ArgStack) Pop() (ast.Expr, bool) {
	if stack.IsEmpty() {
		return nil, false
	} else {
		i := len(*stack) - 1
		node := (*stack)[i]
		*stack = (*stack)[:i]

		return node, true
	}
}
