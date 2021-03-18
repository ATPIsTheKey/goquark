package runtime

import (
	"fmt"
	"math/cmplx"
)

type ComplexObject struct {
	Val ObjectValue
}

func NewComplexObject(val complex128) ComplexObject {
	return ComplexObject{Val: ObjectValue{Complex: val}}
}

func (obj ComplexObject) Repr() string { return fmt.Sprintf("ComplexType(%v)", obj.Val.Complex) }

func (obj ComplexObject) Inspect() ObjectDescription { return ObjectDescription{ComplexType} }

func (obj ComplexObject) GetVal() ObjectValue { return obj.Val }

func (obj ComplexObject) AsBool(frame *Frame) Object {
	return NewBoolObject(!(obj.Val.Complex == (0 + 0i)))
}

func (obj ComplexObject) AsFun(frame *Frame) Object {
	return NewPoisonObject("Can't convert ComplexType to FunType", frame)
}

func (obj ComplexObject) AsInt(frame *Frame) Object {
	return NewPoisonObject("Can't convert ComplexType to IntType", frame)
}

func (obj ComplexObject) AsReal(frame *Frame) Object {
	return NewPoisonObject("Can't convert ComplexType to RealType", frame)
}

func (obj ComplexObject) AsComplex(frame *Frame) Object { return obj }

func (obj ComplexObject) AsString(frame *Frame) Object {
	return NewPoisonObject("Can't convert ComplexType to StringType", frame)
}

func (obj ComplexObject) AsList(frame *Frame) Object {
	return NewPoisonObject("Can't convert ComplexType to ListType", frame)
}

func (obj ComplexObject) AsTuple(frame *Frame) Object {
	return NewPoisonObject("Can't convert ComplexType to TupleType", frame)
}

func (obj ComplexObject) Equal(other Object, frame *Frame) Object {
	if other.Inspect().Type == ComplexType {
		return NewBoolObject(obj.Val.Complex == other.GetVal().Complex)
	} else {
		return NewBoolObject(false)
	}
}

func (obj ComplexObject) NotEqual(other Object, frame *Frame) Object {
	if other.Inspect().Type == ComplexType {
		return NewBoolObject(obj.Val.Complex != other.GetVal().Complex)
	} else {
		return NewBoolObject(false)
	}
}

func (obj ComplexObject) Greater(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Greater not defined for ComplexType", frame)
}

func (obj ComplexObject) GreaterEqual(_ Object, frame *Frame) Object {
	return NewPoisonObject("__GreaterEqual not defined for ComplexType", frame)
}

func (obj ComplexObject) Less(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Less not defined for ComplexType", frame)
}

func (obj ComplexObject) LessEqual(_ Object, frame *Frame) Object {
	return NewPoisonObject("__LessEqual not defined for ComplexType", frame)
}

func (obj ComplexObject) LNot(frame *Frame) Object {
	return NewBoolObject(!obj.AsBool(frame.New("")).GetVal().Bool)
}

func (obj ComplexObject) LAnd(other Object, frame *Frame) Object {
	return NewBoolObject(obj.AsBool(frame.New("")).GetVal().Bool && other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj ComplexObject) LOr(other Object, frame *Frame) Object {
	return NewBoolObject(obj.AsBool(frame.New("")).GetVal().Bool || other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj ComplexObject) LXor(other Object, frame *Frame) Object {
	return NewBoolObject(obj.AsBool(frame.New("")).GetVal().Bool || other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj ComplexObject) BNot(frame *Frame) Object {
	return NewPoisonObject("__BNot not defined for ComplexType", frame)
}

func (obj ComplexObject) BAnd(_ Object, frame *Frame) Object {
	return NewPoisonObject("__BAnd not defined for ComplexType", frame)
}

func (obj ComplexObject) BOr(_ Object, frame *Frame) Object {
	return NewPoisonObject("__BOr not defined for ComplexType", frame)
}

func (obj ComplexObject) BXor(_ Object, frame *Frame) Object {
	return NewPoisonObject("__BXor not defined for ComplexType", frame)
}

func (obj ComplexObject) Add(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType, ComplexType:
		return NewComplexObject(obj.Val.Complex + other.AsComplex(frame.New("")).GetVal().Complex)
	default:
		return NewPoisonObject(fmt.Sprintf("__Add not defined for ComplexType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj ComplexObject) Sub(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType, ComplexType:
		return NewComplexObject(obj.Val.Complex - other.AsComplex(frame.New("")).GetVal().Complex)
	default:
		return NewPoisonObject(fmt.Sprintf("__Sub not defined for ComplexType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj ComplexObject) Mod(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Mod not defined for ComplexType", frame)
}

func (obj ComplexObject) Mul(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType, ComplexType:
		return NewComplexObject(obj.Val.Complex * other.AsComplex(frame.New("")).GetVal().Complex)
	default:
		return NewPoisonObject(fmt.Sprintf("__Mul not defined for ComplexType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj ComplexObject) Div(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType, ComplexType:
		return NewComplexObject(obj.Val.Complex / other.AsComplex(frame.New("")).GetVal().Complex)
	default:
		return NewPoisonObject(fmt.Sprintf("__Div not defined for ComplexType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj ComplexObject) FloorDiv(_ Object, frame *Frame) Object { return nil }

func (obj ComplexObject) Pow(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType, ComplexType:
		return NewComplexObject(cmplx.Pow(obj.GetVal().Complex, other.AsComplex(frame.New("")).GetVal().Complex))
	default:
		return NewPoisonObject(fmt.Sprintf("__Pow not defined for ComplexType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj ComplexObject) Abs(frame *Frame) Object {
	return NewRealObject(cmplx.Abs(obj.Val.Complex))
}

func (obj ComplexObject) Length(frame *Frame) Object {
	return NewPoisonObject("__Length not defined for ComplexType", frame)
}

func (obj ComplexObject) GetItem(_ Object, frame *Frame) Object {
	return NewPoisonObject("__GetItem not defined for ComplexType", frame)
}

func (obj ComplexObject) Concatenate(other Object, frame *Frame) Object {
	return NewPoisonObject("__Concatenate not defined for ComplexType", frame)
}

func (obj ComplexObject) Apply(frame *Frame) Object { return obj }
