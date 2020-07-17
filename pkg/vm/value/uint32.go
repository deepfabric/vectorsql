package value

import (
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewUint32(v uint32) *Uint32 {
	r := Uint32(v)
	return &r
}

func (a *Uint32) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Uint32) ResolvedType() types.T {
	return types.T_uint32
}

// ParseUint32 parses and returns the *Uint32 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseUint32(s string) (*Uint32, error) {
	i, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return nil, makeParseError(s, types.T_uint32, err)
	}
	return NewUint32(uint32(i)), nil
}

func (a *Uint32) Compare(v Value) int {
	var x, y uint32

	x = uint32(*a)
	switch b := v.(type) {
	case *Int:
		y = uint32(*b)
	case *Uint32:
		y = uint32(*b)
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

func (_ *Uint32) Size() int            { return 2 }
func (_ *Uint32) IsLogical() bool      { return false }
func (_ *Uint32) IsAndOnly() bool      { return true }
func (_ *Uint32) Attributes() []string { return []string{} }
