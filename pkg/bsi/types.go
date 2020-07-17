package bsi

import (
	"github.com/pilosa/pilosa/roaring"
)

type Bsi interface {
	Clone() Bsi
	Map() *roaring.Bitmap

	Read([]byte) error
	Show() ([]byte, error)

	Del(uint64) error
	Set(uint64, interface{}) error
	Get(uint64) (interface{}, bool)

	Count(*roaring.Bitmap) uint64
	Min(*roaring.Bitmap) (interface{}, uint64)
	Max(*roaring.Bitmap) (interface{}, uint64)
	Sum(*roaring.Bitmap) (interface{}, uint64)

	Eq(interface{}) (*roaring.Bitmap, error)
	Ne(interface{}) (*roaring.Bitmap, error)
	Lt(interface{}) (*roaring.Bitmap, error)
	Le(interface{}) (*roaring.Bitmap, error)
	Gt(interface{}) (*roaring.Bitmap, error)
	Ge(interface{}) (*roaring.Bitmap, error)
}
