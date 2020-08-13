package types

import "fmt"

const (
	T_null = iota
	T_int  // int64
	T_bool
	T_float // float64
	T_array
	T_string
	T_timestamp
	T_uint8
	T_uint16
	T_uint32
	T_uint64
	T_int8
	T_int16
	T_int32
	T_int64
	T_float32
	T_float64
)

type T uint32

func (t T) Size() int {
	switch t {
	case T_null:
		return 1
	case T_int:
		return 8
	case T_bool:
		return 1
	case T_timestamp:
		return 8
	case T_float:
		return 8
	case T_string:
		return 10
	case T_uint8:
		return 1
	case T_uint16:
		return 2
	case T_uint32:
		return 4
	case T_uint64:
		return 8
	case T_int8:
		return 1
	case T_int16:
		return 2
	case T_int32:
		return 4
	case T_int64:
		return 8
	case T_float32:
		return 4
	case T_float64:
		return 8
	}
	return -1
}

func (t T) String() string {
	switch t {
	case T_null:
		return "NULL"
	case T_int:
		return "INT"
	case T_bool:
		return "BOOL"
	case T_timestamp:
		return "TIMESTAMP"
	case T_float:
		return "FLOAT"
	case T_array:
		return "ARRAY"
	case T_string:
		return "STRING"
	case T_uint8:
		return "UINT8"
	case T_uint16:
		return "UINT16"
	case T_uint32:
		return "UINT32"
	case T_uint64:
		return "UINT64"
	case T_int8:
		return "INT8"
	case T_int16:
		return "INT16"
	case T_int32:
		return "INT32"
	case T_int64:
		return "INT64"
	case T_float32:
		return "FLOAT32"
	case T_float64:
		return "FLOAT64"
	}
	panic(fmt.Errorf("unexpected oid: %d", t))
}
