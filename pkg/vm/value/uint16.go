package value

import (
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewUint16(v uint16) *Uint16 {
	r := Uint16(v)
	return &r
}

func (a *Uint16) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Uint16) ResolvedType() types.T {
	return types.T_uint16
}

// ParseUint16 parses and returns the *Uint16 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseUint16(s string) (*Uint16, error) {
	i, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return nil, makeParseError(s, types.T_uint16, err)
	}
	return NewUint16(uint16(i)), nil
}

func (a *Uint16) Compare(v Value) int {
	var x, y uint16

	x = uint16(*a)
	switch b := v.(type) {
	case *Int:
		y = uint16(*b)
	case *Uint16:
		y = uint16(*b)
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

func (_ *Uint16) Size() int            { return 2 }
func (_ *Uint16) IsLogical() bool      { return false }
func (_ *Uint16) IsAndOnly() bool      { return true }
func (_ *Uint16) Attributes() []string { return []string{} }
