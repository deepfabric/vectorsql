package value

import (
	"fmt"
	"math"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewFloat32(v float32) *Float32 {
	r := Float32(v)
	return &r
}

func (a *Float32) String() string {
	f := float32(*a)
	if _, frac := math.Modf(float64(f)); frac == 0 && -1000000 < *a && *a < 1000000 {
		return fmt.Sprintf("%.1f", f)
	} else {
		return fmt.Sprintf("%g", f)
	}
}

func (_ *Float32) ResolvedType() types.T {
	return types.T_float32
}

// ParseFloat32 parses and returns the *Float32 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseFloat32(s string) (*Float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return nil, makeParseError(s, types.T_float32, err)
	}
	return NewFloat32(float32(f)), nil
}

func (a *Float32) Compare(v Value) int {
	var x, y float32

	x = float32(*a)
	switch b := v.(type) {
	case *Float:
		y = float32(*b)
	case *Float32:
		y = float32(*b)
	default:
		panic(makeUnsupportedComparisonMessage(a, v))
	}
	// NaN sorts before non-NaN.
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	case x == y:
		return 0
	}
	if math.IsNaN(float64(x)) {
		if math.IsNaN(float64(y)) {
			return 0
		}
		return -1
	}
	return 1
}

func (_ *Float32) Size() int            { return 9 }
func (_ *Float32) IsLogical() bool      { return false }
func (_ *Float32) IsAndOnly() bool      { return true }
func (_ *Float32) Attributes() []string { return []string{} }
