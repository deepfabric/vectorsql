package value

import (
	"fmt"
	"math"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewFloat64(v float64) *Float64 {
	r := Float64(v)
	return &r
}

func (a *Float64) String() string {
	f := float64(*a)
	if _, frac := math.Modf(f); frac == 0 && -1000000 < *a && *a < 1000000 {
		return fmt.Sprintf("%.1f", f)
	} else {
		return fmt.Sprintf("%g", f)
	}
}

func (_ *Float64) ResolvedType() types.T {
	return types.T_float64
}

// ParseFloat64 parses and returns the *Float64 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseFloat64(s string) (*Float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, makeParseError(s, types.T_float64, err)
	}
	return NewFloat64(f), nil
}

func (a *Float64) Compare(v Value) int {
	var x, y float64

	x = float64(*a)
	switch b := v.(type) {
	case *Float:
		y = float64(*b)
	case *Float64:
		y = float64(*b)
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
	if math.IsNaN(x) {
		if math.IsNaN(y) {
			return 0
		}
		return -1
	}
	return 1
}

func (_ *Float64) Size() int            { return 9 }
func (_ *Float64) IsLogical() bool      { return false }
func (_ *Float64) IsAndOnly() bool      { return true }
func (_ *Float64) Attributes() []string { return []string{} }
