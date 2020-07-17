package dynamic

func (a *Strings) Count() int {
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

func (a *Strings) Slice() ([]uint64, [][]byte) {
	e := []byte{}
	if n := len(a.Is); n > 0 {
		os := make([]uint64, 0, n)
		rs := make([][]byte, 0, n)
		for _, o := range a.Is {
			if a.Dp != nil && a.Dp.Contains(o) {
				continue
			}
			os = append(os, o)
			if a.Np == nil || !a.Np.Contains(o) {
				rs = append(rs, []byte(a.Vs[o]))
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
			rs = append(rs, []byte(a.Vs[i]))
		} else {
			rs = append(rs, e)
		}
	}
	return os, rs
}
