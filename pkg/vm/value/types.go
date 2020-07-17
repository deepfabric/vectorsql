package value

import (
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/pilosa/pilosa/roaring"
)

type Values interface {
	Size() int

	Show() ([]byte, error)
	Read(int, []byte) error

	Count() int
	Slice() ([]uint64, [][]byte)

	Filter([]uint64) interface{}
	MergeFilter(interface{}) interface{}

	MarkNull(int) error
	Append(interface{}) error
	Update([]int, interface{}) error
	Merge(*roaring.Bitmap, *roaring.Bitmap) error
}

type Value interface {
	Size() int

	ToValues() Values

	String() string
	Compare(Value) int
	ResolvedType() types.T

	IsLogical() bool
	IsAndOnly() bool
	ReturnType() uint32
	Attributes() []string
	Eval(map[string]Values) (Values, uint32, error)
}

type Bool bool // true = 1, false = 0
type Int int64
type Float float64
type String string

type Uint8 uint8
type Uint16 uint16
type Uint32 uint32
type Uint64 uint64

type Int8 int8
type Int16 int16
type Int32 int32
type Int64 int64

type Float32 float32
type Float64 float64

type Null struct{}
type Array []Value

type Data uint32
type Time uint32 // 0 ~ 24 * 3600
type Timestamp int64

var (
	ConstTrue  Bool = true
	ConstFalse Bool = false
	ConstNull  Null = Null{}
)

// time.Time formats.
const (
	TimeOutputFormat = "00:00:00"

	DataOutputFormat = "2000-00-00"

	TimestampOutputFormat = "2000-00-00 00:00:00"
)
