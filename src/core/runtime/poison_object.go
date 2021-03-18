package runtime

import (
	"fmt"
)

type PoisonObject struct {
	ErrorMsg     string
	ReleaseFrame *Frame
}

func NewPoisonObject(errorMsg string, releaseFrame *Frame) PoisonObject {
	return PoisonObject{ErrorMsg: errorMsg, ReleaseFrame: releaseFrame}
}

func (obj PoisonObject) Repr() string {
	return fmt.Sprintf("PoisonObject()")
}

func (obj PoisonObject) Inspect() ObjectDescription { return ObjectDescription{PoisonType} }

func (obj PoisonObject) GetVal() ObjectValue { return ObjectValue{} }

func (obj PoisonObject) AsFun(_ *Frame) Object { return obj }

func (obj PoisonObject) AsBool(_ *Frame) Object { return obj }

func (obj PoisonObject) AsInt(_ *Frame) Object { return obj }

func (obj PoisonObject) AsReal(_ *Frame) Object { return obj }

func (obj PoisonObject) AsComplex(_ *Frame) Object { return obj }

func (obj PoisonObject) AsString(_ *Frame) Object { return obj }

func (obj PoisonObject) AsList(_ *Frame) Object { return obj }

func (obj PoisonObject) AsTuple(_ *Frame) Object { return obj }

func (obj PoisonObject) Equal(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) NotEqual(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Greater(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) GreaterEqual(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Less(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) LessEqual(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) LNot(_ *Frame) Object { return obj }

func (obj PoisonObject) LAnd(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) LOr(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) LXor(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) BNot(_ *Frame) Object { return obj }

func (obj PoisonObject) BAnd(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) BOr(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) BXor(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Add(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Sub(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Mod(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Mul(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Div(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) FloorDiv(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Pow(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Abs(_ *Frame) Object { return obj }

func (obj PoisonObject) Length(_ *Frame) Object { return obj }

func (obj PoisonObject) GetItem(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Concatenate(_ Object, _ *Frame) Object { return obj }

func (obj PoisonObject) Apply(_ *Frame) Object { return obj }
