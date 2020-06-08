package value

import (
	"fmt"
	"time"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewTime(t time.Time) *Time {
	return &Time{Time: t.Round(time.Second)}
}

func AsTime(v interface{}) (Time, bool) {
	switch t := v.(type) {
	case *Time:
		return *t, true
	default:
		return Time{}, false
	}
}

func MustBeTime(v interface{}) time.Time {
	t, ok := AsTime(v)
	if !ok {
		panic(fmt.Errorf("expected *Time, found %T", v))
	}
	return t.Time
}

func GetTime(v Value) (Time, error) {
	if t, ok := v.(*Time); ok {
		return *t, nil
	}
	return Time{}, fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.Time)
}

func (a *Time) String() string {
	return a.UTC().Format(TimeOutputFormat)
}

func (_ *Time) ResolvedType() *types.T {
	return types.Time
}

// ParseTime parses and returns the *Time value represented by
// the provided string in UTC, or an error if parsing is unsuccessful.
func ParseTime(s string) (*Time, error) {
	t, err := time.Parse(TimeOutputFormat, s)
	if err != nil {
		return nil, err
	}
	return NewTime(t), nil
}

func (a *Time) Compare(v Value) int {
	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
	return compareTime(a, v)
}

func compareTime(a, b Value) int {
	aTime, aErr := GetTime(a)
	bTime, bErr := GetTime(b)
	if aErr != nil || bErr != nil {
		panic(makeUnsupportedComparisonMessage(a, b))
	}
	if aTime.Before(bTime.Time) {
		return -1
	}
	if bTime.Before(aTime.Time) {
		return 1
	}
	return 0
}

func (_ *Time) Size() int                              { return 9 }
func (_ *Time) IsLogical() bool                        { return false }
func (_ *Time) IsAndOnly() bool                        { return true }
func (_ *Time) Attributes() []string                   { return []string{} }
func (a *Time) Eval(_ map[string]Value) (Value, error) { return a, nil }
