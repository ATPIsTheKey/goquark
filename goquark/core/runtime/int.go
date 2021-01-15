package runtime

import (
	"fmt"
	"math"
)

type IntObject struct {
	Val ObjectValue
}

func NewIntObject(val int64) IntObject {
	return IntObject{Val: ObjectValue{Integer: val}}
}

func (obj IntObject) Repr() string { return fmt.Sprintf("Integer(%v)", obj.Val.Integer) }

func (obj IntObject) Inspect() ObjectType { return INTEGER }

func (obj IntObject) GetVal() ObjectValue { return obj.Val }

func (obj IntObject) AsBool() Object { return NewBooleanObject(!(obj.Val.Integer == 0)) }

func (obj IntObject) AsInt() Object { return obj }

func (obj IntObject) AsReal() Object {
	return RealObject{Val: ObjectValue{Real: float64(obj.Val.Integer)}}
}

// automatically converts to complex128 with float64 as real part
func (obj IntObject) AsComplex() Object {
	return ComplexObject{Val: ObjectValue{Complex: complex(float64(obj.Val.Integer), 0)}}
}

func (obj IntObject) AsString() Object { return nil /* todo */ }

func (obj IntObject) AsCell() Object { return nil /* todo */ }

func (obj IntObject) AsList() Object { return nil /* todo */ }

func (obj IntObject) Equal(other Object) Object {
	return NewBooleanObject(HaveEqualVal(obj.Val, other.GetVal()))
}

func (obj IntObject) NotEqual(other Object) Object {
	return NewBooleanObject(!HaveEqualVal(obj.Val, other.GetVal()))
}

func (obj IntObject) Greater(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewBooleanObject(obj.Val.Integer > other.GetVal().Integer)
	case REAL:
		return other.Greater(obj.AsReal())
	case COMPLEX:
		return nil // todo: runtime error - cannot compare int and complex types
	default:
		return nil // todo: runtime error
	}
}

func (obj IntObject) GreaterEqual(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewBooleanObject(obj.Val.Integer >= other.GetVal().Integer)
	case REAL:
		return other.GreaterEqual(obj.AsReal())
	case COMPLEX:
		return nil // todo: runtime error - cannot compare int and complex types
	default:
		return nil // todo: runtime error
	}
}

func (obj IntObject) Less(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewBooleanObject(obj.Val.Integer < other.GetVal().Integer)
	case REAL:
		return other.Less(obj.AsReal())
	case COMPLEX:
		return nil // todo: runtime error - cannot compare int and complex types
	default:
		return nil // todo: runtime error
	}
}

func (obj IntObject) LessEqual(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewBooleanObject(obj.Val.Integer <= other.GetVal().Integer)
	case REAL:
		return other.LessEqual(obj.AsReal())
	case COMPLEX:
		return nil // todo: runtime error - cannot compare int and complex types
	default:
		return nil // todo: runtime error - cannot compare types
	}
}

func (obj IntObject) LNot() Object {
	return NewBooleanObject(!obj.AsBool().GetVal().Boolean)
}

func (obj IntObject) LAnd(other Object) Object {
	return NewBooleanObject(obj.AsBool().GetVal().Boolean && other.AsBool().GetVal().Boolean)
}

func (obj IntObject) LOr(other Object) Object {
	return NewBooleanObject(obj.AsBool().GetVal().Boolean || other.AsBool().GetVal().Boolean)
}

func (obj IntObject) LXor(other Object) Object {
	return NewBooleanObject(obj.AsBool().GetVal().Boolean != other.AsBool().GetVal().Boolean)
}

func (obj IntObject) BNot() Object {
	return NewIntObject(^obj.Val.Integer)
}

func (obj IntObject) BAnd(other Object) Object {
	if other.Inspect() == INTEGER {
		return NewIntObject(obj.Val.Integer & other.GetVal().Integer)
	} else {
		// todo: runtime error - bitwise and not defined on types
		return nil
	}
}

func (obj IntObject) BOr(other Object) Object {
	if other.Inspect() == INTEGER {
		return NewIntObject(obj.Val.Integer | other.GetVal().Integer)
	} else {
		// todo: runtime error - bitwise or not defined on types
		return nil
	}
}

func (obj IntObject) BXor(other Object) Object {
	if other.Inspect() == INTEGER {
		return NewIntObject(obj.Val.Integer ^ other.GetVal().Integer)
	} else {
		// todo: runtime error - bitwise xor not defined on types
		return nil
	}
}

func (obj IntObject) Add(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewIntObject(obj.Val.Integer + other.GetVal().Integer)
	case REAL:
		return other.Add(obj.AsReal())
	case COMPLEX:
		return other.Add(obj.AsComplex())
	default:
		return nil // todo: runtime error - cannot add types
	}
}

func (obj IntObject) Sub(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewIntObject(obj.Val.Integer - other.GetVal().Integer)
	case REAL:
		return other.Sub(obj.AsReal())
	case COMPLEX:
		return other.Sub(obj.AsComplex())
	default:
		return nil // todo: runtime error - cannot add types
	}
}

func (obj IntObject) Mod(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewIntObject(obj.Val.Integer % other.GetVal().Integer)
	default:
		return nil // todo: runtime error - cannot mod types
	}
}

func (obj IntObject) Mul(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewIntObject(obj.Val.Integer * other.GetVal().Integer)
	case REAL:
		return other.Mul(obj.AsReal())
	case COMPLEX:
		return other.Mul(obj.AsComplex())
	default:
		return nil // todo: runtime error - cannot mul types
	}
}

func (obj IntObject) Div(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewIntObject(obj.Val.Integer / other.GetVal().Integer)
	case REAL:
		return other.Div(obj.AsReal())
	case COMPLEX:
		return other.Div(obj.AsComplex())
	default:
		return nil // todo: runtime error - cannot div types
	}
}

func (obj IntObject) FloorDiv(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewIntObject(int64(math.Floor(obj.AsReal().GetVal().Real / other.AsReal().GetVal().Real)))
	case REAL:
		return other.FloorDiv(obj.AsReal())
	case COMPLEX:
		return nil // todo: runtime error - cant floor complex type
	default:
		return nil // todo: runtime error - cannot add types
	}
}

func (obj IntObject) Pow(other Object) Object {
	switch other.Inspect() {
	case INTEGER:
		return NewIntObject(int64(math.Pow(obj.AsReal().GetVal().Real, other.AsReal().GetVal().Real)))
	case REAL:
		return other.Pow(obj.AsReal())
	case COMPLEX:
		return other.Sub(obj.AsComplex())
	default:
		return nil // todo: runtime error - cannot pow types
	}
}

func (obj IntObject) Abs() Object {
	if obj.Val.Integer < 0 {
		return NewIntObject(-obj.Val.Integer)
	} else {
		return obj
	}
}

func (obj IntObject) Nil() Object {
	// todo: runtime error - nil not defined for int object
	return nil
}

func (obj IntObject) Tail() Object {
	// todo: runtime error - tail not defined for int object
	return nil
}

func (obj IntObject) Head() Object {
	// todo: runtime error - head not defined for int object
	return nil
}
