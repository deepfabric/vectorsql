package static

import (
	"reflect"
	"unsafe"

	"github.com/pilosa/pilosa/roaring"
)

func NewUint8s(vs []uint8, np, dp *roaring.Bitmap) *Uint8s {
	return &Uint8s{
		Vs: vs,
		Np: np,
		Dp: dp,
	}
}

func (a *Uint8s) Size() int {
	return len(a.Vs)
}

func (a *Uint8s) Show() ([]byte, error) {
	v, err := show(a.Np)
	if err != nil {
		return nil, err
	}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	return append(v, *(*[]byte)(unsafe.Pointer(&hp))...), nil
}

func (a *Uint8s) Read(cnt int, data []byte) error {
	data, np, err := read(data)
	if err != nil {
		return err
	}
	a.Np = np
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&data))
	hp.Len = cnt
	hp.Cap = cnt
	a.Vs = *(*[]uint8)(unsafe.Pointer(&hp))
	return nil
}

func (a *Uint8s) MarkNull(row int) error {
	a.Np.DirectAdd(uint64(row))
	return nil
}

func (a *Uint8s) Append(v interface{}) error {
	a.Vs = append(a.Vs, v.([]uint8)...)
	return nil
}

func (a *Uint8s) Merge(np, dp *roaring.Bitmap) error {
	a.Dp = dp
	if a.Np == nil {
		a.Np = np
		return nil
	}
	a.Np = a.Np.Union(np)
	return nil
}

func (a *Uint8s) Update(rows []int, v interface{}) error {
	vs := v.([]uint8)
	for _, i := range rows {
		a.Vs[i] = vs[i]
	}
	return nil
}

func (a *Uint8s) Filter(is []uint64) interface{} {
	if len(is) == 0 {
		return &Uint8s{}
	}
	return &Uint8s{
		Is: is,
		Vs: a.Vs,
		Np: a.Np,
		Dp: a.Dp,
	}
}

func (a *Uint8s) MergeFilter(v interface{}) interface{} {
	b := v.(*Bools)
	r := &Uint8s{
		Vs: a.Vs,
	}
	switch {
	case a.Np != nil && b.Np == nil:
		r.Np = a.Np
	case a.Np == nil && b.Np != nil:
		r.Np = b.Np
	case a.Np != nil && b.Np != nil:
		r.Np = a.Np.Union(b.Np)
	}
	switch {
	case a.Dp != nil && b.Dp == nil:
		r.Dp = a.Dp
	case a.Dp == nil && b.Dp != nil:
		r.Dp = b.Dp
	case a.Dp != nil && b.Dp != nil:
		r.Dp = a.Np.Union(b.Dp)
	}
	switch {
	case len(a.Is) > 0 && len(b.Is) > 0:
		mp := make(map[uint64]struct{})
		{
			for _, o := range a.Is {
				mp[o] = struct{}{}
			}
		}
		r.Is = make([]uint64, 0, len(b.Is))
		for _, o := range b.Is {
			if _, ok := mp[o]; ok && b.Vs[o] {
				r.Is = append(r.Is, o)
			}
		}
	case len(a.Is) > 0 && len(b.Is) == 0:
		r.Is = make([]uint64, 0, len(a.Is))
		for _, o := range a.Is {
			if b.Vs[o] {
				r.Is = append(r.Is, o)
			}
		}
	case len(a.Is) == 0 && len(b.Is) > 0:
		r.Is = make([]uint64, 0, len(b.Is))
		for _, o := range b.Is {
			if b.Vs[o] {
				r.Is = append(r.Is, o)
			}
		}
	case len(a.Is) == 0 && len(b.Is) == 0:
		r.Is = make([]uint64, 0, len(a.Vs))
		for i := range a.Vs {
			if b.Vs[i] {
				r.Is = append(r.Is, uint64(i))
			}
		}
	}
	return r
}
