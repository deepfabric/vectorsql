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

func AsFloat64(v interface{}) (Float64, bool) {
	switch t := v.(type) {
	case *Float64:
		return *t, true
	default:
		return 0.0, false
	}
}

func MustBeFloat64(v interface{}) float64 {
	f, ok := AsFloat64(v)
	if !ok {
		panic(fmt.Errorf("expected *Float64, found %T", v))
	}
	return float64(f)
}

func GetFloat64(v Value) (Float64, error) {
	if f, ok := v.(*Float64); ok {
		return *f, nil
	}
	return 0, fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.Float64)
}

func (a *Float64) String() string {
	f := float64(*a)
	if _, frac := math.Modf(f); frac == 0 && -1000000 < *a && *a < 1000000 {
		return fmt.Sprintf("%.1f", f)
	} else {
		return fmt.Sprintf("%g", f)
	}
}

func (_ *Float64) ResolvedType() *types.T {
	return types.Float64
}

// ParseFloat64 parses and returns the *Float64 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseFloat64(s string) (*Float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, makeParseError(s, types.Float64, err)
	}
	return NewFloat64(f), nil
}

func (a *Float64) Compare(v Value) int {
	var x, y float64

	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
	x = float64(*a)
	switch b := v.(type) {
	case *Float:
		y = float64(*b)
	case *Float64:
		y = float64(*b)
	default:
		panic(makeUnsupportedComparisonMessage(a, v))
	}
	// NaN sorts before non-NaN (#10109).
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

func (_ *Float64) Size() int                              { return 9 }
func (_ *Float64) IsLogical() bool                        { return false }
func (_ *Float64) IsAndOnly() bool                        { return true }
func (_ *Float64) Attributes() []string                   { return []string{} }
func (a *Float64) Eval(_ map[string]Value) (Value, error) { return a, nil }
