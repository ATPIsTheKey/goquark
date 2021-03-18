package runtime

import (
	"fmt"
	"strings"
)

type ListObject struct {
	Val ObjectValue
}

func NewListObject(objects ...Object) ListObject {
	return ListObject{Val: ObjectValue{Sequence: objects}}
}

func (obj ListObject) Repr() string {
	var reprs []string

	for _, item := range obj.Val.Sequence {
		reprs = append(reprs, item.Repr())
	}

	return fmt.Sprintf("ListObject(%s)", strings.Join(reprs, ", "))
}

func (obj ListObject) Inspect() ObjectDescription { return ObjectDescription{ListType} }

func (obj ListObject) GetVal() ObjectValue { return obj.Val }

func (obj ListObject) AsFun(frame *Frame) Object {
	return NewPoisonObject("Can't convert ListObject to FunObject", frame)
}

func (obj ListObject) AsBool(_ *Frame) Object {
	return NewBoolObject(len(obj.Val.Sequence) != 0)
}

func (obj ListObject) AsInt(frame *Frame) Object {
	return NewPoisonObject("Can't convert ListObject to IntObject", frame)
}

func (obj ListObject) AsReal(frame *Frame) Object {
	return NewPoisonObject("Can't convert ListObject to RealObject", frame)
}

func (obj ListObject) AsComplex(frame *Frame) Object {
	return NewPoisonObject("Can't convert ListObject to ComplexObject", frame)
}

func (obj ListObject) AsString(frame *Frame) Object {
	return NewPoisonObject("Can't convert ListObject to StringObject", frame)
}

func (obj ListObject) AsList(frame *Frame) Object { return obj }

func (obj ListObject) AsTuple(frame *Frame) Object { return nil } // todo

func (obj ListObject) Equal(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case ListType:
		if len(obj.Val.Sequence) != len(other.GetVal().Sequence) {
			return NewBoolObject(false)
		} else {
			for i := range obj.Val.Sequence {
				if obj.Val.Sequence[i].Equal(other.GetVal().Sequence[i], frame.New("todo")).GetVal().Bool == false {
					return NewBoolObject(false)
				}
			}
			return NewBoolObject(true)
		}
	default:
		return NewPoisonObject(fmt.Sprintf("__Equal not defined for ListObject and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj ListObject) NotEqual(other Object, frame *Frame) Object {
	return NewBoolObject(!obj.Equal(other, frame.New("todo")).GetVal().Bool)
}

func (obj ListObject) Greater(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Greater not defined for ListObject", frame)
}

func (obj ListObject) GreaterEqual(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Greater not defined for ListObject", frame)
}

func (obj ListObject) Less(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Greater not defined for ListObject", frame)
}

func (obj ListObject) LessEqual(other Object, frame *Frame) Object {
	return NewPoisonObject("__Greater not defined for ListObject", frame)
}

func (obj ListObject) LNot(frame *Frame) Object {
	return NewBoolObject(len(obj.Val.Sequence) == 0)
}

func (obj ListObject) LAnd(_ Object, frame *Frame) Object {
	return NewPoisonObject("__LAnd not defined for ListObject", frame)
}

func (obj ListObject) LOr(_ Object, frame *Frame) Object {
	return NewPoisonObject("__LOr not defined for ListObject", frame)
}

func (obj ListObject) LXor(_ Object, frame *Frame) Object {
	return NewPoisonObject("__LXor not defined for ListObject", frame)
}

func (obj ListObject) BNot(frame *Frame) Object { return nil } // todo

func (obj ListObject) BAnd(_ Object, frame *Frame) Object {
	return NewPoisonObject("__BAnd not defined for ListObject", frame)
}

func (obj ListObject) BOr(_ Object, frame *Frame) Object {
	return NewPoisonObject("__BOr not defined for ListObject", frame)
}

func (obj ListObject) BXor(_ Object, frame *Frame) Object {
	return NewPoisonObject("__BXor not defined for ListObject", frame)
}

func (obj ListObject) Add(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Add not defined for ListObject", frame)
}

func (obj ListObject) Sub(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Sub not defined for ListObject", frame)
}

func (obj ListObject) Mod(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Mod not defined for ListObject", frame)
}

func (obj ListObject) Mul(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Mul not defined for ListObject", frame)
}

func (obj ListObject) Div(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Div not defined for ListObject", frame)
}

func (obj ListObject) FloorDiv(_ Object, frame *Frame) Object {
	return NewPoisonObject("__FloorDiv not defined for ListObject", frame)
}

func (obj ListObject) Pow(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Pow not defined for ListObject", frame)
}

func (obj ListObject) Abs(frame *Frame) Object {
	return NewPoisonObject("__Abs not defined for ListObject", frame)
}

func (obj ListObject) Length(frame *Frame) Object {
	return NewIntObject(int64(len(obj.Val.Sequence)))
}

func (obj ListObject) GetItem(item Object, frame *Frame) Object {
	switch item.Inspect().Type {
	case IntType:
		if int64(len(obj.Val.Sequence)) < item.GetVal().Int+1 {
			return NewPoisonObject("todo", frame)
		} else {
			return obj.Val.Sequence[item.GetVal().Int]
		}
	default:
		return NewPoisonObject(fmt.Sprintf("__GetItem not defined for ListObject and %s", item.Inspect().Type), frame)
	}
}

func (obj ListObject) Concatenate(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case ListType:
		return NewListObject(append(obj.Val.Sequence, other.GetVal().Sequence...)...)
	default:
		return NewPoisonObject(fmt.Sprintf("__Concatenate not defined for ListObject and %s", other.Inspect().Type), frame)
	}
}

func (obj ListObject) Apply(_ *Frame) Object { return obj }
