package bv

import (
	"math/rand"
	"time"

	"github.com/RoaringBitmap/roaring"
)

func New() *bv {
	return &bv{}
}

func (b *bv) Fvectors(n int, _ []float32) (*roaring.Bitmap, error) {
	return randVector(n), nil
}

func (b *bv) Vectors(_ *roaring.Bitmap, n int, _ []float32) (*roaring.Bitmap, error) {
	return randVector(n), nil
}

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func randVector(n int) *roaring.Bitmap {
	mp := roaring.New()
	for i := 0; i < n; i++ {
		mp.Add(rand.Uint32() % 1000)
	}
	return mp
}
