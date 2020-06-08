package ifilter

import (
	"bytes"
	"fmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/vectorsql/pkg/vm/container/relation"
)

func New(cs []*Condition, r relation.Relation) *filter {
	return &filter{cs, r}
}

func (f *filter) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("SELECT groupBitmapState(seq) FROM %s WHERE ", f.r))
	for i, c := range f.cs {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		switch c.Op {
		case EQ:
			buf.WriteString(fmt.Sprintf("%s = %s", c.Name, c.Val))
		case NE:
			buf.WriteString(fmt.Sprintf("%s <> %s", c.Name, c.Val))
		case LT:
			buf.WriteString(fmt.Sprintf("%s < %s", c.Name, c.Val))
		case LE:
			buf.WriteString(fmt.Sprintf("%s <= %s", c.Name, c.Val))
		case GT:
			buf.WriteString(fmt.Sprintf("%s > %s", c.Name, c.Val))
		case GE:
			buf.WriteString(fmt.Sprintf("%s >= %s", c.Name, c.Val))
		}
	}
	return buf.String()
}

func (f *filter) Bitmap(mcpu int) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, len(f.cs))
	for _, c := range f.cs {
		switch c.Op {
		case EQ:
			mp, err := f.r.Eq(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			ms = append(ms, mp)
		case NE:
			mp, err := f.r.Ne(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			ms = append(ms, mp)
		case LT:
			mp, err := f.r.Lt(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			ms = append(ms, mp)
		case LE:
			mp, err := f.r.Le(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			ms = append(ms, mp)
		case GT:
			mp, err := f.r.Gt(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			ms = append(ms, mp)
		case GE:
			mp, err := f.r.Ge(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			ms = append(ms, mp)
		}
	}
	if !f.r.IsEvent() {
		return roaring.ParAnd(mcpu, ms...), nil
	}
	mp, err := f.r.IdBitmap()
	if err != nil {
		return nil, err
	}
	is := roaring.ParAnd(mcpu, ms...).ToArray()
	var xs []uint32
	for _, i := range is {
		v, ok := mp.Get(i)
		if ok {
			xs = append(xs, v.(uint32))
		}
	}
	return roaring.BitmapOf(xs...), nil
}
