package value

import (
	"fmt"
	"strings"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func Compare(a, b Value) int {
	if ta, tb := a.ResolvedType().Oid, b.ResolvedType().Oid; (ta == types.T_int || ta == types.T_float) &&
		(tb == types.T_int || tb == types.T_float) {
		return a.Compare(b)
	}
	if r := int(a.ResolvedType().Oid - b.ResolvedType().Oid); r != 0 {
		return r
	}
	return a.Compare(b)
}

// makeParseError returns a parse error using the provided string and type. An
// optional error can be provided, which will be appended to the end of the
// error string.
func makeParseError(s string, typ *types.T, err error) error {
	if err != nil {
		return fmt.Errorf("could not parse %q as type %s: %v", s, typ, err)
	}
	return fmt.Errorf("could not parse %q as type %s", s, typ)
}

func makeUnsupportedComparisonMessage(d1, d2 Value) error {
	return fmt.Errorf("unsupported comparison: %s to %s", d1.ResolvedType(), d2.ResolvedType())
}

func isCaseInsensitivePrefix(prefix, s string) bool {
	if len(prefix) > len(s) {
		return false
	}
	return strings.EqualFold(prefix, s[:len(prefix)])
}
