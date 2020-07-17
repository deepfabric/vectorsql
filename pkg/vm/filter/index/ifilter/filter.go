package ifilter

import (
	"bytes"
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/pilosa/pilosa/roaring"
)

func New(cs []*Condition, r storage.Relation) *filter {
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

func (f *filter) Bitmap() (*roaring.Bitmap, error) {
	var m *roaring.Bitmap

	for _, c := range f.cs {
		switch c.Op {
		case EQ:
			mp, err := f.r.Eq(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			if m == nil {
				m = mp
			} else {
				m = m.Intersect(mp)
			}
		case NE:
			mp, err := f.r.Ne(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			if m == nil {
				m = mp
			} else {
				m = m.Intersect(mp)
			}
		case LT:
			mp, err := f.r.Lt(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			if m == nil {
				m = mp
			} else {
				m = m.Intersect(mp)
			}
		case LE:
			mp, err := f.r.Le(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			if m == nil {
				m = mp
			} else {
				m = m.Intersect(mp)
			}
		case GT:
			mp, err := f.r.Gt(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			if m == nil {
				m = mp
			} else {
				m = m.Intersect(mp)
			}
		case GE:
			mp, err := f.r.Ge(c.Name, c.Val)
			if err != nil {
				return nil, err
			}
			if m == nil {
				m = mp
			} else {
				m = m.Intersect(mp)
			}
		}
	}
	if !f.r.IsEvent() {
		return m, nil
	}
	mp, err := f.r.IdMap()
	if err != nil {
		return nil, err
	}
	is := m.Slice()
	var xs []uint64
	for _, i := range is {
		if v, ok := mp.Get(i); ok {
			xs = append(xs, v.(uint64))
		}
	}
	return roaring.NewBitmap(xs...), nil
}
