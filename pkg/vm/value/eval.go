package value

import (
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/value/dynamic"
	"github.com/deepfabric/vectorsql/pkg/vm/value/static"
)

func (a *Bool) Eval(mp map[string]Values) (Values, uint32, error) {
	v := bool(*a)
	vs := make([]bool, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewBools(vs, nil, nil), types.T_bool, nil
}

func (a *Int) Eval(mp map[string]Values) (Values, uint32, error) {
	v := int64(*a)
	vs := make([]int64, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewInts(vs, nil, nil), types.T_int, nil
}

func (a *Int8) Eval(mp map[string]Values) (Values, uint32, error) {
	v := int8(*a)
	vs := make([]int8, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewInt8s(vs, nil, nil), types.T_int8, nil
}

func (a *Int16) Eval(mp map[string]Values) (Values, uint32, error) {
	v := int16(*a)
	vs := make([]int16, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewInt16s(vs, nil, nil), types.T_int16, nil
}

func (a *Int32) Eval(mp map[string]Values) (Values, uint32, error) {
	v := int32(*a)
	vs := make([]int32, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewInt32s(vs, nil, nil), types.T_int32, nil
}

func (a *Int64) Eval(mp map[string]Values) (Values, uint32, error) {
	v := int64(*a)
	vs := make([]int64, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewInt64s(vs, nil, nil), types.T_int64, nil
}

func (a *Uint8) Eval(mp map[string]Values) (Values, uint32, error) {
	v := uint8(*a)
	vs := make([]uint8, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewUint8s(vs, nil, nil), types.T_uint8, nil
}

func (a *Uint16) Eval(mp map[string]Values) (Values, uint32, error) {
	v := uint16(*a)
	vs := make([]uint16, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewUint16s(vs, nil, nil), types.T_uint16, nil
}

func (a *Uint32) Eval(mp map[string]Values) (Values, uint32, error) {
	v := uint32(*a)
	vs := make([]uint32, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewUint32s(vs, nil, nil), types.T_uint32, nil
}

func (a *Uint64) Eval(mp map[string]Values) (Values, uint32, error) {
	v := uint64(*a)
	vs := make([]uint64, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewUint64s(vs, nil, nil), types.T_uint64, nil
}

func (a *Float) Eval(mp map[string]Values) (Values, uint32, error) {
	v := float64(*a)
	vs := make([]float64, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewFloats(vs, nil, nil), types.T_float, nil
}

func (a *Float32) Eval(mp map[string]Values) (Values, uint32, error) {
	v := float32(*a)
	vs := make([]float32, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewFloat32s(vs, nil, nil), types.T_float32, nil
}

func (a *Float64) Eval(mp map[string]Values) (Values, uint32, error) {
	v := float64(*a)
	vs := make([]float64, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewFloat64s(vs, nil, nil), types.T_float64, nil
}

func (a *String) Eval(mp map[string]Values) (Values, uint32, error) {
	v := string(*a)
	vs := make([]string, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return dynamic.NewStrings(vs, nil, nil), types.T_string, nil
}

// Not yet supported
func (a Array) Eval(mp map[string]Values) (Values, uint32, error) {
	return nil, types.T_array, nil
}

func (a *Timestamp) Eval(mp map[string]Values) (Values, uint32, error) {
	v := int64(*a)
	vs := make([]int64, length(mp))
	for i := range vs {
		vs[i] = v
	}
	return static.NewTimestamps(vs, nil, nil), types.T_timestamp, nil
}

func length(mp map[string]Values) int {
	for _, v := range mp {
		return v.Count()
	}
	return 0
}
