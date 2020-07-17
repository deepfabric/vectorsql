package dynamic

import (
	"bytes"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/deepfabric/vectorsql/pkg/vm/value/static"
	"github.com/pilosa/pilosa/roaring"
)

func NewStrings(vs []string, np, dp *roaring.Bitmap) *Strings {
	return &Strings{
		Vs: vs,
		Np: np,
		Dp: dp,
	}
}

func (a *Strings) Size() int {
	var n int

	for _, v := range a.Vs {
		n += len(v)
	}
	return n + len(a.Vs)*4
}

func (a *Strings) Show() ([]byte, error) {
	var buf bytes.Buffer

	data, err := show(a.Np)
	if err != nil {
		return nil, err
	}
	buf.Write(data)
	os := make([]uint32, len(a.Vs))
	for i, x := range a.Vs {
		os[i] = uint32(len(x))
	}
	{
		hp := *(*reflect.SliceHeader)(unsafe.Pointer(&os))
		hp.Len *= 4
		hp.Cap *= 4
		buf.Write(*(*[]byte)(unsafe.Pointer(&hp)))
	}
	for _, v := range a.Vs {
		buf.WriteString(v)
	}
	return buf.Bytes(), nil
}

func (a *Strings) Read(cnt int, data []byte) error {
	data, np, err := read(data)
	if err != nil {
		return err
	}
	a.Np = np
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&data))
	hp.Len = cnt
	hp.Cap = cnt
	a.Vs = make([]string, cnt)
	os := *(*[]uint32)(unsafe.Pointer(&hp))
	data = data[cnt*4:]
	for i := 0; i < cnt; i++ {
		a.Vs[i] = string(data[:os[i]])
		data = data[os[i]:]
	}
	return nil
}

func (a *Strings) MarkNull(row int) error {
	a.Np.DirectAdd(uint64(row))
	return nil
}

func (a *Strings) Append(v interface{}) error {
	a.Vs = append(a.Vs, v.([]string)...)
	return nil
}

func (a *Strings) Merge(np, dp *roaring.Bitmap) error {
	a.Dp = dp
	if a.Np == nil {
		a.Np = np
		return nil
	}
	a.Np = a.Np.Union(np)
	return nil
}

func (a *Strings) Update(rows []int, v interface{}) error {
	vs := v.([]string)
	for _, i := range rows {
		a.Vs[i] = vs[i]
	}
	return nil
}

func (a *Strings) Filter(is []uint64) interface{} {
	if len(is) == 0 {
		return &Strings{}
	}
	return &Strings{
		Is: is,
		Vs: a.Vs,
		Np: a.Np,
		Dp: a.Dp,
	}
}

func (a *Strings) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]string, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Strings) MergeFilter(v interface{}) interface{} {
	b := v.(*static.Bools)
	r := &Strings{
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
