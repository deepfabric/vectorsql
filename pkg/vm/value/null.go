package value

import "github.com/deepfabric/vectorsql/pkg/vm/types"

func (_ Null) Size() int             { return 1 }
func (_ Null) String() string        { return "null" }
func (_ Null) IsLogical() bool       { return false }
func (_ Null) IsAndOnly() bool       { return true }
func (_ Null) ResolvedType() types.T { return types.T_null }
func (_ Null) Attributes() []string  { return []string{} }

func (_ Null) ReturnType() uint32 {
	return types.T_null
}

func (_ Null) Compare(v Value) int {
	if v == ConstNull {
		return 0
	}
	return 1 // NULL is less than any non-NULL value
}

func (a Null) Eval(_ map[string]Values) (Values, uint32, error) {
	return nil, 0, nil
}
