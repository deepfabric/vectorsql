package static

import (
	"reflect"
	"unsafe"
)

func (a *Bools) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Floats) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Float32s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Float64s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Ints) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Int16s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Int32s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Int64s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Int8s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Timestamps) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Uint16s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Uint32s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Uint64s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Uint8s) Count() int {
	var cnt int

	if cnt = len(a.Is); cnt == 0 {
		cnt = len(a.Vs)
	}
	switch {
	case a.Dp == nil && a.Np == nil:
		return cnt
	case a.Dp != nil && a.Np == nil:
		return cnt + int(a.Np.Count())
	case a.Dp == nil && a.Np != nil:
		return cnt - int(a.Dp.Count())
	default: // a.Dp != nil && a.Np != nil
		return cnt + int(a.Np.Count()-a.Dp.Count())
	}
}

func (a *Bools) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+1])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+1])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Floats) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 8
	hp.Cap *= 8
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+8])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+8])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Float32s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 4
	hp.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+4])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+4])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Float64s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 8
	hp.Cap *= 8
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+8])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+8])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Ints) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 8
	hp.Cap *= 8
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+8])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+8])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Int8s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+1])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+1])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Int16s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 2
	hp.Cap *= 2
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+2])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+2])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Int32s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 4
	hp.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+4])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+4])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Int64s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 8
	hp.Cap *= 8
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+8])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+8])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Timestamps) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 8
	hp.Cap *= 8
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+8])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+8])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Uint8s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+1])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+1])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Uint16s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 2
	hp.Cap *= 2
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+2])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+2])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Uint32s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 4
	hp.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+4])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+4])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}

func (a *Uint64s) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&a.Vs))
	hp.Len *= 8
	hp.Cap *= 8
	data := *(*[]byte)(unsafe.Pointer(&hp))
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, data[o:o+8])
			} else {
				rs = append(rs, e)
			}
		}
		return os, rs
	}
	n := uint64(len(a.Vs))
	os := make([]uint64, 0, n)
	rs := make([][]byte, 0, n)
	for i := uint64(0); i < n; i++ {
		if a.Dp != nil && a.Dp.Contains(i) {
			continue
		}
		os = append(os, i)
		if a.Np == nil || !a.Np.Contains(i) {
			rs = append(rs, data[i:i+8])
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}
