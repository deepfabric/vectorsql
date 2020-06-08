package unsigned

import (
	"bytes"

	Roaring "github.com/RoaringBitmap/roaring"
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

func (u *ubsi) Map() *Roaring.Bitmap {
	return convert(u.subMap(bsiExistsBit))
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

func (u *ubsi) Get(k uint32) (interface{}, bool) {
	var v uint64

	if !u.bit(bsiExistsBit, k) {
		return 0, false
	}
	for i, j := uint(0), uint(u.bitSize); i < j; i++ {
		if u.bit(uint32(bsiOffsetBit+i), k) {
			v |= (1 << i)
		}
	}
	return v, true
}

func (u *ubsi) Set(k uint32, e interface{}) error {
	v := e.(uint64)
	for i, j := uint(0), uint(u.bitSize); i < j; i++ {
		if v&(1<<i) != 0 {
			u.setBit(uint32(bsiOffsetBit+i), k)
		} else {
			u.clearBit(uint32(bsiOffsetBit+i), k)
		}
	}
	u.setBit(uint32(bsiExistsBit), k)
	return nil
}

func (u *ubsi) Del(k uint32) error {
	u.clearBit(uint32(bsiExistsBit), k)
	return nil
}

func (u *ubsi) Eq(e interface{}) (*Roaring.Bitmap, error) {
	v := e.(uint64)
	mp := u.subMap(bsiExistsBit)
	for i := u.bitSize - 1; i >= 0; i-- {
		if (v>>uint(i))&1 == 1 {
			mp = mp.Intersect(u.subMap(uint32(bsiOffsetBit + i)))
		} else {
			mp = mp.Difference(u.subMap(uint32(bsiOffsetBit + i)))
		}
	}
	return convert(mp), nil
}

func (u *ubsi) Ne(e interface{}) (*Roaring.Bitmap, error) {
	v := e.(uint64)
	mp := u.subMap(bsiExistsBit)
	for i := u.bitSize - 1; i >= 0; i-- {
		if (v>>uint(i))&1 == 1 {
			mp = mp.Intersect(u.subMap(uint32(bsiOffsetBit + i)))
		} else {
			mp = mp.Difference(u.subMap(uint32(bsiOffsetBit + i)))
		}
	}
	return convert(u.subMap(bsiExistsBit).Difference(mp)), nil
}

func (u *ubsi) Lt(e interface{}) (*Roaring.Bitmap, error) {
	return convert(u.lt(e.(uint64), u.subMap(bsiExistsBit), false)), nil
}

func (u *ubsi) Le(e interface{}) (*Roaring.Bitmap, error) {
	return convert(u.lt(e.(uint64), u.subMap(bsiExistsBit), true)), nil
}

func (u *ubsi) Gt(e interface{}) (*Roaring.Bitmap, error) {
	return convert(u.gt(e.(uint64), u.subMap(bsiExistsBit), false)), nil
}

func (u *ubsi) Ge(e interface{}) (*Roaring.Bitmap, error) {
	return convert(u.gt(e.(uint64), u.subMap(bsiExistsBit), true)), nil
}

func (u *ubsi) lt(v uint64, mp *roaring.Bitmap, eq bool) *roaring.Bitmap {
	zflg := true // leading zero flag
	mq := roaring.NewBitmap()
	for i := u.bitSize - 1; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		if zflg {
			if bit == 0 {
				mp = mp.Difference(u.subMap(uint32(bsiOffsetBit + i)))
				continue
			} else {
				zflg = false
			}
		}
		if i == 0 && !eq {
			if bit == 0 {
				return mq
			}
			return mp.Difference(u.subMap(uint32(bsiOffsetBit + i)).Difference(mq))
		}
		if bit == 0 {
			mp = mp.Difference(u.subMap(uint32(bsiOffsetBit + i)).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Difference(u.subMap(uint32(bsiOffsetBit + i))))
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
			return mp.Difference(mp.Difference(u.subMap(uint32(bsiOffsetBit + i))).Difference(mq))
		}
		if bit == 1 {
			mp = mp.Difference(mp.Difference(u.subMap(uint32(bsiOffsetBit + i))).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Intersect(u.subMap(uint32(bsiOffsetBit + i))))
		}
	}
	return mp
}

// x is the bit offset, y is row id
func (u *ubsi) setBit(x, y uint32)   { u.ms[x].Add(uint64(y)) }
func (u *ubsi) clearBit(x, y uint32) { u.ms[x].Remove(uint64(y)) }
func (u *ubsi) bit(x, y uint32) bool { return u.ms[x].Contains(uint64(y)) }

func (u *ubsi) subMap(x uint32) *roaring.Bitmap { return u.ms[x] }

func show(mp *roaring.Bitmap) ([]byte, error) {
	var buf bytes.Buffer

	if _, err := mp.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func convert(mp *roaring.Bitmap) *Roaring.Bitmap {
	var xs []uint32

	{
		itr := mp.Iterator()
		itr.Seek(0)
		for v, eof := itr.Next(); !eof; v, eof = itr.Next() {
			xs = append(xs, uint32(v))
		}
	}
	return Roaring.BitmapOf(xs...)
}
