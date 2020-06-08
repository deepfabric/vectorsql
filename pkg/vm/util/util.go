package util

import (
	"fmt"

	"github.com/deepfabric/thinkbase/pkg/vm/value"
)

func Contain(xs, ys []string) error {
	mp := make(map[string]struct{})
	for _, y := range ys {
		mp[y] = struct{}{}
	}
	for _, x := range xs {
		if _, ok := mp[x]; !ok {
			return fmt.Errorf("'%s' not in '%v'", x, ys)
		}
	}
	return nil
}

func Indexs(xs, ys []string) []int {
	var rs []int

	mp := make(map[string]int)
	for i, y := range ys {
		mp[y] = i
	}
	for _, x := range xs {
		if i, ok := mp[x]; ok {
			rs = append(rs, i)
		}
	}
	return rs
}

func Tuple2Map(a value.Array, attrs []string) map[string]value.Value {
	mp := make(map[string]value.Value)
	for i, j := 0, len(a); i < j; i++ {
		mp[attrs[i]] = a[i]
	}
	return mp
}

func Map2Tuple(mp map[string]value.Array, attrs []string, index int) value.Array {
	var r value.Array

	for _, attr := range attrs {
		r = append(r, mp[attr][index])
	}
	return r
}

func Tuples2Map(ts value.Array, attrs []string) map[string]value.Array {
	mp := make(map[string]value.Array)
	for _, t := range ts {
		a := t.(value.Array)
		for i, v := range a {
			mp[attrs[i]] = append(mp[attrs[i]], v)
		}
	}
	return mp
}

func Map2Tuples(mp map[string]value.Array, attrs []string) value.Array {
	var r value.Array

	for i, j := 0, len(mp[attrs[0]]); i < j; i++ {
		r = append(r, Map2Tuple(mp, attrs, i))
	}
	return r
}

func SubTuple(a value.Array, is []int) value.Array {
	var r value.Array

	for _, i := range is {
		r = append(r, a[i])
	}
	return r
}

func SubMap(mp map[string]value.Array, attrs []string, index int) map[string]value.Value {
	rq := make(map[string]value.Value)
	for _, attr := range attrs {
		rq[attr] = mp[attr][index]
	}
	return rq
}

func MergeAttributes(xs, ys []string) []string {
	var rs []string

	mp := make(map[string]struct{})
	for i, j := 0, len(xs); i < j; i++ {
		if _, ok := mp[xs[i]]; !ok {
			mp[xs[i]] = struct{}{}
			rs = append(rs, xs[i])
		}
	}
	for i, j := 0, len(ys); i < j; i++ {
		if _, ok := mp[ys[i]]; !ok {
			mp[ys[i]] = struct{}{}
			rs = append(rs, ys[i])
		}
	}
	return rs
}
