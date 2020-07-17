package value

import (
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewUint64(v uint64) *Uint64 {
	r := Uint64(v)
	return &r
}

func (a *Uint64) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Uint64) ResolvedType() types.T {
	return types.T_uint64
}

// ParseUint64 parses and returns the *Uint64 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseUint64(s string) (*Uint64, error) {
	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return nil, makeParseError(s, types.T_uint64, err)
	}
	return NewUint64(uint64(i)), nil
}

func (a *Uint64) Compare(v Value) int {
	var x, y uint64

	x = uint64(*a)
	switch b := v.(type) {
	case *Int:
		y = uint64(*b)
	case *Uint64:
		y = uint64(*b)
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

func (_ *Uint64) Size() int            { return 2 }
func (_ *Uint64) IsLogical() bool      { return false }
func (_ *Uint64) IsAndOnly() bool      { return true }
func (_ *Uint64) Attributes() []string { return []string{} }
