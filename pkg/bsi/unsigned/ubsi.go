package unsigned

import (
	"bytes"

	Bsi "github.com/deepfabric/vectorsql/pkg/bsi"
	"github.com/deepfabric/vectorsql/pkg/vm/util/encoding"
	"github.com/pilosa/pilosa/roaring"
)

func New(bitSize int) *ubsi {
	ms := make([]*roaring.Bitmap, bitSize+1)
	for i := 0; i <= bitSize; i++ {
		ms[i] = roaring.NewBitmap()
	}
	return &ubsi{
		ms:      ms,
		bitSize: bitSize,
	}
}

func (u *ubsi) Map() *roaring.Bitmap {
	return u.subMap(bsiExistsBit)
}

func (u *ubsi) Clone() Bsi.Bsi {
	ms := make([]*roaring.Bitmap, u.bitSize+1)
	for i := 0; i <= u.bitSize; i++ {
		ms[i] = u.ms[i].Clone()
	}
	return &ubsi{
		ms:      ms,
		bitSize: u.bitSize,
	}
}

func (u *ubsi) Show() ([]byte, error) {
	var body []byte
	var buf bytes.Buffer

	os := make([]uint32, 0, len(u.ms))
	buf.WriteByte(byte(u.bitSize & 0xFF))
	for _, m := range u.ms {
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

func (u *ubsi) Read(data []byte) error {
	u.bitSize = int(data[0])
	data = data[1:]
	n := encoding.DecodeUint32(data[:4])
	data = data[4:]
	os := encoding.DecodeUint32Slice(data[:n])
	data = data[n:]
	u.ms = make([]*roaring.Bitmap, u.bitSize+1)
	for i := 0; i <= u.bitSize; i++ {
		u.ms[i] = roaring.NewBitmap()
		if i < u.bitSize {
			if err := u.ms[i].UnmarshalBinary(data[os[i]:os[i+1]]); err != nil {
				return err
			}
		} else {
			if err := u.ms[i].UnmarshalBinary(data[os[i]:]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *ubsi) Count(filter *roaring.Bitmap) uint64 {
	mp := u.subMap(bsiExistsBit)
	if filter != nil {
		mp = mp.Intersect(filter)
	}
	return mp.Count()
}

func (u *ubsi) Min(filter *roaring.Bitmap) (interface{}, uint64) {
	mp := u.subMap(bsiExistsBit)
	if filter != nil {
		mp = mp.Intersect(filter)
	}
	if mp.Count() == 0 {
		return 0, 0
	}
	return u.min(mp)
}

func (u *ubsi) Max(filter *roaring.Bitmap) (interface{}, uint64) {
	mp := u.subMap(bsiExistsBit)
	if filter != nil {
		mp = mp.Intersect(filter)
	}
	if !mp.Any() {
		return 0, 0
	}
	return u.max(mp)
}

func (u *ubsi) Sum(filter *roaring.Bitmap) (interface{}, uint64) {
	var sum uint64

	mp := u.subMap(bsiExistsBit)
	if filter != nil {
		mp = mp.Intersect(filter)
	}
	count := mp.Count()
	for i, j := uint(0), uint(u.bitSize); i < j; i++ {
		mq := u.subMap(uint64(bsiOffsetBit + i))
		sum += (1 << i) * mq.IntersectionCount(mp)
	}
	return sum, count
}

func (u *ubsi) Get(k uint64) (interface{}, bool) {
	var v uint64

	if !u.bit(bsiExistsBit, k) {
		return 0, false
	}
	for i, j := uint(0), uint(u.bitSize); i < j; i++ {
		if u.bit(uint64(bsiOffsetBit+i), k) {
			v |= (1 << i)
		}
	}
	return v, true
}

func (u *ubsi) Set(k uint64, e interface{}) error {
	v := e.(uint64)
	for i, j := uint(0), uint(u.bitSize); i < j; i++ {
		if v&(1<<i) != 0 {
			u.setBit(uint64(bsiOffsetBit+i), k)
		} else {
			u.clearBit(uint64(bsiOffsetBit+i), k)
		}
	}
	u.setBit(bsiExistsBit, k)
	return nil
}

func (u *ubsi) Del(k uint64) error {
	u.clearBit(bsiExistsBit, k)
	return nil
}

func (u *ubsi) Eq(e interface{}) (*roaring.Bitmap, error) {
	v := e.(uint64)
	mp := u.subMap(bsiExistsBit)
	for i := u.bitSize - 1; i >= 0; i-- {
		if (v>>uint(i))&1 == 1 {
			mp = mp.Intersect(u.subMap(uint64(bsiOffsetBit + i)))
		} else {
			mp = mp.Difference(u.subMap(uint64(bsiOffsetBit + i)))
		}
	}
	return mp, nil
}

func (u *ubsi) Ne(e interface{}) (*roaring.Bitmap, error) {
	v := e.(uint64)
	mp := u.subMap(bsiExistsBit)
	for i := u.bitSize - 1; i >= 0; i-- {
		if (v>>uint(i))&1 == 1 {
			mp = mp.Intersect(u.subMap(uint64(bsiOffsetBit + i)))
		} else {
			mp = mp.Difference(u.subMap(uint64(bsiOffsetBit + i)))
		}
	}
	return u.subMap(bsiExistsBit).Difference(mp), nil
}

func (u *ubsi) Lt(e interface{}) (*roaring.Bitmap, error) {
	return u.lt(e.(uint64), u.subMap(bsiExistsBit), false), nil
}

func (u *ubsi) Le(e interface{}) (*roaring.Bitmap, error) {
	return u.lt(e.(uint64), u.subMap(bsiExistsBit), true), nil
}

func (u *ubsi) Gt(e interface{}) (*roaring.Bitmap, error) {
	return u.gt(e.(uint64), u.subMap(bsiExistsBit), false), nil
}

func (u *ubsi) Ge(e interface{}) (*roaring.Bitmap, error) {
	return u.gt(e.(uint64), u.subMap(bsiExistsBit), true), nil
}

func (u *ubsi) min(filter *roaring.Bitmap) (uint64, uint64) {
	var min uint64
	var count uint64

	for i := u.bitSize - 1; i >= 0; i-- {
		mp := filter.Difference(u.subMap(uint64(bsiOffsetBit + i)))
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

func (u *ubsi) max(filter *roaring.Bitmap) (uint64, uint64) {
	var max uint64
	var count uint64

	for i := u.bitSize - 1; i >= 0; i-- {
		mp := u.subMap(uint64(bsiOffsetBit + i)).Intersect(filter)
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

func (u *ubsi) lt(v uint64, mp *roaring.Bitmap, eq bool) *roaring.Bitmap {
	zflg := true // leading zero flag
	mq := roaring.NewBitmap()
	for i := u.bitSize - 1; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		if zflg {
			if bit == 0 {
				mp = mp.Difference(u.subMap(uint64(bsiOffsetBit + i)))
				continue
			} else {
				zflg = false
			}
		}
		if i == 0 && !eq {
			if bit == 0 {
				return mq
			}
			return mp.Difference(u.subMap(uint64(bsiOffsetBit + i)).Difference(mq))
		}
		if bit == 0 {
			mp = mp.Difference(u.subMap(uint64(bsiOffsetBit + i)).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Difference(u.subMap(uint64(bsiOffsetBit + i))))
		}
	}
	return mp
}

func (u *ubsi) gt(v uint64, mp *roaring.Bitmap, eq bool) *roaring.Bitmap {
	mq := roaring.NewBitmap()
	for i := u.bitSize - 1; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		if i == 0 && !eq {
			if bit == 1 {
				return mq
			}
			return mp.Difference(mp.Difference(u.subMap(uint64(bsiOffsetBit + i))).Difference(mq))
		}
		if bit == 1 {
			mp = mp.Difference(mp.Difference(u.subMap(uint64(bsiOffsetBit + i))).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Intersect(u.subMap(uint64(bsiOffsetBit + i))))
		}
	}
	return mp
}

// x is the bit offset, y is row id
func (u *ubsi) setBit(x, y uint64)   { u.ms[x].Add(y) }
func (u *ubsi) clearBit(x, y uint64) { u.ms[x].Remove(y) }
func (u *ubsi) bit(x, y uint64) bool { return u.ms[x].Contains(y) }

func (u *ubsi) subMap(x uint64) *roaring.Bitmap { return u.ms[x] }

func show(mp *roaring.Bitmap) ([]byte, error) {
	var buf bytes.Buffer

	if _, err := mp.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
