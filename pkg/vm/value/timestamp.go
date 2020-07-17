package value

import (
	"time"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewTimestamp(v time.Time) *Timestamp {
	r := Timestamp(v.Unix())
	return &r
}

func (a *Timestamp) String() string {
	return time.Unix(int64(*a), 0).UTC().Format(TimestampOutputFormat)
}

func (_ *Timestamp) ResolvedType() types.T {
	return types.T_timestamp
}

// ParseTimestamp parses and returns the *Timestamp value represented by
// the provided string in UTC, or an error if parsing is unsuccessful.
func ParseTimestamp(s string) (*Timestamp, error) {
	t, err := time.Parse(TimestampOutputFormat, s)
	if err != nil {
		return nil, err
	}
	return NewTimestamp(t), nil
}

func (a *Timestamp) Compare(v Value) int {
	var x, y int64

	x = int64(*a)
	switch b := v.(type) {
	case *Timestamp:
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

func (_ *Timestamp) Size() int            { return 9 }
func (_ *Timestamp) IsLogical() bool      { return false }
func (_ *Timestamp) IsAndOnly() bool      { return true }
func (_ *Timestamp) Attributes() []string { return []string{} }
