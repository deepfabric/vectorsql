package bv

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/beevector/pkg/sdk"
	"github.com/deepfabric/vectorsql/pkg/logger"
)

type BV interface {
	Add([]float32, []int64) error
	Fvectors(int64, []float32) (*roaring.Bitmap, []uint64, error)
	Vectors(int64, *roaring.Bitmap, []float32) (*roaring.Bitmap, []uint64, error)
}

type bv struct {
	cli sdk.Client
	log logger.Log
}
