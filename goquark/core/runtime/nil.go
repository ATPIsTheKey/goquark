package runtime

import "fmt"

type NilObject struct {
}

func NewNilObject() NilObject {
	return NilObject{}
}

func (obj NilObject) Repr() string {
	return fmt.Sprintf("Nil()")
}

func (obj NilObject) Inspect() ObjectType {
	return NIL
}

func (obj NilObject) GetVal() ObjectValue {
	return ObjectValue{}
}

func (obj NilObject) AsBool() Object { return NewBooleanObject(false) }

func (obj NilObject) AsInt() Object { return nil /* todo: runtime error */ }

func (obj NilObject) AsReal() Object { return nil /* todo: runtime error */ }

func (obj NilObject) AsComplex() Object { return nil /* todo: runtime error */ }

func (obj NilObject) AsString() Object { return nil /* todo */ }

func (obj NilObject) AsCell() Object { return nil /* todo */ }

func (obj NilObject) AsList() Object { return nil /* todo: runtime error */ }

func (obj NilObject) Equal(other Object) Object {
	if other.Inspect() == NIL {
		return NewBooleanObject(true)
	} else {
		return NewBooleanObject(false)
	}
}

func (obj NilObject) NotEqual(other Object) Object {
	if other.Inspect() == NIL {
		return NewBooleanObject(false)
	} else {
		return NewBooleanObject(true)
	}
}

func (obj NilObject) Greater(other Object) Object {
	return nil // todo: runtime error - greater not defined on nil
}

func (obj NilObject) GreaterEqual(other Object) Object {
	return nil // todo: runtime error - greater equal not defined on nil
}

func (obj NilObject) Less(other Object) Object {
	return nil // todo: runtime error - less not defined on nil
}

func (obj NilObject) LessEqual(other Object) Object {
	return nil // todo: runtime error - less equal not defined on nil
}

func (obj NilObject) LNot() Object {
	return nil // todo: runtime error - less equal not defined on nil
}

func (obj NilObject) LAnd(other Object) Object {
	return nil // todo: runtime error - less equal not defined on nil
}

func (obj NilObject) LOr(other Object) Object {
	return nil // todo: runtime error - less equal not defined on nil
}

func (obj NilObject) LXor(other Object) Object {
	return nil // todo: runtime error - less equal not defined on nil
}

func (obj NilObject) BNot() Object {
	return nil // todo: runtime error - bitwise not not defined on nil
}

func (obj NilObject) BAnd(other Object) Object {
	return nil // todo: runtime error - bitwise and not defined on nil
}

func (obj NilObject) BOr(other Object) Object {
	return nil // todo: runtime error - bitwise or not defined on nil
}

func (obj NilObject) BXor(other Object) Object {
	return nil // todo: runtime error - bitwise xor not defined on nil
}

func (obj NilObject) Add(other Object) Object {
	return nil // todo: runtime error - add not defined on nil
}

func (obj NilObject) Sub(other Object) Object {
	return nil // todo: runtime error - sub not defined on nil
}

func (obj NilObject) Mod(other Object) Object {
	return nil // todo: runtime error - mod not defined on nil
}

func (obj NilObject) Mul(other Object) Object {
	return nil // todo: runtime error - mul not defined on nil
}

func (obj NilObject) Div(other Object) Object {
	return nil // todo: runtime error - div not defined on nil
}

func (obj NilObject) FloorDiv(other Object) Object {
	return nil // todo: runtime error - floordiv not defined on nil
}

func (obj NilObject) Pow(other Object) Object {
	return nil // todo: runtime error - pow not defined on nil
}

func (obj NilObject) Abs() Object {
	return nil // todo: runtime error - abs not defined on nil
}

func (obj NilObject) Nil() Object {
	return nil // todo: runtime error - nil not defined on nil
}

func (obj NilObject) Tail() Object {
	return nil // todo: runtime error - tail not defined on nil
}

func (obj NilObject) Head() Object {
	return nil // todo: runtime error - head not defined on nil
}
