package runtime

import (
	"fmt"
	"math"
)

type RealObject struct {
	Val ObjectValue
}

func NewRealObject(val float64) RealObject {
	return RealObject{Val: ObjectValue{Real: val}}
}

func (obj RealObject) Repr() string { return fmt.Sprintf("RealType(%v)", obj.Val.Real) }

func (obj RealObject) Inspect() ObjectDescription { return ObjectDescription{RealType} }

func (obj RealObject) GetVal() ObjectValue { return obj.Val }

func (obj RealObject) AsFun(frame *Frame) Object {
	return NewPoisonObject("Can't convert RealType to FunType", frame)
}

func (obj RealObject) AsBool(frame *Frame) Object { return NewBoolObject(!(obj.Val.Real == 0.0)) }

func (obj RealObject) AsInt(frame *Frame) Object { return NewIntObject(int64(obj.Val.Real)) }

func (obj RealObject) AsReal(frame *Frame) Object { return obj }

func (obj RealObject) AsComplex(frame *Frame) Object {
	return NewComplexObject(complex(obj.Val.Real, 0.0))
}

func (obj RealObject) AsString(frame *Frame) Object {
	return NewPoisonObject("Can't convert RealType to FunType", frame)
}

func (obj RealObject) AsList(frame *Frame) Object {
	return NewPoisonObject("Can't convert RealType to ListType", frame)
}

func (obj RealObject) AsTuple(frame *Frame) Object {
	return NewPoisonObject("Can't convert RealType to TupleType", frame)
}

// we can compare structs field by field in golang
func (obj RealObject) Equal(other Object, frame *Frame) Object {
	if other.Inspect().Type == RealType {
		return NewBoolObject(obj.Val.Real == other.GetVal().Real)
	} else {
		return NewBoolObject(false)
	}
}

func (obj RealObject) NotEqual(other Object, frame *Frame) Object {
	if other.Inspect().Type == RealType {
		return NewBoolObject(obj.Val.Real != other.GetVal().Real)
	} else {
		return NewBoolObject(false)
	}
}

func (obj RealObject) Greater(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewBoolObject(obj.Val.Real > other.AsReal(frame.New("")).GetVal().Real)
	default:
		return NewPoisonObject(fmt.Sprintf("__Greater not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) GreaterEqual(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewBoolObject(obj.Val.Real >= other.AsReal(frame.New("")).GetVal().Real)
	default:
		return NewPoisonObject(fmt.Sprintf("__GreaterEqual not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) Less(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewBoolObject(obj.Val.Real < other.AsReal(frame.New("")).GetVal().Real)
	default:
		return NewPoisonObject(fmt.Sprintf("__Less not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) LessEqual(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewBoolObject(obj.Val.Real <= other.AsReal(frame.New("")).GetVal().Real)
	default:
		return NewPoisonObject(fmt.Sprintf("__LessEqual not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) LNot(frame *Frame) Object {
	return NewBoolObject(!obj.AsBool(frame.New("")).GetVal().Bool)
}

func (obj RealObject) LAnd(other Object, frame *Frame) Object {
	return NewBoolObject(obj.AsBool(frame.New("")).GetVal().Bool && other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj RealObject) LOr(other Object, frame *Frame) Object {
	return NewBoolObject(obj.AsBool(frame.New("")).GetVal().Bool || other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj RealObject) LXor(other Object, frame *Frame) Object {
	return NewBoolObject(obj.AsBool(frame.New("")).GetVal().Bool != other.AsBool(frame.New("")).GetVal().Bool)
}

func (obj RealObject) BNot(frame *Frame) Object { return nil } // todo

func (obj RealObject) BAnd(_ Object, frame *Frame) Object { return nil } // todo

func (obj RealObject) BOr(_ Object, frame *Frame) Object { return nil } // todo

func (obj RealObject) BXor(_ Object, frame *Frame) Object { return nil } // todo

func (obj RealObject) Add(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewRealObject(obj.Val.Real + other.AsReal(frame.New("")).GetVal().Real)
	case ComplexType:
		return other.Add(obj, frame.New("")) // complex object will take care of type conversion
	default:
		return NewPoisonObject(fmt.Sprintf("__Add not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) Sub(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewRealObject(obj.Val.Real - other.AsReal(frame.New("")).GetVal().Real)
	case ComplexType:
		return other.Sub(obj, frame.New("")) // complex object will take care of type conversion
	default:
		return NewPoisonObject(fmt.Sprintf("__Sub not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) Mod(_ Object, frame *Frame) Object {
	return NewPoisonObject("__Mod not defined for RealType", frame)
}

func (obj RealObject) Mul(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewRealObject(obj.Val.Real * other.AsReal(frame.New("")).GetVal().Real)
	case ComplexType:
		return other.Mul(obj, frame.New("")) // complex object will take care of type conversion
	default:
		return NewPoisonObject(fmt.Sprintf("__Mul not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) Div(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewRealObject(obj.Val.Real / other.AsReal(frame.New("")).GetVal().Real)
	case ComplexType:
		return other.Div(obj, frame.New("")) // complex object will take care of type conversion
	default:
		return NewPoisonObject(fmt.Sprintf("__LessEqual not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) FloorDiv(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewIntObject(int64(math.Floor(obj.GetVal().Real / other.AsReal(frame.New("")).GetVal().Real)))
	default:
		return NewPoisonObject(fmt.Sprintf("__FloorDiv not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) Pow(other Object, frame *Frame) Object {
	switch other.Inspect().Type {
	case IntType, RealType:
		return NewRealObject(math.Pow(obj.GetVal().Real, other.AsReal(frame.New("")).GetVal().Real))
	default:
		return NewPoisonObject(fmt.Sprintf("__Pow not defined for RealType and %s", other.Inspect().Type.String()), frame)
	}
}

func (obj RealObject) Abs(frame *Frame) Object {
	if obj.Val.Real < 0 {
		return NewRealObject(-obj.Val.Real)
	} else {
		return obj
	}
}

func (obj RealObject) Length(frame *Frame) Object {
	return NewPoisonObject("__Length not defined for RealType", frame)
}

func (obj RealObject) GetItem(_ Object, frame *Frame) Object {
	return NewPoisonObject("__GetItem not defined for RealType", frame)
}

func (obj RealObject) Concatenate(other Object, frame *Frame) Object {
	return NewPoisonObject("__Concatenate not defined for RealType", frame)
}

func (obj RealObject) Apply(_ *Frame) Object { return obj }
