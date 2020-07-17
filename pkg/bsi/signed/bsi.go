package signed

import (
	"bytes"

	Bsi "github.com/deepfabric/vectorsql/pkg/bsi"
	"github.com/deepfabric/vectorsql/pkg/vm/util/encoding"
	"github.com/pilosa/pilosa/roaring"
)

func New(bitSize int) *bsi {
	ms := make([]*roaring.Bitmap, bitSize+2)
	for i, j := 0, bitSize+2; i < j; i++ {
		ms[i] = roaring.NewBitmap()
	}
	return &bsi{
		ms:      ms,
		bitSize: bitSize,
	}
}

func (b *bsi) Map() *roaring.Bitmap {
	return b.subMap(bsiExistsBit)
}

func (b *bsi) Clone() Bsi.Bsi {
	ms := make([]*roaring.Bitmap, b.bitSize+2)
	for i, j := 0, b.bitSize+2; i < j; i++ {
		ms[i] = b.ms[i].Clone()
	}
	return &bsi{
		ms:      ms,
		bitSize: b.bitSize,
	}
}

func (b *bsi) Show() ([]byte, error) {
	var body []byte
	var buf bytes.Buffer

	os := make([]uint32, 0, len(b.ms))
	buf.WriteByte(byte(b.bitSize & 0xFF))
	for _, m := range b.ms {
		data, err := show(m)
		if err != nil {
			return nil, err
		}
		os = append(os, uint32(len(body)))
		body = append(body, data...)
	}
	{
		data := encoding.EncodeUint32Slice(os)
		buf.Write(encoding.EncodeUint32(uint32(len(data))))
		buf.Write(data)
	}
	buf.Write(body)
	return buf.Bytes(), nil
}

func (b *bsi) Read(data []byte) error {
	b.bitSize = int(data[0])
	data = data[1:]
	n := encoding.DecodeUint32(data[:4])
	data = data[4:]
	os := encoding.DecodeUint32Slice(data[:n])
	data = data[n:]
	b.ms = make([]*roaring.Bitmap, b.bitSize+2)
	for i, j := 0, b.bitSize+2; i < j; i++ {
		b.ms[i] = roaring.NewBitmap()
		if i < j-1 {
			if err := b.ms[i].UnmarshalBinary(data[os[i]:os[i+1]]); err != nil {
				return err
			}
		} else {
			if err := b.ms[i].UnmarshalBinary(data[os[i]:]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *bsi) Count(filter *roaring.Bitmap) uint64 {
	mp := b.subMap(bsiExistsBit)
	if filter != nil {
		mp = mp.Intersect(filter)
	}
	return mp.Count()
}

func (b *bsi) Min(filter *roaring.Bitmap) (interface{}, uint64) {
	mp := b.subMap(bsiExistsBit)
	if filter != nil {
		mp = mp.Intersect(filter)
	}
	if mp.Count() == 0 {
		return 0, 0
	}
	if mq := b.subMap(bsiSignBit).Intersect(mp); mq.Any() {
		min, count := b.max(mq)
		return -min, count
	}
	return b.min(mp)
}

func (b *bsi) Max(filter *roaring.Bitmap) (interface{}, uint64) {
	mp := b.subMap(bsiExistsBit)
	if filter != nil {
		mp = mp.Intersect(filter)
	}
	if !mp.Any() {
		return 0, 0
	}
	mq := mp.Difference(b.subMap(bsiSignBit))
	if !mq.Any() {
		max, count := b.min(mp)
		return -max, count
	}
	return b.max(mq)
}

func (b *bsi) Sum(filter *roaring.Bitmap) (interface{}, uint64) {
	var sum int64

	mp := b.subMap(bsiExistsBit)
	if filter != nil {
		mp = mp.Intersect(filter)
	}
	count := mp.Count()
	nmp := b.subMap(bsiSignBit)
	pmp := mp.Difference(nmp)
	for i, j := uint(0), uint(b.bitSize); i < j; i++ {
		mq := b.subMap(uint64(bsiOffsetBit + i))
		n := int64((1 << i) * mq.IntersectionCount(pmp))
		m := int64((1 << i) * mq.IntersectionCount(nmp))
		sum += n - m
	}
	return sum, count
}

func (b *bsi) Get(k uint64) (interface{}, bool) {
	var v int64

	if !b.bit(bsiExistsBit, k) {
		return -1, false
	}
	for i, j := uint(0), uint(b.bitSize); i < j; i++ {
		if b.bit(uint64(bsiOffsetBit+i), k) {
			v |= (1 << i)
		}
	}
	if b.bit(bsiSignBit, k) {
		v = -v
	}
	return v, true
}

func (b *bsi) Set(k uint64, e interface{}) error {
	v := e.(int64)
	uv := uint64(v)
	if v < 0 {
		uv = uint64(-v)
	}
	for i, j := uint(0), uint(b.bitSize); i < j; i++ {
		if uv&(1<<i) != 0 {
			b.setBit(uint64(bsiOffsetBit+i), k)
		} else {
			b.clearBit(uint64(bsiOffsetBit+i), k)
		}
	}
	b.setBit(bsiExistsBit, k)
	if v < 0 {
		b.setBit(bsiSignBit, k)
	} else {
		b.clearBit(bsiSignBit, k)
	}
	return nil
}

func (b *bsi) Del(k uint64) error {
	b.clearBit(bsiExistsBit, k)
	return nil
}

func (b *bsi) Eq(e interface{}) (*roaring.Bitmap, error) {
	v := e.(int64)
	uv := uint64(v)
	mp := b.subMap(bsiExistsBit)
	if v < 0 {
		uv = uint64(-v)
		mp = mp.Intersect(b.subMap(bsiSignBit))
	} else {
		mp = mp.Difference(b.subMap(bsiSignBit))
	}
	for i := b.bitSize - 1; i >= 0; i-- {
		if (uv>>uint(i))&1 == 1 {
			mp = mp.Intersect(b.subMap(uint64(bsiOffsetBit + i)))
		} else {
			mp = mp.Difference(b.subMap(uint64(bsiOffsetBit + i)))
		}
	}
	return mp, nil
}

func (b *bsi) Ne(e interface{}) (*roaring.Bitmap, error) {
	v := e.(int64)
	uv := uint64(v)
	mp := b.subMap(bsiExistsBit)
	if v < 0 {
		uv = uint64(-v)
		mp = mp.Intersect(b.subMap(bsiSignBit))
	} else {
		mp = mp.Difference(b.subMap(bsiSignBit))
	}
	for i := b.bitSize - 1; i >= 0; i-- {
		if (uv>>uint(i))&1 == 1 {
			mp = mp.Intersect(b.subMap(uint64(bsiOffsetBit + i)))
		} else {
			mp = mp.Difference(b.subMap(uint64(bsiOffsetBit + i)))
		}
	}
	return b.subMap(bsiExistsBit).Difference(mp), nil
}

func (b *bsi) Lt(e interface{}) (*roaring.Bitmap, error) {
	v := e.(int64)
	uv := uint64(v)
	if v < 0 {
		uv = uint64(-v)
	}
	mp := b.subMap(bsiExistsBit)
	if v >= 0 {
		mp = b.lt(uv, mp.Difference(b.subMap(bsiSignBit)), false)
		return b.subMap(bsiSignBit).Union(mp), nil
	}
	return b.gt(uv, mp.Intersect(b.subMap(bsiSignBit)), false), nil
}

func (b *bsi) Le(e interface{}) (*roaring.Bitmap, error) {
	v := e.(int64)
	uv := uint64(v)
	if v < 0 {
		uv = uint64(-v)
	}
	mp := b.subMap(bsiExistsBit)
	if v >= 0 {
		mp = b.lt(uv, mp.Difference(b.subMap(bsiSignBit)), true)
		return b.subMap(bsiSignBit).Union(mp), nil
	}
	return b.gt(uv, mp.Intersect(b.subMap(bsiSignBit)), true), nil
}

func (b *bsi) Gt(e interface{}) (*roaring.Bitmap, error) {
	v := e.(int64)
	uv := uint64(v)
	if v < 0 {
		uv = uint64(-v)
	}
	mp := b.subMap(bsiExistsBit)
	if v >= 0 {
		return b.gt(uv, mp.Difference(b.subMap(bsiSignBit)), false), nil
	}
	mq := b.lt(uv, mp.Intersect(b.subMap(bsiSignBit)), false)
	return mp.Difference(b.subMap(bsiSignBit)).Union(mq), nil
}

func (b *bsi) Ge(e interface{}) (*roaring.Bitmap, error) {
	v := e.(int64)
	uv := uint64(v)
	if v < 0 {
		uv = uint64(-v)
	}
	mp := b.subMap(bsiExistsBit)
	if v >= 0 {
		return b.gt(uv, mp.Difference(b.subMap(uint64(bsiSignBit))), true), nil
	}
	mq := b.lt(uv, mp.Intersect(b.subMap(uint64(bsiSignBit))), true)
	return mp.Difference(b.subMap(uint64(bsiSignBit))).Union(mq), nil
}

func (b *bsi) min(filter *roaring.Bitmap) (int64, uint64) {
	var min int64
	var count uint64

	for i := b.bitSize - 1; i >= 0; i-- {
		mp := filter.Difference(b.subMap(uint64(bsiOffsetBit + i)))
		count = mp.Count()
		if count > 0 {
			filter = mp
		} else {
			min += (1 << uint(i))
			if i == 0 {
				count = filter.Count()
			}
		}
	}
	return min, count
}

func (b *bsi) max(filter *roaring.Bitmap) (int64, uint64) {
	var max int64
	var count uint64

	for i := b.bitSize - 1; i >= 0; i-- {
		mp := b.subMap(uint64(bsiOffsetBit + i)).Intersect(filter)
		count = mp.Count()
		if count > 0 {
			max += (1 << uint(i))
			filter = mp
		} else if i == 0 {
			count = filter.Count()
		}
	}
	return max, count
}

func (b *bsi) lt(v uint64, mp *roaring.Bitmap, eq bool) *roaring.Bitmap {
	zflg := true // leading zero flag
	mq := roaring.NewBitmap()
	for i := b.bitSize - 1; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		if zflg {
			if bit == 0 {
				mp = mp.Difference(b.subMap(uint64(bsiOffsetBit + i)))
				continue
			} else {
				zflg = false
			}
		}
		if i == 0 && !eq {
			if bit == 0 {
				return mq
			}
			return mp.Difference(b.subMap(uint64(bsiOffsetBit + i)).Difference(mq))
		}
		if bit == 0 {
			mp = mp.Difference(b.subMap(uint64(bsiOffsetBit + i)).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Difference(b.subMap(uint64(bsiOffsetBit + i))))
		}
	}
	return mp
}

func (b *bsi) gt(v uint64, mp *roaring.Bitmap, eq bool) *roaring.Bitmap {
	mq := roaring.NewBitmap()
	for i := b.bitSize - 1; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		if i == 0 && !eq {
			if bit == 1 {
				return mq
			}
			return mp.Difference(mp.Difference(b.subMap(uint64(bsiOffsetBit + i))).Difference(mq))
		}
		if bit == 1 {
			mp = mp.Difference(mp.Difference(b.subMap(uint64(bsiOffsetBit + i))).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Intersect(b.subMap(uint64(bsiOffsetBit + i))))
		}
	}
	return mp
}

// x is the bit offset, y is row id
func (b *bsi) setBit(x, y uint64)   { b.ms[x].Add(y) }
func (b *bsi) clearBit(x, y uint64) { b.ms[x].Remove(y) }
func (b *bsi) bit(x, y uint64) bool { return b.ms[x].Contains(y) }

func (b *bsi) subMap(x uint64) *roaring.Bitmap { return b.ms[x] }

func show(mp *roaring.Bitmap) ([]byte, error) {
	var buf bytes.Buffer

	if _, err := mp.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
