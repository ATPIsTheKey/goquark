package runtime

import (
	"github.com/scylladb/go-set/strset"
	"strings"
)

type Env map[string]Object

type ObjectStack struct {
	objects []Object
}

func NewObjectStack(init ...Object) *ObjectStack {
	return &ObjectStack{objects: init}
}

func (stack *ObjectStack) Size() int { return len(stack.objects) }

func (stack *ObjectStack) Push(x ...Object) { stack.objects = append(stack.objects, x...) }

func (stack *ObjectStack) Pop() (Object, bool) {
	if len(stack.objects) == 0 {
		return nil, false
	}
	i := len(stack.objects) - 1
	ret := stack.objects[i]
	stack.objects[i] = nil
	stack.objects = stack.objects[:i]
	return ret, true
}
func (stack *ObjectStack) Merge(other *ObjectStack) {
	stack.Push(other.objects...)
}

func (stack *ObjectStack) GetMerged(other *ObjectStack) *ObjectStack {
	newStack := NewObjectStack(stack.objects...)
	newStack.Push(other.objects...)
	return newStack
}

type Frame struct {
	Parent      *Frame
	Env         Env
	ArgStack    *ObjectStack
	Description string
}

func NewRootFrame() *Frame {
	return &Frame{
		Parent:      nil,
		Env:         make(map[string]Object),
		ArgStack:    NewObjectStack(),
		Description: "__NewRootClosure()",
	}
}

func (frame *Frame) GetFromEnv(name string) (Object, bool) {
	if entry, ok := frame.Env[name]; ok {
		return entry, true
	} else if frame.Parent != nil {
		return frame.Parent.GetFromEnv(name)
	} else {
		return nil, false
	}
}

func (frame *Frame) PushToEnv(name string, obj Object) {
	frame.Env[name] = obj
}

func (frame *Frame) Copy() *Frame {
	envCopy := make(Env)

	for k, val := range frame.Env {
		envCopy[k] = val
	}

	return &Frame{
		Parent:      frame.Parent,
		Env:         envCopy,
		ArgStack:    frame.ArgStack,
		Description: frame.Description,
	}
}

func (frame *Frame) New(description string) *Frame {
	return &Frame{
		Parent:      frame,
		Env:         make(Env),
		ArgStack:    frame.ArgStack,
		Description: description,
	}
}

func (frame *Frame) EnvNames() *strset.Set {
	names := strset.New()
	for k := range frame.Env {
		names.Add(k)
	}
	if frame.Parent != nil {
		return strset.Union(names, frame.Parent.EnvNames())
	} else {
		return names
	}
}

func (frame *Frame) BuildTraceback() string {
	var tracebackDump strings.Builder

	for currFrame := frame; currFrame != nil; {
		tracebackDump.WriteString(currFrame.Description + "\n")
		currFrame = currFrame.Parent
	}

	return tracebackDump.String()
}
