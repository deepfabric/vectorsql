package bv

import (
	"encoding/binary"

	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/beevector/pkg/sdk"
	"github.com/deepfabric/vectorsql/pkg/logger"
)

func New(addrs []string, log logger.Log) *bv {
	return &bv{sdk.NewClient(addrs...), log}
}

func (b *bv) Add(xbs []float32, xids []int64) error {
	return b.cli.Add(xbs, xids)
}

func (b *bv) Fvectors(n int64, v []float32) (*roaring.Bitmap, []uint64, error) {
	_, vs, err := b.cli.Search(n, v, nil)
	if err != nil {
		return nil, nil, err
	}
	mp, ids := genIds(vs)
	return mp, ids, nil
}

func (b *bv) Vectors(n int64, mp *roaring.Bitmap, v []float32) (*roaring.Bitmap, []uint64, error) {
	if mp != nil {
		data, err := mp.ToBytes()
		if err != nil {
			return nil, nil, err
		}
		buf := make([]byte, 11)
		buf[0] = 1
		num := binary.PutUvarint(buf[1:], uint64(len(data)))
		_, vs, err := b.cli.Search(n, v, append(buf[:1+num], data...))
		if err != nil {
			return nil, nil, err
		}
		mp, ids := genIds(vs)
		return mp, ids, nil
	}
	_, vs, err := b.cli.Search(n, v, nil)
	if err != nil {
		return nil, nil, err
	}
	mp, ids := genIds(vs)
	return mp, ids, nil
}

func genIds(vs []int64) (*roaring.Bitmap, []uint64) {
	xs := make([]uint32, len(vs))
	ys := make([]uint64, len(vs))
	for i, v := range vs {
		xs[i] = uint32(v >> 34)
		ys[i] = uint64(v)
	}
	return roaring.BitmapOf(xs...), ys
}
