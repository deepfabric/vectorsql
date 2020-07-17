package value

import (
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewInt64(v int64) *Int64 {
	r := Int64(v)
	return &r
}

func (a *Int64) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Int64) ResolvedType() types.T {
	return types.T_int64
}

// ParseInt64 parses and returns the *Int64 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseInt64(s string) (*Int64, error) {
	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return nil, makeParseError(s, types.T_int64, err)
	}
	return NewInt64(int64(i)), nil
}

func (a *Int64) Compare(v Value) int {
	var x, y int64

	x = int64(*a)
	switch b := v.(type) {
	case *Int:
		y = int64(*b)
	case *Int64:
		y = int64(*b)
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

func (_ *Int64) Size() int            { return 2 }
func (_ *Int64) IsLogical() bool      { return false }
func (_ *Int64) IsAndOnly() bool      { return true }
func (_ *Int64) Attributes() []string { return []string{} }
