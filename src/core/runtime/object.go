package runtime

type ObjectType int

//go:generate stringer -type=ObjectType
const (
	PoisonType ObjectType = iota
	NilType
	ThunkType
	FunType
	BoolType
	IntType
	RealType
	ComplexType
	ListType
	DataType // todo
)

type ObjectValue struct {
	Bool     bool
	Int      int64
	Real     float64
	Complex  complex128
	String   string
	Fun      func(callFrame *Frame) Object
	Sequence []Object
}

type ObjectDescription struct {
	// todo: add more object descriptors in the future
	Type ObjectType
}

type Object interface {
	Repr() string
	Inspect() ObjectDescription
	GetVal() ObjectValue

	AsFun(frame *Frame) Object
	AsBool(frame *Frame) Object
	AsInt(frame *Frame) Object
	AsReal(frame *Frame) Object
	AsComplex(frame *Frame) Object
	AsString(frame *Frame) Object
	AsList(frame *Frame) Object
	AsTuple(frame *Frame) Object

	Equal(other Object, frame *Frame) Object
	NotEqual(other Object, frame *Frame) Object
	Greater(other Object, frame *Frame) Object
	GreaterEqual(other Object, frame *Frame) Object
	Less(other Object, frame *Frame) Object
	LessEqual(other Object, frame *Frame) Object
	LNot(frame *Frame) Object
	LAnd(other Object, frame *Frame) Object
	LOr(other Object, frame *Frame) Object
	LXor(other Object, frame *Frame) Object
	BNot(frame *Frame) Object
	BAnd(other Object, frame *Frame) Object
	BOr(other Object, frame *Frame) Object
	BXor(other Object, frame *Frame) Object
	Add(other Object, frame *Frame) Object
	Sub(other Object, frame *Frame) Object
	Mod(other Object, frame *Frame) Object
	Mul(other Object, frame *Frame) Object
	Div(other Object, frame *Frame) Object
	FloorDiv(other Object, frame *Frame) Object
	Pow(other Object, frame *Frame) Object
	Abs(frame *Frame) Object

	Length(frame *Frame) Object
	GetItem(item Object, frame *Frame) Object
	Concatenate(other Object, frame *Frame) Object

	Apply(frame *Frame) Object
}
