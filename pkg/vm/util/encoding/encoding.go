package encoding

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"unsafe"
)

func Encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
}

func EncodeUint32(v uint32) []byte {
	hp := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&v)), Len: 4, Cap: 4}
	return *(*[]byte)(unsafe.Pointer(&hp))
}

func DecodeUint32(v []byte) uint32 {
	return *(*uint32)(unsafe.Pointer(&v[0]))
}

func EncodeUint64(v uint64) []byte {
	hp := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&v)), Len: 8, Cap: 8}
	return *(*[]byte)(unsafe.Pointer(&hp))
}

func DecodeUint64(v []byte) uint64 {
	return *(*uint64)(unsafe.Pointer(&v[0]))
}

func EncodeUint32Slice(v []uint32) []byte {
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	hp.Len *= 4
	hp.Cap *= 4
	return *(*[]byte)(unsafe.Pointer(&hp))
}

func DecodeUint32Slice(v []byte) []uint32 {
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	hp.Len /= 4
	hp.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&hp))
}
