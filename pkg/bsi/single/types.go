package single

import "github.com/pilosa/pilosa/roaring"

type bsi struct {
	ms []*roaring.Bitmap
}

const (
	bsiExistsBit  = 0
	bsiSignBit    = 1
	bsiOffsetBit  = 2
	bsiFoffsetBit = 10
)
