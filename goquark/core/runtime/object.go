package runtime

type ObjectType int

const (
	NIL ObjectType = iota
	BOOLEAN
	INTEGER
	REAL
	COMPLEX
	STRING
	CELL
	ARRAY
)

type ObjectValue struct {
	Boolean bool
	Integer int64
	Real    float64
	Complex complex128
	String  string
	Cell    *CellObject
	List    *ListObject
}

type Object interface {
	Repr() string
	Inspect() ObjectType
	GetVal() ObjectValue

	AsBool() Object
	AsInt() Object
	AsReal() Object
	AsComplex() Object
	AsString() Object
	AsCell() Object
	AsList() Object

	Equal(other Object) Object
	NotEqual(other Object) Object
	Greater(other Object) Object
	GreaterEqual(other Object) Object
	Less(other Object) Object
	LessEqual(other Object) Object
	LNot() Object
	LAnd(other Object) Object
	LOr(other Object) Object
	LXor(other Object) Object
	BNot() Object
	BAnd(other Object) Object
	BOr(other Object) Object
	BXor(other Object) Object
	Add(other Object) Object
	Sub(other Object) Object
	Mod(other Object) Object
	Mul(other Object) Object
	Div(other Object) Object
	FloorDiv(other Object) Object
	Pow(other Object) Object
	Abs() Object
	Nil() Object
	Tail() Object
	Head() Object
}

func IsNumeric(obj Object) bool {
	return obj.Inspect() == INTEGER || obj.Inspect() == REAL || obj.Inspect() == COMPLEX
}

func GetPromotedNumericObject(obj Object) Object {
	switch obj.Inspect() {
	case INTEGER:
		return obj.AsReal()
	case REAL:
		return obj.AsComplex()
	case COMPLEX:
		return obj
	default:
		return nil
	}
}

func HaveEqualVal(obj1 ObjectValue, obj2 ObjectValue) bool {
	return obj1 == obj2
}
