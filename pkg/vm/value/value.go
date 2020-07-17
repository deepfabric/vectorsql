package value

import (
	"fmt"
	"strings"
	"time"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/value/dynamic"
	"github.com/deepfabric/vectorsql/pkg/vm/value/static"
)

func NewValues(v interface{}) Values {
	switch a := v.(type) {
	case bool:
		return &static.Bools{Vs: []bool{a}}
	case int:
		return &static.Ints{Vs: []int64{int64(a)}}
	case int8:
		return &static.Int8s{Vs: []int8{a}}
	case int16:
		return &static.Int16s{Vs: []int16{a}}
	case int32:
		return &static.Int32s{Vs: []int32{a}}
	case int64:
		return &static.Int64s{Vs: []int64{a}}
	case uint8:
		return &static.Uint8s{Vs: []uint8{a}}
	case uint16:
		return &static.Uint16s{Vs: []uint16{a}}
	case uint32:
		return &static.Uint32s{Vs: []uint32{a}}
	case uint64:
		return &static.Uint64s{Vs: []uint64{a}}
	case float32:
		return &static.Float32s{Vs: []float32{a}}
	case float64:
		return &static.Float64s{Vs: []float64{a}}
	case string:
		return &dynamic.Strings{Vs: []string{a}}
	}
	return nil
}

// makeParseError returns a parse error using the provided string and type. An
// optional error can be provided, which will be appended to the end of the
// error string.
func makeParseError(s string, typ types.T, err error) error {
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

func MustBeBool(v interface{}) bool {
	return bool(*(v.(*Bool)))
}

func MustBeInt(v interface{}) int64 {
	return int64(*(v.(*Int)))
}

func MustBeInt8(v interface{}) int8 {
	return int8(*(v.(*Int8)))
}

func MustBeInt16(v interface{}) int16 {
	return int16(*(v.(*Int16)))
}

func MustBeInt32(v interface{}) int32 {
	return int32(*(v.(*Int32)))
}

func MustBeInt64(v interface{}) int64 {
	return int64(*(v.(*Int)))
}

func MustBeUint8(v interface{}) uint8 {
	return uint8(*(v.(*Uint8)))
}

func MustBeUint16(v interface{}) uint16 {
	return uint16(*(v.(*Uint16)))
}

func MustBeUint32(v interface{}) uint32 {
	return uint32(*(v.(*Uint32)))
}

func MustBeUint64(v interface{}) uint64 {
	return uint64(*(v.(*Uint64)))
}

func MustBeFloat(v interface{}) float64 {
	return float64(*(v.(*Float)))
}

func MustBeFloat32(v interface{}) float32 {
	return float32(*(v.(*Float32)))
}

func MustBeFloat64(v interface{}) float64 {
	return float64(*(v.(*Float64)))
}

func MustBeString(v interface{}) string {
	return string(*(v.(*String)))
}

func MustBeTimestamp(v interface{}) time.Time {
	return time.Unix(int64(*(v.(*Timestamp))), 0)
}
