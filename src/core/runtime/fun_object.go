package runtime

import (
	"fmt"
)

type FunObject struct {
	Val ObjectValue
}

func NewFunObject(callback func(frame *Frame) Object) FunObject {
	return FunObject{ObjectValue{Fun: callback}}
}

func (obj FunObject) Repr() string { return fmt.Sprintf("FunType(arity=%d)") } // todo

func (obj FunObject) Inspect() ObjectDescription { return ObjectDescription{FunType} }

func (obj FunObject) GetVal() ObjectValue { return obj.Val }

func (obj FunObject) AsFun(frame *Frame) Object { return obj }

func (obj FunObject) AsBool(frame *Frame) Object { return NewBoolObject(true) }

func (obj FunObject) AsInt(frame *Frame) Object { return nil }

func (obj FunObject) AsReal(frame *Frame) Object { return nil }

// automatically converts to complex128 with float64 as real part
func (obj FunObject) AsComplex(frame *Frame) Object { return nil }

func (obj FunObject) AsString(frame *Frame) Object { return nil }

func (obj FunObject) AsList(frame *Frame) Object { return nil }

func (obj FunObject) AsTuple(frame *Frame) Object { return nil }

// we can compare structs field by field in golang
func (obj FunObject) Equal(_ Object, frame *Frame) Object { return nil }

// we can compare structs field by field in golang
func (obj FunObject) NotEqual(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Greater(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) GreaterEqual(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Less(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) LessEqual(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) LNot(frame *Frame) Object { return nil }

func (obj FunObject) LAnd(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) LOr(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) LXor(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) BNot(frame *Frame) Object { return nil }

func (obj FunObject) BAnd(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) BOr(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) BXor(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Add(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Sub(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Mod(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Mul(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Div(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) FloorDiv(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Pow(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Abs(frame *Frame) Object { return nil }

func (obj FunObject) Length(frame *Frame) Object { return nil }

func (obj FunObject) GetItem(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Concatenate(_ Object, frame *Frame) Object { return nil }

func (obj FunObject) Apply(callFrame *Frame) Object { return obj.Val.Fun(callFrame) }
