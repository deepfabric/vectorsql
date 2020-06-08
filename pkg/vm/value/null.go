package value

import "github.com/deepfabric/vectorsql/pkg/vm/types"

func (_ Null) Size() int                              { return 1 }
func (_ Null) String() string                         { return "null" }
func (_ Null) IsLogical() bool                        { return false }
func (_ Null) IsAndOnly() bool                        { return true }
func (_ Null) ResolvedType() *types.T                 { return types.Null }
func (_ Null) Attributes() []string                   { return []string{} }
func (a Null) Eval(_ map[string]Value) (Value, error) { return a, nil }

func (a Null) Compare(v Value) int {
	if v == ConstNull {
		return 0
	}
	return 1 // NULL is less than any non-NULL value
}
