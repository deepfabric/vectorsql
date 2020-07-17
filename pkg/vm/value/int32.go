package value

import (
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewInt32(v int32) *Int32 {
	r := Int32(v)
	return &r
}

func (a *Int32) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Int32) ResolvedType() types.T {
	return types.T_int32
}

// ParseInt32 parses and returns the *Int32 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseInt32(s string) (*Int32, error) {
	i, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return nil, makeParseError(s, types.T_int32, err)
	}
	return NewInt32(int32(i)), nil
}

func (a *Int32) Compare(v Value) int {
	var x, y int32

	x = int32(*a)
	switch b := v.(type) {
	case *Int:
		y = int32(*b)
	case *Int32:
		y = int32(*b)
	default:
		panic(makeUnsupportedComparisonMessage(a, v))
	}
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	default:
		return 0
	}
}

func (_ *Int32) Size() int            { return 2 }
func (_ *Int32) IsLogical() bool      { return false }
func (_ *Int32) IsAndOnly() bool      { return true }
func (_ *Int32) Attributes() []string { return []string{} }
