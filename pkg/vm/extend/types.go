package extend

import "github.com/deepfabric/vectorsql/pkg/vm/value"

type Extend interface {
	String() string
	IsAndOnly() bool
	IsLogical() bool
	ReturnType() uint32
	Attributes() []string
	Eval(map[string]value.Values) (value.Values, uint32, error)
}

type UnaryExtend struct {
	Op int
	E  Extend
}

type BinaryExtend struct {
	Op          int
	Left, Right Extend
}

type MultiExtend struct {
	Op   int
	Args []Extend
}

type ParenExtend struct {
	E Extend
}

type Attribute struct {
	Type uint32
	Name string
}
