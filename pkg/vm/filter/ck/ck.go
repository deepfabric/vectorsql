package ck

import (
	"encoding/binary"
	"errors"
	"reflect"
	"unsafe"

	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
)

func New(cli client.Client, query string) *ck {
	return &ck{cli: cli, query: query}
}

func (c *ck) String() string {
	return c.query
}

func (c *ck) Bitmap(_ int) (*roaring.Bitmap, error) {
	rows, err := c.cli.Query(c.query, "bitmap")
	if err != nil {
		return nil, err
	}
	return unmarshalBitMap([]byte(rows[0]))
}

func unmarshalBitMap(data []byte) (*roaring.Bitmap, error) {
	switch data[0] {
	case 0:
		data = data[1:]
		_, n := binary.Uvarint(data)
		if n < 0 {
			return nil, errors.New("overflow")
		}
		mp := roaring.New()
		mp.AddMany(decodeVector(data[n:]))
		return mp, nil
	case 1:
		data = data[1:]
		_, n := binary.Uvarint(data)
		if n < 0 {
			return nil, errors.New("overflow")
		}
		data = data[n:]
		mp := roaring.New()
		if err := mp.UnmarshalBinary(data); err != nil {
			return nil, err
		}
		return mp, nil
	}
	return nil, nil
}

func decodeVector(v []byte) []uint32 {
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	hp.Len /= 4
	hp.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&hp))
}
