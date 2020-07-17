package value

import (
	"errors"
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

func (a *Bool) String() string {
	return strconv.FormatBool(bool(*a))
}

func (_ *Bool) ResolvedType() types.T {
	return types.T_bool
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
	return nil, makeParseError(s, types.T_bool, errors.New("invalid bool value"))
}

func (a *Bool) Compare(v Value) int {
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

func (_ *Bool) Size() int            { return 2 }
func (_ *Bool) IsLogical() bool      { return true }
func (_ *Bool) IsAndOnly() bool      { return true }
func (_ *Bool) Attributes() []string { return []string{} }
