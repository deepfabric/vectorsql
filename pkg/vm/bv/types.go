package bv

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/beevector/pkg/sdk"
)

type BV interface {
	Fvectors(int64, []float32) ([]int64, error)
	Vectors(int64, *roaring.Bitmap, []float32) ([]int64, error)
}

type bv struct {
	cli sdk.Client
}
