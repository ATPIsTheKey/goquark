package evaluate

import (
	"goquark/goquark/core/ast"
)

type ClosureEntry struct {
	Node        ast.Node
	EvalClosure *Closure
}

type Closure struct {
	Parent   *Closure
	EntryMap map[string]ClosureEntry
}

func MakeVoidClosure() Closure {
	return Closure{Parent: nil, EntryMap: map[string]ClosureEntry{}}
}

func (closure *Closure) LookupNode(nameVal string) (ClosureEntry, bool) {
	entry, ok := closure.EntryMap[nameVal]
	return entry, ok
}

func (closure *Closure) PushNode(nameVal string, node ast.Node, evalClosure *Closure) {
	closure.EntryMap[nameVal] = ClosureEntry{
		Node:        node,
		EvalClosure: evalClosure,
	}
}

func (closure *Closure) Copy() Closure {
	entryMapCopy := make(map[string]ClosureEntry)

	for k, val := range closure.EntryMap {
		entryMapCopy[k] = val
	}

	return Closure{closure.Parent, entryMapCopy}
}

func (closure *Closure) MakeChild() Closure {
	entryMapCopy := make(map[string]ClosureEntry)

	for k, val := range closure.EntryMap {
		entryMapCopy[k] = val
	}

	return Closure{closure, entryMapCopy}
}
