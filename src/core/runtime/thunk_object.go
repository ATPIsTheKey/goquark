package runtime

type ThunkObject struct {
	Memo      Object
	Evaluator func() Object
}

func NewThunkObject(evaluator func() Object) ThunkObject {
	return ThunkObject{Memo: nil, Evaluator: evaluator}
}

func (obj ThunkObject) GetActualObject() Object {
	if obj.Memo == nil {
		// todo
		return obj.Evaluator()
	} else {
		return obj.Memo
	}
}

func (obj ThunkObject) Repr() string { return obj.GetActualObject().Repr() }

func (obj ThunkObject) Inspect() ObjectDescription { return obj.GetActualObject().Inspect() }

func (obj ThunkObject) GetVal() ObjectValue { return obj.GetActualObject().GetVal() }

func (obj ThunkObject) AsFun(frame *Frame) Object { return obj.GetActualObject().AsFun(frame) }

func (obj ThunkObject) AsBool(frame *Frame) Object { return obj.GetActualObject().AsBool(frame) }

func (obj ThunkObject) AsInt(frame *Frame) Object { return obj.GetActualObject().AsInt(frame) }

func (obj ThunkObject) AsReal(frame *Frame) Object { return obj.GetActualObject().AsReal(frame) }

// automatically converts to complex128 with float64 as real part
func (obj ThunkObject) AsComplex(frame *Frame) Object { return obj.GetActualObject().AsComplex(frame) }

func (obj ThunkObject) AsString(frame *Frame) Object { return obj.GetActualObject().AsString(frame) }

func (obj ThunkObject) AsList(frame *Frame) Object { return obj.GetActualObject().AsList(frame) }

func (obj ThunkObject) AsTuple(frame *Frame) Object { return obj.GetActualObject().AsTuple(frame) }

// we can compare structs field by field in golang
func (obj ThunkObject) Equal(other Object, frame *Frame) Object {
	return obj.GetActualObject().Equal(other, frame)
}

// we can compare structs field by field in golang
func (obj ThunkObject) NotEqual(other Object, frame *Frame) Object {
	return obj.GetActualObject().NotEqual(other, frame)
}

func (obj ThunkObject) Greater(other Object, frame *Frame) Object {
	return obj.GetActualObject().Greater(other, frame)
}

func (obj ThunkObject) GreaterEqual(other Object, frame *Frame) Object {
	return obj.GetActualObject().GreaterEqual(other, frame)
}

func (obj ThunkObject) Less(other Object, frame *Frame) Object {
	return obj.GetActualObject().Less(other, frame)
}

func (obj ThunkObject) LessEqual(other Object, frame *Frame) Object {
	return obj.GetActualObject().LessEqual(other, frame)
}

func (obj ThunkObject) LNot(frame *Frame) Object { return obj.GetActualObject().LNot(frame) }

func (obj ThunkObject) LAnd(other Object, frame *Frame) Object {
	return obj.GetActualObject().LAnd(other, frame)
}

func (obj ThunkObject) LOr(other Object, frame *Frame) Object {
	return obj.GetActualObject().LOr(other, frame)
}

func (obj ThunkObject) LXor(other Object, frame *Frame) Object {
	return obj.GetActualObject().LXor(other, frame)
}

func (obj ThunkObject) BNot(frame *Frame) Object { return obj.GetActualObject().BNot(frame) }

func (obj ThunkObject) BAnd(other Object, frame *Frame) Object {
	return obj.GetActualObject().BAnd(other, frame)
}

func (obj ThunkObject) BOr(other Object, frame *Frame) Object {
	return obj.GetActualObject().BOr(other, frame)
}

func (obj ThunkObject) BXor(other Object, frame *Frame) Object {
	return obj.GetActualObject().BXor(other, frame)
}

func (obj ThunkObject) Add(other Object, frame *Frame) Object {
	return obj.GetActualObject().Add(other, frame)
}

func (obj ThunkObject) Sub(other Object, frame *Frame) Object {
	return obj.GetActualObject().Sub(other, frame)
}

func (obj ThunkObject) Mod(other Object, frame *Frame) Object {
	return obj.GetActualObject().Mod(other, frame)
}

func (obj ThunkObject) Mul(other Object, frame *Frame) Object {
	return obj.GetActualObject().Mul(other, frame)
}

func (obj ThunkObject) Div(other Object, frame *Frame) Object {
	return obj.GetActualObject().Div(other, frame)
}

func (obj ThunkObject) FloorDiv(other Object, frame *Frame) Object {
	return obj.GetActualObject().FloorDiv(other, frame)
}

func (obj ThunkObject) Pow(other Object, frame *Frame) Object {
	return obj.GetActualObject().Pow(other, frame)
}

func (obj ThunkObject) Abs(frame *Frame) Object { return obj.GetActualObject().Abs(frame) }

func (obj ThunkObject) Length(frame *Frame) Object { return obj.GetActualObject().Length(frame) }

func (obj ThunkObject) GetItem(item Object, frame *Frame) Object {
	return obj.GetActualObject().GetItem(item, frame)
}

func (obj ThunkObject) Concatenate(other Object, frame *Frame) Object {
	return obj.GetActualObject().Concatenate(other, frame)
}

func (obj ThunkObject) Apply(frame *Frame) Object { return obj.GetActualObject().Apply(frame) }
