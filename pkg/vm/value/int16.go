package value

import (
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewInt16(v int16) *Int16 {
	r := Int16(v)
	return &r
}

func (a *Int16) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Int16) ResolvedType() types.T {
	return types.T_int16
}

// ParseInt16 parses and returns the *Int16 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseInt16(s string) (*Int16, error) {
	i, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return nil, makeParseError(s, types.T_int16, err)
	}
	return NewInt16(int16(i)), nil
}

func (a *Int16) Compare(v Value) int {
	var x, y int16

	x = int16(*a)
	switch b := v.(type) {
	case *Int:
		y = int16(*b)
	case *Int16:
		y = int16(*b)
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

func (_ *Int16) Size() int            { return 2 }
func (_ *Int16) IsLogical() bool      { return false }
func (_ *Int16) IsAndOnly() bool      { return true }
func (_ *Int16) Attributes() []string { return []string{} }
