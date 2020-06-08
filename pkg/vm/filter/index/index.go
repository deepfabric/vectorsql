package index

import (
	"bytes"

	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/vectorsql/pkg/vm/filter"
)

func New(fs []filter.Filter) *index {
	return &index{fs}
}

func (i *index) String() string {
	if len(i.fs) == 0 {
		return i.fs[0].String()
	}
	var buf bytes.Buffer
	for i, f := range i.fs {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString("(")
		buf.WriteString(f.String())
		buf.WriteString(")")
	}
	return buf.String()
}

func (i *index) Bitmap(mcpu int) (*roaring.Bitmap, error) {
	ms := make([]*roaring.Bitmap, 0, len(i.fs))
	for _, v := range i.fs {
		mp, err := v.Bitmap(mcpu)
		if err != nil {
			return nil, err
		}
		ms = append(ms, mp)
	}
	return roaring.ParAnd(mcpu, ms...), nil
}
