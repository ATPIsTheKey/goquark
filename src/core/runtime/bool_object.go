package runtime

import (
	"fmt"
)

type BoolObject struct {
	Val ObjectValue
}

func NewBoolObject(val bool) BoolObject {
	return BoolObject{Val: ObjectValue{Bool: val}}
}

func (obj BoolObject) Repr() string { return fmt.Sprintf("BoolType(%v)", obj.Val.Bool) }

func (obj BoolObject) Inspect() ObjectDescription { return ObjectDescription{BoolType} }

func (obj BoolObject) GetVal() ObjectValue { return obj.Val }

func (obj BoolObject) AsFun(frame *Frame) Object { return nil }

func (obj BoolObject) AsBool(frame *Frame) Object { return obj }

func (obj BoolObject) AsInt(frame *Frame) Object {
	return NewPoisonObject("Can't convert BoolType to IntType", frame)
}

func (obj BoolObject) AsReal(frame *Frame) Object {
	return NewPoisonObject("Can't convert BoolType to RealType", frame)
}

func (obj BoolObject) AsComplex(frame *Frame) Object {
	return NewPoisonObject("Can't convert BoolType to ComplexType", frame)
}

func (obj BoolObject) AsString(frame *Frame) Object {
	return NewPoisonObject("Can't convert BoolType to StringType", frame)
}

func (obj BoolObject) AsList(frame *Frame) Object {
	return NewPoisonObject("Can't convert BoolType to ListType", frame)
}

func (obj BoolObject) AsTuple(frame *Frame) Object {
	return NewPoisonObject("Can't convert BoolType to TupleType", frame)
}

func (obj BoolObject) Equal(other Object, frame *Frame) Object {
	if other.Inspect().Type == BoolType {
		return NewBoolObject(obj.Val.Bool == other.GetVal().Bool)
	} else {
		return NewBoolObject(false)
	}
}

func (obj BoolObject) NotEqual(other Object, frame *Frame) Object {
	if other.Inspect().Type == BoolType {
		return NewBoolObject(obj.Val.Bool != other.GetVal().Bool)
	} else {
		return NewBoolObject(false)
	}
}

func (obj BoolObject) Greater(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for BoolType", frame)
}

func (obj BoolObject) GreaterEqual(_ Object, frame *Frame) Object {
	return NewPoisonObject("GreaterEqual not defined for BoolType", frame)
}

func (obj BoolObject) Less(_ Object, frame *Frame) Object {
	return NewPoisonObject("Less not defined for BoolType", frame)
}

func (obj BoolObject) LessEqual(_ Object, frame *Frame) Object {
	return NewPoisonObject("LessEqual not defined for BoolType", frame)
}

func (obj BoolObject) LNot(frame *Frame) Object {
	return NewBoolObject(!obj.Val.Bool)
}

func (obj BoolObject) LAnd(other Object, frame *Frame) Object {
	return NewBoolObject(obj.Val.Bool && other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj BoolObject) LOr(other Object, frame *Frame) Object {
	return NewBoolObject(obj.Val.Bool || other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj BoolObject) LXor(other Object, frame *Frame) Object {
	return NewBoolObject(obj.Val.Bool != other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj BoolObject) BNot(frame *Frame) Object {
	return NewPoisonObject("BNot not defined for BoolType", frame)
}

func (obj BoolObject) BAnd(_ Object, frame *Frame) Object {
	return NewPoisonObject("BAnd not defined for BoolType", frame)
}

func (obj BoolObject) BOr(_ Object, frame *Frame) Object {
	return NewPoisonObject("BOr not defined for BoolType", frame)
}

func (obj BoolObject) BXor(_ Object, frame *Frame) Object {
	return NewPoisonObject("BXor not defined for BoolType", frame)
}

func (obj BoolObject) Add(_ Object, frame *Frame) Object {
	return NewPoisonObject("Add not defined for BoolType", frame)
}

func (obj BoolObject) Sub(_ Object, frame *Frame) Object {
	return NewPoisonObject("Sub not defined for BoolType", frame)
}

func (obj BoolObject) Mod(_ Object, frame *Frame) Object {
	return NewPoisonObject("Mod not defined for BoolType", frame)
}

func (obj BoolObject) Mul(_ Object, frame *Frame) Object {
	return NewPoisonObject("Mul not defined for BoolType", frame)
}

func (obj BoolObject) Div(_ Object, frame *Frame) Object {
	return NewPoisonObject("Div not defined for BoolType", frame)
}

func (obj BoolObject) FloorDiv(_ Object, frame *Frame) Object {
	return NewPoisonObject("FloorDiv not defined for BoolType", frame)
}

func (obj BoolObject) Pow(_ Object, frame *Frame) Object {
	return NewPoisonObject("Greater not defined for BoolType", frame)
}

func (obj BoolObject) Abs(frame *Frame) Object {
	return NewPoisonObject("Abs not defined for BoolType", frame)
}

func (obj BoolObject) Length(frame *Frame) Object {
	return NewPoisonObject("Length not defined for BoolType", frame)
}

func (obj BoolObject) GetItem(item Object, frame *Frame) Object {
	return NewPoisonObject("GetItem not defined for BoolType", frame)
}

func (obj BoolObject) Concatenate(other Object, frame *Frame) Object {
	return NewPoisonObject("Concatenate not defined for BoolType", frame)
}

func (obj BoolObject) Apply(_ *Frame) Object { return obj }
