package runtime

import "fmt"

type RealObject struct {
	Val ObjectValue
}

func NewRealObject(val float64) RealObject {
	return RealObject{Val: ObjectValue{Real: val}}
}

func (obj RealObject) Repr() string { return fmt.Sprintf("AsReal(%v)", obj.Val.Real) }

func (obj RealObject) Inspect() ObjectType { return REAL }

func (obj RealObject) GetVal() ObjectValue { return obj.Val }

func (obj RealObject) AsBool() Object { return NewBooleanObject(!(obj.Val.Real == 0.0)) }

func (obj RealObject) AsInt() Object { return nil /* todo: runtime error */ }

func (obj RealObject) AsReal() Object { return obj }

// complex automatically converts to complex128 with float64 as real part
func (obj RealObject) AsComplex() Object {
	return ComplexObject{Val: ObjectValue{Complex: complex(obj.Val.Real, 0)}}
}

func (obj RealObject) AsString() Object { return nil /* todo */ }

func (obj RealObject) AsCell() Object { return nil /* todo */ }

func (obj RealObject) AsList() Object { return nil /* todo */ }

func (obj RealObject) Equal(other Object) Object { return nil /* todo */ }

func (obj RealObject) NotEqual(other Object) Object { return nil /* todo */ }

func (obj RealObject) Greater(other Object) Object { return nil /* todo */ }

func (obj RealObject) GreaterEqual(other Object) Object { return nil /* todo */ }

func (obj RealObject) Less(other Object) Object { return nil /* todo */ }

func (obj RealObject) LessEqual(other Object) Object { return nil /* todo */ }

func (obj RealObject) LNot() Object { return nil /* todo */ }

func (obj RealObject) LAnd(other Object) Object { return nil /* todo */ }

func (obj RealObject) LOr(other Object) Object { return nil /* todo */ }

func (obj RealObject) LXor(other Object) Object { return nil /* todo */ }

func (obj RealObject) BNot() Object { return nil /* todo */ }

func (obj RealObject) BAnd(other Object) Object { return nil /* todo */ }

func (obj RealObject) BOr(other Object) Object { return nil /* todo */ }

func (obj RealObject) BXor(other Object) Object { return nil /* todo */ }

func (obj RealObject) Add(other Object) Object { return nil /* todo */ }

func (obj RealObject) Sub(other Object) Object { return nil /* todo */ }

func (obj RealObject) Mod(other Object) Object { return nil /* todo */ }

func (obj RealObject) Mul(other Object) Object { return nil /* todo */ }

func (obj RealObject) Div(other Object) Object { return nil /* todo */ }

func (obj RealObject) FloorDiv(other Object) Object { return nil /* todo */ }

func (obj RealObject) Pow(other Object) Object { return nil /* todo */ }

func (obj RealObject) Abs() Object { return nil /* todo */ }

func (obj RealObject) Nil() Object { return nil /* todo */ }

func (obj RealObject) Tail() Object { return nil /* todo */ }

func (obj RealObject) Head() Object { return nil /* todo */ }
