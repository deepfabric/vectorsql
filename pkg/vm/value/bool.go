package value

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewBool(v bool) *Bool {
	if v {
		return &ConstTrue
	}
	return &ConstFalse
}

func AsBool(v interface{}) (Bool, bool) {
	switch t := v.(type) {
	case *Bool:
		return *t, true
	default:
		return false, false
	}
}

// MustBeBool attempts to retrieve a Bool from a value, panicking if the
// assertion fails.
func MustBeBool(v interface{}) bool {
	b, ok := AsBool(v)
	if !ok {
		panic(fmt.Errorf("expected *Bool, found %T", v))
	}
	return bool(b)
}

// GetBool get Bool or an error.
func GetBool(v Value) (Bool, error) {
	if b, ok := v.(*Bool); ok {
		return *b, nil
	}
	return false, fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.Bool)
}

func (a *Bool) String() string {
	return strconv.FormatBool(bool(*a))
}

func (_ *Bool) ResolvedType() *types.T {
	return types.Bool
}

func ParseBool(s string) (*Bool, error) {
	s = strings.TrimSpace(s)
	if len(s) >= 1 {
		switch s[0] {
		case 't', 'T':
			if isCaseInsensitivePrefix(s, "true") {
				return &ConstTrue, nil
			}
		case 'f', 'F':
			if isCaseInsensitivePrefix(s, "false") {
				return &ConstFalse, nil
			}
		}
	}
	return nil, makeParseError(s, types.Bool, errors.New("invalid bool value"))
}

func (a *Bool) Compare(v Value) int {
	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
	b, ok := v.(*Bool)
	if !ok {
		panic(makeUnsupportedComparisonMessage(a, v))
	}
	return CompareBool(bool(*a), bool(*b))
}

// CompareBool compare the input bools according to the SQL comparison rules.
func CompareBool(d, v bool) int {
	if !d && v {
		return -1
	}
	if d && !v {
		return 1
	}
	return 0
}

func (_ *Bool) Size() int                              { return 2 }
func (_ *Bool) IsLogical() bool                        { return true }
func (_ *Bool) IsAndOnly() bool                        { return true }
func (_ *Bool) Attributes() []string                   { return []string{} }
func (a *Bool) Eval(_ map[string]Value) (Value, error) { return a, nil }
