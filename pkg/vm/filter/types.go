package filter

import "github.com/RoaringBitmap/roaring"

type Filter interface {
	String() string
	Bitmap() (*roaring.Bitmap, error)
}
