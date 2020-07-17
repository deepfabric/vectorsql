package value

import "github.com/deepfabric/vectorsql/pkg/vm/types"

func (a *Bool) ReturnType() uint32 {
	return types.T_bool
}

func (a *Int) ReturnType() uint32 {
	return types.T_int
}

func (a *Int8) ReturnType() uint32 {
	return types.T_int8
}

func (a *Int16) ReturnType() uint32 {
	return types.T_int16
}

func (a *Int32) ReturnType() uint32 {
	return types.T_int32
}

func (a *Int64) ReturnType() uint32 {
	return types.T_int64
}

func (a *Uint8) ReturnType() uint32 {
	return types.T_uint8
}

func (a *Uint16) ReturnType() uint32 {
	return types.T_uint16
}

func (a *Uint32) ReturnType() uint32 {
	return types.T_uint32
}

func (a *Uint64) ReturnType() uint32 {
	return types.T_uint64
}

func (a *Float) ReturnType() uint32 {
	return types.T_float
}

func (a *Float32) ReturnType() uint32 {
	return types.T_float32
}

func (a *Float64) ReturnType() uint32 {
	return types.T_float64
}

func (a *String) ReturnType() uint32 {
	return types.T_string
}

func (a *Timestamp) ReturnType() uint32 {
	return types.T_timestamp
}
