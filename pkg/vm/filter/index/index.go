package index

import (
	"bytes"

	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/vectorsql/pkg/vm/filter/index/ifilter"
	Roaring "github.com/pilosa/pilosa/roaring"
)

func New(fs []ifilter.Filter) *index {
	return &index{fs}
}

func (r *index) String() string {
	if len(r.fs) == 0 {
		return r.fs[0].String()
	}
	var buf bytes.Buffer
	for i, f := range r.fs {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString("(")
		buf.WriteString(f.String())
		buf.WriteString(")")
	}
	return buf.String()
}

func (r *index) Bitmap() (*roaring.Bitmap, error) {
	var m *Roaring.Bitmap

	for _, v := range r.fs {
		mp, err := v.Bitmap()
		if err != nil {
			return nil, err
		}
		if m == nil {
			m = mp
		} else {
			m = m.Intersect(mp)
		}
	}
	return convert(m), nil
}

func convert(mp *Roaring.Bitmap) *roaring.Bitmap {
	var xs []uint32

	{
		itr := mp.Iterator()
		itr.Seek(0)
		for v, eof := itr.Next(); !eof; v, eof = itr.Next() {
			xs = append(xs, uint32(v))
		}
	}
	return roaring.BitmapOf(xs...)
}
