package bv

import "github.com/RoaringBitmap/roaring"

type BV interface {
	Fvectors(int, []float32) (*roaring.Bitmap, error)
	Vectors(*roaring.Bitmap, int, []float32) (*roaring.Bitmap, error)
}

type bv struct {
}
