package runtime

import "fmt"

type BooleanObject struct {
	Val ObjectValue
}

func NewBooleanObject(val bool) BooleanObject {
	return BooleanObject{Val: ObjectValue{Boolean: val}}
}

func (obj BooleanObject) Repr() string { return fmt.Sprintf("Boolean(%v)", obj.Val.Boolean) }

func (obj BooleanObject) Inspect() ObjectType { return BOOLEAN }

func (obj BooleanObject) GetVal() ObjectValue { return obj.Val }

func (obj BooleanObject) AsBool() Object { return obj }

func (obj BooleanObject) AsInt() Object { return nil /* todo: runtime error */ }

func (obj BooleanObject) AsReal() Object { return nil /* todo: runtime error */ }

func (obj BooleanObject) AsComplex() Object { return nil /* todo: runtime error */ }

func (obj BooleanObject) AsString() Object { return nil /* todo */ }

func (obj BooleanObject) AsCell() Object { return nil /* todo */ }

func (obj BooleanObject) AsList() Object { return nil /* todo */ }

func (obj BooleanObject) Equal(other Object) Object {
	switch other.Inspect() {
	case BOOLEAN:
		return NewBooleanObject(obj.Val.Boolean == other.GetVal().Boolean)
	default:
		return NewBooleanObject(false)
	}
}

func (obj BooleanObject) NotEqual(other Object) Object {
	if other.Inspect() == BOOLEAN {
		return NewBooleanObject(obj.Val.Boolean != other.GetVal().Boolean)
	} else {
		return NewBooleanObject(false)
	}
}

func (obj BooleanObject) Greater(other Object) Object {
	return nil // todo: runtime error - greater not defined on bool
}

func (obj BooleanObject) GreaterEqual(other Object) Object {
	return nil // todo: runtime error - greater equal not defined on bool
}

func (obj BooleanObject) Less(other Object) Object {
	return nil // todo: runtime error - less not defined on bool
}

func (obj BooleanObject) LessEqual(other Object) Object {
	return nil // todo: runtime error - less equal not defined on bool
}

func (obj BooleanObject) LNot() Object {
	return NewBooleanObject(!obj.Val.Boolean)
}

func (obj BooleanObject) LAnd(other Object) Object {
	return NewBooleanObject(obj.Val.Boolean && other.AsBool().GetVal().Boolean)
}

func (obj BooleanObject) LOr(other Object) Object {
	return NewBooleanObject(obj.Val.Boolean || other.AsBool().GetVal().Boolean)
}

func (obj BooleanObject) LXor(other Object) Object {
	return NewBooleanObject(obj.Val.Boolean != other.AsBool().GetVal().Boolean)
}

func (obj BooleanObject) BNot() Object {
	return nil // todo: runtime error - bitwise not not defined on bool
}

func (obj BooleanObject) BAnd(other Object) Object {
	return nil // todo: runtime error - bitwise and not defined on bool
}

func (obj BooleanObject) BOr(other Object) Object {
	return nil // todo: runtime error - bitwise or not defined on bool
}

func (obj BooleanObject) BXor(other Object) Object {
	return nil // todo: runtime error - bitwise xor not defined on bool
}

func (obj BooleanObject) Add(other Object) Object {
	return nil // todo: runtime error - add not defined on bool
}

func (obj BooleanObject) Sub(other Object) Object {
	return nil // todo: runtime error - sub not defined on bool
}

func (obj BooleanObject) Mod(other Object) Object {
	return nil // todo: runtime error - mod not defined on bool
}

func (obj BooleanObject) Mul(other Object) Object {
	return nil // todo: runtime error - mul not defined on bool
}

func (obj BooleanObject) Div(other Object) Object {
	return nil // todo: runtime error - div not defined on bool
}

func (obj BooleanObject) FloorDiv(other Object) Object {
	return nil // todo: runtime error - floordiv not defined on bool
}

func (obj BooleanObject) Pow(other Object) Object {
	return nil // todo: runtime error - pow not defined on bool
}

func (obj BooleanObject) Abs() Object {
	return nil // todo: runtime error - abs not defined on bool
}

func (obj BooleanObject) Nil() Object {
	return nil // todo: runtime error - nil not defined on bool
}

func (obj BooleanObject) Tail() Object {
	return nil // todo: runtime error - tail not defined on bool
}

func (obj BooleanObject) Head() Object {
	return nil // todo: runtime error - head not defined on bool
}
