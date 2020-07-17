package value

import (
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewUint8(v uint8) *Uint8 {
	r := Uint8(v)
	return &r
}

func (a *Uint8) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Uint8) ResolvedType() types.T {
	return types.T_uint8
}

// ParseUint8 parses and returns the *Uint8 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseUint8(s string) (*Uint8, error) {
	i, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return nil, makeParseError(s, types.T_uint8, err)
	}
	return NewUint8(uint8(i)), nil
}

func (a *Uint8) Compare(v Value) int {
	var x, y uint8

	x = uint8(*a)
	switch b := v.(type) {
	case *Int:
		y = uint8(*b)
	case *Uint8:
		y = uint8(*b)
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

func (_ *Uint8) Size() int            { return 2 }
func (_ *Uint8) IsLogical() bool      { return false }
func (_ *Uint8) IsAndOnly() bool      { return true }
func (_ *Uint8) Attributes() []string { return []string{} }
