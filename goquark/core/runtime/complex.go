package runtime

import (
	"fmt"
)

type ComplexObject struct {
	Val ObjectValue
}

func NewComplexObject(val complex128) ComplexObject {
	return ComplexObject{Val: ObjectValue{Complex: val}}
}

func (obj ComplexObject) Repr() string { return fmt.Sprintf("AsComplex(%v)", obj.Val.Complex) }

func (obj ComplexObject) Inspect() ObjectType { return COMPLEX }

func (obj ComplexObject) GetVal() ObjectValue { return obj.Val }

func (obj ComplexObject) AsBool() Object { return NewBooleanObject(!(obj.Val.Complex == (0 + 0i))) }

func (obj ComplexObject) AsInt() Object { return nil /* todo: runtime error */ }

func (obj ComplexObject) AsReal() Object { return nil /* todo: runtime error */ }

func (obj ComplexObject) AsComplex() Object { return obj }

func (obj ComplexObject) AsString() Object { return nil /* todo */ }

func (obj ComplexObject) AsCell() Object { return nil /* todo */ }

func (obj ComplexObject) AsList() Object { return nil /* todo */ }

func (obj ComplexObject) Equal(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) NotEqual(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) Greater(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) GreaterEqual(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) Less(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) LessEqual(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) LNot() Object { return nil /* todo */ }

func (obj ComplexObject) LAnd(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) LOr(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) LXor(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) BNot() Object { return nil /* todo */ }

func (obj ComplexObject) BAnd(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) BOr(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) BXor(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) Add(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) Sub(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) Mod(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) Mul(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) Div(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) FloorDiv(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) Pow(other Object) Object { return nil /* todo */ }

func (obj ComplexObject) Abs() Object { return nil /* todo */ }

func (obj ComplexObject) Nil() Object { return nil /* todo */ }

func (obj ComplexObject) Tail() Object { return nil /* todo */ }

func (obj ComplexObject) Head() Object { return nil /* todo */ }
