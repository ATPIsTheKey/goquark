package runtime

import (
	"fmt"
)

type NilObject struct {
}

func NewNilObject() NilObject {
	return NilObject{}
}

func (obj NilObject) Repr() string {
	return fmt.Sprintf("NilObject()")
}

func (obj NilObject) Inspect() ObjectDescription { return ObjectDescription{NilType} }

func (obj NilObject) GetVal() ObjectValue { return ObjectValue{} }

func (obj NilObject) AsFun(frame *Frame) Object {
	return NewPoisonObject("Can't convert NilType to FunType", frame)
}

func (obj NilObject) AsBool(frame *Frame) Object {
	return NewPoisonObject("Can't convert NilType to BoolType", frame)
}

func (obj NilObject) AsInt(frame *Frame) Object {
	return NewPoisonObject("Can't convert NilType to IntType", frame)
}

func (obj NilObject) AsReal(frame *Frame) Object {
	return NewPoisonObject("Can't convert NilType to RealType", frame)
}

func (obj NilObject) AsComplex(frame *Frame) Object {
	return NewPoisonObject("Can't convert NilType to ComplexType", frame)
}

func (obj NilObject) AsString(frame *Frame) Object {
	return NewPoisonObject("Can't convert NilType to StringType", frame)
}

func (obj NilObject) AsList(frame *Frame) Object {
	return NewPoisonObject("Can't convert NilType to ListType", frame)
}

func (obj NilObject) AsTuple(frame *Frame) Object {
	return NewPoisonObject("Can't convert NilType to FunType", frame)
}

func (obj NilObject) Equal(other Object, frame *Frame) Object {
	if other.Inspect().Type == NilType {
		return NewBoolObject(true)
	} else {
		return NewBoolObject(false)
	}
}

func (obj NilObject) NotEqual(other Object, frame *Frame) Object {
	if other.Inspect().Type == NilType {
		return NewBoolObject(false)
	} else {
		return NewBoolObject(true)
	}
}

func (obj NilObject) Greater(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) GreaterEqual(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Less(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) LessEqual(other Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) LNot(frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) LAnd(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) LOr(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) LXor(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) BNot(frame *Frame) Object { return nil } // todo

func (obj NilObject) BAnd(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) BOr(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) BXor(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Add(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Sub(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Mod(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Mul(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Div(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) FloorDiv(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Pow(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Abs(frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Length(frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) GetItem(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Concatenate(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for NilType", frame)
}

func (obj NilObject) Apply(_ *Frame) Object { return obj }
