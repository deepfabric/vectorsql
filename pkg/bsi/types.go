package bsi

import (
	"github.com/RoaringBitmap/roaring"
)

type Bsi interface {
	Clone() Bsi
	Map() *roaring.Bitmap

	Read([]byte) error
	Show() ([]byte, error)

	Del(uint32) error
	Set(uint32, interface{}) error
	Get(uint32) (interface{}, bool)

	Eq(interface{}) (*roaring.Bitmap, error)
	Ne(interface{}) (*roaring.Bitmap, error)
	Lt(interface{}) (*roaring.Bitmap, error)
	Le(interface{}) (*roaring.Bitmap, error)
	Gt(interface{}) (*roaring.Bitmap, error)
	Ge(interface{}) (*roaring.Bitmap, error)
}
