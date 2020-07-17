package value

import (
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewInt(v int64) *Int {
	r := Int(v)
	return &r
}

func (a *Int) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Int) ResolvedType() types.T {
	return types.T_int
}

// ParseInt parses and returns the *Int value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseInt(s string) (*Int, error) {
	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return nil, makeParseError(s, types.T_int, err)
	}
	return NewInt(i), nil
}

func (a *Int) Compare(v Value) int {
	var x, y int64

	x = int64(*a)
	switch b := v.(type) {
	case *Int:
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

func (_ *Int) Size() int            { return 9 }
func (_ *Int) IsLogical() bool      { return false }
func (_ *Int) IsAndOnly() bool      { return true }
func (_ *Int) Attributes() []string { return []string{} }
