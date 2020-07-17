package bv

import (
	"encoding/binary"

	"github.com/RoaringBitmap/roaring"
)

func New() *bv {
	return &bv{}
}

func (b *bv) Fvectors(n int64, v []float32) ([]int64, error) {
	_, ids, err := b.cli.Search(n, v, nil)
	return ids, err
}

func (b *bv) Vectors(n int64, mp *roaring.Bitmap, v []float32) ([]int64, error) {
	data, err := mp.ToBytes()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 11)
	buf[0] = 1
	n := binary.PutUvarint(buf[1:], uint64(len(data)))
	_, ids, err := b.cli.Search(n, v, append(buf[:1+n], data...))
	return ids, err
}
