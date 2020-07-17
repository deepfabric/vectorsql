package static

import (
	"bytes"

	"github.com/deepfabric/vectorsql/pkg/vm/util/encoding"
	"github.com/pilosa/pilosa/roaring"
)

func show(np *roaring.Bitmap) ([]byte, error) {
	var buf bytes.Buffer

	if np == nil {
		buf.Write(encoding.EncodeUint32(0))
	} else {
		buf.Write(encoding.EncodeUint32(0))
		n, err := np.WriteTo(&buf)
		if err != nil {
			return nil, err
		}
		copy(buf.Bytes(), encoding.EncodeUint32(uint32(n)))
	}
	return buf.Bytes(), nil
}

func read(v []byte) ([]byte, *roaring.Bitmap, error) {
	var np *roaring.Bitmap

	{
		n := encoding.DecodeUint32(v[:4])
		v = v[4:]
		if n != 0 {
			np = roaring.NewBitmap()
			if err := np.UnmarshalBinary(v[:n]); err != nil {
				return nil, nil, err
			}
			v = v[n:]
		}
	}
	return v, np, nil
}
