package single

import (
	"bytes"
	"fmt"
	"math"

	Roaring "github.com/RoaringBitmap/roaring"
	"github.com/pilosa/pilosa/roaring"
)

func New() *bsi {
	ms := make([]*roaring.Bitmap, 34)
	for i := 0; i < 34; i++ {
		ms[i] = roaring.NewBitmap()
	}
	return &bsi{
		ms: ms,
	}
}

func (b *bsi) Map() *Roaring.Bitmap {
	return convert(b.subMap(bsiExistsBit))
}

func (b *bsi) Clone() *bsi {
	ms := make([]*roaring.Bitmap, 34)
	for i := 0; i < 34; i++ {
		ms[i] = b.ms[i].Clone()
	}
	return &bsi{
		ms: ms,
	}
}

/*
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

func (b *bsi) Get(k uint32) (interface{}, bool) {
	var v int64

	if !b.bit(bsiExistsBit, k) {
		return -1, false
	}
	for i, j := uint(0), uint(b.bitSize); i < j; i++ {
		if b.bit(uint32(bsiOffsetBit+i), k) {
			v |= (1 << i)
		}
	}
	if b.bit(uint32(bsiSignBit), k) {
		v = -v
	}
	return v, true
}
*/

func (b *bsi) Set(k uint32, x interface{}) error {
	v := x.(float32)
	uv := math.Float32bits(v)
	s := (uv >> 31)
	e := (uv>>23)&0xFF - 127
	m := uint64(uv&0x7FFFFF) + 1<<23

	{
		fmt.Printf("\te: %v, m: %x\n", e, m)
	}
	for i := uint(0); i < 8; i++ {
		if e&(1<<i) != 0 {
			b.setBit(uint32(bsiOffsetBit+i), k)
		} else {
			b.clearBit(uint32(bsiOffsetBit+i), k)
		}
	}
	for i := uint(0); i < 24; i++ {
		if m&(1<<i) != 0 {
			b.setBit(uint32(bsiFoffsetBit+i), k)
		} else {
			b.clearBit(uint32(bsiFoffsetBit+i), k)
		}
	}
	b.setBit(uint32(bsiExistsBit), k)
	if s != 0 {
		b.setBit(uint32(bsiSignBit), k)
	} else {
		b.clearBit(uint32(bsiSignBit), k)
	}
	return nil
}

func (b *bsi) Del(k uint32) error {
	b.clearBit(uint32(bsiExistsBit), k)
	return nil
}

/*
func (b *bsi) Eq(e interface{}) (*Roaring.Bitmap, error) {
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
			mp = mp.Intersect(b.subMap(uint32(bsiOffsetBit + i)))
		} else {
			mp = mp.Difference(b.subMap(uint32(bsiOffsetBit + i)))
		}
	}
	return convert(mp), nil
}

func (b *bsi) Ne(e interface{}) (*Roaring.Bitmap, error) {
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
			mp = mp.Intersect(b.subMap(uint32(bsiOffsetBit + i)))
		} else {
			mp = mp.Difference(b.subMap(uint32(bsiOffsetBit + i)))
		}
	}
	return convert(b.subMap(bsiExistsBit).Difference(mp)), nil
}
*/

func (b *bsi) Lt(x interface{}) (*Roaring.Bitmap, error) {
	v := x.(float32)
	uv := math.Float32bits(v)
	//s := (uv >> 31)
	e := (uv>>23)&0xFF - 127
	m := uint64(uv&0x7FFFFF) + 1<<23
	{
		fmt.Printf("\tlt -> e = %v, m = %x\n", e, m)
	}

	m0 := b.eeq(uint64(e), b.subMap(bsiExistsBit))
	m1 := b.lt(m, b.subMap(bsiExistsBit), false)
	m2 := b.elt(uint64(e), b.subMap(bsiExistsBit))
	return convert(m2.Union(m0.Intersect(m1))), nil
	/*
		v := e.(int64)
		uv := uint64(v)
		if v < 0 {
			uv = uint64(-v)
		}
		mp := b.subMap(bsiExistsBit)
		if v >= 0 {
			mp = b.lt(uv, mp.Difference(b.subMap(bsiSignBit)), false)
			return convert(b.subMap(bsiSignBit).Union(mp)), nil
		}
		return convert(b.gt(uv, mp.Intersect(b.subMap(bsiSignBit)), false)), nil
	*/
}

/*
func (b *bsi) Le(e interface{}) (*Roaring.Bitmap, error) {
	v := e.(int64)
	uv := uint64(v)
	if v < 0 {
		uv = uint64(-v)
	}
	mp := b.subMap(bsiExistsBit)
	if v >= 0 {
		mp = b.lt(uv, mp.Difference(b.subMap(bsiSignBit)), true)
		return convert(b.subMap(bsiSignBit).Union(mp)), nil
	}
	return convert(b.gt(uv, mp.Intersect(b.subMap(bsiSignBit)), true)), nil
}
*/

func (b *bsi) Gt(x interface{}) (*Roaring.Bitmap, error) {
	v := x.(float32)
	uv := math.Float32bits(v)
	//s := (uv >> 31)
	e := (uv>>23)&0xFF - 127
	m := uint64(uv&0x7FFFFF) + 1<<23

	m0 := b.eeq(uint64(e), b.subMap(bsiExistsBit))
	m1 := b.gt(m, b.subMap(bsiExistsBit), false)
	m2 := b.egt(uint64(e), b.subMap(bsiExistsBit))
	return convert(m2.Union(m0.Intersect(m1))), nil
	/*
		v := e.(int64)
		uv := uint64(v)
		if v < 0 {
			uv = uint64(-v)
		}
		mp := b.subMap(bsiExistsBit)
		if v >= 0 {
			return convert(b.gt(uv, mp.Difference(b.subMap(uint32(bsiSignBit))), false)), nil
		}
		mq := b.lt(uv, mp.Intersect(b.subMap(uint32(bsiSignBit))), false)
		return convert(mp.Difference(b.subMap(uint32(bsiSignBit))).Union(mq)), nil
	*/
}

/*
func (b *bsi) Ge(e interface{}) (*Roaring.Bitmap, error) {
	v := e.(int64)
	uv := uint64(v)
	if v < 0 {
		uv = uint64(-v)
	}
	mp := b.subMap(bsiExistsBit)
	if v >= 0 {
		return convert(b.gt(uv, mp.Difference(b.subMap(uint32(bsiSignBit))), true)), nil
	}
	mq := b.lt(uv, mp.Intersect(b.subMap(uint32(bsiSignBit))), true)
	return convert(mp.Difference(b.subMap(uint32(bsiSignBit))).Union(mq)), nil
}
*/

func (b *bsi) eeq(v uint64, mp *roaring.Bitmap) *roaring.Bitmap {
	for i := 7; i >= 0; i-- {
		if (v>>uint(i))&1 == 1 {
			mp = mp.Intersect(b.subMap(uint32(bsiOffsetBit + i)))
		} else {
			mp = mp.Difference(b.subMap(uint32(bsiOffsetBit + i)))
		}
	}
	return mp
}

func (b *bsi) elt(v uint64, mp *roaring.Bitmap) *roaring.Bitmap {
	zflg := true // leading zero flag
	mq := roaring.NewBitmap()
	for i := 7; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		if zflg {
			if bit == 0 {
				mp = mp.Difference(b.subMap(uint32(bsiOffsetBit + i)))
				continue
			} else {
				zflg = false
			}
		}
		if i == 0 {
			if bit == 0 {
				return mq
			}
			return mp.Difference(b.subMap(uint32(bsiFoffsetBit + i)).Difference(mq))
		}
		if bit == 0 {
			mp = mp.Difference(b.subMap(uint32(bsiOffsetBit + i)).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Difference(b.subMap(uint32(bsiOffsetBit + i))))
		}
	}
	return mp
}

func (b *bsi) egt(v uint64, mp *roaring.Bitmap) *roaring.Bitmap {
	mq := roaring.NewBitmap()
	for i := 7; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		if i == 0 {
			if bit == 1 {
				return mq
			}
			return mp.Difference(mp.Difference(b.subMap(uint32(bsiFoffsetBit + i))).Difference(mq))
		}
		if bit == 1 {
			mp = mp.Difference(mp.Difference(b.subMap(uint32(bsiOffsetBit + i))).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Intersect(b.subMap(uint32(bsiOffsetBit + i))))
		}
	}
	return mp
}

func (b *bsi) lt(v uint64, mp *roaring.Bitmap, eq bool) *roaring.Bitmap {
	zflg := true // leading zero flag
	mq := roaring.NewBitmap()
	for i := 23; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		if zflg {
			if bit == 0 {
				mp = mp.Difference(b.subMap(uint32(bsiFoffsetBit + i)))
				continue
			} else {
				zflg = false
			}
		}
		if i == 0 && !eq {
			if bit == 0 {
				return mq
			}
			return mp.Difference(b.subMap(uint32(bsiFoffsetBit + i)).Difference(mq))
		}
		if bit == 0 {
			mp = mp.Difference(b.subMap(uint32(bsiFoffsetBit + i)).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Difference(b.subMap(uint32(bsiFoffsetBit + i))))
		}
	}
	return mp
}

func (b *bsi) gt(v uint64, mp *roaring.Bitmap, eq bool) *roaring.Bitmap {
	mq := roaring.NewBitmap()
	for i := 23; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		if i == 0 && !eq {
			if bit == 1 {
				return mq
			}
			return mp.Difference(mp.Difference(b.subMap(uint32(bsiFoffsetBit + i))).Difference(mq))
		}
		if bit == 1 {
			mp = mp.Difference(mp.Difference(b.subMap(uint32(bsiFoffsetBit + i))).Difference(mq))
			continue
		}
		if i > 0 {
			mq = mq.Union(mp.Intersect(b.subMap(uint32(bsiFoffsetBit + i))))
		}
	}
	return mp
}

// x is the bit offset, y is row id
func (b *bsi) setBit(x, y uint32)   { b.ms[x].Add(uint64(y)) }
func (b *bsi) clearBit(x, y uint32) { b.ms[x].Remove(uint64(y)) }
func (b *bsi) bit(x, y uint32) bool { return b.ms[x].Contains(uint64(y)) }

func (b *bsi) subMap(x uint32) *roaring.Bitmap { return b.ms[x] }

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
