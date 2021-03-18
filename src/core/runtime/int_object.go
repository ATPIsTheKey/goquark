package runtime

import (
	"fmt"
	"math"
)

type IntObject struct {
	Val ObjectValue
}

func NewIntObject(val int64) IntObject {
	return IntObject{Val: ObjectValue{Int: val}}
}

func (obj IntObject) Repr() string { return fmt.Sprintf("IntObject(%v)", obj.Val.Int) }

func (obj IntObject) Inspect() ObjectDescription { return ObjectDescription{IntType} }

func (obj IntObject) GetVal() ObjectValue { return obj.Val }

func (obj IntObject) AsFun(frame *Frame) Object {
	return NewPoisonObject("Can't convert IntObject to FunType", frame)
}

func (obj IntObject) AsBool(_ *Frame) Object { return NewBoolObject(!(obj.Val.Int == 0)) }

func (obj IntObject) AsInt(_ *Frame) Object { return obj }

func (obj IntObject) AsReal(_ *Frame) Object { return NewRealObject(float64(obj.Val.Int)) }

// automatically converts to complex128 with float64 as real part
func (obj IntObject) AsComplex(_ *Frame) Object {
	return NewComplexObject(complex(float64(obj.Val.Int), 0))
}

func (obj IntObject) AsString(frame *Frame) Object {
	return NewPoisonObject("Can't convert IntObject to StringType", frame)
}

func (obj IntObject) AsList(frame *Frame) Object {
	return NewPoisonObject("Can't convert IntObject to ListType", frame)
}

func (obj IntObject) AsTuple(frame *Frame) Object {
	return NewPoisonObject("Can't convert IntObject to TupleType", frame)
}

// we can compare structs field by field in golang
func (obj IntObject) Equal(other Object, _ *Frame) Object {
	if other.Inspect().Type == IntType {
		return NewBoolObject(obj.Val.Int == other.GetVal().Int)
	} else {
		return NewBoolObject(false)
	}
}

// we can compare structs field by field in golang
func (obj IntObject) NotEqual(other Object, _ *Frame) Object {
	if other.Inspect().Type == IntType {
		return NewBoolObject(obj.Val.Int != other.GetVal().Int)
	} else {
		return NewBoolObject(false)
	}
}

func (obj IntObject) Greater(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType:
		return NewBoolObject(obj.Val.Int > other.GetVal().Int)
	case RealType:
		return other.Greater(obj, frame.New(""))
	default:
		return NewPoisonObject(fmt.Sprintf("__Greater not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) GreaterEqual(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType:
		return NewBoolObject(obj.Val.Int >= other.GetVal().Int)
	case RealType:
		return other.GreaterEqual(obj, frame.New(""))
	default:
		return NewPoisonObject(fmt.Sprintf("__GreaterEqual not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) Less(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType:
		return NewBoolObject(obj.Val.Int < other.GetVal().Int)
	case RealType:
		return other.Less(obj, frame.New(""))
	default:
		return NewPoisonObject(fmt.Sprintf("__Less not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) LessEqual(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType:
		return NewBoolObject(obj.Val.Int <= other.GetVal().Int)
	case RealType:
		return other.LessEqual(obj, frame.New(""))
	default:
		return NewPoisonObject(fmt.Sprintf("__LessEqual not defined for IntType and %s", other.Inspect().Type.String()), frame)
		return nil
	}
}

func (obj IntObject) LNot(frame *Frame) Object {
	return NewBoolObject(!obj.AsBool(frame.New("")).GetVal().Bool)
}

func (obj IntObject) LAnd(other Object, frame *Frame) Object {
	return NewBoolObject(obj.AsBool(frame.New("")).GetVal().Bool && other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj IntObject) LOr(other Object, frame *Frame) Object {
	return NewBoolObject(obj.AsBool(frame.New("")).GetVal().Bool || other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj IntObject) LXor(other Object, frame *Frame) Object {
	return NewBoolObject(obj.AsBool(frame.New("")).GetVal().Bool != other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj IntObject) BNot(_ *Frame) Object {
	return NewIntObject(^obj.Val.Int)
}

func (obj IntObject) BAnd(other Object, frame *Frame) Object {
	if other.Inspect().Type == IntType {
		return NewIntObject(obj.Val.Int & other.GetVal().Int)
	} else {
		return NewPoisonObject(fmt.Sprintf("__BAnd not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) BOr(other Object, frame *Frame) Object {
	if other.Inspect().Type == IntType {
		return NewIntObject(obj.Val.Int | other.GetVal().Int)
	} else {
		return NewPoisonObject(fmt.Sprintf("__BOr not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) BXor(other Object, frame *Frame) Object {
	if other.Inspect().Type == IntType {
		return NewIntObject(obj.Val.Int ^ other.GetVal().Int)
	} else {
		return NewPoisonObject(fmt.Sprintf("__BXor not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) Add(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType:
		return NewIntObject(obj.Val.Int + other.GetVal().Int)
	case RealType, ComplexType:
		return other.Add(obj, frame.New(""))
	default:
		return NewPoisonObject(fmt.Sprintf("__Add not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) Sub(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType:
		return NewIntObject(obj.Val.Int - other.GetVal().Int)
	case RealType, ComplexType:
		return other.Sub(obj, frame.New(""))
	default:
		return NewPoisonObject(fmt.Sprintf("__Sub not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) Mod(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType:
		return NewIntObject(obj.Val.Int % other.GetVal().Int)
	default:
		return NewPoisonObject(fmt.Sprintf("__Mod not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) Mul(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType:
		return NewIntObject(obj.Val.Int * other.GetVal().Int)
	case RealType, ComplexType:
		return other.Mul(obj, frame.New(""))
	default:
		return NewPoisonObject(fmt.Sprintf("__Mul not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) Div(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType:
		return NewIntObject(obj.Val.Int / other.GetVal().Int)
	case RealType, ComplexType:
		return other.Div(obj, frame.New(""))
	default:
		return NewPoisonObject(fmt.Sprintf("__Div not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) FloorDiv(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewIntObject(int64(obj.AsReal(frame.New("")).GetVal().Real / other.AsReal(frame.New("")).GetVal().Real))
	default:
		return NewPoisonObject(fmt.Sprintf("__FloorDiv not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) Pow(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewIntObject(int64(math.Pow(obj.AsReal(frame.New("")).GetVal().Real, other.AsReal(frame.New("")).GetVal().Real)))
	case ComplexType:
		return other.Pow(obj, frame.New(""))
	default:
		return NewPoisonObject(fmt.Sprintf("__Pow not defined for IntType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj IntObject) Abs(frame *Frame) Object {
	if obj.Val.Int < 0 {
		return NewIntObject(-obj.Val.Int)
	} else {
		return obj
	}
}

func (obj IntObject) Length(frame *Frame) Object {
	return NewPoisonObject("__Length not defined for IntType", frame)
}

func (obj IntObject) GetItem(_ Object, frame *Frame) Object {
	return NewPoisonObject("__GetItem not defined for IntType", frame)
}

func (obj IntObject) Concatenate(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Concatenate not defined for IntType", frame)
}

func (obj IntObject) Apply(_ *Frame) Object { return obj }
